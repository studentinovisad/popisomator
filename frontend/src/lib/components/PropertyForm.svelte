<script lang="ts">
	import {
		api,
		ApiError,
		type CreatePropertyRequest,
		type Property,
		type PropertyValueType
	} from '$lib/api';
	import { propertyValueTypeOptions } from '$lib/items';
	import { Button, Label, Select } from 'bits-ui';

	let {
		property,
		onsaved
	}: {
		property?: Property;
		onsaved: () => void;
	} = $props();

	let name = $state('');
	let description = $state('');
	let valueType = $state<PropertyValueType>('string');
	let hasDefaultValue = $state(false);
	let scalarDefaultValue = $state('');
	let booleanDefaultValue = $state(false);
	let saving = $state(false);
	let error = $state('');

	$effect(() => {
		name = property?.name ?? '';
		description = property?.description ?? '';
		valueType = property?.value_type ?? 'string';
		hasDefaultValue = property?.default_value !== null && property?.default_value !== undefined;
		setDefaultValue(property?.default_value ?? '');
	});

	function setDefaultValue(value: string) {
		try {
			const parsed = JSON.parse(value) as unknown;
			scalarDefaultValue =
				typeof parsed === 'string' || typeof parsed === 'number' ? String(parsed) : '';
			booleanDefaultValue = parsed === true;
		} catch {
			scalarDefaultValue = value;
			booleanDefaultValue = false;
		}
	}

	function changeValueType(value: string) {
		valueType = value as PropertyValueType;
		hasDefaultValue = false;
		scalarDefaultValue = '';
		booleanDefaultValue = false;
	}

	function serializedDefaultValue() {
		if (!hasDefaultValue) return null;
		if (valueType === 'string') return JSON.stringify(scalarDefaultValue);
		if (valueType === 'number') return scalarDefaultValue;
		return String(booleanDefaultValue);
	}

	async function save(event: SubmitEvent) {
		event.preventDefault();
		saving = true;
		error = '';

		try {
			const defaultValue = serializedDefaultValue();
			if (property) {
				await api.updateProperty(property.id, {
					name,
					description,
					default_value: defaultValue || null
				});
			} else {
				const payload: CreatePropertyRequest = {
					name,
					description,
					value_type: valueType,
					default_value: defaultValue || null
				};
				await api.createProperty(payload);
				name = '';
				description = '';
				valueType = 'string';
				hasDefaultValue = false;
				scalarDefaultValue = '';
				booleanDefaultValue = false;
			}
			onsaved();
		} catch (reason) {
			error = reason instanceof ApiError ? reason.message : 'Svojstvo nije sačuvano.';
		} finally {
			saving = false;
		}
	}
</script>

<form class="grid gap-4 sm:grid-cols-2" onsubmit={save}>
	<div>
		<Label.Root class="text-sm font-medium text-ink" for="property-name">Naziv</Label.Root>
		<input id="property-name" class="mt-1 block w-full" bind:value={name} required />
	</div>
	<div>
		<Label.Root class="text-sm font-medium text-ink" for="property-value-type"
			>Tip vrednosti</Label.Root
		>
		{#if property}
			<input
				class="mt-1 block w-full"
				value={propertyValueTypeOptions.find((option) => option.value === valueType)?.label}
				disabled
			/>
		{:else}
			<Select.Root
				type="single"
				value={valueType}
				items={propertyValueTypeOptions}
				onValueChange={changeValueType}
			>
				<Select.Trigger
					id="property-value-type"
					class="mt-1 flex h-10 w-full items-center justify-between rounded-md border border-line bg-surface px-3 text-sm text-ink"
				>
					<Select.Value />
				</Select.Trigger>
				<Select.Portal>
					<Select.Content
						class="z-40 w-[var(--bits-select-anchor-width)] rounded-md border border-line bg-surface p-1 shadow-lg shadow-ink/10"
					>
						<Select.Viewport>
							{#each propertyValueTypeOptions as option (option.value)}
								<Select.Item
									value={option.value}
									label={option.label}
									class="cursor-pointer rounded px-3 py-2 text-sm outline-none data-highlighted:bg-brand-soft"
								>
									{option.label}
								</Select.Item>
							{/each}
						</Select.Viewport>
					</Select.Content>
				</Select.Portal>
			</Select.Root>
		{/if}
	</div>
	<div class="sm:col-span-2">
		<Label.Root class="text-sm font-medium text-ink" for="property-description">Opis</Label.Root>
		<textarea id="property-description" class="mt-1 block min-h-20 w-full" bind:value={description}
		></textarea>
	</div>
	<div class="sm:col-span-2">
		<label class="flex items-center gap-2 text-sm font-medium text-ink">
			<input type="checkbox" bind:checked={hasDefaultValue} />
			Podrazumevana vrednost
		</label>
		{#if hasDefaultValue}
			{#if valueType === 'string'}
				<Label.Root class="sr-only" for="property-default-value">Podrazumevani tekst</Label.Root>
				<input
					id="property-default-value"
					class="mt-2 block w-full"
					bind:value={scalarDefaultValue}
					placeholder="npr. crno"
				/>
			{:else if valueType === 'number'}
				<Label.Root class="sr-only" for="property-default-value">Podrazumevani broj</Label.Root>
				<input
					id="property-default-value"
					class="mt-2 block w-full"
					type="number"
					bind:value={scalarDefaultValue}
					required
				/>
			{:else if valueType === 'boolean'}
				<label
					class="mt-2 flex h-10 items-center gap-2 rounded-md border border-line bg-surface px-3 text-sm text-ink"
				>
					<input type="checkbox" bind:checked={booleanDefaultValue} />
					{booleanDefaultValue ? 'Da' : 'Ne'}
				</label>
			{/if}
		{/if}
	</div>
	<div class="sm:col-span-2">
		<Button.Root
			class="rounded-md bg-brand px-4 py-2 text-sm font-medium text-on-brand hover:bg-brand-strong disabled:opacity-60"
			disabled={saving}
			type="submit"
		>
			{saving ? 'Čuvanje…' : property ? 'Sačuvaj izmene' : 'Dodaj svojstvo'}
		</Button.Root>
		{#if error}<p class="mt-3 text-sm text-danger" role="alert">{error}</p>{/if}
	</div>
</form>
