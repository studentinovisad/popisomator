<script lang="ts">
	import {
		api,
		ApiError,
		type ItemType,
		type ItemTypeOption,
		type PropertyOption,
		type PropertyValue
	} from '$lib/api';
	import ItemPropertyValueInput from '$lib/components/inventory/ItemPropertyValueInput.svelte';
	import NumberInput from '$lib/components/shared/NumberInput.svelte';
	import OptionCombobox from '$lib/components/shared/OptionCombobox.svelte';
	import FormAlert from '$lib/components/shared/FormAlert.svelte';
	import { propertyValueError, requiredTextError } from '$lib/domain/form-validation';
	import { Button, Label, Portal, Progress, Tabs } from 'bits-ui';
	import { toast } from 'svelte-sonner';

	const steps = [
		{ value: 'setup', label: 'Tip i količina' },
		{ value: 'item', label: 'Stavka' }
	] as const;

	let {
		itemTypes,
		properties,
		oncreated,
		oncancel
	}: {
		itemTypes: ItemTypeOption[];
		properties: PropertyOption[];
		oncreated: () => void;
		oncancel?: () => void;
	} = $props();

	let selectedTypeID = $state('');
	let amount = $state(1);
	let propertyValues = $state<Record<number, PropertyValue | null>>({});
	let selectedPropertyIDs = $state<number[]>([]);
	let fieldErrors = $state<{ type?: string; amount?: string }>({});
	let propertyErrors = $state<Record<number, string>>({});
	let creating = $state(false);
	let loadingType = $state(false);
	let selectedType = $state<ItemType | null>(null);
	let activeStep = $state<(typeof steps)[number]['value']>('setup');
	let typeLoadVersion = 0;

	let propertiesByID = $derived(new Map(properties.map((property) => [property.id, property])));

	async function selectType(value: string) {
		selectedTypeID = value;
		if (fieldErrors.type) fieldErrors = { ...fieldErrors, type: undefined };
		selectedType = null;
		propertyValues = {};
		selectedPropertyIDs = [];
		propertyErrors = {};

		const version = ++typeLoadVersion;
		const typeID = Number(value);
		if (!Number.isSafeInteger(typeID)) {
			loadingType = false;
			activeStep = 'setup';
			return;
		}

		loadingType = true;
		try {
			const itemType = await api.getItemType(typeID);
			if (version !== typeLoadVersion) return;

			selectedType = itemType;
			selectedPropertyIDs = itemType.properties.map((itemProperty) => itemProperty.id);
			propertyValues = Object.fromEntries(
				itemType.properties.flatMap((itemProperty) => {
					const property = propertiesByID.get(itemProperty.id);
					return property ? [[itemProperty.id, itemProperty.default_value]] : [];
				})
			);
		} catch (reason) {
			if (version !== typeLoadVersion) return;
			toast.error(reason instanceof ApiError ? reason.message : 'Tip stavke nije učitan.');
		} finally {
			if (version === typeLoadVersion) loadingType = false;
		}
	}

	function validateSetup() {
		const parsedAmount = Number(amount);
		fieldErrors = {
			type: requiredTextError(selectedTypeID, 'tip stavke'),
			amount:
				Number.isInteger(parsedAmount) && parsedAmount >= 1 && parsedAmount <= 100
					? undefined
					: 'Količina mora biti ceo broj od 1 do 100.'
		};

		return Object.values(fieldErrors).every((fieldError) => fieldError === undefined);
	}

	function validateProperties() {
		const nextPropertyErrors: Record<number, string> = {};
		for (const propertyID of selectedPropertyIDs) {
			const property = propertiesByID.get(propertyID);
			if (!property) continue;
			const message = propertyValueError(property, propertyValues[propertyID]);
			if (message) nextPropertyErrors[propertyID] = message;
		}
		propertyErrors = nextPropertyErrors;

		return Object.keys(nextPropertyErrors).length === 0;
	}

	function nextStep() {
		if (activeStep === 'setup' && !validateSetup()) return;
		activeStep = 'item';
	}

	function previousStep() {
		activeStep = 'setup';
	}

	async function createItem(event: SubmitEvent) {
		event.preventDefault();
		if (activeStep === 'setup') {
			nextStep();
			return;
		}

		if (!validateSetup()) {
			activeStep = 'setup';
			return;
		}
		if (!validateProperties()) {
			activeStep = 'item';
			return;
		}
		creating = true;

		try {
			await api.createItem({
				type_id: Number(selectedTypeID),
				properties: Object.entries(propertyValues)
					.filter(([id]) => selectedPropertyIDs.includes(Number(id)))
					.flatMap(([id, value]) => (value === null ? [] : [{ id: Number(id), value }])),
				amount
			});
			selectedTypeID = '';
			propertyValues = {};
			selectedPropertyIDs = [];
			toast.success('Stavka je dodata.');
			oncreated();
		} catch (reason) {
			toast.error(reason instanceof ApiError ? reason.message : 'Stavka nije sačuvana.');
		} finally {
			creating = false;
		}
	}

	function clearPropertyError(propertyID: number) {
		if (!propertyErrors[propertyID]) return;
		const nextErrors = { ...propertyErrors };
		delete nextErrors[propertyID];
		propertyErrors = nextErrors;
	}
