<script lang="ts">
	import { onMount } from 'svelte';
	import { SvelteMap } from 'svelte/reactivity';
	import { page } from '$app/state';
	import { api, type Item, type ItemRequest } from '$lib/api';
	import PaginationFooter from '$lib/components/shared/PaginationFooter.svelte';
	import ProtectedPageState from '$lib/components/shared/ProtectedPageState.svelte';
	import {
		formatRequestDate,
		itemRequestStatusClass,
		itemRequestStatusLabel
	} from '$lib/domain/item-requests';
	import { createAuthPage } from '$lib/state/auth-page.svelte';
	import { createServerPagination } from '$lib/state/server-pagination.svelte';
	import { getTablePage, updateTableQuery } from '$lib/state/table-query';

	const authPage = createAuthPage({ unavailableMessage: 'Zahtevi trenutno nisu dostupni.' });

	const requestsPage = createServerPagination<ItemRequest>({
		loadPage: ({ limit, offset }) => api.listPersonalItemRequests({ limit, offset }),
		unavailableMessage: 'Zahtevi nisu učitani.'
	});

	const itemCache = new SvelteMap<number, Item>();

	onMount(() => void authPage.load());

	$effect(() => {
		requestsPage.sync({ page: getTablePage(page.url) });
	});

	$effect(() => {
		void loadMissingItems(requestsPage.items.map((r) => r.item_id));
	});

	async function loadMissingItems(itemIDs: number[]) {
		const missing = [...new Set(itemIDs)].filter((id) => !itemCache.has(id));
		if (missing.length === 0) return;

		const fetched = await Promise.all(missing.map((id) => api.getItem(id).catch(() => null)));
		for (const item of fetched) if (item) itemCache.set(item.id, item);
	}

	function itemName(itemID: number) {
		return itemCache.get(itemID)?.derived_name || `Stavka #${itemID}`;
	}

	function goToPage(nextPage: number) {
		updateTableQuery({ page: nextPage });
	}
</script>

<svelte:head>
	<title>Zahtevi | Popisomator</title>
</svelte:head>

<main class="px-4 pt-4 pb-8 sm:px-6">
	<ProtectedPageState
		loading={authPage.state.loading}
		error={authPage.state.error}
		authorized={authPage.state.authorized}
	>
		<p class="font-mono text-xs font-medium tracking-wide text-muted">
			UKUPNO: {requestsPage.total}
		</p>
		{#if requestsPage.error}
			<p class="mt-3 text-sm text-danger" role="alert">{requestsPage.error}</p>
		{/if}
		<div class="-mx-4 mt-4 border-y border-line bg-surface sm:-mx-6">
			<table class="hidden w-full table-fixed text-left text-sm lg:table">
				<colgroup>
					<col />
					<col />
					<col class="w-36" />
					<col class="w-44" />
				</colgroup>
				<thead class="border-b border-line bg-soft text-muted">
					<tr class="h-12">
						<th class="px-4 py-3 font-medium">Stavka</th>
						<th class="px-4 py-3 font-medium">Razlog</th>
						<th class="px-4 py-3 font-medium">Status</th>
						<th class="px-4 py-3 font-medium">Datum</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-line text-ink">
					{#each requestsPage.items as itemRequest (itemRequest.item_id)}
						<tr class="h-16">
							<td class="px-4 py-3 align-middle">
								<span class="block truncate">{itemName(itemRequest.item_id)}</span>
							</td>
							<td class="px-4 py-3 align-middle">
								<span class="block truncate" title={itemRequest.reason}>{itemRequest.reason}</span>
							</td>
							<td class="px-4 py-3 align-middle">
								<span
									class={`inline-flex rounded px-2 py-1 text-xs font-medium ${itemRequestStatusClass(itemRequest.status)}`}
								>
									{itemRequestStatusLabel(itemRequest.status)}
								</span>
							</td>
							<td class="px-4 py-3 align-middle text-muted">
								{formatRequestDate(itemRequest.created_at)}
							</td>
						</tr>
					{/each}
					{#if requestsPage.items.length === 0}
						<tr class="h-16"><td class="px-4 py-3 text-muted" colspan="4">Nema zahteva.</td></tr>
					{/if}
				</tbody>
			</table>
			<ul class="divide-y divide-line lg:hidden" aria-label="Zahtevi">
				{#each requestsPage.items as itemRequest (itemRequest.item_id)}
					<li class="px-4 py-3">
						<div class="flex items-start justify-between gap-3">
							<p class="min-w-0 truncate text-sm font-medium text-ink">
								{itemName(itemRequest.item_id)}
							</p>
							<span
								class={`shrink-0 rounded px-2 py-1 text-xs font-medium ${itemRequestStatusClass(itemRequest.status)}`}
							>
								{itemRequestStatusLabel(itemRequest.status)}
							</span>
						</div>
						<p class="mt-0.5 truncate text-sm text-muted" title={itemRequest.reason}>
							{itemRequest.reason}
						</p>
						<p class="mt-1 text-xs text-muted">{formatRequestDate(itemRequest.created_at)}</p>
					</li>
				{/each}
				{#if requestsPage.items.length === 0}
					<li class="px-4 py-3 text-sm text-muted">Nema zahteva.</li>
				{/if}
			</ul>
		</div>
		<PaginationFooter
			total={requestsPage.total}
			perPage={requestsPage.perPage}
			page={requestsPage.currentPage}
			hasPreviousPage={requestsPage.hasPreviousPage}
			hasNextPage={requestsPage.hasNextPage}
			loading={requestsPage.loading}
			onpagechange={goToPage}
		/>
	</ProtectedPageState>
</main>
