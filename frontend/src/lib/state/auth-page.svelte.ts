import { goto } from '$app/navigation';
import { resolve } from '$app/paths';
import { ApiError, type UserRole } from '$lib/api';
import type { AppPath } from '$lib/domain/navigation';
import { session } from '$lib/state/session.svelte';

type AuthPageOptions = {
	unavailableMessage: string;
	requiredRole?: UserRole;
	requiredRoles?: UserRole[];
	unauthorizedRedirect?: AppPath;
};

export function createAuthPage({
	unavailableMessage,
	requiredRole,
	requiredRoles,
	unauthorizedRedirect = '/account'
}: AuthPageOptions) {
	const state = $state({
		loading: true,
		error: '',
		authorized: false,
		user: null as Awaited<ReturnType<typeof session.getCurrentUser>> | null
	});

	async function load() {
		state.loading = true;
		state.error = '';
		state.authorized = false;

		try {
			const user = await session.getCurrentUser();
			session.setUser(user);
			state.user = user;

			const allowedRoles = requiredRoles ?? (requiredRole ? [requiredRole] : undefined);
			if (allowedRoles && !allowedRoles.includes(user.role)) {
				await goto(resolve(unauthorizedRedirect));
				return;
			}

			state.authorized = true;
		} catch (reason) {
			state.user = null;

			if (reason instanceof ApiError && reason.status === 401) {
				await goto(resolve('/login'));
				return;
			}

			state.error = unavailableMessage;
		} finally {
			state.loading = false;
		}
	}

	return { state, load };
}
