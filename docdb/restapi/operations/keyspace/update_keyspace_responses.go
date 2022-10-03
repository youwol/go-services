// Code generated by go-swagger; DO NOT EDIT.

package keyspace

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"platform/services/docdb/models"
)

// UpdateKeyspaceOKCode is the HTTP code returned for type UpdateKeyspaceOK
const UpdateKeyspaceOKCode int = 200

/*UpdateKeyspaceOK Ok

swagger:response updateKeyspaceOK
*/
type UpdateKeyspaceOK struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewUpdateKeyspaceOK creates UpdateKeyspaceOK with default headers values
func NewUpdateKeyspaceOK() *UpdateKeyspaceOK {

	return &UpdateKeyspaceOK{}
}

// WithPayload adds the payload to the update keyspace o k response
func (o *UpdateKeyspaceOK) WithPayload(payload *models.APIResponse) *UpdateKeyspaceOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update keyspace o k response
func (o *UpdateKeyspaceOK) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateKeyspaceOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateKeyspaceBadRequestCode is the HTTP code returned for type UpdateKeyspaceBadRequest
const UpdateKeyspaceBadRequestCode int = 400

/*UpdateKeyspaceBadRequest Bad request

swagger:response updateKeyspaceBadRequest
*/
type UpdateKeyspaceBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewUpdateKeyspaceBadRequest creates UpdateKeyspaceBadRequest with default headers values
func NewUpdateKeyspaceBadRequest() *UpdateKeyspaceBadRequest {

	return &UpdateKeyspaceBadRequest{}
}

// WithPayload adds the payload to the update keyspace bad request response
func (o *UpdateKeyspaceBadRequest) WithPayload(payload *models.APIResponse) *UpdateKeyspaceBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update keyspace bad request response
func (o *UpdateKeyspaceBadRequest) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateKeyspaceBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateKeyspaceUnauthorizedCode is the HTTP code returned for type UpdateKeyspaceUnauthorized
const UpdateKeyspaceUnauthorizedCode int = 401

/*UpdateKeyspaceUnauthorized Unauthorized

swagger:response updateKeyspaceUnauthorized
*/
type UpdateKeyspaceUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewUpdateKeyspaceUnauthorized creates UpdateKeyspaceUnauthorized with default headers values
func NewUpdateKeyspaceUnauthorized() *UpdateKeyspaceUnauthorized {

	return &UpdateKeyspaceUnauthorized{}
}

// WithPayload adds the payload to the update keyspace unauthorized response
func (o *UpdateKeyspaceUnauthorized) WithPayload(payload *models.APIResponse) *UpdateKeyspaceUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update keyspace unauthorized response
func (o *UpdateKeyspaceUnauthorized) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateKeyspaceUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateKeyspaceNotFoundCode is the HTTP code returned for type UpdateKeyspaceNotFound
const UpdateKeyspaceNotFoundCode int = 404

/*UpdateKeyspaceNotFound Not found

swagger:response updateKeyspaceNotFound
*/
type UpdateKeyspaceNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewUpdateKeyspaceNotFound creates UpdateKeyspaceNotFound with default headers values
func NewUpdateKeyspaceNotFound() *UpdateKeyspaceNotFound {

	return &UpdateKeyspaceNotFound{}
}

// WithPayload adds the payload to the update keyspace not found response
func (o *UpdateKeyspaceNotFound) WithPayload(payload *models.APIResponse) *UpdateKeyspaceNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update keyspace not found response
func (o *UpdateKeyspaceNotFound) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateKeyspaceNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateKeyspaceInternalServerErrorCode is the HTTP code returned for type UpdateKeyspaceInternalServerError
const UpdateKeyspaceInternalServerErrorCode int = 500

/*UpdateKeyspaceInternalServerError Internal error

swagger:response updateKeyspaceInternalServerError
*/
type UpdateKeyspaceInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewUpdateKeyspaceInternalServerError creates UpdateKeyspaceInternalServerError with default headers values
func NewUpdateKeyspaceInternalServerError() *UpdateKeyspaceInternalServerError {

	return &UpdateKeyspaceInternalServerError{}
}

// WithPayload adds the payload to the update keyspace internal server error response
func (o *UpdateKeyspaceInternalServerError) WithPayload(payload *models.APIResponse) *UpdateKeyspaceInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update keyspace internal server error response
func (o *UpdateKeyspaceInternalServerError) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateKeyspaceInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
