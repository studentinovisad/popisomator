<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { onMount } from 'svelte';
	import { api, ApiError, type ItemType, type Property } from '$lib/api';
	import { createAuthPage } from '$lib/state/auth-page.svelte';
	import CreateItemTypeForm from '$lib/components/catalog/CreateItemTypeForm.svelte';
	import ProtectedPageState from '$lib/components/shared/ProtectedPageState.svelte';

	const authPage = createAuthPage({
		unavailableMessage: 'Izmena tipa stavke trenutno nije dostupna.',
		requiredRoles: ['admin']
	});

	let itemType = $state<ItemType | null>(null);
	let properties = $state<Property[]>([]);
	let loading = $state(false);
	let error = $state('');

	onMount(() => {
		void authPage.load().then(() => {
			if (authPage.state.authorized) void loadItemType();
		});
	});

	async function loadItemType() {
		const id = Number(page.params.id);
		if (!Number.isSafeInteger(id) || id < 1) {
			error = 'Tip stavke nije pronađen.';
			return;
		}

		loading = true;
		error = '';
		try {
			const [nextItemType, nextProperties] = await Promise.all([
				api.getItemType(id),
				api.listProperties()
			]);
			itemType = nextItemType;
			properties = nextProperties;
		} catch (reason) {
			error =
				reason instanceof ApiError && reason.status === 404
					? 'Tip stavke nije pronađen.'
					: 'Tip stavke nije učitan.';
		} finally {
			loading = false;
		}
	}

	function itemTypeSaved() {
		void goto(resolve('/catalog/item-types'));
	}
</script>

<svelte:head>
	<title>Izmeni tip stavke | Popisomator</title>
</svelte:head>

<main class="px-4 pt-4 pb-8 sm:px-6">
	<ProtectedPageState
		loading={authPage.state.loading || (authPage.state.authorized && loading)}
		error={authPage.state.error || error}
		authorized={authPage.state.authorized}
	>
		{#if itemType}
			<section class="mx-auto max-w-3xl" aria-labelledby="edit-item-type-heading">
				<div class="border-b border-line pb-4">
					<h2 id="edit-item-type-heading" class="text-lg font-semibold text-ink">
						Izmeni tip stavke
					</h2>
					<p class="mt-1 text-sm text-muted">{itemType.name}</p>
				</div>
				<div class="mt-6">
					<CreateItemTypeForm {itemType} {properties} onsaved={itemTypeSaved} />
				</div>
			</section>
		{/if}
	</ProtectedPageState>
</main>
