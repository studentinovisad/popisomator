<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { onMount } from 'svelte';
	import { api, ApiError, type Property } from '$lib/api';
	import { createAuthPage } from '$lib/auth-page.svelte';
	import CreateItemTypeForm from '$lib/components/CreateItemTypeForm.svelte';
	import PageLoader from '$lib/components/PageLoader.svelte';

	const authPage = createAuthPage({
		unavailableMessage: 'Dodavanje tipa stavke trenutno nije dostupno.',
		requiredRoles: ['admin']
	});

	let properties = $state<Property[]>([]);
	let loading = $state(false);
	let error = $state('');

	onMount(() => {
		void authPage.load().then(() => {
			if (authPage.state.authorized) void loadProperties();
		});
	});

	async function loadProperties() {
		loading = true;
		error = '';

		try {
			properties = await api.listProperties();
		} catch (reason) {
			error = reason instanceof ApiError ? reason.message : 'Svojstva nisu učitana.';
		} finally {
			loading = false;
		}
	}

	function itemTypeCreated() {
		void goto(resolve('/catalog/item-types'));
	}

	function cancelItemTypeCreation() {
		void goto(resolve('/catalog/item-types'));
	}
</script>

<svelte:head>
	<title>Novi tip stavke | Popisomator</title>
</svelte:head>

<main class="px-4 pt-4 pb-8 sm:px-6">
	{#if authPage.state.loading || (authPage.state.authorized && loading)}
		<PageLoader />
	{:else if authPage.state.error || error}
		<div class="grid min-h-[calc(100svh-14rem)] place-items-center">
			<p class="text-danger" role="alert">{authPage.state.error || error}</p>
		</div>
	{:else if authPage.state.authorized}
		<section class="mx-auto max-w-3xl" aria-labelledby="new-item-type-heading">
			<div class="border-b border-line pb-4">
				<h2 id="new-item-type-heading" class="text-lg font-semibold text-ink">Novi tip stavke</h2>
				<p class="mt-1 text-sm text-muted">Odaberite svojstva koja pripadaju ovom tipu.</p>
			</div>
			<div class="mt-6">
				<CreateItemTypeForm
					{properties}
					onsaved={itemTypeCreated}
					oncancel={cancelItemTypeCreation}
				/>
			</div>
		</section>
	{/if}
</main>
