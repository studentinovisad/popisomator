<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import Plus from '@lucide/svelte/icons/plus';
	import { onMount } from 'svelte';
	import {
		api,
		ApiError,
		type ConsumptionStatus,
		type Item,
		type ItemTypeOption,
		type PropertyOption
	} from '$lib/api';
	import { createAuthPage } from '$lib/state/auth-page.svelte';
	import PaginationFooter from '$lib/components/shared/PaginationFooter.svelte';
	import InventoryList from '$lib/components/inventory/InventoryList.svelte';
	import InventoryToolbar from '$lib/components/inventory/InventoryToolbar.svelte';
	import ProtectedPageState from '$lib/components/shared/ProtectedPageState.svelte';
	import { pagination } from '$lib/state/pagination.svelte';
	import {
		getTableFilter,
		getTablePage,
		getTableSearch,
		updateTableQuery
	} from '$lib/state/table-query';
	import { Portal } from 'bits-ui';

	const authPage = createAuthPage({ unavailableMessage: 'Inventar trenutno nije dostupan.' });

	let items = $state<Item[]>([]);
	let itemTypes = $state<ItemTypeOption[]>([]);
	let properties = $state<PropertyOption[]>([]);
	let loadingInventory = $state(false);
	let loadingInventoryOptions = $state(false);
	let inventoryLoaded = $state(false);
	let inventoryOptionsLoaded = $state(false);
	let inventoryError = $state('');
	let derivedNameSearch = $state('');
	let itemTypeFilter = $state('all');
	let loadVersion = 0;
	let inventoryQueryKey = '';
	let itemsPerPage = $derived(pagination.perPage);
	let itemOffset = $state(0);
	let itemsTotal = $state(0);
	let canManage = $derived(
		authPage.state.user?.role === 'admin' || authPage.state.user?.role === 'manager'
	);
	let currentPage = $derived(Math.floor(itemOffset / itemsPerPage) + 1);
	let hasPreviousPage = $derived(itemOffset > 0);
	let hasNextPage = $derived(itemOffset + items.length < itemsTotal);

	onMount(() => {
		void authPage.load().then(() => {
			if (authPage.state.authorized) {
				void loadInventoryOptions();
			}
		});
	});

	$effect(() => {
		if (!authPage.state.authorized || !inventoryOptionsLoaded) return;

		const url = page.url;
		const search = getTableSearch(url);
		const requestedTypeID = Number.parseInt(getTableFilter(url, 'type_id'), 10);
		const itemTypeID =
			Number.isSafeInteger(requestedTypeID) && requestedTypeID > 0 ? requestedTypeID : undefined;
		const currentPage = getTablePage(url);
		const queryKey = JSON.stringify({ currentPage, search, itemTypeID, itemsPerPage });

		if (queryKey === inventoryQueryKey) return;

		inventoryQueryKey = queryKey;
		derivedNameSearch = search;
		itemTypeFilter = itemTypeID === undefined ? 'all' : String(itemTypeID);
		void loadInventory((currentPage - 1) * itemsPerPage, search, itemTypeID);
	});

	async function loadInventory(offset: number, search: string, itemTypeID: number | undefined) {
		const version = ++loadVersion;
		loadingInventory = true;
		inventoryError = '';

		try {
			const nextItems = await api.listItems({
				limit: itemsPerPage,
				offset,
				search,
				typeID: itemTypeID
			});
			if (version !== loadVersion) return;

			items = nextItems.items;
			itemOffset = nextItems.offset;
			itemsTotal = nextItems.total;
			inventoryLoaded = true;
		} catch (reason) {
			if (version !== loadVersion) return;
			inventoryError = reason instanceof ApiError ? reason.message : 'Stavke nisu učitane.';
		} finally {
			if (version === loadVersion) loadingInventory = false;
		}
	}

	async function loadInventoryOptions() {
		loadingInventoryOptions = true;
		inventoryError = '';

		try {
			const [nextItemTypes, nextProperties] = await Promise.all([
				api.getItemTypeOptions(),
				api.getPropertyOptions()
			]);
			itemTypes = nextItemTypes;
			properties = nextProperties;
			const requestedTypeID = Number.parseInt(getTableFilter(page.url, 'type_id'), 10);
			if (
				nextItemTypes.length === 1 &&
				(!Number.isSafeInteger(requestedTypeID) || requestedTypeID <= 0)
			) {
				updateTableQuery({ type_id: nextItemTypes[0].id });
			}
			inventoryOptionsLoaded = true;
		} catch (reason) {
			inventoryError =
				reason instanceof ApiError ? reason.message : 'Podaci kataloga nisu učitani.';
		} finally {
			loadingInventoryOptions = false;
		}
	}

	async function changeConsumption(item: Item, status: ConsumptionStatus) {
		if (item.consumption === status) return;
		inventoryError = '';

		try {
			await api.consumeItem(item.id, status);
			items = items.map((currentItem) =>
				currentItem.id === item.id ? { ...currentItem, consumption: status } : currentItem
			);
		} catch (reason) {
			inventoryError =
				reason instanceof ApiError ? reason.message : 'Stanje stavke nije promenjeno.';
		}
	}

	function goToPage(page: number) {
		updateTableQuery({ page });
	}

	function searchItems(search: string) {
		updateTableQuery({ search, page: 1 });
	}

	function filterByItemType(itemTypeID: number | undefined) {
		updateTableQuery({ type_id: itemTypeID, page: 1 });
	}

	function refreshInventory() {
		const requestedTypeID = Number.parseInt(getTableFilter(page.url, 'type_id'), 10);
		const itemTypeID =
			Number.isSafeInteger(requestedTypeID) && requestedTypeID > 0 ? requestedTypeID : undefined;
		void loadInventory(itemOffset, derivedNameSearch, itemTypeID);
	}
</script>

<svelte:head>
	<title>Stavke | Popisomator</title>
</svelte:head>

<main class="px-4 pt-4 sm:px-6">
	<ProtectedPageState
		loading={authPage.state.loading ||
			(authPage.state.authorized && (loadingInventory || loadingInventoryOptions))}
		contentLoaded={inventoryLoaded && inventoryOptionsLoaded}
		error={authPage.state.error}
		authorized={authPage.state.authorized && authPage.state.user !== null}
	>
		{#if canManage}
			<Portal to="#page-header-actions">
				<a
					class="inline-flex size-10 items-center justify-center rounded-md bg-brand text-on-brand hover:bg-brand-strong"
					href={resolve('/items/new')}
					aria-label="Dodaj stavku"
					title="Dodaj stavku"
				>
					<Plus class="size-4" aria-hidden="true" />
				</a>
			</Portal>
		{/if}

		{#if inventoryError}<p class="mt-3 text-sm text-danger" role="alert">{inventoryError}</p>{/if}
		<InventoryToolbar
			total={itemsTotal}
			{itemTypes}
			typeFilter={itemTypeFilter}
			bind:search={derivedNameSearch}
			loading={loadingInventory}
			onitemtypechange={filterByItemType}
			onsearch={searchItems}
		/>
		<InventoryList
			{items}
			{itemTypes}
			{properties}
			{canManage}
			currentUserID={authPage.state.user?.id}
			onconsumptionchange={changeConsumption}
			onrequested={refreshInventory}
		/>
		<PaginationFooter
			total={itemsTotal}
			perPage={itemsPerPage}
			page={currentPage}
			{hasPreviousPage}
			{hasNextPage}
			loading={loadingInventory}
			onpagechange={goToPage}
		/>
	</ProtectedPageState>
</main>
