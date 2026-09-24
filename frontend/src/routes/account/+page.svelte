<script lang="ts">
	import { onMount } from 'svelte';
	import KeyRound from '@lucide/svelte/icons/key-round';
	import { Button } from 'bits-ui';
	import ChangePasswordDialog from '$lib/components/auth/ChangePasswordDialog.svelte';
	import { createAuthPage } from '$lib/state/auth-page.svelte';
	import ProtectedPageState from '$lib/components/shared/ProtectedPageState.svelte';
	import { userRoleLabel } from '$lib/domain/users';

	const authPage = createAuthPage({ unavailableMessage: 'Nalog trenutno nije dostupan.' });

	let passwordDialogOpen = $state(false);

	onMount(() => {
		void authPage.load();
	});
</script>

<svelte:head>
	<title>Moj nalog | Popisomator</title>
</svelte:head>

<main class="px-4 pt-4 pb-8 sm:px-6">
	<ProtectedPageState
		loading={authPage.state.loading}
		error={authPage.state.error}
		authorized={authPage.state.authorized}
	>
		{#if authPage.state.user}
			<section class="mx-auto max-w-3xl" aria-labelledby="account-heading">
				<div class="border-b border-line pb-5">
					<div class="min-w-0">
						<p class="text-sm text-muted">{userRoleLabel(authPage.state.user.role)}</p>
						<h2 id="account-heading" class="mt-1 truncate text-xl font-semibold text-ink">
							{authPage.state.user.full_name}
						</h2>
					</div>
					<div class="mt-3 flex items-center gap-2 sm:justify-end">
						<Button.Root
							class="inline-grid size-8 place-items-center rounded text-muted transition-colors hover:bg-soft hover:text-ink"
							type="button"
							onclick={() => (passwordDialogOpen = true)}
							aria-label="Promeni lozinku"
							title="Promeni lozinku"
						>
							<KeyRound class="size-4" aria-hidden="true" />
						</Button.Root>
					</div>
				</div>

				<ChangePasswordDialog bind:open={passwordDialogOpen} />

				<section class="mt-3" aria-labelledby="account-details-heading">
					<h3 id="account-details-heading" class="text-base font-semibold text-ink">
						Podaci naloga
					</h3>
					<dl class="mt-3 divide-y divide-line border-y border-line">
						<div
							class="grid grid-cols-[minmax(0,1fr)_minmax(0,2fr)] items-center py-1.5 sm:grid-cols-[minmax(12rem,1fr)_minmax(0,2fr)]"
						>
							<dt class="pr-3 text-sm text-muted sm:pr-6">Ime i prezime</dt>
							<dd
								class="col-start-2 flex h-8 items-center rounded-md border border-transparent pl-3 text-sm text-ink"
							>
								{authPage.state.user.full_name}
							</dd>
						</div>
						<div
							class="grid grid-cols-[minmax(0,1fr)_minmax(0,2fr)] items-center py-1.5 sm:grid-cols-[minmax(12rem,1fr)_minmax(0,2fr)]"
						>
							<dt class="pr-3 text-sm text-muted sm:pr-6">Email</dt>
							<dd
								class="col-start-2 flex h-8 items-center rounded-md border border-transparent pl-3 text-sm text-ink"
							>
								{authPage.state.user.email}
							</dd>
						</div>
						<div
							class="grid grid-cols-[minmax(0,1fr)_minmax(0,2fr)] items-center py-1.5 sm:grid-cols-[minmax(12rem,1fr)_minmax(0,2fr)]"
						>
							<dt class="pr-3 text-sm text-muted sm:pr-6">Uloga</dt>
							<dd
								class="col-start-2 flex h-8 items-center rounded-md border border-transparent pl-3 text-sm text-ink"
							>
								{userRoleLabel(authPage.state.user.role)}
							</dd>
						</div>
					</dl>
				</section>
			</section>
		{/if}
	</ProtectedPageState>
</main>
