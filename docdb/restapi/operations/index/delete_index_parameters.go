// Code generated by go-swagger; DO NOT EDIT.

package index

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// NewDeleteIndexParams creates a new DeleteIndexParams object
// no default values defined in spec.
func NewDeleteIndexParams() DeleteIndexParams {

	return DeleteIndexParams{}
}

// DeleteIndexParams contains all the bound params for the delete index operation
// typically these are obtained from a http.Request
//
// swagger:parameters deleteIndex
type DeleteIndexParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: query
	*/
	IndexName string
	/*Name of keyspace
	  Required: true
	  In: path
	*/
	KeyspaceName string
	/*Name of table
	  Required: true
	  In: path
	*/
	TableName string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewDeleteIndexParams() beforehand.
func (o *DeleteIndexParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qIndexName, qhkIndexName, _ := qs.GetOK("indexName")
	if err := o.bindIndexName(qIndexName, qhkIndexName, route.Formats); err != nil {
		res = append(res, err)
	}

	rKeyspaceName, rhkKeyspaceName, _ := route.Params.GetOK("keyspaceName")
	if err := o.bindKeyspaceName(rKeyspaceName, rhkKeyspaceName, route.Formats); err != nil {
		res = append(res, err)
	}

	rTableName, rhkTableName, _ := route.Params.GetOK("tableName")
	if err := o.bindTableName(rTableName, rhkTableName, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindIndexName binds and validates parameter IndexName from query.
func (o *DeleteIndexParams) bindIndexName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("indexName", "query", rawData)
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// AllowEmptyValue: false
	if err := validate.RequiredString("indexName", "query", raw); err != nil {
		return err
	}

	o.IndexName = raw

	return nil
}

// bindKeyspaceName binds and validates parameter KeyspaceName from path.
func (o *DeleteIndexParams) bindKeyspaceName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.KeyspaceName = raw

	return nil
}

// bindTableName binds and validates parameter TableName from path.
func (o *DeleteIndexParams) bindTableName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.TableName = raw

	return nil
}
