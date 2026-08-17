<script lang="ts">
	import { api, ApiError, type ItemType, type Property } from '$lib/api';
	import ItemPropertyValueInput from '$lib/components/inventory/ItemPropertyValueInput.svelte';
	import { defaultJsonValue, propertyValueTypeLabel } from '$lib/domain/items';
	import { Button, Label, ScrollArea, Separator } from 'bits-ui';

	let {
		itemType,
		properties,
		onsaved,
		oncancel
	}: {
		itemType?: ItemType;
		properties: Property[];
		onsaved: () => void;
		oncancel?: () => void;
	} = $props();

	let name = $state('');
	let description = $state('');
	let selectedPropertyIDs = $state<number[]>([]);
	let defaultValues = $state<Record<number, string>>({});
	let editedDefaultPropertyIDs = $state<Set<number>>(new Set());
	let creating = $state(false);
	let error = $state('');
	let initializedItemTypeID = $state<number | undefined>(undefined);
	let propertiesByID = $derived(new Map(properties.map((property) => [property.id, property])));
	let originalProperties = $derived(
		new Map(itemType?.properties.map((property) => [property.id, property]))
	);

	$effect(() => {
		if (initializedItemTypeID === itemType?.id) return;

		initializedItemTypeID = itemType?.id;
		name = itemType?.name ?? '';
		description = itemType?.description ?? '';
		selectedPropertyIDs = itemType?.properties.map((property) => property.id) ?? [];
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

	function selectProperty(id: number, checked: boolean) {
		selectedPropertyIDs = checked
			? [...selectedPropertyIDs, id]
			: selectedPropertyIDs.filter((propertyID) => propertyID !== id);

		if (checked && defaultValues[id] === undefined) {
			const property = propertiesByID.get(id);
			if (property) {
				defaultValues = {
					...defaultValues,
					[id]: defaultJsonValue(property.value_type, property.default_value)
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
						visibility: "overview"
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
							api.addItemTypeProperty(itemType.id, {property_id: propertyID, default_value: defaultValues[propertyID]})
						);
					} else if (editedDefaultPropertyIDs.has(propertyID)) {
						changes.push(
							api.updateItemTypeProperty(itemType.id, propertyID, {default_value: defaultValues[propertyID]})
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
		<ScrollArea.Root class="mt-3 h-64 overflow-hidden rounded-md border border-line" type="auto">
			<ScrollArea.Viewport class="h-full w-full">
				<div class="divide-y divide-line">
					{#each properties as property (property.id)}
						<div class="px-3 py-2">
							<label class="flex cursor-pointer items-center justify-between gap-4">
								<span>
									<span class="block text-sm font-medium text-ink">{property.name}</span>
									<span class="block text-xs text-muted"
										>{propertyValueTypeLabel(property.value_type)}</span
									>
								</span>
								<input
									type="checkbox"
									checked={selectedPropertyIDs.includes(property.id)}
									onchange={(event) => selectProperty(property.id, event.currentTarget.checked)}
								/>
							</label>
							{#if selectedPropertyIDs.includes(property.id)}
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
							{/if}
						</div>
					{/each}
					{#if properties.length === 0}<p class="px-3 py-3 text-sm text-muted">
							Najpre dodajte svojstvo.
						</p>{/if}
				</div>
			</ScrollArea.Viewport>
			<ScrollArea.Scrollbar class="flex w-2.5 touch-none bg-soft p-0.5" orientation="vertical">
				<ScrollArea.Thumb class="flex-1 rounded-full bg-line" />
			</ScrollArea.Scrollbar>
		</ScrollArea.Root>
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
