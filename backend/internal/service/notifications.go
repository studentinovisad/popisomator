package service

import (
	"context"

	"github.com/studentinovisad/popisomator/backend/internal/db"
	"github.com/studentinovisad/popisomator/backend/internal/dto"
	"github.com/studentinovisad/popisomator/backend/internal/repository"
)

// CreateItemRequestNotifications gives each recipient a notification about userID's request for
// itemID, returning the new notification IDs in recipient order. The notification and its
// descriptor are written under one transaction, since a notification whose descriptor is missing
// has nothing to render. Recipients are taken as given; picking them is the caller's decision.
func CreateItemRequestNotifications(ctx context.Context, recipientIDs []int64, userID, itemID int64) ([]int64, error) {
	tx, err := db.BeginTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	queriesTx := db.Queries.WithTx(tx)

	notifications, err := queriesTx.CreateNotifications(ctx, repository.CreateNotificationsParams{
		Kind:         repository.NotificationKindItemRequest,
		RecipientIds: recipientIDs,
	})
	if err != nil {
		return nil, err
	}

	notificationIDs := notificationIDsOf(notifications)
	if _, err := queriesTx.CreateNotificationDescriptors_ItemRequest(ctx, repository.CreateNotificationDescriptors_ItemRequestParams{
		NotificationIds: notificationIDs,
		UserID:          userID,
		ItemID:          itemID,
	}); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return notificationIDs, nil
}

// CreateItemExpiryNotifications gives each recipient a notification that itemID is near or past its
// expiry date, expiryType saying which, and returns the new notification IDs in recipient order.
// Written under one transaction for the same reason as CreateItemRequestNotifications.
func CreateItemExpiryNotifications(ctx context.Context, recipientIDs []int64, itemID int64, expiryType repository.NotifdescExpiryType) ([]int64, error) {
	tx, err := db.BeginTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	queriesTx := db.Queries.WithTx(tx)

	notifications, err := queriesTx.CreateNotifications(ctx, repository.CreateNotificationsParams{
		Kind:         repository.NotificationKindItemExpiry,
		RecipientIds: recipientIDs,
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
			itemRequest, err := GetItemRequest(ctx, dto.ItemRequestIdentifierRequest{
				UserID: row.ItemRequestUserID.Int64,
				ItemID: row.ItemRequestItemID.Int64,
			})
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
