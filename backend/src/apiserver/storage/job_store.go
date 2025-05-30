// Copyright 2018 The Kubeflow Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package storage

import (
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/golang/glog"
	"github.com/kubeflow/pipelines/backend/src/apiserver/common"
	"github.com/kubeflow/pipelines/backend/src/apiserver/list"
	"github.com/kubeflow/pipelines/backend/src/apiserver/model"
	"github.com/kubeflow/pipelines/backend/src/common/util"
)

var jobColumns = []string{
	"UUID",
	"DisplayName",
	"Name",
	"Namespace",
	"ServiceAccount",
	"Description",
	"MaxConcurrency",
	"NoCatchup",
	"CreatedAtInSec",
	"UpdatedAtInSec",
	"Enabled",
	"CronScheduleStartTimeInSec",
	"CronScheduleEndTimeInSec",
	"Schedule",
	"PeriodicScheduleStartTimeInSec",
	"PeriodicScheduleEndTimeInSec",
	"IntervalSecond",
	"PipelineId",
	"PipelineName",
	"PipelineSpecManifest",
	"WorkflowSpecManifest",
	"Parameters",
	"Conditions",
	"RuntimeParameters",
	"PipelineRoot",
	"ExperimentUUID",
	"PipelineVersionId",
}

type JobStoreInterface interface {
	// Create a recurring run entry in the database.
	CreateJob(*model.Job) (*model.Job, error)

	// Fetches a recurring run with a given id from the database.
	GetJob(id string) (*model.Job, error)

	// Fetches recurring runs from the database with the specified filtering and listing options.
	ListJobs(filterContext *model.FilterContext, opts *list.Options) ([]*model.Job, int, string, error)

	// Enable or disables a recurring run in the database.
	ChangeJobMode(id string, enabled bool) error

	// Update a recurring run entry in the database.
	UpdateJob(swf *util.ScheduledWorkflow) error

	// Removes a recurring run entry from the database.
	DeleteJob(id string) error
}

type JobStore struct {
	db                     *DB
	resourceReferenceStore *ResourceReferenceStore
	time                   util.TimeInterface
}

// Runs two SQL queries in a transaction to return a list of matching jobs, as well as their
// total_size. The total_size does not reflect the page size, but it does reflect the number of jobs
// matching the supplied filters and resource references.
func (s *JobStore) ListJobs(
	filterContext *model.FilterContext, opts *list.Options,
) ([]*model.Job, int, string, error) {
	errorF := func(err error) ([]*model.Job, int, string, error) {
		return nil, 0, "", util.NewInternalServerError(err, "Failed to list jobs: %v", err)
	}

	rowsSql, rowsArgs, err := s.buildSelectJobsQuery(false, opts, filterContext)
	if err != nil {
		return errorF(err)
	}

	sizeSql, sizeArgs, err := s.buildSelectJobsQuery(true, opts, filterContext)
	if err != nil {
		return errorF(err)
	}

	// Use a transaction to make sure we're returning the total_size of the same rows queried
	tx, err := s.db.Begin()
	if err != nil {
		glog.Errorf("Failed to start transaction to list jobs")
		return errorF(err)
	}

	rows, err := tx.Query(rowsSql, rowsArgs...)
	if err != nil {
		return errorF(err)
	}
	if err := rows.Err(); err != nil {
		return errorF(err)
	}
	jobs, err := s.scanRows(rows)
	if err != nil {
		tx.Rollback()
		return errorF(err)
	}
	defer rows.Close()

	sizeRow, err := tx.Query(sizeSql, sizeArgs...)
	if err != nil {
		tx.Rollback()
		return errorF(err)
	}
	if err := sizeRow.Err(); err != nil {
		tx.Rollback()
		return errorF(err)
	}
	totalSize, err := list.ScanRowToTotalSize(sizeRow)
	if err != nil {
		tx.Rollback()
		return errorF(err)
	}
	defer sizeRow.Close()

	err = tx.Commit()
	if err != nil {
		glog.Errorf("Failed to commit transaction to list jobs")
		return errorF(err)
	}

	if len(jobs) <= opts.PageSize {
		return jobs, totalSize, "", nil
	}

	npt, err := opts.NextPageToken(jobs[opts.PageSize])
	return jobs[:opts.PageSize], totalSize, npt, err
}

