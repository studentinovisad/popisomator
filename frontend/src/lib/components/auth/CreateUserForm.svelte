<script lang="ts">
	import { api, ApiError, type UserRole } from '$lib/api';
	import PasswordInput from '$lib/components/auth/PasswordInput.svelte';
	import RoleSelect from '$lib/components/auth/RoleSelect.svelte';
	import { emailError, passwordError, requiredTextError } from '$lib/domain/form-validation';
	import { Button, Label, Portal } from 'bits-ui';
	import Plus from '@lucide/svelte/icons/plus';
	import X from '@lucide/svelte/icons/x';
	import { toast } from 'svelte-sonner';

	let {
		oncreated,
		oncancel
	}: {
		oncreated: () => void;
		oncancel?: () => void;
	} = $props();

	let fieldErrors = $state<{ fullName?: string; email?: string; password?: string }>({});
	let creating = $state(false);
	let email = $state('');
	let fullName = $state('');
	let password = $state('');
	let role = $state<UserRole>('user');

	function clearFieldError(field: 'fullName' | 'email' | 'password') {
		if (!fieldErrors[field]) return;
		fieldErrors = { ...fieldErrors, [field]: undefined };
	}

	function validate() {
		fieldErrors = {
			fullName: requiredTextError(fullName, 'ime i prezime'),
			email: emailError(email),
			password: passwordError(password)
		};
		return Object.values(fieldErrors).every((fieldError) => fieldError === undefined);
	}

	async function createUser(event: SubmitEvent) {
		event.preventDefault();
		if (!validate()) return;
		creating = true;

		try {
			await api.createUser({ email, full_name: fullName, password, role });
			email = '';
			fullName = '';
			password = '';
			role = 'user';
			toast.success('Korisnik je dodat.');
			oncreated();
		} catch (reason) {
			toast.error(reason instanceof ApiError ? reason.message : 'Korisnik nije sačuvan.');
		} finally {
			creating = false;
		}
	}
</script>

<Portal to="#page-header-actions">
	<Button.Root
		class="inline-grid size-10 place-items-center rounded-md bg-brand text-on-brand hover:bg-brand-strong disabled:opacity-60"
		disabled={creating}
		form="create-user-form"
		type="submit"
		aria-label="Dodaj korisnika"
		title="Dodaj korisnika"
	>
		<Plus class="size-4" aria-hidden="true" />
	</Button.Root>
	{#if oncancel}
		<Button.Root
			class="inline-grid size-10 place-items-center rounded-md border border-line bg-surface text-ink hover:border-brand/40 hover:bg-brand-soft hover:text-brand disabled:opacity-60"
			disabled={creating}
			type="button"
			onclick={oncancel}
			aria-label="Otkaži dodavanje"
			title="Otkaži dodavanje"
		>
			<X class="size-4" aria-hidden="true" />
		</Button.Root>
	{/if}
</Portal>

<form id="create-user-form" class="grid gap-4 sm:grid-cols-2" novalidate onsubmit={createUser}>
	<div class="block">
		<Label.Root class="text-sm font-medium text-ink" for="new-user-full-name">
			Ime i prezime
		</Label.Root>
		<input
			id="new-user-full-name"
			class={`mt-1 block w-full ${fieldErrors.fullName ? 'field-invalid' : ''}`}
			bind:value={fullName}
			aria-invalid={Boolean(fieldErrors.fullName)}
			aria-describedby={fieldErrors.fullName ? 'new-user-full-name-error' : undefined}
			oninput={() => clearFieldError('fullName')}
		/>
		{#if fieldErrors.fullName}
			<p id="new-user-full-name-error" class="mt-1 text-xs text-danger" role="alert">
				{fieldErrors.fullName}
			</p>
		{/if}
	</div>
	<div class="block">
		<Label.Root class="text-sm font-medium text-ink" for="new-user-email">Email</Label.Root>
		<input
			id="new-user-email"
			class={`mt-1 block w-full ${fieldErrors.email ? 'field-invalid' : ''}`}
			type="email"
			bind:value={email}
			aria-invalid={Boolean(fieldErrors.email)}
			aria-describedby={fieldErrors.email ? 'new-user-email-error' : undefined}
			oninput={() => clearFieldError('email')}
		/>
		{#if fieldErrors.email}
			<p id="new-user-email-error" class="mt-1 text-xs text-danger" role="alert">
				{fieldErrors.email}
			</p>
		{/if}
	</div>
	<div class="block">
		<Label.Root class="text-sm font-medium text-ink" for="new-user-password">Lozinka</Label.Root>
		<PasswordInput
			id="new-user-password"
			className="mt-1"
			bind:value={password}
			autocomplete="new-password"
			minlength={8}
			invalid={Boolean(fieldErrors.password)}
			describedBy={fieldErrors.password ? 'new-user-password-error' : undefined}
			oninput={() => clearFieldError('password')}
		/>
		<span class="mt-1 block text-xs text-muted">
			Najmanje 8 znakova, veliko i malo slovo, i broj.
		</span>
		{#if fieldErrors.password}
			<p id="new-user-password-error" class="mt-1 text-xs text-danger" role="alert">
				{fieldErrors.password}
			</p>
		{/if}
	</div>
	<div class="block">
		<Label.Root class="text-sm font-medium text-ink" for="new-user-role">Uloga</Label.Root>
		<div class="mt-1"><RoleSelect id="new-user-role" bind:value={role} ariaLabel="Uloga" /></div>
	</div>
</form>
