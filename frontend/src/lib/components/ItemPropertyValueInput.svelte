<script lang="ts">
	import type { Property } from '$lib/api';

	let {
		property,
		id,
		value = $bindable(),
		required = false
	}: {
		property: Property;
		id: string;
		value: string;
		required?: boolean;
	} = $props();

	let scalarValue = $state('');
	let booleanValue = $state(false);
	let lastCommittedValue = '\u0000';

	$effect(() => {
		if (value === lastCommittedValue) return;

		try {
			const parsed = JSON.parse(value) as unknown;
			scalarValue = typeof parsed === 'string' || typeof parsed === 'number' ? String(parsed) : '';
			booleanValue = parsed === true;
		} catch {
			scalarValue = value;
			booleanValue = false;
		}

		lastCommittedValue = value;
	});

	function commitScalar() {
		const nextValue = property.value_type === 'string' ? JSON.stringify(scalarValue) : scalarValue;
		lastCommittedValue = nextValue;
		value = nextValue;
	}

	function commitBoolean() {
		const nextValue = String(booleanValue);
		lastCommittedValue = nextValue;
		value = nextValue;
	}
</script>

{#if property.value_type === 'string'}
	<input
		{id}
		class="mt-2 block w-full"
		bind:value={scalarValue}
		oninput={commitScalar}
		placeholder="Unesite tekst"
		{required}
	/>
{:else if property.value_type === 'number'}
	<input
		{id}
		class="mt-2 block w-full"
		type="number"
		bind:value={scalarValue}
		oninput={commitScalar}
		placeholder="Unesite broj"
		{required}
	/>
{:else if property.value_type === 'boolean'}
	<label
		class="mt-2 flex h-10 items-center gap-2 rounded-md border border-line bg-surface px-3 text-sm text-ink"
	>
		<input type="checkbox" bind:checked={booleanValue} onchange={commitBoolean} />
		{booleanValue ? 'Da' : 'Ne'}
	</label>
{/if}
