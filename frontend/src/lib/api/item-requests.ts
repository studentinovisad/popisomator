import { jsonRequest, request } from '$lib/api/client';
import type {
	CreatePersonalItemRequest,
	ItemRequest,
	ItemRequestPreparationReport,
	ItemRequestPreparationReportParams,
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
	listItemRequests: ({
		limit = 20,
		offset = 0,
		status,
		userIDs,
		createdFrom,
		createdTo
	}: ListItemRequestsParams = {}) => {
		const query = new URLSearchParams({ limit: String(limit), offset: String(offset) });
		if (status) query.set('status', status);
		for (const userID of userIDs ?? []) query.append('user_id', String(userID));
		if (createdFrom) query.set('created_from', createdFrom);
		if (createdTo) query.set('created_to', createdTo);
		return request<ItemRequestsPage>(`/item-requests?${query}`);
	},
	listItemRequestUsers: () => request<ItemRequestUserOption[]>('/item-requests/users'),
	getItemRequestPreparationReport: (
		userID: number,
		{ itemIDs, status, createdFrom, createdTo }: ItemRequestPreparationReportParams = {}
	) => {
		const query = new URLSearchParams({ user_id: String(userID) });
		for (const itemID of itemIDs ?? []) query.append('item_id', String(itemID));
		if (status) query.set('status', status);
		if (createdFrom) query.set('created_from', createdFrom);
		if (createdTo) query.set('created_to', createdTo);
		return request<ItemRequestPreparationReport>(`/item-requests/preparation-report?${query}`);
	},
	approveItemRequest: (userID: number, itemID: number) =>
		request<ItemRequest>(
			'/item-requests/approve',
			jsonRequest('POST', { user_id: userID, item_id: itemID })
		),
	denyItemRequest: (userID: number, itemID: number) =>
		request<void>('/item-requests', jsonRequest('DELETE', { user_id: userID, item_id: itemID }))
};
