export type UserRole = 'admin' | 'manager' | 'user';

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
};

type LoginRequest = {
	email: string;
	password: string;
};

type CreateUserRequest = LoginRequest & {
	full_name: string;
	role: UserRole;
};

export type CreateItemRequest = {
	type_id: number;
	properties: ItemProperty[];
};

export type CreateItemTypeRequest = {
	name: string;
	description: string;
	properties: ItemTypeProperty[];
};

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
	listUsers: ({ limit = 25, offset = 0, search = '', role }: ListUsersParams = {}) => {
		const query = new URLSearchParams({ limit: String(limit), offset: String(offset) });
		if (search) {
			query.set('search', search);
		}
		if (role !== undefined) {
			query.set('role', role);
		}

		return request<UsersPage>(`/users?${query}`);
	},
	createUser: (payload: CreateUserRequest) =>
		request<User>('/auth/register', jsonRequest('POST', payload)),
	updateUserRole: (id: number, role: UserRole) =>
		request<User>(`/user/${id}/role`, jsonRequest('PATCH', { role })),
	listItems: () => request<Item[]>('/item'),
	createItem: (payload: CreateItemRequest) => request<Item>('/item', jsonRequest('POST', payload)),
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
	createItemType: (payload: CreateItemTypeRequest) =>
		request<ItemType>('/item/types', jsonRequest('POST', payload)),
	deleteItemType: (id: number) => request<void>(`/item/types/${id}`, { method: 'DELETE' }),
	listProperties: () => request<Property[]>('/item/properties'),
	getProperty: (id: number) => request<Property>(`/item/properties/${id}`),
	createProperty: (payload: CreatePropertyRequest) =>
		request<Property>('/item/properties', jsonRequest('POST', payload)),
	updateProperty: (id: number, payload: UpdatePropertyRequest) =>
		request<Property>(`/item/properties/${id}`, jsonRequest('PATCH', payload)),
	deleteProperty: (id: number) => request<void>(`/item/properties/${id}`, { method: 'DELETE' })
};
