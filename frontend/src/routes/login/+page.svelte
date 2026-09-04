<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { api, ApiError } from '$lib/api';
	import PasswordInput from '$lib/components/auth/PasswordInput.svelte';
	import { emailError, requiredTextError } from '$lib/domain/form-validation';
	import { session } from '$lib/state/session.svelte';
	import { Button, Label } from 'bits-ui';
	import { toast } from 'svelte-sonner';

	let email = $state('');
	let password = $state('');
	let fieldErrors = $state<{ email?: string; password?: string }>({});
	let submitting = $state(false);

	function clearFieldError(field: 'email' | 'password') {
		if (!fieldErrors[field]) return;
		fieldErrors = { ...fieldErrors, [field]: undefined };
	}

	function validate() {
		const nextErrors: typeof fieldErrors = {};
		nextErrors.email = emailError(email);
		nextErrors.password = requiredTextError(password, 'lozinku');
		fieldErrors = nextErrors;
		return Object.values(nextErrors).every((fieldError) => fieldError === undefined);
	}

	async function submit(event: SubmitEvent) {
		event.preventDefault();
		if (!validate()) return;
		submitting = true;

		try {
			await api.login({ email, password });
			await session.refresh();
			await goto(resolve('/'));
		} catch (reason) {
			toast.error(reason instanceof ApiError ? reason.message : 'Prijava trenutno nije dostupna.');
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
		<form class="mt-8 space-y-5 text-left" novalidate onsubmit={submit}>
			<div class="block">
				<Label.Root class="text-sm font-medium text-ink" for="login-email">Email</Label.Root>
				<input
					id="login-email"
					class={`mt-1 block w-full ${fieldErrors.email ? 'field-invalid' : ''}`}
					type="email"
					bind:value={email}
					autocomplete="email"
					aria-invalid={Boolean(fieldErrors.email)}
					aria-describedby={fieldErrors.email ? 'login-email-error' : undefined}
					oninput={() => clearFieldError('email')}
				/>
				{#if fieldErrors.email}
					<p id="login-email-error" class="mt-1 text-xs text-danger" role="alert">
						{fieldErrors.email}
					</p>
				{/if}
			</div>

			<div class="block">
				<Label.Root class="text-sm font-medium text-ink" for="login-password">Lozinka</Label.Root>
				<PasswordInput
					id="login-password"
					className="mt-1"
					bind:value={password}
					autocomplete="current-password"
					invalid={Boolean(fieldErrors.password)}
					describedBy={fieldErrors.password ? 'login-password-error' : undefined}
					oninput={() => clearFieldError('password')}
				/>
				{#if fieldErrors.password}
					<p id="login-password-error" class="mt-1 text-xs text-danger" role="alert">
						{fieldErrors.password}
					</p>
				{/if}
			</div>

			<Button.Root
				class="w-full rounded-md bg-brand px-4 py-2 font-medium text-on-brand hover:bg-brand-strong disabled:opacity-60"
				disabled={submitting}
				type="submit"
			>
				{submitting ? 'Prijavljivanje…' : 'Prijavi se'}
			</Button.Root>

			<div class="space-y-2 pt-2">
				<p class="text-sm font-medium text-ink">Ako još uvek nemate nalog:</p>
				<Button.Root
					class="w-full rounded-md bg-brand px-4 py-2 font-medium text-on-brand hover:bg-brand-strong"
					type="button"
					onclick={registration_page}
				>
					Zatražite pristup
				</Button.Root>
			</div>
		</form>
	</div>
</main>
