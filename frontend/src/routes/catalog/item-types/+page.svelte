<script lang="ts">
	import { resolve } from '$app/paths';
	import { onMount } from 'svelte';
	import { api, ApiError, type ItemType } from '$lib/api';
	import { createAuthPage } from '$lib/state/auth-page.svelte';
	import PaginationFooter from '$lib/components/shared/PaginationFooter.svelte';
	import ItemTypesList from '$lib/components/catalog/ItemTypesList.svelte';
	import ProtectedPageState from '$lib/components/shared/ProtectedPageState.svelte';
	import { createServerPagination } from '$lib/state/server-pagination.svelte';
	import { Portal } from 'bits-ui';

	const authPage = createAuthPage({
		unavailableMessage: 'Tipovi stavki trenutno nisu dostupni.',
		requiredRoles: ['admin']
	});

	let error = $state('');
	const itemTypesPage = createServerPagination<ItemType>({
		loadPage: api.listItemTypesPage,
		unavailableMessage: 'Tipovi stavki nisu učitani.'
	});

	onMount(() => {
		void authPage.load().then(() => {
			if (authPage.state.authorized) void itemTypesPage.load();
		});
	});

	async function deleteItemType(itemType: ItemType) {
		if (!confirm(`Obrisati tip ${itemType.name}?`)) return;
		error = '';

		try {
			await api.deleteItemType(itemType.id);
			itemTypesPage.reloadAfterDelete();
		} catch (reason) {
			error = reason instanceof ApiError ? reason.message : 'Tip stavke nije obrisan.';
		}
	}
</script>

<svelte:head>
	<title>Tipovi stavki | Popisomator</title>
</svelte:head>

<main class="px-4 pt-4 pb-8 sm:px-6">
	<ProtectedPageState
		loading={authPage.state.loading || (authPage.state.authorized && itemTypesPage.loading)}
		error={authPage.state.error || itemTypesPage.error}
		authorized={authPage.state.authorized}
	>
		<p class="font-mono text-xs leading-none font-medium tracking-wide text-muted">
			UKUPNO: {itemTypesPage.total}
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
		<ItemTypesList itemTypes={itemTypesPage.items} deleteitemtype={deleteItemType} />
		<PaginationFooter
			total={itemTypesPage.total}
			perPage={itemTypesPage.perPage}
			page={itemTypesPage.currentPage}
			hasPreviousPage={itemTypesPage.hasPreviousPage}
			hasNextPage={itemTypesPage.hasNextPage}
			loading={itemTypesPage.loading}
			onpagechange={itemTypesPage.goToPage}
		/>
	</ProtectedPageState>
</main>
