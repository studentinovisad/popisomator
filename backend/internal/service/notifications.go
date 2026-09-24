package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/studentinovisad/popisomator/backend/internal/db"
	"github.com/studentinovisad/popisomator/backend/internal/dto"
	"github.com/studentinovisad/popisomator/backend/internal/repository"
)

func getRecipientsByRoles(ctx context.Context, roles ...repository.UserRole) ([]int64, error) {
	users, err := db.Queries.GetActiveUsersByRoles(ctx, roles)
	if err != nil {
		return nil, err
	}

	recipientIDs := make([]int64, len(users))
	for index, user := range users {
		recipientIDs[index] = user.ID
	}

	return recipientIDs, nil
}

func DeleteItemRequestNotifications(ctx context.Context, queries repository.Querier, itemID int64) error {
	_, err := queries.DeleteItemRequestNotifications(ctx, itemID)
	return err
}

func CreateItemRequestNotifications(ctx context.Context, queries repository.Querier, userID, itemID int64, approvedNotification bool) ([]int64, error) {
	var recipientIDs []int64
	if approvedNotification {
		recipientIDs = []int64{userID}

		if err := DeleteItemRequestNotifications(ctx, queries, itemID); err != nil {
			return nil, err
		}
	} else {
		var err error
		recipientIDs, err = getRecipientsByRoles(ctx, "admin", "manager")
		if err != nil {
			return nil, err
		}
	}

	notifications, err := queries.CreateNotifications(ctx, repository.CreateNotificationsParams{
		Kind:         repository.NotificationKindItemRequest,
		RecipientIds: recipientIDs,
	})
	if err != nil {
		return nil, err
	}

	notificationIDs := notificationIDsOf(notifications)
	if _, err := queries.CreateNotificationDescriptors_ItemRequest(ctx, repository.CreateNotificationDescriptors_ItemRequestParams{
		NotificationIds: notificationIDs,
		UserID:          userID,
		ItemID:          itemID,
	}); err != nil {
		return nil, err
	}

	return notificationIDs, nil
}

func CreateItemExpiryNotifications(ctx context.Context, itemID int64, expiryType repository.NotifdescExpiryType) ([]int64, error) {
	recipientIDs, err := getRecipientsByRoles(ctx, "admin", "manager")
	if err != nil {
		return nil, err
	}

	tx, err := db.BeginTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	queriesTx := db.Queries.WithTx(tx)

	// Filter out recipient IDs so that users don't receive the same notifications twice
	recipients := make(map[int64]struct{}, 0)
	for _, recipientID := range recipientIDs {
		recipients[recipientID] = struct{}{}
	}

	existingNotifications, err := queriesTx.GetExistingExpiryNotifications(ctx, repository.GetExistingExpiryNotificationsParams{
		RecipientIds: recipientIDs,
		ItemID:       itemID,
		ExpiryType:   expiryType,
	})
	if err != nil {
		return nil, err
	}
	for _, recipientID := range existingNotifications {
		delete(recipients, recipientID)
	}
	filteredRecipientIDs := make([]int64, 0)
	for recipientID := range recipients {
		filteredRecipientIDs = append(filteredRecipientIDs, recipientID)
	}

	// Create notifications only for filtered recipient IDs
	notifications, err := queriesTx.CreateNotifications(ctx, repository.CreateNotificationsParams{
		Kind:         repository.NotificationKindItemExpiry,
		RecipientIds: filteredRecipientIDs,
	})
	if err != nil {
		return nil, err
	}

	notificationIDs := notificationIDsOf(notifications)
	if _, err := queriesTx.CreateNotificationDescriptors_ItemExpiry(ctx, repository.CreateNotificationDescriptors_ItemExpiryParams{
		NotificationIds: notificationIDs,
		ItemID:          itemID,
		ExpiryType:      expiryType,
	}); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return notificationIDs, nil
}

