<script lang="ts">
	import type { PropertyOption, PropertyValue } from '$lib/api';
	import NumberInput from '$lib/components/shared/NumberInput.svelte';
	import { measureUnits } from '$lib/domain/items';
	import MeasureInput from './MeasureInput.svelte';
	import PriceInput from './PriceInput.svelte';

	let {
		property,
		id,
		value = $bindable(),
		className = 'mt-2',
		inputClassName = '',
		compact = false,
		required = false,
		onvaluechange
	}: {
		property: PropertyOption;
		id: string;
		value: PropertyValue;
		className?: string;
		inputClassName?: string;
		compact?: boolean;
		required?: boolean;
		onvaluechange?: () => void;
	} = $props();

	let stringValue = $state('');
	let numberValue = $state(0);
	let booleanValue = $state(false);
	let objectValue = $state({});
	let lastCommittedValue = {};

	$effect(() => {
		if (value === lastCommittedValue) return;

		try {
			stringValue = typeof value === 'string' ? value : '';
			numberValue = typeof value === 'number' ? value : 0;
			objectValue = typeof value === 'object' && value != null ? value : {};
			booleanValue = value === true;
		} catch {
			stringValue = String(value);
			numberValue = 0;
			objectValue = typeof value === 'object' && value != null ? value : {};
			booleanValue = false;
		}

		lastCommittedValue = value;
	});

	function commitValue(nextValue: PropertyValue) {
		lastCommittedValue = nextValue;
		value = nextValue;
		onvaluechange?.();
	}

	function commitString() {
		commitValue(stringValue);
	}

	function commitNumber() {
		commitValue(numberValue);
	}

	function commitBoolean() {
		commitValue(booleanValue);
	}

	function commitObject() {
		commitValue(objectValue);
	}
</script>

{#if property.value_type === 'string'}
	<input
		{id}
		class={`${className} block w-full ${compact ? 'h-8' : 'h-10'} ${inputClassName}`}
		bind:value={stringValue}
		oninput={commitString}
		placeholder="Unesite tekst"
		{required}
	/>
{:else if property.value_type === 'number'}
	<NumberInput
		{id}
		bind:value={numberValue}
		ariaLabel={`Vrednost za ${property.name}`}
		{className}
		{inputClassName}
		{compact}
		placeholder="Unesite broj"
		{required}
		onvaluechange={commitNumber}
	/>
{:else if property.value_type === 'boolean'}
	<label
		class={`${className} flex items-center gap-2 rounded-md border border-line bg-surface px-3 text-sm text-ink ${compact ? 'h-8' : 'h-10'} ${inputClassName}`}
	>
		<input type="checkbox" bind:checked={booleanValue} onchange={commitBoolean} />
		{booleanValue ? 'Da' : 'Ne'}
	</label>
{:else if property.value_type === 'price'}
	<PriceInput
		{id}
		bind:value={objectValue}
		ariaLabel={property.name}
		{className}
		{inputClassName}
		{compact}
		{required}
		onvaluechange={commitObject}
	/>
{:else if property.value_type === 'expiry'}
	<input
		type="date"
		{id}
		class={`${className} block w-full ${compact ? 'h-8' : 'h-10'} ${inputClassName}`}
		bind:value={stringValue}
		oninput={commitString}
		placeholder="Unesite datum"
		{required}
	/>
{:else if property.value_type === 'mass' || property.value_type === 'volume'}
	<MeasureInput
		{id}
		bind:value={objectValue}
		units={measureUnits(property.value_type)}
		ariaLabel={property.name}
		{className}
		{inputClassName}
		{compact}
		{required}
		onvaluechange={commitObject}
	/>
{/if}
