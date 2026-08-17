<script lang="ts">
	import ChevronDown from '@lucide/svelte/icons/chevron-down';
	import ChevronUp from '@lucide/svelte/icons/chevron-up';
	import { Button } from 'bits-ui';

	type NumericValue = number | string;

	let {
		id,
		value = $bindable<NumericValue>(),
		ariaLabel,
		placeholder,
		min,
		max,
		step = 1,
		required = false,
		disabled = false,
		className = '',
		onvaluechange
	}: {
		id: string;
		value: NumericValue;
		ariaLabel: string;
		placeholder?: string;
		min?: number;
		max?: number;
		step?: number;
		required?: boolean;
		disabled?: boolean;
		className?: string;
		onvaluechange?: () => void;
	} = $props();

	function adjust(direction: 1 | -1) {
		const currentValue = Number(value);
		let nextValue = Number.isFinite(currentValue) ? currentValue + direction * step : (min ?? 0);

		if (min !== undefined) nextValue = Math.max(min, nextValue);
		if (max !== undefined) nextValue = Math.min(max, nextValue);

		value = typeof value === 'number' ? nextValue : String(nextValue);
		onvaluechange?.();
	}
</script>

<div class={`relative ${className}`}>
	<input
		{id}
		class="number-input-field block h-10 w-full pr-9"
		type="number"
		bind:value
		aria-label={ariaLabel}
		{placeholder}
		{min}
		{max}
		{step}
		{required}
		{disabled}
		oninput={onvaluechange}
	/>
	<div
		class="absolute inset-y-px right-px flex w-8 flex-col overflow-hidden rounded-r-[calc(0.375rem-1px)] border-l border-line bg-soft"
	>
		<Button.Root
			class="inline-flex flex-1 items-center justify-center text-muted transition-colors hover:bg-brand-soft hover:text-brand disabled:cursor-not-allowed disabled:opacity-50"
			type="button"
			aria-label={`Povećaj: ${ariaLabel}`}
			{disabled}
			onclick={() => adjust(1)}
		>
			<ChevronUp class="size-3.5" aria-hidden="true" />
		</Button.Root>
		<Button.Root
			class="inline-flex flex-1 items-center justify-center border-t border-line text-muted transition-colors hover:bg-brand-soft hover:text-brand disabled:cursor-not-allowed disabled:opacity-50"
			type="button"
			aria-label={`Smanji: ${ariaLabel}`}
			{disabled}
			onclick={() => adjust(-1)}
		>
			<ChevronDown class="size-3.5" aria-hidden="true" />
		</Button.Root>
	</div>
</div>
