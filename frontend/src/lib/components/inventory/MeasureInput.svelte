<script lang="ts">
	import { type PTMeasure } from '$lib/api';
	import { MeasureMultiplier } from '$lib/domain/items';
	import { Select } from 'bits-ui';

	let {
		id,
		value = $bindable<PTMeasure>(),
		units,
		ariaLabel,
		required = false,
		disabled = false,
		className = '',
		inputClassName = '',
		compact = false,
		onvaluechange
	}: {
		id: string;
		value: Partial<PTMeasure>;
		units: string[];
		ariaLabel: string;
		required?: boolean;
		disabled?: boolean;
		className?: string;
		inputClassName?: string;
		compact?: boolean;
		onvaluechange?: () => void;
	} = $props();

	let amount = $derived(value.amount != undefined ? value.amount / MeasureMultiplier : 0);

	function onAmountInput() {
		value.amount = Math.floor(amount * MeasureMultiplier);
		onvaluechange?.();
	}
</script>

<div {id} class={`relative flex gap-2 ${className}`}>
	<input
		class={`number-input-field block w-full pr-9 ${compact ? 'h-8' : 'h-10'} ${inputClassName}`}
		type="number"
		bind:value={amount}
		aria-label="Vrednost za {ariaLabel}"
		placeholder="Unesite vrednost"
		min="0"
		max="1000000000000"
		step="0.0001"
		{required}
		{disabled}
		oninput={onAmountInput}
	/>
	<Select.Root
		type="single"
		bind:value={value.unit}
		onValueChange={onvaluechange}
		{required}
		{disabled}
	>
		<Select.Trigger
			class={`flex ${compact ? 'h-8' : 'h-10'} w-40 items-center justify-between rounded-md border border-line bg-surface px-3 text-sm text-ink transition-colors hover:border-brand/40 ${inputClassName}`}
			aria-label="Jedinica mere za {ariaLabel}"
		>
			<Select.Value />
		</Select.Trigger>
		<Select.Portal>
			<Select.Content
				class="z-30 w-44 rounded-md border border-line bg-surface p-1 shadow-lg shadow-black/15"
				sideOffset={4}
			>
				<Select.Viewport>
					{#each units as option (option)}
						<Select.Item
							value={option}
							label={option}
							class="cursor-pointer rounded px-3 py-2 text-sm outline-none data-highlighted:bg-brand-soft"
						>
							{option}
						</Select.Item>
					{/each}
				</Select.Viewport>
			</Select.Content>
		</Select.Portal>
	</Select.Root>
</div>