func (s *JobStore) buildSelectJobsQuery(selectCount bool, opts *list.Options,
	filterContext *model.FilterContext,
) (string, []interface{}, error) {
	var filteredSelectBuilder sq.SelectBuilder
	var err error

	refKey := filterContext.ReferenceKey
	if refKey != nil && refKey.Type == model.ExperimentResourceType && (refKey.ID != "" || common.IsMultiUserMode()) {
		filteredSelectBuilder, err = list.FilterOnExperiment("jobs", jobColumns,
			selectCount, refKey.ID)
	} else if refKey != nil && refKey.Type == model.NamespaceResourceType && (refKey.ID != "" || common.IsMultiUserMode()) {
		filteredSelectBuilder, err = list.FilterOnNamespace("jobs", jobColumns,
			selectCount, refKey.ID)
	} else {
		filteredSelectBuilder, err = list.FilterOnResourceReference("jobs", jobColumns,
			model.JobResourceType, selectCount, filterContext)
	}
	if err != nil {
		return "", nil, util.NewInternalServerError(err, "Failed to list jobs: %v", err)
	}
	sqlBuilder := opts.AddFilterToSelect(filteredSelectBuilder)

	// If we're not just counting, then also add select columns and perform a left join
	// to get resource reference information. Also add pagination.
	if !selectCount {
		sqlBuilder = opts.AddPaginationToSelect(sqlBuilder)
		sqlBuilder = s.addResourceReferences(sqlBuilder)
		sqlBuilder = opts.AddSortingToSelect(sqlBuilder)
	}
	sql, args, err := sqlBuilder.ToSql()
	if err != nil {
		return "", nil, util.NewInternalServerError(err, "Failed to list jobs: %v", err)
	}

	return sql, args, err
}

func (s *JobStore) GetJob(id string) (*model.Job, error) {
	sql, args, err := s.addResourceReferences(sq.Select(jobColumns...).From("jobs")).
		Where(sq.Eq{"uuid": id}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, util.NewInternalServerError(err, "Failed to create query to get job: %v",
			err.Error())
	}
	row, err := s.db.Query(sql, args...)
	if err != nil {
		return nil, util.NewInternalServerError(err, "Failed to get job: %v",
			err.Error())
	}
	defer row.Close()
	jobs, err := s.scanRows(row)
	if err != nil || len(jobs) > 1 {
		return nil, util.NewInternalServerError(err, "Failed to get job: %v", err.Error())
	}
	if len(jobs) == 0 {
		return nil, util.NewResourceNotFoundError("Job", fmt.Sprint(id))
	}
	return jobs[0], nil
}

func (s *JobStore) addResourceReferences(filteredSelectBuilder sq.SelectBuilder) sq.SelectBuilder {
	resourceRefConcatQuery := s.db.Concat([]string{`"["`, s.db.GroupConcat("r.Payload", ","), `"]"`}, "")
	return sq.
		Select("jobs.*", resourceRefConcatQuery+" AS refs").
		FromSelect(filteredSelectBuilder, "jobs").
		// Append all the resource references for the run as a json column
		LeftJoin("(select * from resource_references where ResourceType='Job') AS r ON jobs.UUID=r.ResourceUUID").
		GroupBy("jobs.UUID")
}

