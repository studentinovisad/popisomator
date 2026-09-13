<script lang="ts">
	import { onMount } from 'svelte';
	import Plus from '@lucide/svelte/icons/plus';
	import { resolve } from '$app/paths';
	import { createAuthPage } from '$lib/state/auth-page.svelte';
	import ProtectedPageState from '$lib/components/shared/ProtectedPageState.svelte';
	import UserList from '$lib/components/users/UserList.svelte';
	import { Portal } from 'bits-ui';

	const authPage = createAuthPage({
		unavailableMessage: 'Stranica Korisnici trenutno nije dostupna.',
		requiredRole: 'admin'
	});
	let usersError = $state('');

	onMount(() => {
		void authPage.load();
	});
</script>

<svelte:head>
	<title>Korisnici | Popisomator</title>
</svelte:head>

<main class="px-4 pt-4 pb-8 sm:px-6">
	<ProtectedPageState
		loading={authPage.state.loading}
		error={authPage.state.error || usersError}
		authorized={authPage.state.authorized}
	>
		<div id="users-summary"></div>

		<Portal to="#page-header-actions">
			<a
				class="inline-flex size-10 items-center justify-center rounded-md bg-brand text-on-brand hover:bg-brand-strong"
				href={resolve('/users/new')}
				aria-label="Dodaj korisnika"
				title="Dodaj korisnika"
			>
				<Plus class="size-4" aria-hidden="true" />
			</a>
		</Portal>
		<UserList onloaderror={(message) => (usersError = message)} />
	</ProtectedPageState>
</main>
