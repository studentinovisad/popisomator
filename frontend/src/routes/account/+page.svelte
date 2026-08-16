<script lang="ts">
	import { onMount } from 'svelte';
	import { createAuthPage } from '$lib/auth-page.svelte';
	import PageLoader from '$lib/components/PageLoader.svelte';
	import { userRoleLabel } from '$lib/users';

	const authPage = createAuthPage({ unavailableMessage: 'Nalog trenutno nije dostupan.' });

	onMount(() => {
		void authPage.load();
	});
</script>

<svelte:head>
	<title>Moj nalog | Popisomator</title>
</svelte:head>

<main class="px-4 sm:px-6">
	{#if authPage.state.loading}
		<PageLoader />
	{:else if authPage.state.error}
		<div class="grid min-h-[calc(100svh-14rem)] place-items-center">
			<p class="text-danger" role="alert">{authPage.state.error}</p>
		</div>
	{:else if authPage.state.authorized && authPage.state.user}
		<div class="grid min-h-[calc(100svh-14rem)] place-items-center">
			<div class="w-full max-w-md">
				<dl
					class="divide-y divide-line rounded-lg border border-line bg-surface shadow-sm shadow-black/5"
				>
					<div class="px-4 py-3">
						<dt class="font-mono text-xs tracking-wide text-muted">IME I PREZIME</dt>
						<dd class="mt-1 font-medium text-ink">{authPage.state.user.full_name}</dd>
					</div>
					<div class="px-4 py-3">
						<dt class="font-mono text-xs tracking-wide text-muted">EMAIL</dt>
						<dd class="mt-1 font-medium text-ink">{authPage.state.user.email}</dd>
					</div>
					<div class="px-4 py-3">
						<dt class="font-mono text-xs tracking-wide text-muted">ULOGA</dt>
						<dd class="mt-1 font-medium text-ink">{userRoleLabel(authPage.state.user.role)}</dd>
					</div>
				</dl>
			</div>
		</div>
	{/if}
</main>
