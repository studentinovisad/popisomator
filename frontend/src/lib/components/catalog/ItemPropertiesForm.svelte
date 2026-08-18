<script lang="ts">
	import { api, ApiError, type Item, type ItemType, type PropertyOption } from '$lib/api';
	import ItemPropertyValueInput from '$lib/components/inventory/ItemPropertyValueInput.svelte';
	import { defaultJsonValue } from '$lib/domain/items';
	import { Button, Label } from 'bits-ui';

	let {
		item,
		itemType,
		properties,
		onsaved
	}: {
		item: Item;
		itemType: ItemType;
		properties: PropertyOption[];
		onsaved: () => void;
	} = $props();

	let editablePropertyIDs = $derived([
		...new Set([
			...item.properties.map((property) => property.id),
			...itemType.properties.map((property) => property.id)
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

<form onsubmit={save}>
	<div class="divide-y divide-line border-y border-line">
		{#if editablePropertyIDs.length}
			{#each editablePropertyIDs as propertyID (propertyID)}
				{@const property = propertyByID.get(propertyID)}
				{#if property}
					<div
						class="grid grid-cols-[minmax(0,1fr)_minmax(0,2fr)_auto] items-center py-1.5 sm:grid-cols-[minmax(12rem,1fr)_minmax(0,2fr)_auto]"
					>
						<p class="pr-3 text-sm text-muted sm:pr-6">{property.name}</p>
						<div class="col-start-2">
							{#if selectedPropertyIDs.includes(propertyID)}
								<Label.Root class="sr-only" for={`item-${item.id}-property-${propertyID}`}>
									Vrednost za {property.name}
								</Label.Root>
								<ItemPropertyValueInput
									id={`item-${item.id}-property-${propertyID}`}
									bind:value={values[propertyID]}
									className=""
									inputClassName="text-sm"
									compact={true}
									{property}
									required
								/>
							{:else}
								<p
									class="flex h-8 items-center rounded-md border border-transparent pl-3 text-sm text-muted"
								>
									Nema.
								</p>
							{/if}
						</div>
						<label class="col-start-3 inline-flex size-8 items-center justify-center">
							<input
								type="checkbox"
								bind:group={selectedPropertyIDs}
								value={propertyID}
								aria-label={`Uključi svojstvo ${property.name}`}
							/>
						</label>
					</div>
				{/if}
			{/each}
		{:else}
			<p class="py-3 text-sm text-muted">Ovaj tip nema dostupna svojstva.</p>
		{/if}
	</div>

	<div class="mt-4">
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
