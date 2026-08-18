package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/studentinovisad/popisomator/backend/internal/db"
	"github.com/studentinovisad/popisomator/backend/internal/dto"
	"github.com/studentinovisad/popisomator/backend/internal/repository"
)

// validatePropertyValue looks up the property's declared value_type and checks that rawValue's
// JSON shape matches it. Takes the generated Querier interface so it works identically whether
// called with db.Queries or a transaction's db.Queries.WithTx(tx).
func validatePropertyValue(ctx context.Context, q repository.Querier, propertyID int64, rawValue string) error {
	prop, err := q.GetPropertyByID(ctx, propertyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrInvalidReference
		}
		return err
	}

	return dto.Validate(dto.PropertyValueCheck{Value: rawValue, ValueType: prop.ValueType})
}

func GetAllProperties(ctx context.Context) ([]dto.Property, error) {
	props, err := db.Queries.GetAllProperties(ctx)
	if err != nil {
		return nil, err
	}

	propsDTO := make([]dto.Property, len(props))
	for i, prop := range props {
		propsDTO[i] = dto.ToPropertyDTO(prop)
	}

	return propsDTO, nil
}

func GetPropertyOptions(ctx context.Context) ([]dto.PropertyOption, error) {
	props, err := db.Queries.ListPropertyOptions(ctx)
	if err != nil {
		return nil, err
	}

	propOptionsDTO := make([]dto.PropertyOption, len(props))
	for i, prop := range props {
		propOptionsDTO[i] = dto.PropertyOption{
			ID:           prop.ID,
			Name:         prop.Name,
			ValueType:    prop.ValueType,
			DefaultValue: prop.DefaultValue,
		}
	}

	return propOptionsDTO, nil
}

func ListProperties(ctx context.Context, limit, offset int32) (dto.PropertiesPage, error) {
	total, err := db.Queries.CountProperties(ctx)
	if err != nil {
		return dto.PropertiesPage{}, err
	}

	properties, err := db.Queries.ListProperties(ctx, repository.ListPropertiesParams{
		PageLimit:  limit,
		PageOffset: offset,
	})
	if err != nil {
		return dto.PropertiesPage{}, err
	}

	pageItems := make([]dto.Property, len(properties))
	for index, property := range properties {
		pageItems[index] = dto.ToPropertyDTO(property)
	}

	return dto.PropertiesPage{Items: pageItems, Limit: limit, Offset: offset, Total: total}, nil
}

func GetPropertyByID(ctx context.Context, id int64) (dto.Property, error) {
	prop, err := db.Queries.GetPropertyByID(ctx, id)
	if err != nil {
		return dto.Property{}, err
	}

	propDTO := dto.ToPropertyDTO(prop)

	return propDTO, nil
}

func CreateProperty(ctx context.Context, req dto.CreatePropertyRequest) (dto.Property, error) {
	if err := dto.Validate(req); err != nil {
		return dto.Property{}, err
	}

	if req.DefaultValue != nil {
		if err := dto.Validate(dto.PropertyValueCheck{Value: *req.DefaultValue, ValueType: req.ValueType}); err != nil {
			return dto.Property{}, err
		}
	}

	description := pgtype.Text{String: "", Valid: false}
	if len(req.Description) > 0 {
		description = pgtype.Text{String: req.Description, Valid: true}
	}
	prop, err := db.Queries.CreateProperty(ctx, repository.CreatePropertyParams{
		Name:         req.Name,
		Description:  description,
		ValueType:    req.ValueType,
		DefaultValue: req.DefaultValue,
	})
	if err != nil {
		return dto.Property{}, err
	}

	propDTO := dto.ToPropertyDTO(prop)

	return propDTO, nil
}

func UpdateProperty(ctx context.Context, req dto.UpdatePropertyRequest) (dto.Property, error) {
	if err := dto.Validate(req); err != nil {
		return dto.Property{}, err
	}

	if req.Name != nil {
		existing, err := db.Queries.GetPropertyByID(ctx, req.ID)
		if err != nil {
			return dto.Property{}, err
		}
		if existing.Name != *req.Name {
			itemTypes, err := db.Queries.GetAllItemTypes(ctx)
			if err != nil {
				return dto.Property{}, err
			}
			for _, itemType := range itemTypes {
				if derivedNameUsesProperty(itemType.DerivedNameFormat.String, existing.Name) {
					return dto.Property{}, ErrDerivedNamePropertyInUse
				}
			}
		}

		if err := db.Queries.UpdateProperty_Name(ctx, repository.UpdateProperty_NameParams{
			ID:   req.ID,
			Name: *req.Name,
		}); err != nil {
			return dto.Property{}, err
		}
	}

	if req.Description != nil {
		description := pgtype.Text{String: *req.Description, Valid: true}

		if err := db.Queries.UpdateProperty_Description(ctx, repository.UpdateProperty_DescriptionParams{
			ID:          req.ID,
			Description: description,
		}); err != nil {
			return dto.Property{}, err
		}
	}

	if req.DefaultValueSet {
		existing, err := db.Queries.GetPropertyByID(ctx, req.ID)
		if err != nil {
			return dto.Property{}, err
		}

		if req.DefaultValue != nil {
			if err := dto.Validate(dto.PropertyValueCheck{Value: *req.DefaultValue, ValueType: existing.ValueType}); err != nil {
				return dto.Property{}, err
			}
		}

		if err := db.Queries.UpdateProperty_DefaultValue(ctx, repository.UpdateProperty_DefaultValueParams{
			ID:           req.ID,
			DefaultValue: req.DefaultValue,
		}); err != nil {
			return dto.Property{}, err
		}
	}

	prop, err := db.Queries.GetPropertyByID(ctx, req.ID)
	if err != nil {
		return dto.Property{}, err
	}

	propDTO := dto.ToPropertyDTO(prop)

	return propDTO, nil
}

func DeleteProperty(ctx context.Context, id int64) error {
	rowsAffected, err := db.Queries.DeleteProperty(ctx, id)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
