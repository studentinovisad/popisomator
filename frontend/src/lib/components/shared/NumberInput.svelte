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
		stepIncrement,
		required = false,
		disabled = false,
		className = '',
		inputClassName = '',
		compact = false,
		onvaluechange
	}: {
		id: string;
		value: NumericValue;
		ariaLabel: string;
		placeholder?: string;
		min?: number;
		max?: number;
		step?: number;
		stepIncrement?: number;
		required?: boolean;
		disabled?: boolean;
		className?: string;
		inputClassName?: string;
		compact?: boolean;
		onvaluechange?: () => void;
	} = $props();

	// `step` is how fine a value the field accepts - a price takes four decimals - which is not the
	// same question as how far one press of the arrows moves it. A price steps by a whole dinar, not
	// by the 0.0001 it can store, so stepIncrement carries that separately and falls back to `step`
	// for the fields where the two coincide.
	let incrementAmount = $derived(stepIncrement ?? step);

	// A fractional step drifts when it is added repeatedly - 0.1 + 0.0001 lands on
	// 0.10009999999999999 - so each adjustment is rounded back to the precision `step` allows.
	let stepDecimals = $derived((String(step).split('.')[1] ?? '').length);

	function adjust(direction: 1 | -1) {
		const currentValue = Number(value);
		let nextValue = Number.isFinite(currentValue)
			? currentValue + direction * incrementAmount
			: (min ?? 0);

		if (min !== undefined) nextValue = Math.max(min, nextValue);
		if (max !== undefined) nextValue = Math.min(max, nextValue);
		nextValue = Number(nextValue.toFixed(stepDecimals));

		value = typeof value === 'number' ? nextValue : String(nextValue);
		onvaluechange?.();
	}

	// The arrow keys move a number input by its step attribute, which would leave the keyboard
	// creeping along in 0.0001s while a click on the chevron beside it moves a whole unit. Drive
	// both from adjust() so they agree.
	function onArrowKey(event: KeyboardEvent) {
		if (event.key !== 'ArrowUp' && event.key !== 'ArrowDown') return;

		event.preventDefault();
		adjust(event.key === 'ArrowUp' ? 1 : -1);
	}
</script>

<div class={`relative ${className}`}>
	<input
		{id}
		class={`number-input-field block w-full pr-9 ${compact ? 'h-8' : 'h-10'} ${inputClassName}`}
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
		onkeydown={onArrowKey}
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
