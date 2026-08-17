import type { UserRole } from '$lib/api';

export type NavigationIconName = 'inventory' | 'catalog' | 'settings' | 'users';

export type AppPath =
	| '/'
	| '/items/new'
	| '/account'
	| '/admin/users'
	| '/admin/users/pending'
	| '/catalog/item-types'
	| '/catalog/properties'
	| '/settings'
	| '/login'
	| '/register';

export type NavigationItem = {
	path: AppPath;
	label: string;
	icon: NavigationIconName;
	requiredRoles?: UserRole[];
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
		title: 'Stavke',
		description: 'Pratite stanje stavki i evidentirajte njihovu potrošnju.'
	},
	'/items/new': {
		title: 'Nova stavka',
		description: 'Dodajte stavku i njene početne vrednosti svojstava.'
	},
	'/account': {
		title: 'Moj nalog',
		description: 'Pregledajte podatke svog naloga.'
	},
	'/admin/users': {
		title: 'Korisnici',
		description: 'Upravljajte pristupom i ulogama korisnika sistema.'
	},
	'/admin/users/pending': {
		title: 'Zahtevi za registraciju',
		description: 'Odobrite ili odbijte zahteve za pristup sistemu.'
	},
	'/catalog/item-types': {
		title: 'Tipovi stavki',
		description: 'Upravljajte tipovima stavki i njihovim pripadajućim svojstvima.'
	},
	'/catalog/properties': {
		title: 'Svojstva',
		description: 'Upravljajte svojstvima koja se mogu dodeliti stavkama.'
	},
	'/settings': {
		title: 'Podešavanja',
		description: 'Prilagodite prikaz i proverite stanje sistema.'
	},
	'/login': {
		title: 'Prijava',
		description: 'Prijavite se nalogom koji je napravio ili odobrio administrator.'
	},
	'/register': {
		title: 'Registracija',
		description: 'Pošaljite zahtev za pravljenje naloga.'
	}
};

export const primaryNavigation: NavigationItem[] = [
	{ path: '/', label: 'Stavke', icon: 'inventory' },
	{
		path: '/catalog/item-types',
		label: 'Katalog',
		icon: 'catalog',
		requiredRoles: ['admin']
	},
	{ path: '/admin/users', label: 'Korisnici', icon: 'users', requiredRoles: ['admin'] }
];

export const secondaryNavigation: NavigationItem[] = [
	{ path: '/settings', label: 'Podešavanja', icon: 'settings' }
];

export function getPageMetadata(pathname: string): PageMetadata {
	if (pathname.startsWith('/catalog/item-types/')) {
		return pageMetadata['/catalog/item-types'];
	}

	if (pathname.startsWith('/catalog/properties/')) {
		return pageMetadata['/catalog/properties'];
	}

	return pageMetadata[pathname as AppPath] ?? fallbackPageMetadata;
}
