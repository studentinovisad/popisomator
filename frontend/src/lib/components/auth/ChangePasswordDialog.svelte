<script lang="ts">
	import { Button, Dialog, Label } from 'bits-ui';
	import { api, ApiError } from '$lib/api';
	import PasswordInput from '$lib/components/auth/PasswordInput.svelte';
	import { passwordError } from '$lib/domain/form-validation';
	import { toast } from 'svelte-sonner';

	let { open = $bindable(false) }: { open?: boolean } = $props();

	let oldPassword = $state('');
	let newPassword = $state('');
	let confirmPassword = $state('');
	let fieldErrors = $state<{
		oldPassword?: string;
		newPassword?: string;
		confirmPassword?: string;
	}>({});
	let saving = $state(false);

	function resetForm() {
		oldPassword = '';
		newPassword = '';
		confirmPassword = '';
		fieldErrors = {};
	}

	$effect(() => {
		if (!open) resetForm();
	});

	function clearFieldError(field: 'oldPassword' | 'newPassword' | 'confirmPassword') {
		if (!fieldErrors[field]) return;
		fieldErrors = { ...fieldErrors, [field]: undefined };
	}

	function validate() {
		fieldErrors = {
			oldPassword: oldPassword ? undefined : 'Unesite trenutnu lozinku.',
			newPassword: passwordError(newPassword),
			confirmPassword: confirmPassword === newPassword ? undefined : 'Lozinke se ne poklapaju.'
		};
		return Object.values(fieldErrors).every((fieldError) => fieldError === undefined);
	}

	async function changePassword(event: SubmitEvent) {
		event.preventDefault();
		if (!validate()) return;
		saving = true;

		try {
			await api.changePassword({ old_password: oldPassword, new_password: newPassword });
			toast.success('Lozinka je promenjena.');
			open = false;
		} catch (reason) {
			if (reason instanceof ApiError && reason.status === 400) {
				fieldErrors = { ...fieldErrors, oldPassword: reason.message };
			} else {
				toast.error(reason instanceof ApiError ? reason.message : 'Lozinka nije promenjena.');
			}
		} finally {
			saving = false;
		}
	}
</script>

<Dialog.Root bind:open>
	<Dialog.Portal>
		<Dialog.Overlay class="fixed inset-0 z-20 bg-black/35 backdrop-blur-sm" />
		<Dialog.Content
			class="fixed top-1/2 left-1/2 z-30 w-[calc(100%-2rem)] max-w-md -translate-x-1/2 -translate-y-1/2 rounded-lg border border-line bg-surface p-6 shadow-black/20"
		>
			<div class="flex items-start justify-between gap-4">
				<div>
					<Dialog.Title class="text-xl font-semibold text-ink">Promena lozinke</Dialog.Title>
					<Dialog.Description class="mt-1 text-sm text-muted">
						Unesite trenutnu i novu lozinku.
					</Dialog.Description>
				</div>
				<Dialog.Close class="rounded-md px-2 py-1 text-sm text-muted hover:bg-soft hover:text-ink">
					Zatvori
				</Dialog.Close>
			</div>
			<form class="mt-6 grid gap-4" novalidate onsubmit={changePassword}>
				<div class="block">
					<Label.Root class="text-sm font-medium text-ink" for="change-password-old">
						Trenutna lozinka
					</Label.Root>
					<PasswordInput
						id="change-password-old"
						className="mt-1"
						bind:value={oldPassword}
						autocomplete="current-password"
						invalid={Boolean(fieldErrors.oldPassword)}
						describedBy={fieldErrors.oldPassword ? 'change-password-old-error' : undefined}
						oninput={() => clearFieldError('oldPassword')}
					/>
					{#if fieldErrors.oldPassword}
						<p id="change-password-old-error" class="mt-1 text-xs text-danger" role="alert">
							{fieldErrors.oldPassword}
						</p>
					{/if}
				</div>
				<div class="block">
					<Label.Root class="text-sm font-medium text-ink" for="change-password-new">
						Nova lozinka
					</Label.Root>
					<PasswordInput
						id="change-password-new"
						className="mt-1"
						bind:value={newPassword}
						autocomplete="new-password"
						minlength={8}
						invalid={Boolean(fieldErrors.newPassword)}
						describedBy={fieldErrors.newPassword ? 'change-password-new-error' : undefined}
						oninput={() => clearFieldError('newPassword')}
					/>
					<span class="mt-1 block text-xs text-muted">
						Najmanje 8 znakova, veliko i malo slovo, i broj.
					</span>
					{#if fieldErrors.newPassword}
						<p id="change-password-new-error" class="mt-1 text-xs text-danger" role="alert">
							{fieldErrors.newPassword}
						</p>
					{/if}
				</div>
				<div class="block">
					<Label.Root class="text-sm font-medium text-ink" for="change-password-confirm">
						Potvrda nove lozinke
					</Label.Root>
					<PasswordInput
						id="change-password-confirm"
						className="mt-1"
						bind:value={confirmPassword}
						autocomplete="new-password"
						invalid={Boolean(fieldErrors.confirmPassword)}
						describedBy={fieldErrors.confirmPassword ? 'change-password-confirm-error' : undefined}
						oninput={() => clearFieldError('confirmPassword')}
					/>
					{#if fieldErrors.confirmPassword}
						<p id="change-password-confirm-error" class="mt-1 text-xs text-danger" role="alert">
							{fieldErrors.confirmPassword}
						</p>
					{/if}
				</div>
				<Button.Root
					class="mt-2 w-full rounded-md bg-brand px-4 py-2 font-medium text-on-brand hover:bg-brand-strong disabled:opacity-60"
					type="submit"
					disabled={saving}
				>
					{saving ? 'Čuvanje…' : 'Promeni lozinku'}
				</Button.Root>
			</form>
		</Dialog.Content>
	</Dialog.Portal>
</Dialog.Root>
