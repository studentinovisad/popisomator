<script lang="ts">
	import { onMount } from 'svelte';
	import Plus from '@lucide/svelte/icons/plus';
	import { createAuthPage } from '$lib/state/auth-page.svelte';
	import CreateUserForm from '$lib/components/auth/CreateUserForm.svelte';
	import ProtectedPageState from '$lib/components/shared/ProtectedPageState.svelte';
	import UserList from '$lib/components/users/UserList.svelte';
	import { Dialog, Portal } from 'bits-ui';

	const authPage = createAuthPage({
		unavailableMessage: 'Administracija trenutno nije dostupna.',
		requiredRole: 'admin'
	});
	let usersRefreshKey = $state(0);
	let createUserDialogOpen = $state(false);

	onMount(() => {
		void authPage.load();
	});

	function refreshUsers() {
		usersRefreshKey += 1;
		createUserDialogOpen = false;
	}
</script>

<svelte:head>
	<title>Administracija | Popisomator</title>
</svelte:head>

<main class="px-4 pt-4 pb-8 sm:px-6">
	<ProtectedPageState
		loading={authPage.state.loading}
		error={authPage.state.error}
		authorized={authPage.state.authorized}
	>
		<div id="users-summary"></div>

		<Dialog.Root bind:open={createUserDialogOpen}>
			<Portal to="#page-header-actions">
				<Dialog.Trigger
					class="inline-flex size-10 items-center justify-center rounded-md bg-brand text-on-brand hover:bg-brand-strong"
					aria-label="Dodaj korisnika"
					title="Dodaj korisnika"
				>
					<Plus class="size-4" aria-hidden="true" />
				</Dialog.Trigger>
			</Portal>
			<Dialog.Portal>
				<Dialog.Overlay class="fixed inset-0 z-20 bg-black/35 backdrop-blur-sm" />
				<Dialog.Content
					class="fixed top-1/2 left-1/2 z-30 w-[calc(100%-2rem)] max-w-xl -translate-x-1/2 -translate-y-1/2 rounded-lg border border-line bg-surface p-6 shadow-black/20"
				>
					<div class="flex items-start justify-between gap-4">
						<div>
							<Dialog.Title class="text-xl font-semibold text-ink">Novi korisnik</Dialog.Title>
							<Dialog.Description class="mt-1 text-sm text-muted">
								Korisnik će moći da se prijavi odmah nakon kreiranja naloga.
							</Dialog.Description>
						</div>
						<Dialog.Close
							class="rounded-md px-2 py-1 text-sm text-muted hover:bg-soft hover:text-ink"
						>
							Zatvori
						</Dialog.Close>
					</div>
					<div class="mt-6"><CreateUserForm oncreated={refreshUsers} /></div>
				</Dialog.Content>
			</Dialog.Portal>
		</Dialog.Root>
		<UserList refreshKey={usersRefreshKey} currentUserID={authPage.state.user!.id} />
	</ProtectedPageState>
</main>
