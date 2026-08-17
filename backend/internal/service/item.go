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

//
// Items
//

func GetAllItems(ctx context.Context, req dto.ListItemsRequest) (dto.ItemsPage, error) {
	if err := dto.Validate(req); err != nil {
		return dto.ItemsPage{}, err
	}

	typeID := pgtype.Int8{}
	if req.TypeID != nil {
		typeID = pgtype.Int8{Int64: *req.TypeID, Valid: true}
	}

	createdFrom := pgtype.Timestamptz{}
	if req.CreatedFrom != nil {
		createdFrom = pgtype.Timestamptz{Time: *req.CreatedFrom, Valid: true}
	}

	createdTo := pgtype.Timestamptz{}
	if req.CreatedTo != nil {
		createdTo = pgtype.Timestamptz{Time: *req.CreatedTo, Valid: true}
	}

	totalItems, err := db.Queries.CountItems(ctx, repository.CountItemsParams{
		TypeID:      typeID,
		Consumption: req.Consumption,
		CreatedFrom: createdFrom,
		CreatedTo:   createdTo,
	})
	if err != nil {
		return dto.ItemsPage{}, err
	}

	items, err := db.Queries.ListItems(ctx, repository.ListItemsParams{
		TypeID:      typeID,
		Consumption: req.Consumption,
		CreatedFrom: createdFrom,
		CreatedTo:   createdTo,
		LimitVal:    req.Limit,
		OffsetVal:   req.Offset,
		OrderAsc:    req.Order == "asc",
	})
	if err != nil {
		return dto.ItemsPage{}, err
	}

	itemsDTO := make([]dto.Item, len(items))
	index := make(map[int64]int, len(items))
	itemIDs := make([]int64, len(items))
	for i, item := range items {
		itemsDTO[i] = dto.ToItemDTO(item)
		index[item.ID] = i
		itemIDs[i] = item.ID
	}

	if len(itemIDs) > 0 {
		propRows, err := db.Queries.GetItemPropertiesForItems(ctx, itemIDs)
		if err != nil {
			return dto.ItemsPage{}, err
		}
		for _, propRow := range propRows {
			idx := index[propRow.ItemID]
			itemsDTO[idx].Properties = append(itemsDTO[idx].Properties, dto.RowToItemPropertyDTO(propRow))
		}
	}

	return dto.ItemsPage{
		Items:  itemsDTO,
		Limit:  req.Limit,
		Offset: req.Offset,
		Total:  totalItems,
	}, nil
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

func CreateItem(ctx context.Context, req dto.CreateItemRequest) ([]dto.Item, error) {
	if err := dto.Validate(req); err != nil {
		return nil, err
	}

	tx, err := db.BeginTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	queriesTx := db.Queries.WithTx(tx)

	items, err := queriesTx.CreateItemBulk(ctx, repository.CreateItemBulkParams{
		TypeID: req.TypeID,
		Amount: req.Amount,
	})
	if err != nil {
		return nil, err
	}

	itemsDTO := make([]dto.Item, len(items))
	itemIDs := make([]int64, len(items))
	for i, item := range items {
		itemsDTO[i] = dto.ToItemDTO(item)
		itemIDs[i] = item.ID
	}

	if req.Properties != nil {
		for _, propRequest := range req.Properties {
			if err := validatePropertyValue(ctx, queriesTx, propRequest.ID, propRequest.Value); err != nil {
				return nil, err
			}

			props, err := queriesTx.AddItemPropertyBulk(ctx, repository.AddItemPropertyBulkParams{
				ItemIds:       itemIDs,
				PropertyID:    propRequest.ID,
				PropertyValue: propRequest.Value,
			})
			if err != nil {
				return nil, err
			}
			if props != nil && len(props) == 0 {
				return nil, errors.New("No item properties returned")
			}

			propDTO := dto.ToItemPropertyDTO(props[0])
			for _, itemDTO := range itemsDTO {
				itemDTO.Properties = append(itemDTO.Properties, propDTO)
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return itemsDTO, nil
}

func ConsumeItem(ctx context.Context, req dto.ConsumeItemRequest) error {
	consumption := repository.ConsumptionStatusFullyConsumed
	if req.Status != nil {
		consumption = *req.Status
	}

	rowsAffected, err := db.Queries.UpdateItemConsumption(ctx, repository.UpdateItemConsumptionParams{
		ID:          req.ID,
		Consumption: consumption,
	})
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func SetItemType(ctx context.Context, req dto.SetItemTypeRequest) (dto.Item, error) {
	if err := dto.Validate(req); err != nil {
		return dto.Item{}, err
	}

	// db.Queries.UpdateItemType sets items.type_id (reassigns this item's type) — unrelated to
	// service.UpdateItemType below, which edits an item_type's own name/description.
	item, err := db.Queries.UpdateItemType(ctx, repository.UpdateItemTypeParams{
		ID:     req.ID,
		TypeID: req.TypeID,
	})
	if err != nil {
		return dto.Item{}, err
	}

	itemProps, err := GetItemProperties(ctx, item.ID)
	if err != nil {
		return dto.Item{}, err
	}

	itemDTO := dto.ToItemDTO(item)
	itemDTO.Properties = itemProps

	return itemDTO, nil
}

func DeleteItem(ctx context.Context, id int64) error {
	rowsAffected, err := db.Queries.DeleteItem(ctx, id)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
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

	items := make([]dto.Property, len(properties))
	for index, property := range properties {
		items[index] = dto.ToPropertyDTO(property)
	}

	return dto.PropertiesPage{Items: items, Limit: limit, Offset: offset, Total: total}, nil
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
	rowsAffected, err := db.Queries.DeleteProperty(ctx, id)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
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

	if err := validatePropertyValue(ctx, db.Queries, req.PropertyID, req.Value); err != nil {
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

	if err := validatePropertyValue(ctx, db.Queries, req.PropertyID, req.Value); err != nil {
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
	rowsAffected, err := db.Queries.RemoveItemProperty(ctx, repository.RemoveItemPropertyParams{
		ItemID:     itemId,
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

//
// Item Types
//

func GetAllItemTypes(ctx context.Context) ([]dto.ItemType, error) {
	rows, err := db.Queries.GetAllItemTypesWithProperties(ctx)
	if err != nil {
		return nil, err
	}

	itemTypes := make([]dto.ItemType, 0, len(rows))
	index := make(map[int64]int, len(rows))
	for _, row := range rows {
		idx, ok := index[row.ItemType.ID]
		if !ok {
			idx = len(itemTypes)
			index[row.ItemType.ID] = idx
			itemTypes = append(itemTypes, dto.ToItemTypeDTO(row.ItemType))
		}

		if row.PropertyID.Valid {
			itemTypes[idx].Properties = append(itemTypes[idx].Properties, dto.ItemTypeProperty{
				ID:           row.PropertyID.Int64,
				DefaultValue: row.DefaultValue,
				Name:         row.PropertyName.String,
			})
		}
	}

	return itemTypes, nil
}

func ListItemTypes(ctx context.Context, limit, offset int32) (dto.ItemTypesPage, error) {
	total, err := db.Queries.CountItemTypes(ctx)
	if err != nil {
		return dto.ItemTypesPage{}, err
	}

	rows, err := db.Queries.ListItemTypesWithProperties(ctx, repository.ListItemTypesWithPropertiesParams{
		PageLimit:  limit,
		PageOffset: offset,
	})
	if err != nil {
		return dto.ItemTypesPage{}, err
	}

	itemTypes := make([]dto.ItemType, 0, len(rows))
	index := make(map[int64]int, len(rows))
	for _, row := range rows {
		idx, ok := index[row.ItemType.ID]
		if !ok {
			idx = len(itemTypes)
			index[row.ItemType.ID] = idx
			itemTypes = append(itemTypes, dto.ToItemTypeDTO(row.ItemType))
		}

		if row.PropertyID.Valid {
			itemTypes[idx].Properties = append(itemTypes[idx].Properties, dto.ItemTypeProperty{
				ID:           row.PropertyID.Int64,
				DefaultValue: row.DefaultValue,
				Name:         row.PropertyName.String,
			})
		}
	}

	return dto.ItemTypesPage{Items: itemTypes, Limit: limit, Offset: offset, Total: total}, nil
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

func UpdateItemType(ctx context.Context, req dto.UpdateItemTypeRequest) (dto.ItemType, error) {
	if err := dto.Validate(req); err != nil {
		return dto.ItemType{}, err
	}

	if req.Name != nil || req.Description != nil {
		tx, err := db.BeginTransaction(ctx)
		if err != nil {
			return dto.ItemType{}, err
		}
		defer tx.Rollback(ctx)
		queriesTx := db.Queries.WithTx(tx)

		if req.Name != nil {
			if _, err := queriesTx.UpdateItemTypeName(ctx, repository.UpdateItemTypeNameParams{
				ID:   req.ID,
				Name: *req.Name,
			}); err != nil {
				return dto.ItemType{}, err
			}
		}

		if req.Description != nil {
			description := pgtype.Text{String: *req.Description, Valid: true}

			if _, err := queriesTx.UpdateItemTypeDescription(ctx, repository.UpdateItemTypeDescriptionParams{
				ID:          req.ID,
				Description: description,
			}); err != nil {
				return dto.ItemType{}, err
			}
		}

		if req.DerivedNameFormat != nil {
			derived_name_format := pgtype.Text{String: *req.DerivedNameFormat, Valid: true}

			if _, err := queriesTx.UpdateItemTypeDerivedNameFormat(ctx, repository.UpdateItemTypeDerivedNameFormatParams{
				ID:                req.ID,
				DerivedNameFormat: derived_name_format,
			}); err != nil {
				return dto.ItemType{}, err
			}
		}

		if err := tx.Commit(ctx); err != nil {
			return dto.ItemType{}, err
		}
	}

	itemType, err := db.Queries.GetItemTypeByID(ctx, req.ID)
	if err != nil {
		return dto.ItemType{}, err
	}

	typeProps, err := GetItemTypeProperties(ctx, req.ID)
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

	if req.DefaultValue != nil {
		if err := validatePropertyValue(ctx, db.Queries, req.PropertyID, *req.DefaultValue); err != nil {
			return dto.ItemTypeProperty{}, err
		}
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

	var typeProp repository.ItemTypeProperty

	if req.DefaultValue != nil {
		if err := validatePropertyValue(ctx, db.Queries, req.PropertyID, *req.DefaultValue); err != nil {
			return dto.ItemTypeProperty{}, err
		}

		var err error
		typeProp, err = db.Queries.UpdateItemTypePropertyDefaultValue(ctx, repository.UpdateItemTypePropertyDefaultValueParams{
			TypeID:       req.TypeID,
			PropertyID:   req.PropertyID,
			DefaultValue: req.DefaultValue,
		})
		if err != nil {
			return dto.ItemTypeProperty{}, err
		}
	}

	if req.Visibility != nil {
		if err := validatePropertyValue(ctx, db.Queries, req.PropertyID, string(*req.Visibility)); err != nil {
			return dto.ItemTypeProperty{}, err
		}

		var err error
		typeProp, err = db.Queries.UpdateItemTypePropertyVisibility(ctx, repository.UpdateItemTypePropertyVisibilityParams{
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
