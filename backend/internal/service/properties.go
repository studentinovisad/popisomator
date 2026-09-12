package service

import (
	"context"
	"encoding/json"
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
//
// The property itself comes back because every caller goes on to write an audit entry naming it,
// and the lookup has already happened here - returning it saves each of them a second round trip.
func validatePropertyValue(ctx context.Context, q repository.Querier, propertyID int64, rawValue []byte) (repository.Property, error) {
	prop, err := q.GetPropertyByID(ctx, propertyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.Property{}, ErrInvalidReference
		}
		return repository.Property{}, err
	}

	if err := dto.Validate(dto.PropertyValueCheck{Value: rawValue, ValueType: prop.ValueType}); err != nil {
		return repository.Property{}, err
	}

	return prop, nil
}

func GetAllProperties(ctx context.Context) ([]dto.Property, error) {
	// nil names no ids, which the query reads as no filter.
	props, err := db.Queries.GetProperties(ctx, nil)
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

func ListProperties(ctx context.Context, limit, offset int32, search string) (dto.PropertiesPage, error) {
	total, err := db.Queries.CountProperties(ctx, search)
	if err != nil {
		return dto.PropertiesPage{}, err
	}

	properties, err := db.Queries.ListProperties(ctx, repository.ListPropertiesParams{
		PageLimit:  limit,
		PageOffset: offset,
		Search:     search,
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

	tx, err := db.BeginTransaction(ctx)
	if err != nil {
		return dto.Property{}, err
	}
	defer tx.Rollback(ctx)
	queriesTx := db.Queries.WithTx(tx)

	prop, err := queriesTx.CreateProperty(ctx, repository.CreatePropertyParams{
		Name:         req.Name,
		Description:  description,
		ValueType:    req.ValueType,
		DefaultValue: req.DefaultValue,
	})
	if err != nil {
		return dto.Property{}, err
	}

	changes := auditNew(nil, "name", dto.AuditValueTypeText, prop.Name)
	changes = auditNew(changes, "value_type", dto.AuditValueTypeText, prop.ValueType)
	if prop.Description.Valid {
		changes = auditNew(changes, "description", dto.AuditValueTypeText, prop.Description.String)
	}
	if prop.DefaultValue != nil {
		changes = append(changes, dto.AuditChange{
			Key:       "default_value",
			ValueType: prop.ValueType,
			New:       *prop.DefaultValue,
		})
	}

	if err := writeAudit(ctx, queriesTx, repository.AuditActionPropertyCreate, repository.AuditTargetTypeProperty,
		prop.ID, prop.Name, changes, dto.AuditContext{}); err != nil {
		return dto.Property{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return dto.Property{}, err
	}

	return dto.ToPropertyDTO(prop), nil
}

func UpdateProperty(ctx context.Context, req dto.UpdatePropertyRequest) (dto.Property, error) {
	if err := dto.Validate(req); err != nil {
		return dto.Property{}, err
	}

	// Up to three separate UPDATEs, so they belong in one transaction whether or not anything is
	// being audited: a failure partway used to leave the property half edited.
	tx, err := db.BeginTransaction(ctx)
	if err != nil {
		return dto.Property{}, err
	}
	defer tx.Rollback(ctx)
	queriesTx := db.Queries.WithTx(tx)

	// One read of the pre-image up front serves both the derived-name guard and the audit diff.
	existing, err := queriesTx.GetPropertyByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.Property{}, ErrNotFound
		}
		return dto.Property{}, err
	}

	if req.Name != nil {
		if existing.Name != *req.Name {
			itemTypes, err := queriesTx.GetAllItemTypes(ctx)
			if err != nil {
				return dto.Property{}, err
			}
			for _, itemType := range itemTypes {
				if derivedNameUsesProperty(itemType.DerivedNameFormat.String, existing.Name) {
					return dto.Property{}, ErrDerivedNamePropertyInUse
				}
			}
		}

		if err := queriesTx.UpdateProperty_Name(ctx, repository.UpdateProperty_NameParams{
			ID:   req.ID,
			Name: *req.Name,
		}); err != nil {
			return dto.Property{}, err
		}
	}

	if req.Description != nil {
		description := pgtype.Text{String: *req.Description, Valid: true}

		if err := queriesTx.UpdateProperty_Description(ctx, repository.UpdateProperty_DescriptionParams{
			ID:          req.ID,
			Description: description,
		}); err != nil {
			return dto.Property{}, err
		}
	}

	if req.DefaultValueSet {
		if req.DefaultValue != nil {
			if err := dto.Validate(dto.PropertyValueCheck{Value: *req.DefaultValue, ValueType: existing.ValueType}); err != nil {
				return dto.Property{}, err
			}
		}

		if err := queriesTx.UpdateProperty_DefaultValue(ctx, repository.UpdateProperty_DefaultValueParams{
			ID:           req.ID,
			DefaultValue: req.DefaultValue,
		}); err != nil {
			return dto.Property{}, err
		}
	}

	prop, err := queriesTx.GetPropertyByID(ctx, req.ID)
	if err != nil {
		return dto.Property{}, err
	}

	changes := auditDiff(nil, "name", dto.AuditValueTypeText, existing.Name, prop.Name)
	changes = auditDiff(changes, "description", dto.AuditValueTypeText,
		existing.Description.String, prop.Description.String)
	changes = auditRawDiff(changes, "default_value", prop.ValueType,
		optionalRawJSON(existing.DefaultValue), optionalRawJSON(prop.DefaultValue))

	if len(changes) > 0 {
		if err := writeAudit(ctx, queriesTx, repository.AuditActionPropertyUpdate, repository.AuditTargetTypeProperty,
			prop.ID, prop.Name, changes, dto.AuditContext{}); err != nil {
			return dto.Property{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return dto.Property{}, err
	}

	return dto.ToPropertyDTO(prop), nil
}

func DeleteProperty(ctx context.Context, id int64) error {
	tx, err := db.BeginTransaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	queriesTx := db.Queries.WithTx(tx)

	// Read before the delete, so the entry can say which property this was rather than just its id.
	existing, err := queriesTx.GetPropertyByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	rowsAffected, err := queriesTx.DeleteProperty(ctx, id)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	// Identity only, like every other deletion: the name it went by, and who removed it.
	if err := writeAudit(ctx, queriesTx, repository.AuditActionPropertyDelete, repository.AuditTargetTypeProperty,
		id, existing.Name, nil, dto.AuditContext{}); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// optionalRawJSON flattens a nullable jsonb column to the raw JSON auditRawDiff compares, so clearing
// a default value reads as a move to null rather than as no change at all.
func optionalRawJSON(value *json.RawMessage) json.RawMessage {
	if value == nil {
		return json.RawMessage("null")
	}

	return *value
}
