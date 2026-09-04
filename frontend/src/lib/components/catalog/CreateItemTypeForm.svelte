<script lang="ts">
	import { api, ApiError, type ItemType, type PropertyOption, type PropertyValue } from '$lib/api';
	import ItemPropertyValueInput from '$lib/components/inventory/ItemPropertyValueInput.svelte';
	import MultiOptionCombobox from '$lib/components/shared/MultiOptionCombobox.svelte';
	import { defaultJsonValue, propertyValueTypeLabel } from '$lib/domain/items';
	import { requiredTextError } from '$lib/domain/form-validation';
	import { Button, Label, ScrollArea, Separator } from 'bits-ui';
	import { toast } from 'svelte-sonner';

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
	let derivedNameFormat = $state('');
	let selectedPropertyValues = $state<string[]>([]);
	let selectedPropertyIDs = $derived(selectedPropertyValues.map(Number));
	let defaultValues = $state<Record<number, PropertyValue>>({});
	let editedDefaultPropertyIDs = $state<Set<number>>(new Set());
	let creating = $state(false);
	let fieldErrors = $state<{ name?: string; derivedNameFormat?: string }>({});
	let initializedItemTypeID = $state<number | undefined>(undefined);
	let propertiesByID = $derived(new Map(properties.map((property) => [property.id, property])));
	let selectedProperties = $derived(
		selectedPropertyIDs.flatMap((propertyID) => {
			const property = propertiesByID.get(propertyID);
			return property ? [property] : [];
		})
	);
	let originalProperties = $derived(
		new Map(itemType?.properties.map((property) => [property.id, property]))
	);

	$effect(() => {
		if (initializedItemTypeID === itemType?.id) return;

		initializedItemTypeID = itemType?.id;
		name = itemType?.name ?? '';
		description = itemType?.description ?? '';
		derivedNameFormat = itemType?.derived_name_format ?? '';
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

	function clearFieldError(field: 'name' | 'derivedNameFormat') {
		if (!fieldErrors[field]) return;
		fieldErrors = { ...fieldErrors, [field]: undefined };
	}

	function validate() {
		fieldErrors = {
			name: requiredTextError(name, 'naziv tipa stavke'),
			derivedNameFormat: requiredTextError(derivedNameFormat, 'format izvedenog naziva')
		};
		return Object.values(fieldErrors).every((fieldError) => fieldError === undefined);
	}

	async function createItemType(event: SubmitEvent) {
		event.preventDefault();
		if (!validate()) return;
		creating = true;

		try {
			if (!itemType) {
				await api.createItemType({
					name,
					description,
					derived_name_format: derivedNameFormat,
					properties: selectedPropertyIDs.map((id) => ({
						id,
						default_value: defaultValues[id],
						visibility: 'overview'
					}))
				});
			} else {
				await api.updateItemType(itemType.id, { name, description });

				const defaultValueUpdates: Promise<unknown>[] = [];
				for (const propertyID of selectedPropertyIDs) {
					if (!originalProperties.has(propertyID)) {
						// Additions are intentionally sequential: their request order becomes their position.
						await api.addItemTypeProperty(itemType.id, {
							property_id: propertyID,
							default_value: defaultValues[propertyID]
						});
					} else if (editedDefaultPropertyIDs.has(propertyID)) {
						defaultValueUpdates.push(
							api.updateItemTypeProperty(itemType.id, propertyID, {
								default_value: defaultValues[propertyID]
							})
						);
					}
				}
				await Promise.all(defaultValueUpdates);
				await api.updateItemType(itemType.id, { derived_name_format: derivedNameFormat });

				const removals = itemType.properties
					.filter((property) => !selectedPropertyIDs.includes(property.id))
					.map((property) => api.removeItemTypeProperty(itemType.id, property.id));
				await Promise.all(removals);
			}

			toast.success(itemType ? 'Tip stavke je izmenjen.' : 'Tip stavke je dodat.');
			onsaved();
		} catch (reason) {
			toast.error(reason instanceof ApiError ? reason.message : 'Tip stavke nije sačuvan.');
		} finally {
			creating = false;
		}
	}
</script>

<form class="grid gap-4" novalidate onsubmit={createItemType}>
	<div>
		<Label.Root class="text-sm font-medium text-ink" for="item-type-name">Naziv</Label.Root>
		<input
			id="item-type-name"
			class={`mt-1 block w-full ${fieldErrors.name ? 'field-invalid' : ''}`}
			bind:value={name}
			aria-invalid={Boolean(fieldErrors.name)}
			aria-describedby={fieldErrors.name ? 'item-type-name-error' : undefined}
			oninput={() => clearFieldError('name')}
		/>
		{#if fieldErrors.name}
			<p id="item-type-name-error" class="mt-1 text-xs text-danger" role="alert">
				{fieldErrors.name}
			</p>
		{/if}
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
		<div class="mt-4">
			<Label.Root class="text-sm font-medium text-ink" for="item-type-derived-name">
				Izvedeni naziv stavke
			</Label.Root>
			<textarea
				id="item-type-derived-name"
				class="mt-1 block min-h-20 w-full font-mono text-sm"
				bind:value={derivedNameFormat}
				placeholder={'{Naziv stavke} · {Proizvođač}'}
				class:field-invalid={Boolean(fieldErrors.derivedNameFormat)}
				aria-invalid={Boolean(fieldErrors.derivedNameFormat)}
				aria-describedby={fieldErrors.derivedNameFormat
					? 'item-type-derived-name-error'
					: undefined}
				oninput={() => clearFieldError('derivedNameFormat')}></textarea>
			<p class="mt-1 text-xs text-muted">
				Unesite nazive svojstava u vitičastim zagradama, npr. {'{Naziv stavke}'}. Možete dodati i
				običan tekst.
			</p>
			{#if fieldErrors.derivedNameFormat}
				<p id="item-type-derived-name-error" class="mt-1 text-xs text-danger" role="alert">
					{fieldErrors.derivedNameFormat}
				</p>
			{/if}
		</div>
		{#if selectedProperties.length}
			<p class="mt-3 text-xs font-medium text-muted">Podrazumevane vrednosti</p>
			<ScrollArea.Root class="mt-2 h-48 overflow-hidden rounded-md border border-line" type="auto">
				<ScrollArea.Viewport class="h-full w-full">
					<div class="divide-y divide-line">
						{#each selectedProperties as property (property.id)}
							<div
								class="grid gap-2 px-3 py-2 sm:grid-cols-[minmax(0,1fr)_minmax(12rem,0.8fr)] sm:items-center sm:gap-4"
							>
								<div class="min-w-0">
									<p class="truncate text-sm font-medium text-ink">{property.name}</p>
									<p class="text-xs text-muted">{propertyValueTypeLabel(property.value_type)}</p>
								</div>
								<Label.Root class="sr-only" for={`item-type-property-${property.id}`}>
									Podrazumevana vrednost za {property.name}
								</Label.Root>
								<ItemPropertyValueInput
									id={`item-type-property-${property.id}`}
									bind:value={defaultValues[property.id]}
									{property}
									className=""
									compact
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
	</div>
</form>
