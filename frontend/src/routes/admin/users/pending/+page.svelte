<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import ArrowLeft from '@lucide/svelte/icons/arrow-left';
	import { createAuthPage } from '$lib/auth-page.svelte';
	import ProtectedPageState from '$lib/components/ProtectedPageState.svelte';
	import PendingRegistrations from '$lib/components/PendingRegistrations.svelte';
	import { Portal } from 'bits-ui';

	const authPage = createAuthPage({
		unavailableMessage: 'Zahtevi za registraciju trenutno nisu dostupni.',
		requiredRole: 'admin'
	});

	onMount(() => void authPage.load());

	function returnToUsers() {
		void goto(resolve('/admin/users'));
	}
</script>

<svelte:head>
	<title>Zahtevi za registraciju | Popisomator</title>
</svelte:head>

<main class="px-4 pt-4 pb-8 sm:px-6">
	<ProtectedPageState
		loading={authPage.state.loading}
		error={authPage.state.error}
		authorized={authPage.state.authorized}
	>
		<Portal to="#page-header-actions">
			<a
				class="inline-flex size-10 items-center justify-center rounded-md border border-line bg-surface text-ink transition-colors hover:border-brand/40 hover:bg-brand-soft hover:text-brand"
				href={resolve('/admin/users')}
				aria-label="Nazad na korisnike"
				title="Nazad na korisnike"
			>
				<ArrowLeft class="size-4" aria-hidden="true" />
			</a>
		</Portal>
		<PendingRegistrations onempty={returnToUsers} />
	</ProtectedPageState>
</main>
