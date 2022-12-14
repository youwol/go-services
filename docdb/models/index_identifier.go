// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// IndexIdentifier Decription of a table index identifier. If using a partition key, this will refer to a local secondary index, else it will be a global secondary index
//
// swagger:model IndexIdentifier
type IndexIdentifier struct {

	// column name
	// Required: true
	ColumnName *string `json:"column_name"`

	// option
	// Enum: [ keys values entries full]
	Option string `json:"option,omitempty"`

	// partition key
	PartitionKey string `json:"partition_key,omitempty"`
}

// Validate validates this index identifier
func (m *IndexIdentifier) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateColumnName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateOption(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *IndexIdentifier) validateColumnName(formats strfmt.Registry) error {

	if err := validate.Required("column_name", "body", m.ColumnName); err != nil {
		return err
	}

	return nil
}

var indexIdentifierTypeOptionPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["","keys","values","entries","full"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		indexIdentifierTypeOptionPropEnum = append(indexIdentifierTypeOptionPropEnum, v)
	}
}

const (

	// IndexIdentifierOptionEmpty captures enum value ""
	IndexIdentifierOptionEmpty string = ""

	// IndexIdentifierOptionKeys captures enum value "keys"
	IndexIdentifierOptionKeys string = "keys"

	// IndexIdentifierOptionValues captures enum value "values"
	IndexIdentifierOptionValues string = "values"

	// IndexIdentifierOptionEntries captures enum value "entries"
	IndexIdentifierOptionEntries string = "entries"

	// IndexIdentifierOptionFull captures enum value "full"
	IndexIdentifierOptionFull string = "full"
)

// prop value enum
func (m *IndexIdentifier) validateOptionEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, indexIdentifierTypeOptionPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *IndexIdentifier) validateOption(formats strfmt.Registry) error {

	if swag.IsZero(m.Option) { // not required
		return nil
	}

	// value enum
	if err := m.validateOptionEnum("option", "body", m.Option); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *IndexIdentifier) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *IndexIdentifier) UnmarshalBinary(b []byte) error {
	var res IndexIdentifier
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
