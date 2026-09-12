<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { onMount } from 'svelte';
	import { api, ApiError, type Location, type LocationOption } from '$lib/api';
	import { createAuthPage } from '$lib/state/auth-page.svelte';
	import ProtectedPageState from '$lib/components/shared/ProtectedPageState.svelte';
	import LocationForm from '$lib/components/catalog/LocationForm.svelte';

	const authPage = createAuthPage({
		unavailableMessage: 'Izmena lokacije trenutno nije dostupna.',
		requiredRoles: ['admin']
	});

	let location = $state<Location | null>(null);
	let locationOptions = $state<LocationOption[]>([]);
	let loading = $state(false);
	let error = $state('');

	onMount(() => {
		void authPage.load().then(() => {
			if (authPage.state.authorized) void loadLocations();
		});
	});

	async function loadLocations() {
		const id = Number(page.params.id);
		if (!Number.isSafeInteger(id) || id < 1) {
			error = 'Lokacija nije pronađena.';
			return;
		}

		loading = true;
		error = '';
		try {
			const [nextLocation, nextLocationOptions] = await Promise.all([
				api.getLocation(id),
				api.getLocationOptionsFlat()
			]);
			location = nextLocation;
			locationOptions = nextLocationOptions;
		} catch (reason) {
			error =
				reason instanceof ApiError && reason.status === 404
					? 'Lokacija nije pronađena.'
					: 'Lokacija nije učitana.';
		} finally {
			loading = false;
		}
	}

	function onSaved() {
		void goto(resolve('/locations'));
	}

	function cancelLocationEdit() {
		void goto(resolve('/locations'));
	}
</script>

<svelte:head>
	<title>Izmeni lokaciju | Popisomator</title>
</svelte:head>

<main class="flex min-h-full flex-col px-4 pt-4 pb-8 sm:px-6">
	<ProtectedPageState
		loading={authPage.state.loading || (authPage.state.authorized && loading)}
		error={authPage.state.error || error}
		authorized={authPage.state.authorized}
	>
		{#if location}
			<section class="mx-auto flex w-full max-w-3xl flex-1 flex-col" aria-label="Izmeni tip stavke">
				<div class="flex min-h-0 flex-1">
					<LocationForm {location} {locationOptions} onsaved={onSaved} oncancel={cancelLocationEdit}/>
				</div>
			</section>
		{/if}
	</ProtectedPageState>
</main>
