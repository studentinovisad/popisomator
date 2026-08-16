<script lang="ts">
	import { resolve } from '$app/paths';
	import ArrowLeft from '@lucide/svelte/icons/arrow-left';
	import { onMount } from 'svelte';
	import { api, ApiError, type Property } from '$lib/api';
	import { createAuthPage } from '$lib/auth-page.svelte';
	import PaginationFooter from '$lib/components/PaginationFooter.svelte';
	import PropertiesList from '$lib/components/PropertiesList.svelte';
	import ProtectedPageState from '$lib/components/ProtectedPageState.svelte';
	import { createServerPagination } from '$lib/server-pagination.svelte';
	import { Portal } from 'bits-ui';

	const authPage = createAuthPage({
		unavailableMessage: 'Svojstva trenutno nisu dostupna.',
		requiredRoles: ['admin']
	});

	let error = $state('');
	const propertiesPage = createServerPagination<Property>({
		loadPage: api.listPropertiesPage,
		unavailableMessage: 'Svojstva nisu učitana.'
	});

	onMount(() => {
		void authPage.load().then(() => {
			if (authPage.state.authorized) void propertiesPage.load();
		});
	});

	async function deleteProperty(property: Property) {
		if (!confirm(`Obrisati svojstvo ${property.name}?`)) return;
		error = '';

		try {
			await api.deleteProperty(property.id);
			propertiesPage.reloadAfterDelete();
		} catch (reason) {
			error = reason instanceof ApiError ? reason.message : 'Svojstvo nije obrisano.';
		}
	}
</script>

<svelte:head>
	<title>Svojstva | Popisomator</title>
</svelte:head>

<main class="px-4 pt-4 pb-8 sm:px-6">
	<ProtectedPageState
		loading={authPage.state.loading || (authPage.state.authorized && propertiesPage.loading)}
		error={authPage.state.error || propertiesPage.error}
		authorized={authPage.state.authorized}
	>
		<p class="font-mono text-xs leading-none font-medium tracking-wide text-muted">
			UKUPNO: {propertiesPage.total}
		</p>
		<Portal to="#page-header-actions">
			<a
				class="inline-flex h-10 items-center justify-center rounded-md bg-brand px-4 text-sm font-medium text-on-brand hover:bg-brand-strong"
				href={resolve('/catalog/properties/new')}
			>
				Dodaj svojstvo
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

		{#if error}<p class="mt-3 text-sm text-danger" role="alert">{error}</p>{/if}
		<PropertiesList properties={propertiesPage.items} deleteproperty={deleteProperty} />
		<PaginationFooter
			total={propertiesPage.total}
			perPage={propertiesPage.perPage}
			page={propertiesPage.currentPage}
			hasPreviousPage={propertiesPage.hasPreviousPage}
			hasNextPage={propertiesPage.hasNextPage}
			loading={propertiesPage.loading}
			onpagechange={propertiesPage.goToPage}
		/>
	</ProtectedPageState>
</main>
