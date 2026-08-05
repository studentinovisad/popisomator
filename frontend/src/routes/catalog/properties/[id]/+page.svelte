<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { onMount } from 'svelte';
	import { api, ApiError, type Property } from '$lib/api';
	import { createAuthPage } from '$lib/auth-page.svelte';
	import PageLoader from '$lib/components/PageLoader.svelte';
	import PropertyForm from '$lib/components/PropertyForm.svelte';

	const authPage = createAuthPage({
		unavailableMessage: 'Izmena svojstva trenutno nije dostupna.',
		requiredRoles: ['admin']
	});

	let property = $state<Property | null>(null);
	let loadingProperty = $state(false);
	let error = $state('');

	onMount(() => {
		void authPage.load().then(() => {
			if (authPage.state.authorized) void loadProperty();
		});
	});

	async function loadProperty() {
		const id = Number(page.params.id);
		if (!Number.isSafeInteger(id) || id < 1) {
			error = 'Svojstvo nije pronađeno.';
			return;
		}

		loadingProperty = true;
		error = '';
		try {
			property = await api.getProperty(id);
		} catch (reason) {
			error =
				reason instanceof ApiError && reason.status === 404
					? 'Svojstvo nije pronađeno.'
					: 'Svojstvo nije učitano.';
		} finally {
			loadingProperty = false;
		}
	}

	function propertySaved() {
		void goto(resolve('/catalog/properties'));
	}
</script>

<svelte:head>
	<title>Izmeni svojstvo | Popisomator</title>
</svelte:head>

<main class="px-4 pt-4 pb-8 sm:px-6">
	{#if authPage.state.loading || (authPage.state.authorized && loadingProperty)}
		<PageLoader />
	{:else if authPage.state.error || error}
		<div class="grid min-h-[calc(100svh-14rem)] place-items-center">
			<p class="text-danger" role="alert">{authPage.state.error || error}</p>
		</div>
	{:else if authPage.state.authorized && property}
		<section class="mx-auto max-w-3xl" aria-labelledby="edit-property-heading">
			<div class="border-b border-line pb-4">
				<h2 id="edit-property-heading" class="text-lg font-semibold text-ink">Izmeni svojstvo</h2>
				<p class="mt-1 text-sm text-muted">{property.name}</p>
			</div>
			<div class="mt-6"><PropertyForm {property} onsaved={propertySaved} /></div>
		</section>
	{/if}
</main>
