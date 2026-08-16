import { env } from '$env/dynamic/private';
import type { User } from '$lib/api';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ cookies, fetch, setHeaders }) => {
	setHeaders({ 'cache-control': 'private, no-store' });

	const session = cookies.get('session');
	let currentUser: User | null = null;

	if (session) {
		try {
			const response = await fetch(
				`${env.POPISOMATOR_BACKEND_URL ?? 'http://localhost:8080'}/user/details`,
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

	return {
		currentUser,
		sidebarExpanded: cookies.get('popisomator-sidebar') === 'expanded'
	};
};
