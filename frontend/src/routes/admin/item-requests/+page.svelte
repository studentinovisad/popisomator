<script lang="ts">
	import { onMount } from 'svelte';
	import ItemRequestsList from '$lib/components/admin/ItemRequestsList.svelte';
	import ProtectedPageState from '$lib/components/shared/ProtectedPageState.svelte';
	import { createAuthPage } from '$lib/state/auth-page.svelte';

	const authPage = createAuthPage({
		unavailableMessage: 'Zahtevi za stavke trenutno nisu dostupni.',
		requiredRole: 'admin'
	});

	onMount(() => void authPage.load());
</script>

<svelte:head>
	<title>Zahtevi za stavke | Popisomator</title>
</svelte:head>

<main class="px-4 pt-4 pb-8 sm:px-6">
	<ProtectedPageState
		loading={authPage.state.loading}
		error={authPage.state.error}
		authorized={authPage.state.authorized}
	>
		<ItemRequestsList />
	</ProtectedPageState>
</main>
