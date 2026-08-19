<script lang="ts">
	import { resolve } from '$app/paths';
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
	import ProtectedPageState from '$lib/components/shared/ProtectedPageState.svelte';
	import { pagination } from '$lib/state/pagination.svelte';
	import { Portal } from 'bits-ui';

	const authPage = createAuthPage({ unavailableMessage: 'Inventar trenutno nije dostupan.' });

	let items = $state<Item[]>([]);
	let itemTypes = $state<ItemTypeOption[]>([]);
	let properties = $state<PropertyOption[]>([]);
	let loadingInventory = $state(false);
	let inventoryError = $state('');
	let loadVersion = 0;
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
			if (authPage.state.authorized) void loadInventory();
		});
	});

	async function loadInventory(offset = 0) {
		const version = ++loadVersion;
		loadingInventory = true;
		inventoryError = '';

		try {
			const [nextItems, nextItemTypes, nextProperties] = await Promise.all([
				api.listItems({ limit: itemsPerPage, offset }),
				api.getItemTypeOptions(),
				api.getPropertyOptions()
			]);
			if (version !== loadVersion) return;

			items = nextItems.items;
			itemOffset = nextItems.offset;
			itemsTotal = nextItems.total;
			itemTypes = nextItemTypes;
			properties = nextProperties;
		} catch (reason) {
			if (version !== loadVersion) return;
			inventoryError = reason instanceof ApiError ? reason.message : 'Stavke nisu učitane.';
		} finally {
			if (version === loadVersion) loadingInventory = false;
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
		void loadInventory((page - 1) * itemsPerPage);
	}
</script>

<svelte:head>
	<title>Stavke | Popisomator</title>
</svelte:head>

<main class="px-4 pt-4 sm:px-6">
	<ProtectedPageState
		loading={authPage.state.loading || (authPage.state.authorized && loadingInventory)}
		error={authPage.state.error}
		authorized={authPage.state.authorized && authPage.state.user !== null}
	>
		<p class="font-mono text-xs leading-none font-medium tracking-wide text-muted">
			UKUPNO: {itemsTotal}
		</p>

		{#if canManage}
			<Portal to="#page-header-actions">
				<a
					class="inline-flex h-10 items-center justify-center rounded-md bg-brand px-4 text-sm font-medium text-on-brand hover:bg-brand-strong"
					href={resolve('/items/new')}
				>
					Dodaj stavku
				</a>
			</Portal>
		{/if}

		{#if inventoryError}<p class="mt-3 text-sm text-danger" role="alert">{inventoryError}</p>{/if}
		<InventoryList {items} {itemTypes} {properties} onconsumptionchange={changeConsumption} />
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
