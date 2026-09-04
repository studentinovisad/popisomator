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
	#queryKey = '';
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
			const message = reason instanceof ApiError ? reason.message : this.#unavailableMessage;
			this.error = message;
		} finally {
			if (version === this.#loadVersion) this.loading = false;
		}
	};

	sync = ({ page, search, filters }: { page: number; search?: string; filters?: Filters }) => {
		const nextPage = Math.max(1, page);
		const nextSearch = search ?? '';
		const nextFilters = filters ?? this.filters;
		const queryKey = JSON.stringify({
			page: nextPage,
			search: nextSearch,
			filters: nextFilters,
			perPage: this.perPage
		});

		if (queryKey === this.#queryKey) return;

		this.#queryKey = queryKey;
		this.search = nextSearch;
		this.filters = nextFilters;
		void this.load((nextPage - 1) * this.perPage);
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
