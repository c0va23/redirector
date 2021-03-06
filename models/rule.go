package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// Rule rule
// swagger:model Rule
type Rule struct {

	// active from
	// Format: date-time
	ActiveFrom *strfmt.DateTime `json:"activeFrom,omitempty"`

	// active to
	// Format: date-time
	ActiveTo *strfmt.DateTime `json:"activeTo,omitempty"`

	// Rule resolver
	// Required: true
	// Enum: [simple pattern]
	Resolver string `json:"resolver"`

	// Regex for match source path
	// Required: true
	SourcePath string `json:"sourcePath"`

	// target
	// Required: true
	Target Target `json:"target"`
}

// Validate validates this rule
func (m *Rule) Validate(formats strfmt.Registry) error {
	return nil
}

const (

	// RuleResolverSimple captures enum value "simple"
	RuleResolverSimple string = "simple"

	// RuleResolverPattern captures enum value "pattern"
	RuleResolverPattern string = "pattern"
)

// MarshalBinary interface implementation
func (m *Rule) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Rule) UnmarshalBinary(b []byte) error {
	var res Rule
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
