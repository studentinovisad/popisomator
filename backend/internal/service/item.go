package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/studentinovisad/popisomator/backend/internal/db"
	"github.com/studentinovisad/popisomator/backend/internal/dto"
	"github.com/studentinovisad/popisomator/backend/internal/repository"
)

func populateItemRequestInformation(ctx context.Context, items []dto.Item) error {
	itemIndexes := make(map[int64]int, len(items))
	itemIDs := make([]int64, 0, len(items))
	for index, item := range items {
		itemIndexes[item.ID] = index
		itemIDs = append(itemIDs, item.ID)
	}

	itemRequests, err := db.Queries.CheckItemsForRequests(ctx, itemIDs)
	if err != nil {
		return err
	}

	for _, itemRequest := range itemRequests {
		index := itemIndexes[itemRequest.ItemID]
		items[index].Requests = append(items[index].Requests, dto.Item_RequestInformation{
			UserID: itemRequest.UserID,
			Status: itemRequest.Status,
		})
	}

	return nil
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

	items, err := queriesTx.CreateItems(ctx, repository.CreateItemsParams{
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

		}
	}

	propertyRows, err := queriesTx.GetItemProperties(ctx, itemIDs)
	if err != nil {
		return nil, err
	}
	derivedNameRows, err := queriesTx.GetItemsDerivedNames(ctx, itemIDs)
	if err != nil {
		return nil, err
	}
	populateItemDetails(itemsDTO, propertyRows, derivedNameRows)
	if err := populateItemRequestInformation(ctx, itemsDTO); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return itemsDTO, nil
}

func ListItems(ctx context.Context, req dto.ListItemsRequest) (dto.ItemsPage, error) {
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
		Search:      req.Search,
	})
	if err != nil {
		return dto.ItemsPage{}, err
	}

	items, err := db.Queries.ListItems(ctx, repository.ListItemsParams{
		TypeID:      typeID,
		Consumption: req.Consumption,
		CreatedFrom: createdFrom,
		CreatedTo:   createdTo,
		Search:      req.Search,
		LimitVal:    req.Limit,
		OffsetVal:   req.Offset,
		OrderAsc:    req.Order == "asc",
	})
	if err != nil {
		return dto.ItemsPage{}, err
	}

	itemsDTO := make([]dto.Item, len(items))
	itemIDs := make([]int64, len(items))
	for i, item := range items {
		itemsDTO[i] = dto.ToItemDTO(item)
		itemIDs[i] = item.ID
	}

	if len(itemIDs) > 0 {
		propRows, err := db.Queries.GetItemProperties(ctx, itemIDs)
		if err != nil {
			return dto.ItemsPage{}, err
		}
		derivedNameRows, err := db.Queries.GetItemsDerivedNames(ctx, itemIDs)
		if err != nil {
			return dto.ItemsPage{}, err
		}
		populateItemDetails(itemsDTO, propRows, derivedNameRows)
		if err := populateItemRequestInformation(ctx, itemsDTO); err != nil {
			return dto.ItemsPage{}, err
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

	propertyRows, err := db.Queries.GetItemProperties(ctx, []int64{id})
	if err != nil {
		return dto.Item{}, err
	}

	derivedNameRows, err := db.Queries.GetItemsDerivedNames(ctx, []int64{id})
	if err != nil {
		return dto.Item{}, err
	}

	itemsDTO := []dto.Item{dto.ToItemDTO(item)}
	populateItemDetails(itemsDTO, propertyRows, derivedNameRows)
	if err := populateItemRequestInformation(ctx, itemsDTO); err != nil {
		return dto.Item{}, err
	}

	return itemsDTO[0], nil
}

func UpdateItem(ctx context.Context, req dto.UpdateItemRequest) (dto.Item, error) {
	if err := dto.Validate(req); err != nil {
		return dto.Item{}, err
	}

	var item repository.Item
	if req.TypeID != nil || req.Consumption != nil {
		tx, err := db.BeginTransaction(ctx)
		if err != nil {
			return dto.Item{}, err
		}
		defer tx.Rollback(ctx)
		queriesTx := db.Queries.WithTx(tx)

		if req.TypeID != nil {
			var err error
			if item, err = queriesTx.UpdateItem_Type(ctx, repository.UpdateItem_TypeParams{
				ID:     req.ID,
				TypeID: *req.TypeID,
			}); err != nil {
				return dto.Item{}, err
			}
		}

		if req.Consumption != nil {
			var err error
			if item, err = queriesTx.UpdateItem_Consumption(ctx, repository.UpdateItem_ConsumptionParams{
				ID:          req.ID,
				Consumption: repository.ConsumptionStatus(*req.Consumption),
			}); err != nil {
				return dto.Item{}, err
			}
		}

		if err := tx.Commit(ctx); err != nil {
			return dto.Item{}, err
		}
	} else {
		return dto.Item{}, ErrNoUpdateFields
	}

	propertyRows, err := db.Queries.GetItemProperties(ctx, []int64{item.ID})
	if err != nil {
		return dto.Item{}, err
	}

	derivedNameRows, err := db.Queries.GetItemsDerivedNames(ctx, []int64{item.ID})
	if err != nil {
		return dto.Item{}, err
	}

	itemsDTO := []dto.Item{dto.ToItemDTO(item)}
	populateItemDetails(itemsDTO, propertyRows, derivedNameRows)
	if err := populateItemRequestInformation(ctx, itemsDTO); err != nil {
		return dto.Item{}, err
	}

	return itemsDTO[0], nil
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

func GetItemProperties(ctx context.Context, id int64) ([]dto.ItemProperty, error) {
	propsRows, err := db.Queries.GetItemProperties(ctx, []int64{id})
	if err != nil {
		return nil, err
	}

	itemPropsDTO := make([]dto.ItemProperty, len(propsRows))
	for i, propRow := range propsRows {
		itemPropsDTO[i] = dto.ToItemPropertyDTO(propRow.ItemProperty)
		itemPropsDTO[i].Visibility = string(propRow.Visibility)
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
