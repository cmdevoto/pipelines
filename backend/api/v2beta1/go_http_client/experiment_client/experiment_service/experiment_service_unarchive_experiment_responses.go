// Code generated by go-swagger; DO NOT EDIT.

package experiment_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	experiment_model "github.com/kubeflow/pipelines/backend/api/v2beta1/go_http_client/experiment_model"
)

// ExperimentServiceUnarchiveExperimentReader is a Reader for the ExperimentServiceUnarchiveExperiment structure.
type ExperimentServiceUnarchiveExperimentReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ExperimentServiceUnarchiveExperimentReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewExperimentServiceUnarchiveExperimentOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		result := NewExperimentServiceUnarchiveExperimentDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewExperimentServiceUnarchiveExperimentOK creates a ExperimentServiceUnarchiveExperimentOK with default headers values
func NewExperimentServiceUnarchiveExperimentOK() *ExperimentServiceUnarchiveExperimentOK {
	return &ExperimentServiceUnarchiveExperimentOK{}
}

/*
ExperimentServiceUnarchiveExperimentOK handles this case with default header values.

A successful response.
*/
type ExperimentServiceUnarchiveExperimentOK struct {
	Payload interface{}
}

func (o *ExperimentServiceUnarchiveExperimentOK) Error() string {
	return fmt.Sprintf("[POST /apis/v2beta1/experiments/{experiment_id}:unarchive][%d] experimentServiceUnarchiveExperimentOK  %+v", 200, o.Payload)
}

func (o *ExperimentServiceUnarchiveExperimentOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewExperimentServiceUnarchiveExperimentDefault creates a ExperimentServiceUnarchiveExperimentDefault with default headers values
func NewExperimentServiceUnarchiveExperimentDefault(code int) *ExperimentServiceUnarchiveExperimentDefault {
	return &ExperimentServiceUnarchiveExperimentDefault{
		_statusCode: code,
	}
}

/*
ExperimentServiceUnarchiveExperimentDefault handles this case with default header values.

An unexpected error response.
*/
type ExperimentServiceUnarchiveExperimentDefault struct {
	_statusCode int

	Payload *experiment_model.RuntimeError
}

// Code gets the status code for the experiment service unarchive experiment default response
func (o *ExperimentServiceUnarchiveExperimentDefault) Code() int {
	return o._statusCode
}

func (o *ExperimentServiceUnarchiveExperimentDefault) Error() string {
	return fmt.Sprintf("[POST /apis/v2beta1/experiments/{experiment_id}:unarchive][%d] ExperimentService_UnarchiveExperiment default  %+v", o._statusCode, o.Payload)
}

func (o *ExperimentServiceUnarchiveExperimentDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(experiment_model.RuntimeError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
