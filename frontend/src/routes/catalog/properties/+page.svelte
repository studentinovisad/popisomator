<script lang="ts">
	import { resolve } from '$app/paths';
	import ArrowLeft from '@lucide/svelte/icons/arrow-left';
	import { onMount } from 'svelte';
	import { api, ApiError, type Property } from '$lib/api';
	import { createAuthPage } from '$lib/auth-page.svelte';
	import InventoryPagination from '$lib/components/InventoryPagination.svelte';
	import PageLoader from '$lib/components/PageLoader.svelte';
	import PropertiesList from '$lib/components/PropertiesList.svelte';
	import { pagination } from '$lib/pagination.svelte';
	import { Portal } from 'bits-ui';

	const authPage = createAuthPage({
		unavailableMessage: 'Svojstva trenutno nisu dostupna.',
		requiredRoles: ['admin']
	});

	let properties = $state<Property[]>([]);
	let loading = $state(false);
	let error = $state('');
	let loadVersion = 0;
	let propertiesPerPage = $derived(pagination.perPage);
	let propertyOffset = $state(0);
	let propertiesTotal = $state(0);
	let currentPage = $derived(Math.floor(propertyOffset / propertiesPerPage) + 1);
	let hasPreviousPage = $derived(propertyOffset > 0);
	let hasNextPage = $derived(propertyOffset + properties.length < propertiesTotal);

	onMount(() => {
		void authPage.load().then(() => {
			if (authPage.state.authorized) void loadProperties();
		});
	});

	async function loadProperties(offset = 0) {
		const version = ++loadVersion;
		loading = true;
		error = '';

		try {
			const nextProperties = await api.listPropertiesPage({ limit: propertiesPerPage, offset });
			if (version !== loadVersion) return;
			properties = nextProperties.items;
			propertyOffset = nextProperties.offset;
			propertiesTotal = nextProperties.total;
		} catch (reason) {
			if (version !== loadVersion) return;
			error = reason instanceof ApiError ? reason.message : 'Svojstva nisu učitana.';
		} finally {
			if (version === loadVersion) loading = false;
		}
	}

	async function deleteProperty(property: Property) {
		if (!confirm(`Obrisati svojstvo ${property.name}?`)) return;
		error = '';

		try {
			await api.deleteProperty(property.id);
			const nextOffset =
				properties.length === 1 && propertyOffset > 0
					? propertyOffset - propertiesPerPage
					: propertyOffset;
			void loadProperties(nextOffset);
		} catch (reason) {
			error = reason instanceof ApiError ? reason.message : 'Svojstvo nije obrisano.';
		}
	}

	function goToPage(page: number) {
		void loadProperties((page - 1) * propertiesPerPage);
	}
</script>

<svelte:head>
	<title>Svojstva | Popisomator</title>
</svelte:head>

<main class="px-4 pt-4 pb-8 sm:px-6">
	{#if authPage.state.loading || (authPage.state.authorized && loading)}
		<PageLoader />
	{:else if authPage.state.error}
		<div class="grid min-h-[calc(100svh-14rem)] place-items-center">
			<p class="text-danger" role="alert">{authPage.state.error}</p>
		</div>
	{:else if authPage.state.authorized}
		<p class="font-mono text-xs leading-none font-medium tracking-wide text-muted">
			UKUPNO: {propertiesTotal}
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
		<PropertiesList {properties} deleteproperty={deleteProperty} />
		<InventoryPagination
			total={propertiesTotal}
			perPage={propertiesPerPage}
			page={currentPage}
			{hasPreviousPage}
			{hasNextPage}
			{loading}
			onpagechange={goToPage}
		/>
	{/if}
</main>
