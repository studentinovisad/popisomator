import { ApiError } from '$lib/api';
import { pagination } from '$lib/state/pagination.svelte';

type ServerPage<T> = {
	items: T[];
	offset: number;
	total: number;
};

type ServerPageLoader<T, Filters extends object> = (
	parameters: {
		limit: number;
		offset: number;
		search: string;
	} & Filters
) => Promise<ServerPage<T>>;

type ServerPaginationOptions<T, Filters extends object> = {
	loadPage: ServerPageLoader<T, Filters>;
	unavailableMessage: string;
	initialFilters?: Filters;
};

export class ServerPagination<T, Filters extends object = Record<string, never>> {
	items = $state<T[]>([]);
	loading = $state(false);
	loaded = $state(false);
	error = $state('');
	offset = $state(0);
	total = $state(0);
	search = $state('');
	filters = $state<Filters>({} as Filters);
	perPage = $derived(pagination.perPage);
	currentPage = $derived(Math.floor(this.offset / this.perPage) + 1);
	hasPreviousPage = $derived(this.offset > 0);
	hasNextPage = $derived(this.offset + this.items.length < this.total);

	#loadVersion = 0;
	#loadPage: ServerPageLoader<T, Filters>;
	#unavailableMessage: string;

	constructor({
		loadPage,
		unavailableMessage,
		initialFilters
	}: ServerPaginationOptions<T, Filters>) {
		this.#loadPage = loadPage;
		this.#unavailableMessage = unavailableMessage;
		this.filters = initialFilters ?? ({} as Filters);
	}

	load = async (offset = 0) => {
		const version = ++this.#loadVersion;
		this.loading = true;
		this.error = '';

		try {
			const page = await this.#loadPage({
				limit: this.perPage,
				offset,
				search: this.search,
				...this.filters
			});
			if (version !== this.#loadVersion) return;

			this.items = page.items;
			this.offset = page.offset;
			this.total = page.total;
			this.loaded = true;
		} catch (reason) {
			if (version !== this.#loadVersion) return;
			this.error = reason instanceof ApiError ? reason.message : this.#unavailableMessage;
		} finally {
			if (version === this.#loadVersion) this.loading = false;
		}
	};

	goToPage = (page: number) => {
		void this.load((page - 1) * this.perPage);
	};

	searchBy = (search: string) => {
		this.search = search;
		void this.load();
	};

	setFilters = (filters: Partial<Filters>) => {
		this.filters = { ...this.filters, ...filters };
		void this.load();
	};

	reloadAfterDelete = () => {
		const nextOffset =
			this.items.length === 1 && this.offset > 0 ? this.offset - this.perPage : this.offset;
		void this.load(nextOffset);
	};
}

export function createServerPagination<T, Filters extends object = Record<string, never>>(
	options: ServerPaginationOptions<T, Filters>
) {
	return new ServerPagination(options);
}
