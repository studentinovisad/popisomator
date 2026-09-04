<script lang="ts">
	import {
		api,
		ApiError,
		type CreatePropertyRequest,
		type Property,
		type PropertyValue,
		type PropertyValueType
	} from '$lib/api';
	import { defaultJsonValue, propertyValueTypeOptions } from '$lib/domain/items';
	import { propertyValueError, requiredTextError } from '$lib/domain/form-validation';
	import { Button, Label, Select, Separator } from 'bits-ui';
	import { toast } from 'svelte-sonner';
	import ItemPropertyValueInput from '../inventory/ItemPropertyValueInput.svelte';

	let {
		property,
		onsaved,
		oncancel
	}: {
		property?: Property;
		onsaved: () => void;
		oncancel?: () => void;
	} = $props();

	let name = $state('');
	let description = $state('');
	let valueType = $state<PropertyValueType>('string');
	let hasDefaultValue = $state(false);
	let defaultValue = $state<PropertyValue>('');
	let saving = $state(false);
	let nameError = $state('');
	let defaultValueError = $state('');

	$effect(() => {
		name = property?.name ?? '';
		description = property?.description ?? '';
		valueType = property?.value_type ?? 'string';
		hasDefaultValue = property?.default_value !== null && property?.default_value !== undefined;
		defaultValue = property?.default_value ?? '';
	});

	function changeValueType(value: string) {
		valueType = value as PropertyValueType;
		hasDefaultValue = false;
		defaultValue = defaultJsonValue(valueType, null);
	}

	async function save(event: SubmitEvent) {
		event.preventDefault();
		nameError = requiredTextError(name, 'naziv svojstva') ?? '';
		defaultValueError = hasDefaultValue
			? (propertyValueError({ value_type: valueType }, defaultValue) ?? '')
			: '';
		if (nameError || defaultValueError) return;
		saving = true;

		try {
			const savedDefaultValue = hasDefaultValue ? defaultValue : null;
			if (property) {
				await api.updateProperty(property.id, {
					name,
					description,
					default_value: savedDefaultValue
				});
			} else {
				const payload: CreatePropertyRequest = {
					name,
					description,
					value_type: valueType,
					default_value: savedDefaultValue
				};
				await api.createProperty(payload);
				name = '';
				description = '';
				valueType = 'string';
				hasDefaultValue = false;
				defaultValue = '';
			}
			toast.success(property ? 'Svojstvo je izmenjeno.' : 'Svojstvo je dodato.');
			onsaved();
		} catch (reason) {
			toast.error(reason instanceof ApiError ? reason.message : 'Svojstvo nije sačuvano.');
		} finally {
			saving = false;
		}
	}
</script>

<form class="grid gap-4 sm:grid-cols-2" novalidate onsubmit={save}>
	<div>
		<Label.Root class="text-sm font-medium text-ink" for="property-name">Naziv</Label.Root>
		<input
			id="property-name"
			class={`mt-1 block w-full ${nameError ? 'field-invalid' : ''}`}
			bind:value={name}
			aria-invalid={Boolean(nameError)}
			aria-describedby={nameError ? 'property-name-error' : undefined}
			oninput={() => (nameError = '')}
		/>
		{#if nameError}<p id="property-name-error" class="mt-1 text-xs text-danger" role="alert">
				{nameError}
			</p>{/if}
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
						class="z-40 w-(--bits-select-anchor-width) rounded-md border border-line bg-surface p-1 shadow-lg shadow-black/15"
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
	<Separator.Root class="h-px bg-line sm:col-span-2" decorative />
	<div class="sm:col-span-2">
		<label class="flex items-center gap-2 text-sm font-medium text-ink">
			<input type="checkbox" bind:checked={hasDefaultValue} />
			Podrazumevana vrednost
		</label>
		{#if hasDefaultValue}
			<ItemPropertyValueInput
				id="property-default-value"
				property={{ id: 0, value_type: valueType, name: name, default_value: null }}
				bind:value={defaultValue}
				inputClassName={defaultValueError ? 'field-invalid' : ''}
				onvaluechange={() => (defaultValueError = '')}
			/>
			{#if defaultValueError}
				<p id="property-default-value-error" class="mt-1 text-xs text-danger" role="alert">
					{defaultValueError}
				</p>
			{/if}
		{/if}
	</div>
	<Separator.Root class="h-px bg-line sm:col-span-2" decorative />
	<div class="sm:col-span-2">
		<div class="flex flex-wrap items-center gap-3">
			<Button.Root
				class="rounded-md bg-brand px-4 py-2 text-sm font-medium text-on-brand hover:bg-brand-strong disabled:opacity-60"
				disabled={saving}
				type="submit"
			>
				{saving ? 'Čuvanje…' : property ? 'Sačuvaj izmene' : 'Dodaj svojstvo'}
			</Button.Root>
			{#if !property && oncancel}
				<Button.Root
					class="rounded-md border border-line bg-surface px-4 py-2 text-sm font-medium text-ink hover:bg-soft disabled:opacity-60"
					disabled={saving}
					type="button"
					onclick={oncancel}
				>
					Otkaži dodavanje
				</Button.Root>
			{/if}
		</div>
	</div>
</form>
