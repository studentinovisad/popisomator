import { jsonRequest, request } from '$lib/api/client';
import type {
	CreatePersonalItemRequest,
	ItemRequest,
	ItemRequestsPage,
	ItemRequestUserOption,
	ListItemRequestsParams
} from '$lib/api/types';

export const itemRequestsApi = {
	createPersonalItemRequest: (payload: CreatePersonalItemRequest) =>
		request<ItemRequest>('/item-requests/me', jsonRequest('POST', payload)),
	listPersonalItemRequests: ({ limit = 20, offset = 0 }: ListItemRequestsParams = {}) => {
		const query = new URLSearchParams({ limit: String(limit), offset: String(offset) });
		return request<ItemRequestsPage>(`/item-requests/me?${query}`);
	},
	listItemRequests: ({ limit = 20, offset = 0, status, userID }: ListItemRequestsParams = {}) => {
		const query = new URLSearchParams({ limit: String(limit), offset: String(offset) });
		if (status) query.set('status', status);
		if (userID) query.set('user_id', String(userID));
		return request<ItemRequestsPage>(`/item-requests?${query}`);
	},
	listItemRequestUsers: () => request<ItemRequestUserOption[]>('/item-requests/users'),
	approveItemRequest: (userID: number, itemID: number) =>
		request<ItemRequest>(
			'/item-requests/approve',
			jsonRequest('POST', { user_id: userID, item_id: itemID })
		),
	denyItemRequest: (userID: number, itemID: number) =>
		request<void>('/item-requests', jsonRequest('DELETE', { user_id: userID, item_id: itemID }))
};
