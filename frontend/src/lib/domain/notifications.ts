import { resolve } from '$app/paths';
import type { Notification, UserRole } from '$lib/api';

export type NotificationIconName = 'request' | 'expiring' | 'expired' | 'generic';

function itemRequestLabel(itemID: number, itemName: string | undefined) {
	return itemName ?? `Stavka #${itemID}`;
}

export function notificationTitle(notification: Notification) {
	const itemRequest = notification.desc_item_request;
	if (itemRequest) {
		const itemLabel = itemRequestLabel(itemRequest.item_id, itemRequest.item_name);
		if (itemRequest.status === 'approved') {
			return `Zahtev za stavku „${itemLabel}” je odobren`;
		}

		const userLabel = itemRequest.user_name ?? `Korisnik #${itemRequest.user_id}`;
		return `„${userLabel}” je zatražio stavku „${itemLabel}”`;
	}

	const itemExpiry = notification.desc_item_expiry;
	if (itemExpiry) {
		const itemLabel = itemRequestLabel(itemExpiry.item.id, itemExpiry.item.derived_name);
		return itemExpiry.expiry_type === 'expired'
			? `Stavci „${itemLabel}” je istekao rok`
			: `Stavci „${itemLabel}” uskoro ističe rok`;
	}

	return 'Obaveštenje';
}

export function notificationDetail(notification: Notification) {
	return notification.desc_item_request?.reason ?? '';
}

export function notificationIcon(notification: Notification): NotificationIconName {
	if (notification.desc_item_request) return 'request';
	if (notification.desc_item_expiry) {
		return notification.desc_item_expiry.expiry_type === 'expired' ? 'expired' : 'expiring';
	}

	return 'generic';
}

// Where clicking the notification takes the user, as a ready-to-use href. Item request
// notifications land on whichever request list the recipient is allowed to see; for the managers
// and admins who can act on it, that list is pre-filtered to the person who made the request, so
// the row they were told about is the one in front of them. `user_id` is the filter
// ItemRequestsList already reads off the URL.
export function notificationHref(
	notification: Notification,
	role: UserRole | undefined
): string | null {
	const itemRequest = notification.desc_item_request;
	if (itemRequest) {
		return role === 'admin' || role === 'manager'
			? `${resolve('/item-requests')}?user_id=${itemRequest.user_id}`
			: resolve('/item-requests/me');
	}

	const itemExpiry = notification.desc_item_expiry;
	if (itemExpiry) return resolve(`/items/${itemExpiry.item.id}`);

	return null;
}
