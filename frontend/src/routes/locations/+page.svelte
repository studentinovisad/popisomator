<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import Plus from '@lucide/svelte/icons/plus';
	import { onMount } from 'svelte';
	import { api, ApiError, type ItemType, type LocationOption, type PropertyOption } from '$lib/api';
	import { createAuthPage } from '$lib/state/auth-page.svelte';
	import PaginationFooter from '$lib/components/shared/PaginationFooter.svelte';
	import ItemTypesList from '$lib/components/catalog/ItemTypesList.svelte';
	import ProtectedPageState from '$lib/components/shared/ProtectedPageState.svelte';
	import TableSearch from '$lib/components/shared/TableSearch.svelte';
	import { createServerPagination } from '$lib/state/server-pagination.svelte';
	import { getTablePage, getTableSearch, updateTableQuery } from '$lib/state/table-query';
	import { Portal } from 'bits-ui';
	import { toast } from 'svelte-sonner';
	import LocationBlock from '$lib/components/catalog/LocationBlock.svelte';

	const authPage = createAuthPage({
		unavailableMessage: 'Lokacije trenutno nisu dostupne.',
		requiredRoles: ['admin']
	});

	let locationOptionsError = $state('');
	let locationOptions = $state<LocationOption[]>([]);
	let locationOptionsLoading = $state(true);
	let search = $state('');

	onMount(() => {
		void authPage.load().then(() => {
			if (authPage.state.authorized) {
				void loadLocationOptions();
			}
		});
	});

	$effect(() => {
		if (!authPage.state.authorized) return;

		const url = page.url;
		const nextSearch = getTableSearch(url);
		search = nextSearch;
		//itemTypesPage.sync({ page: getTablePage(url), search: nextSearch });
	});

	async function loadLocationOptions() {
		locationOptionsLoading = true;
		try {
			locationOptions = await api.getLocationOptions();
		} catch (reason) {
			locationOptionsError = reason instanceof ApiError ? reason.message : 'Lokacije nisu učitane.';
		} finally {
			locationOptionsLoading = false;
		}
	}

	async function deleteLocation(location: LocationOption) {
		if (!confirm(`Obrisati lokaciju ${location.name}?`)) return;

		try {
			await api.deleteLocation(location.id);
			toast.success('Lokacija je obrisana.');
			loadLocationOptions();
		} catch (reason) {
			toast.error(reason instanceof ApiError ? reason.message : 'Lokacija nije obrisana.');
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
	<title>Lokacije | Popisomator</title>
</svelte:head>

<main class="px-4 pt-4 pb-8 sm:px-6">
	<ProtectedPageState
		loading={authPage.state.loading ||
			(authPage.state.authorized && locationOptionsLoading)}
		contentLoaded={!locationOptionsLoading}
		error={authPage.state.error || locationOptionsError}
		authorized={authPage.state.authorized}
	>

		<Portal to="#page-header-actions">
			<a
				class="inline-flex size-10 items-center justify-center rounded-md bg-brand text-on-brand hover:bg-brand-strong"
				href={resolve('/locations/new')}
				aria-label="Dodaj lokaciju"
				title="Dodaj lokaciju"
			>
				<Plus class="size-4" aria-hidden="true" />
			</a>
		</Portal>

		{#if locationOptions.length > 0}
			<div class="flex flex-wrap w-full h-full gap-8 justify-center content-center items-center ">
				{#each locationOptions as location}
					<LocationBlock location={location} alternateBg={true} deletelocation={deleteLocation}/>
				{/each}
			</div>
		{:else}
			<p class="px-4 py-3 text-sm text-muted">Nema lokacija.</p>
		{/if}
	</ProtectedPageState>
</main>
