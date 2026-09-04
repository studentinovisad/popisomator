<script lang="ts">
	import Trash2 from '@lucide/svelte/icons/trash-2';
	import Printer from '@lucide/svelte/icons/printer';
	import { Button, Select } from 'bits-ui';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { page } from '$app/state';
	import {
		api,
		ApiError,
		type ItemRequestPreparationReport,
		type ItemRequestSummary,
		type ItemRequestUserOption
	} from '$lib/api';
	import RegistrationApproval from '$lib/components/auth/RegistrationApproval.svelte';
	import PaginationFooter from '$lib/components/shared/PaginationFooter.svelte';
	import {
		itemRequestStatusClass,
		itemRequestStatusFilterOptions,
		itemRequestStatusLabel,
		type ItemRequestStatusFilter
	} from '$lib/domain/item-requests';
	import { createServerPagination } from '$lib/state/server-pagination.svelte';
	import {
		preparationReportPrintContextKey,
		type PreparationReportPrintContext
	} from '$lib/state/preparation-report-print-context';
	import { getTableFilter, getTablePage, updateTableQuery } from '$lib/state/table-query';
	import { toast } from 'svelte-sonner';

	const printContext = getContext<PreparationReportPrintContext>(preparationReportPrintContextKey);
	let { onloaderror }: { onloaderror?: (message: string) => void } = $props();

	const requestsPage = createServerPagination<
		ItemRequestSummary,
		{ status: ItemRequestStatusFilter; userID: string }
	>({
		initialFilters: { status: 'all', userID: 'all' },
		loadPage: ({ limit, offset, status, userID }) =>
			api.listItemRequests({
				limit,
				offset,
				status: status === 'all' ? undefined : status,
				userID: userID === 'all' ? undefined : Number(userID)
			}),
		unavailableMessage: 'Zahtevi nisu učitani.'
	});

	let requestUsers = $state<ItemRequestUserOption[]>([]);
	let preparationReport = $state<ItemRequestPreparationReport | null>(null);
	let preparationLoading = $state(false);
	let preparationRequestVersion = 0;
	const userFilterOptions = $derived([
		{ value: 'all', label: 'Svi korisnici' },
		...requestUsers.map((user) => ({ value: String(user.id), label: user.name }))
	]);

	onMount(() => {
		void loadRequestUsers();
	});

	onDestroy(() => printContext.setPreparationReport(null));

	$effect(() => {
		const url = page.url;
		const requestedStatus = getTableFilter(url, 'status');
		const status = itemRequestStatusFilterOptions.some((option) => option.value === requestedStatus)
			? (requestedStatus as ItemRequestStatusFilter)
			: 'all';
		const requestedUserID = getTableFilter(url, 'user_id');
		const parsedUserID = Number(requestedUserID);
		const userID = Number.isSafeInteger(parsedUserID) && parsedUserID > 0 ? requestedUserID : 'all';

		requestsPage.sync({ page: getTablePage(url), filters: { status, userID } });
	});

	$effect(() => {
		if (requestsPage.error) onloaderror?.(requestsPage.error);
	});

	$effect(() => {
		const userID =
			requestsPage.filters.userID === 'all' ? undefined : Number(requestsPage.filters.userID);
		refreshPreparationReport(userID);
	});

	function refreshPreparationReport(userID?: number) {
		const version = ++preparationRequestVersion;
		preparationReport = null;
		printContext.setPreparationReport(null);
		preparationLoading = false;

		if (!userID) return;

		preparationLoading = true;
		void loadPreparationReport(userID, version);
	}

	async function loadRequestUsers() {
		try {
			requestUsers = await api.listItemRequestUsers();
		} catch {
			requestUsers = [];
		}
	}

	async function loadPreparationReport(userID: number, version: number) {
		try {
			const report = await api.getItemRequestPreparationReport(userID);
			if (version === preparationRequestVersion) {
				preparationReport = report;
				printContext.setPreparationReport(report);
			}
		} catch {
			if (version === preparationRequestVersion) {
				preparationReport = null;
				printContext.setPreparationReport(null);
			}
		} finally {
			if (version === preparationRequestVersion) preparationLoading = false;
		}
	}

	async function decide(itemRequest: ItemRequestSummary, approve: boolean) {
		try {
			if (approve) {
				await api.approveItemRequest(itemRequest.user_id, itemRequest.item_id);
			} else {
				await api.denyItemRequest(itemRequest.user_id, itemRequest.item_id);
			}
			toast.success(approve ? 'Zahtev je odobren.' : 'Zahtev je odbijen.');
			requestsPage.reloadAfterDelete();
			refreshPreparationReport(selectedUserID);
		} catch (reason) {
			toast.error(reason instanceof ApiError ? reason.message : 'Zahtev nije obrađen.');
		}
	}

	function filterByStatus(status: ItemRequestStatusFilter) {
		updateTableQuery({ status: status === 'all' ? undefined : status, page: 1 });
	}

	function filterByUser(userID: string) {
		updateTableQuery({ user_id: userID === 'all' ? undefined : userID, page: 1 });
	}

	function goToPage(nextPage: number) {
		updateTableQuery({ page: nextPage });
	}

	let selectedUserID = $derived(
		requestsPage.filters.userID === 'all' ? undefined : Number(requestsPage.filters.userID)
	);
