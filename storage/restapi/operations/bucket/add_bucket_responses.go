// Code generated by go-swagger; DO NOT EDIT.

package bucket

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"platform/services/storage/models"
)

// AddBucketCreatedCode is the HTTP code returned for type AddBucketCreated
const AddBucketCreatedCode int = 201

/*AddBucketCreated Created

swagger:response addBucketCreated
*/
type AddBucketCreated struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewAddBucketCreated creates AddBucketCreated with default headers values
func NewAddBucketCreated() *AddBucketCreated {

	return &AddBucketCreated{}
}

// WithPayload adds the payload to the add bucket created response
func (o *AddBucketCreated) WithPayload(payload *models.APIResponse) *AddBucketCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add bucket created response
func (o *AddBucketCreated) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddBucketCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// AddBucketBadRequestCode is the HTTP code returned for type AddBucketBadRequest
const AddBucketBadRequestCode int = 400

/*AddBucketBadRequest Bad request

swagger:response addBucketBadRequest
*/
type AddBucketBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewAddBucketBadRequest creates AddBucketBadRequest with default headers values
func NewAddBucketBadRequest() *AddBucketBadRequest {

	return &AddBucketBadRequest{}
}

// WithPayload adds the payload to the add bucket bad request response
func (o *AddBucketBadRequest) WithPayload(payload *models.APIResponse) *AddBucketBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add bucket bad request response
func (o *AddBucketBadRequest) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddBucketBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// AddBucketInternalServerErrorCode is the HTTP code returned for type AddBucketInternalServerError
const AddBucketInternalServerErrorCode int = 500

/*AddBucketInternalServerError Internal error

swagger:response addBucketInternalServerError
*/
type AddBucketInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewAddBucketInternalServerError creates AddBucketInternalServerError with default headers values
func NewAddBucketInternalServerError() *AddBucketInternalServerError {

	return &AddBucketInternalServerError{}
}

// WithPayload adds the payload to the add bucket internal server error response
func (o *AddBucketInternalServerError) WithPayload(payload *models.APIResponse) *AddBucketInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add bucket internal server error response
func (o *AddBucketInternalServerError) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddBucketInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
