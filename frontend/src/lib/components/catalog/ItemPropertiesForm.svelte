<script lang="ts">
	import { untrack } from 'svelte';
	import {
		api,
		ApiError,
		type Item,
		type ItemType,
		type ItemTypeOption,
		type PropertyOption
	} from '$lib/api';
	import ItemPropertyValueInput from '$lib/components/inventory/ItemPropertyValueInput.svelte';
	import OptionCombobox from '$lib/components/shared/OptionCombobox.svelte';
	import { defaultJsonValue } from '$lib/domain/items';
	import { Button, Label, Separator } from 'bits-ui';

	let {
		item,
		itemTypes,
		properties,
		onitemtypechange,
		onsaved
	}: {
		item: Item;
		itemTypes: ItemTypeOption[];
		properties: PropertyOption[];
		onitemtypechange: (typeID: number) => Promise<void>;
		onsaved: () => void;
	} = $props();

	let type = $state<ItemType | null>(null);
	let typeLoadVersion = 0;
	let loadedTypeID = $state<number | null>(null);
	let loadingType = $state(false);
	let editablePropertyIDs = $derived([
		...new Set([
			...item.properties.map((property) => property.id),
			...(type?.properties.map((property) => property.id) ?? [])
		])
	]);
	let propertyByID = $derived(new Map(properties.map((property) => [property.id, property])));
	let originalValues = $derived(
		new Map(item.properties.map((property) => [property.id, property.value]))
	);
	let selectedPropertyIDs = $state<number[]>([]);
	let values = $state<Record<number, string>>({});
	let saving = $state(false);
	let changingType = $state(false);
	let selectedTypeID = $state('');
	let error = $state('');

	$effect(() => {
		const typeID = item.type_id;
		untrack(() => {
			if (loadedTypeID !== typeID) void loadItemType(typeID);
		});
	});

	$effect(() => {
		selectedTypeID = String(item.type_id);
		selectedPropertyIDs = item.properties.map((property) => property.id);
		values = Object.fromEntries(
			editablePropertyIDs.flatMap((id) => {
				const property = propertyByID.get(id);
				return property
					? [[id, originalValues.get(id) ?? defaultJsonValue(property.value_type, null)]]
					: [];
			})
		);
	});

	async function loadItemType(typeID: number) {
		const version = ++typeLoadVersion;
		loadingType = true;
		type = null;

		try {
			const nextType = await api.getItemType(typeID);
			if (version !== typeLoadVersion) return;
			type = nextType;
			loadedTypeID = typeID;
		} catch (reason) {
			if (version !== typeLoadVersion) return;
			error = reason instanceof ApiError ? reason.message : 'Tip stavke nije učitan.';
		} finally {
			if (version === typeLoadVersion) loadingType = false;
		}
	}

	async function changeItemType(value: string) {
		const typeID = Number(value);
		if (!Number.isSafeInteger(typeID) || typeID === item.type_id) return;

		changingType = true;
		error = '';
		try {
			await onitemtypechange(typeID);
			await loadItemType(typeID);
		} catch (reason) {
			selectedTypeID = String(item.type_id);
			error = reason instanceof ApiError ? reason.message : 'Tip stavke nije promenjen.';
		} finally {
			changingType = false;
		}
	}

	async function save(event: SubmitEvent) {
		event.preventDefault();
		saving = true;
		error = '';

		try {
			const changes: Promise<unknown>[] = [];
			for (const propertyID of editablePropertyIDs) {
				const wasSelected = originalValues.has(propertyID);
				const isSelected = selectedPropertyIDs.includes(propertyID);
				if (wasSelected && !isSelected) changes.push(api.removeItemProperty(item.id, propertyID));
				if (!wasSelected && isSelected)
					changes.push(api.addItemProperty(item.id, propertyID, values[propertyID]));
				if (wasSelected && isSelected && originalValues.get(propertyID) !== values[propertyID])
					changes.push(api.updateItemProperty(item.id, propertyID, values[propertyID]));
			}
			await Promise.all(changes);
			onsaved();
		} catch (reason) {
			error = reason instanceof ApiError ? reason.message : 'Svojstva stavke nisu sačuvana.';
		} finally {
			saving = false;
		}
	}
</script>

<form class="grid gap-4" onsubmit={save}>
	<div>
		<Label.Root class="text-sm font-medium text-ink" for="item-type">Tip stavke</Label.Root>
		<div class="mt-1">
			<OptionCombobox
				id="item-type"
				options={itemTypes}
				bind:value={selectedTypeID}
				placeholder="Odaberite tip"
				disabled={changingType || loadingType || saving}
				onvaluechange={changeItemType}
			/>
		</div>
		<p class="mt-1 text-xs text-muted">
			Promenom tipa dostupna svojstva se prilagođavaju novom tipu.
		</p>
	</div>

	{#if loadingType}
		<p class="text-sm text-muted">Učitavanje svojstava tipa…</p>
	{:else if editablePropertyIDs.length}
		{#each editablePropertyIDs as propertyID (propertyID)}
			{@const property = propertyByID.get(propertyID)}
			{#if property}
				<div class="border-b border-line pb-4 last:border-b-0 last:pb-0">
					<label class="flex items-center gap-2 text-sm font-medium text-ink">
						<input type="checkbox" bind:group={selectedPropertyIDs} value={propertyID} />
						{property.name}
					</label>
					{#if selectedPropertyIDs.includes(propertyID)}
						<Label.Root class="sr-only" for={`item-${item.id}-property-${propertyID}`}>
							Vrednost za {property.name}
						</Label.Root>
						<ItemPropertyValueInput
							id={`item-${item.id}-property-${propertyID}`}
							bind:value={values[propertyID]}
							{property}
							required
						/>
					{/if}
				</div>
			{/if}
		{/each}
	{:else}
		<p class="text-sm text-muted">Ovaj tip nema dostupna svojstva.</p>
	{/if}
	<Separator.Root class="h-px bg-line" decorative />
	<div>
		<Button.Root
			class="rounded-md bg-brand px-4 py-2 text-sm font-medium text-on-brand hover:bg-brand-strong disabled:opacity-60"
			disabled={saving}
			type="submit"
		>
			{saving ? 'Čuvanje…' : 'Sačuvaj svojstva'}
		</Button.Root>
		{#if error}<p class="mt-3 text-sm text-danger" role="alert">{error}</p>{/if}
	</div>
</form>
