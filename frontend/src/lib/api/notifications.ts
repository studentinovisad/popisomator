import { request } from '$lib/api/client';
import type { ListNotificationsParams, NotificationsPage } from '$lib/api/types';

export const notificationsApi = {
	listNotifications: ({ limit = 20, offset = 0 }: ListNotificationsParams = {}) => {
		const query = new URLSearchParams({ limit: String(limit), offset: String(offset) });
		return request<NotificationsPage>(`/notifications?${query}`);
	},
	countUnreadNotifications: () => request<number>('/notifications/unread-count'),
	// Marks every notification of the current user as read, returning how many were still unread.
	readNotifications: () => request<number>('/notifications/read', { method: 'POST' }),
	deleteNotification: (id: number) => request<void>(`/notifications/${id}`, { method: 'DELETE' })
};
