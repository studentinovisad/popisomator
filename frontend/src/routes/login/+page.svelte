<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { api, ApiError } from '$lib/api';
	import { session } from '$lib/session.svelte';
	import { Button, Label } from 'bits-ui';

	let email = $state('');
	let password = $state('');
	let error = $state('');
	let submitting = $state(false);

	async function submit(event: SubmitEvent) {
		event.preventDefault();
		error = '';
		submitting = true;

		try {
			await api.login({ email, password });
			await session.refresh();
			await goto(resolve('/account'));
		} catch (reason) {
			error = reason instanceof ApiError ? reason.message : 'Prijava trenutno nije dostupna.';
		} finally {
			submitting = false;
		}
	}

	async function registration_page() {
		await goto(resolve('/register'));
	}
</script>

<svelte:head>
	<title>Prijava | Popisomator</title>
</svelte:head>

<main class="grid min-h-[calc(100svh-14rem)] place-items-center px-4 sm:px-6">
	<div class="w-full max-w-md text-center">
		<form class="mt-8 space-y-5 text-left" onsubmit={submit}>
			<div class="block">
				<Label.Root class="text-sm font-medium text-ink" for="login-email">Email</Label.Root>
				<input
					id="login-email"
					class="mt-1 block w-full"
					type="email"
					bind:value={email}
					autocomplete="email"
					required
				/>
			</div>

			<div class="block">
				<Label.Root class="text-sm font-medium text-ink" for="login-password">Lozinka</Label.Root>
				<input
					id="login-password"
					class="mt-1 block w-full"
					type="password"
					bind:value={password}
					autocomplete="current-password"
					required
				/>
			</div>

			{#if error}
				<p class="text-sm text-danger" role="alert">{error}</p>
			{/if}

			<Button.Root
				class="w-full rounded-md bg-brand px-4 py-2 font-medium text-on-brand hover:bg-brand-strong disabled:opacity-60"
				disabled={submitting}
				type="submit"
			>
				{submitting ? 'Prijavljivanje…' : 'Prijavi se'}
			</Button.Root>

			<div class="space-y-2 pt-2">
				<p class="text-sm font-medium text-ink">Ako još uvek nemaš nalog:</p>
				<Button.Root
					class="w-full rounded-md bg-brand px-4 py-2 font-medium text-on-brand hover:bg-brand-strong"
					type="button"
					onclick={registration_page}
				>
					Zatraži pristup
				</Button.Root>
			</div>
		</form>
	</div>
</main>
