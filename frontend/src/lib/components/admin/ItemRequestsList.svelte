<script lang="ts">
	import Trash2 from '@lucide/svelte/icons/trash-2';
	import Printer from '@lucide/svelte/icons/printer';
	import { Button, Select } from 'bits-ui';
	import { getContext, onDestroy, onMount, tick } from 'svelte';
	import { page } from '$app/state';
	import {
		api,
		ApiError,
		type ItemRequestPreparationReport,
		type ItemRequestSummary,
		type ItemRequestUserOption
	} from '$lib/api';
	import RegistrationApproval from '$lib/components/auth/RegistrationApproval.svelte';
	import MultiOptionCombobox from '$lib/components/shared/MultiOptionCombobox.svelte';
	import PaginationFooter from '$lib/components/shared/PaginationFooter.svelte';
	import {
		itemRequestDateFilterOptions,
		itemRequestDateRange,
		itemRequestStatusClass,
		itemRequestStatusFilterOptions,
		itemRequestStatusLabel,
		type ItemRequestDateFilter,
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
	type PreparationReportTarget = {
		userID: number;
		itemIDs?: number[];
	};

	const requestsPage = createServerPagination<
		ItemRequestSummary,
		{ status: ItemRequestStatusFilter; userIDs: string[]; date: ItemRequestDateFilter }
	>({
		initialFilters: { status: 'all', userIDs: [], date: 'all' },
		loadPage: ({ limit, offset, status, userIDs, date }) =>
			api.listItemRequests({
				limit,
				offset,
				status: status === 'all' ? undefined : status,
				userIDs: userIDs.map(Number),
				...itemRequestDateRange(date)
			}),
		unavailableMessage: 'Zahtevi nisu učitani.'
	});

	let requestUsers = $state<ItemRequestUserOption[]>([]);
	let preparationReports = $state<ItemRequestPreparationReport[]>([]);
	let preparationLoading = $state(false);
	let preparationRequestVersion = 0;

	onMount(() => {
		void loadRequestUsers();
	});

	onDestroy(() => {
		preparationRequestVersion += 1;
		printContext.setPreparationReports([]);
	});

	$effect(() => {
		const url = page.url;
		const requestedStatus = getTableFilter(url, 'status');
		const status = itemRequestStatusFilterOptions.some((option) => option.value === requestedStatus)
			? (requestedStatus as ItemRequestStatusFilter)
			: 'all';
		const userIDs = [...new Set(url.searchParams.getAll('user_id'))].filter((userID) => {
			const parsedUserID = Number(userID);
			return Number.isSafeInteger(parsedUserID) && parsedUserID > 0;
		});
		const requestedDate = getTableFilter(url, 'date');
		const date = itemRequestDateFilterOptions.some((option) => option.value === requestedDate)
			? (requestedDate as ItemRequestDateFilter)
			: 'all';

		requestsPage.sync({ page: getTablePage(url), filters: { status, userIDs, date } });
		preparationRequestVersion += 1;
		preparationReports = [];
		printContext.setPreparationReports([]);
		preparationLoading = false;
	});

	$effect(() => {
		if (requestsPage.error) onloaderror?.(requestsPage.error);
	});

	async function loadRequestUsers() {
		try {
			requestUsers = await api.listItemRequestUsers();
		} catch {
			requestUsers = [];
		}
	}

	async function printPreparationReports() {
		const targets = preparationReportTargets;
		if (targets.length === 0) return;

		const status = requestsPage.filters.status;
		const date = requestsPage.filters.date;
		const version = ++preparationRequestVersion;
		preparationLoading = true;
		preparationReports = [];
		printContext.setPreparationReports([]);

		try {
			const reports = await Promise.all(
				targets.map(({ userID, itemIDs }) =>
					api.getItemRequestPreparationReport(userID, {
						itemIDs,
						status: status === 'all' ? undefined : status,
						...itemRequestDateRange(date)
					})
				)
			);
			if (version === preparationRequestVersion) {
				preparationReports = reports.filter((report) => report.items.length > 0);
				printContext.setPreparationReports(preparationReports);
				if (preparationReports.length === 0) {
					toast.info('Nema zahteva za štampu.');
					return;
				}

				await tick();
				if (version === preparationRequestVersion) printContext.print();
			}
		} catch (reason) {
			if (version === preparationRequestVersion) {
				preparationReports = [];
				printContext.setPreparationReports([]);
				toast.error(reason instanceof ApiError ? reason.message : 'Izveštaj nije učitan.');
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
		} catch (reason) {
			toast.error(reason instanceof ApiError ? reason.message : 'Zahtev nije obrađen.');
		}
	}

	function filterByStatus(status: ItemRequestStatusFilter) {
		updateTableQuery({ status: status === 'all' ? undefined : status, page: 1 });
	}

	function filterByUsers(userIDs: string[]) {
		updateTableQuery({ user_id: userIDs, page: 1 });
	}

	function filterByDate(date: ItemRequestDateFilter) {
		updateTableQuery({ date: date === 'all' ? undefined : date, page: 1 });
	}

	function goToPage(nextPage: number) {
		updateTableQuery({ page: nextPage });
	}

	let selectedUserIDs = $derived(requestsPage.filters.userIDs.map(Number));
	let preparationReportTargets = $derived(
		getPreparationReportTargets(selectedUserIDs, requestsPage.items)
	);

	function getPreparationReportTargets(
		selectedUserIDs: number[],
		shownRequests: ItemRequestSummary[]
	): PreparationReportTarget[] {
		if (selectedUserIDs.length > 0) return selectedUserIDs.map((userID) => ({ userID }));

		const targets: PreparationReportTarget[] = [];
		for (const itemRequest of shownRequests) {
			let target = targets.find(({ userID }) => userID === itemRequest.user_id);
			if (!target) {
				target = { userID: itemRequest.user_id, itemIDs: [] };
				targets.push(target);
			}
			target.itemIDs?.push(itemRequest.item_id);
		}

		return targets;
	}
</script>

<section aria-labelledby="item-requests-heading">
	<h2 id="item-requests-heading" class="sr-only">Zahtevi za stavke</h2>
	<div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between sm:gap-4">
		<p
			class="min-w-0 font-mono text-xs leading-relaxed font-medium tracking-wide text-balance text-muted"
		>
			UKUPNO: {requestsPage.total}
		</p>
		<div
			class="grid w-full grid-cols-2 gap-2 sm:flex sm:w-auto sm:flex-wrap sm:items-center sm:justify-end"
		>
			<Button.Root
				type="button"
				class="order-last col-span-2 inline-flex h-9 w-full items-center justify-center gap-2 rounded-md border border-line bg-surface px-3 text-sm text-ink transition-colors hover:border-brand/40 hover:bg-brand-soft sm:order-none sm:w-auto"
				disabled={requestsPage.items.length === 0 || requestsPage.loading || preparationLoading}
				onclick={() => void printPreparationReports()}
				aria-label="Štampaj"
				title={preparationLoading ? 'Učitavanje izveštaja…' : 'Štampaj'}
			>
				<Printer class="size-4" aria-hidden="true" />
				<span>Štampaj</span>
			</Button.Root>
			<div class="order-first col-span-2 w-full sm:order-none sm:w-64">
				<label class="sr-only" for="item-request-user-filter"
					>Filtriraj zahteve po korisnicima</label
				>
				<MultiOptionCombobox
					id="item-request-user-filter"
					options={requestUsers}
					values={requestsPage.filters.userIDs}
					placeholder="Svi korisnici"
					emptyMessage="Nema korisnika sa zahtevima."
					selectedLabel="Izabrani korisnici"
					showSelected={false}
					showSelectedInMenu
					onvaluechange={filterByUsers}
				/>
			</div>
			<Select.Root
				type="single"
				value={requestsPage.filters.date}
				items={itemRequestDateFilterOptions}
				onValueChange={(value) => filterByDate(value as ItemRequestDateFilter)}
			>
				<Select.Trigger
					class="flex h-9 w-full items-center justify-between rounded-md border border-line bg-surface px-3 text-sm text-ink transition-colors hover:border-brand/40 sm:w-44"
					aria-label="Filtriraj zahteve po datumu"
				>
					<Select.Value />
				</Select.Trigger>
				<Select.Portal>
					<Select.Content
						class="z-30 w-52 rounded-md border border-line bg-surface p-1 shadow-lg shadow-black/15"
						sideOffset={4}
					>
						<Select.Viewport>
							{#each itemRequestDateFilterOptions as option (option.value)}
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
					class="flex h-9 w-full items-center justify-between rounded-md border border-line bg-surface px-3 text-sm text-ink transition-colors hover:border-brand/40 sm:w-40"
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
					<th class="sr-only">Radnje</th>
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
