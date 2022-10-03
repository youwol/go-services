// Code generated by go-swagger; DO NOT EDIT.

package document

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"

	"platform/services/docdb/models"
)

// DeleteDocumentHandlerFunc turns a function with the right signature into a delete document handler
type DeleteDocumentHandlerFunc func(DeleteDocumentParams, *models.Principal) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteDocumentHandlerFunc) Handle(params DeleteDocumentParams, principal *models.Principal) middleware.Responder {
	return fn(params, principal)
}

// DeleteDocumentHandler interface for that can handle valid delete document params
type DeleteDocumentHandler interface {
	Handle(DeleteDocumentParams, *models.Principal) middleware.Responder
}

// NewDeleteDocument creates a new http.Handler for the delete document operation
func NewDeleteDocument(ctx *middleware.Context, handler DeleteDocumentHandler) *DeleteDocument {
	return &DeleteDocument{Context: ctx, Handler: handler}
}

/*DeleteDocument swagger:route DELETE /{keyspaceName}/{tableName}/document document deleteDocument

Deletes document

*/
type DeleteDocument struct {
	Context *middleware.Context
	Handler DeleteDocumentHandler
}

func (o *DeleteDocument) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDeleteDocumentParams()

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
