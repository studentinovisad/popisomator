package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/studentinovisad/popisomator/backend/internal/db"
	"github.com/studentinovisad/popisomator/backend/internal/dto"
	"github.com/studentinovisad/popisomator/backend/internal/repository"
)

func GetLocationOptionsFlat(ctx context.Context) ([]dto.LocationOption, error) {
	locations, err := db.Queries.ListLocationOptions(ctx)
	if err != nil {
		return nil, err
	}

	locationOptionsDTO := make([]dto.LocationOption, len(locations))
	for i, location := range locations {
		locationOptionsDTO[i] = dto.LocationOption{
			ID:       location.ID,
			Name:     location.Name,
			Children: nil,
		}
	}

	return locationOptionsDTO, nil
}

func GetLocationOptions(ctx context.Context) ([]*dto.LocationOption, error) {
	locations, err := db.Queries.ListLocationOptions(ctx)
	if err != nil {
		return nil, err
	}

	// First pass: gather all locations as location options
	allOptionsDTO := make(map[int64]*dto.LocationOption, len(locations))
	for _, location := range locations {
		allOptionsDTO[location.ID] = &dto.LocationOption{
			ID:   location.ID,
			Name: location.Name,
		}
	}

	// Second pass: put location options into children arrays, fill final array
	locationOptionsDTO := make([]*dto.LocationOption, 0)
	for _, location := range locations {
		if location.ParentID.Valid {
			parentOption := allOptionsDTO[location.ParentID.Int64]
			parentOption.Children = append(parentOption.Children, allOptionsDTO[location.ID])
		} else {
			locationOptionsDTO = append(locationOptionsDTO, allOptionsDTO[location.ID])
		}
	}

	return locationOptionsDTO, nil
}

func GetLocationByID(ctx context.Context, id int64) (dto.Location, error) {
	location, err := db.Queries.GetLocationByID(ctx, id)
	if err != nil {
		return dto.Location{}, err
	}

	locationDTO := dto.ToLocationDTO(location)

	return locationDTO, nil
}

func CreateLocation(ctx context.Context, req dto.CreateLocationRequest) (dto.Location, error) {
	if err := dto.Validate(req); err != nil {
		return dto.Location{}, err
	}

	description := pgtype.Text{String: "", Valid: false}
	if len(req.Description) > 0 {
		description = pgtype.Text{String: req.Description, Valid: true}
	}
	parentID := pgtype.Int8{Valid: false}
	if req.ParentID != nil {
		parentID = pgtype.Int8{Int64: *req.ParentID, Valid: true}
	}

	location, err := db.Queries.CreateLocation(ctx, repository.CreateLocationParams{
		Name:        req.Name,
		Description: description,
		ParentID:    parentID,
	})
	if err != nil {
		return dto.Location{}, err
	}

	locationDTO := dto.ToLocationDTO(location)

	return locationDTO, nil
}

func UpdateLocation(ctx context.Context, req dto.UpdateLocationRequest) (dto.Location, error) {
	if err := dto.Validate(req); err != nil {
		return dto.Location{}, err
	}

	if req.Name != nil {
		if err := db.Queries.UpdateLocation_Name(ctx, repository.UpdateLocation_NameParams{
			ID:   req.ID,
			Name: *req.Name,
		}); err != nil {
			return dto.Location{}, err
		}
	}

	if req.Description != nil {
		description := pgtype.Text{String: *req.Description, Valid: true}

		if err := db.Queries.UpdateLocation_Description(ctx, repository.UpdateLocation_DescriptionParams{
			ID:          req.ID,
			Description: description,
		}); err != nil {
			return dto.Location{}, err
		}
	}

	if req.ParentIDSet {
		parentID := pgtype.Int8{Valid: false}
		if req.ParentID != nil {
			res, err := db.Queries.CheckLocationCycle(ctx, repository.CheckLocationCycleParams{
				LocationID: req.ID,
				ParentID:   *req.ParentID,
			})
			if err != nil {
				return dto.Location{}, err
			}
			if res == true {
				return dto.Location{}, ErrLocationCycleDetected
			}

			parentID = pgtype.Int8{Int64: *req.ParentID, Valid: true}
		}

		if err := db.Queries.UpdateLocation_ParentID(ctx, repository.UpdateLocation_ParentIDParams{
			ID:       req.ID,
			ParentID: parentID,
		}); err != nil {
			return dto.Location{}, err
		}
	}

	location, err := db.Queries.GetLocationByID(ctx, req.ID)
	if err != nil {
		return dto.Location{}, err
	}

	locationDTO := dto.ToLocationDTO(location)

	return locationDTO, nil
}

func DeleteLocation(ctx context.Context, id int64) error {
	rowsAffected, err := db.Queries.DeleteLocation(ctx, id)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
