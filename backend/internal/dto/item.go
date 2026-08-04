package dto

import (
	"github.com/studentinovisad/popisomator/backend/internal/repository"
)

type ItemType struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Item struct {
	ID          int64                        `json:"id"`
	Consumption repository.ConsumptionStatus `json:"consumption"`
	Properties  []ItemProperty               `json:"properties"`
}

func ToItemDTO(item repository.Item) Item {
	return Item{
		ID:          item.ID,
		Consumption: item.Consumption,
		Properties:  make([]ItemProperty, 0),
	}
}

// Property added to an item
type ItemProperty struct {
	ID    int64  `json:"id"`
	Value string `json:"value"`
}

func ToItemPropertyDTO(itemProp repository.ItemProperty) ItemProperty {
	return ItemProperty{
		ID:    itemProp.PropertyID,
		Value: itemProp.PropertyValue,
	}
}

type Property struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	ValueType    string `json:"value_type"`
	DefaultValue string `json:"default_value"`
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

type CreateItemRequest struct {
	Properties []ItemProperty `json:"properties"`
}

type CreatePropertyRequest struct {
	Name         string `json:"name" validate:"required"`
	Description  string `json:"description"`
	ValueType    string `json:"value_type" validate:"required"`
	DefaultValue string `json:"default_value"`
}

type UpdatePropertyRequest struct {
	ID           int64   `json:"id" validate:"required"`
	Name         *string `json:"name"`
	Description  *string `json:"description"`
	DefaultValue *string `json:"default_value"`
}

type AddUpdateItemPropertyRequest struct {
	ItemID     int64  `json:"item_id" validate:"required"`
	PropertyID int64  `json:"property_id" validate:"required"`
	Value      string `json:"value" validate:"required"`
}
