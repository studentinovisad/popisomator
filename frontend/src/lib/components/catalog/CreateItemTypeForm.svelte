<script lang="ts">
	import X from '@lucide/svelte/icons/x';
	import {
		api,
		ApiError,
		type ItemType,
		type PropertyOption,
		type PropertyValue,
		type PropertyVisibility
	} from '$lib/api';
	import ItemPropertyValueInput from '$lib/components/inventory/ItemPropertyValueInput.svelte';
	import MultiOptionCombobox from '$lib/components/shared/MultiOptionCombobox.svelte';
	import { defaultJsonValue, propertyValueTypeLabel } from '$lib/domain/items';
	import { requiredTextError } from '$lib/domain/form-validation';
	import { Button, Label, Portal, Separator, Tabs } from 'bits-ui';
	import { toast } from 'svelte-sonner';

	const steps = [
		{ value: 'basic', label: 'Osnovno' },
		{ value: 'properties', label: 'Svojstva' },
		{ value: 'display', label: 'Prikaz' }
	] as const;

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
	let visibilities = $state<Record<number, PropertyVisibility>>({});
	let editedDefaultPropertyIDs = $state<Set<number>>(new Set());
	let creating = $state(false);
	let fieldErrors = $state<{ name?: string; derivedNameFormat?: string }>({});
	let activeStep = $state<(typeof steps)[number]['value']>('basic');
	let activePropertyID = $state<number | undefined>(undefined);
	let initializedItemTypeID = $state<number | undefined>(undefined);
	let propertiesByID = $derived(new Map(properties.map((property) => [property.id, property])));
	let selectedProperties = $derived(
		selectedPropertyIDs.flatMap((propertyID) => {
			const property = propertiesByID.get(propertyID);
			return property ? [property] : [];
		})
	);
	let activeProperty = $derived(
		selectedProperties.find((property) => property.id === activePropertyID) ?? null
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
		visibilities = Object.fromEntries(
			(itemType?.properties ?? []).map((property) => [property.id, property.visibility])
		);
		editedDefaultPropertyIDs = new Set();
	});

	$effect(() => {
		if (selectedPropertyIDs.length === 0) {
			activePropertyID = undefined;
			return;
		}

		if (!activePropertyID || !selectedPropertyIDs.includes(activePropertyID)) {
			activePropertyID = selectedPropertyIDs[0];
		}
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

		for (const propertyID of values.map(Number)) {
			if (visibilities[propertyID] !== undefined) continue;
			visibilities = { ...visibilities, [propertyID]: 'overview' };
		}
	}

	function markDefaultEdited(id: number) {
		editedDefaultPropertyIDs = new Set([...editedDefaultPropertyIDs, id]);
	}

	function removeSelectedProperty(propertyID: number) {
		updateSelectedProperties(
			selectedPropertyValues.filter(
				(selectedPropertyID) => Number(selectedPropertyID) !== propertyID
			)
		);
	}

	function setVisibility(propertyID: number, visibility: PropertyVisibility) {
		visibilities = { ...visibilities, [propertyID]: visibility };
	}

	function toggleOverview(propertyID: number, checked: boolean) {
		setVisibility(propertyID, checked ? 'overview' : 'details');
	}

	function clearFieldError(field: 'name' | 'derivedNameFormat') {
		if (!fieldErrors[field]) return;
		fieldErrors = { ...fieldErrors, [field]: undefined };
	}

	function appendPropertyToken(propertyName: string) {
		const token = `{${propertyName}}`;
		derivedNameFormat = derivedNameFormat.trim() ? `${derivedNameFormat} · ${token}` : token;
		clearFieldError('derivedNameFormat');
	}

	function validate() {
		fieldErrors = {
			name: requiredTextError(name, 'naziv tipa stavke'),
			derivedNameFormat: requiredTextError(derivedNameFormat, 'format izvedenog naziva')
		};

		if (fieldErrors.name) activeStep = 'basic';
		else if (fieldErrors.derivedNameFormat) activeStep = 'display';

		return Object.values(fieldErrors).every((fieldError) => fieldError === undefined);
	}

	function previousStep() {
		const index = steps.findIndex((step) => step.value === activeStep);
		activeStep = steps[Math.max(0, index - 1)].value;
	}

	function nextStep() {
		const index = steps.findIndex((step) => step.value === activeStep);
		activeStep = steps[Math.min(steps.length - 1, index + 1)].value;
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
						visibility: visibilities[id] ?? 'overview'
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
							default_value: defaultValues[propertyID],
							visibility: visibilities[propertyID] ?? 'overview'
						});
					} else {
						const originalProperty = originalProperties.get(propertyID);
						const visibilityChanged = originalProperty?.visibility !== visibilities[propertyID];
						if (!editedDefaultPropertyIDs.has(propertyID) && !visibilityChanged) continue;

						defaultValueUpdates.push(
							api.updateItemTypeProperty(itemType.id, propertyID, {
								...(editedDefaultPropertyIDs.has(propertyID)
									? { default_value: defaultValues[propertyID] }
									: {}),
								...(visibilityChanged ? { visibility: visibilities[propertyID] } : {})
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

<form id="item-type-form" class="flex min-h-0 flex-1 flex-col" novalidate onsubmit={createItemType}>
	<Tabs.Root bind:value={activeStep} activationMode="manual" class="flex min-h-0 flex-1 flex-col">
		<Tabs.List
			class="grid grid-cols-3 gap-1 rounded-lg border border-line bg-soft p-1"
			aria-label="Koraci za tip stavke"
		>
			{#each steps as step, index (step.value)}
				<Tabs.Trigger
					value={step.value}
					class="flex min-w-0 items-center justify-center gap-2 rounded-md px-2 py-2.5 text-sm text-muted transition-colors hover:bg-surface hover:text-ink data-[state=active]:bg-surface data-[state=active]:font-medium data-[state=active]:text-ink data-[state=active]:shadow-sm"
				>
					<span
						class={`grid size-5 shrink-0 place-items-center rounded-full border font-mono text-[0.6875rem] ${activeStep === step.value ? 'border-brand bg-brand text-on-brand' : 'border-line'}`}
					>
						{index + 1}
					</span>
					<span class="truncate">{step.label}</span>
				</Tabs.Trigger>
			{/each}
		</Tabs.List>

		<Tabs.Content value="basic" class="min-h-0 flex-1 pt-6">
			<div class="grid gap-5">
				<div>
					<h2 class="text-base font-medium text-ink">Osnovni podaci</h2>
					<p class="mt-1 text-sm text-muted">Naziv i kratak opis će pomoći pri radu sa stavkama.</p>
				</div>
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
					<Label.Root class="text-sm font-medium text-ink" for="item-type-description"
						>Opis</Label.Root
					>
					<textarea
						id="item-type-description"
						class="mt-1 block min-h-28 w-full"
						bind:value={description}></textarea>
				</div>
			</div>
		</Tabs.Content>

		<Tabs.Content value="properties" class="min-h-0 flex-1 pt-6">
			<div class="flex min-h-0 flex-1 flex-col">
				<div>
					<h2 class="text-base font-medium text-ink">Svojstva</h2>
					<p class="mt-1 text-sm text-muted">
						Dodajte svojstva, zatim za svako podesite podrazumevanu vrednost i prikaz na pregledu.
					</p>
					<div class="mt-5">
						<Label.Root class="text-sm font-medium text-ink" for="item-type-properties">
							Dodaj svojstvo
						</Label.Root>
						<div class="mt-1">
							<MultiOptionCombobox
								id="item-type-properties"
								options={properties}
								bind:values={selectedPropertyValues}
								placeholder="Pretražite i dodajte svojstvo"
								emptyMessage="Nema odgovarajućih svojstava."
								showSelected={false}
								onvaluechange={updateSelectedProperties}
							/>
						</div>
					</div>
				</div>

				{#if selectedProperties.length}
					<Separator.Root class="mt-6 h-px bg-line" decorative />
					<div class="mt-4 min-h-0 flex-1">
						<div class="flex items-center justify-between gap-3">
							<p class="text-sm font-medium text-ink">Izabrana svojstva</p>
							<span class="text-xs text-muted">{selectedProperties.length}</span>
						</div>
						<ul class="mt-1 flex flex-wrap gap-2" aria-label="Izabrana svojstva">
							{#each selectedProperties as property, index (property.id)}
								<li
									class={`flex min-w-0 items-center rounded-md border ${activePropertyID === property.id ? 'border-brand bg-brand-soft' : 'border-line bg-soft'}`}
								>
									<button
										class="flex min-w-0 cursor-pointer items-center gap-1.5 py-1.5 pr-1.5 pl-2.5 text-sm text-ink focus-visible:outline-1 focus-visible:outline-brand"
										type="button"
										aria-pressed={activePropertyID === property.id}
										onclick={() => (activePropertyID = property.id)}
									>
										<span class="font-mono text-xs text-muted">{index + 1}</span>
										<span class="max-w-48 truncate">{property.name}</span>
									</button>
									<button
										class="mr-1 grid size-6 shrink-0 cursor-pointer place-items-center rounded text-muted hover:bg-surface hover:text-ink"
										type="button"
										onclick={() => removeSelectedProperty(property.id)}
										aria-label={`Ukloni ${property.name}`}
									>
										<X class="size-3.5" aria-hidden="true" />
									</button>
								</li>
							{/each}
						</ul>

						{#if activeProperty}
							<Separator.Root class="mt-6 h-px bg-line" decorative />
							<p class="mt-4 text-sm font-medium text-ink">Podešavanja svojstva</p>
							<section class="mt-1 rounded-md border border-line bg-surface p-4">
								<div class="flex flex-wrap items-start justify-between gap-3">
									<div>
										<h3 class="text-sm font-medium text-ink">{activeProperty.name}</h3>
										<p class="mt-0.5 text-xs text-muted">
											{propertyValueTypeLabel(activeProperty.value_type)}
										</p>
									</div>
									<label class="flex cursor-pointer items-center gap-2 text-sm text-ink">
										<input
											type="checkbox"
											checked={visibilities[activeProperty.id] !== 'details'}
											onchange={(event) =>
												toggleOverview(
													activeProperty.id,
													(event.currentTarget as HTMLInputElement).checked
												)}
										/>
										Prikaži direkt u tabeli
									</label>
								</div>
								<div class="mt-4">
									<Label.Root
										class="text-sm font-medium text-ink"
										for={`item-type-property-${activeProperty.id}`}
									>
										Podrazumevana vrednost
									</Label.Root>
									<ItemPropertyValueInput
										id={`item-type-property-${activeProperty.id}`}
										bind:value={defaultValues[activeProperty.id]}
										property={activeProperty}
										className="mt-1"
										onvaluechange={() => markDefaultEdited(activeProperty.id)}
									/>
								</div>
							</section>
						{/if}
					</div>
				{:else if properties.length === 0}
					<p class="mt-5 text-sm text-muted">Najpre dodajte svojstvo.</p>
				{/if}
			</div>
		</Tabs.Content>

		<Tabs.Content value="display" class="min-h-0 flex-1 pt-6">
			<div class="grid gap-5">
				<div>
					<h2 class="text-base font-medium text-ink">Izvedeni naziv</h2>
					<p class="mt-1 text-sm text-muted">
						Sastavite naziv po kom će se stavke prepoznavati u katalogu i na pregledu.
					</p>
				</div>
				<div>
					<p class="text-sm font-medium text-ink">Izabrana svojstva</p>
					{#if selectedProperties.length}
						<ul class="mt-1 flex flex-wrap gap-2" aria-label="Izabrana svojstva za izvedeni naziv">
							{#each selectedProperties as property (property.id)}
								<li>
									<button
										class="cursor-pointer rounded-md border border-line bg-soft px-2.5 py-1.5 text-sm text-ink transition-colors hover:border-brand hover:bg-brand-soft focus-visible:outline-1 focus-visible:outline-brand"
										type="button"
										onclick={() => appendPropertyToken(property.name)}
										aria-label={`Dodaj ${property.name} u format izvedenog naziva`}
									>
										{property.name}
									</button>
								</li>
							{/each}
						</ul>
					{:else}
						<p class="mt-1 text-sm text-muted">Najpre dodajte svojstva na prethodnom koraku.</p>
					{/if}
				</div>
				<Separator.Root class="h-px bg-line" decorative />
				<div>
					<Label.Root class="text-sm font-medium text-ink" for="item-type-derived-name">
						Format izvedenog naziva
					</Label.Root>
					<textarea
						id="item-type-derived-name"
						class={`mt-1 block min-h-28 w-full font-mono text-sm ${fieldErrors.derivedNameFormat ? 'field-invalid' : ''}`}
						bind:value={derivedNameFormat}
						placeholder={'{Naziv stavke} · {Proizvođač}'}
						aria-invalid={Boolean(fieldErrors.derivedNameFormat)}
						aria-describedby={fieldErrors.derivedNameFormat
							? 'item-type-derived-name-error'
							: undefined}
						oninput={() => clearFieldError('derivedNameFormat')}></textarea>
					<p class="mt-2 text-sm text-muted">
						Koristite nazive izabranih svojstava u vitičastim zagradama. Možete dodati i običan
						tekst.
					</p>
					{#if fieldErrors.derivedNameFormat}
						<p id="item-type-derived-name-error" class="mt-1 text-xs text-danger" role="alert">
							{fieldErrors.derivedNameFormat}
						</p>
					{/if}
				</div>
			</div>
		</Tabs.Content>
	</Tabs.Root>

	<Portal to="#page-footer-actions">
		<footer class="flex h-full items-center justify-between gap-3">
			<div class="flex items-center gap-3">
				<Button.Root
					class="rounded-md border border-line bg-surface px-4 py-2 text-sm font-medium text-ink hover:bg-soft disabled:opacity-50"
					disabled={activeStep === 'basic' || creating}
					type="button"
					onclick={previousStep}
				>
					Nazad
				</Button.Root>
				{#if activeStep !== 'display'}
					<Button.Root
						class="rounded-md bg-brand px-4 py-2 text-sm font-medium text-on-brand hover:bg-brand-strong"
						type="button"
						onclick={nextStep}
					>
						Dalje
					</Button.Root>
				{:else}
					<Button.Root
						class="rounded-md bg-brand px-4 py-2 text-sm font-medium text-on-brand hover:bg-brand-strong disabled:opacity-60"
						disabled={creating}
						form="item-type-form"
						type="submit"
					>
						{creating ? 'Čuvanje…' : itemType ? 'Sačuvaj izmene' : 'Dodaj tip'}
					</Button.Root>
				{/if}
			</div>
			{#if !itemType && oncancel}
				<Button.Root
					class="rounded-md px-3 py-2 text-sm font-medium text-muted hover:bg-soft hover:text-ink"
					disabled={creating}
					type="button"
					onclick={oncancel}
				>
					Otkaži dodavanje
				</Button.Root>
			{/if}
		</footer>
	</Portal>
</form>
