import { jsonRequest, request } from '$lib/api/client';
import type {
	CreateUserRequest,
	UpdateUserRequest,
	ListUsersParams,
	LoginRequest,
	RegistrationRequest,
	User,
	UsersPage
} from '$lib/api/types';

export const usersApi = {
	login: (payload: LoginRequest) => request<void>('/auth/login', jsonRequest('POST', payload)),
	logout: () => request<void>('/auth/logout', { method: 'POST' }),
	currentUser: () => request<User>('/users/me'),
	listUsers: ({ limit = 25, offset = 0, search = '', role, status }: ListUsersParams = {}) => {
		const query = new URLSearchParams({ limit: String(limit), offset: String(offset) });
		if (search) query.set('search', search);
		if (role !== undefined) query.set('role', role);
		if (status !== undefined) query.set('status', status);

		return request<UsersPage>(`/users?${query}`);
	},
	createUser: (payload: CreateUserRequest) => request<User>('/users', jsonRequest('POST', payload)),
	register: (payload: RegistrationRequest) =>
		request<User>('/auth/register', jsonRequest('POST', payload)),
	updateUser: (id: number, payload: UpdateUserRequest) =>
		request<User>(`/users/${id}`, jsonRequest('PATCH', payload)),
	deleteUser: (id: number) => request<void>(`/users/${id}`, { method: 'DELETE' })
};
