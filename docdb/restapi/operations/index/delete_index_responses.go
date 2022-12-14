// Code generated by go-swagger; DO NOT EDIT.

package index

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"platform/services/docdb/models"
)

// DeleteIndexOKCode is the HTTP code returned for type DeleteIndexOK
const DeleteIndexOKCode int = 200

/*DeleteIndexOK Ok

swagger:response deleteIndexOK
*/
type DeleteIndexOK struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewDeleteIndexOK creates DeleteIndexOK with default headers values
func NewDeleteIndexOK() *DeleteIndexOK {

	return &DeleteIndexOK{}
}

// WithPayload adds the payload to the delete index o k response
func (o *DeleteIndexOK) WithPayload(payload *models.APIResponse) *DeleteIndexOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete index o k response
func (o *DeleteIndexOK) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteIndexOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// DeleteIndexBadRequestCode is the HTTP code returned for type DeleteIndexBadRequest
const DeleteIndexBadRequestCode int = 400

/*DeleteIndexBadRequest Bad request

swagger:response deleteIndexBadRequest
*/
type DeleteIndexBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewDeleteIndexBadRequest creates DeleteIndexBadRequest with default headers values
func NewDeleteIndexBadRequest() *DeleteIndexBadRequest {

	return &DeleteIndexBadRequest{}
}

// WithPayload adds the payload to the delete index bad request response
func (o *DeleteIndexBadRequest) WithPayload(payload *models.APIResponse) *DeleteIndexBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete index bad request response
func (o *DeleteIndexBadRequest) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteIndexBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// DeleteIndexUnauthorizedCode is the HTTP code returned for type DeleteIndexUnauthorized
const DeleteIndexUnauthorizedCode int = 401

/*DeleteIndexUnauthorized Unauthorized

swagger:response deleteIndexUnauthorized
*/
type DeleteIndexUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewDeleteIndexUnauthorized creates DeleteIndexUnauthorized with default headers values
func NewDeleteIndexUnauthorized() *DeleteIndexUnauthorized {

	return &DeleteIndexUnauthorized{}
}

// WithPayload adds the payload to the delete index unauthorized response
func (o *DeleteIndexUnauthorized) WithPayload(payload *models.APIResponse) *DeleteIndexUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete index unauthorized response
func (o *DeleteIndexUnauthorized) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteIndexUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// DeleteIndexNotFoundCode is the HTTP code returned for type DeleteIndexNotFound
const DeleteIndexNotFoundCode int = 404

/*DeleteIndexNotFound Not found

swagger:response deleteIndexNotFound
*/
type DeleteIndexNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewDeleteIndexNotFound creates DeleteIndexNotFound with default headers values
func NewDeleteIndexNotFound() *DeleteIndexNotFound {

	return &DeleteIndexNotFound{}
}

// WithPayload adds the payload to the delete index not found response
func (o *DeleteIndexNotFound) WithPayload(payload *models.APIResponse) *DeleteIndexNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete index not found response
func (o *DeleteIndexNotFound) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteIndexNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// DeleteIndexInternalServerErrorCode is the HTTP code returned for type DeleteIndexInternalServerError
const DeleteIndexInternalServerErrorCode int = 500

/*DeleteIndexInternalServerError Internal error

swagger:response deleteIndexInternalServerError
*/
type DeleteIndexInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewDeleteIndexInternalServerError creates DeleteIndexInternalServerError with default headers values
func NewDeleteIndexInternalServerError() *DeleteIndexInternalServerError {

	return &DeleteIndexInternalServerError{}
}

// WithPayload adds the payload to the delete index internal server error response
func (o *DeleteIndexInternalServerError) WithPayload(payload *models.APIResponse) *DeleteIndexInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete index internal server error response
func (o *DeleteIndexInternalServerError) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteIndexInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
