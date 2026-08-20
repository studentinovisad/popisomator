package dto

import (
	"time"

	"github.com/studentinovisad/popisomator/backend/internal/repository"
)

type ItemRequest struct {
	UserID    int64     `json:"user_id"`
	ItemID    int64     `json:"item_id"`
	CreatedAt time.Time `json:"created_at"`
	Status    string    `json:"status"`
	Reason    string    `json:"reason"`
}

type ItemRequestsPage struct {
	// ItemRequest page items
	Items  []ItemRequest `json:"items"`
	Limit  int32         `json:"limit"`
	Offset int32         `json:"offset"`
	Total  int64         `json:"total"`
}

type ItemRequestCreateRequest struct {
	UserID int64  `json:"user_id" validate:"required,gt=0"`
	ItemID int64  `json:"item_id" validate:"required,gt=0"`
	Reason string `json:"reason" validate:"required,max=400"`
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
	Limit  int32
	Offset int32
	Status *string `validate:"omitempty,oneof=requested approved"`
	UserID *int64  `validate:"omitempty,gt=0"`
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