</script>

<section aria-labelledby="item-requests-heading">
	<h2 id="item-requests-heading" class="sr-only">Zahtevi za stavke</h2>
	<div class="flex items-center justify-between gap-4">
		<p class="font-mono text-xs font-medium tracking-wide text-muted">
			UKUPNO: {requestsPage.total}
		</p>
		<div class="flex items-center gap-2">
			{#if selectedUserID}
				<Button.Root
					type="button"
					class="inline-flex h-9 items-center justify-center gap-2 rounded-md border border-line bg-surface px-3 text-sm text-ink transition-colors hover:border-brand/40 hover:bg-brand-soft"
					disabled={!preparationReport || preparationLoading}
					onclick={() => printContext.print()}
					aria-label="Štampaj pripremu"
					title={preparationLoading ? 'Učitavanje pripreme…' : 'Štampaj pripremu'}
				>
					<Printer class="size-4" aria-hidden="true" />
					<span class="max-sm:sr-only"
						>{preparationLoading ? 'Učitavanje…' : 'Štampaj pripremu'}</span
					>
				</Button.Root>
			{/if}
			<Select.Root
				type="single"
				value={requestsPage.filters.userID}
				items={userFilterOptions}
				onValueChange={(value) => filterByUser(value)}
			>
				<Select.Trigger
					class="flex h-9 w-44 items-center justify-between rounded-md border border-line bg-surface px-3 text-sm text-ink transition-colors hover:border-brand/40"
					aria-label="Filtriraj zahteve po korisniku"
				>
					<Select.Value />
				</Select.Trigger>
				<Select.Portal>
					<Select.Content
						class="z-30 w-52 rounded-md border border-line bg-surface p-1 shadow-lg shadow-black/15"
						sideOffset={4}
					>
						<Select.Viewport>
							{#each userFilterOptions as option (option.value)}
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
	</div>
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
							<span class="block truncate"
								>{itemRequest.item_name ?? `Stavka #${itemRequest.item_id}`}</span
							>
						</td>
						<td class="px-4 py-3 align-middle">
							<span class="block truncate"
								>{itemRequest.user_name ?? `Korisnik #${itemRequest.user_id}`}</span
							>
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
							{itemRequest.item_name ?? `Stavka #${itemRequest.item_id}`}
						</p>
						<span
							class={`shrink-0 rounded px-2 py-1 text-xs font-medium ${itemRequestStatusClass(itemRequest.status)}`}
						>
							{itemRequestStatusLabel(itemRequest.status)}
						</span>
					</div>
					<p class="mt-0.5 truncate text-sm text-muted">
						{itemRequest.user_name ?? `Korisnik #${itemRequest.user_id}`}
					</p>
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
