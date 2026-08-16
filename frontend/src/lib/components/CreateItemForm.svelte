<script lang="ts">
	import { api, ApiError, type ItemType, type Property } from '$lib/api';
	import ItemPropertyValueInput from '$lib/components/ItemPropertyValueInput.svelte';
	import { defaultJsonValue } from '$lib/items';
	import { Button, Label, Select } from 'bits-ui';

	let {
		itemTypes,
		properties,
		oncreated
	}: {
		itemTypes: ItemType[];
		properties: Property[];
		oncreated: () => void;
	} = $props();

	let selectedTypeID = $state('');
	let amount = $state(1);
	let propertyValues = $state<Record<number, string>>({});
	let error = $state('');
	let creating = $state(false);

	let typeOptions = $derived(
		itemTypes.map((itemType) => ({ value: String(itemType.id), label: itemType.name }))
	);
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
				properties: Object.entries(propertyValues).map(([id, value]) => ({ id: Number(id), value })),
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
		<Select.Root
			type="single"
			value={selectedTypeID}
			items={typeOptions}
			onValueChange={selectType}
		>
			<Select.Trigger
				id="new-item-type"
				class="mt-1 flex h-10 w-full items-center justify-between rounded-md border border-line bg-surface px-3 text-sm text-ink"
				aria-required="true"
			>
				<Select.Value placeholder="Odaberite tip" />
			</Select.Trigger>
			<Select.Portal>
				<Select.Content
					class="z-40 w-[var(--bits-select-anchor-width)] rounded-md border border-line bg-surface p-1 shadow-lg shadow-ink/10"
					sideOffset={4}
				>
					<Select.Viewport>
						{#each typeOptions as option (option.value)}
							<Select.Item
								value={option.value}
								label={option.label}
								class="cursor-pointer rounded px-3 py-2 text-sm outline-none data-highlighted:bg-brand-soft"
							>
								{option.label}
							</Select.Item>
						{/each}
					</Select.Viewport>
				</Select.Content>
			</Select.Portal>
		</Select.Root>
		<p class="mt-1 text-xs text-muted">Tip određuje skup dostupnih svojstava.</p>
		{#if itemTypes.length === 0}
			<p class="mt-2 text-sm text-danger" role="alert">Prvo dodajte tip stavke u katalogu.</p>
		{/if}
	</div>

	<Label.Root class="text-sm font-medium text-ink" for="new-item-amount">Količina</Label.Root>
	<input
		id="new-item-amount"
		class="block w-full"
		type="number"
		bind:value={amount}
		placeholder="Unesite količinu"
		aria-required="true"
		min="1"
		max="100"
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

	<div class="flex items-center gap-3 border-t border-line pt-4">
		<Button.Root
			class="rounded-md bg-brand px-4 py-2 text-sm font-medium text-on-brand hover:bg-brand-strong disabled:opacity-60"
			disabled={creating || !selectedTypeID}
			type="submit"
		>
			{creating ? 'Čuvanje…' : 'Dodaj stavku'}
		</Button.Root>
		{#if error}<p class="text-sm text-danger" role="alert">{error}</p>{/if}
	</div>
</form>
