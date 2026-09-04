<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import ArrowLeft from '@lucide/svelte/icons/arrow-left';
	import Plus from '@lucide/svelte/icons/plus';
	import { onMount } from 'svelte';
	import { api, ApiError, type Property } from '$lib/api';
	import { createAuthPage } from '$lib/state/auth-page.svelte';
	import PaginationFooter from '$lib/components/shared/PaginationFooter.svelte';
	import PropertiesList from '$lib/components/catalog/PropertiesList.svelte';
	import ProtectedPageState from '$lib/components/shared/ProtectedPageState.svelte';
	import TableSearch from '$lib/components/shared/TableSearch.svelte';
	import { createServerPagination } from '$lib/state/server-pagination.svelte';
	import { getTablePage, getTableSearch, updateTableQuery } from '$lib/state/table-query';
	import { Portal } from 'bits-ui';
	import { toast } from 'svelte-sonner';

	const authPage = createAuthPage({
		unavailableMessage: 'Svojstva trenutno nisu dostupna.',
		requiredRoles: ['admin']
	});

	const propertiesPage = createServerPagination<Property>({
		loadPage: api.listProperties,
		unavailableMessage: 'Svojstva nisu učitana.'
	});
	let search = $state('');

	onMount(() => {
		void authPage.load();
	});

	$effect(() => {
		if (!authPage.state.authorized) return;

		const url = page.url;
		const nextSearch = getTableSearch(url);
		search = nextSearch;
		propertiesPage.sync({ page: getTablePage(url), search: nextSearch });
	});

	async function deleteProperty(property: Property) {
		if (!confirm(`Obrisati svojstvo ${property.name}?`)) return;
		try {
			await api.deleteProperty(property.id);
			toast.success('Svojstvo je obrisano.');
			propertiesPage.reloadAfterDelete();
		} catch (reason) {
			toast.error(reason instanceof ApiError ? reason.message : 'Svojstvo nije obrisano.');
		}
	}

	function searchProperties(nextSearch: string) {
		updateTableQuery({ search: nextSearch, page: 1 });
	}

	function goToPage(nextPage: number) {
		updateTableQuery({ page: nextPage });
	}
</script>

<svelte:head>
	<title>Svojstva | Popisomator</title>
</svelte:head>

<main class="px-4 pt-4 pb-8 sm:px-6">
	<ProtectedPageState
		loading={authPage.state.loading || (authPage.state.authorized && propertiesPage.loading)}
		contentLoaded={propertiesPage.loaded}
		error={authPage.state.error || propertiesPage.error}
		authorized={authPage.state.authorized}
	>
		<p class="font-mono text-xs leading-none font-medium tracking-wide text-muted">
			UKUPNO: {propertiesPage.total}
		</p>
		<Portal to="#page-header-actions">
			<a
				class="inline-flex size-10 items-center justify-center rounded-md bg-brand text-on-brand hover:bg-brand-strong"
				href={resolve('/catalog/properties/new')}
				aria-label="Dodaj svojstvo"
				title="Dodaj svojstvo"
			>
				<Plus class="size-4" aria-hidden="true" />
			</a>
			<a
				class="inline-flex size-10 items-center justify-center rounded-md border border-line bg-surface text-ink transition-colors hover:border-brand/40 hover:bg-brand-soft hover:text-brand"
				href={resolve('/catalog/item-types')}
				aria-label="Nazad na tipove stavki"
				title="Nazad na tipove stavki"
			>
				<ArrowLeft class="size-4" aria-hidden="true" />
			</a>
		</Portal>

		<TableSearch
			id="property-name-search"
			placeholder="Pretraži po nazivu"
			bind:search
			loading={propertiesPage.loading}
			onsearch={searchProperties}
		/>
		<PropertiesList properties={propertiesPage.items} deleteproperty={deleteProperty} />
		<PaginationFooter
			total={propertiesPage.total}
			perPage={propertiesPage.perPage}
			page={propertiesPage.currentPage}
			hasPreviousPage={propertiesPage.hasPreviousPage}
			hasNextPage={propertiesPage.hasNextPage}
			loading={propertiesPage.loading}
			onpagechange={goToPage}
		/>
	</ProtectedPageState>
</main>
