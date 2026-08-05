export type UserRole = 'admin' | 'manager' | 'user';

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

function jsonRequest(method: 'PATCH' | 'POST', body: unknown): RequestInit {
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
		request<User>(`/user/${id}/role`, jsonRequest('PATCH', { role }))
};
