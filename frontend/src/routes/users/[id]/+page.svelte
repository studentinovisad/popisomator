<script lang="ts">
	import { onMount } from 'svelte';
	import ArrowLeft from '@lucide/svelte/icons/arrow-left';
	import Pencil from '@lucide/svelte/icons/pencil';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { Button, Portal } from 'bits-ui';
	import { api, ApiError, type UpdateUserRequest, type User, type UserRole } from '$lib/api';
	import RoleSelect from '$lib/components/auth/RoleSelect.svelte';
	import ProtectedPageState from '$lib/components/shared/ProtectedPageState.svelte';
	import { emailError, requiredTextError } from '$lib/domain/form-validation';
	import { userRoleLabel } from '$lib/domain/users';
	import { createAuthPage } from '$lib/state/auth-page.svelte';
	import { session } from '$lib/state/session.svelte';
	import { toast } from 'svelte-sonner';

	const authPage = createAuthPage({
		unavailableMessage: 'Korisnik trenutno nije dostupan.',
		requiredRole: 'admin'
	});

	let editing = $state(false);
	let email = $state('');
	let emailErrorMessage = $state('');
	let fullName = $state('');
	let fullNameError = $state('');
	let loading = $state(false);
	let loadError = $state('');
	let role = $state<UserRole>('user');
	let saving = $state(false);
	let user = $state<User | null>(null);
	let canChangeRole = $derived(user?.id !== authPage.state.user?.id);

	onMount(() => {
		void authPage.load().then(() => {
			if (authPage.state.authorized) void loadUser();
		});
	});

	function userID() {
		const id = Number(page.params.id);
		return Number.isSafeInteger(id) && id > 0 ? id : null;
	}

	async function loadUser() {
		const id = userID();
		if (id === null) {
			loadError = 'Korisnik nije pronađen.';
			return;
		}

		loading = true;
		loadError = '';
		try {
			user = await api.getUser(id);
		} catch (reason) {
			loadError =
				reason instanceof ApiError && reason.status === 404
					? 'Korisnik nije pronađen.'
					: 'Korisnik nije učitan.';
		} finally {
			loading = false;
		}
	}

	function startEditing() {
		if (!user) return;
		fullName = user.full_name;
		email = user.email;
		role = user.role;
		fullNameError = '';
		emailErrorMessage = '';
		editing = true;
	}

	function cancelEditing() {
		editing = false;
		fullNameError = '';
		emailErrorMessage = '';
	}

	function validate() {
		fullNameError = requiredTextError(fullName, 'ime i prezime') ?? '';
		emailErrorMessage = emailError(email) ?? '';
		return !fullNameError && !emailErrorMessage;
	}

	async function save(event: SubmitEvent) {
		event.preventDefault();
		fullName = fullName.trim();
		email = email.trim();
		if (!validate() || !user) {
			toast.error(fullNameError || emailErrorMessage || 'Podaci korisnika nisu ispravni.');
			return;
		}

		const update: UpdateUserRequest = {};
		if (fullName !== user.full_name) update.full_name = fullName;
		if (email !== user.email) update.email = email;
		if (canChangeRole && role !== user.role) update.role = role;
		if (Object.keys(update).length === 0) {
			cancelEditing();
			return;
		}

		saving = true;
		try {
			const updatedUser = await api.updateUser(user.id, update);
			user = updatedUser;
			if (updatedUser.id === authPage.state.user?.id) {
				authPage.state.user = updatedUser;
				session.setUser(updatedUser);
			}
			editing = false;
			toast.success('Podaci korisnika su sačuvani.');
		} catch (reason) {
			if (reason instanceof ApiError && reason.status === 409) {
				emailErrorMessage = 'Email adresa je već zauzeta.';
			}
			toast.error(reason instanceof ApiError ? reason.message : 'Podaci korisnika nisu sačuvani.');
		} finally {
			saving = false;
		}
	}
</script>

<svelte:head>
	<title>Korisnik | Popisomator</title>
</svelte:head>

