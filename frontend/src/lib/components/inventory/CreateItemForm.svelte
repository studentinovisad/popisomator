<script lang="ts">
	import { api, ApiError, type ItemType, type ItemTypeOption, type PropertyOption } from '$lib/api';
	import ItemPropertyValueInput from '$lib/components/inventory/ItemPropertyValueInput.svelte';
	import NumberInput from '$lib/components/shared/NumberInput.svelte';
	import OptionCombobox from '$lib/components/shared/OptionCombobox.svelte';
	import { defaultJsonValue } from '$lib/domain/items';
	import { Button, Label, Separator } from 'bits-ui';

	let {
		itemTypes,
		properties,
		oncreated,
		oncancel
	}: {
		itemTypes: ItemTypeOption[];
		properties: PropertyOption[];
		oncreated: () => void;
		oncancel?: () => void;
	} = $props();

	let selectedTypeID = $state('');
	let amount = $state(1);
	let propertyValues = $state<Record<number, string>>({});
	let error = $state('');
	let creating = $state(false);
	let loadingType = $state(false);
	let selectedType = $state<ItemType | null>(null);
	let typeLoadVersion = 0;

	let propertiesByID = $derived(new Map(properties.map((property) => [property.id, property])));

	async function selectType(value: string) {
		selectedTypeID = value;
		selectedType = null;
		propertyValues = {};

		const typeID = Number(value);
		if (!Number.isSafeInteger(typeID)) return;

		const version = ++typeLoadVersion;
		loadingType = true;
		error = '';
		try {
			const itemType = await api.getItemType(typeID);
			if (version !== typeLoadVersion) return;

			selectedType = itemType;
			propertyValues = Object.fromEntries(
				itemType.properties.flatMap((itemProperty) => {
					const property = propertiesByID.get(itemProperty.id);
					return property
						? [[itemProperty.id, defaultJsonValue(property.value_type, itemProperty.default_value)]]
						: [];
				})
			);
		} catch (reason) {
			if (version !== typeLoadVersion) return;
			error = reason instanceof ApiError ? reason.message : 'Tip stavke nije učitan.';
		} finally {
			if (version === typeLoadVersion) loadingType = false;
		}
	}

	async function createItem(event: SubmitEvent) {
		event.preventDefault();
		error = '';
		if (!selectedTypeID) {
			error = 'Odaberite tip stavke.';
			return;
		}
		creating = true;

		try {
			await api.createItem({
				type_id: Number(selectedTypeID),
				properties: Object.entries(propertyValues).map(([id, value]) => ({
					id: Number(id),
					value
				})),
				amount
			});
			selectedTypeID = '';
			propertyValues = {};
			oncreated();
		} catch (reason) {
			error = reason instanceof ApiError ? reason.message : 'Stavka nije sačuvana.';
		} finally {
			creating = false;
		}
	}
</script>

<form class="grid gap-4" onsubmit={createItem}>
	<div>
		<Label.Root class="text-sm font-medium text-ink" for="new-item-type">Tip stavke</Label.Root>
		<div class="mt-1">
			<OptionCombobox
				id="new-item-type"
				options={itemTypes}
				bind:value={selectedTypeID}
				placeholder="Odaberite tip"
				onvaluechange={selectType}
			/>
		</div>
		<p class="mt-1 text-xs text-muted">Tip određuje skup dostupnih svojstava.</p>
		{#if itemTypes.length === 0}
			<p class="mt-2 text-sm text-danger" role="alert">Prvo dodajte tip stavke u katalogu.</p>
		{/if}
	</div>

	<Label.Root class="text-sm font-medium text-ink" for="new-item-amount">Količina</Label.Root>
	<NumberInput
		id="new-item-amount"
		bind:value={amount}
		ariaLabel="Količina"
		placeholder="Unesite količinu"
		min={1}
		max={100}
		required
	/>

	{#if loadingType}
		<p class="text-sm text-muted">Učitavanje svojstava tipa…</p>
	{:else if selectedType?.properties.length}
		<fieldset class="grid gap-3 border-t border-line pt-4">
			<legend class="text-sm font-medium text-ink">Svojstva</legend>
			{#each selectedType.properties as itemProperty (itemProperty.id)}
				{@const property = propertiesByID.get(itemProperty.id)}
				{#if property}
					<div>
						<Label.Root
							class="text-sm font-medium text-ink"
							for={`new-item-property-${property.id}`}
						>
							{property.name}
						</Label.Root>
						<ItemPropertyValueInput
							id={`new-item-property-${property.id}`}
							bind:value={propertyValues[property.id]}
							{property}
							required
						/>
					</div>
				{/if}
			{/each}
		</fieldset>
	{/if}

	<Separator.Root class="h-px bg-line" decorative />
	<div class="flex items-center gap-3">
		<Button.Root
			class="rounded-md bg-brand px-4 py-2 text-sm font-medium text-on-brand hover:bg-brand-strong disabled:opacity-60"
			disabled={creating || loadingType || !selectedType}
			type="submit"
		>
			{creating ? 'Čuvanje…' : 'Dodaj stavku'}
		</Button.Root>
		{#if oncancel}
			<Button.Root
				class="rounded-md border border-line bg-surface px-4 py-2 text-sm font-medium text-ink hover:bg-soft disabled:opacity-60"
				disabled={creating}
				type="button"
				onclick={oncancel}
			>
				Otkaži dodavanje
			</Button.Root>
		{/if}
		{#if error}<p class="text-sm text-danger" role="alert">{error}</p>{/if}
	</div>
</form>
