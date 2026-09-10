package service

import (
	"context"
	"encoding/json"
	"errors"
	"slices"
	"strconv"

	"github.com/jackc/pgx/v5"
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

func ListItemTypeFilterableProperties(ctx context.Context, itemTypeID int64) ([]dto.ItemTypeFilterableProperty, error) {
	if _, err := db.Queries.GetItemTypeByID(ctx, itemTypeID); err != nil {
		return nil, err
	}

	properties, err := db.Queries.ListItemTypeFilterableProperties(ctx, itemTypeID)
	if err != nil {
		return nil, err
	}

	result := make([]dto.ItemTypeFilterableProperty, len(properties))
	for index, property := range properties {
		result[index] = dto.ItemTypeFilterableProperty{
			PropertyID: property.PropertyID,
			ValueCount: property.ValueCount,
		}
	}

	return result, nil
}

func ListItemTypePropertyValues(
	ctx context.Context,
	itemTypeID, propertyID int64,
	search string,
	limit int32,
) ([]json.RawMessage, error) {
	if _, err := db.Queries.GetItemTypeByID(ctx, itemTypeID); err != nil {
		return nil, err
	}

	return db.Queries.ListItemTypePropertyValues(ctx, repository.ListItemTypePropertyValuesParams{
		TypeID:     itemTypeID,
		PropertyID: propertyID,
		Search:     search,
		LimitVal:   limit,
	})
}

func ListItemTypes(ctx context.Context, limit, offset int32, search string) (dto.ItemTypesPage, error) {
	total, err := db.Queries.CountItemTypes(ctx, search)
	if err != nil {
		return dto.ItemTypesPage{}, err
	}

	itemTypes, err := db.Queries.ListItemTypes(ctx, repository.ListItemTypesParams{
		LimitVal:  limit,
		OffsetVal: offset,
		Search:    search,
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

	derivedNameFormat := pgtype.Text{String: "", Valid: false}
	if len(req.DerivedNameFormat) > 0 {
		derivedNameFormat = pgtype.Text{String: req.DerivedNameFormat, Valid: true}
	}

	expiringSoonDays := pgtype.Int2{Valid: false}
	if req.ExpiringSoonDays != nil && *req.ExpiringSoonDays > 0 {
		expiringSoonDays = pgtype.Int2{Int16: *req.ExpiringSoonDays, Valid: true}
	}

	itemType, err := queriesTx.CreateItemType(ctx, repository.CreateItemTypeParams{
		Name:              req.Name,
		Description:       description,
		DerivedNameFormat: derivedNameFormat,
		ExpiringSoonDays:  expiringSoonDays,
	})
	if err != nil {
		return dto.ItemType{}, err
	}

	itemTypeDTO := dto.ToItemTypeDTO(itemType)

	// The properties a type is created with are part of the type, so they go in the creation entry's
	// context rather than each getting an item_type_property_add of its own.
	auditProperties := make([]dto.AuditProperty, 0, len(req.Properties))

	if req.Properties != nil {
		for _, propRequest := range req.Properties {
			var property repository.Property
			if propRequest.DefaultValue != nil {
				if property, err = validatePropertyValue(ctx, queriesTx, propRequest.ID, *propRequest.DefaultValue); err != nil {
					return dto.ItemType{}, err
				}
			} else {
				// Without a default value nothing has looked the property up yet, but the entry
				// still needs its name.
				if property, err = queriesTx.GetPropertyByID(ctx, propRequest.ID); err != nil {
					if errors.Is(err, pgx.ErrNoRows) {
						return dto.ItemType{}, ErrInvalidReference
					}
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

			auditProperties = append(auditProperties, dto.AuditProperty{
				ID:        property.ID,
				Name:      property.Name,
				ValueType: property.ValueType,
			})

			propDTO := dto.ToItemTypePropertyDTO(prop)
			itemTypeDTO.Properties = append(itemTypeDTO.Properties, propDTO)
		}
	}

	changes := auditNew(nil, "name", dto.AuditValueTypeText, itemType.Name)
	if itemType.Description.Valid {
		changes = auditNew(changes, "description", dto.AuditValueTypeText, itemType.Description.String)
	}
	if itemType.DerivedNameFormat.Valid {
		changes = auditNew(changes, "derived_name_format", dto.AuditValueTypeText, itemType.DerivedNameFormat.String)
	}
	if itemType.ExpiringSoonDays.Valid {
		changes = auditNew(changes, "expiring_soon_days", dto.AuditValueTypeCount, itemType.ExpiringSoonDays.Int16)
	}

	if err := writeAudit(ctx, queriesTx, repository.AuditActionItemTypeCreate, repository.AuditTargetTypeItemType,
		itemType.ID, itemType.Name, changes, dto.AuditContext{Properties: auditProperties}); err != nil {
		return dto.ItemType{}, err
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

		// The UPDATE statements below return only the new row, so the prior state has to be read
		// first or the audit entry has nothing to diff against.
		before, err := queriesTx.GetItemTypeByID(ctx, req.ID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return dto.ItemType{}, ErrNotFound
			}
			return dto.ItemType{}, err
		}
		itemType = before

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

		if req.ExpiringSoonDays != nil {
			expiringSoonDays := pgtype.Int2{Valid: false}
			if *req.ExpiringSoonDays > 0 {
				expiringSoonDays = pgtype.Int2{Int16: *req.ExpiringSoonDays, Valid: true}
			}

			var err error
			if itemType, err = queriesTx.UpdateItemType_ExpiringSoonDays(ctx, repository.UpdateItemType_ExpiringSoonDaysParams{
				ID:               req.ID,
				ExpiringSoonDays: expiringSoonDays,
			}); err != nil {
				return dto.ItemType{}, err
			}
		}

		changes := auditDiff(nil, "name", dto.AuditValueTypeText, before.Name, itemType.Name)
		changes = auditDiff(changes, "description", dto.AuditValueTypeText,
			before.Description.String, itemType.Description.String)
		changes = auditDiff(changes, "derived_name_format", dto.AuditValueTypeText,
			before.DerivedNameFormat.String, itemType.DerivedNameFormat.String)
		changes = auditDiff(changes, "expiring_soon_days", dto.AuditValueTypeCount,
			before.ExpiringSoonDays.Int16, itemType.ExpiringSoonDays.Int16)

		if len(changes) > 0 {
			if err := writeAudit(ctx, queriesTx, repository.AuditActionItemTypeUpdate, repository.AuditTargetTypeItemType,
				itemType.ID, itemType.Name, changes, dto.AuditContext{}); err != nil {
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
	tx, err := db.BeginTransaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	queriesTx := db.Queries.WithTx(tx)

	// Read before the delete, so the entry names the type rather than just its id.
	existing, err := queriesTx.GetItemTypeByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	rowsAffected, err := queriesTx.DeleteItemType(ctx, id)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	// Identity only, like every other deletion: the name it went by, and who removed it. What it was
	// configured with is on its own timeline, which outlives it.
	if err := writeAudit(ctx, queriesTx, repository.AuditActionItemTypeDelete, repository.AuditTargetTypeItemType,
		id, existing.Name, nil, dto.AuditContext{}); err != nil {
		return err
	}

	return tx.Commit(ctx)
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

	tx, err := db.BeginTransaction(ctx)
	if err != nil {
		return dto.ItemTypeProperty{}, err
	}
	defer tx.Rollback(ctx)
	queriesTx := db.Queries.WithTx(tx)

	// validatePropertyValue only runs when a default value was given, so the property is looked up
	// unconditionally here - the entry needs its name either way.
	property, err := queriesTx.GetPropertyByID(ctx, req.PropertyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.ItemTypeProperty{}, ErrInvalidReference
		}
		return dto.ItemTypeProperty{}, err
	}

	if req.DefaultValue != nil {
		if _, err := validatePropertyValue(ctx, queriesTx, req.PropertyID, *req.DefaultValue); err != nil {
			return dto.ItemTypeProperty{}, err
		}
	}

	visibility := repository.PropertyVisibilityOverview
	if req.Visibility != nil {
		visibility = *req.Visibility
	}

	typeProp, err := queriesTx.AddItemTypeProperty(ctx, repository.AddItemTypePropertyParams{
		TypeID:       req.TypeID,
		PropertyID:   req.PropertyID,
		DefaultValue: req.DefaultValue,
		Visibility:   visibility,
	})
	if err != nil {
		return dto.ItemTypeProperty{}, err
	}

	itemType, err := queriesTx.GetItemTypeByID(ctx, req.TypeID)
	if err != nil {
		return dto.ItemTypeProperty{}, err
	}

	changes := auditNew(nil, "visibility", dto.AuditValueTypeVisibility, string(visibility))
	if typeProp.DefaultValue != nil {
		changes = append(changes, dto.AuditChange{
			Key:       "default_value",
			ValueType: property.ValueType,
			New:       *typeProp.DefaultValue,
		})
	}

	if err := writeAudit(ctx, queriesTx, repository.AuditActionItemTypePropertyAdd, repository.AuditTargetTypeItemType,
		req.TypeID, itemType.Name, changes, dto.AuditContext{
			PropertyID:   &property.ID,
			PropertyName: property.Name,
		}); err != nil {
		return dto.ItemTypeProperty{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return dto.ItemTypeProperty{}, err
	}

	return dto.ToItemTypePropertyDTO(typeProp), nil
}

func UpdateItemTypeProperty(ctx context.Context, req dto.AddUpdateItemTypePropertyRequest) (dto.ItemTypeProperty, error) {
	if err := dto.Validate(req); err != nil {
		return dto.ItemTypeProperty{}, err
	}

	if req.DefaultValue == nil && req.Visibility == nil {
		return dto.ItemTypeProperty{}, ErrNoUpdateFields
	}

	tx, err := db.BeginTransaction(ctx)
	if err != nil {
		return dto.ItemTypeProperty{}, err
	}
	defer tx.Rollback(ctx)
	queriesTx := db.Queries.WithTx(tx)

	property, err := queriesTx.GetPropertyByID(ctx, req.PropertyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.ItemTypeProperty{}, ErrInvalidReference
		}
		return dto.ItemTypeProperty{}, err
	}

	// Both UPDATEs return only the new row, so the pre-image is read first.
	before, err := itemTypeProperty(ctx, queriesTx, req.TypeID, req.PropertyID)
	if err != nil {
		return dto.ItemTypeProperty{}, err
	}
	typeProp := before

	if req.DefaultValue != nil {
		if _, err := validatePropertyValue(ctx, queriesTx, req.PropertyID, *req.DefaultValue); err != nil {
			return dto.ItemTypeProperty{}, err
		}

		var err error
		typeProp, err = queriesTx.UpdateItemTypeProperty_DefaultValue(ctx, repository.UpdateItemTypeProperty_DefaultValueParams{
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
		typeProp, err = queriesTx.UpdateItemTypeProperty_Visibility(ctx, repository.UpdateItemTypeProperty_VisibilityParams{
			TypeID:     req.TypeID,
			PropertyID: req.PropertyID,
			Visibility: *req.Visibility,
		})
		if err != nil {
			return dto.ItemTypeProperty{}, err
		}
	}

	changes := auditDiff(nil, "visibility", dto.AuditValueTypeVisibility,
		string(before.Visibility), string(typeProp.Visibility))
	changes = auditRawDiff(changes, "default_value", property.ValueType,
		optionalRawJSON(before.DefaultValue), optionalRawJSON(typeProp.DefaultValue))

	if len(changes) > 0 {
		itemType, err := queriesTx.GetItemTypeByID(ctx, req.TypeID)
		if err != nil {
			return dto.ItemTypeProperty{}, err
		}

		if err := writeAudit(ctx, queriesTx, repository.AuditActionItemTypePropertyUpdate, repository.AuditTargetTypeItemType,
			req.TypeID, itemType.Name, changes, dto.AuditContext{
				PropertyID:   &property.ID,
				PropertyName: property.Name,
			}); err != nil {
			return dto.ItemTypeProperty{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return dto.ItemTypeProperty{}, err
	}

	return dto.ToItemTypePropertyDTO(typeProp), nil
}

// itemTypeProperty reads one property's settings on a type, for the before half of a diff.
func itemTypeProperty(ctx context.Context, q repository.Querier, typeID, propertyID int64) (repository.ItemTypeProperty, error) {
	rows, err := q.GetItemTypeProperties(ctx, []int64{typeID})
	if err != nil {
		return repository.ItemTypeProperty{}, err
	}

	for _, row := range rows {
		if row.PropertyID == propertyID {
			return row, nil
		}
	}

	return repository.ItemTypeProperty{}, ErrNotFound
}

func ReorderItemTypeProperties(ctx context.Context, req dto.ReorderItemTypePropertiesRequest) error {
	if err := dto.Validate(req); err != nil {
		return err
	}

	tx, err := db.BeginTransaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	queriesTx := db.Queries.WithTx(tx)

	// The lock already returns the whole row, so keeping it saves looking the name up again for the
	// audit entry.
	itemType, err := queriesTx.LockItemType(ctx, req.TypeID)
	if err != nil {
		return err
	}

	typeProperties, err := queriesTx.GetItemTypeProperties(ctx, []int64{req.TypeID})
	if err != nil {
		return err
	}
	if len(typeProperties) != len(req.PropertyIDs) {
		return ErrInvalidItemTypePropertyOrder
	}

	existingPropertyIDs := make(map[int64]struct{}, len(typeProperties))
	for _, typeProperty := range typeProperties {
		existingPropertyIDs[typeProperty.PropertyID] = struct{}{}
	}
	for _, propertyID := range req.PropertyIDs {
		if _, exists := existingPropertyIDs[propertyID]; !exists {
			return ErrInvalidItemTypePropertyOrder
		}
		delete(existingPropertyIDs, propertyID)
	}
	if len(existingPropertyIDs) != 0 {
		return ErrInvalidItemTypePropertyOrder
	}

	if err := queriesTx.OffsetItemTypePropertyPositions(ctx, req.TypeID); err != nil {
		return err
	}

	rowsAffected, err := queriesTx.SetItemTypePropertyPositions(ctx, repository.SetItemTypePropertyPositionsParams{
		TypeID:      req.TypeID,
		PropertyIds: req.PropertyIDs,
	})

	if err != nil {
		return err
	}
	if rowsAffected != int64(len(req.PropertyIDs)) {
		return ErrInvalidItemTypePropertyOrder
	}

	// typeProperties came back in position order, so it is the old arrangement; req.PropertyIDs is
	// the new one. Recorded as names, because a list of ids says nothing to whoever reads the log.
	previousIDs := make([]int64, len(typeProperties))
	for index, typeProperty := range typeProperties {
		previousIDs[index] = typeProperty.PropertyID
	}

	propertyNames, err := propertyNamesByID(ctx, queriesTx, req.PropertyIDs)
	if err != nil {
		return err
	}

	// A drag that ended where it started is not worth an entry.
	if !slices.Equal(previousIDs, req.PropertyIDs) {
		changes := []dto.AuditChange{{
			Key:       "order",
			ValueType: dto.AuditValueTypeOrder,
			Old:       mustJSON(namesInOrder(previousIDs, propertyNames)),
			New:       mustJSON(namesInOrder(req.PropertyIDs, propertyNames)),
		}}

		if err := writeAudit(ctx, queriesTx, repository.AuditActionItemTypePropertyReorder,
			repository.AuditTargetTypeItemType, req.TypeID, itemType.Name, changes, dto.AuditContext{}); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// propertyNamesByID resolves a set of property ids to their names in one round trip.
func propertyNamesByID(ctx context.Context, q repository.Querier, propertyIDs []int64) (map[int64]string, error) {
	rows, err := q.GetPropertiesByIDs(ctx, propertyIDs)
	if err != nil {
		return nil, err
	}

	names := make(map[int64]string, len(rows))
	for _, row := range rows {
		names[row.ID] = row.Name
	}

	return names, nil
}

// namesInOrder maps an ordering of property ids onto their names, falling back to the id for a
// property that has since gone away.
func namesInOrder(propertyIDs []int64, names map[int64]string) []string {
	ordered := make([]string, len(propertyIDs))
	for index, propertyID := range propertyIDs {
		name, ok := names[propertyID]
		if !ok {
			name = strconv.FormatInt(propertyID, 10)
		}
		ordered[index] = name
	}

	return ordered
}

// mustJSON encodes a value that cannot fail to encode - a slice of strings - so the call sites stay
// readable. A failure yields null rather than panicking, since a half-rendered audit entry still
// beats taking the request down with it.
func mustJSON(value any) json.RawMessage {
	encoded, err := json.Marshal(value)
	if err != nil {
		return json.RawMessage("null")
	}

	return encoded
}

func RemoveItemTypeProperty(ctx context.Context, typeId int64, propId int64) error {
	tx, err := db.BeginTransaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	queriesTx := db.Queries.WithTx(tx)

	itemType, err := queriesTx.GetItemTypeByID(ctx, typeId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	property, err := queriesTx.GetPropertyByID(ctx, propId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	if derivedNameUsesProperty(itemType.DerivedNameFormat.String, property.Name) {
		return ErrInvalidDerivedNameFormat
	}

	rowsAffected, err := queriesTx.RemoveItemTypeProperty(ctx, repository.RemoveItemTypePropertyParams{
		TypeID:     typeId,
		PropertyID: propId,
	})
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	if err := writeAudit(ctx, queriesTx, repository.AuditActionItemTypePropertyRemove,
		repository.AuditTargetTypeItemType, typeId, itemType.Name, nil, dto.AuditContext{
			PropertyID:   &property.ID,
			PropertyName: property.Name,
		}); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
