package dto

import (
	"encoding/json"
	"time"

	"github.com/studentinovisad/popisomator/backend/internal/repository"
)

type ItemType struct {
	ID                int64              `json:"id"`
	Name              string             `json:"name"`
	Description       string             `json:"description"`
	Properties        []ItemTypeProperty `json:"properties"`
	DerivedNameFormat string             `json:"derived_name_format"`
}

func ToItemTypeDTO(itemType repository.ItemType) ItemType {
	return ItemType{
		ID:                itemType.ID,
		Name:              itemType.Name,
		Description:       itemType.Description.String,
		Properties:        make([]ItemTypeProperty, 0),
		DerivedNameFormat: itemType.DerivedNameFormat.String,
	}
}

type ItemTypeOption struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Item struct {
	ID          int64                        `json:"id"`
	Consumption repository.ConsumptionStatus `json:"consumption"`
	Properties  []ItemProperty               `json:"properties"`
	TypeID      int64                        `json:"type_id"`
	DerivedName string                       `json:"derived_name,omitempty"`
}

func ToItemDTO(item repository.Item) Item {
	return Item{
		ID:          item.ID,
		Consumption: item.Consumption,
		Properties:  make([]ItemProperty, 0),
		TypeID:      item.TypeID,
	}
}

// Property added to an item
type ItemProperty struct {
	ID         int64  `json:"id" validate:"required"`
	Value      string `json:"value" validate:"required"`
	Visibility string `json:"visibility,omitempty"`
}

func ToItemPropertyDTO(itemProp repository.ItemProperty) ItemProperty {
	return ItemProperty{
		ID:    itemProp.PropertyID,
		Value: itemProp.PropertyValue,
	}
}

func RowToItemPropertyDTO(propRow repository.GetItemPropertiesForItemsRow) ItemProperty {
	return ItemProperty{
		ID:         propRow.PropertyID,
		Value:      propRow.PropertyValue,
		Visibility: string(propRow.Visibility),
	}
}

// Property added to an item type
type ItemTypeProperty struct {
	ID           int64   `json:"id" validate:"required"`
	DefaultValue *string `json:"default_value"`
	Name         string  `json:"name,omitempty"`
	Visibility   string  `json:"visibility" validate:"omitempty,oneof=overview details"`
}

func ToItemTypePropertyDTO(itemTypeProp repository.ItemTypeProperty) ItemTypeProperty {
	return ItemTypeProperty{
		ID:           itemTypeProp.PropertyID,
		DefaultValue: itemTypeProp.DefaultValue,
		Visibility:   string(itemTypeProp.Visibility),
	}
}

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

type ListItemsRequest struct {
	TypeID      *int64
	Consumption []repository.ConsumptionStatus `validate:"omitempty,dive,oneof=not_consumed partially_consumed fully_consumed damaged"`
	CreatedFrom *time.Time
	CreatedTo   *time.Time
	Limit       int32
	Offset      int32
	Order       string `validate:"oneof=asc desc"`
}

type ItemsPage struct {
	Items  []Item `json:"items"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
	Total  int64  `json:"total"`
}

type ItemTypesPage struct {
	Items  []ItemType `json:"items"`
	Limit  int32      `json:"limit"`
	Offset int32      `json:"offset"`
	Total  int64      `json:"total"`
}

type PropertiesPage struct {
	Items  []Property `json:"items"`
	Limit  int32      `json:"limit"`
	Offset int32      `json:"offset"`
	Total  int64      `json:"total"`
}

type CreateItemRequest struct {
	Properties []ItemProperty `json:"properties" validate:"dive"`
	TypeID     int64          `json:"type_id" validate:"required,gt=0"`
	Amount     int32          `json:"amount" validate:"required,gt=0,lte=100"`
}

type ConsumeItemRequest struct {
	ID     int64                         `json:"id" validate:"required"`
	Status *repository.ConsumptionStatus `json:"status"`
}

type CreateItemTypeRequest struct {
	Name              string             `json:"name" validate:"required"`
	Description       string             `json:"description"`
	Properties        []ItemTypeProperty `json:"properties" validate:"dive"`
	DerivedNameFormat string             `json:"derived_name_format"`
}

type UpdateItemTypeRequest struct {
	ID                int64   `json:"id" validate:"required"`
	Name              *string `json:"name"`
	Description       *string `json:"description"`
	DerivedNameFormat *string `json:"derived_name_format"`
}

type SetItemTypeRequest struct {
	ID     int64 `json:"id" validate:"required"`
	TypeID int64 `json:"type_id" validate:"required,gt=0"`
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

type AddUpdateItemPropertyRequest struct {
	ItemID     int64  `json:"item_id" validate:"required"`
	PropertyID int64  `json:"property_id" validate:"required"`
	Value      string `json:"value" validate:"required"`
}

type AddUpdateItemTypePropertyRequest struct {
	TypeID       int64                          `json:"type_id" validate:"required"`
	PropertyID   int64                          `json:"property_id" validate:"required"`
	DefaultValue *string                        `json:"default_value"`
	Visibility   *repository.PropertyVisibility `json:"visibility" validate:"omitempty,oneof=overview details"`
}
