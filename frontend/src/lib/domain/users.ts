import type { UserRole } from '$lib/api';

export type UserRoleFilter = 'all' | UserRole;

export const userRoleOptions: { value: UserRole; label: string }[] = [
	{ value: 'admin', label: 'Administrator' },
	{ value: 'manager', label: 'Menadžer' },
	{ value: 'user', label: 'Korisnik' }
];

export const roleFilterOptions: { value: UserRoleFilter; label: string }[] = [
	{ value: 'all', label: 'Sve uloge' },
	...userRoleOptions
];

export function userRoleLabel(role: UserRole) {
	return userRoleOptions.find((option) => option.value === role)?.label ?? role;
}
