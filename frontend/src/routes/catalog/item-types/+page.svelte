<script lang="ts">
	import { resolve } from '$app/paths';
	import { onMount } from 'svelte';
	import { api, ApiError, type ItemType } from '$lib/api';
	import { createAuthPage } from '$lib/auth-page.svelte';
	import PaginationFooter from '$lib/components/PaginationFooter.svelte';
	import ItemTypesList from '$lib/components/ItemTypesList.svelte';
	import ProtectedPageState from '$lib/components/ProtectedPageState.svelte';
	import { pagination } from '$lib/pagination.svelte';
	import { Portal } from 'bits-ui';

	const authPage = createAuthPage({
		unavailableMessage: 'Tipovi stavki trenutno nisu dostupni.',
		requiredRoles: ['admin']
	});

	let itemTypes = $state<ItemType[]>([]);
	let loading = $state(false);
	let error = $state('');
	let loadVersion = 0;
	let itemTypesPerPage = $derived(pagination.perPage);
	let itemTypeOffset = $state(0);
	let itemTypesTotal = $state(0);
	let currentPage = $derived(Math.floor(itemTypeOffset / itemTypesPerPage) + 1);
	let hasPreviousPage = $derived(itemTypeOffset > 0);
	let hasNextPage = $derived(itemTypeOffset + itemTypes.length < itemTypesTotal);

	onMount(() => {
		void authPage.load().then(() => {
			if (authPage.state.authorized) void loadItemTypes();
		});
	});

	async function loadItemTypes(offset = 0) {
		const version = ++loadVersion;
		loading = true;
		error = '';

		try {
			const nextItemTypes = await api.listItemTypesPage({ limit: itemTypesPerPage, offset });
			if (version !== loadVersion) return;
			itemTypes = nextItemTypes.items;
			itemTypeOffset = nextItemTypes.offset;
			itemTypesTotal = nextItemTypes.total;
		} catch (reason) {
			if (version !== loadVersion) return;
			error = reason instanceof ApiError ? reason.message : 'Tipovi stavki nisu učitani.';
		} finally {
			if (version === loadVersion) loading = false;
		}
	}

	async function deleteItemType(itemType: ItemType) {
		if (!confirm(`Obrisati tip ${itemType.name}?`)) return;
		error = '';

		try {
			await api.deleteItemType(itemType.id);
			const nextOffset =
				itemTypes.length === 1 && itemTypeOffset > 0
					? itemTypeOffset - itemTypesPerPage
					: itemTypeOffset;
			void loadItemTypes(nextOffset);
		} catch (reason) {
			error = reason instanceof ApiError ? reason.message : 'Tip stavke nije obrisan.';
		}
	}

	function goToPage(page: number) {
		void loadItemTypes((page - 1) * itemTypesPerPage);
	}
</script>

<svelte:head>
	<title>Tipovi stavki | Popisomator</title>
</svelte:head>

<main class="px-4 pt-4 pb-8 sm:px-6">
	<ProtectedPageState
		loading={authPage.state.loading || (authPage.state.authorized && loading)}
		error={authPage.state.error}
		authorized={authPage.state.authorized}
	>
		<p class="font-mono text-xs leading-none font-medium tracking-wide text-muted">
			UKUPNO: {itemTypesTotal}
		</p>

		<Portal to="#page-header-actions">
			<a
				class="inline-flex h-10 items-center justify-center rounded-md bg-brand px-4 text-sm font-medium text-on-brand hover:bg-brand-strong"
				href={resolve('/catalog/item-types/new')}
			>
				Dodaj tip
			</a>
			<a
				class="inline-flex h-10 items-center justify-center rounded-md border border-line bg-surface px-4 text-sm font-medium text-ink hover:bg-soft"
				href={resolve('/catalog/properties')}
			>
				Svojstva
			</a>
		</Portal>

		{#if error}<p class="mt-3 text-sm text-danger" role="alert">{error}</p>{/if}
		<ItemTypesList {itemTypes} deleteitemtype={deleteItemType} />
		<PaginationFooter
			total={itemTypesTotal}
			perPage={itemTypesPerPage}
			page={currentPage}
			{hasPreviousPage}
			{hasNextPage}
			{loading}
			onpagechange={goToPage}
		/>
	</ProtectedPageState>
</main>
