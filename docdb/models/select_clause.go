// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// SelectClause select clause
//
// swagger:model SelectClause
type SelectClause struct {

	// identifier
	Identifier string `json:"identifier,omitempty"`

	// selector
	Selector string `json:"selector,omitempty"`
}

// Validate validates this select clause
func (m *SelectClause) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *SelectClause) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SelectClause) UnmarshalBinary(b []byte) error {
	var res SelectClause
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