func CreateLowStockNotifications(
	ctx context.Context,
	queries repository.Querier,
	typeID int64,
	groupName string,
	threshold, observed int32,
) ([]int64, error) {
	recipientIDs, err := getRecipientsByRoles(ctx, "admin", "manager")
	if err != nil {
		return nil, err
	}

	notifications, err := queries.CreateNotifications(ctx, repository.CreateNotificationsParams{
		Kind:         repository.NotificationKindItemLowStock,
		RecipientIds: recipientIDs,
	})
	if err != nil {
		return nil, err
	}

	notificationIDs := notificationIDsOf(notifications)
	if _, err := queries.CreateNotificationDescriptors_LowStock(ctx, repository.CreateNotificationDescriptors_LowStockParams{
		NotificationIds: notificationIDs,
		TypeID:          pgtype.Int8{Int64: typeID, Valid: true},
		GroupName:       groupName,
		Threshold:       threshold,
		Observed:        observed,
	}); err != nil {
		return nil, err
	}

	return notificationIDs, nil
}

func notificationIDsOf(notifications []repository.Notification) []int64 {
	ids := make([]int64, len(notifications))
	for index, notification := range notifications {
		ids[index] = notification.ID
	}

	return ids
}

func ListNotifications(ctx context.Context, recipient_id int64, limit, offset int32) (dto.NotificationsPage, error) {
	total, err := db.Queries.CountNotifications(ctx, recipient_id)
	if err != nil {
		return dto.NotificationsPage{}, err
	}

	totalUnread, err := db.Queries.CountUnreadNotifications(ctx, recipient_id)
	if err != nil {
		return dto.NotificationsPage{}, err
	}

	notifications, err := db.Queries.ListNotifications(ctx, repository.ListNotificationsParams{
		RecipientID: recipient_id,
		PageLimit:   limit,
		PageOffset:  offset,
	})
	if err != nil {
		return dto.NotificationsPage{}, err
	}

	pageItems := make([]dto.Notification, len(notifications))
	for index, row := range notifications {
		notif := dto.ToNotificationDTO(row.Notification)
		switch notif.Kind {
		case repository.NotificationKindItemRequest:
			itemRequest, err := GetItemRequest(ctx, row.ItemRequestUserID.Int64, row.ItemRequestItemID.Int64)
			if err != nil {
				return dto.NotificationsPage{}, err
			}
			notif.Descriptor_ItemRequest = &itemRequest
		case repository.NotificationKindItemExpiry:
			item, err := GetItem(ctx, row.ItemExpiryItemID.Int64, recipient_id)
			if err != nil {
				return dto.NotificationsPage{}, err
			}
			notif.Descriptor_ItemExpiry = &dto.NotificationDescriptor_ItemExpiry{
				Item: item,
				Type: row.ItemExpiryType.NotifdescExpiryType,
			}
		case repository.NotificationKindItemLowStock:
			descriptor := dto.NotificationDescriptor_LowStock{
				TypeName:  row.LowStockTypeName.String,
				GroupName: row.LowStockGroupName.String,
				Threshold: row.LowStockThreshold.Int32,
				Observed:  row.LowStockObserved.Int32,
			}
			if row.LowStockTypeID.Valid {
				descriptor.TypeID = &row.LowStockTypeID.Int64
			}
			notif.Descriptor_LowStock = &descriptor
		}
		pageItems[index] = notif
	}

	return dto.NotificationsPage{Items: pageItems, Limit: limit, Offset: offset, Total: total, TotalUnread: totalUnread}, nil
}

func CountUnreadNotifications(ctx context.Context, recipient_id int64) (int64, error) {
	totalUnread, err := db.Queries.CountUnreadNotifications(ctx, recipient_id)
	if err != nil {
		return 0, err
	}

	return totalUnread, nil
}

func ReadNotifications(ctx context.Context, recipient_id int64) (int64, error) {
	rowsAffected, err := db.Queries.ReadNotifications(ctx, recipient_id)
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func DeleteNotification(ctx context.Context, id int64, recipient_id int64) error {
	rowsAffected, err := db.Queries.DeleteNotification(ctx, repository.DeleteNotificationParams{
		ID:          id,
		RecipientID: recipient_id,
	})
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
