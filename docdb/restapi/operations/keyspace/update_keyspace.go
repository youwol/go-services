// Code generated by go-swagger; DO NOT EDIT.

package keyspace

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"

	"platform/services/docdb/models"
)

// UpdateKeyspaceHandlerFunc turns a function with the right signature into a update keyspace handler
type UpdateKeyspaceHandlerFunc func(UpdateKeyspaceParams, *models.Principal) middleware.Responder

// Handle executing the request and returning a response
func (fn UpdateKeyspaceHandlerFunc) Handle(params UpdateKeyspaceParams, principal *models.Principal) middleware.Responder {
	return fn(params, principal)
}

// UpdateKeyspaceHandler interface for that can handle valid update keyspace params
type UpdateKeyspaceHandler interface {
	Handle(UpdateKeyspaceParams, *models.Principal) middleware.Responder
}

// NewUpdateKeyspace creates a new http.Handler for the update keyspace operation
func NewUpdateKeyspace(ctx *middleware.Context, handler UpdateKeyspaceHandler) *UpdateKeyspace {
	return &UpdateKeyspace{Context: ctx, Handler: handler}
}

/*UpdateKeyspace swagger:route PUT /keyspace keyspace updateKeyspace

Updates a keyspace in the store (keyspace name cannot be updated)

*/
type UpdateKeyspace struct {
	Context *middleware.Context
	Handler UpdateKeyspaceHandler
}

func (o *UpdateKeyspace) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewUpdateKeyspaceParams()

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
