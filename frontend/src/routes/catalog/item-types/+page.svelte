<script lang="ts">
	import { resolve } from '$app/paths';
	import { onMount } from 'svelte';
	import { api, ApiError, type ItemType, type Property } from '$lib/api';
	import { createAuthPage } from '$lib/auth-page.svelte';
	import ItemTypesList from '$lib/components/ItemTypesList.svelte';
	import PageLoader from '$lib/components/PageLoader.svelte';
	import { Portal } from 'bits-ui';

	const authPage = createAuthPage({
		unavailableMessage: 'Tipovi stavki trenutno nisu dostupni.',
		requiredRoles: ['admin']
	});

	let itemTypes = $state<ItemType[]>([]);
	let properties = $state<Property[]>([]);
	let loading = $state(false);
	let error = $state('');
	let loadVersion = 0;
	let propertyNames = $derived(new Map(properties.map((property) => [property.id, property.name])));

	onMount(() => {
		void authPage.load().then(() => {
			if (authPage.state.authorized) void loadItemTypes();
		});
	});

	async function loadItemTypes() {
		const version = ++loadVersion;
		loading = true;
		error = '';

		try {
			const [nextItemTypes, nextProperties] = await Promise.all([
				api.listItemTypes(),
				api.listProperties()
			]);
			if (version !== loadVersion) return;
			itemTypes = nextItemTypes;
			properties = nextProperties;
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
			itemTypes = itemTypes.filter((currentItemType) => currentItemType.id !== itemType.id);
		} catch (reason) {
			error = reason instanceof ApiError ? reason.message : 'Tip stavke nije obrisan.';
		}
	}
</script>

<svelte:head>
	<title>Tipovi stavki | Popisomator</title>
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
			UKUPNO: {itemTypes.length}
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
		<ItemTypesList {itemTypes} {propertyNames} deleteitemtype={deleteItemType} />
	{/if}
</main>
