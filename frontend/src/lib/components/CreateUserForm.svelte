<script lang="ts">
	import { api, ApiError, type UserRole } from '$lib/api';
	import RoleSelect from '$lib/components/RoleSelect.svelte';
	import { Button, Label } from 'bits-ui';

	let { oncreated }: { oncreated: () => void } = $props();

	let error = $state('');
	let creating = $state(false);
	let email = $state('');
	let fullName = $state('');
	let password = $state('');
	let role = $state<UserRole>('user');

	async function createUser(event: SubmitEvent) {
		event.preventDefault();
		error = '';
		creating = true;

		try {
			await api.createUser({ email, full_name: fullName, password, role });
			email = '';
			fullName = '';
			password = '';
			role = 'user';
			oncreated();
		} catch (reason) {
			error = reason instanceof ApiError ? reason.message : 'Korisnik nije sačuvan.';
		} finally {
			creating = false;
		}
	}
</script>

<form class="grid gap-4 sm:grid-cols-2" onsubmit={createUser}>
	<div class="block">
		<Label.Root class="text-sm font-medium text-ink" for="new-user-full-name">
			Ime i prezime
		</Label.Root>
		<input id="new-user-full-name" class="mt-1 block w-full" bind:value={fullName} required />
	</div>
	<div class="block">
		<Label.Root class="text-sm font-medium text-ink" for="new-user-email">Email</Label.Root>
		<input id="new-user-email" class="mt-1 block w-full" type="email" bind:value={email} required />
	</div>
	<div class="block">
		<Label.Root class="text-sm font-medium text-ink" for="new-user-password">Lozinka</Label.Root>
		<input
			id="new-user-password"
			class="mt-1 block w-full"
			type="password"
			bind:value={password}
			minlength="8"
			required
		/>
		<span class="mt-1 block text-xs text-muted">
			Najmanje 8 znakova, veliko i malo slovo, i broj.
		</span>
	</div>
	<div class="block">
		<Label.Root class="text-sm font-medium text-ink" for="new-user-role">Uloga</Label.Root>
		<div class="mt-1"><RoleSelect id="new-user-role" bind:value={role} ariaLabel="Uloga" /></div>
	</div>
	<div class="sm:col-span-2">
		<Button.Root
			class="rounded-md bg-brand px-4 py-2 font-medium text-on-brand hover:bg-brand-strong disabled:opacity-60"
			disabled={creating}
			type="submit"
		>
			{creating ? 'Čuvanje…' : 'Dodaj korisnika'}
		</Button.Root>
	</div>
</form>
{#if error}
	<p class="mt-3 text-sm text-danger" role="alert">{error}</p>
{/if}
