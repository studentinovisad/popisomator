import type { UserRole } from '$lib/api';

export type NavigationIconName =
	'inventory' | 'catalog' | 'settings' | 'users' | 'requests' | 'notifications';

export type AppPath =
	| '/'
	| '/items/new'
	| '/item-requests'
	| '/item-requests/me'
	| '/account'
	| '/admin/users'
	| '/admin/users/pending'
	| '/catalog/item-types'
	| '/catalog/item-types/new'
	| '/catalog/properties'
	| '/catalog/properties/new'
	| '/notifications'
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

const itemTypeEditPageMetadata: PageMetadata = {
	title: 'Izmeni tip stavke',
	description: 'Izmenite svojstva, podrazumevane vrednosti i način prikaza stavki ovog tipa.'
};

const propertyEditPageMetadata: PageMetadata = {
	title: 'Izmeni svojstvo',
	description: 'Izmenite naziv, opis i podrazumevanu vrednost svojstva.'
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
	'/item-requests/me': {
		title: 'Moji zahtevi',
		description: 'Pratite status svojih zahteva za korišćenje stavki.'
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
	'/item-requests': {
		title: 'Zahtevi',
		description: 'Odobrite ili odbijte zahteve korisnika za korišćenje stavki.'
	},
	'/catalog/item-types': {
		title: 'Tipovi stavki',
		description: 'Upravljajte tipovima stavki i njihovim pripadajućim svojstvima.'
	},
	'/catalog/item-types/new': {
		title: 'Novi tip stavke',
		description: 'Odaberite svojstva koja pripadaju ovom tipu.'
	},
	'/catalog/properties': {
		title: 'Svojstva',
		description: 'Upravljajte svojstvima koja se mogu dodeliti stavkama.'
	},
	'/catalog/properties/new': {
		title: 'Novo svojstvo',
		description: 'Odaberite tip vrednosti i opcionalnu podrazumevanu vrednost.'
	},
	'/notifications': {
		title: 'Obaveštenja',
		description: 'Pregledajte obaveštenja o zahtevima i rokovima stavki.'
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
		path: '/item-requests',
		label: 'Zahtevi',
		icon: 'requests',
		requiredRoles: ['manager', 'admin']
	},
	{ path: '/item-requests/me', label: 'Moji zahtevi', icon: 'requests', requiredRoles: ['user'] },
	{ path: '/admin/users', label: 'Korisnici', icon: 'users', requiredRoles: ['admin'] },
	{
		path: '/catalog/item-types',
		label: 'Katalog',
		icon: 'catalog',
		requiredRoles: ['admin']
	}
];

export const secondaryNavigation: NavigationItem[] = [
	// Every signed-in role has notifications; listing them all is what hides the link from signed-out
	// visitors, who can still reach Podešavanja below.
	{
		path: '/notifications',
		label: 'Obaveštenja',
		icon: 'notifications',
		requiredRoles: ['admin', 'manager', 'user']
	},
	{ path: '/settings', label: 'Podešavanja', icon: 'settings' }
];

export function getPageMetadata(pathname: string): PageMetadata {
	if (pathname.startsWith('/items/') && pathname !== '/items/new') {
		return {
			title: 'Stavka',
			description: 'Pregledajte podatke, svojstva i stanje odabrane stavke.'
		};
	}

	if (pathname === '/catalog/item-types/new') {
		return pageMetadata['/catalog/item-types/new'];
	}

	if (pathname.startsWith('/catalog/item-types/')) {
		return itemTypeEditPageMetadata;
	}

	if (pathname === '/catalog/properties/new') {
		return pageMetadata['/catalog/properties/new'];
	}

	if (pathname.startsWith('/catalog/properties/')) {
		return propertyEditPageMetadata;
	}

	return pageMetadata[pathname as AppPath] ?? fallbackPageMetadata;
}