<main class="px-4 pt-4 pb-8 sm:px-6">
	<ProtectedPageState
		loading={authPage.state.loading || (authPage.state.authorized && loading)}
		error={authPage.state.error || loadError}
		authorized={authPage.state.authorized}
	>
		<Portal to="#page-header-actions">
			<a
				class="inline-flex size-10 items-center justify-center rounded-md border border-line bg-surface text-ink transition-colors hover:border-brand/40 hover:bg-brand-soft hover:text-brand"
				href={resolve('/users')}
				aria-label="Nazad na korisnike"
				title="Nazad na korisnike"
			>
				<ArrowLeft class="size-4" aria-hidden="true" />
			</a>
		</Portal>

		{#if user}
			<section class="mx-auto max-w-3xl" aria-labelledby="user-heading">
				<div class="border-b border-line pb-5">
					<div class="min-w-0">
						<p class="text-sm text-muted">{userRoleLabel(user.role)}</p>
						<h2 id="user-heading" class="mt-1 truncate text-xl font-semibold text-ink">
							{user.full_name}
						</h2>
					</div>
					<div class="mt-3 flex items-center gap-2 sm:justify-end">
						<Button.Root
							class={`inline-grid size-8 place-items-center rounded transition-colors ${
								editing ? 'bg-brand-soft text-brand' : 'text-muted hover:bg-soft hover:text-ink'
							}`}
							type="button"
							disabled={saving}
							onclick={() => (editing ? cancelEditing() : startEditing())}
							aria-label={editing ? 'Zatvori izmenu korisnika' : 'Izmeni korisnika'}
							aria-pressed={editing}
							title={editing ? 'Zatvori izmenu' : 'Izmeni'}
						>
							<Pencil class="size-4" aria-hidden="true" />
						</Button.Root>
					</div>
				</div>

				<section class="mt-3" aria-labelledby="user-details-heading">
					<h3 id="user-details-heading" class="text-base font-semibold text-ink">Podaci naloga</h3>
					<form class="mt-3" novalidate onsubmit={save}>
						<dl class="divide-y divide-line border-y border-line">
							<div
								class="grid grid-cols-[minmax(0,1fr)_minmax(0,2fr)] items-center py-1.5 sm:grid-cols-[minmax(12rem,1fr)_minmax(0,2fr)]"
							>
								<dt class="pr-3 text-sm text-muted sm:pr-6">Ime i prezime</dt>
								<dd class="col-start-2">
									{#if editing}
										<input
											class={`block h-8 w-full text-sm ${fullNameError ? 'field-invalid' : ''}`}
											aria-label="Ime i prezime"
											aria-invalid={Boolean(fullNameError)}
											title={fullNameError}
											bind:value={fullName}
											oninput={() => (fullNameError = '')}
										/>
									{:else}
										<span
											class="flex h-8 items-center rounded-md border border-transparent pl-3 text-sm text-ink"
										>
											{user.full_name}
										</span>
									{/if}
								</dd>
							</div>
							<div
								class="grid grid-cols-[minmax(0,1fr)_minmax(0,2fr)] items-center py-1.5 sm:grid-cols-[minmax(12rem,1fr)_minmax(0,2fr)]"
							>
								<dt class="pr-3 text-sm text-muted sm:pr-6">Email</dt>
								<dd class="col-start-2">
									{#if editing}
										<input
											class={`block h-8 w-full text-sm ${emailErrorMessage ? 'field-invalid' : ''}`}
											type="email"
											aria-label="Email"
											aria-invalid={Boolean(emailErrorMessage)}
											title={emailErrorMessage}
											bind:value={email}
											oninput={() => (emailErrorMessage = '')}
										/>
									{:else}
										<span
											class="flex h-8 items-center rounded-md border border-transparent pl-3 text-sm text-ink"
										>
											{user.email}
										</span>
									{/if}
								</dd>
							</div>
							<div
								class="grid grid-cols-[minmax(0,1fr)_minmax(0,2fr)] items-center py-1.5 sm:grid-cols-[minmax(12rem,1fr)_minmax(0,2fr)]"
							>
								<dt class="pr-3 text-sm text-muted sm:pr-6">Uloga</dt>
								<dd class="col-start-2">
									{#if editing}
										<RoleSelect
											value={role}
											ariaLabel="Uloga korisnika"
											compact={true}
											disabled={!canChangeRole || saving}
											onvaluechange={(nextRole) => (role = nextRole)}
										/>
									{:else}
										<span
											class="flex h-8 items-center rounded-md border border-transparent pl-3 text-sm text-ink"
										>
											{userRoleLabel(user.role)}
										</span>
									{/if}
								</dd>
							</div>
						</dl>

						{#if editing}
							<div class="mt-4">
								<Button.Root
									class="rounded-md bg-brand px-4 py-2 text-sm font-medium text-on-brand hover:bg-brand-strong disabled:opacity-60"
									type="submit"
									disabled={saving}
								>
									{saving ? 'Čuvanje…' : 'Sačuvaj izmene'}
								</Button.Root>
							</div>
						{/if}
					</form>
				</section>
			</section>
		{/if}
	</ProtectedPageState>
</main>
