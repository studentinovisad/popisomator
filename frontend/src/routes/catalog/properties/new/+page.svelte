<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { onMount } from 'svelte';
	import { createAuthPage } from '$lib/state/auth-page.svelte';
	import ProtectedPageState from '$lib/components/shared/ProtectedPageState.svelte';
	import PropertyForm from '$lib/components/catalog/PropertyForm.svelte';

	const authPage = createAuthPage({
		unavailableMessage: 'Dodavanje svojstva trenutno nije dostupno.',
		requiredRoles: ['admin']
	});

	onMount(() => void authPage.load());

	function propertyCreated() {
		void goto(resolve('/catalog/properties'));
	}

	function cancelPropertyCreation() {
		void goto(resolve('/catalog/properties'));
	}
</script>

<svelte:head>
	<title>Novo svojstvo | Popisomator</title>
</svelte:head>

<main class="flex min-h-full flex-col px-4 pt-4 pb-8 sm:px-6">
	<ProtectedPageState
		loading={authPage.state.loading}
		error={authPage.state.error}
		authorized={authPage.state.authorized}
	>
		<section class="mx-auto flex w-full max-w-2xl flex-1 flex-col" aria-label="Novo svojstvo">
			<div class="flex min-h-0 flex-1">
				<PropertyForm onsaved={propertyCreated} oncancel={cancelPropertyCreation} />
			</div>
		</section>
	</ProtectedPageState>
</main>
