export type UserRole = 'admin' | 'manager' | 'user';
export type UserStatus = 'requested' | 'active';

export type ConsumptionStatus =
	'not_consumed' | 'partially_consumed' | 'fully_consumed' | 'damaged';

export type ItemProperty = {
	id: number;
	value: string;
};

export type Item = {
	id: number;
	consumption: ConsumptionStatus;
	properties: ItemProperty[];
	type_id: number;
};

export type ItemsPage = {
	items: Item[];
	limit: number;
	offset: number;
	total: number;
};

export type ListItemsParams = {
	limit?: number;
	offset?: number;
};

export type ItemTypeProperty = {
	id: number;
	default_value: string | null;
};

export type ItemType = {
	id: number;
	name: string;
	description: string;
	properties: ItemTypeProperty[];
};

export type PropertyValueType = 'string' | 'number' | 'boolean';

export type Property = {
	id: number;
	name: string;
	description: string;
	value_type: PropertyValueType;
	default_value: string | null;
};

export type User = {
	id: number;
	email: string;
	full_name: string;
	role: UserRole;
	status: UserStatus;
};

export type UsersPage = {
	items: User[];
	limit: number;
	offset: number;
	total: number;
};

export type ListUsersParams = {
	limit?: number;
	offset?: number;
	search?: string;
	role?: UserRole;
	status?: UserStatus;
};

type LoginRequest = {
	email: string;
	password: string;
};

type RegistrationRequest = LoginRequest & {
	full_name: string;
};

type CreateUserRequest = RegistrationRequest & {
	role: UserRole;
};

export type CreateItemRequest = {
	type_id: number;
	properties: ItemProperty[];
	amount: number;
};

export type CreateItemTypeRequest = {
	name: string;
	description: string;
	properties: ItemTypeProperty[];
};

export type UpdateItemTypeRequest = Partial<Pick<ItemType, 'name' | 'description'>>;

export type CreatePropertyRequest = {
	name: string;
	description: string;
	value_type: PropertyValueType;
	default_value: string | null;
};

export type UpdatePropertyRequest = Partial<
	Pick<CreatePropertyRequest, 'name' | 'description' | 'default_value'>
>;

export class ApiError extends Error {
	constructor(
		message: string,
		readonly status: number
	) {
		super(message);
	}
}

async function request<T>(path: string, init: RequestInit = {}): Promise<T> {
	const response = await fetch(`/api${path}`, {
		credentials: 'include',
		...init
	});

	const contentType = response.headers.get('content-type');
	const body = contentType?.includes('application/json') ? await response.json() : undefined;

	if (!response.ok) {
		const message =
			typeof body === 'object' && body !== null && 'error' in body && typeof body.error === 'string'
				? body.error
				: 'Request failed';
		throw new ApiError(message, response.status);
	}

	return body as T;
}

function jsonRequest(method: 'PATCH' | 'POST' | 'PUT', body: unknown): RequestInit {
	return {
		method,
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(body)
	};
}

export const api = {
	login: (payload: LoginRequest) => request<void>('/auth/login', jsonRequest('POST', payload)),
	logout: () => request<void>('/auth/logout', { method: 'POST' }),
	currentUser: () => request<User>('/user/details'),
	listUsers: ({ limit = 25, offset = 0, search = '', role, status }: ListUsersParams = {}) => {
		const query = new URLSearchParams({ limit: String(limit), offset: String(offset) });
		if (search) {
			query.set('search', search);
		}
		if (role !== undefined) {
			query.set('role', role);
		}
		if (status !== undefined) {
			query.set('status', status);
		}

		return request<UsersPage>(`/users?${query}`);
	},
	createUser: (payload: CreateUserRequest) =>
		request<User>('/auth/create', jsonRequest('POST', payload)),
	register: (payload: RegistrationRequest) =>
		request<User>('/auth/register', jsonRequest('POST', payload)),
	updateUserRole: (id: number, role: UserRole) =>
		request<User>(`/user/${id}/role`, jsonRequest('PATCH', { role })),
	approveRegistration: (id: number) => request<User>(`/user/${id}/approve`, { method: 'POST' }),
	declineRegistration: (id: number) => request<void>(`/user/${id}/decline`, { method: 'POST' }),
	listItems: ({ limit = 20, offset = 0 }: ListItemsParams = {}) => {
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
		request<void>(`/item/${itemID}/properties/${propertyID}`, { method: 'DELETE' }),
	listItemTypes: () => request<ItemType[]>('/item/types'),
	getItemType: (id: number) => request<ItemType>(`/item/types/${id}`),
	createItemType: (payload: CreateItemTypeRequest) =>
		request<ItemType>('/item/types', jsonRequest('POST', payload)),
	updateItemType: (id: number, payload: UpdateItemTypeRequest) =>
		request<ItemType>(`/item/types/${id}`, jsonRequest('PATCH', payload)),
	deleteItemType: (id: number) => request<void>(`/item/types/${id}`, { method: 'DELETE' }),
	addItemTypeProperty: (itemTypeID: number, propertyID: number, defaultValue: string | null) =>
		request<ItemTypeProperty>(
			`/item/types/${itemTypeID}/properties`,
			jsonRequest('POST', { property_id: propertyID, default_value: defaultValue })
		),
	updateItemTypeProperty: (itemTypeID: number, propertyID: number, defaultValue: string | null) =>
		request<ItemTypeProperty>(
			`/item/types/${itemTypeID}/properties/${propertyID}`,
			jsonRequest('PUT', { default_value: defaultValue })
		),
	removeItemTypeProperty: (itemTypeID: number, propertyID: number) =>
		request<void>(`/item/types/${itemTypeID}/properties/${propertyID}`, { method: 'DELETE' }),
	listProperties: () => request<Property[]>('/item/properties'),
	getProperty: (id: number) => request<Property>(`/item/properties/${id}`),
	createProperty: (payload: CreatePropertyRequest) =>
		request<Property>('/item/properties', jsonRequest('POST', payload)),
	updateProperty: (id: number, payload: UpdatePropertyRequest) =>
		request<Property>(`/item/properties/${id}`, jsonRequest('PATCH', payload)),
	deleteProperty: (id: number) => request<void>(`/item/properties/${id}`, { method: 'DELETE' })
};
