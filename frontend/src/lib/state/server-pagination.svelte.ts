import { ApiError } from '$lib/api';
import { pagination } from '$lib/state/pagination.svelte';

type ServerPage<T> = {
	items: T[];
	offset: number;
	total: number;
};

type ServerPageLoader<T> = (parameters: {
	limit: number;
	offset: number;
}) => Promise<ServerPage<T>>;

type ServerPaginationOptions<T> = {
	loadPage: ServerPageLoader<T>;
	unavailableMessage: string;
};

export class ServerPagination<T> {
	items = $state<T[]>([]);
	loading = $state(false);
	error = $state('');
	offset = $state(0);
	total = $state(0);
	perPage = $derived(pagination.perPage);
	currentPage = $derived(Math.floor(this.offset / this.perPage) + 1);
	hasPreviousPage = $derived(this.offset > 0);
	hasNextPage = $derived(this.offset + this.items.length < this.total);

	#loadVersion = 0;
	#loadPage: ServerPageLoader<T>;
	#unavailableMessage: string;

	constructor({ loadPage, unavailableMessage }: ServerPaginationOptions<T>) {
		this.#loadPage = loadPage;
		this.#unavailableMessage = unavailableMessage;
	}

	load = async (offset = 0) => {
		const version = ++this.#loadVersion;
		this.loading = true;
		this.error = '';

		try {
			const page = await this.#loadPage({ limit: this.perPage, offset });
			if (version !== this.#loadVersion) return;

			this.items = page.items;
			this.offset = page.offset;
			this.total = page.total;
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

	reloadAfterDelete = () => {
		const nextOffset =
			this.items.length === 1 && this.offset > 0 ? this.offset - this.perPage : this.offset;
		void this.load(nextOffset);
	};
}

export function createServerPagination<T>(options: ServerPaginationOptions<T>) {
	return new ServerPagination(options);
}
