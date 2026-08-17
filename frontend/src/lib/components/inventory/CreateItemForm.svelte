<script lang="ts">
	import { api, ApiError, type ItemType, type ItemTypeOption, type Property } from '$lib/api';
	import ItemPropertyValueInput from '$lib/components/inventory/ItemPropertyValueInput.svelte';
	import NumberInput from '$lib/components/shared/NumberInput.svelte';
	import { defaultJsonValue } from '$lib/domain/items';
	import { Button, Label, Select, Separator } from 'bits-ui';

	let {
		itemTypes,
		properties,
		oncreated,
		oncancel
	}: {
		itemTypes: ItemTypeOption[];
		properties: Property[];
		oncreated: () => void;
		oncancel?: () => void;
	} = $props();

	let selectedTypeID = $state('');
	let amount = $state(1);
	let propertyValues = $state<Record<number, string>>({});
	let error = $state('');
	let creating = $state(false);

	let typeOptions = $derived(
		itemTypes.map((itemType) => ({ value: String(itemType.id), label: itemType.name }))
	);
	let selectedType = $state<ItemType | null>(null);
	let loadingType = $state(false);
	let typeLoadVersion = 0;
	let propertiesByID = $derived(new Map(properties.map((property) => [property.id, property])));

	async function selectType(value: string) {
		const version = ++typeLoadVersion;
		selectedTypeID = value;
		selectedType = null;
		propertyValues = {};
		if (!value) {
			loadingType = false;
			return;
		}

		loadingType = true;
		error = '';
		try {
			const itemType = await api.getItemType(Number(value));
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
			error = reason instanceof ApiError ? reason.message : 'Svojstva tipa nisu učitana.';
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
					class="z-40 w-(--bits-select-anchor-width) rounded-md border border-line bg-surface p-1 shadow-lg shadow-black/15"
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
