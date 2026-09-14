package dto

import (
	"encoding/json"
	"time"

	"github.com/studentinovisad/popisomator/backend/internal/repository"
)

type ItemRequest struct {
	UserID    int64     `json:"user_id"`
	ItemID    int64     `json:"item_id"`
	CreatedAt time.Time `json:"created_at"`
	Status    string    `json:"status"`
	Reason    string    `json:"reason"`
	UserName  string    `json:"user_name,omitempty"`
	ItemName  string    `json:"item_name,omitempty"`
}

type ItemRequestsPage struct {
	// ItemRequest page items
	Items  []ItemRequestSummary `json:"items"`
	Limit  int32                `json:"limit"`
	Offset int32                `json:"offset"`
	Total  int64                `json:"total"`
}

type ItemRequestSummary struct {
	ItemRequest
	UserName string `json:"user_name"`
	ItemName string `json:"item_name"`
}

// ItemRequestUserOption is a requester available in the manager filter.
type ItemRequestUserOption struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type ItemRequestPreparationReport struct {
	User  ItemRequestUserOption        `json:"user"`
	Items []ItemRequestPreparationItem `json:"items"`
}

type ItemRequestPreparationReportRequest struct {
	UserID      int64   `validate:"required,gt=0"`
	ItemIDs     []int64 `validate:"omitempty,dive,gt=0"`
	Status      *string `validate:"omitempty,oneof=requested approved"`
	CreatedFrom *time.Time
	CreatedTo   *time.Time
}

type ItemRequestPreparationItem struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	TypeName string `json:"type_name"`
	// Returns name of the location where the item is located in,
	// as well as the names of it's ancestor locations, starting with
	// the root location.
	LocationNames     []string                         `json:"location_names,omitempty"`
	DerivedNameFormat string                           `json:"derived_name_format"`
	Status            repository.RequestStatus         `json:"status"`
	Consumption       repository.ConsumptionStatus     `json:"consumption"`
	Reason            string                           `json:"reason"`
	RequestedAt       time.Time                        `json:"requested_at"`
	Properties        []ItemRequestPreparationProperty `json:"properties"`
}

type ItemRequestPreparationProperty struct {
	Name       string                        `json:"name"`
	Value      json.RawMessage               `json:"value"`
	ValueType  string                        `json:"value_type"`
	Visibility repository.PropertyVisibility `json:"visibility"`
	Position   int32                         `json:"position"`
}

type ItemRequestCreateRequest struct {
	UserID int64  `json:"user_id" validate:"required,gt=0"`
	ItemID int64  `json:"item_id" validate:"required,gt=0"`
	Reason string `json:"reason" validate:"max=400"`
}

type ItemRequestCreatePersonalRequest struct {
	ItemID int64  `json:"item_id"`
	Reason string `json:"reason"`
}

type ItemRequestIdentifierRequest struct {
	UserID int64 `json:"user_id" validate:"required,gt=0"`
	ItemID int64 `json:"item_id" validate:"required,gt=0"`
}

type ItemRequestsListRequest struct {
	Limit       int32
	Offset      int32
	Status      *string `validate:"omitempty,oneof=requested approved"`
	UserIDs     []int64 `validate:"omitempty,dive,gt=0"`
	CreatedFrom *time.Time
	CreatedTo   *time.Time
}

func ToItemRequestDTO(request repository.ItemRequest) ItemRequest {
	return ItemRequest{
		UserID:    request.UserID,
		ItemID:    request.ItemID,
		CreatedAt: request.CreatedAt.Time,
		Status:    string(request.Status),
		Reason:    request.Reason,
	}
}

func ToItemRequestSummaryDTO(request repository.ListItemRequestsRow) ItemRequestSummary {
	return ItemRequestSummary{
		ItemRequest: ToItemRequestDTO(repository.ItemRequest{
			UserID: request.UserID, ItemID: request.ItemID, CreatedAt: request.CreatedAt,
			Status: request.Status, Reason: request.Reason,
		}),
		UserName: request.UserName,
		ItemName: request.ItemName,
	}
}
