// Code generated by go-swagger; DO NOT EDIT.

package query

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"

	"platform/services/docdb/models"
)

// NewDeleteQueryParams creates a new DeleteQueryParams object
// with the default values initialized.
func NewDeleteQueryParams() DeleteQueryParams {

	var (
		// initialize parameters with default values

		ownerDefault = string("")
	)

	return DeleteQueryParams{
		Owner: &ownerDefault,
	}
}

// DeleteQueryParams contains all the bound params for the delete query operation
// typically these are obtained from a http.Request
//
// swagger:parameters deleteQuery
type DeleteQueryParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Query to select the entities that will be deleted
	  Required: true
	  In: body
	*/
	Delete *models.DeleteStatement
	/*Name of keyspace to use
	  Required: true
	  In: path
	*/
	KeyspaceName string
	/*For which owner do we do the delete ? Defaults to current user. Indicate a group path to act as a group (e.g.: /youwol-users/subgroup)
	  In: query
	  Default: ""
	*/
	Owner *string
	/*Name of table to return
	  Required: true
	  In: path
	*/
	TableName string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewDeleteQueryParams() beforehand.
func (o *DeleteQueryParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body models.DeleteStatement
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("delete", "body", ""))
			} else {
				res = append(res, errors.NewParseError("delete", "body", "", err))
			}
		} else {
			// validate body object
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.Delete = &body
			}
		}
	} else {
		res = append(res, errors.Required("delete", "body", ""))
	}
	rKeyspaceName, rhkKeyspaceName, _ := route.Params.GetOK("keyspaceName")
	if err := o.bindKeyspaceName(rKeyspaceName, rhkKeyspaceName, route.Formats); err != nil {
		res = append(res, err)
	}

	qOwner, qhkOwner, _ := qs.GetOK("owner")
	if err := o.bindOwner(qOwner, qhkOwner, route.Formats); err != nil {
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

// bindKeyspaceName binds and validates parameter KeyspaceName from path.
func (o *DeleteQueryParams) bindKeyspaceName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.KeyspaceName = raw

	return nil
}

// bindOwner binds and validates parameter Owner from query.
func (o *DeleteQueryParams) bindOwner(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewDeleteQueryParams()
		return nil
	}

	o.Owner = &raw

	return nil
}

// bindTableName binds and validates parameter TableName from path.
func (o *DeleteQueryParams) bindTableName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.TableName = raw

	return nil
}
