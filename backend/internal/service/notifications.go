package service

import (
	"context"

	"github.com/studentinovisad/popisomator/backend/internal/db"
	"github.com/studentinovisad/popisomator/backend/internal/dto"
	"github.com/studentinovisad/popisomator/backend/internal/repository"
)

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
				UserID: row.NotifdescItemRequest.UserID,
				ItemID: row.NotifdescItemRequest.ItemID,
			})
			if err != nil {
				return dto.NotificationsPage{}, err
			}
			notif.Descriptor_ItemRequest = &itemRequest
		case repository.NotificationKindItemExpiry:
			item, err := GetItem(ctx, row.NotifdescItemExpiry.ItemID, recipient_id)
			if err != nil {
				return dto.NotificationsPage{}, err
			}
			notif.Descriptor_ItemExpiry = &dto.NotificationDescriptor_ItemExpiry{
				Item: item,
				Type: row.NotifdescItemExpiry.ExpiryType,
			}
		}
		pageItems[index] = notif
	}

	return dto.NotificationsPage{Items: pageItems, Limit: limit, Offset: offset, Total: total, TotalUnread: totalUnread}, nil
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
