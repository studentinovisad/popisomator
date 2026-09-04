<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { api, ApiError } from '$lib/api';
	import PasswordInput from '$lib/components/auth/PasswordInput.svelte';
	import { emailError, passwordError, requiredTextError } from '$lib/domain/form-validation';
	import { Button, Label } from 'bits-ui';
	import { toast } from 'svelte-sonner';

	let full_name = $state('');
	let email = $state('');
	let password = $state('');
	let fieldErrors = $state<{ fullName?: string; email?: string; password?: string }>({});
	let submitting = $state(false);

	function clearFieldError(field: 'fullName' | 'email' | 'password') {
		if (!fieldErrors[field]) return;
		fieldErrors = { ...fieldErrors, [field]: undefined };
	}

	function validate() {
		const nextErrors: typeof fieldErrors = {};
		nextErrors.fullName = requiredTextError(full_name, 'ime i prezime');
		nextErrors.email = emailError(email);
		nextErrors.password = passwordError(password);
		fieldErrors = nextErrors;
		return Object.values(nextErrors).every((fieldError) => fieldError === undefined);
	}

	async function submit(event: SubmitEvent) {
		event.preventDefault();
		if (!validate()) return;
		submitting = true;

		try {
			await api.register({ full_name, email, password });
			await goto(resolve('/login'));
			toast.success(
				'Uspešno ste poslali zahtev za registraciju. Administrator mora da ga odobri pre prijave.'
			);
		} catch (reason) {
			toast.error(
				reason instanceof ApiError ? reason.message : 'Registracija trenutno nije dostupna.'
			);
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
		<form class="mt-8 space-y-5 text-left" novalidate onsubmit={submit}>
			<div class="block">
				<Label.Root class="text-sm font-medium text-ink" for="registration-name"
					>Ime i prezime</Label.Root
				>
				<input
					id="registration-name"
					class={`mt-1 block w-full ${fieldErrors.fullName ? 'field-invalid' : ''}`}
					type="text"
					bind:value={full_name}
					autocomplete="name"
					aria-invalid={Boolean(fieldErrors.fullName)}
					aria-describedby={fieldErrors.fullName ? 'registration-name-error' : undefined}
					oninput={() => clearFieldError('fullName')}
				/>
				{#if fieldErrors.fullName}
					<p id="registration-name-error" class="mt-1 text-xs text-danger" role="alert">
						{fieldErrors.fullName}
					</p>
				{/if}
			</div>

			<div class="block">
				<Label.Root class="text-sm font-medium text-ink" for="registration-email">Email</Label.Root>
				<input
					id="registration-email"
					class={`mt-1 block w-full ${fieldErrors.email ? 'field-invalid' : ''}`}
					type="email"
					bind:value={email}
					autocomplete="email"
					aria-invalid={Boolean(fieldErrors.email)}
					aria-describedby={fieldErrors.email ? 'registration-email-error' : undefined}
					oninput={() => clearFieldError('email')}
				/>
				{#if fieldErrors.email}
					<p id="registration-email-error" class="mt-1 text-xs text-danger" role="alert">
						{fieldErrors.email}
					</p>
				{/if}
			</div>

			<div class="block">
				<Label.Root class="text-sm font-medium text-ink" for="registration-password"
					>Lozinka</Label.Root
				>
				<PasswordInput
					id="registration-password"
					className="mt-1"
					bind:value={password}
					autocomplete="new-password"
					invalid={Boolean(fieldErrors.password)}
					describedBy={fieldErrors.password ? 'registration-password-error' : undefined}
					oninput={() => clearFieldError('password')}
				/>
				{#if fieldErrors.password}
					<p id="registration-password-error" class="mt-1 text-xs text-danger" role="alert">
						{fieldErrors.password}
					</p>
				{/if}
			</div>

			<Button.Root
				class="w-full rounded-md bg-brand px-4 py-2 font-medium text-on-brand hover:bg-brand-strong disabled:opacity-60"
				disabled={submitting}
				type="submit"
			>
				{submitting ? 'Registrovanje…' : 'Registruj se'}
			</Button.Root>
		</form>
	</div>
</main>
