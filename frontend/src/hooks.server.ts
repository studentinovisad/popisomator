import { redirect, type Handle } from '@sveltejs/kit';

const publicPaths = new Set(['/login', '/settings']);

export const handle: Handle = async ({ event, resolve }) => {
	if (!publicPaths.has(event.url.pathname) && !event.cookies.get('session')) {
		throw redirect(303, '/login');
	}

	return resolve(event);
};
