package dto

import (
	"github.com/studentinovisad/popisomator/backend/internal/repository"
)

type ItemType struct {
	ID          int64              `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Properties  []ItemTypeProperty `json:"properties"`
}

func ToItemTypeDTO(itemType repository.ItemType) ItemType {
	return ItemType{
		ID:          itemType.ID,
		Name:        itemType.Name,
		Description: itemType.Description.String,
		Properties:  make([]ItemTypeProperty, 0),
	}
}

type Item struct {
	ID          int64                        `json:"id"`
	Consumption repository.ConsumptionStatus `json:"consumption"`
	Properties  []ItemProperty               `json:"properties"`
	TypeID      *int64                       `json:"type_id"`
}

func ToItemDTO(item repository.Item) Item {
	var type_id *int64
	if item.TypeID.Valid {
		type_id = &item.TypeID.Int64
	}
	return Item{
		ID:          item.ID,
		Consumption: item.Consumption,
		Properties:  make([]ItemProperty, 0),
		TypeID:      type_id,
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

// Property added to an item type
type ItemTypeProperty struct {
	ID           int64   `json:"id"`
	DefaultValue *string `json:"default_value"`
}

func ToItemTypePropertyDTO(itemTypeProp repository.ItemTypeProperty) ItemTypeProperty {
	return ItemTypeProperty{
		ID:           itemTypeProp.PropertyID,
		DefaultValue: itemTypeProp.DefaultValue,
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

type CreateItemRequest struct {
	Properties []ItemProperty `json:"properties"`
	TypeID     *int64         `json:"type_id"`
}

type CreateItemTypeRequest struct {
	Name        string             `json:"name" validate:"required"`
	Description string             `json:"description"`
	Properties  []ItemTypeProperty `json:"properties"`
}

type CreatePropertyRequest struct {
	Name         string  `json:"name" validate:"required"`
	Description  string  `json:"description"`
	ValueType    string  `json:"value_type" validate:"required"`
	DefaultValue *string `json:"default_value"`
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

type AddUpdateItemTypePropertyRequest struct {
	TypeID       int64   `json:"type_id" validate:"required"`
	PropertyID   int64   `json:"property_id" validate:"required"`
	DefaultValue *string `json:"default_value"`
}
