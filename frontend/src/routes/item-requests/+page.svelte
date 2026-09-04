<script lang="ts">
	import { onMount } from 'svelte';
	import ItemRequestsList from '$lib/components/admin/ItemRequestsList.svelte';
	import ProtectedPageState from '$lib/components/shared/ProtectedPageState.svelte';
	import { createAuthPage } from '$lib/state/auth-page.svelte';

	const authPage = createAuthPage({
		unavailableMessage: 'Zahtevi trenutno nisu dostupni.',
		requiredRoles: ['manager', 'admin']
	});
	let requestsError = $state('');

	onMount(() => void authPage.load());
</script>

<svelte:head>
	<title>Zahtevi | Popisomator</title>
</svelte:head>

<main class="px-4 pt-4 pb-8 sm:px-6">
	<ProtectedPageState
		loading={authPage.state.loading}
		error={authPage.state.error || requestsError}
		authorized={authPage.state.authorized}
	>
		<ItemRequestsList onloaderror={(message) => (requestsError = message)} />
	</ProtectedPageState>
</main>
