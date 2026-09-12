<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { onMount } from 'svelte';
	import { api, ApiError, type LocationOption } from '$lib/api';
	import { createAuthPage } from '$lib/state/auth-page.svelte';
	import ProtectedPageState from '$lib/components/shared/ProtectedPageState.svelte';
	import LocationForm from '$lib/components/catalog/LocationForm.svelte';

	const authPage = createAuthPage({
		unavailableMessage: 'Dodavanje lokacije trenutno nije dostupno.',
		requiredRoles: ['admin']
	});

	let locationOptions = $state<LocationOption[]>([]);
	let loading = $state(false);
	let error = $state('');

	onMount(() => {
		void authPage.load().then(() => {
			if (authPage.state.authorized) void loadLocations();
		});
	});

	async function loadLocations() {
		loading = true;
		error = '';

		try {
			locationOptions = await api.getLocationOptionsFlat();
		} catch (reason) {
			error = reason instanceof ApiError ? reason.message : 'Lokacije nisu učitane.';
		} finally {
			loading = false;
		}
	}

	function onCreated() {
		void goto(resolve('/locations'));
	}

	function onCancelCreation() {
		void goto(resolve('/locations'));
	}
</script>

<svelte:head>
	<title>Nova lokacija | Popisomator</title>
</svelte:head>

<main class="flex min-h-full flex-col px-4 pt-4 pb-8 sm:px-6">
	<ProtectedPageState
		loading={authPage.state.loading || (authPage.state.authorized && loading)}
		error={authPage.state.error || error}
		authorized={authPage.state.authorized}
	>
		<section class="mx-auto flex w-full max-w-3xl flex-1 flex-col">
			<div class="flex min-h-0 flex-1">
				<LocationForm
					{locationOptions}
					onsaved={onCreated}
					oncancel={onCancelCreation}
				/>
			</div>
		</section>
	</ProtectedPageState>
</main>
