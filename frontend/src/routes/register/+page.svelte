<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { api, ApiError } from '$lib/api';
	import { session } from '$lib/session.svelte';
	import { Button, Label } from 'bits-ui';

	let full_name = $state('');
	let email = $state('');
	let password = $state('');
	let error = $state('');
	let submit_message = $state('');
	let submitting = $state(false);

	async function submit(event: SubmitEvent) {
		event.preventDefault();
		error = '';
		submit_message = '';
		submitting = true;

		try {
			await api.register({ full_name, email, password });
			submit_message = 'Uspešno poslat zahtev za registraciju. Administrator će morati da Vam odobri zahtev za registraciju.'
		} catch (reason) {
			error = reason instanceof ApiError ? reason.message : 'Registracija trenutno nije dostupna.';
		} finally {
			submitting = false;
		}
	}
</script>

<svelte:head>
	<title>Registracija | Popisomator</title>
</svelte:head>

<main class="grid min-h-[calc(100svh-14rem)] place-items-center px-4 sm:px-6">
	<div class="w-full max-w-md text-center">
		<form class="mt-8 space-y-5 text-left" onsubmit={submit}>
		<div class="block">
				<Label.Root class="text-sm font-medium text-ink" for="registration-name">Ime i prezime</Label.Root>
				<input
					id="registration-name"
					class="mt-1 block w-full"
					type="text"
					bind:value={full_name}
					autocomplete="name"
					required
				/>
			</div>

			<div class="block">
				<Label.Root class="text-sm font-medium text-ink" for="registration-email">Email</Label.Root>
				<input
					id="registration-email"
					class="mt-1 block w-full"
					type="email"
					bind:value={email}
					autocomplete="email"
					required
				/>
			</div>

			<div class="block">
				<Label.Root class="text-sm font-medium text-ink" for="registration-password">Lozinka</Label.Root>
				<input
					id="registration-password"
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
			{#if submit_message}
				<p class="text-sm text-success" role="alert">{submit_message}</p>
			{:else}
				<Button.Root
					class="w-full rounded-md bg-brand px-4 py-2 font-medium text-on-brand hover:bg-brand-strong disabled:opacity-60"
					disabled={submitting}
					type="submit"
				>
					{submitting ? 'Registrovanje…' : 'Registruj se'}
				</Button.Root>
			{/if}
		</form>
	</div>
</main>
