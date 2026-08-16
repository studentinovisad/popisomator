<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { onMount } from 'svelte';
	import { createAuthPage } from '$lib/auth-page.svelte';
	import ProtectedPageState from '$lib/components/ProtectedPageState.svelte';
	import PropertyForm from '$lib/components/PropertyForm.svelte';

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

<main class="px-4 pt-4 pb-8 sm:px-6">
	<ProtectedPageState
		loading={authPage.state.loading}
		error={authPage.state.error}
		authorized={authPage.state.authorized}
	>
		<section class="mx-auto max-w-3xl" aria-labelledby="new-property-heading">
			<div class="border-b border-line pb-4">
				<h2 id="new-property-heading" class="text-lg font-semibold text-ink">Novo svojstvo</h2>
				<p class="mt-1 text-sm text-muted">Podesite tip i podrazumevanu vrednost svojstva.</p>
			</div>
			<div class="mt-6">
				<PropertyForm onsaved={propertyCreated} oncancel={cancelPropertyCreation} />
			</div>
		</section>
	</ProtectedPageState>
</main>
