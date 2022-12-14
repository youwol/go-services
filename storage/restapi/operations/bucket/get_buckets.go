// Code generated by go-swagger; DO NOT EDIT.

package bucket

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"

	"platform/services/storage/models"
)

// GetBucketsHandlerFunc turns a function with the right signature into a get buckets handler
type GetBucketsHandlerFunc func(GetBucketsParams, *models.Principal) middleware.Responder

// Handle executing the request and returning a response
func (fn GetBucketsHandlerFunc) Handle(params GetBucketsParams, principal *models.Principal) middleware.Responder {
	return fn(params, principal)
}

// GetBucketsHandler interface for that can handle valid get buckets params
type GetBucketsHandler interface {
	Handle(GetBucketsParams, *models.Principal) middleware.Responder
}

// NewGetBuckets creates a new http.Handler for the get buckets operation
func NewGetBuckets(ctx *middleware.Context, handler GetBucketsHandler) *GetBuckets {
	return &GetBuckets{Context: ctx, Handler: handler}
}

/*GetBuckets swagger:route GET /buckets bucket getBuckets

Get the list of buckets

*/
type GetBuckets struct {
	Context *middleware.Context
	Handler GetBucketsHandler
}

func (o *GetBuckets) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetBucketsParams()

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
