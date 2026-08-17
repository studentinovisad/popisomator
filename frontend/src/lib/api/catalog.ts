import { jsonRequest, request } from '$lib/api/client';
import type {
	AddUpdateItemTypePropertyRequest,
	CreateItemTypeRequest,
	CreatePropertyRequest,
	ItemType,
	ItemTypeOption,
	ItemTypeProperty,
	ItemTypesPage,
	PageRequest,
	PropertiesPage,
	Property,
	UpdateItemTypeRequest,
	UpdatePropertyRequest
} from '$lib/api/types';

export const catalogApi = {
	listItemTypes: () => request<ItemTypeOption[]>('/item/types'),
	listItemTypesPage: ({ limit = 20, offset = 0 }: PageRequest = {}) => {
		const query = new URLSearchParams({ limit: String(limit), offset: String(offset) });
		return request<ItemTypesPage>(`/item/types/page?${query}`);
	},
	getItemType: (id: number) => request<ItemType>(`/item/types/${id}`),
	createItemType: (payload: CreateItemTypeRequest) =>
		request<ItemType>('/item/types', jsonRequest('POST', payload)),
	updateItemType: (id: number, payload: UpdateItemTypeRequest) =>
		request<ItemType>(`/item/types/${id}`, jsonRequest('PATCH', payload)),
	deleteItemType: (id: number) => request<void>(`/item/types/${id}`, { method: 'DELETE' }),
	addItemTypeProperty: (itemTypeID: number, payload: AddUpdateItemTypePropertyRequest) =>
		request<ItemTypeProperty>(`/item/types/${itemTypeID}/properties`, jsonRequest('POST', payload)),
	updateItemTypeProperty: (
		itemTypeID: number,
		propertyID: number,
		payload: AddUpdateItemTypePropertyRequest
	) =>
		request<ItemTypeProperty>(
			`/item/types/${itemTypeID}/properties/${propertyID}`,
			jsonRequest('PATCH', payload)
		),
	removeItemTypeProperty: (itemTypeID: number, propertyID: number) =>
		request<void>(`/item/types/${itemTypeID}/properties/${propertyID}`, { method: 'DELETE' }),
	listProperties: () => request<Property[]>('/item/properties'),
	listPropertiesPage: ({ limit = 20, offset = 0 }: PageRequest = {}) => {
		const query = new URLSearchParams({ limit: String(limit), offset: String(offset) });
		return request<PropertiesPage>(`/item/properties/page?${query}`);
	},
	getProperty: (id: number) => request<Property>(`/item/properties/${id}`),
	createProperty: (payload: CreatePropertyRequest) =>
		request<Property>('/item/properties', jsonRequest('POST', payload)),
	updateProperty: (id: number, payload: UpdatePropertyRequest) =>
		request<Property>(`/item/properties/${id}`, jsonRequest('PATCH', payload)),
	deleteProperty: (id: number) => request<void>(`/item/properties/${id}`, { method: 'DELETE' })
};
