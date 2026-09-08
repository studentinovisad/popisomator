<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { onMount } from 'svelte';
	import { api, ApiError, type PropertyOption } from '$lib/api';
	import { createAuthPage } from '$lib/state/auth-page.svelte';
	import CreateItemTypeForm from '$lib/components/catalog/CreateItemTypeForm.svelte';
	import ProtectedPageState from '$lib/components/shared/ProtectedPageState.svelte';

	const authPage = createAuthPage({
		unavailableMessage: 'Dodavanje tipa stavke trenutno nije dostupno.',
		requiredRoles: ['admin']
	});

	let properties = $state<PropertyOption[]>([]);
	let loading = $state(false);
	let error = $state('');

	onMount(() => {
		void authPage.load().then(() => {
			if (authPage.state.authorized) void loadProperties();
		});
	});

	async function loadProperties() {
		loading = true;
		error = '';

		try {
			properties = await api.getPropertyOptions();
		} catch (reason) {
			error = reason instanceof ApiError ? reason.message : 'Svojstva nisu učitana.';
		} finally {
			loading = false;
		}
	}

	function itemTypeCreated() {
		void goto(resolve('/catalog/item-types'));
	}

	function cancelItemTypeCreation() {
		void goto(resolve('/catalog/item-types'));
	}
</script>

<svelte:head>
	<title>Novi tip stavke | Popisomator</title>
</svelte:head>

<main class="flex min-h-full flex-col px-4 pt-4 pb-8 sm:px-6">
	<ProtectedPageState
		loading={authPage.state.loading || (authPage.state.authorized && loading)}
		error={authPage.state.error || error}
		authorized={authPage.state.authorized}
	>
		<section class="mx-auto flex w-full max-w-3xl flex-1 flex-col">
			<div class="flex min-h-0 flex-1">
				<CreateItemTypeForm
					{properties}
					onsaved={itemTypeCreated}
					oncancel={cancelItemTypeCreation}
				/>
			</div>
		</section>
	</ProtectedPageState>
</main>
