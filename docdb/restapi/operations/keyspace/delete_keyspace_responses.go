// Code generated by go-swagger; DO NOT EDIT.

package keyspace

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"platform/services/docdb/models"
)

// DeleteKeyspaceOKCode is the HTTP code returned for type DeleteKeyspaceOK
const DeleteKeyspaceOKCode int = 200

/*DeleteKeyspaceOK Ok

swagger:response deleteKeyspaceOK
*/
type DeleteKeyspaceOK struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewDeleteKeyspaceOK creates DeleteKeyspaceOK with default headers values
func NewDeleteKeyspaceOK() *DeleteKeyspaceOK {

	return &DeleteKeyspaceOK{}
}

// WithPayload adds the payload to the delete keyspace o k response
func (o *DeleteKeyspaceOK) WithPayload(payload *models.APIResponse) *DeleteKeyspaceOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete keyspace o k response
func (o *DeleteKeyspaceOK) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteKeyspaceOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// DeleteKeyspaceBadRequestCode is the HTTP code returned for type DeleteKeyspaceBadRequest
const DeleteKeyspaceBadRequestCode int = 400

/*DeleteKeyspaceBadRequest Bad request

swagger:response deleteKeyspaceBadRequest
*/
type DeleteKeyspaceBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewDeleteKeyspaceBadRequest creates DeleteKeyspaceBadRequest with default headers values
func NewDeleteKeyspaceBadRequest() *DeleteKeyspaceBadRequest {

	return &DeleteKeyspaceBadRequest{}
}

// WithPayload adds the payload to the delete keyspace bad request response
func (o *DeleteKeyspaceBadRequest) WithPayload(payload *models.APIResponse) *DeleteKeyspaceBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete keyspace bad request response
func (o *DeleteKeyspaceBadRequest) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteKeyspaceBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// DeleteKeyspaceUnauthorizedCode is the HTTP code returned for type DeleteKeyspaceUnauthorized
const DeleteKeyspaceUnauthorizedCode int = 401

/*DeleteKeyspaceUnauthorized Unauthorized

swagger:response deleteKeyspaceUnauthorized
*/
type DeleteKeyspaceUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewDeleteKeyspaceUnauthorized creates DeleteKeyspaceUnauthorized with default headers values
func NewDeleteKeyspaceUnauthorized() *DeleteKeyspaceUnauthorized {

	return &DeleteKeyspaceUnauthorized{}
}

// WithPayload adds the payload to the delete keyspace unauthorized response
func (o *DeleteKeyspaceUnauthorized) WithPayload(payload *models.APIResponse) *DeleteKeyspaceUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete keyspace unauthorized response
func (o *DeleteKeyspaceUnauthorized) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteKeyspaceUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// DeleteKeyspaceNotFoundCode is the HTTP code returned for type DeleteKeyspaceNotFound
const DeleteKeyspaceNotFoundCode int = 404

/*DeleteKeyspaceNotFound Not found

swagger:response deleteKeyspaceNotFound
*/
type DeleteKeyspaceNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewDeleteKeyspaceNotFound creates DeleteKeyspaceNotFound with default headers values
func NewDeleteKeyspaceNotFound() *DeleteKeyspaceNotFound {

	return &DeleteKeyspaceNotFound{}
}

// WithPayload adds the payload to the delete keyspace not found response
func (o *DeleteKeyspaceNotFound) WithPayload(payload *models.APIResponse) *DeleteKeyspaceNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete keyspace not found response
func (o *DeleteKeyspaceNotFound) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteKeyspaceNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// DeleteKeyspaceInternalServerErrorCode is the HTTP code returned for type DeleteKeyspaceInternalServerError
const DeleteKeyspaceInternalServerErrorCode int = 500

/*DeleteKeyspaceInternalServerError Internal error

swagger:response deleteKeyspaceInternalServerError
*/
type DeleteKeyspaceInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewDeleteKeyspaceInternalServerError creates DeleteKeyspaceInternalServerError with default headers values
func NewDeleteKeyspaceInternalServerError() *DeleteKeyspaceInternalServerError {

	return &DeleteKeyspaceInternalServerError{}
}

// WithPayload adds the payload to the delete keyspace internal server error response
func (o *DeleteKeyspaceInternalServerError) WithPayload(payload *models.APIResponse) *DeleteKeyspaceInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete keyspace internal server error response
func (o *DeleteKeyspaceInternalServerError) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteKeyspaceInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
