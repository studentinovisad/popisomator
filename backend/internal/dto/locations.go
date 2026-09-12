package dto

import (
	"encoding/json"

	"github.com/studentinovisad/popisomator/backend/internal/repository"
)

type LocationOption struct {
	ID       int64             `json:"id"`
	Name     string            `json:"name"`
	Children []*LocationOption `json:"children"`
}

type Location struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	ParentID    *int64  `json:"parent_id,omitempty"`
}

func ToLocationDTO(location repository.Location) Location {
	locationDTO := Location{
		ID:          location.ID,
		Name:        location.Name,
		Description: &location.Description.String,
	}
	if location.ParentID.Valid {
		locationDTO.ParentID = &location.ParentID.Int64
	}
	return locationDTO
}

type CreateLocationRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	ParentID    *int64 `json:"parent_id"`
}

type UpdateLocationRequest struct {
	ID          int64   `json:"id" validate:"required"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	ParentID    *int64  `json:"parent_id"`
	ParentIDSet bool    `json:"-"`
}

func (r *UpdateLocationRequest) UnmarshalJSON(data []byte) error {
	type requestAlias UpdateLocationRequest
	var request requestAlias
	if err := json.Unmarshal(data, &request); err != nil {
		return err
	}

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	*r = UpdateLocationRequest(request)
	_, r.ParentIDSet = fields["parent_id"]
	return nil
}
