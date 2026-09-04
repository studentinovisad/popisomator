<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import Plus from '@lucide/svelte/icons/plus';
	import { onMount } from 'svelte';
	import { SvelteMap } from 'svelte/reactivity';
	import { toast } from 'svelte-sonner';
	import {
		api,
		ApiError,
		type ConsumptionStatus,
		type Item,
		type ItemPropertyTotal,
		type ItemType,
		type ItemTypeOption,
		type ItemTypeFilterableProperty,
		type PropertyOption,
		type PropertyValue
	} from '$lib/api';
	import { createAuthPage } from '$lib/state/auth-page.svelte';
	import PaginationFooter from '$lib/components/shared/PaginationFooter.svelte';
	import InventoryList from '$lib/components/inventory/InventoryList.svelte';
	import InventoryToolbar from '$lib/components/inventory/InventoryToolbar.svelte';
	import InventoryPropertyFilters, {
		type PropertyFilter
	} from '$lib/components/inventory/InventoryPropertyFilters.svelte';
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
	let selectedItemType = $state<ItemType | null>(null);
	let itemTypeFilterableProperties = $state<ItemTypeFilterableProperty[]>([]);
	let propertyFilterOptions = $state<Record<number, PropertyValue[]>>({});
	let loadingInventory = $state(false);
	let loadingInventoryOptions = $state(false);
	let inventoryLoaded = $state(false);
	let inventoryOptionsLoaded = $state(false);
	let inventoryError = $state('');
	let derivedNameSearch = $state('');
	let loadVersion = 0;
	let inventoryQueryKey = '';
	let itemTypeFiltersLoadVersion = 0;
	let propertyFilterValueLoadVersions = new SvelteMap<string, number>();
	let itemsPerPage = $derived(pagination.perPage);
	let itemOffset = $state(0);
	let itemsTotal = $state(0);
	let itemPropertyTotals = $state<ItemPropertyTotal[]>([]);
	let canManage = $derived(
		authPage.state.user?.role === 'admin' || authPage.state.user?.role === 'manager'
	);
	let currentPage = $derived(Math.floor(itemOffset / itemsPerPage) + 1);
	let hasPreviousPage = $derived(itemOffset > 0);
	let hasNextPage = $derived(itemOffset + items.length < itemsTotal);
	let selectedItemTypeID = $derived(getSelectedItemTypeID(page.url));
	let itemTypeFilter = $derived(
		selectedItemTypeID === undefined ? 'all' : String(selectedItemTypeID)
	);
	let propertyOptionsByID = $derived(
		new Map(properties.map((property) => [property.id, property]))
	);
	let propertyFilters = $derived.by((): PropertyFilter[] => {
		if (!selectedItemType) return [];

		const valueCounts = new Map(
			itemTypeFilterableProperties.map((property) => [property.property_id, property.value_count])
		);
		const selectedFilters = getPropertyFilters(page.url);
		return selectedItemType.properties.flatMap((itemTypeProperty) => {
			const property = propertyOptionsByID.get(itemTypeProperty.id);
			if (!property || (valueCounts.get(itemTypeProperty.id) ?? 0) < 2) return [];
			return [
				{
					id: property.id,
					name: property.name,
					value: selectedFilters[property.id] ?? undefined,
					value_type: property.value_type
				}
			];
		});
	});

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
		const itemTypeID = getSelectedItemTypeID(url);
		const selectedPropertyFilters = getPropertyFilters(url);
		const currentPage = getTablePage(url);
		const queryKey = JSON.stringify({
			currentPage,
			search,
			itemTypeID,
			selectedPropertyFilters,
			itemsPerPage
		});

		if (queryKey === inventoryQueryKey) return;

		inventoryQueryKey = queryKey;
		derivedNameSearch = search;
		void loadInventory(
			(currentPage - 1) * itemsPerPage,
			search,
			itemTypeID,
			selectedPropertyFilters
		);
	});

	$effect(() => {
		if (!authPage.state.authorized || !inventoryOptionsLoaded) return;

		const itemTypeID = getSelectedItemTypeID(page.url);
		if (itemTypeID === undefined) {
			itemTypeFiltersLoadVersion += 1;
			selectedItemType = null;
			itemTypeFilterableProperties = [];
			propertyFilterOptions = {};
			propertyFilterValueLoadVersions = new SvelteMap();
			return;
		}

		void loadItemTypeFilters(itemTypeID);
	});

	async function loadInventory(
		offset: number,
		search: string,
		itemTypeID: number | undefined,
		selectedPropertyFilters: Record<number, PropertyValue>
	) {
		const version = ++loadVersion;
		loadingInventory = true;
		inventoryError = '';

		try {
			const nextItems = await api.listItems({
				limit: itemsPerPage,
				offset,
				search,
				typeID: itemTypeID,
				propertyFilters: selectedPropertyFilters
			});
			if (version !== loadVersion) return;

			items = nextItems.items;
			itemOffset = nextItems.offset;
			itemsTotal = nextItems.total;
			itemPropertyTotals = nextItems.property_totals ?? [];
			inventoryLoaded = true;
		} catch (reason) {
			if (version !== loadVersion) return;
			reportInventoryError(reason instanceof ApiError ? reason.message : 'Stavke nisu učitane.');
		} finally {
			if (version === loadVersion) loadingInventory = false;
		}
	}

	async function loadItemTypeFilters(itemTypeID: number) {
		const version = ++itemTypeFiltersLoadVersion;
		let itemTypeLoaded = false;
		try {
			const itemType = await api.getItemType(itemTypeID);
			if (version !== itemTypeFiltersLoadVersion) return;
			selectedItemType = itemType;
			itemTypeLoaded = true;
			const filterableProperties = await api.listItemTypeFilterableProperties(itemTypeID);
			if (version !== itemTypeFiltersLoadVersion) return;
			itemTypeFilterableProperties = filterableProperties;
			propertyFilterOptions = {};
			propertyFilterValueLoadVersions = new SvelteMap();
		} catch (reason) {
			if (version !== itemTypeFiltersLoadVersion) return;
			if (!itemTypeLoaded) selectedItemType = null;
			itemTypeFilterableProperties = [];
			propertyFilterOptions = {};
			reportInventoryError(
				reason instanceof ApiError ? reason.message : 'Filteri svojstava nisu učitani.'
			);
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
			const requestedTypeID = getSelectedItemTypeID(page.url);
			if (nextItemTypes.length === 1 && requestedTypeID === undefined) {
				updateTableQuery({ type_id: nextItemTypes[0].id });
			}
			inventoryOptionsLoaded = true;
		} catch (reason) {
			reportInventoryError(
				reason instanceof ApiError ? reason.message : 'Podaci kataloga nisu učitani.'
			);
		} finally {
			loadingInventoryOptions = false;
		}
	}

	async function changeConsumption(item: Item, status: ConsumptionStatus) {
		if (item.consumption === status) return;

		try {
			await api.consumeItem(item.id, status);
			items = items.map((currentItem) =>
				currentItem.id === item.id ? { ...currentItem, consumption: status } : currentItem
			);
			toast.success('Stanje stavke je promenjeno.');
		} catch (reason) {
			toast.error(reason instanceof ApiError ? reason.message : 'Stanje stavke nije promenjeno.');
		}
	}

	function goToPage(page: number) {
		updateTableQuery({ page });
	}

	function searchItems(search: string) {
		updateTableQuery({ search, page: 1 });
	}

	function filterByItemType(itemTypeID: number | undefined) {
		const propertyFilterUpdates: Record<string, string | null> = {};
		for (const key of page.url.searchParams.keys()) {
			if (key.startsWith('property.')) propertyFilterUpdates[key] = null;
		}
		updateTableQuery({ type_id: itemTypeID, page: 1, ...propertyFilterUpdates });
	}

	function filterByProperty(propertyID: number, value: PropertyValue | undefined) {
		updateTableQuery({ [`property.${propertyID}`]: JSON.stringify(value), page: 1 });
	}

	async function loadPropertyFilterValues(propertyID: number, search: string) {
		const itemTypeID = selectedItemTypeID;
		if (itemTypeID === undefined) return;

		const key = `${itemTypeID}:${propertyID}`;
		const version = (propertyFilterValueLoadVersions.get(key) ?? 0) + 1;
		propertyFilterValueLoadVersions.set(key, version);

		try {
			const values = await api.searchItemTypePropertyValues(itemTypeID, propertyID, search);
			if (selectedItemTypeID !== itemTypeID || propertyFilterValueLoadVersions.get(key) !== version)
				return;
			propertyFilterOptions = { ...propertyFilterOptions, [propertyID]: values };
		} catch (reason) {
			if (selectedItemTypeID !== itemTypeID || propertyFilterValueLoadVersions.get(key) !== version)
				return;
			reportInventoryError(
				reason instanceof ApiError ? reason.message : 'Vrednosti filtera nisu učitane.'
			);
		}
	}

	function reportInventoryError(message: string) {
		inventoryError = message;
	}

	function refreshInventory() {
		const itemTypeID = getSelectedItemTypeID(page.url);
		void loadInventory(itemOffset, derivedNameSearch, itemTypeID, getPropertyFilters(page.url));
	}

	function getSelectedItemTypeID(url: URL) {
		const requestedTypeID = Number.parseInt(getTableFilter(url, 'type_id'), 10);
		return Number.isSafeInteger(requestedTypeID) && requestedTypeID > 0
			? requestedTypeID
			: undefined;
	}

	function getPropertyFilters(url: URL): Record<number, PropertyValue> {
		const filters: Record<number, PropertyValue> = {};
		for (const [key, value] of url.searchParams) {
			if (!key.startsWith('property.')) continue;
			const propertyID = Number.parseInt(key.slice('property.'.length), 10);
			if (Number.isSafeInteger(propertyID) && propertyID > 0 && value)
				filters[propertyID] = JSON.parse(value);
		}
		return filters;
	}

	async function requestItemUsage(itemID: number, reason: string) {
		await api.createPersonalItemRequest({ item_id: itemID, reason });
		refreshInventory();
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
		error={authPage.state.error || inventoryError}
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

		<InventoryToolbar
			total={itemsTotal}
			propertyTotals={itemPropertyTotals}
			{properties}
			{itemTypes}
			typeFilter={itemTypeFilter}
			bind:search={derivedNameSearch}
			loading={loadingInventory}
			onitemtypechange={filterByItemType}
			onsearch={searchItems}
		/>
		{#if selectedItemType && selectedItemType.id === selectedItemTypeID}
			<InventoryPropertyFilters
				filters={propertyFilters}
				filterOptions={propertyFilterOptions}
				onfilterchange={filterByProperty}
				onfiltervaluesearch={loadPropertyFilterValues}
			/>
		{/if}
		<InventoryList
			{items}
			{itemTypes}
			{properties}
			{canManage}
			onconsumptionchange={changeConsumption}
			onrequest={requestItemUsage}
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
