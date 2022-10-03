// Code generated by go-swagger; DO NOT EDIT.

package table

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"

	"platform/services/docdb/models"
)

// AddTableHandlerFunc turns a function with the right signature into a add table handler
type AddTableHandlerFunc func(AddTableParams, *models.Principal) middleware.Responder

// Handle executing the request and returning a response
func (fn AddTableHandlerFunc) Handle(params AddTableParams, principal *models.Principal) middleware.Responder {
	return fn(params, principal)
}

// AddTableHandler interface for that can handle valid add table params
type AddTableHandler interface {
	Handle(AddTableParams, *models.Principal) middleware.Responder
}

// NewAddTable creates a new http.Handler for the add table operation
func NewAddTable(ctx *middleware.Context, handler AddTableHandler) *AddTable {
	return &AddTable{Context: ctx, Handler: handler}
}

/*AddTable swagger:route POST /{keyspaceName}/table table addTable

Add a new table (@see: https://docs.scylladb.com/getting-started/ddl/#create-table-statement)

*/
type AddTable struct {
	Context *middleware.Context
	Handler AddTableHandler
}

func (o *AddTable) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewAddTableParams()

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
