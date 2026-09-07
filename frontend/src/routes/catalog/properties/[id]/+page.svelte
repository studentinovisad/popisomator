<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { onMount } from 'svelte';
	import { api, ApiError, type Property } from '$lib/api';
	import { createAuthPage } from '$lib/state/auth-page.svelte';
	import PropertyForm from '$lib/components/catalog/PropertyForm.svelte';
	import ProtectedPageState from '$lib/components/shared/ProtectedPageState.svelte';

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

	function cancelPropertyEdit() {
		void goto(resolve('/catalog/properties'));
	}
</script>

<svelte:head>
	<title>Izmeni svojstvo | Popisomator</title>
</svelte:head>

<main class="flex min-h-full flex-col px-4 pt-4 pb-8 sm:px-6">
	<ProtectedPageState
		loading={authPage.state.loading || (authPage.state.authorized && loadingProperty)}
		error={authPage.state.error || error}
		authorized={authPage.state.authorized}
	>
		{#if property}
			<section class="mx-auto flex w-full max-w-2xl flex-1 flex-col" aria-label="Izmeni svojstvo">
				<div class="flex min-h-0 flex-1">
					<PropertyForm {property} onsaved={propertySaved} oncancel={cancelPropertyEdit} />
				</div>
			</section>
		{/if}
	</ProtectedPageState>
</main>
