// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// RelationOperator relation operator
//
// swagger:model RelationOperator
type RelationOperator string

const (

	// RelationOperatorEq captures enum value "eq"
	RelationOperatorEq RelationOperator = "eq"

	// RelationOperatorLt captures enum value "lt"
	RelationOperatorLt RelationOperator = "lt"

	// RelationOperatorLeq captures enum value "leq"
	RelationOperatorLeq RelationOperator = "leq"

	// RelationOperatorGt captures enum value "gt"
	RelationOperatorGt RelationOperator = "gt"

	// RelationOperatorGeq captures enum value "geq"
	RelationOperatorGeq RelationOperator = "geq"

	// RelationOperatorIn captures enum value "in"
	RelationOperatorIn RelationOperator = "in"

	// RelationOperatorCnt captures enum value "cnt"
	RelationOperatorCnt RelationOperator = "cnt"

	// RelationOperatorCntKey captures enum value "cntKey"
	RelationOperatorCntKey RelationOperator = "cntKey"

	// RelationOperatorLike captures enum value "like"
	RelationOperatorLike RelationOperator = "like"
)

// for schema
var relationOperatorEnum []interface{}

func init() {
	var res []RelationOperator
	if err := json.Unmarshal([]byte(`["eq","lt","leq","gt","geq","in","cnt","cntKey","like"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		relationOperatorEnum = append(relationOperatorEnum, v)
	}
}

func (m RelationOperator) validateRelationOperatorEnum(path, location string, value RelationOperator) error {
	if err := validate.EnumCase(path, location, value, relationOperatorEnum, true); err != nil {
		return err
	}
	return nil
}

// Validate validates this relation operator
func (m RelationOperator) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateRelationOperatorEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
