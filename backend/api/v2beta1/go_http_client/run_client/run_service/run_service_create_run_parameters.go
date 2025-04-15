// Code generated by go-swagger; DO NOT EDIT.

package run_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"

	run_model "github.com/kubeflow/pipelines/backend/api/v2beta1/go_http_client/run_model"
)

// NewRunServiceCreateRunParams creates a new RunServiceCreateRunParams object
// with the default values initialized.
func NewRunServiceCreateRunParams() *RunServiceCreateRunParams {
	var ()
	return &RunServiceCreateRunParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewRunServiceCreateRunParamsWithTimeout creates a new RunServiceCreateRunParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewRunServiceCreateRunParamsWithTimeout(timeout time.Duration) *RunServiceCreateRunParams {
	var ()
	return &RunServiceCreateRunParams{

		timeout: timeout,
	}
}

// NewRunServiceCreateRunParamsWithContext creates a new RunServiceCreateRunParams object
// with the default values initialized, and the ability to set a context for a request
func NewRunServiceCreateRunParamsWithContext(ctx context.Context) *RunServiceCreateRunParams {
	var ()
	return &RunServiceCreateRunParams{

		Context: ctx,
	}
}

// NewRunServiceCreateRunParamsWithHTTPClient creates a new RunServiceCreateRunParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewRunServiceCreateRunParamsWithHTTPClient(client *http.Client) *RunServiceCreateRunParams {
	var ()
	return &RunServiceCreateRunParams{
		HTTPClient: client,
	}
}

/*
RunServiceCreateRunParams contains all the parameters to send to the API endpoint
for the run service create run operation typically these are written to a http.Request
*/
type RunServiceCreateRunParams struct {

	/*Body
	  Run to be created.

	*/
	Body *run_model.V2beta1Run
	/*ExperimentID
	  The ID of the parent experiment.

	*/
	ExperimentID *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the run service create run params
func (o *RunServiceCreateRunParams) WithTimeout(timeout time.Duration) *RunServiceCreateRunParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the run service create run params
func (o *RunServiceCreateRunParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the run service create run params
func (o *RunServiceCreateRunParams) WithContext(ctx context.Context) *RunServiceCreateRunParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the run service create run params
func (o *RunServiceCreateRunParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the run service create run params
func (o *RunServiceCreateRunParams) WithHTTPClient(client *http.Client) *RunServiceCreateRunParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the run service create run params
func (o *RunServiceCreateRunParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the run service create run params
func (o *RunServiceCreateRunParams) WithBody(body *run_model.V2beta1Run) *RunServiceCreateRunParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the run service create run params
func (o *RunServiceCreateRunParams) SetBody(body *run_model.V2beta1Run) {
	o.Body = body
}

// WithExperimentID adds the experimentID to the run service create run params
func (o *RunServiceCreateRunParams) WithExperimentID(experimentID *string) *RunServiceCreateRunParams {
	o.SetExperimentID(experimentID)
	return o
}

// SetExperimentID adds the experimentId to the run service create run params
func (o *RunServiceCreateRunParams) SetExperimentID(experimentID *string) {
	o.ExperimentID = experimentID
}

// WriteToRequest writes these params to a swagger request
func (o *RunServiceCreateRunParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
		}
	}

	if o.ExperimentID != nil {

		// query param experiment_id
		var qrExperimentID string
		if o.ExperimentID != nil {
			qrExperimentID = *o.ExperimentID
		}
		qExperimentID := qrExperimentID
		if qExperimentID != "" {
			if err := r.SetQueryParam("experiment_id", qExperimentID); err != nil {
				return err
			}
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
