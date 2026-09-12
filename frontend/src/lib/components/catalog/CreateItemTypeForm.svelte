<script lang="ts">
	import { flip } from 'svelte/animate';
	import { onMount } from 'svelte';
	import X from '@lucide/svelte/icons/x';
	import {
		api,
		ApiError,
		type ItemType,
		type PropertyOption,
		type PropertyValue,
		type PropertyVisibility
	} from '$lib/api';
	import ExpiringSoonDaysInput from '$lib/components/catalog/ExpiringSoonDaysInput.svelte';
	import ItemPropertyValueInput from '$lib/components/inventory/ItemPropertyValueInput.svelte';
	import MultiOptionCombobox from '$lib/components/shared/MultiOptionCombobox.svelte';
	import {
		defaultExpiringSoonDays,
		defaultJsonValue,
		propertyValueTypeLabel
	} from '$lib/domain/items';
	import { requiredTextError } from '$lib/domain/form-validation';
	import { Button, Label, Portal, Separator, Tabs } from 'bits-ui';
	import { dndzone, setAriaStrings, type DndEvent } from 'svelte-dnd-action';
	import { toast } from 'svelte-sonner';

	const steps = [
		{ value: 'basic', label: 'Osnovno' },
		{ value: 'properties', label: 'Svojstva' },
		{ value: 'display', label: 'Prikaz' }
	] as const;
	const reorderAnimationDuration = 150;

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
	let expiringSoonDays = $state(defaultExpiringSoonDays);
	let comboboxSelectedPropertyValues = $state<string[]>([]);
	let selectedPropertyIDs = $state<number[]>([]);
	let orderedProperties = $state<PropertyOption[]>([]);
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
	let hasExpiryProperty = $derived(
		selectedProperties.some((property) => property.value_type === 'expiry')
	);
	let originalProperties = $derived(
		new Map(itemType?.properties.map((property) => [property.id, property]))
	);
	let propertyOrderChanged = $derived(
		itemType
			? itemType.properties.length !== selectedPropertyIDs.length ||
					itemType.properties.some((property, index) => property.id !== selectedPropertyIDs[index])
			: false
	);

	$effect(() => {
		if (initializedItemTypeID === itemType?.id) return;

		initializedItemTypeID = itemType?.id;
		name = itemType?.name ?? '';
		description = itemType?.description ?? '';
		derivedNameFormat = itemType?.derived_name_format ?? '';
		// A type with no window stored has none by intent - it is what the backend writes when the
		// count is cleared - so it opens at nought rather than back at the default.
		expiringSoonDays = itemType ? (itemType.expiring_soon_days ?? 0) : defaultExpiringSoonDays;
		comboboxSelectedPropertyValues =
			itemType?.properties.map((property) => String(property.id)) ?? [];
		selectedPropertyIDs = itemType?.properties.map((property) => property.id) ?? [];
		orderedProperties = selectedPropertyIDs.flatMap((propertyID) => {
			const property = propertiesByID.get(propertyID);
			return property ? [property] : [];
		});
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
		comboboxSelectedPropertyValues = values;
		const nextPropertyIDs = values.map(Number);
		selectedPropertyIDs = [
			...selectedPropertyIDs.filter((propertyID) => nextPropertyIDs.includes(propertyID)),
			...nextPropertyIDs.filter((propertyID) => !selectedPropertyIDs.includes(propertyID))
		];
		orderedProperties = selectedPropertyIDs.flatMap((propertyID) => {
			const property = propertiesByID.get(propertyID);
			return property ? [property] : [];
		});

		for (const propertyID of selectedPropertyIDs) {
			if (defaultValues[propertyID] !== undefined) continue;

			const property = propertiesByID.get(propertyID);
			if (property) {
				defaultValues = {
					...defaultValues,
					[propertyID]: defaultJsonValue(property.value_type, property.default_value)
				};
			}
		}

		for (const propertyID of selectedPropertyIDs) {
			if (visibilities[propertyID] !== undefined) continue;
			visibilities = { ...visibilities, [propertyID]: 'overview' };
		}
	}

	function markDefaultEdited(id: number) {
		editedDefaultPropertyIDs = new Set([...editedDefaultPropertyIDs, id]);
	}

	function removeSelectedProperty(propertyID: number) {
		updateSelectedProperties(
			comboboxSelectedPropertyValues.filter(
				(selectedPropertyID) => Number(selectedPropertyID) !== propertyID
			)
		);
	}

	function reorderSelectedProperties(event: CustomEvent<DndEvent<PropertyOption>>) {
		orderedProperties = event.detail.items;

		if (event.detail.items.every((property) => typeof property.id === 'number')) {
			selectedPropertyIDs = event.detail.items.map((property) => property.id);
		}
	}

	function listenForPropertyReorder(node: HTMLElement) {
		const handleReorder = (event: Event) =>
			reorderSelectedProperties(event as CustomEvent<DndEvent<PropertyOption>>);

		node.addEventListener('consider', handleReorder);
		node.addEventListener('finalize', handleReorder);

		return {
			destroy() {
				node.removeEventListener('consider', handleReorder);
				node.removeEventListener('finalize', handleReorder);
			}
		};
	}

	function keepPillInteraction(event: MouseEvent | TouchEvent) {
		event.stopPropagation();
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
					expiring_soon_days: expiringSoonDays,
					properties: selectedPropertyIDs.map((id) => ({
						id,
						default_value: defaultValues[id],
						visibility: visibilities[id] ?? 'overview'
					}))
				});
			} else {
				await api.updateItemType(itemType.id, {
					name,
					description,
					expiring_soon_days: expiringSoonDays
				});

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
				if (propertyOrderChanged) {
					await api.reorderItemTypeProperties(itemType.id, selectedPropertyIDs);
				}
			}

			toast.success(itemType ? 'Tip stavke je izmenjen.' : 'Tip stavke je dodat.');
			onsaved();
		} catch (reason) {
			toast.error(reason instanceof ApiError ? reason.message : 'Tip stavke nije sačuvan.');
		} finally {
			creating = false;
		}
	}

	onMount(() => {
		setAriaStrings({
			dragStarted: ({ itemLabel, zoneLabel }) =>
				`Početo je premeštanje svojstva ${itemLabel} u listi ${zoneLabel}.`,
			movedToPosition: ({ itemLabel, zoneLabel, position, count }) =>
				`Svojstvo ${itemLabel} je na poziciji ${position} od ${count} u listi ${zoneLabel}.`,
			movedToZoneEnd: ({ itemLabel, zoneLabel }) =>
				`Svojstvo ${itemLabel} je premešteno na kraj liste ${zoneLabel}.`,
			movedToZoneStart: ({ itemLabel, zoneLabel }) =>
				`Svojstvo ${itemLabel} je premešteno na početak liste ${zoneLabel}.`,
			dropped: ({ itemLabel, zoneLabel, position, count }) =>
				`Svojstvo ${itemLabel} je postavljeno na poziciju ${position} od ${count} u listi ${zoneLabel}.`,
			zoneActiveInstruction:
				'Pritisnite razmak ili Enter da započnete premeštanje. Strelicama promenite položaj, a zatim razmakom, Enterom ili Escape tasterom završite.'
		});

		return () => setAriaStrings(null);
	});
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
								bind:values={comboboxSelectedPropertyValues}
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
						<ul
							class="mt-1 flex flex-wrap gap-2"
							aria-label="Izabrana svojstva"
							use:listenForPropertyReorder
							use:dndzone={{
								items: orderedProperties,
								type: 'item-type-properties',
								flipDurationMs: reorderAnimationDuration,
								dropTargetStyle: {},
								dragDisabled: creating
							}}
						>
							{#each orderedProperties as property (property.id)}
								<li
									animate:flip={{ duration: reorderAnimationDuration }}
									aria-label={`${property.name}${visibilities[property.id] === 'overview' ? ', prikazuje se direktno u tabeli' : ''}`}
									class={`flex min-w-0 cursor-grab items-center rounded-md border bg-soft transition-colors hover:border-brand hover:bg-brand-soft active:cursor-grabbing ${activePropertyID === property.id ? 'border-brand ring-1 ring-brand' : 'border-line'}`}
								>
									<span
										class="flex min-w-0 items-center gap-1.5 py-1.5 pr-1.5 pl-2.5 text-sm text-ink"
										role="button"
										tabindex="0"
										aria-pressed={activePropertyID === property.id}
										onclick={() => (activePropertyID = property.id)}
										onkeydown={(event) => {
											if (event.key !== 'Enter' && event.key !== ' ') return;
											event.preventDefault();
											event.stopPropagation();
											activePropertyID = property.id;
										}}
									>
										<span
											class={`max-w-48 truncate ${visibilities[property.id] === 'overview' ? '' : 'text-ink/70'}`}
											>{property.name}</span
										>
									</span>
									<button
										class="mr-1 grid size-6 shrink-0 cursor-pointer place-items-center rounded text-muted hover:bg-surface hover:text-ink"
										type="button"
										onmousedown={keepPillInteraction}
										ontouchstart={keepPillInteraction}
										onclick={(event) => {
											event.stopPropagation();
											removeSelectedProperty(property.id);
										}}
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
										class="inline-flex cursor-pointer items-center gap-1.5 rounded-md border border-line bg-soft px-2.5 py-1.5 text-sm text-ink transition-colors hover:border-brand hover:bg-brand-soft focus-visible:outline-1 focus-visible:outline-brand"
										type="button"
										onclick={() => appendPropertyToken(property.name)}
										aria-label={`Dodaj ${property.name} u format izvedenog naziva${visibilities[property.id] === 'overview' ? ', prikazuje se direktno u tabeli' : ''}`}
									>
										<span class={visibilities[property.id] === 'overview' ? '' : 'text-ink/70'}
											>{property.name}</span
										>
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
				<!-- Only types that record an expiry have anything to warn about, and which properties
				     they record was settled on the step before. -->
				{#if hasExpiryProperty}
					<Separator.Root class="h-px bg-line" decorative />
					<div>
						<Label.Root class="text-sm font-medium text-ink" for="item-type-expiring-soon-days">
							Upozorenje pred istek roka
						</Label.Root>
						<p class="mt-1 text-sm text-muted">
							Koliko dana pre isteka roka stavke ovog tipa dobijaju žutu oznaku.
						</p>
						<ExpiringSoonDaysInput id="item-type-expiring-soon-days" bind:days={expiringSoonDays} />
					</div>
				{/if}
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
