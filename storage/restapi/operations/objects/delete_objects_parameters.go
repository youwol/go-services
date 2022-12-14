// Code generated by go-swagger; DO NOT EDIT.

package objects

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// NewDeleteObjectsParams creates a new DeleteObjectsParams object
// with the default values initialized.
func NewDeleteObjectsParams() DeleteObjectsParams {

	var (
		// initialize parameters with default values

		isolationDefault = bool(true)
		ownerDefault     = string("")
	)

	return DeleteObjectsParams{
		Isolation: &isolationDefault,

		Owner: &ownerDefault,
	}
}

// DeleteObjectsParams contains all the bound params for the delete objects operation
// typically these are obtained from a http.Request
//
// swagger:parameters deleteObjects
type DeleteObjectsParams struct {

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
	/*Which owner does the object belong to ? Defaults to current user. Indicate a group path to act on a shared object
	  In: query
	  Default: ""
	*/
	Owner *string
	/*Prefix of names to retrieve
	  In: query
	*/
	Prefix *string
	/*Whether to look in subfolders
	  In: query
	*/
	Recursive *bool
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewDeleteObjectsParams() beforehand.
func (o *DeleteObjectsParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
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

	qOwner, qhkOwner, _ := qs.GetOK("owner")
	if err := o.bindOwner(qOwner, qhkOwner, route.Formats); err != nil {
		res = append(res, err)
	}

	qPrefix, qhkPrefix, _ := qs.GetOK("prefix")
	if err := o.bindPrefix(qPrefix, qhkPrefix, route.Formats); err != nil {
		res = append(res, err)
	}

	qRecursive, qhkRecursive, _ := qs.GetOK("recursive")
	if err := o.bindRecursive(qRecursive, qhkRecursive, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindBucketName binds and validates parameter BucketName from path.
func (o *DeleteObjectsParams) bindBucketName(rawData []string, hasKey bool, formats strfmt.Registry) error {
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
func (o *DeleteObjectsParams) bindIsolation(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewDeleteObjectsParams()
		return nil
	}

	value, err := swag.ConvertBool(raw)
	if err != nil {
		return errors.InvalidType("isolation", "query", "bool", raw)
	}
	o.Isolation = &value

	return nil
}

// bindOwner binds and validates parameter Owner from query.
func (o *DeleteObjectsParams) bindOwner(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewDeleteObjectsParams()
		return nil
	}

	o.Owner = &raw

	return nil
}

// bindPrefix binds and validates parameter Prefix from query.
func (o *DeleteObjectsParams) bindPrefix(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.Prefix = &raw

	return nil
}

// bindRecursive binds and validates parameter Recursive from query.
func (o *DeleteObjectsParams) bindRecursive(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	value, err := swag.ConvertBool(raw)
	if err != nil {
		return errors.InvalidType("recursive", "query", "bool", raw)
	}
	o.Recursive = &value

	return nil
}
