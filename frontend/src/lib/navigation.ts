export type NavigationIconName = 'overview' | 'settings' | 'users';

export type AppPath = '/' | '/account' | '/admin/users' | '/settings' | '/login';

export type NavigationItem = {
	path: AppPath;
	label: string;
	icon: NavigationIconName;
	adminOnly?: boolean;
};

type PageMetadata = {
	title: string;
	description: string;
};

const fallbackPageMetadata: PageMetadata = {
	title: 'Popisomator',
	description: 'Popisomator'
};

export const pageMetadata: Record<AppPath, PageMetadata> = {
	'/': {
		title: 'Pregled',
		description: 'Pregled trenutnog stanja sistema.'
	},
	'/account': {
		title: 'Moj nalog',
		description: 'Pregledajte podatke svog naloga.'
	},
	'/admin/users': {
		title: 'Korisnici',
		description: 'Upravljajte pristupom i ulogama korisnika sistema.'
	},
	'/settings': {
		title: 'Podešavanja',
		description: 'Prilagodite prikaz i proverite stanje sistema.'
	},
	'/login': {
		title: 'Prijava',
		description: 'Prijavite se nalogom koji je napravio administrator.'
	}
};

export const primaryNavigation: NavigationItem[] = [
	{ path: '/', label: 'Pregled', icon: 'overview' },
	{ path: '/admin/users', label: 'Korisnici', icon: 'users', adminOnly: true }
];

export const secondaryNavigation: NavigationItem[] = [
	{ path: '/settings', label: 'Podešavanja', icon: 'settings' }
];

export function getPageMetadata(pathname: string): PageMetadata {
	return pageMetadata[pathname as AppPath] ?? fallbackPageMetadata;
}
