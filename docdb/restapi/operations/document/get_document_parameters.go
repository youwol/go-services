// Code generated by go-swagger; DO NOT EDIT.

package document

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

// NewGetDocumentParams creates a new GetDocumentParams object
// with the default values initialized.
func NewGetDocumentParams() GetDocumentParams {

	var (
		// initialize parameters with default values

		ownerDefault = string("")
	)

	return GetDocumentParams{
		Owner: &ownerDefault,
	}
}

// GetDocumentParams contains all the bound params for the get document operation
// typically these are obtained from a http.Request
//
// swagger:parameters getDocument
type GetDocumentParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Key of entity to return
	  In: query
	*/
	ClusteringKey []string
	/*Name of keyspace to use
	  Required: true
	  In: path
	*/
	KeyspaceName string
	/*Which owner does the document belong to ? Defaults to current user. Indicate a group path to act as a group (e.g.: /youwol-users/subgroup)
	  In: query
	  Default: ""
	*/
	Owner *string
	/*Primary key of entity to return
	  Required: true
	  In: query
	*/
	PartitionKey []string
	/*Name of table to return
	  Required: true
	  In: path
	*/
	TableName string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetDocumentParams() beforehand.
func (o *GetDocumentParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qClusteringKey, qhkClusteringKey, _ := qs.GetOK("clusteringKey")
	if err := o.bindClusteringKey(qClusteringKey, qhkClusteringKey, route.Formats); err != nil {
		res = append(res, err)
	}

	rKeyspaceName, rhkKeyspaceName, _ := route.Params.GetOK("keyspaceName")
	if err := o.bindKeyspaceName(rKeyspaceName, rhkKeyspaceName, route.Formats); err != nil {
		res = append(res, err)
	}

	qOwner, qhkOwner, _ := qs.GetOK("owner")
	if err := o.bindOwner(qOwner, qhkOwner, route.Formats); err != nil {
		res = append(res, err)
	}

	qPartitionKey, qhkPartitionKey, _ := qs.GetOK("partitionKey")
	if err := o.bindPartitionKey(qPartitionKey, qhkPartitionKey, route.Formats); err != nil {
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

// bindClusteringKey binds and validates array parameter ClusteringKey from query.
//
// Arrays are parsed according to CollectionFormat: "" (defaults to "csv" when empty).
func (o *GetDocumentParams) bindClusteringKey(rawData []string, hasKey bool, formats strfmt.Registry) error {

	var qvClusteringKey string
	if len(rawData) > 0 {
		qvClusteringKey = rawData[len(rawData)-1]
	}

	// CollectionFormat:
	clusteringKeyIC := swag.SplitByFormat(qvClusteringKey, "")
	if len(clusteringKeyIC) == 0 {
		return nil
	}

	var clusteringKeyIR []string
	for _, clusteringKeyIV := range clusteringKeyIC {
		clusteringKeyI := clusteringKeyIV

		clusteringKeyIR = append(clusteringKeyIR, clusteringKeyI)
	}

	o.ClusteringKey = clusteringKeyIR

	return nil
}

// bindKeyspaceName binds and validates parameter KeyspaceName from path.
func (o *GetDocumentParams) bindKeyspaceName(rawData []string, hasKey bool, formats strfmt.Registry) error {
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
func (o *GetDocumentParams) bindOwner(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewGetDocumentParams()
		return nil
	}

	o.Owner = &raw

	return nil
}

// bindPartitionKey binds and validates array parameter PartitionKey from query.
//
// Arrays are parsed according to CollectionFormat: "" (defaults to "csv" when empty).
func (o *GetDocumentParams) bindPartitionKey(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("partitionKey", "query", rawData)
	}

	var qvPartitionKey string
	if len(rawData) > 0 {
		qvPartitionKey = rawData[len(rawData)-1]
	}

	// CollectionFormat:
	partitionKeyIC := swag.SplitByFormat(qvPartitionKey, "")

	if len(partitionKeyIC) == 0 {
		return errors.Required("partitionKey", "query", partitionKeyIC)
	}

	var partitionKeyIR []string
	for _, partitionKeyIV := range partitionKeyIC {
		partitionKeyI := partitionKeyIV

		partitionKeyIR = append(partitionKeyIR, partitionKeyI)
	}

	o.PartitionKey = partitionKeyIR

	return nil
}

// bindTableName binds and validates parameter TableName from path.
func (o *GetDocumentParams) bindTableName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.TableName = raw

	return nil
}
