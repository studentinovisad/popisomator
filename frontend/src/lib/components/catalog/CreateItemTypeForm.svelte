<script lang="ts">
	import { api, ApiError, type ItemType, type PropertyOption } from '$lib/api';
	import ItemPropertyValueInput from '$lib/components/inventory/ItemPropertyValueInput.svelte';
	import MultiOptionCombobox from '$lib/components/shared/MultiOptionCombobox.svelte';
	import { defaultJsonValue, propertyValueTypeLabel } from '$lib/domain/items';
	import { Button, Label, ScrollArea, Separator } from 'bits-ui';

	let {
		itemType,
		properties,
		onsaved,
		oncancel
	}: {
		itemType?: ItemType;
		properties: PropertyOption[];
		onsaved: () => void;
		oncancel?: () => void;
	} = $props();

	let name = $state('');
	let description = $state('');
	let selectedPropertyValues = $state<string[]>([]);
	let selectedPropertyIDs = $derived(selectedPropertyValues.map(Number));
	let defaultValues = $state<Record<number, string>>({});
	let editedDefaultPropertyIDs = $state<Set<number>>(new Set());
	let creating = $state(false);
	let error = $state('');
	let initializedItemTypeID = $state<number | undefined>(undefined);
	let propertiesByID = $derived(new Map(properties.map((property) => [property.id, property])));
	let selectedProperties = $derived(
		properties.filter((property) => selectedPropertyIDs.includes(property.id))
	);
	let originalProperties = $derived(
		new Map(itemType?.properties.map((property) => [property.id, property]))
	);

	$effect(() => {
		if (initializedItemTypeID === itemType?.id) return;

		initializedItemTypeID = itemType?.id;
		name = itemType?.name ?? '';
		description = itemType?.description ?? '';
		selectedPropertyValues = itemType?.properties.map((property) => String(property.id)) ?? [];
		defaultValues = Object.fromEntries(
			(itemType?.properties ?? []).flatMap((itemTypeProperty) => {
				const property = propertiesByID.get(itemTypeProperty.id);
				return property
					? [
							[
								itemTypeProperty.id,
								defaultJsonValue(
									property.value_type,
									itemTypeProperty.default_value ?? property.default_value
								)
							]
						]
					: [];
			})
		);
		editedDefaultPropertyIDs = new Set();
	});

	function updateSelectedProperties(values: string[]) {
		selectedPropertyValues = values;
		for (const propertyID of values.map(Number)) {
			if (defaultValues[propertyID] !== undefined) continue;

			const property = propertiesByID.get(propertyID);
			if (property) {
				defaultValues = {
					...defaultValues,
					[propertyID]: defaultJsonValue(property.value_type, property.default_value)
				};
			}
		}
	}

	function markDefaultEdited(id: number) {
		editedDefaultPropertyIDs = new Set([...editedDefaultPropertyIDs, id]);
	}

	async function createItemType(event: SubmitEvent) {
		event.preventDefault();
		creating = true;
		error = '';

		try {
			if (!itemType) {
				await api.createItemType({
					name,
					description,
					properties: selectedPropertyIDs.map((id) => ({
						id,
						default_value: defaultValues[id],
						visibility: 'overview'
					}))
				});
			} else {
				await api.updateItemType(itemType.id, { name, description });

				const changes: Promise<unknown>[] = [];
				for (const property of itemType.properties) {
					if (!selectedPropertyIDs.includes(property.id)) {
						changes.push(api.removeItemTypeProperty(itemType.id, property.id));
					}
				}
				for (const propertyID of selectedPropertyIDs) {
					if (!originalProperties.has(propertyID)) {
						changes.push(
							api.addItemTypeProperty(itemType.id, {
								property_id: propertyID,
								default_value: defaultValues[propertyID]
							})
						);
					} else if (editedDefaultPropertyIDs.has(propertyID)) {
						changes.push(
							api.updateItemTypeProperty(itemType.id, propertyID, {
								default_value: defaultValues[propertyID]
							})
						);
					}
				}
				await Promise.all(changes);
			}

			onsaved();
		} catch (reason) {
			error = reason instanceof ApiError ? reason.message : 'Tip stavke nije sačuvan.';
		} finally {
			creating = false;
		}
	}
</script>

<form class="grid gap-4" onsubmit={createItemType}>
	<div>
		<Label.Root class="text-sm font-medium text-ink" for="item-type-name">Naziv</Label.Root>
		<input id="item-type-name" class="mt-1 block w-full" bind:value={name} required />
	</div>
	<div>
		<Label.Root class="text-sm font-medium text-ink" for="item-type-description">Opis</Label.Root>
		<textarea id="item-type-description" class="mt-1 block min-h-20 w-full" bind:value={description}
		></textarea>
	</div>
	<fieldset class="border-t border-line pt-4">
		<legend class="text-sm font-medium text-ink">Svojstva tipa</legend>
		<p class="mt-1 text-xs text-muted">
			Podesite svojstva i njihove podrazumevane vrednosti za tip.
		</p>
		<MultiOptionCombobox
			options={properties}
			bind:values={selectedPropertyValues}
			placeholder="Pretražite svojstva"
			emptyMessage="Nema odgovarajućih svojstava."
			onvaluechange={updateSelectedProperties}
		/>
		{#if selectedProperties.length}
			<ScrollArea.Root class="mt-3 h-64 overflow-hidden rounded-md border border-line" type="auto">
				<ScrollArea.Viewport class="h-full w-full">
					<div class="divide-y divide-line">
						{#each selectedProperties as property (property.id)}
							<div class="px-3 py-2">
								<p class="text-sm font-medium text-ink">{property.name}</p>
								<p class="text-xs text-muted">{propertyValueTypeLabel(property.value_type)}</p>
								<Label.Root
									class="mt-3 block text-xs font-medium text-muted"
									for={`item-type-property-${property.id}`}
								>
									Podrazumevana vrednost
								</Label.Root>
								<ItemPropertyValueInput
									id={`item-type-property-${property.id}`}
									bind:value={defaultValues[property.id]}
									{property}
									onvaluechange={() => markDefaultEdited(property.id)}
								/>
							</div>
						{/each}
					</div>
				</ScrollArea.Viewport>
				<ScrollArea.Scrollbar class="flex w-2.5 touch-none bg-soft p-0.5" orientation="vertical">
					<ScrollArea.Thumb class="flex-1 rounded-full bg-line" />
				</ScrollArea.Scrollbar>
			</ScrollArea.Root>
		{:else if properties.length === 0}
			<p class="mt-3 text-sm text-muted">Najpre dodajte svojstvo.</p>
		{/if}
	</fieldset>
	<Separator.Root class="h-px bg-line" decorative />
	<div>
		<div class="flex flex-wrap items-center gap-3">
			<Button.Root
				class="rounded-md bg-brand px-4 py-2 text-sm font-medium text-on-brand hover:bg-brand-strong disabled:opacity-60"
				disabled={creating}
				type="submit"
			>
				{creating ? 'Čuvanje…' : itemType ? 'Sačuvaj izmene' : 'Dodaj tip'}
			</Button.Root>
			{#if !itemType && oncancel}
				<Button.Root
					class="rounded-md border border-line bg-surface px-4 py-2 text-sm font-medium text-ink hover:bg-soft disabled:opacity-60"
					disabled={creating}
					type="button"
					onclick={oncancel}
				>
					Otkaži dodavanje
				</Button.Root>
			{/if}
		</div>
		{#if error}<p class="mt-3 text-sm text-danger" role="alert">{error}</p>{/if}
	</div>
</form>
