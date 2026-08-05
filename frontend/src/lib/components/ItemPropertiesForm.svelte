<script lang="ts">
	import { api, ApiError, type Item, type ItemType, type Property } from '$lib/api';
	import ItemPropertyValueInput from '$lib/components/ItemPropertyValueInput.svelte';
	import { defaultJsonValue } from '$lib/items';
	import { Button, Label } from 'bits-ui';

	let {
		item,
		itemTypes,
		properties,
		onsaved
	}: {
		item: Item;
		itemTypes: ItemType[];
		properties: Property[];
		onsaved: () => void;
	} = $props();

	let type = $derived(itemTypes.find((itemType) => itemType.id === item.type_id));
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
	let error = $state('');

	$effect(() => {
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
	{#if editablePropertyIDs.length}
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
					{#if property.description}<p class="mt-1 text-xs text-muted">
							{property.description}
						</p>{/if}
				</div>
			{/if}
		{/each}
	{:else}
		<p class="text-sm text-muted">Ovaj tip nema dostupna svojstva.</p>
	{/if}
	<div class="border-t border-line pt-4">
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
