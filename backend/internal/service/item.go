package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
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

	// The initial property values ride along on each item's own creation entry rather than getting
	// item_property_add entries of their own: they are part of what was created, not later edits to it.
	initialChanges := make([]dto.AuditChange, 0, len(req.Properties))

	if req.Properties != nil {
		for _, propRequest := range req.Properties {
			property, err := validatePropertyValue(ctx, queriesTx, propRequest.ID, propRequest.Value)
			if err != nil {
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

			initialChanges = append(initialChanges, dto.AuditChange{
				Key:       "property",
				Label:     property.Name,
				ValueType: property.ValueType,
				New:       propRequest.Value,
			})
		}
	}

	if err := populateItemDetails(ctx, queriesTx, itemsDTO); err != nil {
		return nil, err
	}

	itemType, err := queriesTx.GetItemTypeByID(ctx, req.TypeID)
	if err != nil {
		// The insert above would already have failed on a bad type_id, so this is a genuine lookup
		// failure rather than a bad request.
		return nil, err
	}

	// One entry per item, not one for the batch: an item whose own timeline does not start with its
	// creation is the surface this exists for.
	auditTargets := make([]auditTarget, len(itemsDTO))
	for index, item := range itemsDTO {
		auditTargets[index] = auditTarget{
			ID:      item.ID,
			Label:   item.DerivedName,
			Changes: initialChanges,
			Context: dto.AuditContext{
				TypeID:    &req.TypeID,
				TypeName:  itemType.Name,
				BatchSize: int32(len(itemsDTO)),
			},
		}
	}
	if err := writeAuditBulk(ctx, queriesTx, repository.AuditActionItemCreate, repository.AuditTargetTypeItem, auditTargets); err != nil {
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

		// The UPDATE statements below only return the new row, so the prior state has to be read
		// first or the audit entry has nothing to diff against.
		before, err := queriesTx.GetItemByID(ctx, req.ID)
		if err != nil {
			return dto.Item{}, err
		}

		// PATCH /items/{id} and POST /items/{id}/consume both land here. Today their controllers hand
		// over disjoint fields - the first can only carry a type, the second only a consumption - but
		// the service does not require that, so each kind of change is recorded under its own action
		// rather than one action being picked for the call. Setting both yields both entries instead
		// of a type change mislabelled as a consumption.
		var typeChanges, consumptionChanges []dto.AuditChange

		if req.TypeID != nil {
			var err error
			if item, err = queriesTx.UpdateItem_Type(ctx, repository.UpdateItem_TypeParams{
				ID:     req.ID,
				TypeID: *req.TypeID,
			}); err != nil {
				return dto.Item{}, err
			}

			if before.TypeID != item.TypeID {
				oldType, err := queriesTx.GetItemTypeByID(ctx, before.TypeID)
				if err != nil {
					return dto.Item{}, err
				}
				newType, err := queriesTx.GetItemTypeByID(ctx, item.TypeID)
				if err != nil {
					return dto.Item{}, err
				}

				typeChanges = auditDiff(typeChanges, "type_id", dto.AuditValueTypeReference,
					auditReference{ID: oldType.ID, Name: oldType.Name},
					auditReference{ID: newType.ID, Name: newType.Name})
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

			consumptionChanges = auditDiff(consumptionChanges, "consumption", dto.AuditValueTypeConsumption,
				string(before.Consumption), string(item.Consumption))
		}

		// An update that set every field to what it already held is not worth an entry, so the name
		// is only resolved once something actually moved.
		if len(typeChanges) > 0 || len(consumptionChanges) > 0 {
			label, err := itemDerivedName(ctx, queriesTx, req.ID)
			if err != nil {
				return dto.Item{}, err
			}

			for _, entry := range []struct {
				action  repository.AuditAction
				changes []dto.AuditChange
			}{
				{repository.AuditActionItemUpdate, typeChanges},
				{repository.AuditActionItemConsume, consumptionChanges},
			} {
				if len(entry.changes) == 0 {
					continue
				}
				if err := writeAudit(ctx, queriesTx, entry.action, repository.AuditTargetTypeItem,
					req.ID, label, entry.changes, dto.AuditContext{}); err != nil {
					return dto.Item{}, err
				}
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
	tx, err := db.BeginTransaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	queriesTx := db.Queries.WithTx(tx)

	// Take the same lock every path that touches this item's requests takes. Without it a request
	// created between reading them below and the DELETE would cascade away unrecorded, which is
	// exactly the claim this function exists to name.
	if _, err := queriesTx.LockItemForRequest(ctx, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	// Everything the audit entries need has to be read before the DELETE. An item's name is derived
	// from its properties, and those cascade away with it, so afterwards there is nothing left to
	// identify the item by.
	//
	// What the item held is deliberately not snapshotted here: every value it ever carried is already
	// on its own timeline, which outlives it, so recording the final state again would only duplicate
	// entries that are still queryable.
	item, err := queriesTx.GetItemByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	derivedNames, err := queriesTx.GetItemsDerivedNames(ctx, []int64{id})
	if err != nil {
		return err
	}
	label := ""
	if len(derivedNames) > 0 {
		label = derivedNames[0].DerivedName
	}

	itemType, err := queriesTx.GetItemTypeByID(ctx, item.TypeID)
	if err != nil {
		return err
	}

	// Requests standing against this item cascade away with it, silently cancelling someone's claim.
	// Read them while they still exist so each person who loses one is named.
	requests, err := queriesTx.ListItemRequestsForItem(ctx, id)
	if err != nil {
		return err
	}

	rowsAffected, err := queriesTx.DeleteItem(ctx, id)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	// Written before the deletion entry so that, sharing one now(), the id tiebreaker leaves the
	// deletion at the top of the feed with the claims it cancelled underneath.
	cascadeTargets := make([]auditTarget, len(requests))
	for index, request := range requests {
		userID := request.UserID
		cascadeTargets[index] = auditTarget{
			ID:    id,
			Label: label,
			Context: dto.AuditContext{
				SubjectUserID:   &userID,
				SubjectUserName: request.UserName,
				Reason:          request.Reason,
				RequestStatus:   string(request.Status),
				Cascaded:        true,
			},
		}
	}
	if err := writeAuditBulk(ctx, queriesTx, repository.AuditActionItemRequestDelete,
		repository.AuditTargetTypeItem, cascadeTargets); err != nil {
		return err
	}

	if err := writeAudit(ctx, queriesTx, repository.AuditActionItemDelete, repository.AuditTargetTypeItem,
		id, label, nil, dto.AuditContext{
			TypeID:   &item.TypeID,
			TypeName: itemType.Name,
		}); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func AddItemProperty(ctx context.Context, req dto.AddUpdateItemPropertyRequest) (dto.ItemProperty, error) {
	if err := dto.Validate(req); err != nil {
		return dto.ItemProperty{}, err
	}

	tx, err := db.BeginTransaction(ctx)
	if err != nil {
		return dto.ItemProperty{}, err
	}
	defer tx.Rollback(ctx)
	queriesTx := db.Queries.WithTx(tx)

	property, err := validatePropertyValue(ctx, queriesTx, req.PropertyID, req.Value)
	if err != nil {
		return dto.ItemProperty{}, err
	}

	itemProp, err := queriesTx.AddItemProperty(ctx, repository.AddItemPropertyParams{
		ItemID:        req.ItemID,
		PropertyID:    req.PropertyID,
		PropertyValue: req.Value,
	})
	if err != nil {
		return dto.ItemProperty{}, err
	}

	// Read the name after the insert: the value just added may itself be part of the derived name.
	label, err := itemDerivedName(ctx, queriesTx, req.ItemID)
	if err != nil {
		return dto.ItemProperty{}, err
	}

	changes := auditNew(nil, "property", property.ValueType, req.Value)
	if err := writeAudit(ctx, queriesTx, repository.AuditActionItemPropertyAdd, repository.AuditTargetTypeItem,
		req.ItemID, label, changes, dto.AuditContext{
			PropertyID:   &property.ID,
			PropertyName: property.Name,
		}); err != nil {
		return dto.ItemProperty{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return dto.ItemProperty{}, err
	}

	return dto.ToItemPropertyDTO(itemProp), nil
}

func UpdateItemProperty(ctx context.Context, req dto.AddUpdateItemPropertyRequest) (dto.ItemProperty, error) {
	if err := dto.Validate(req); err != nil {
		return dto.ItemProperty{}, err
	}

	tx, err := db.BeginTransaction(ctx)
	if err != nil {
		return dto.ItemProperty{}, err
	}
	defer tx.Rollback(ctx)
	queriesTx := db.Queries.WithTx(tx)

	property, err := validatePropertyValue(ctx, queriesTx, req.PropertyID, req.Value)
	if err != nil {
		return dto.ItemProperty{}, err
	}

	// The update returns only the new value, so the old one has to be read first.
	previousValue, err := itemPropertyValue(ctx, queriesTx, req.ItemID, req.PropertyID)
	if err != nil {
		return dto.ItemProperty{}, err
	}

	itemProp, err := queriesTx.UpdateItemProperty(ctx, repository.UpdateItemPropertyParams{
		ItemID:        req.ItemID,
		PropertyID:    req.PropertyID,
		PropertyValue: req.Value,
	})
	if err != nil {
		return dto.ItemProperty{}, err
	}

	// The UPDATE succeeds whether or not the value moved, so an edit that set a property to what it
	// already held would otherwise fill the log with entries recording nothing.
	changes := auditRawDiff(nil, "property", property.ValueType, previousValue, itemProp.PropertyValue)
	if len(changes) > 0 {
		label, err := itemDerivedName(ctx, queriesTx, req.ItemID)
		if err != nil {
			return dto.ItemProperty{}, err
		}

		if err := writeAudit(ctx, queriesTx, repository.AuditActionItemPropertyUpdate, repository.AuditTargetTypeItem,
			req.ItemID, label, changes, dto.AuditContext{
				PropertyID:   &property.ID,
				PropertyName: property.Name,
			}); err != nil {
			return dto.ItemProperty{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return dto.ItemProperty{}, err
	}

	return dto.ToItemPropertyDTO(itemProp), nil
}

func RemoveItemProperty(ctx context.Context, itemId int64, propId int64) error {
	tx, err := db.BeginTransaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	queriesTx := db.Queries.WithTx(tx)

	property, err := queriesTx.GetPropertyByID(ctx, propId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	// Both reads happen before the delete: the value is about to be gone, and removing it can change
	// the derived name, so the name recorded should be the one the item was known by at the time.
	removedValue, err := itemPropertyValue(ctx, queriesTx, itemId, propId)
	if err != nil {
		return err
	}

	label, err := itemDerivedName(ctx, queriesTx, itemId)
	if err != nil {
		return err
	}

	rowsAffected, err := queriesTx.RemoveItemProperty(ctx, repository.RemoveItemPropertyParams{
		ItemID:     itemId,
		PropertyID: propId,
	})
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	changes := []dto.AuditChange{{
		Key:       "property",
		ValueType: property.ValueType,
		Old:       removedValue,
	}}
	if err := writeAudit(ctx, queriesTx, repository.AuditActionItemPropertyRemove, repository.AuditTargetTypeItem,
		itemId, label, changes, dto.AuditContext{
			PropertyID:   &property.ID,
			PropertyName: property.Name,
		}); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// itemDerivedName is the name an item goes by, for an audit entry's target_label. A missing row
// yields an empty label rather than an error: an item with no derived name format still has a
// history worth recording, and the reader falls back to the id.
func itemDerivedName(ctx context.Context, q repository.Querier, itemID int64) (string, error) {
	derivedNames, err := q.GetItemsDerivedNames(ctx, []int64{itemID})
	if err != nil {
		return "", err
	}
	if len(derivedNames) == 0 {
		return "", nil
	}

	return derivedNames[0].DerivedName, nil
}

// itemPropertyValue reads one property's current value off an item, for the before half of a diff.
// A property that is not set yet is not an error - it is simply an add rather than an edit.
func itemPropertyValue(ctx context.Context, q repository.Querier, itemID, propertyID int64) (json.RawMessage, error) {
	rows, err := q.GetItemProperties(ctx, []int64{itemID})
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		if row.ItemProperty.PropertyID == propertyID {
			return row.ItemProperty.PropertyValue, nil
		}
	}

	return nil, nil
}
