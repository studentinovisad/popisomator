package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/studentinovisad/popisomator/backend/internal/db"
	"github.com/studentinovisad/popisomator/backend/internal/dto"
	"github.com/studentinovisad/popisomator/backend/internal/repository"
)

func GetAllItemTypes(ctx context.Context) ([]dto.ItemType, error) {
	itemTypes, err := db.Queries.GetAllItemTypes(ctx)
	if err != nil {
		return nil, err
	}

	typesDTO := make([]dto.ItemType, len(itemTypes))
	for i, itemType := range itemTypes {
		typesDTO[i] = dto.ToItemTypeDTO(itemType)
	}

	return typesDTO, nil
}

func GetItemTypeOptions(ctx context.Context) ([]dto.ItemTypeOption, error) {
	itemTypes, err := db.Queries.ListItemTypeOptions(ctx)
	if err != nil {
		return nil, err
	}

	typeOptionsDTO := make([]dto.ItemTypeOption, len(itemTypes))
	for i, itemType := range itemTypes {
		typeOptionsDTO[i] = dto.ItemTypeOption{ID: itemType.ID, Name: itemType.Name}
	}

	return typeOptionsDTO, nil
}

func ListItemTypes(ctx context.Context, limit, offset int32, search string) (dto.ItemTypesPage, error) {
	searchText := pgtype.Text{String: search, Valid: true}
	total, err := db.Queries.CountItemTypes(ctx, searchText)
	if err != nil {
		return dto.ItemTypesPage{}, err
	}

	itemTypes, err := db.Queries.ListItemTypes(ctx, repository.ListItemTypesParams{
		LimitVal:  limit,
		OffsetVal: offset,
		Search:    searchText,
	})
	if err != nil {
		return dto.ItemTypesPage{}, err
	}

	typesDTO := make([]dto.ItemType, len(itemTypes))
	index := make(map[int64]int, len(itemTypes))
	typeIDs := make([]int64, len(itemTypes))
	for i, itemType := range itemTypes {
		typesDTO[i] = dto.ToItemTypeDTO(itemType)
		index[itemType.ID] = i
		typeIDs[i] = itemType.ID
	}

	if len(typeIDs) > 0 {
		props, err := db.Queries.GetItemTypeProperties(ctx, typeIDs)
		if err != nil {
			return dto.ItemTypesPage{}, err
		}
		for _, prop := range props {
			idx := index[prop.TypeID]
			propDTO := dto.ToItemTypePropertyDTO(prop)
			typesDTO[idx].Properties = append(typesDTO[idx].Properties, propDTO)
		}
	}

	return dto.ItemTypesPage{
		Items:  typesDTO,
		Limit:  limit,
		Offset: offset,
		Total:  total,
	}, nil
}

func GetItemType(ctx context.Context, id int64) (dto.ItemType, error) {
	itemType, err := db.Queries.GetItemTypeByID(ctx, id)
	if err != nil {
		return dto.ItemType{}, err
	}

	typeProps, err := GetItemTypeProperties(ctx, id)
	if err != nil {
		return dto.ItemType{}, err
	}

	itemTypeDTO := dto.ToItemTypeDTO(itemType)
	itemTypeDTO.Properties = typeProps

	return itemTypeDTO, nil
}

