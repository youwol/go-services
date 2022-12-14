// Code generated by go-swagger; DO NOT EDIT.

package bucket

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"errors"
	"net/url"
	golangswaggerpaths "path"
	"strings"

	"github.com/go-openapi/swag"
)

// DeleteBucketURL generates an URL for the delete bucket operation
type DeleteBucketURL struct {
	BucketName string

	ForceNotEmpty *bool

	_basePath string
	// avoid unkeyed usage
	_ struct{}
}

// WithBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *DeleteBucketURL) WithBasePath(bp string) *DeleteBucketURL {
	o.SetBasePath(bp)
	return o
}

// SetBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *DeleteBucketURL) SetBasePath(bp string) {
	o._basePath = bp
}

// Build a url path and query string
func (o *DeleteBucketURL) Build() (*url.URL, error) {
	var _result url.URL

	var _path = "/bucket/{bucketName}"

	bucketName := o.BucketName
	if bucketName != "" {
		_path = strings.Replace(_path, "{bucketName}", bucketName, -1)
	} else {
		return nil, errors.New("bucketName is required on DeleteBucketURL")
	}

	_basePath := o._basePath
	if _basePath == "" {
		_basePath = "/api/v0-alpha1"
	}
	_result.Path = golangswaggerpaths.Join(_basePath, _path)

	qs := make(url.Values)

	var forceNotEmptyQ string
	if o.ForceNotEmpty != nil {
		forceNotEmptyQ = swag.FormatBool(*o.ForceNotEmpty)
	}
	if forceNotEmptyQ != "" {
		qs.Set("forceNotEmpty", forceNotEmptyQ)
	}

	_result.RawQuery = qs.Encode()

	return &_result, nil
}

// Must is a helper function to panic when the url builder returns an error
func (o *DeleteBucketURL) Must(u *url.URL, err error) *url.URL {
	if err != nil {
		panic(err)
	}
	if u == nil {
		panic("url can't be nil")
	}
	return u
}

// String returns the string representation of the path with query string
func (o *DeleteBucketURL) String() string {
	return o.Must(o.Build()).String()
}

// BuildFull builds a full url with scheme, host, path and query string
func (o *DeleteBucketURL) BuildFull(scheme, host string) (*url.URL, error) {
	if scheme == "" {
		return nil, errors.New("scheme is required for a full url on DeleteBucketURL")
	}
	if host == "" {
		return nil, errors.New("host is required for a full url on DeleteBucketURL")
	}

	base, err := o.Build()
	if err != nil {
		return nil, err
	}

	base.Scheme = scheme
	base.Host = host
	return base, nil
}

// StringFull returns the string representation of a complete url
func (o *DeleteBucketURL) StringFull(scheme, host string) string {
	return o.Must(o.BuildFull(scheme, host)).String()
}
