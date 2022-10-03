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

// Replication replication
//
// swagger:model Replication
type Replication struct {

	// class
	// Required: true
	// Enum: [SimpleStrategy NetworkTopologyStrategy]
	Class *string `json:"class"`

	// replication factor
	ReplicationFactor *int64 `json:"replication_factor,omitempty"`
}

// Validate validates this replication
func (m *Replication) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateClass(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var replicationTypeClassPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["SimpleStrategy","NetworkTopologyStrategy"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		replicationTypeClassPropEnum = append(replicationTypeClassPropEnum, v)
	}
}

const (

	// ReplicationClassSimpleStrategy captures enum value "SimpleStrategy"
	ReplicationClassSimpleStrategy string = "SimpleStrategy"

	// ReplicationClassNetworkTopologyStrategy captures enum value "NetworkTopologyStrategy"
	ReplicationClassNetworkTopologyStrategy string = "NetworkTopologyStrategy"
)

// prop value enum
func (m *Replication) validateClassEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, replicationTypeClassPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *Replication) validateClass(formats strfmt.Registry) error {

	if err := validate.Required("class", "body", m.Class); err != nil {
		return err
	}

	// value enum
	if err := m.validateClassEnum("class", "body", *m.Class); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Replication) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Replication) UnmarshalBinary(b []byte) error {
	var res Replication
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
