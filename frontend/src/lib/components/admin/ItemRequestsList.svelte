<script lang="ts">
	import Trash2 from '@lucide/svelte/icons/trash-2';
	import { Button, Select } from 'bits-ui';
	import { SvelteMap } from 'svelte/reactivity';
	import { page } from '$app/state';
	import { api, ApiError, type Item, type ItemRequest, type User } from '$lib/api';
	import RegistrationApproval from '$lib/components/auth/RegistrationApproval.svelte';
	import PaginationFooter from '$lib/components/shared/PaginationFooter.svelte';
	import {
		itemRequestStatusClass,
		itemRequestStatusFilterOptions,
		itemRequestStatusLabel,
		type ItemRequestStatusFilter
	} from '$lib/domain/item-requests';
	import { createServerPagination } from '$lib/state/server-pagination.svelte';
	import { getTableFilter, getTablePage, updateTableQuery } from '$lib/state/table-query';

	const requestsPage = createServerPagination<ItemRequest, { status: ItemRequestStatusFilter }>({
		initialFilters: { status: 'requested' },
		loadPage: ({ limit, offset, status }) =>
			api.listItemRequests({ limit, offset, status: status === 'all' ? undefined : status }),
		unavailableMessage: 'Zahtevi nisu učitani.'
	});

	const itemCache = new SvelteMap<number, Item>();
	const userCache = new SvelteMap<number, User>();
	let userDirectoryLoaded = $state(false);

	$effect(() => {
		const url = page.url;
		const requestedStatus = getTableFilter(url, 'status');
		const status = itemRequestStatusFilterOptions.some((option) => option.value === requestedStatus)
			? (requestedStatus as ItemRequestStatusFilter)
			: 'requested';

		requestsPage.sync({ page: getTablePage(url), filters: { status } });
	});

	$effect(() => {
		void loadMissingItems(requestsPage.items.map((itemRequest) => itemRequest.item_id));
	});

	$effect(() => {
		if (!userDirectoryLoaded) void loadUserDirectory();
	});

	async function loadMissingItems(itemIDs: number[]) {
		const missing = [...new Set(itemIDs)].filter((id) => !itemCache.has(id));
		if (missing.length === 0) return;

		const fetched = await Promise.all(missing.map((id) => api.getItem(id).catch(() => null)));
		for (const item of fetched) if (item) itemCache.set(item.id, item);
	}

	async function loadUserDirectory() {
		userDirectoryLoaded = true;
		let offset = 0;

		for (let iteration = 0; iteration < 10; iteration++) {
			let batch;
			try {
				batch = await api.listUsers({ limit: 50, offset });
			} catch {
				break;
			}

			for (const user of batch.items) userCache.set(user.id, user);
			offset += batch.items.length;
			if (batch.items.length === 0 || offset >= batch.total) break;
		}
	}

	function itemName(itemID: number) {
		return itemCache.get(itemID)?.derived_name || `Stavka #${itemID}`;
	}

	function userName(userID: number) {
		const user = userCache.get(userID);
		return user ? `${user.full_name}` : `Korisnik #${userID}`;
	}

	async function decide(itemRequest: ItemRequest, approve: boolean) {
		requestsPage.error = '';

		try {
			if (approve) {
				await api.approveItemRequest(itemRequest.user_id, itemRequest.item_id);
			} else {
				await api.denyItemRequest(itemRequest.user_id, itemRequest.item_id);
			}
			requestsPage.reloadAfterDelete();
		} catch (reason) {
			requestsPage.error = reason instanceof ApiError ? reason.message : 'Zahtev nije obrađen.';
		}
	}

	function filterByStatus(status: ItemRequestStatusFilter) {
		updateTableQuery({ status: status === 'requested' ? undefined : status, page: 1 });
	}

	function goToPage(nextPage: number) {
		updateTableQuery({ page: nextPage });
	}
</script>

