import { jsonRequest, request } from '$lib/api/client';
import type {
	ConsumptionStatus,
	CreateItemRequest,
	Item,
	ItemProperty,
	ItemsPage,
	ListItemsParams,
	UpdateItemRequest
} from '$lib/api/types';

export const itemsApi = {
	listItems: ({
		limit = 20,
		offset = 0,
		search = '',
		typeID,
		propertyFilters = {}
	}: ListItemsParams = {}) => {
		const query = new URLSearchParams({ limit: String(limit), offset: String(offset) });
		if (search) query.set('search', search);
		if (typeID) query.set('type_id', String(typeID));
		for (const [propertyID, value] of Object.entries(propertyFilters)) {
			if (value) query.set(`property.${propertyID}`, JSON.stringify(value));
		}
		return request<ItemsPage>(`/items?${query}`);
	},
	getItem: (id: number) => request<Item>(`/items/${id}`),
	createItem: (payload: CreateItemRequest) =>
		request<Item[]>('/items', jsonRequest('POST', payload)),
	updateItem: (id: number, payload: UpdateItemRequest) =>
		request<Item>(`/items/${id}`, jsonRequest('PATCH', payload)),
	consumeItem: (id: number, status: ConsumptionStatus) =>
		request<Item>(`/items/${id}/consume`, jsonRequest('POST', { consumption: status })),
	deleteItem: (id: number) => request<void>(`/items/${id}`, { method: 'DELETE' }),
	addItemProperty: (itemID: number, propertyID: number, value: string) =>
		request<ItemProperty>(
			`/items/${itemID}/properties`,
			jsonRequest('POST', { property_id: propertyID, value })
		),
	updateItemProperty: (itemID: number, propertyID: number, value: string) =>
		request<ItemProperty>(
			`/items/${itemID}/properties/${propertyID}`,
			jsonRequest('PUT', { value })
		),
	removeItemProperty: (itemID: number, propertyID: number) =>
		request<void>(`/items/${itemID}/properties/${propertyID}`, { method: 'DELETE' })
};
