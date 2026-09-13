<script lang="ts">
	import { onMount } from 'svelte';
	import AuditLogList from '$lib/components/admin/AuditLogList.svelte';
	import ProtectedPageState from '$lib/components/shared/ProtectedPageState.svelte';
	import { createAuthPage } from '$lib/state/auth-page.svelte';

	const authPage = createAuthPage({
		unavailableMessage: 'Dnevnik izmena trenutno nije dostupan.',
		requiredRole: 'admin'
	});
	let auditError = $state('');

	onMount(() => void authPage.load());
</script>

<svelte:head>
	<title>Dnevnik izmena | Popisomator</title>
</svelte:head>

<main class="px-4 pt-4 pb-8 sm:px-6">
	<ProtectedPageState
		loading={authPage.state.loading}
		error={authPage.state.error || auditError}
		authorized={authPage.state.authorized}
	>
		<AuditLogList onloaderror={(message) => (auditError = message)} />
	</ProtectedPageState>
</main>
