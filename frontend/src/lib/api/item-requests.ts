import { jsonRequest, request } from '$lib/api/client';
import type {
	CreatePersonalItemRequest,
	ItemRequest,
	ItemRequestsPage,
	ListItemRequestsParams
} from '$lib/api/types';

export const itemRequestsApi = {
	createPersonalItemRequest: (payload: CreatePersonalItemRequest) =>
		request<ItemRequest>('/item-requests/me', jsonRequest('POST', payload)),
	listPersonalItemRequests: ({ limit = 20, offset = 0 }: ListItemRequestsParams = {}) => {
		const query = new URLSearchParams({ limit: String(limit), offset: String(offset) });
		return request<ItemRequestsPage>(`/item-requests/me?${query}`);
	},
	listItemRequests: ({ limit = 20, offset = 0, status }: ListItemRequestsParams = {}) => {
		const query = new URLSearchParams({ limit: String(limit), offset: String(offset) });
		if (status) query.set('status', status);
		return request<ItemRequestsPage>(`/item-requests?${query}`);
	},
	approveItemRequest: (userID: number, itemID: number) =>
		request<ItemRequest>(
			'/item-requests/approve',
			jsonRequest('POST', { user_id: userID, item_id: itemID })
		),
	denyItemRequest: (userID: number, itemID: number) =>
		request<void>('/item-requests', jsonRequest('DELETE', { user_id: userID, item_id: itemID }))
};
