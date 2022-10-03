// Code generated by go-swagger; DO NOT EDIT.

package object

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// NewGetObjectParams creates a new GetObjectParams object
// with the default values initialized.
func NewGetObjectParams() GetObjectParams {

	var (
		// initialize parameters with default values

		isolationDefault = bool(true)

		ownerDefault = string("")
	)

	return GetObjectParams{
		Isolation: &isolationDefault,

		Owner: &ownerDefault,
	}
}

// GetObjectParams contains all the bound params for the get object operation
// typically these are obtained from a http.Request
//
// swagger:parameters getObject
type GetObjectParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: path
	*/
	BucketName string
	/*Should the server automatically isolate data between owners
	  In: query
	  Default: true
	*/
	Isolation *bool
	/*
	  Required: true
	  In: query
	*/
	ObjectName string
	/*Which owner does the object belong to ? Defaults to current user. Indicate a group path to act on a shared object
	  In: query
	  Default: ""
	*/
	Owner *string
	/*Optional encryption setting
	  In: query
	*/
	ServerSideEncryption *string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetObjectParams() beforehand.
func (o *GetObjectParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	rBucketName, rhkBucketName, _ := route.Params.GetOK("bucketName")
	if err := o.bindBucketName(rBucketName, rhkBucketName, route.Formats); err != nil {
		res = append(res, err)
	}

	qIsolation, qhkIsolation, _ := qs.GetOK("isolation")
	if err := o.bindIsolation(qIsolation, qhkIsolation, route.Formats); err != nil {
		res = append(res, err)
	}

	qObjectName, qhkObjectName, _ := qs.GetOK("objectName")
	if err := o.bindObjectName(qObjectName, qhkObjectName, route.Formats); err != nil {
		res = append(res, err)
	}

	qOwner, qhkOwner, _ := qs.GetOK("owner")
	if err := o.bindOwner(qOwner, qhkOwner, route.Formats); err != nil {
		res = append(res, err)
	}

	qServerSideEncryption, qhkServerSideEncryption, _ := qs.GetOK("server_side_encryption")
	if err := o.bindServerSideEncryption(qServerSideEncryption, qhkServerSideEncryption, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindBucketName binds and validates parameter BucketName from path.
func (o *GetObjectParams) bindBucketName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.BucketName = raw

	return nil
}

// bindIsolation binds and validates parameter Isolation from query.
func (o *GetObjectParams) bindIsolation(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewGetObjectParams()
		return nil
	}

	value, err := swag.ConvertBool(raw)
	if err != nil {
		return errors.InvalidType("isolation", "query", "bool", raw)
	}
	o.Isolation = &value

	return nil
}

// bindObjectName binds and validates parameter ObjectName from query.
func (o *GetObjectParams) bindObjectName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("objectName", "query", rawData)
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// AllowEmptyValue: false
	if err := validate.RequiredString("objectName", "query", raw); err != nil {
		return err
	}

	o.ObjectName = raw

	return nil
}

// bindOwner binds and validates parameter Owner from query.
func (o *GetObjectParams) bindOwner(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewGetObjectParams()
		return nil
	}

	o.Owner = &raw

	return nil
}

// bindServerSideEncryption binds and validates parameter ServerSideEncryption from query.
func (o *GetObjectParams) bindServerSideEncryption(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.ServerSideEncryption = &raw

	return nil
}
