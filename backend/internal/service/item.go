package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/studentinovisad/popisomator/backend/internal/db"
	"github.com/studentinovisad/popisomator/backend/internal/dto"
	"github.com/studentinovisad/popisomator/backend/internal/repository"
)

//
// Items
//

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

func CreateItem(ctx context.Context, req dto.CreateItemRequest) (dto.Item, error) {
	if err := dto.Validate(req); err != nil {
		return dto.Item{}, err
	}

	tx, err := db.BeginTransaction(ctx)
	if err != nil {
		return dto.Item{}, err
	}
	defer tx.Rollback(ctx)
	queriesTx := db.Queries.WithTx(tx)

	typeId := pgtype.Int8{}
	if req.TypeID != nil {
		typeId = pgtype.Int8{Int64: *req.TypeID, Valid: true}
	}

	item, err := queriesTx.CreateItem(ctx, typeId)
	if err != nil {
		return dto.Item{}, err
	}

	itemDTO := dto.ToItemDTO(item)

	if req.Properties != nil {
		for _, propRequest := range req.Properties {
			prop, err := queriesTx.AddItemProperty(ctx, repository.AddItemPropertyParams{
				ItemID:        item.ID,
				PropertyID:    propRequest.ID,
				PropertyValue: propRequest.Value,
			})
			if err != nil {
				return dto.Item{}, err
			}

			propDTO := dto.ToItemPropertyDTO(prop)
			itemDTO.Properties = append(itemDTO.Properties, propDTO)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return dto.Item{}, err
	}

	return itemDTO, nil
}

func DeleteItem(ctx context.Context, id int64) error {
	if err := db.Queries.DeleteItem(ctx, id); err != nil {
		return err
	}

	return nil
}

//
// Properties
//

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
	if err := db.Queries.DeleteProperty(ctx, id); err != nil {
		return err
	}

	return nil
}

//
//	Item Properties
//

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

//
// Item Types
//

func GetAllItemTypes(ctx context.Context) ([]dto.ItemType, error) {
	itemTypes, err := db.Queries.GetAllItemTypes(ctx)
	if err != nil {
		return nil, err
	}

	itemTypesDTO := make([]dto.ItemType, len(itemTypes))
	for i, itemType := range itemTypes {
		typeProps, err := GetItemTypeProperties(ctx, itemType.ID)
		if err != nil {
			return nil, err
		}
		itemTypesDTO[i] = dto.ToItemTypeDTO(itemType)
		itemTypesDTO[i].Properties = typeProps
	}

	return itemTypesDTO, nil
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

	description := pgtype.Text{String: "", Valid: false}
	if len(req.Description) > 0 {
		description = pgtype.Text{String: req.Description, Valid: true}
	}

	itemType, err := queriesTx.CreateItemType(ctx, repository.CreateItemTypeParams{
		Name:        req.Name,
		Description: description,
	})
	if err != nil {
		return dto.ItemType{}, err
	}

	itemTypeDTO := dto.ToItemTypeDTO(itemType)

	if req.Properties != nil {
		for _, propRequest := range req.Properties {
			prop, err := queriesTx.AddItemTypeProperty(ctx, repository.AddItemTypePropertyParams{
				TypeID:       itemType.ID,
				PropertyID:   propRequest.ID,
				DefaultValue: propRequest.DefaultValue,
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

func DeleteItemType(ctx context.Context, id int64) error {
	if err := db.Queries.DeleteItemType(ctx, id); err != nil {
		return err
	}

	return nil
}

//
//	Item Type Properties
//

func GetItemTypeProperties(ctx context.Context, id int64) ([]dto.ItemTypeProperty, error) {
	typeProps, err := db.Queries.GetItemTypeProperties(ctx, id)
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

	typeProp, err := db.Queries.AddItemTypeProperty(ctx, repository.AddItemTypePropertyParams{
		TypeID:       req.TypeID,
		PropertyID:   req.PropertyID,
		DefaultValue: req.DefaultValue,
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

	typeProp, err := db.Queries.UpdateItemTypeProperty(ctx, repository.UpdateItemTypePropertyParams{
		TypeID:       req.TypeID,
		PropertyID:   req.PropertyID,
		DefaultValue: req.DefaultValue,
	})
	if err != nil {
		return dto.ItemTypeProperty{}, err
	}

	typePropDTO := dto.ToItemTypePropertyDTO(typeProp)

	return typePropDTO, nil
}

func RemoveItemTypeProperty(ctx context.Context, typeId int64, propId int64) error {
	if err := db.Queries.RemoveItemTypeProperty(ctx, repository.RemoveItemTypePropertyParams{
		TypeID:     typeId,
		PropertyID: propId,
	}); err != nil {
		return err
	}

	return nil
}
