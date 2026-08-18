import { env } from '$env/dynamic/private';
import { redirect } from '@sveltejs/kit';
import type { User } from '$lib/api';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ cookies, fetch, setHeaders, url }) => {
	setHeaders({ 'cache-control': 'private, no-store' });

	const session = cookies.get('session');
	let currentUser: User | null = null;

	if (session) {
		try {
			const response = await fetch(
				`${env.POPISOMATOR_BACKEND_URL ?? 'http://localhost:8080'}/users/me`,
				{
					headers: { Cookie: `session=${session}` }
				}
			);

			if (response.ok) {
				currentUser = (await response.json()) as User;
			}
		} catch {
			currentUser = null;
		}
	}

	if (currentUser && (url.pathname === '/login' || url.pathname === '/register')) {
		redirect(303, '/');
	}

	return {
		currentUser,
		sidebarExpanded: cookies.get('popisomator-sidebar') === 'expanded'
	};
};
