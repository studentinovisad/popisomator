package dto

import (
	"encoding/json"

	"github.com/studentinovisad/popisomator/backend/internal/repository"
)

type Property struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	ValueType    string  `json:"value_type"`
	DefaultValue *string `json:"default_value"`
}

func ToPropertyDTO(prop repository.Property) Property {
	return Property{
		ID:           prop.ID,
		Name:         prop.Name,
		Description:  prop.Description.String,
		ValueType:    prop.ValueType,
		DefaultValue: prop.DefaultValue,
	}
}

type PropertyOption struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	ValueType    string  `json:"value_type"`
	DefaultValue *string `json:"default_value"`
}

type PropertiesPage struct {
	// Property page items
	Items  []Property `json:"items"`
	Limit  int32      `json:"limit"`
	Offset int32      `json:"offset"`
	Total  int64      `json:"total"`
}

type CreatePropertyRequest struct {
	Name         string  `json:"name" validate:"required"`
	Description  string  `json:"description"`
	ValueType    string  `json:"value_type" validate:"required,oneof=string number boolean"`
	DefaultValue *string `json:"default_value"`
}

// PropertyValueCheck is an internal validation carrier (not a request/response DTO) used to
// run a raw property value through the "valuetype" tag against a known ValueType.
type PropertyValueCheck struct {
	Value     string `validate:"required,valuetype"`
	ValueType string
}

type UpdatePropertyRequest struct {
	ID              int64   `json:"id" validate:"required"`
	Name            *string `json:"name"`
	Description     *string `json:"description"`
	DefaultValue    *string `json:"default_value"`
	DefaultValueSet bool    `json:"-"`
}

func (r *UpdatePropertyRequest) UnmarshalJSON(data []byte) error {
	type requestAlias UpdatePropertyRequest
	var request requestAlias
	if err := json.Unmarshal(data, &request); err != nil {
		return err
	}

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	*r = UpdatePropertyRequest(request)
	_, r.DefaultValueSet = fields["default_value"]
	return nil
}
