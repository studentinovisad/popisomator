import { jsonRequest, request } from '$lib/api/client';
import type {
	CreateUserRequest,
	ListUsersParams,
	LoginRequest,
	RegistrationRequest,
	User,
	UserRole,
	UsersPage
} from '$lib/api/types';

export const usersApi = {
	login: (payload: LoginRequest) => request<void>('/auth/login', jsonRequest('POST', payload)),
	logout: () => request<void>('/auth/logout', { method: 'POST' }),
	currentUser: () => request<User>('/user/details'),
	listUsers: ({ limit = 25, offset = 0, search = '', role, status }: ListUsersParams = {}) => {
		const query = new URLSearchParams({ limit: String(limit), offset: String(offset) });
		if (search) query.set('search', search);
		if (role !== undefined) query.set('role', role);
		if (status !== undefined) query.set('status', status);

		return request<UsersPage>(`/users?${query}`);
	},
	createUser: (payload: CreateUserRequest) =>
		request<User>('/auth/create', jsonRequest('POST', payload)),
	register: (payload: RegistrationRequest) =>
		request<User>('/auth/register', jsonRequest('POST', payload)),
	updateUserRole: (id: number, role: UserRole) =>
		request<User>(`/user/${id}/role`, jsonRequest('PATCH', { role })),
	approveRegistration: (id: number) => request<User>(`/user/${id}/approve`, { method: 'POST' }),
	declineRegistration: (id: number) => request<void>(`/user/${id}/decline`, { method: 'POST' })
};
