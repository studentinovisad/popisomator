package service

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/studentinovisad/popisomator/backend/internal/db"
	"github.com/studentinovisad/popisomator/backend/internal/dto"
	"github.com/studentinovisad/popisomator/backend/internal/repository"
)

// ListItemStock breaks an item type down into stock groups, marking the ones at or below its
// threshold. A type with no threshold still reports its groups; none of them are ever low.
func ListItemStock(ctx context.Context, typeID int64) (dto.ItemTypeStock, error) {
	itemType, err := db.Queries.GetItemTypeByID(ctx, typeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.ItemTypeStock{}, ErrNotFound
		}
		return dto.ItemTypeStock{}, err
	}

	rows, err := db.Queries.ListStockGroups(ctx, pgtype.Int8{Int64: typeID, Valid: true})
	if err != nil {
		return dto.ItemTypeStock{}, err
	}

	stock := dto.ItemTypeStock{
		TypeID:   itemType.ID,
		TypeName: itemType.Name,
		Groups:   make([]dto.StockGroup, len(rows)),
	}
	if itemType.LowStockCount.Valid {
		stock.Threshold = &itemType.LowStockCount.Int32
	}

	for index, row := range rows {
		stock.Groups[index] = dto.StockGroup{
			Name:         stockGroupLabel(row.GroupName, itemType.Name),
			InStockCount: row.InStockCount,
			TotalCount:   row.TotalCount,
			Low:          isLowStock(row.InStockCount, itemType.LowStockCount.Int32, itemType.LowStockCount.Valid),
		}
	}

	return stock, nil
}

