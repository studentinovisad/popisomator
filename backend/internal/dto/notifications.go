package dto

import (
	"time"

	"github.com/studentinovisad/popisomator/backend/internal/repository"
)

type Notification struct {
	ID                     int64                              `json:"id"`
	RecipientID            int64                              `json:"recipient_id"`
	CreatedAt              time.Time                          `json:"created_at"`
	Kind                   repository.NotificationKind        `json:"kind"`
	Read                   bool                               `json:"read"`
	Descriptor_ItemRequest *ItemRequest                       `json:"desc_item_request,omitempty"`
	Descriptor_ItemExpiry  *NotificationDescriptor_ItemExpiry `json:"desc_item_expiry,omitempty"`
}

type NotificationDescriptor_ItemExpiry struct {
	Item Item                           `json:"item"`
	Type repository.NotifdescExpiryType `json:"expiry_type"`
}

type NotificationsPage struct {
	Items       []Notification `json:"items"`
	Limit       int32          `json:"limit"`
	Offset      int32          `json:"offset"`
	Total       int64          `json:"total"`
	TotalUnread int64          `json:"total_unread"`
}

func ToNotificationDTO(notif repository.Notification) Notification {
	return Notification{
		ID:          notif.ID,
		RecipientID: notif.RecipientID,
		CreatedAt:   notif.CreatedAt.Time,
		Kind:        notif.Kind,
		Read:        notif.Read,
	}
}
