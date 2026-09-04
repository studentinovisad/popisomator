<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import ArrowLeft from '@lucide/svelte/icons/arrow-left';
	import Pencil from '@lucide/svelte/icons/pencil';
	import Trash2 from '@lucide/svelte/icons/trash-2';
	import { onMount } from 'svelte';
	import { Button, Portal } from 'bits-ui';
	import {
		api,
		ApiError,
		type ConsumptionStatus,
		type Item,
		type ItemType,
		type PropertyOption
	} from '$lib/api';
	import ItemConsumptionControl from '$lib/components/inventory/ItemConsumptionControl.svelte';
	import ItemPropertiesForm from '$lib/components/catalog/ItemPropertiesForm.svelte';
	import ProtectedPageState from '$lib/components/shared/ProtectedPageState.svelte';
	import { displayJson } from '$lib/domain/items';
	import { createAuthPage } from '$lib/state/auth-page.svelte';
	import { toast } from 'svelte-sonner';

	const authPage = createAuthPage({ unavailableMessage: 'Stavka trenutno nije dostupna.' });

	let item = $state<Item | null>(null);
	let itemType = $state<ItemType | null>(null);
	let properties = $state<PropertyOption[]>([]);
	let loading = $state(false);
	let loadError = $state('');
	let editing = $state(false);
	let changingConsumption = $state(false);
	let deleting = $state(false);
	let canManage = $derived(
		authPage.state.user?.role === 'admin' || authPage.state.user?.role === 'manager'
	);
	let propertyNames = $derived(new Map(properties.map((property) => [property.id, property.name])));

	onMount(() => {
		void authPage.load().then(() => {
			if (authPage.state.authorized) void loadItem();
		});
	});

	function itemID() {
		const id = Number(page.params.id);
		return Number.isSafeInteger(id) && id > 0 ? id : null;
	}

	async function loadItem() {
		const id = itemID();
		if (id === null) {
			loadError = 'Stavka nije pronađena.';
			return;
		}

		loading = true;
		loadError = '';
		try {
			const nextItem = await api.getItem(id);
			const [nextItemType, nextProperties] = await Promise.all([
				api.getItemType(nextItem.type_id),
				api.getPropertyOptions()
			]);
			item = nextItem;
			itemType = nextItemType;
			properties = nextProperties;
		} catch (reason) {
			loadError =
				reason instanceof ApiError && reason.status === 404
					? 'Stavka nije pronađena.'
					: 'Stavka nije učitana.';
		} finally {
			loading = false;
		}
	}

	async function refreshItem() {
		const id = itemID();
		if (id === null) return;
		item = await api.getItem(id);
	}

	async function changeConsumption(status: ConsumptionStatus) {
		if (!item || item.consumption === status) return;

		changingConsumption = true;
		try {
			item = await api.consumeItem(item.id, status);
			toast.success('Stanje stavke je promenjeno.');
		} catch (reason) {
			toast.error(reason instanceof ApiError ? reason.message : 'Stanje stavke nije promenjeno.');
		} finally {
			changingConsumption = false;
		}
	}

	async function requestItemUsage(itemID: number, reason: string) {
		await api.createPersonalItemRequest({ item_id: itemID, reason });
		await refreshItem();
	}

	async function deleteItem() {
		if (!item) return;

		const name = item.derived_name || typeName();
		if (!confirm(`Obrisati stavku „${name}“?`)) return;

		deleting = true;
		try {
			await api.deleteItem(item.id);
			toast.success('Stavka je obrisana.');
			await goto(resolve('/'));
		} catch (reason) {
			toast.error(reason instanceof ApiError ? reason.message : 'Stavka nije obrisana.');
		} finally {
			deleting = false;
		}
	}

	function propertiesSaved() {
		editing = false;
		void refreshItem().catch((reason) => {
			toast.error(reason instanceof ApiError ? reason.message : 'Stavka nije osvežena.');
		});
	}

	function typeName() {
		return itemType?.name ?? 'Nepoznat tip';
	}
</script>

<svelte:head>
	<title>Stavka | Popisomator</title>
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
				href={resolve('/')}
				aria-label="Nazad na stavke"
				title="Nazad na stavke"
			>
				<ArrowLeft class="size-4" aria-hidden="true" />
			</a>
		</Portal>

		{#if item}
			<section class="mx-auto max-w-4xl" aria-labelledby="item-heading">
				<div class="border-b border-line pb-5">
					<div class="min-w-0">
						<p class="text-sm text-muted">{typeName()}</p>
						<h2 id="item-heading" class="mt-1 truncate text-xl font-semibold text-ink">
							{item.derived_name || typeName()}
						</h2>
					</div>
					<div class="mt-3 flex min-w-0 items-center gap-2">
						<ItemConsumptionControl
							{item}
							{canManage}
							disabled={changingConsumption}
							onconsumptionchange={(_, status) => void changeConsumption(status)}
							onrequest={requestItemUsage}
						/>
						{#if canManage}
							<Button.Root
								class={`inline-grid size-8 place-items-center rounded transition-colors sm:ml-auto ${
									editing ? 'bg-brand-soft text-brand' : 'text-muted hover:bg-soft hover:text-ink'
								}`}
								type="button"
								onclick={() => (editing = !editing)}
								aria-label={editing ? 'Zatvori izmenu stavke' : 'Izmeni stavku'}
								aria-pressed={editing}
								title={editing ? 'Zatvori izmenu' : 'Izmeni'}
							>
								<Pencil class="size-4" aria-hidden="true" />
							</Button.Root>
							<Button.Root
								class="inline-grid size-8 place-items-center rounded text-danger transition-colors hover:bg-danger-soft disabled:opacity-60"
								disabled={deleting}
								type="button"
								onclick={() => void deleteItem()}
								aria-label="Obriši stavku"
								title="Obriši"
							>
								<Trash2 class="size-4" aria-hidden="true" />
							</Button.Root>
						{/if}
					</div>
				</div>

				<section class="mt-3" aria-labelledby="item-properties-heading">
					<h3 id="item-properties-heading" class="text-base font-semibold text-ink">Svojstva</h3>
					<div class="mt-3">
						{#if editing && canManage && itemType}
							<ItemPropertiesForm {item} {itemType} {properties} onsaved={propertiesSaved} />
						{:else if item.properties.length}
							<dl class="divide-y divide-line border-y border-line">
								{#each item.properties as property (property.id)}
									<div
										class="grid grid-cols-[minmax(0,1fr)_minmax(0,2fr)_auto] items-center py-1.5 sm:grid-cols-[minmax(12rem,1fr)_minmax(0,2fr)_auto]"
									>
										<dt class="pr-3 text-sm text-muted sm:pr-6">
											{propertyNames.get(property.id) ?? `Svojstvo #${property.id}`}
										</dt>
										<dd
											class="col-start-2 flex min-h-8 items-center rounded-md border border-transparent pl-3 text-sm wrap-break-word text-ink"
										>
											{displayJson(property.value_type, property.value)}
										</dd>
										<span class="col-start-3 size-8" aria-hidden="true"></span>
									</div>
								{/each}
							</dl>
						{:else}
							<p class="text-sm text-muted">Ova stavka nema uneta svojstva.</p>
						{/if}
					</div>
				</section>
			</section>
		{/if}
	</ProtectedPageState>
</main>
