package dto

import (
	"encoding/json"
	"time"

	"github.com/studentinovisad/popisomator/backend/internal/repository"
)

type Item struct {
	ID            int64                        `json:"id"`
	Consumption   repository.ConsumptionStatus `json:"consumption"`
	Properties    []ItemProperty               `json:"properties"`
	TypeID        int64                        `json:"type_id"`
	DerivedName   string                       `json:"derived_name,omitempty"`
	RequestStatus *repository.RequestStatus    `json:"request_status,omitempty"`
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
	ID         int64           `json:"id" validate:"required"`
	Value      json.RawMessage `json:"value" validate:"required"`
	ValueType  string          `json:"value_type,omitempty"`
	Visibility string          `json:"visibility,omitempty"`
	SmartData  string          `json:"smart_data,omitempty"`
}

func ToItemPropertyDTO(itemProp repository.ItemProperty) ItemProperty {
	return ItemProperty{
		ID:    itemProp.PropertyID,
		Value: itemProp.PropertyValue,
	}
}

type ListItemsRequest struct {
	TypeID          *int64
	PropertyFilters map[int64]json.RawMessage
	Consumption     []repository.ConsumptionStatus `validate:"omitempty,dive,oneof=not_consumed partially_consumed fully_consumed damaged"`
	CreatedFrom     *time.Time
	CreatedTo       *time.Time
	Limit           int32
	Offset          int32
	Order           string `validate:"oneof=asc desc"`
	Search          string `validate:"max=100"`
	ViewerID        int64
}

type ItemsPage struct {
	// Item page items
	Items  []Item `json:"items"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
	Total  int64  `json:"total"`
}

type CreateItemRequest struct {
	Properties []ItemProperty `json:"properties" validate:"dive"`
	TypeID     int64          `json:"type_id" validate:"required,gt=0"`
	Amount     int32          `json:"amount" validate:"required,gt=0,lte=100"`
}

type UpdateItemRequest struct {
	ID          int64   `json:"id" validate:"required"`
	TypeID      *int64  `json:"type_id" validate:"omitempty,gt=0"`
	Consumption *string `json:"consumption" validate:"omitempty,oneof=not_consumed partially_consumed fully_consumed damaged"`
	ViewerID    int64   `json:"-"`
}

type AddUpdateItemPropertyRequest struct {
	ItemID     int64           `json:"item_id" validate:"required"`
	PropertyID int64           `json:"property_id" validate:"required"`
	Value      json.RawMessage `json:"value" validate:"required"`
}
