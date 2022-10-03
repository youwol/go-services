// Code generated by go-swagger; DO NOT EDIT.

package query

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"platform/services/docdb/models"
)

// DeleteQueryOKCode is the HTTP code returned for type DeleteQueryOK
const DeleteQueryOKCode int = 200

/*DeleteQueryOK Ok

swagger:response deleteQueryOK
*/
type DeleteQueryOK struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewDeleteQueryOK creates DeleteQueryOK with default headers values
func NewDeleteQueryOK() *DeleteQueryOK {

	return &DeleteQueryOK{}
}

// WithPayload adds the payload to the delete query o k response
func (o *DeleteQueryOK) WithPayload(payload *models.APIResponse) *DeleteQueryOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete query o k response
func (o *DeleteQueryOK) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteQueryOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// DeleteQueryBadRequestCode is the HTTP code returned for type DeleteQueryBadRequest
const DeleteQueryBadRequestCode int = 400

/*DeleteQueryBadRequest Bad request

swagger:response deleteQueryBadRequest
*/
type DeleteQueryBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewDeleteQueryBadRequest creates DeleteQueryBadRequest with default headers values
func NewDeleteQueryBadRequest() *DeleteQueryBadRequest {

	return &DeleteQueryBadRequest{}
}

// WithPayload adds the payload to the delete query bad request response
func (o *DeleteQueryBadRequest) WithPayload(payload *models.APIResponse) *DeleteQueryBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete query bad request response
func (o *DeleteQueryBadRequest) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteQueryBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// DeleteQueryUnauthorizedCode is the HTTP code returned for type DeleteQueryUnauthorized
const DeleteQueryUnauthorizedCode int = 401

/*DeleteQueryUnauthorized Unauthorized

swagger:response deleteQueryUnauthorized
*/
type DeleteQueryUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewDeleteQueryUnauthorized creates DeleteQueryUnauthorized with default headers values
func NewDeleteQueryUnauthorized() *DeleteQueryUnauthorized {

	return &DeleteQueryUnauthorized{}
}

// WithPayload adds the payload to the delete query unauthorized response
func (o *DeleteQueryUnauthorized) WithPayload(payload *models.APIResponse) *DeleteQueryUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete query unauthorized response
func (o *DeleteQueryUnauthorized) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteQueryUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// DeleteQueryNotFoundCode is the HTTP code returned for type DeleteQueryNotFound
const DeleteQueryNotFoundCode int = 404

/*DeleteQueryNotFound Not found

swagger:response deleteQueryNotFound
*/
type DeleteQueryNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewDeleteQueryNotFound creates DeleteQueryNotFound with default headers values
func NewDeleteQueryNotFound() *DeleteQueryNotFound {

	return &DeleteQueryNotFound{}
}

// WithPayload adds the payload to the delete query not found response
func (o *DeleteQueryNotFound) WithPayload(payload *models.APIResponse) *DeleteQueryNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete query not found response
func (o *DeleteQueryNotFound) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteQueryNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// DeleteQueryInternalServerErrorCode is the HTTP code returned for type DeleteQueryInternalServerError
const DeleteQueryInternalServerErrorCode int = 500

/*DeleteQueryInternalServerError Internal error

swagger:response deleteQueryInternalServerError
*/
type DeleteQueryInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewDeleteQueryInternalServerError creates DeleteQueryInternalServerError with default headers values
func NewDeleteQueryInternalServerError() *DeleteQueryInternalServerError {

	return &DeleteQueryInternalServerError{}
}

// WithPayload adds the payload to the delete query internal server error response
func (o *DeleteQueryInternalServerError) WithPayload(payload *models.APIResponse) *DeleteQueryInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete query internal server error response
func (o *DeleteQueryInternalServerError) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteQueryInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
