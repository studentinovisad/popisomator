<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { onMount } from 'svelte';
	import CreateUserForm from '$lib/components/auth/CreateUserForm.svelte';
	import ProtectedPageState from '$lib/components/shared/ProtectedPageState.svelte';
	import { createAuthPage } from '$lib/state/auth-page.svelte';

	const authPage = createAuthPage({
		unavailableMessage: 'Dodavanje korisnika trenutno nije dostupno.',
		requiredRole: 'admin'
	});

	onMount(() => void authPage.load());

	function userCreated() {
		void goto(resolve('/users'));
	}

	function cancelUserCreation() {
		void goto(resolve('/users'));
	}
</script>

<svelte:head>
	<title>Novi korisnik | Popisomator</title>
</svelte:head>

<main class="flex min-h-full flex-col px-4 pt-4 pb-8 sm:px-6">
	<ProtectedPageState
		loading={authPage.state.loading}
		error={authPage.state.error}
		authorized={authPage.state.authorized}
	>
		<section class="mx-auto flex w-full max-w-2xl flex-1 flex-col" aria-label="Novi korisnik">
			<CreateUserForm oncreated={userCreated} oncancel={cancelUserCreation} />
		</section>
	</ProtectedPageState>
</main>