<section aria-labelledby="item-requests-heading">
	<h2 id="item-requests-heading" class="sr-only">Zahtevi za stavke</h2>
	<div class="flex items-center justify-between gap-4">
		<p class="font-mono text-xs font-medium tracking-wide text-muted">
			UKUPNO: {requestsPage.total}
		</p>
		<Select.Root
			type="single"
			value={requestsPage.filters.status}
			items={itemRequestStatusFilterOptions}
			onValueChange={(value) => filterByStatus(value as ItemRequestStatusFilter)}
		>
			<Select.Trigger
				class="flex h-9 w-40 items-center justify-between rounded-md border border-line bg-surface px-3 text-sm text-ink transition-colors hover:border-brand/40"
				aria-label="Filtriraj zahteve po statusu"
			>
				<Select.Value />
			</Select.Trigger>
			<Select.Portal>
				<Select.Content
					class="z-30 w-44 rounded-md border border-line bg-surface p-1 shadow-lg shadow-black/15"
					sideOffset={4}
				>
					<Select.Viewport>
						{#each itemRequestStatusFilterOptions as option (option.value)}
							<Select.Item
								value={option.value}
								label={option.label}
								class="cursor-pointer rounded px-3 py-2 text-sm outline-none data-highlighted:bg-brand-soft"
							>
								{option.label}
							</Select.Item>
						{/each}
					</Select.Viewport>
				</Select.Content>
			</Select.Portal>
		</Select.Root>
	</div>
	{#if requestsPage.error}
		<p class="mt-3 text-sm text-danger" role="alert">{requestsPage.error}</p>
	{/if}
	<div class="-mx-4 mt-4 border-y border-line bg-surface sm:-mx-6">
		<table class="hidden w-full table-fixed text-left text-sm lg:table">
			<colgroup>
				<col />
				<col class="w-48" />
				<col />
				<col class="w-32" />
				<col class="w-32" />
			</colgroup>
			<thead class="border-b border-line bg-soft text-muted">
				<tr class="h-12">
					<th class="px-4 py-3 font-medium">Stavka</th>
					<th class="px-4 py-3 font-medium">Korisnik</th>
					<th class="px-4 py-3 font-medium">Razlog</th>
					<th class="px-4 py-3 font-medium">Status</th>
					<th class="px-4 py-3 text-right font-medium">Radnje</th>
				</tr>
			</thead>
			<tbody class="divide-y divide-line text-ink">
				{#each requestsPage.items as itemRequest (`${itemRequest.user_id}:${itemRequest.item_id}`)}
					<tr class="h-16">
						<td class="px-4 py-3 align-middle">
							<span class="block truncate">{itemName(itemRequest.item_id)}</span>
						</td>
						<td class="px-4 py-3 align-middle">
							<span class="block truncate">{userName(itemRequest.user_id)}</span>
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
						<td class="px-4 py-3 align-middle">
							{#if itemRequest.status === 'requested'}
								<RegistrationApproval onclick={(approved) => void decide(itemRequest, approved)} />
							{:else}
								<div class="flex justify-end">
									<Button.Root
										class="inline-flex size-9 items-center justify-center rounded-md border border-danger bg-surface text-danger transition-colors hover:bg-danger-soft"
										onclick={() => void decide(itemRequest, false)}
										aria-label="Ukloni odobrenje"
										title="Ukloni odobrenje"
									>
										<Trash2 class="size-4" aria-hidden="true" />
									</Button.Root>
								</div>
							{/if}
						</td>
					</tr>
				{/each}
				{#if requestsPage.items.length === 0}
					<tr class="h-16"><td class="px-4 py-3 text-muted" colspan="5">Nema zahteva.</td></tr>
				{/if}
			</tbody>
		</table>
		<ul class="divide-y divide-line lg:hidden" aria-label="Zahtevi za stavke">
			{#each requestsPage.items as itemRequest (`${itemRequest.user_id}:${itemRequest.item_id}`)}
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
					<p class="mt-0.5 truncate text-sm text-muted">{userName(itemRequest.user_id)}</p>
					<p class="mt-0.5 truncate text-sm text-muted" title={itemRequest.reason}>
						{itemRequest.reason}
					</p>
					<div class="mt-3">
						{#if itemRequest.status === 'requested'}
							<RegistrationApproval onclick={(approved) => void decide(itemRequest, approved)} />
						{:else}
							<div class="flex justify-end">
								<Button.Root
									class="inline-flex size-9 items-center justify-center rounded-md border border-danger bg-surface text-danger transition-colors hover:bg-danger-soft"
									onclick={() => void decide(itemRequest, false)}
									aria-label="Ukloni odobrenje"
									title="Ukloni odobrenje"
								>
									<Trash2 class="size-4" aria-hidden="true" />
								</Button.Root>
							</div>
						{/if}
					</div>
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
</section>
