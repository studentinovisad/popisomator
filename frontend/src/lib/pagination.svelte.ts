const storageKey = 'popisomator-rows-per-page';
const defaultPerPage = 20;
const minimumPerPage = 5;
const maximumPerPage = 50;

function normalizePerPage(value: number) {
	return Math.min(maximumPerPage, Math.max(minimumPerPage, value));
}

class PaginationPreference {
	perPage = $state(defaultPerPage);

	constructor() {
		if (typeof window === 'undefined') return;

		const storedPerPage = Number(localStorage.getItem(storageKey));
		if (Number.isInteger(storedPerPage)) {
			this.perPage = normalizePerPage(storedPerPage);
		}
	}

	set(nextPerPage: number) {
		if (!Number.isInteger(nextPerPage)) return;

		this.perPage = normalizePerPage(nextPerPage);
		localStorage.setItem(storageKey, String(this.perPage));
	}
}

export const pagination = new PaginationPreference();