// ReconcileLowStock evaluates every group of an item type. Configuration changes use the full pass:
// changing a threshold or the derived-name format can affect every group at once.
func ReconcileLowStock(ctx context.Context, typeID int64) error {
	itemType, err := db.Queries.GetItemTypeByID(ctx, typeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	if !itemType.LowStockCount.Valid {
		_, err := db.Queries.ClearLowStockAlerts(ctx, typeID)
		return err
	}

	rows, err := db.Queries.ListStockGroups(ctx, pgtype.Int8{Int64: typeID, Valid: true})
	if err != nil {
		return err
	}

	groups := make([]lowStockGroup, len(rows))
	groupNames := make(map[string]struct{}, len(rows))
	for index, row := range rows {
		groups[index] = lowStockGroup{
			Name:         row.GroupName,
			InStockCount: row.InStockCount,
			TotalCount:   row.TotalCount,
		}
		groupNames[row.GroupName] = struct{}{}
	}

	if err := reconcileLowStockGroups(ctx, itemType, groups); err != nil {
		return err
	}

	alerts, err := db.Queries.ListLowStockAlerts(ctx, typeID)
	if err != nil {
		return err
	}
	staleGroupNames := make([]string, 0)
	for _, alert := range alerts {
		if _, exists := groupNames[alert.GroupName]; !exists {
			staleGroupNames = append(staleGroupNames, alert.GroupName)
		}
	}
	if len(staleGroupNames) == 0 {
		return nil
	}

	_, err = db.Queries.DeleteLowStockAlerts(ctx, repository.DeleteLowStockAlertsParams{
		TypeID:     typeID,
		GroupNames: staleGroupNames,
	})
	return err
}

// ReconcileLowStockGroups evaluates only named groups. Item writes pass the group before and after
// the write, so a renamed or deleted item still clears an alert for the group it left.
func ReconcileLowStockGroups(ctx context.Context, typeID int64, groupNames []string) error {
	groupNames = uniqueGroupNames(groupNames)
	if len(groupNames) == 0 {
		return nil
	}

	itemType, err := db.Queries.GetItemTypeByID(ctx, typeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	if !itemType.LowStockCount.Valid {
		return nil
	}

	rows, err := db.Queries.GroupItemCountsForGroups(ctx, repository.GroupItemCountsForGroupsParams{
		TypeID:     typeID,
		GroupNames: groupNames,
	})
	if err != nil {
		return err
	}

	groups := make([]lowStockGroup, len(rows))
	for index, row := range rows {
		groups[index] = lowStockGroup{
			Name:         row.GroupName,
			InStockCount: row.InStockCount,
			TotalCount:   row.TotalCount,
		}
	}

	return reconcileLowStockGroups(ctx, itemType, groups)
}

type lowStockGroup struct {
	Name         string
	InStockCount int64
	TotalCount   int64
}

func reconcileLowStockGroups(ctx context.Context, itemType repository.ItemType, groups []lowStockGroup) error {
	if len(groups) == 0 {
		return nil
	}

	groupNames := make([]string, len(groups))
	for index, group := range groups {
		groupNames[index] = group.Name
	}
	alerts, err := db.Queries.ListLowStockAlertsForGroups(ctx, repository.ListLowStockAlertsForGroupsParams{
		TypeID:     itemType.ID,
		GroupNames: groupNames,
	})
	if err != nil {
		return err
	}

	alerted := make(map[string]struct{}, len(alerts))
	for _, alert := range alerts {
		alerted[alert.GroupName] = struct{}{}
	}

	recoveredGroupNames := make([]string, 0)
	for _, group := range groups {
		_, wasAlerted := alerted[group.Name]
		if group.TotalCount == 0 || !isLowStock(group.InStockCount, itemType.LowStockCount.Int32, true) {
			if wasAlerted {
				recoveredGroupNames = append(recoveredGroupNames, group.Name)
			}
			continue
		}
		if wasAlerted {
			continue
		}

		if _, err := createLowStockAlert(ctx, itemType.ID, group.Name,
			stockGroupLabel(group.Name, itemType.Name), itemType.LowStockCount.Int32, int32(group.InStockCount)); err != nil {
			return err
		}
	}

	if len(recoveredGroupNames) == 0 {
		return nil
	}
	_, err = db.Queries.DeleteLowStockAlerts(ctx, repository.DeleteLowStockAlertsParams{
		TypeID:     itemType.ID,
		GroupNames: recoveredGroupNames,
	})
	return err
}

func reconcileLowStockAfterTypeChange(ctx context.Context, typeID int64) {
	if err := ReconcileLowStock(ctx, typeID); err != nil {
		log.Printf("Couldn't reconcile low stock for item type %v. Error: %v", typeID, err)
	}
}

func reconcileLowStockAfterItemChange(ctx context.Context, typeID int64, groupNames ...string) {
	if err := ReconcileLowStockGroups(ctx, typeID, groupNames); err != nil {
		log.Printf("Couldn't reconcile low stock for item type %v. Error: %v", typeID, err)
	}
}

func reconcileLowStockAfterItemRequestChange(ctx context.Context, itemID int64) {
	item, err := db.Queries.GetItemByID(ctx, itemID)
	if err != nil {
		log.Printf("Couldn't read item %v to reconcile low stock. Error: %v", itemID, err)
		return
	}

	groupName, err := itemDerivedName(ctx, db.Queries, itemID)
	if err != nil {
		log.Printf("Couldn't read item %v's derived name to reconcile low stock. Error: %v", itemID, err)
		return
	}

	reconcileLowStockAfterItemChange(ctx, item.TypeID, groupName)
}

func createLowStockAlert(
	ctx context.Context,
	typeID int64,
	groupName, groupLabel string,
	threshold, observed int32,
) ([]int64, error) {
	tx, err := db.BeginTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	queriesTx := db.Queries.WithTx(tx)

	claimed, err := queriesTx.InsertLowStockAlert(ctx, repository.InsertLowStockAlertParams{
		TypeID:    typeID,
		GroupName: groupName,
	})
	if err != nil || claimed == 0 {
		return nil, err
	}

	notificationIDs, err := CreateLowStockNotifications(ctx, queriesTx,
		typeID, groupLabel, threshold, observed)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return notificationIDs, nil
}

func uniqueGroupNames(groupNames []string) []string {
	unique := make([]string, 0, len(groupNames))
	seen := make(map[string]struct{}, len(groupNames))
	for _, groupName := range groupNames {
		if _, exists := seen[groupName]; exists {
			continue
		}
		seen[groupName] = struct{}{}
		unique = append(unique, groupName)
	}

	return unique
}

func isLowStock(inStockCount int64, threshold int32, hasThreshold bool) bool {
	return hasThreshold && inStockCount <= int64(threshold)
}

func stockGroupLabel(groupName, typeName string) string {
	if groupName == "" {
		return typeName
	}

	return groupName
}
