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
	PropertyOption,
	UpdateItemTypeRequest,
	UpdatePropertyRequest
} from '$lib/api/types';

export const catalogApi = {
	listItemTypes: ({ limit = 20, offset = 0 }: PageRequest = {}) => {
		const query = new URLSearchParams({ limit: String(limit), offset: String(offset) });
		return request<ItemTypesPage>(`/item-types?${query}`);
	},
	getItemTypeOptions: () => request<ItemTypeOption[]>('/item-types/options'),
	getItemType: (id: number) => request<ItemType>(`/item-types/${id}`),
	createItemType: (payload: CreateItemTypeRequest) =>
		request<ItemType>('/item-types', jsonRequest('POST', payload)),
	updateItemType: (id: number, payload: UpdateItemTypeRequest) =>
		request<ItemType>(`/item-types/${id}`, jsonRequest('PATCH', payload)),
	deleteItemType: (id: number) => request<void>(`/item-types/${id}`, { method: 'DELETE' }),
	addItemTypeProperty: (itemTypeID: number, payload: AddUpdateItemTypePropertyRequest) =>
		request<ItemTypeProperty>(
			`/item-types/${itemTypeID}/properties`,
			jsonRequest('POST', payload)
		),
	updateItemTypeProperty: (itemTypeID: number, propertyID: number, payload: AddUpdateItemTypePropertyRequest) =>
		request<ItemTypeProperty>(
			`/item-types/${itemTypeID}/properties/${propertyID}`,
			jsonRequest('PATCH', payload)
		),
	removeItemTypeProperty: (itemTypeID: number, propertyID: number) =>
		request<void>(`/item-types/${itemTypeID}/properties/${propertyID}`, { method: 'DELETE' }),
	listProperties: ({ limit = 20, offset = 0 }: PageRequest = {}) => {
		const query = new URLSearchParams({ limit: String(limit), offset: String(offset) });
		return request<PropertiesPage>(`/properties?${query}`);
	},
	getPropertyOptions: () => request<PropertyOption[]>('/properties/options'),
	getProperty: (id: number) => request<Property>(`/properties/${id}`),
	createProperty: (payload: CreatePropertyRequest) =>
		request<Property>('/properties', jsonRequest('POST', payload)),
	updateProperty: (id: number, payload: UpdatePropertyRequest) =>
		request<Property>(`/properties/${id}`, jsonRequest('PATCH', payload)),
	deleteProperty: (id: number) => request<void>(`/properties/${id}`, { method: 'DELETE' })
};