</script>

<form id="item-form" class="flex min-h-0 flex-1 flex-col" novalidate onsubmit={createItem}>
	<Tabs.Root bind:value={activeStep} activationMode="manual" class="flex min-h-0 flex-1 flex-col">
		<Tabs.List
			class="grid grid-cols-2 gap-1 rounded-lg border border-line bg-soft p-1"
			aria-label="Koraci za novu stavku"
		>
			{#each steps as step, index (step.value)}
				<Tabs.Trigger
					value={step.value}
					disabled={step.value === 'item' && !selectedType}
					class="flex min-w-0 items-center justify-center gap-2 rounded-md px-2 py-2.5 text-sm text-muted transition-colors hover:bg-surface hover:text-ink disabled:cursor-not-allowed disabled:opacity-50 data-[state=active]:bg-surface data-[state=active]:font-medium data-[state=active]:text-ink data-[state=active]:shadow-sm"
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

		<Tabs.Content value="setup" class="min-h-0 flex-1 pt-6">
			<div class="grid gap-5">
				<div>
					<h2 class="text-base font-medium text-ink">Tip i količina</h2>
					<p class="mt-1 text-sm text-muted">
						Odaberite tip stavke i broj istih stavki koje želite da dodate.
					</p>
				</div>
				<div class="grid gap-5 sm:grid-cols-[minmax(0,1fr)_12rem]">
					<div>
						<Label.Root class="text-sm font-medium text-ink" for="new-item-type"
							>Tip stavke</Label.Root
						>
						<div class="mt-1">
							<OptionCombobox
								id="new-item-type"
								options={itemTypes}
								bind:value={selectedTypeID}
								placeholder="Odaberite tip"
								invalid={Boolean(fieldErrors.type)}
								describedBy={fieldErrors.type ? 'new-item-type-error' : undefined}
								onvaluechange={selectType}
							/>
						</div>
						<p class="mt-1 text-xs text-muted">Tip određuje skup dostupnih svojstava.</p>
						<p id="new-item-type-error" class="min-h-4 text-xs text-danger" aria-live="polite">
							{fieldErrors.type}
						</p>
					</div>
					<div>
						<Label.Root class="text-sm font-medium text-ink" for="new-item-amount"
							>Količina</Label.Root
						>
						<div class="mt-1">
							<NumberInput
								id="new-item-amount"
								bind:value={amount}
								ariaLabel="Količina"
								placeholder="Unesite količinu"
								min={1}
								max={100}
								required
								invalid={Boolean(fieldErrors.amount)}
								describedBy={fieldErrors.amount ? 'new-item-amount-error' : undefined}
								onvaluechange={() => {
									if (fieldErrors.amount) fieldErrors = { ...fieldErrors, amount: undefined };
								}}
							/>
							<p class="mt-1 text-xs text-muted">Broj istih stavki za unos.</p>
							<p id="new-item-amount-error" class="min-h-4 text-xs text-danger" aria-live="polite">
								{fieldErrors.amount}
							</p>
						</div>
					</div>
				</div>
				{#if itemTypes.length === 0}
					<FormAlert message="Prvo dodajte tip stavke u katalogu." />
				{/if}
			</div>
		</Tabs.Content>

		<Tabs.Content value="item" class="min-h-0 flex-1 pt-6">
			<div class="flex min-h-0 flex-1 flex-col">
				<div>
					<h2 class="text-base font-medium text-ink">Svojstva stavke</h2>
					<p class="mt-1 text-sm text-muted">
						Unesite vrednosti svojstava za stavku koju dodajete.
					</p>
				</div>

				{#if loadingType}
					<div class="grid min-h-48 flex-1 place-items-center" aria-busy="true">
						<Progress.Root
							value={null}
							class="size-8 animate-spin rounded-full border-4 border-line border-t-brand"
							aria-label="Učitavanje svojstava tipa"
						/>
					</div>
				{:else if selectedType?.properties.length}
					<section class="mt-5" aria-labelledby="new-item-properties-heading">
						<h3 id="new-item-properties-heading" class="sr-only">Svojstva stavke</h3>
						<div class="divide-y divide-line border-y border-line">
							{#each selectedType.properties as itemProperty (itemProperty.id)}
								{@const property = propertiesByID.get(itemProperty.id)}
								{#if property}
									<div
										class="grid grid-cols-[minmax(0,1fr)_minmax(0,2fr)_auto] items-center py-1.5 sm:grid-cols-[minmax(12rem,1fr)_minmax(0,2fr)_auto]"
									>
										<p class="flex h-8 items-center pr-3 text-sm text-muted sm:pr-6">
											{property.name}
										</p>
										<div class="col-start-2">
											{#if selectedPropertyIDs.includes(property.id)}
												<Label.Root class="sr-only" for={`new-item-property-${property.id}`}>
													Vrednost za {property.name}
												</Label.Root>
												<ItemPropertyValueInput
													id={`new-item-property-${property.id}`}
													bind:value={propertyValues[property.id]}
													className=""
													inputClassName={`text-sm ${propertyErrors[property.id] ? 'field-invalid' : ''}`}
													compact={true}
													{property}
													required
													onvaluechange={() => clearPropertyError(property.id)}
												/>
											{:else}
												<p
													class="flex h-8 items-center rounded-md border border-transparent pl-3 text-sm text-muted"
												>
													Nije uključeno.
												</p>
											{/if}
										</div>
										<label class="col-start-3 inline-flex size-8 items-center justify-center">
											<input
												type="checkbox"
												bind:group={selectedPropertyIDs}
												value={property.id}
												aria-label={`Uključi svojstvo ${property.name}`}
											/>
										</label>
									</div>
								{/if}
							{/each}
						</div>
					</section>
				{:else if selectedType}
					<p class="mt-5 text-sm text-muted">Izabrani tip nema svojstva za unos.</p>
				{:else}
					<p class="mt-5 text-sm text-muted">Najpre odaberite tip stavke.</p>
				{/if}
			</div>
		</Tabs.Content>
	</Tabs.Root>

	<Portal to="#page-footer-actions">
		<footer class="flex h-full items-center justify-between gap-3">
			<div class="flex items-center gap-3">
				<Button.Root
					class="rounded-md border border-line bg-surface px-4 py-2 text-sm font-medium text-ink hover:bg-soft disabled:opacity-50"
					disabled={activeStep === 'setup' || creating}
					type="button"
					onclick={previousStep}
				>
					Nazad
				</Button.Root>
				{#if activeStep === 'setup'}
					<Button.Root
						class="rounded-md bg-brand px-4 py-2 text-sm font-medium text-on-brand hover:bg-brand-strong disabled:opacity-60"
						disabled={creating || loadingType}
						type="button"
						onclick={nextStep}
					>
						Dalje
					</Button.Root>
				{:else}
					<Button.Root
						class="rounded-md bg-brand px-4 py-2 text-sm font-medium text-on-brand hover:bg-brand-strong disabled:opacity-60"
						disabled={creating || loadingType || !selectedType}
						form="item-form"
						type="submit"
					>
						{creating ? 'Čuvanje…' : 'Dodaj stavku'}
					</Button.Root>
				{/if}
			</div>
			{#if oncancel}
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
