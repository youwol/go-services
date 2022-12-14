// Code generated by go-swagger; DO NOT EDIT.

package query

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"

	"platform/services/docdb/models"
)

// SelectQueryHandlerFunc turns a function with the right signature into a select query handler
type SelectQueryHandlerFunc func(SelectQueryParams, *models.Principal) middleware.Responder

// Handle executing the request and returning a response
func (fn SelectQueryHandlerFunc) Handle(params SelectQueryParams, principal *models.Principal) middleware.Responder {
	return fn(params, principal)
}

// SelectQueryHandler interface for that can handle valid select query params
type SelectQueryHandler interface {
	Handle(SelectQueryParams, *models.Principal) middleware.Responder
}

// NewSelectQuery creates a new http.Handler for the select query operation
func NewSelectQuery(ctx *middleware.Context, handler SelectQueryHandler) *SelectQuery {
	return &SelectQuery{Context: ctx, Handler: handler}
}

/*SelectQuery swagger:route POST /{keyspaceName}/{tableName}/query query selectQuery

Retrieves custom columns for a group of entities

Returns a list of enities or compact (array) data

*/
type SelectQuery struct {
	Context *middleware.Context
	Handler SelectQueryHandler
}

func (o *SelectQuery) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewSelectQueryParams()

	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		r = aCtx
	}
	var principal *models.Principal
	if uprinc != nil {
		principal = uprinc.(*models.Principal) // this is really a models.Principal, I promise
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
