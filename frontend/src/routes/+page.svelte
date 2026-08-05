<script lang="ts">
	import { resolve } from '$app/paths';
	import { onMount } from 'svelte';
	import {
		api,
		ApiError,
		type ConsumptionStatus,
		type Item,
		type ItemType,
		type Property
	} from '$lib/api';
	import { createAuthPage } from '$lib/auth-page.svelte';
	import ItemPropertiesForm from '$lib/components/ItemPropertiesForm.svelte';
	import InventoryList from '$lib/components/InventoryList.svelte';
	import PageLoader from '$lib/components/PageLoader.svelte';
	import { Dialog, Portal } from 'bits-ui';

	const authPage = createAuthPage({ unavailableMessage: 'Inventar trenutno nije dostupan.' });

	let items = $state<Item[]>([]);
	let itemTypes = $state<ItemType[]>([]);
	let properties = $state<Property[]>([]);
	let loadingInventory = $state(false);
	let inventoryError = $state('');
	let editItemDialogOpen = $state(false);
	let editingItem = $state<Item | null>(null);
	let loadVersion = 0;
	let canManage = $derived(
		authPage.state.user?.role === 'admin' || authPage.state.user?.role === 'manager'
	);

	onMount(() => {
		void authPage.load().then(() => {
			if (authPage.state.authorized) void loadInventory();
		});
	});

	async function loadInventory() {
		const version = ++loadVersion;
		loadingInventory = true;
		inventoryError = '';

		try {
			const [nextItems, nextItemTypes, nextProperties] = await Promise.all([
				api.listItems(),
				api.listItemTypes(),
				api.listProperties()
			]);
			if (version !== loadVersion) return;

			items = nextItems;
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

	async function deleteItem(item: Item) {
		if (!confirm(`Obrisati stavku #${item.id}?`)) return;
		inventoryError = '';

		try {
			await api.deleteItem(item.id);
			items = items.filter((currentItem) => currentItem.id !== item.id);
		} catch (reason) {
			inventoryError = reason instanceof ApiError ? reason.message : 'Stavka nije obrisana.';
		}
	}

	function editItem(item: Item) {
		editingItem = item;
		editItemDialogOpen = true;
	}

	function itemPropertiesSaved() {
		editItemDialogOpen = false;
		void loadInventory();
	}
</script>

<svelte:head>
	<title>Stavke | Popisomator</title>
</svelte:head>

<main class="px-4 pt-4 pb-8 sm:px-6">
	{#if authPage.state.loading || (authPage.state.authorized && loadingInventory)}
		<PageLoader />
	{:else if authPage.state.error}
		<div class="grid min-h-[calc(100svh-14rem)] place-items-center">
			<p class="text-danger" role="alert">{authPage.state.error}</p>
		</div>
	{:else if authPage.state.authorized && authPage.state.user}
		<div class="flex items-center justify-between gap-4">
			<p class="font-mono text-xs leading-none font-medium tracking-wide text-muted">
				UKUPNO: {items.length}
			</p>
		</div>

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

		<Dialog.Root bind:open={editItemDialogOpen}>
			<Dialog.Portal>
				<Dialog.Overlay class="fixed inset-0 z-20 bg-ink/35 backdrop-blur-sm" />
				<Dialog.Content
					class="fixed top-1/2 left-1/2 z-30 max-h-[calc(100svh-2rem)] w-[calc(100%-2rem)] max-w-xl -translate-x-1/2 -translate-y-1/2 overflow-y-auto rounded-lg border border-line bg-surface p-6 shadow-xl shadow-ink/20"
				>
					<div class="flex items-start justify-between gap-4">
						<div>
							<Dialog.Title class="text-xl font-semibold text-ink">Svojstva stavke</Dialog.Title>
							<Dialog.Description class="mt-1 text-sm text-muted"
								>Izmenite vrednosti za stavku #{editingItem?.id}.</Dialog.Description
							>
						</div>
						<Dialog.Close
							class="rounded-md px-2 py-1 text-sm text-muted hover:bg-soft hover:text-ink"
							>Zatvori</Dialog.Close
						>
					</div>
					{#if editingItem}
						<div class="mt-6">
							<ItemPropertiesForm
								item={editingItem}
								{itemTypes}
								{properties}
								onsaved={itemPropertiesSaved}
							/>
						</div>
					{/if}
				</Dialog.Content>
			</Dialog.Portal>
		</Dialog.Root>

		{#if inventoryError}<p class="mt-3 text-sm text-danger" role="alert">{inventoryError}</p>{/if}
		<InventoryList
			{items}
			{itemTypes}
			{properties}
			user={authPage.state.user}
			onconsumptionchange={changeConsumption}
			onedititem={editItem}
			deleteitem={deleteItem}
		/>
	{/if}
</main>
