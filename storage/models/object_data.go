// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// ObjectData object data
//
// swagger:model ObjectData
type ObjectData struct {

	// base64 encoded file data
	// Format: byte
	Data strfmt.Base64 `json:"data,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// size
	Size int64 `json:"size,omitempty"`
}

// Validate validates this object data
func (m *ObjectData) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ObjectData) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ObjectData) UnmarshalBinary(b []byte) error {
	var res ObjectData
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
