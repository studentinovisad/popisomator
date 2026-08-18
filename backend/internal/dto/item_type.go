package dto

import "github.com/studentinovisad/popisomator/backend/internal/repository"

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

// Minified ItemType object for things such as dropdown lists
type ItemTypeOption struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func ToItemTypeOptionDTO(itemType repository.ItemType) ItemTypeOption {
	return ItemTypeOption{
		ID:   itemType.ID,
		Name: itemType.Name,
	}
}

// Property added to an item type
type ItemTypeProperty struct {
	ID           int64   `json:"id" validate:"required"`
	DefaultValue *string `json:"default_value"`
	Visibility   string  `json:"visibility" validate:"omitempty,oneof=overview details"`
}

func ToItemTypePropertyDTO(itemTypeProp repository.ItemTypeProperty) ItemTypeProperty {
	return ItemTypeProperty{
		ID:           itemTypeProp.PropertyID,
		DefaultValue: itemTypeProp.DefaultValue,
		Visibility:   string(itemTypeProp.Visibility),
	}
}

type ItemTypesPage struct {
	// ItemType page items
	Items  []ItemType `json:"items"`
	Limit  int32      `json:"limit"`
	Offset int32      `json:"offset"`
	Total  int64      `json:"total"`
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

type AddUpdateItemTypePropertyRequest struct {
	TypeID       int64                          `json:"type_id" validate:"required"`
	PropertyID   int64                          `json:"property_id" validate:"required"`
	DefaultValue *string                        `json:"default_value"`
	Visibility   *repository.PropertyVisibility `json:"visibility" validate:"omitempty,oneof=overview details"`
}