func (s *JobStore) scanRows(r *sql.Rows) ([]*model.Job, error) {
	var jobs []*model.Job
	for r.Next() {
		var uuid, displayName, name, namespace, pipelineId, pipelineName, conditions, serviceAccount,
			description, parameters, pipelineSpecManifest, workflowSpecManifest string
		var cronScheduleStartTimeInSec, cronScheduleEndTimeInSec, createdAtInSec,
			periodicScheduleStartTimeInSec, periodicScheduleEndTimeInSec, intervalSecond, updatedAtInSec sql.NullInt64
		var cron, resourceReferencesInString, runtimeParameters, pipelineRoot sql.NullString
		var experimentId, pipelineVersionId sql.NullString
		var enabled, noCatchup bool
		var maxConcurrency int64
		err := r.Scan(
			&uuid, &displayName, &name, &namespace, &serviceAccount, &description,
			&maxConcurrency, &noCatchup, &createdAtInSec, &updatedAtInSec, &enabled,
			&cronScheduleStartTimeInSec, &cronScheduleEndTimeInSec, &cron,
			&periodicScheduleStartTimeInSec, &periodicScheduleEndTimeInSec, &intervalSecond,
			&pipelineId, &pipelineName, &pipelineSpecManifest, &workflowSpecManifest, &parameters,
			&conditions, &runtimeParameters, &pipelineRoot, &experimentId,
			&pipelineVersionId, &resourceReferencesInString)
		if err != nil {
			return nil, err
		}
		resourceReferences, _ := parseResourceReferences(resourceReferencesInString)
		expId := experimentId.String
		pvId := pipelineVersionId.String
		if len(resourceReferences) > 0 {
			if expId == "" {
				expId = model.GetRefIdFromResourceReferences(resourceReferences, model.ExperimentResourceType)
			}
			if namespace == "" {
				namespace = model.GetRefIdFromResourceReferences(resourceReferences, model.NamespaceResourceType)
			}
			if pipelineId == "" {
				pipelineId = model.GetRefIdFromResourceReferences(resourceReferences, model.PipelineResourceType)
			}
			if pvId == "" {
				pvId = model.GetRefIdFromResourceReferences(resourceReferences, model.PipelineVersionResourceType)
			}
		}
		runtimeConfig := parseRuntimeConfig(runtimeParameters, pipelineRoot)
		job := &model.Job{
			UUID:           uuid,
			DisplayName:    displayName,
			K8SName:        name,
			Namespace:      namespace,
			ServiceAccount: serviceAccount,
			Description:    description,
			Enabled:        enabled,
			Conditions:     conditions,
			ExperimentId:   expId,
			MaxConcurrency: maxConcurrency,
			NoCatchup:      noCatchup,
			// ResourceReferences: resourceReferences,
			Trigger: model.Trigger{
				CronSchedule: model.CronSchedule{
					CronScheduleStartTimeInSec: NullInt64ToPointer(cronScheduleStartTimeInSec),
					CronScheduleEndTimeInSec:   NullInt64ToPointer(cronScheduleEndTimeInSec),
					Cron:                       NullStringToPointer(cron),
				},
				PeriodicSchedule: model.PeriodicSchedule{
					PeriodicScheduleStartTimeInSec: NullInt64ToPointer(periodicScheduleStartTimeInSec),
					PeriodicScheduleEndTimeInSec:   NullInt64ToPointer(periodicScheduleEndTimeInSec),
					IntervalSecond:                 NullInt64ToPointer(intervalSecond),
				},
			},
			PipelineSpec: model.PipelineSpec{
				PipelineId:           pipelineId,
				PipelineVersionId:    pvId,
				PipelineName:         pipelineName,
				PipelineSpecManifest: pipelineSpecManifest,
				WorkflowSpecManifest: workflowSpecManifest,
				Parameters:           parameters,
				RuntimeConfig:        runtimeConfig,
			},
			CreatedAtInSec: createdAtInSec.Int64,
			UpdatedAtInSec: updatedAtInSec.Int64,
		}
		job = job.ToV2()
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (s *JobStore) DeleteJob(id string) error {
	jobSql, jobArgs, err := sq.Delete("jobs").Where(sq.Eq{"UUID": id}).ToSql()
	if err != nil {
		return util.NewInternalServerError(err,
			"Failed to create query to delete job: %s", id)
	}
	// Use a transaction to make sure both run and its resource references are stored.
	tx, err := s.db.Begin()
	if err != nil {
		return util.NewInternalServerError(err, "Failed to create a new transaction to delete job")
	}
	_, err = tx.Exec(jobSql, jobArgs...)
	if err != nil {
		tx.Rollback()
		return util.NewInternalServerError(err, "Failed to delete job %s from table", id)
	}
	err = s.resourceReferenceStore.DeleteResourceReferences(tx, id, model.JobResourceType)
	if err != nil {
		tx.Rollback()
		return util.NewInternalServerError(err, "Failed to delete resource references from table for job %v ", id)
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return util.NewInternalServerError(err, "Failed to delete job %v and its resource references from table", id)
	}
	return nil
}

func (s *JobStore) CreateJob(j *model.Job) (*model.Job, error) {
	// Add creation/update time.
	j = j.ToV1().ToV2()
	now := s.time.Now().Unix()
	j.CreatedAtInSec = now
	j.UpdatedAtInSec = now

	jobSql, jobArgs, err := sq.
		Insert("jobs").
		SetMap(sq.Eq{
			"UUID":                           j.UUID,
			"DisplayName":                    j.DisplayName,
			"Name":                           j.K8SName,
			"Namespace":                      j.Namespace,
			"ServiceAccount":                 j.ServiceAccount,
			"Description":                    j.Description,
			"MaxConcurrency":                 j.MaxConcurrency,
			"NoCatchup":                      j.NoCatchup,
			"Enabled":                        j.Enabled,
			"Conditions":                     j.Conditions,
			"CronScheduleStartTimeInSec":     PointerToNullInt64(j.Trigger.CronSchedule.CronScheduleStartTimeInSec),
			"CronScheduleEndTimeInSec":       PointerToNullInt64(j.Trigger.CronSchedule.CronScheduleEndTimeInSec),
			"Schedule":                       PointerToNullString(j.Trigger.CronSchedule.Cron),
			"PeriodicScheduleStartTimeInSec": PointerToNullInt64(j.Trigger.PeriodicSchedule.PeriodicScheduleStartTimeInSec),
			"PeriodicScheduleEndTimeInSec":   PointerToNullInt64(j.Trigger.PeriodicSchedule.PeriodicScheduleEndTimeInSec),
			"IntervalSecond":                 PointerToNullInt64(j.Trigger.PeriodicSchedule.IntervalSecond),
			"CreatedAtInSec":                 j.CreatedAtInSec,
			"UpdatedAtInSec":                 j.UpdatedAtInSec,
			"PipelineId":                     j.PipelineSpec.PipelineId,
			"PipelineName":                   j.PipelineSpec.PipelineName,
			"PipelineSpecManifest":           j.PipelineSpec.PipelineSpecManifest,
			"WorkflowSpecManifest":           j.PipelineSpec.WorkflowSpecManifest,
			"Parameters":                     j.PipelineSpec.Parameters,
			"RuntimeParameters":              j.PipelineSpec.RuntimeConfig.Parameters,
			"PipelineRoot":                   j.PipelineSpec.RuntimeConfig.PipelineRoot,
			"ExperimentUUID":                 j.ExperimentId,
			"PipelineVersionId":              j.PipelineSpec.PipelineVersionId,
		}).ToSql()
	if err != nil {
		return nil, util.NewInternalServerError(err, "Failed to create query to add job to job table: %v",
			err.Error())
	}

	// Use a transaction to make sure both job and its resource references are stored.
	tx, err := s.db.Begin()
	if err != nil {
		return nil, util.NewInternalServerError(err, "Failed to create a new transaction to create job")
	}
	_, err = tx.Exec(jobSql, jobArgs...)
	if err != nil {
		tx.Rollback()
		return nil, util.NewInternalServerError(err, "Failed to store job %v to table", j.DisplayName)
	}

	// TODO(gkcalat): remove this workflow once we fully deprecate resource references
	// and provide logic for data migration for v1beta1 data.
	err = s.resourceReferenceStore.CreateResourceReferences(tx, j.ResourceReferences)
	if err != nil {
		tx.Rollback()
		return nil, util.NewInternalServerError(err, "Failed to store resource references to table for job %v ", j.DisplayName)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, util.NewInternalServerError(err, "Failed to store job %v and its resource references to table", j.DisplayName)
	}
	return j, nil
}

func (s *JobStore) ChangeJobMode(id string, enabled bool) error {
	now := s.time.Now().Unix()
	sql, args, err := sq.
		Update("jobs").
		SetMap(sq.Eq{
			"Enabled":        enabled,
			"UpdatedAtInSec": now,
		}).
		Where(sq.Eq{"UUID": string(id)}).
		Where(sq.Eq{"Enabled": !enabled}).
		ToSql()
	if err != nil {
		return util.NewInternalServerError(err, "Error when creating query to enable job %v to %v", id, enabled)
	}
	_, err = s.db.Exec(sql, args...)
	if err != nil {
		return util.NewInternalServerError(err, "Error when enabling job %v to %v", id, enabled)
	}
	return nil
}

func (s *JobStore) UpdateJob(swf *util.ScheduledWorkflow) error {
	now := s.time.Now().Unix()
	parameters, err := swf.ParametersAsString()
	if err != nil {
		return err
	}
	updateSql := sq.
		Update("jobs").
		SetMap(sq.Eq{
			"Name": swf.Name,
			// Namespace changes for recurring runs is forbidden
			// "Namespace":                      swf.Namespace,
			"Enabled":                        swf.Spec.Enabled,
			"Conditions":                     model.StatusState(swf.ConditionSummary()).ToString(),
			"MaxConcurrency":                 swf.MaxConcurrencyOr0(),
			"NoCatchup":                      swf.NoCatchupOrFalse(),
			"UpdatedAtInSec":                 now,
			"CronScheduleStartTimeInSec":     PointerToNullInt64(swf.CronScheduleStartTimeInSecOrNull()),
			"CronScheduleEndTimeInSec":       PointerToNullInt64(swf.CronScheduleEndTimeInSecOrNull()),
			"Schedule":                       swf.CronOrEmpty(),
			"PeriodicScheduleStartTimeInSec": PointerToNullInt64(swf.PeriodicScheduleStartTimeInSecOrNull()),
			"PeriodicScheduleEndTimeInSec":   PointerToNullInt64(swf.PeriodicScheduleEndTimeInSecOrNull()),
			"IntervalSecond":                 swf.IntervalSecondOr0(),
		})
	if len(parameters) > 0 {
		if swf.GetVersion() == util.SWFv1 {
			updateSql = updateSql.SetMap(sq.Eq{"Parameters": parameters})
		} else if swf.GetVersion() == util.SWFv2 {
			updateSql = updateSql.SetMap(sq.Eq{"RuntimeParameters": parameters})
		} else {
			return util.NewInternalServerError(util.NewInvalidInputError("ScheduledWorkflow has an invalid version: %v", swf.GetVersion()), "Failed to update job %v", swf.UID)
		}
	}
	sql, args, err := updateSql.Where(sq.Eq{"UUID": string(swf.UID)}).ToSql()
	if err != nil {
		return util.NewInternalServerError(err,
			"Error while creating query to update job with scheduled workflow: %v: %+v",
			err, swf.ScheduledWorkflow)
	}
	r, err := s.db.Exec(sql, args...)
	if err != nil {
		return util.NewInternalServerError(err,
			"Error while updating job with scheduled workflow: %v: %+v",
			err, swf.ScheduledWorkflow)
	}

	rowsAffected, err := r.RowsAffected()
	if err != nil {
		return util.NewInternalServerError(err,
			"Error getting affected rows while updating job with scheduled workflow: %v: %+v",
			err, swf.ScheduledWorkflow)
	}
	if rowsAffected <= 0 {
		return util.NewInvalidInputError(
			"There is no job corresponding to this scheduled workflow: %v/%v/%v",
			swf.UID, swf.Namespace, swf.Name)
	}

	return nil
}

// If pipelineStore is provided, it will be used instead of direct database queries for getting pipelines
// and pipeline versions.
func NewJobStore(db *DB, time util.TimeInterface, pipelineStore PipelineStoreInterface) *JobStore {
	return &JobStore{
		db:                     db,
		resourceReferenceStore: NewResourceReferenceStore(db, pipelineStore),
		time:                   time,
	}
}
