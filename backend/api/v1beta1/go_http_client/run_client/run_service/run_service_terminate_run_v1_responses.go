// Code generated by go-swagger; DO NOT EDIT.

package run_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	run_model "github.com/kubeflow/pipelines/backend/api/v1beta1/go_http_client/run_model"
)

// RunServiceTerminateRunV1Reader is a Reader for the RunServiceTerminateRunV1 structure.
type RunServiceTerminateRunV1Reader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *RunServiceTerminateRunV1Reader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewRunServiceTerminateRunV1OK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		result := NewRunServiceTerminateRunV1Default(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewRunServiceTerminateRunV1OK creates a RunServiceTerminateRunV1OK with default headers values
func NewRunServiceTerminateRunV1OK() *RunServiceTerminateRunV1OK {
	return &RunServiceTerminateRunV1OK{}
}

/*
RunServiceTerminateRunV1OK handles this case with default header values.

A successful response.
*/
type RunServiceTerminateRunV1OK struct {
	Payload interface{}
}

func (o *RunServiceTerminateRunV1OK) Error() string {
	return fmt.Sprintf("[POST /apis/v1beta1/runs/{run_id}/terminate][%d] runServiceTerminateRunV1OK  %+v", 200, o.Payload)
}

func (o *RunServiceTerminateRunV1OK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewRunServiceTerminateRunV1Default creates a RunServiceTerminateRunV1Default with default headers values
func NewRunServiceTerminateRunV1Default(code int) *RunServiceTerminateRunV1Default {
	return &RunServiceTerminateRunV1Default{
		_statusCode: code,
	}
}

/*
RunServiceTerminateRunV1Default handles this case with default header values.

An unexpected error response.
*/
type RunServiceTerminateRunV1Default struct {
	_statusCode int

	Payload *run_model.GatewayruntimeError
}

// Code gets the status code for the run service terminate run v1 default response
func (o *RunServiceTerminateRunV1Default) Code() int {
	return o._statusCode
}

func (o *RunServiceTerminateRunV1Default) Error() string {
	return fmt.Sprintf("[POST /apis/v1beta1/runs/{run_id}/terminate][%d] RunService_TerminateRunV1 default  %+v", o._statusCode, o.Payload)
}

func (o *RunServiceTerminateRunV1Default) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(run_model.GatewayruntimeError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
