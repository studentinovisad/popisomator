<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import Plus from '@lucide/svelte/icons/plus';
	import { onMount } from 'svelte';
	import { api, ApiError, type ItemType, type PropertyOption } from '$lib/api';
	import { createAuthPage } from '$lib/state/auth-page.svelte';
	import PaginationFooter from '$lib/components/shared/PaginationFooter.svelte';
	import ItemTypesList from '$lib/components/catalog/ItemTypesList.svelte';
	import ProtectedPageState from '$lib/components/shared/ProtectedPageState.svelte';
	import TableSearch from '$lib/components/shared/TableSearch.svelte';
	import { createServerPagination } from '$lib/state/server-pagination.svelte';
	import { getTablePage, getTableSearch, updateTableQuery } from '$lib/state/table-query';
	import { Portal } from 'bits-ui';
	import { toast } from 'svelte-sonner';

	const authPage = createAuthPage({
		unavailableMessage: 'Tipovi stavki trenutno nisu dostupni.',
		requiredRoles: ['admin']
	});

	let propertyOptionsError = $state('');
	const itemTypesPage = createServerPagination<ItemType>({
		loadPage: api.listItemTypes,
		unavailableMessage: 'Tipovi stavki nisu učitani.'
	});
	let propertyOptions = $state<PropertyOption[]>([]);
	let propertyOptionsLoading = $state(true);
	let search = $state('');

	onMount(() => {
		void authPage.load().then(() => {
			if (authPage.state.authorized) {
				void loadPropertyOptions();
			}
		});
	});

	$effect(() => {
		if (!authPage.state.authorized) return;

		const url = page.url;
		const nextSearch = getTableSearch(url);
		search = nextSearch;
		itemTypesPage.sync({ page: getTablePage(url), search: nextSearch });
	});

	async function loadPropertyOptions() {
		propertyOptionsLoading = true;
		try {
			propertyOptions = await api.getPropertyOptions();
		} catch (reason) {
			propertyOptionsError = reason instanceof ApiError ? reason.message : 'Svojstva nisu učitana.';
		} finally {
			propertyOptionsLoading = false;
		}
	}

	async function deleteItemType(itemType: ItemType) {
		if (!confirm(`Obrisati tip ${itemType.name}?`)) return;

		try {
			await api.deleteItemType(itemType.id);
			toast.success('Tip stavke je obrisan.');
			itemTypesPage.reloadAfterDelete();
		} catch (reason) {
			toast.error(reason instanceof ApiError ? reason.message : 'Tip stavke nije obrisan.');
		}
	}

	function searchItemTypes(nextSearch: string) {
		updateTableQuery({ search: nextSearch, page: 1 });
	}

	function goToPage(nextPage: number) {
		updateTableQuery({ page: nextPage });
	}
</script>

<svelte:head>
	<title>Tipovi stavki | Popisomator</title>
</svelte:head>

<main class="px-4 pt-4 pb-8 sm:px-6">
	<ProtectedPageState
		loading={authPage.state.loading ||
			(authPage.state.authorized && (itemTypesPage.loading || propertyOptionsLoading))}
		contentLoaded={itemTypesPage.loaded && !propertyOptionsLoading}
		error={authPage.state.error || itemTypesPage.error || propertyOptionsError}
		authorized={authPage.state.authorized}
	>
		<p class="font-mono text-xs leading-none font-medium tracking-wide text-muted">
			UKUPNO: {itemTypesPage.total}
		</p>

		<Portal to="#page-header-actions">
			<a
				class="inline-flex size-10 items-center justify-center rounded-md bg-brand text-on-brand hover:bg-brand-strong"
				href={resolve('/catalog/item-types/new')}
				aria-label="Dodaj tip"
				title="Dodaj tip"
			>
				<Plus class="size-4" aria-hidden="true" />
			</a>
			<a
				class="inline-flex h-10 items-center justify-center rounded-md border border-line bg-surface px-4 text-sm font-medium text-ink hover:bg-soft"
				href={resolve('/catalog/properties')}
			>
				Svojstva
			</a>
		</Portal>

		<TableSearch
			id="item-type-name-search"
			placeholder="Pretraži po nazivu"
			bind:search
			loading={itemTypesPage.loading}
			onsearch={searchItemTypes}
		/>
		<ItemTypesList
			itemTypes={itemTypesPage.items}
			deleteitemtype={deleteItemType}
			{propertyOptions}
		/>
		<PaginationFooter
			total={itemTypesPage.total}
			perPage={itemTypesPage.perPage}
			page={itemTypesPage.currentPage}
			hasPreviousPage={itemTypesPage.hasPreviousPage}
			hasNextPage={itemTypesPage.hasNextPage}
			loading={itemTypesPage.loading}
			onpagechange={goToPage}
		/>
	</ProtectedPageState>
</main>
