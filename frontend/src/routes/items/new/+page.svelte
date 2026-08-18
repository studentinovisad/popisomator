<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { onMount } from 'svelte';
	import { api, ApiError, type ItemTypeOption, type PropertyOption } from '$lib/api';
	import { createAuthPage } from '$lib/state/auth-page.svelte';
	import CreateItemForm from '$lib/components/inventory/CreateItemForm.svelte';
	import ProtectedPageState from '$lib/components/shared/ProtectedPageState.svelte';

	const authPage = createAuthPage({
		unavailableMessage: 'Dodavanje stavke trenutno nije dostupno.',
		requiredRoles: ['admin', 'manager']
	});

	let itemTypes = $state<ItemTypeOption[]>([]);
	let properties = $state<PropertyOption[]>([]);
	let loading = $state(false);
	let error = $state('');

	onMount(() => {
		void authPage.load().then(() => {
			if (authPage.state.authorized) void loadFormData();
		});
	});

	async function loadFormData() {
		loading = true;
		error = '';

		try {
			[itemTypes, properties] = await Promise.all([
				api.getItemTypeOptions(),
				api.getPropertyOptions()
			]);
		} catch (reason) {
			error = reason instanceof ApiError ? reason.message : 'Podaci za stavku nisu učitani.';
		} finally {
			loading = false;
		}
	}

	function itemCreated() {
		void goto(resolve('/'));
	}

	function cancelItemCreation() {
		void goto(resolve('/'));
	}
</script>

<svelte:head>
	<title>Nova stavka | Popisomator</title>
</svelte:head>

<main class="px-4 pt-4 pb-8 sm:px-6">
	<ProtectedPageState
		loading={authPage.state.loading || (authPage.state.authorized && loading)}
		error={authPage.state.error || error}
		authorized={authPage.state.authorized}
	>
		<section class="mx-auto max-w-3xl" aria-labelledby="new-item-heading">
			<div class="border-b border-line pb-4">
				<h2 id="new-item-heading" class="text-lg font-semibold text-ink">Nova stavka</h2>
				<p class="mt-1 text-sm text-muted">Dodajte stavku i njene početne vrednosti svojstava.</p>
			</div>
			<div class="mt-6">
				<CreateItemForm
					{itemTypes}
					{properties}
					oncreated={itemCreated}
					oncancel={cancelItemCreation}
				/>
			</div>
		</section>
	</ProtectedPageState>
</main>
