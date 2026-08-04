package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/studentinovisad/popisomator/backend/internal/db"
	"github.com/studentinovisad/popisomator/backend/internal/dto"
	"github.com/studentinovisad/popisomator/backend/internal/repository"
)

func GetAllItems(ctx context.Context) ([]dto.Item, error) {
	items, err := db.Queries.GetAllItems(ctx)
	if err != nil {
		return nil, err
	}

	itemsDTO := make([]dto.Item, len(items))
	for i, item := range items {
		itemProps, err := GetItemProperties(ctx, item.ID)
		if err != nil {
			return nil, err
		}
		itemsDTO[i] = dto.ToItemDTO(item)
		itemsDTO[i].Properties = itemProps
	}

	return itemsDTO, nil
}

func GetItem(ctx context.Context, id int64) (dto.Item, error) {
	item, err := db.Queries.GetItemByID(ctx, id)
	if err != nil {
		return dto.Item{}, err
	}

	itemProps, err := GetItemProperties(ctx, id)
	if err != nil {
		return dto.Item{}, err
	}

	itemDTO := dto.ToItemDTO(item)
	itemDTO.Properties = itemProps

	return itemDTO, nil
}

func CreateItem(ctx context.Context) (dto.Item, error) {
	item, err := db.Queries.CreateItem(ctx)
	if err != nil {
		return dto.Item{}, err
	}

	itemDTO := dto.ToItemDTO(item)

	return itemDTO, nil
}

func DeleteItem(ctx context.Context, id int64) error {
	if err := db.Queries.DeleteItem(ctx, id); err != nil {
		return err
	}

	return nil
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
		if err := db.Queries.UpdatePropertyName(ctx, repository.UpdatePropertyNameParams{
			ID:   req.ID,
			Name: *req.Name,
		}); err != nil {
			return dto.Property{}, err
		}
	}

	if req.Description != nil {
		description := pgtype.Text{String: *req.Description, Valid: true}

		if err := db.Queries.UpdatePropertyDescription(ctx, repository.UpdatePropertyDescriptionParams{
			ID:          req.ID,
			Description: description,
		}); err != nil {
			return dto.Property{}, err
		}
	}

	if req.DefaultValue != nil {
		if err := db.Queries.UpdatePropertyDefaultValue(ctx, repository.UpdatePropertyDefaultValueParams{
			ID:           req.ID,
			DefaultValue: *req.DefaultValue,
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
	if err := db.Queries.DeleteProperty(ctx, id); err != nil {
		return err
	}

	return nil
}

func GetItemProperties(ctx context.Context, id int64) ([]dto.ItemProperty, error) {
	itemProps, err := db.Queries.GetItemProperties(ctx, id)
	if err != nil {
		return nil, err
	}

	itemPropsDTO := make([]dto.ItemProperty, len(itemProps))
	for i, itemProp := range itemProps {
		itemPropsDTO[i] = dto.ToItemPropertyDTO(itemProp)
	}

	return itemPropsDTO, nil
}

func AddItemProperty(ctx context.Context, req dto.AddUpdateItemPropertyRequest) (dto.ItemProperty, error) {
	if err := dto.Validate(req); err != nil {
		return dto.ItemProperty{}, err
	}

	itemProp, err := db.Queries.AddItemProperty(ctx, repository.AddItemPropertyParams{
		ItemID:        req.ItemID,
		PropertyID:    req.PropertyID,
		PropertyValue: req.Value,
	})
	if err != nil {
		return dto.ItemProperty{}, err
	}

	itemPropDTO := dto.ToItemPropertyDTO(itemProp)

	return itemPropDTO, nil
}

func UpdateItemProperty(ctx context.Context, req dto.AddUpdateItemPropertyRequest) (dto.ItemProperty, error) {
	if err := dto.Validate(req); err != nil {
		return dto.ItemProperty{}, err
	}

	itemProp, err := db.Queries.UpdateItemProperty(ctx, repository.UpdateItemPropertyParams{
		ItemID:        req.ItemID,
		PropertyID:    req.PropertyID,
		PropertyValue: req.Value,
	})
	if err != nil {
		return dto.ItemProperty{}, err
	}

	itemPropDTO := dto.ToItemPropertyDTO(itemProp)

	return itemPropDTO, nil
}

func RemoveItemProperty(ctx context.Context, itemId int64, propId int64) error {
	if err := db.Queries.RemoveItemProperty(ctx, repository.RemoveItemPropertyParams{
		ItemID:     itemId,
		PropertyID: propId,
	}); err != nil {
		return err
	}

	return nil
}
