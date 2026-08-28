<script lang="ts">
	import type { PropertyOption } from '$lib/api';
	import NumberInput from '$lib/components/shared/NumberInput.svelte';
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
		value: string;
		className?: string;
		inputClassName?: string;
		compact?: boolean;
		required?: boolean;
		onvaluechange?: () => void;
	} = $props();

	let scalarValue = $state('');
	let booleanValue = $state(false);
	let objectValue = $state({});
	let lastCommittedValue = '\u0000';

	$effect(() => {
		if (value === lastCommittedValue) return;

		try {
			const parsed = JSON.parse(value) as unknown;
			scalarValue = typeof parsed === 'string' || typeof parsed === 'number' ? String(parsed) : '';
			objectValue = typeof parsed === 'object' && parsed != null ? parsed : {};
			booleanValue = parsed === true;
		} catch {
			scalarValue = value;
			objectValue = value;
			booleanValue = false;
		}

		lastCommittedValue = value;
	});

	function commitScalar() {
		const nextValue = JSON.stringify(scalarValue);
		lastCommittedValue = nextValue;
		value = nextValue;
		onvaluechange?.();
	}

	function commitBoolean() {
		const nextValue = String(booleanValue);
		lastCommittedValue = nextValue;
		value = nextValue;
		onvaluechange?.();
	}

	function commitObject() {
		const nextValue = JSON.stringify(objectValue);
		lastCommittedValue = nextValue;
		value = nextValue;
		onvaluechange?.();
	}
</script>

{#if property.value_type === 'string'}
	<input
		{id}
		class={`${className} block w-full ${compact ? 'h-8' : 'h-10'} ${inputClassName}`}
		bind:value={scalarValue}
		oninput={commitScalar}
		placeholder="Unesite tekst"
		{required}
	/>
{:else if property.value_type === 'number'}
	<NumberInput
		{id}
		bind:value={scalarValue}
		ariaLabel={`Vrednost za ${property.name}`}
		{className}
		{inputClassName}
		{compact}
		placeholder="Unesite broj"
		{required}
		onvaluechange={commitScalar}
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
{/if}
