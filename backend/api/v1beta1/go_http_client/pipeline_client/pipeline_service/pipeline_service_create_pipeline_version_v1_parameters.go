// Code generated by go-swagger; DO NOT EDIT.

package pipeline_service

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

	pipeline_model "github.com/kubeflow/pipelines/backend/api/v1beta1/go_http_client/pipeline_model"
)

// NewPipelineServiceCreatePipelineVersionV1Params creates a new PipelineServiceCreatePipelineVersionV1Params object
// with the default values initialized.
func NewPipelineServiceCreatePipelineVersionV1Params() *PipelineServiceCreatePipelineVersionV1Params {
	var ()
	return &PipelineServiceCreatePipelineVersionV1Params{

		timeout: cr.DefaultTimeout,
	}
}

// NewPipelineServiceCreatePipelineVersionV1ParamsWithTimeout creates a new PipelineServiceCreatePipelineVersionV1Params object
// with the default values initialized, and the ability to set a timeout on a request
func NewPipelineServiceCreatePipelineVersionV1ParamsWithTimeout(timeout time.Duration) *PipelineServiceCreatePipelineVersionV1Params {
	var ()
	return &PipelineServiceCreatePipelineVersionV1Params{

		timeout: timeout,
	}
}

// NewPipelineServiceCreatePipelineVersionV1ParamsWithContext creates a new PipelineServiceCreatePipelineVersionV1Params object
// with the default values initialized, and the ability to set a context for a request
func NewPipelineServiceCreatePipelineVersionV1ParamsWithContext(ctx context.Context) *PipelineServiceCreatePipelineVersionV1Params {
	var ()
	return &PipelineServiceCreatePipelineVersionV1Params{

		Context: ctx,
	}
}

// NewPipelineServiceCreatePipelineVersionV1ParamsWithHTTPClient creates a new PipelineServiceCreatePipelineVersionV1Params object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewPipelineServiceCreatePipelineVersionV1ParamsWithHTTPClient(client *http.Client) *PipelineServiceCreatePipelineVersionV1Params {
	var ()
	return &PipelineServiceCreatePipelineVersionV1Params{
		HTTPClient: client,
	}
}

/*
PipelineServiceCreatePipelineVersionV1Params contains all the parameters to send to the API endpoint
for the pipeline service create pipeline version v1 operation typically these are written to a http.Request
*/
type PipelineServiceCreatePipelineVersionV1Params struct {

	/*Body
	  ResourceReference inside PipelineVersion specifies the pipeline that this
	version belongs to.

	*/
	Body *pipeline_model.APIPipelineVersion

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the pipeline service create pipeline version v1 params
func (o *PipelineServiceCreatePipelineVersionV1Params) WithTimeout(timeout time.Duration) *PipelineServiceCreatePipelineVersionV1Params {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the pipeline service create pipeline version v1 params
func (o *PipelineServiceCreatePipelineVersionV1Params) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the pipeline service create pipeline version v1 params
func (o *PipelineServiceCreatePipelineVersionV1Params) WithContext(ctx context.Context) *PipelineServiceCreatePipelineVersionV1Params {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the pipeline service create pipeline version v1 params
func (o *PipelineServiceCreatePipelineVersionV1Params) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the pipeline service create pipeline version v1 params
func (o *PipelineServiceCreatePipelineVersionV1Params) WithHTTPClient(client *http.Client) *PipelineServiceCreatePipelineVersionV1Params {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the pipeline service create pipeline version v1 params
func (o *PipelineServiceCreatePipelineVersionV1Params) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the pipeline service create pipeline version v1 params
func (o *PipelineServiceCreatePipelineVersionV1Params) WithBody(body *pipeline_model.APIPipelineVersion) *PipelineServiceCreatePipelineVersionV1Params {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the pipeline service create pipeline version v1 params
func (o *PipelineServiceCreatePipelineVersionV1Params) SetBody(body *pipeline_model.APIPipelineVersion) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *PipelineServiceCreatePipelineVersionV1Params) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
