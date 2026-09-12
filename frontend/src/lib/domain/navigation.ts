import type { UserRole } from '$lib/api';
import { notifications } from '$lib/state/notifications.svelte';

export type NavigationIconName =
	'inventory' | 'catalog' | 'settings' | 'users' | 'requests' | 'notifications' | 'locations';

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
	| '/locations'
	| '/locations/new'
	| '/notifications'
	| '/settings'
	| '/login'
	| '/register';

export type NavigationItem = {
	path: AppPath;
	label: string;
	icon: NavigationIconName;
	requiredRoles?: UserRole[];
	// A count to pin on the item's icon, left out by items that do not have one. Where the value
	// changes, declare it as a getter rather than a fixed number: the read then happens while the
	// link renders, which is what keeps it up to date.
	badgeCount?: number;
	// How to voice that count for screen readers. The wording belongs to whatever is being counted -
	// NavigationLinks only ever knows there is a number - so each item brings its own.
	badgeLabel?: (count: number) => string;
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

const locationEditPageMetadata: PageMetadata = {
	title: 'Izmeni lokaciju',
	description: 'Izmenite detalje o lokaciji i roditeljsku lokaciju.'
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
	'/locations': {
		title: 'Lokacije',
		description: 'Upravljajte lokacijama na kojima se mogu nalaziti stavke.'
	},
	'/locations/new': {
		title: 'Nova lokacija',
		description: 'Odaberite ime lokacije i roditeljsku lokaciju.'
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
	},
	{
		path: '/locations',
		label: 'Lokacije',
		icon: 'locations',
		requiredRoles: ['admin']
	}
];

// Exported on its own as well as through secondaryNavigation, because the mobile header renders the
// bell outside the sidebar list and should not restate the count or its wording.
export const notificationsNavigationItem: NavigationItem = {
	path: '/notifications',
	label: 'Obaveštenja',
	icon: 'notifications',
	// Every signed-in role has notifications; listing them all is what hides the link from
	// signed-out visitors, who can still reach Podešavanja below.
	requiredRoles: ['admin', 'manager', 'user'],
	get badgeCount() {
		return notifications.unreadCount;
	},
	badgeLabel: (count) => `${count} nepročitanih`
};

export const secondaryNavigation: NavigationItem[] = [
	notificationsNavigationItem,
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

	if (pathname === '/locations/new') {
		return pageMetadata['/locations/new'];
	}

	if (pathname.startsWith('/locations/')) {
		return locationEditPageMetadata;
	}

	return pageMetadata[pathname as AppPath] ?? fallbackPageMetadata;
}