func CreateItemType(ctx context.Context, req dto.CreateItemTypeRequest) (dto.ItemType, error) {
	if err := dto.Validate(req); err != nil {
		return dto.ItemType{}, err
	}

	tx, err := db.BeginTransaction(ctx)
	if err != nil {
		return dto.ItemType{}, err
	}
	defer tx.Rollback(ctx)
	queriesTx := db.Queries.WithTx(tx)

	if err := validateDerivedNameFormat(ctx, queriesTx, req.DerivedNameFormat, req.Properties); err != nil {
		return dto.ItemType{}, err
	}

	description := pgtype.Text{String: "", Valid: false}
	if len(req.Description) > 0 {
		description = pgtype.Text{String: req.Description, Valid: true}
	}

	derived_name_format := pgtype.Text{String: "", Valid: false}
	if len(req.DerivedNameFormat) > 0 {
		derived_name_format = pgtype.Text{String: req.DerivedNameFormat, Valid: true}
	}

	itemType, err := queriesTx.CreateItemType(ctx, repository.CreateItemTypeParams{
		Name:              req.Name,
		Description:       description,
		DerivedNameFormat: derived_name_format,
	})
	if err != nil {
		return dto.ItemType{}, err
	}

	itemTypeDTO := dto.ToItemTypeDTO(itemType)

	if req.Properties != nil {
		for _, propRequest := range req.Properties {
			if propRequest.DefaultValue != nil {
				if err := validatePropertyValue(ctx, queriesTx, propRequest.ID, *propRequest.DefaultValue); err != nil {
					return dto.ItemType{}, err
				}
			}

			visibility := repository.PropertyVisibilityOverview
			if propRequest.Visibility != "" {
				visibility = repository.PropertyVisibility(propRequest.Visibility)
			}

			prop, err := queriesTx.AddItemTypeProperty(ctx, repository.AddItemTypePropertyParams{
				TypeID:       itemType.ID,
				PropertyID:   propRequest.ID,
				DefaultValue: propRequest.DefaultValue,
				Visibility:   visibility,
			})
			if err != nil {
				return dto.ItemType{}, err
			}

			propDTO := dto.ToItemTypePropertyDTO(prop)
			itemTypeDTO.Properties = append(itemTypeDTO.Properties, propDTO)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return dto.ItemType{}, err
	}

	return itemTypeDTO, nil
}

func UpdateItemType(ctx context.Context, req dto.UpdateItemTypeRequest) (dto.ItemType, error) {
	if err := dto.Validate(req); err != nil {
		return dto.ItemType{}, err
	}
	if req.DerivedNameFormat != nil {
		typeProperties, err := GetItemTypeProperties(ctx, req.ID)
		if err != nil {
			return dto.ItemType{}, err
		}
		if err := validateDerivedNameFormat(ctx, db.Queries, *req.DerivedNameFormat, typeProperties); err != nil {
			return dto.ItemType{}, err
		}
	}

	var itemType repository.ItemType
	if req.Name != nil || req.Description != nil || req.DerivedNameFormat != nil {
		tx, err := db.BeginTransaction(ctx)
		if err != nil {
			return dto.ItemType{}, err
		}
		defer tx.Rollback(ctx)
		queriesTx := db.Queries.WithTx(tx)

		if req.Name != nil {
			var err error
			if itemType, err = queriesTx.UpdateItemType_Name(ctx, repository.UpdateItemType_NameParams{
				ID:   req.ID,
				Name: *req.Name,
			}); err != nil {
				return dto.ItemType{}, err
			}
		}

		if req.Description != nil {
			description := pgtype.Text{String: *req.Description, Valid: true}

			var err error
			if itemType, err = queriesTx.UpdateItemType_Description(ctx, repository.UpdateItemType_DescriptionParams{
				ID:          req.ID,
				Description: description,
			}); err != nil {
				return dto.ItemType{}, err
			}
		}

		if req.DerivedNameFormat != nil {
			derived_name_format := pgtype.Text{String: *req.DerivedNameFormat, Valid: true}

			var err error
			if itemType, err = queriesTx.UpdateItemType_DerivedNameFormat(ctx, repository.UpdateItemType_DerivedNameFormatParams{
				ID:                req.ID,
				DerivedNameFormat: derived_name_format,
			}); err != nil {
				return dto.ItemType{}, err
			}
		}

		if err := tx.Commit(ctx); err != nil {
			return dto.ItemType{}, err
		}
	} else {
		return dto.ItemType{}, ErrNoUpdateFields
	}

	typeProps, err := GetItemTypeProperties(ctx, itemType.ID)
	if err != nil {
		return dto.ItemType{}, err
	}

	itemTypeDTO := dto.ToItemTypeDTO(itemType)
	itemTypeDTO.Properties = typeProps

	return itemTypeDTO, nil
}

func DeleteItemType(ctx context.Context, id int64) error {
	rowsAffected, err := db.Queries.DeleteItemType(ctx, id)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

//
//	Item Type Properties
//

func GetItemTypeProperties(ctx context.Context, id int64) ([]dto.ItemTypeProperty, error) {
	typeProps, err := db.Queries.GetItemTypeProperties(ctx, []int64{id})
	if err != nil {
		return nil, err
	}

	typePropsDTO := make([]dto.ItemTypeProperty, len(typeProps))
	for i, typeProp := range typeProps {
		typePropsDTO[i] = dto.ToItemTypePropertyDTO(typeProp)
	}

	return typePropsDTO, nil
}

func AddItemTypeProperty(ctx context.Context, req dto.AddUpdateItemTypePropertyRequest) (dto.ItemTypeProperty, error) {
	if err := dto.Validate(req); err != nil {
		return dto.ItemTypeProperty{}, err
	}

	if req.DefaultValue != nil {
		if err := validatePropertyValue(ctx, db.Queries, req.PropertyID, *req.DefaultValue); err != nil {
			return dto.ItemTypeProperty{}, err
		}
	}

	visibility := repository.PropertyVisibilityOverview
	if req.Visibility != nil {
		visibility = *req.Visibility
	}

	typeProp, err := db.Queries.AddItemTypeProperty(ctx, repository.AddItemTypePropertyParams{
		TypeID:       req.TypeID,
		PropertyID:   req.PropertyID,
		DefaultValue: req.DefaultValue,
		Visibility:   visibility,
	})
	if err != nil {
		return dto.ItemTypeProperty{}, err
	}

	typePropDTO := dto.ToItemTypePropertyDTO(typeProp)

	return typePropDTO, nil
}

func UpdateItemTypeProperty(ctx context.Context, req dto.AddUpdateItemTypePropertyRequest) (dto.ItemTypeProperty, error) {
	if err := dto.Validate(req); err != nil {
		return dto.ItemTypeProperty{}, err
	}

	if req.DefaultValue == nil && req.Visibility == nil {
		return dto.ItemTypeProperty{}, ErrNoUpdateFields
	}

	var typeProp repository.ItemTypeProperty

	if req.DefaultValue != nil {
		if err := validatePropertyValue(ctx, db.Queries, req.PropertyID, *req.DefaultValue); err != nil {
			return dto.ItemTypeProperty{}, err
		}

		var err error
		typeProp, err = db.Queries.UpdateItemTypeProperty_DefaultValue(ctx, repository.UpdateItemTypeProperty_DefaultValueParams{
			TypeID:       req.TypeID,
			PropertyID:   req.PropertyID,
			DefaultValue: req.DefaultValue,
		})
		if err != nil {
			return dto.ItemTypeProperty{}, err
		}
	}

	if req.Visibility != nil {
		var err error
		typeProp, err = db.Queries.UpdateItemTypeProperty_Visibility(ctx, repository.UpdateItemTypeProperty_VisibilityParams{
			TypeID:     req.TypeID,
			PropertyID: req.PropertyID,
			Visibility: *req.Visibility,
		})
		if err != nil {
			return dto.ItemTypeProperty{}, err
		}
	}

	typePropDTO := dto.ToItemTypePropertyDTO(typeProp)

	return typePropDTO, nil
}

func RemoveItemTypeProperty(ctx context.Context, typeId int64, propId int64) error {
	itemType, err := db.Queries.GetItemTypeByID(ctx, typeId)
	if err != nil {
		return err
	}
	property, err := db.Queries.GetPropertyByID(ctx, propId)
	if err != nil {
		return err
	}
	if derivedNameUsesProperty(itemType.DerivedNameFormat.String, property.Name) {
		return ErrInvalidDerivedNameFormat
	}

	rowsAffected, err := db.Queries.RemoveItemTypeProperty(ctx, repository.RemoveItemTypePropertyParams{
		TypeID:     typeId,
		PropertyID: propId,
	})
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
