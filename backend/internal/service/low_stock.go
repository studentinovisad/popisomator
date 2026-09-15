package service

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
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

	rows, err := db.Queries.GroupItemCounts(ctx, typeID)
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

// EvaluateLowStock warns about every group of typeID that has just fallen to its threshold.
//
// Stock only moves when items do, so running this from the mutations that move them is exact and
// needs no scheduler. What it must not do is warn twice about the same shortage: an item edited while
// a group sits below its threshold would otherwise send the warning again. low_stock_alerts holds the
// groups already warned about, so a notification is tied to the crossing rather than to the state,
// and clearing a row on recovery is what makes the next dip notifiable.
//
// The alert row is claimed with an ON CONFLICT DO NOTHING insert and the notification is written only
// if that insert took. Two evaluations racing on the same group therefore produce one warning between
// them, whichever gets there first.
func EvaluateLowStock(ctx context.Context, typeID int64) error {
	itemType, err := db.Queries.GetItemTypeByID(ctx, typeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	// Turning the warning off stops tracking rather than freezing it: leaving the rows behind would
	// silence the groups that were low at the time if it were ever turned back on.
	if !itemType.LowStockCount.Valid {
		_, err := db.Queries.ClearLowStockAlerts(ctx, typeID)
		return err
	}
	threshold := itemType.LowStockCount.Int32

	groups, err := db.Queries.GroupItemCounts(ctx, typeID)
	if err != nil {
		return err
	}

	alerted, err := db.Queries.ListLowStockAlerts(ctx, typeID)
	if err != nil {
		return err
	}
	stale := make(map[string]struct{}, len(alerted))
	for _, alert := range alerted {
		stale[alert.GroupName] = struct{}{}
	}

	// Read on first need rather than up front: this runs after every item write, and the common case
	// is a type where nothing has crossed its threshold and no warning is going out at all.
	var recipientIDs []int64
	recipientsRead := false

	recovered := make([]string, 0, len(alerted))
	for _, group := range groups {
		_, wasAlerted := stale[group.GroupName]
		delete(stale, group.GroupName)

		if !isLowStock(group.InStockCount, threshold, true) {
			// Only groups carrying an alert row are worth naming in the delete; the rest were never
			// low and have nothing to clear.
			if wasAlerted {
				recovered = append(recovered, group.GroupName)
			}
			continue
		}

		if wasAlerted {
			continue
		}

		if !recipientsRead {
			if recipientIDs, err = notificationRecipients(ctx); err != nil {
				return err
			}
			recipientsRead = true
		}

		if _, err := CreateLowStockNotifications(ctx, recipientIDs, itemType,
			stockGroupLabel(group.GroupName, itemType.Name), threshold, int32(group.InStockCount)); err != nil {
			return err
		}
	}

	// Whatever is left has an alert row but no group: every item carrying that name was deleted
	// outright. There is nothing to be short of any more, so the row goes with them.
	for groupName := range stale {
		recovered = append(recovered, groupName)
	}

	if len(recovered) > 0 {
		if _, err := db.Queries.DeleteLowStockAlerts(ctx, repository.DeleteLowStockAlertsParams{
			TypeID:     typeID,
			GroupNames: recovered,
		}); err != nil {
			return err
		}
	}

	return nil
}

// notificationRecipients is who hears about something the system noticed on its own, rather than
// about a request they made themselves.
func notificationRecipients(ctx context.Context) ([]int64, error) {
	users, err := db.Queries.GetUsersByRoles(ctx, repository.GetUsersByRolesParams{
		Roles:        []string{string(repository.UserRoleManager), string(repository.UserRoleAdmin)},
		StatusFilter: string(repository.UserStatusActive),
	})
	if err != nil {
		return nil, err
	}

	recipientIDs := make([]int64, len(users))
	for index, user := range users {
		recipientIDs[index] = user.ID
	}

	return recipientIDs, nil
}

// EvaluateLowStockAsync runs EvaluateLowStock for a type whose stock a just-committed write may have
// moved, and reports failure to the log rather than to the caller. A warning that could not be raised
// is worth knowing about, but it is not worth failing a write the user already completed.
func EvaluateLowStockAsync(ctx context.Context, typeID int64) {
	if err := EvaluateLowStock(ctx, typeID); err != nil {
		log.Printf("Couldn't evaluate low stock for item type %v. Error: %v", typeID, err)
	}
}

// EvaluateLowStockForItemAsync re-counts the stock of whatever type an item belongs to. The property
// mutations work in item ids, and a property edit can rename an item out of one group and into
// another without either the type or the number of items changing - which moves two groups' counts at
// once, both of them inside this one type.
func EvaluateLowStockForItemAsync(ctx context.Context, itemID int64) {
	item, err := db.Queries.GetItemByID(ctx, itemID)
	if err != nil {
		log.Printf("Couldn't read item %v to evaluate low stock. Error: %v", itemID, err)
		return
	}

	EvaluateLowStockAsync(ctx, item.TypeID)
}

// isLowStock is the whole definition of out of stock: at or under the threshold, not merely under it,
// so a threshold of 1 warns on the last item rather than only once it is gone.
func isLowStock(inStockCount int64, threshold int32, hasThreshold bool) bool {
	return hasThreshold && inStockCount <= int64(threshold)
}

// stockGroupLabel names a group for display. derived_name_format is nullable, and
// render_item_derived_name yields an empty string when it is unset - every item of the type then
// falls into one nameless group, which is the correct grouping for a type that gives its items no
// distinguishing name, but needs the type's own name to be readable.
func stockGroupLabel(groupName, typeName string) string {
	if groupName == "" {
		return typeName
	}

	return groupName
}
