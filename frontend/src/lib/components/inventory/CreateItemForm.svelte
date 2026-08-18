<script lang="ts">
	import { api, ApiError, type ItemType, type Property } from '$lib/api';
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
		itemTypes: ItemType[];
		properties: Property[];
		oncreated: () => void;
		oncancel?: () => void;
	} = $props();

	let selectedTypeID = $state('');
	let amount = $state(1);
	let propertyValues = $state<Record<number, string>>({});
	let error = $state('');
	let creating = $state(false);

	let selectedType = $derived(itemTypes.find((itemType) => itemType.id === Number(selectedTypeID)));
	let propertiesByID = $derived(new Map(properties.map((property) => [property.id, property])));

	function selectType(value: string) {
		selectedTypeID = value;
		const itemType = itemTypes.find((candidate) => candidate.id === Number(value));
		propertyValues = Object.fromEntries(
			(itemType?.properties ?? []).flatMap((itemProperty) => {
				const property = propertiesByID.get(itemProperty.id);
				return property
					? [[itemProperty.id, defaultJsonValue(property.value_type, itemProperty.default_value)]]
					: [];
			})
		);
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

	{#if selectedType?.properties.length}
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
						{#if property.description}
							<p class="mt-1 text-xs text-muted">{property.description}</p>
						{/if}
					</div>
				{/if}
			{/each}
		</fieldset>
	{/if}

	<Separator.Root class="h-px bg-line" decorative />
	<div class="flex items-center gap-3">
		<Button.Root
			class="rounded-md bg-brand px-4 py-2 text-sm font-medium text-on-brand hover:bg-brand-strong disabled:opacity-60"
			disabled={creating || !selectedTypeID}
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
