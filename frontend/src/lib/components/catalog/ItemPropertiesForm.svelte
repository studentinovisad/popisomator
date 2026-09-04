<script lang="ts">
	import {
		api,
		ApiError,
		type Item,
		type ItemType,
		type PropertyOption,
		type PropertyValue
	} from '$lib/api';
	import ItemPropertyValueInput from '$lib/components/inventory/ItemPropertyValueInput.svelte';
	import { defaultJsonValue, samePropertyValue } from '$lib/domain/items';
	import { propertyValueError } from '$lib/domain/form-validation';
	import { Button, Label } from 'bits-ui';
	import { toast } from 'svelte-sonner';

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
		...itemType.properties.map((property) => property.id),
		...item.properties
			.filter(
				(property) => !itemType.properties.some((typeProperty) => typeProperty.id === property.id)
			)
			.map((property) => property.id)
	]);
	let propertyByID = $derived(new Map(properties.map((property) => [property.id, property])));
	let originalValues = $derived(
		new Map(item.properties.map((property) => [property.id, property.value]))
	);
	let selectedPropertyIDs = $state<number[]>([]);
	let values = $state<Record<number, PropertyValue>>({});
	let fieldErrors = $state<Record<number, string>>({});
	let saving = $state(false);

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

	function clearFieldError(propertyID: number) {
		if (!fieldErrors[propertyID]) return;
		const nextErrors = { ...fieldErrors };
		delete nextErrors[propertyID];
		fieldErrors = nextErrors;
	}

	function validate() {
		const nextErrors: Record<number, string> = {};
		for (const propertyID of selectedPropertyIDs) {
			const property = propertyByID.get(propertyID);
			if (!property) continue;
			const message = propertyValueError(property, values[propertyID]);
			if (message) nextErrors[propertyID] = message;
		}
		fieldErrors = nextErrors;
		return Object.keys(nextErrors).length === 0;
	}

	async function save(event: SubmitEvent) {
		event.preventDefault();
		if (!validate()) return;
		saving = true;

		try {
			const changes: Promise<unknown>[] = [];
			for (const propertyID of editablePropertyIDs) {
				const wasSelected = originalValues.has(propertyID);
				const isSelected = selectedPropertyIDs.includes(propertyID);
				if (wasSelected && !isSelected) changes.push(api.removeItemProperty(item.id, propertyID));
				if (!wasSelected && isSelected)
					changes.push(api.addItemProperty(item.id, propertyID, values[propertyID]));
				// Structured property types (price, mass, volume) hold objects, so compare by content:
				// a reference check reports "unchanged" for every edit and silently drops it.
				if (
					wasSelected &&
					isSelected &&
					!samePropertyValue(originalValues.get(propertyID), values[propertyID])
				)
					changes.push(api.updateItemProperty(item.id, propertyID, values[propertyID]));
			}
			await Promise.all(changes);
			if (changes.length > 0) toast.success('Svojstva stavke su sačuvana.');
			onsaved();
		} catch (reason) {
			toast.error(reason instanceof ApiError ? reason.message : 'Svojstva stavke nisu sačuvana.');
		} finally {
			saving = false;
		}
	}
</script>

<form novalidate onsubmit={save}>
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
									inputClassName={`text-sm ${fieldErrors[propertyID] ? 'field-invalid' : ''}`}
									compact={true}
									{property}
									onvaluechange={() => clearFieldError(propertyID)}
								/>
								{#if fieldErrors[propertyID]}
									<p class="mt-1 text-xs text-danger" role="alert">{fieldErrors[propertyID]}</p>
								{/if}
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
	</div>
</form>
