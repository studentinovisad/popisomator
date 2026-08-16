import { jsonRequest, request } from '$lib/api/client';
import type {
	ConsumptionStatus,
	CreateItemRequest,
	Item,
	ItemProperty,
	ItemsPage,
	PageRequest
} from '$lib/api/types';

export const itemsApi = {
	listItems: ({ limit = 20, offset = 0 }: PageRequest = {}) => {
		const query = new URLSearchParams({ limit: String(limit), offset: String(offset) });
		return request<ItemsPage>(`/item?${query}`);
	},
	createItem: (payload: CreateItemRequest) => request<Item>('/item', jsonRequest('POST', payload)),
	setItemType: (id: number, typeID: number) =>
		request<Item>(`/item/${id}`, jsonRequest('PATCH', { type_id: typeID })),
	consumeItem: (id: number, status: ConsumptionStatus) =>
		request<void>(`/item/${id}/consume`, jsonRequest('POST', { status })),
	deleteItem: (id: number) => request<void>(`/item/${id}`, { method: 'DELETE' }),
	addItemProperty: (itemID: number, propertyID: number, value: string) =>
		request<ItemProperty>(
			`/item/${itemID}/properties`,
			jsonRequest('POST', { property_id: propertyID, value })
		),
	updateItemProperty: (itemID: number, propertyID: number, value: string) =>
		request<ItemProperty>(
			`/item/${itemID}/properties/${propertyID}`,
			jsonRequest('PUT', { value })
		),
	removeItemProperty: (itemID: number, propertyID: number) =>
		request<void>(`/item/${itemID}/properties/${propertyID}`, { method: 'DELETE' })
};
