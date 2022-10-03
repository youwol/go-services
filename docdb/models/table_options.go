// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// TableOptions @see https://docs.scylladb.com/getting-started/ddl/#create-table-statement
//
// swagger:model TableOptions
type TableOptions struct {

	// bloom filter fp chance
	BloomFilterFpChance *float64 `json:"bloom_filter_fp_chance,omitempty"`

	// clustering order
	ClusteringOrder ClusteringOrder `json:"clustering_order,omitempty"`

	// comment
	Comment string `json:"comment,omitempty"`

	// compaction
	Compaction *CompactionOptions `json:"compaction,omitempty"`

	// compression
	Compression *CompressionOptions `json:"compression,omitempty"`

	// dclocal read repair chance
	DclocalReadRepairChance *float64 `json:"dclocal_read_repair_chance,omitempty"`

	// default time to live
	DefaultTimeToLive int64 `json:"default_time_to_live,omitempty"`

	// gc grace seconds
	GcGraceSeconds *int64 `json:"gc_grace_seconds,omitempty"`

	// memtable flush period in ms
	MemtableFlushPeriodInMs int64 `json:"memtable_flush_period_in_ms,omitempty"`

	// read repair chance
	ReadRepairChance int64 `json:"read_repair_chance,omitempty"`

	// @see https://docs.scylladb.com/getting-started/ddl/#speculative-retry-options
	SpeculativeRetry *string `json:"speculative_retry,omitempty"`
}

// Validate validates this table options
func (m *TableOptions) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateClusteringOrder(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCompaction(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCompression(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *TableOptions) validateClusteringOrder(formats strfmt.Registry) error {

	if swag.IsZero(m.ClusteringOrder) { // not required
		return nil
	}

	if err := m.ClusteringOrder.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("clustering_order")
		}
		return err
	}

	return nil
}

func (m *TableOptions) validateCompaction(formats strfmt.Registry) error {

	if swag.IsZero(m.Compaction) { // not required
		return nil
	}

	if m.Compaction != nil {
		if err := m.Compaction.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("compaction")
			}
			return err
		}
	}

	return nil
}

func (m *TableOptions) validateCompression(formats strfmt.Registry) error {

	if swag.IsZero(m.Compression) { // not required
		return nil
	}

	if m.Compression != nil {
		if err := m.Compression.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("compression")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *TableOptions) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *TableOptions) UnmarshalBinary(b []byte) error {
	var res TableOptions
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
