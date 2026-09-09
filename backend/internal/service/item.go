package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/studentinovisad/popisomator/backend/internal/db"
	"github.com/studentinovisad/popisomator/backend/internal/dto"
	"github.com/studentinovisad/popisomator/backend/internal/repository"
)

func populateItemRequestInformation(ctx context.Context, items []dto.Item, viewerID int64) error {
	itemIndexes := make(map[int64]int, len(items))
	itemIDs := make([]int64, 0, len(items))
	for index, item := range items {
		itemIndexes[item.ID] = index
		itemIDs = append(itemIDs, item.ID)
	}

	rows, err := db.Queries.GetItemsRequestStatuses(ctx, repository.GetItemsRequestStatusesParams{
		ViewerID: viewerID,
		ItemIds:  itemIDs,
	})
	if err != nil {
		return err
	}

	for _, row := range rows {
		index := itemIndexes[row.ItemID]
		if row.UserID == viewerID {
			items[index].RequestStatus = &row.Status
		} else {
			items[index].HolderName = &row.UserFullName
		}
	}

	return nil
}

func populateItemDetails(
	ctx context.Context,
	q *repository.Queries,
	items []dto.Item,
) error {
	itemIDs := make([]int64, len(items))
	itemIndexes := make(map[int64]int, len(items))
	for index, item := range items {
		itemIndexes[item.ID] = index
		itemIDs[index] = item.ID
	}

	propertyRows, err := q.GetItemProperties(ctx, itemIDs)
	if err != nil {
		return err
	}
	itemTypeRows, err := q.GetItemTypesByItemIDs(ctx, itemIDs)
	if err != nil {
		return err
	}
	derivedNameRows, err := q.GetItemsDerivedNames(ctx, itemIDs)
	if err != nil {
		return err
	}

	itemTypesDTO := make(map[int64]dto.ItemType, len(items))
	for _, row := range itemTypeRows {
		itemTypesDTO[row.ItemID] = dto.ToItemTypeDTO(row.ItemType)
	}

	for _, row := range propertyRows {
		itemIndex, exists := itemIndexes[row.ItemProperty.ItemID]
		if !exists {
			continue
		}
		item := &items[itemIndex]
		itemType, itemTypeExists := itemTypesDTO[item.ID]

		property := dto.ToItemPropertyDTO(row.ItemProperty)
		property.Visibility = string(row.Visibility)
		property.ValueType = row.PropertyType
		switch property.ValueType {
		case "expiry":
			var propertyValue string
			if err := json.Unmarshal(property.Value, &propertyValue); err != nil {
				log.Printf("Couldn't unmarshal expiry date value %v. Error: %v", string(property.Value), err)
				continue
			}
			expiryTime, err := time.Parse(time.DateOnly, propertyValue)
			if err != nil {
				log.Printf("Couldn't parse expiry date %v. Error: %v", propertyValue, err)
				continue
			}
			currentTime := time.Now().UTC()
			if currentTime.After(expiryTime) {
				property.SmartData = "expired"
			} else if itemTypeExists && itemType.ExpiringSoonDays != nil {
				difference := expiryTime.Sub(currentTime)
				days := int16(difference.Hours() / 24)
				if days < *itemTypesDTO[item.ID].ExpiringSoonDays {
					property.SmartData = "expiring_soon"
				}
			}
		}
		item.Properties = append(item.Properties, property)
	}

	for _, row := range derivedNameRows {
		itemIndex, exists := itemIndexes[row.ItemID]
		if !exists {
			continue
		}

		items[itemIndex].DerivedName = row.DerivedName
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

	if err := populateItemDetails(ctx, queriesTx, itemsDTO); err != nil {
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

	heldByID := pgtype.Int8{}
	if req.HeldByID != nil {
		heldByID = pgtype.Int8{Int64: *req.HeldByID, Valid: true}
	}

	sortPropertyID := pgtype.Int8{}
	if req.SortPropertyID != nil {
		sortPropertyID = pgtype.Int8{Int64: *req.SortPropertyID, Valid: true}
	}

	// Both ListItems and SumItemProperties need the unit factor table: the first to compare masses
	// and volumes recorded in different units, the second to add them up.
	unitValueTypes, unitNames, unitFactors := dto.MeasureUnitFactorRows()

	propertyIDs := make([]int64, 0, len(req.PropertyFilters))
	for propertyID := range req.PropertyFilters {
		propertyIDs = append(propertyIDs, propertyID)
	}
	sort.Slice(propertyIDs, func(left, right int) bool { return propertyIDs[left] < propertyIDs[right] })
	propertyValues := make([]json.RawMessage, 0, len(propertyIDs))
	for _, propertyID := range propertyIDs {
		propertyValues = append(propertyValues, req.PropertyFilters[propertyID])
	}

	totalItems, err := db.Queries.CountItems(ctx, repository.CountItemsParams{
		TypeID:         typeID,
		Consumption:    req.Consumption,
		CreatedFrom:    createdFrom,
		CreatedTo:      createdTo,
		Search:         req.Search,
		PropertyIds:    propertyIDs,
		PropertyValues: propertyValues,
		HeldBy:         heldByID,
	})
	if err != nil {
		return dto.ItemsPage{}, err
	}

	items, err := db.Queries.ListItems(ctx, repository.ListItemsParams{
		TypeID:         typeID,
		Consumption:    req.Consumption,
		CreatedFrom:    createdFrom,
		CreatedTo:      createdTo,
		Search:         req.Search,
		PropertyIds:    propertyIDs,
		PropertyValues: propertyValues,
		LimitVal:       req.Limit,
		OffsetVal:      req.Offset,
		OrderAsc:       req.Order == "asc",
		SortPropertyID: sortPropertyID,
		UnitValueTypes: unitValueTypes,
		UnitNames:      unitNames,
		UnitFactors:    unitFactors,
		HeldBy:         heldByID,
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
		if err := populateItemDetails(ctx, db.Queries, itemsDTO); err != nil {
			return dto.ItemsPage{}, err
		}
		if err := populateItemRequestInformation(ctx, itemsDTO, req.ViewerID); err != nil {
			return dto.ItemsPage{}, err
		}
	}

	totalRows, err := db.Queries.SumItemProperties(ctx, repository.SumItemPropertiesParams{
		TypeID:         typeID,
		Consumption:    req.Consumption,
		CreatedFrom:    createdFrom,
		CreatedTo:      createdTo,
		Search:         req.Search,
		PropertyIds:    propertyIDs,
		PropertyValues: propertyValues,
		UnitValueTypes: unitValueTypes,
		UnitNames:      unitNames,
		UnitFactors:    unitFactors,
		HeldBy:         heldByID,
	})
	if err != nil {
		return dto.ItemsPage{}, err
	}

	return dto.ItemsPage{
		Items:          itemsDTO,
		Limit:          req.Limit,
		Offset:         req.Offset,
		Total:          totalItems,
		PropertyTotals: toItemPropertyTotals(totalRows),
	}, nil
}

// toItemPropertyTotals turns the summed rows into the JSON shape the property itself uses, so the
// frontend renders a total the same way it renders a single value. Mass and volume arrive summed in
// their base unit and get scaled up to something readable; prices are already per currency.
func toItemPropertyTotals(rows []repository.SumItemPropertiesRow) []dto.ItemPropertyTotal {
	totals := make([]dto.ItemPropertyTotal, 0, len(rows))

	for _, row := range rows {
		summedAmount, err := strconv.ParseInt(row.TotalAmount, 10, 64)
		if err != nil {
			log.Printf("Couldn't parse %v total %v for property %v. Error: %v",
				row.ValueType, row.TotalAmount, row.PropertyID, err)
			continue
		}

		var value any
		switch row.ValueType {
		case "price":
			value = dto.PTPrice{Amount: summedAmount, Currency: row.Currency}
		case "mass":
			amount, unit := dto.ScaleMeasureToLargestUnit(summedAmount, dto.MassUnitFactors)
			value = dto.PTMass{Amount: amount, Unit: unit}
		case "volume":
			amount, unit := dto.ScaleMeasureToLargestUnit(summedAmount, dto.VolumeUnitFactors)
			value = dto.PTVolume{Amount: amount, Unit: unit}
		default:
			continue
		}

		encodedValue, err := json.Marshal(value)
		if err != nil {
			log.Printf("Couldn't marshal %v total for property %v. Error: %v", row.ValueType, row.PropertyID, err)
			continue
		}

		totals = append(totals, dto.ItemPropertyTotal{
			PropertyID: row.PropertyID,
			ValueType:  row.ValueType,
			Value:      encodedValue,
			ValueCount: row.ValueCount,
		})
	}

	return totals
}

func GetItem(ctx context.Context, id, viewerID int64) (dto.Item, error) {
	item, err := db.Queries.GetItemByID(ctx, id)
	if err != nil {
		return dto.Item{}, err
	}

	itemsDTO := []dto.Item{dto.ToItemDTO(item)}

	if err := populateItemDetails(ctx, db.Queries, itemsDTO); err != nil {
		return dto.Item{}, err
	}
	if err := populateItemRequestInformation(ctx, itemsDTO, viewerID); err != nil {
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

	itemsDTO := []dto.Item{dto.ToItemDTO(item)}
	if err := populateItemDetails(ctx, db.Queries, itemsDTO); err != nil {
		return dto.Item{}, err
	}
	if err := populateItemRequestInformation(ctx, itemsDTO, req.ViewerID); err != nil {
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
