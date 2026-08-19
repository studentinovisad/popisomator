import { goto } from '$app/navigation';
import { resolve } from '$app/paths';
import { page } from '$app/state';
import type { AppPath } from '$lib/domain/navigation';

type TableQueryValue = string | number | null | undefined;

export function getTablePage(url: URL) {
	const value = Number.parseInt(url.searchParams.get('page') ?? '', 10);
	return Number.isSafeInteger(value) && value > 0 ? value : 1;
}

export function getTableSearch(url: URL) {
	return url.searchParams.get('search')?.trim() ?? '';
}

export function getTableFilter(url: URL, key: string) {
	return url.searchParams.get(key)?.trim() ?? '';
}

export function updateTableQuery(values: Record<string, TableQueryValue>) {
	const nextURL = new URL(page.url);

	for (const [key, value] of Object.entries(values)) {
		if (value === undefined || value === null || value === '' || value === 1) {
			nextURL.searchParams.delete(key);
		} else {
			nextURL.searchParams.set(key, String(value));
		}
	}

	if (nextURL.href === page.url.href) return;

	void goto(resolve(`${nextURL.pathname}${nextURL.search}${nextURL.hash}` as AppPath), {
		keepFocus: true,
		noScroll: true
	});
}
