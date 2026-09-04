<script lang="ts">
	import { type PTPrice } from '$lib/api';
	import NumberInput from '$lib/components/shared/NumberInput.svelte';
	import { PriceMultiplier, scaleAmount, unscaleAmount } from '$lib/domain/items';
	import { Select } from 'bits-ui';

	let {
		id,
		value = $bindable<PTPrice>(),
		ariaLabel,
		required = false,
		disabled = false,
		className = '',
		inputClassName = '',
		compact = false,
		onvaluechange
	}: {
		id: string;
		value: Partial<PTPrice>;
		ariaLabel: string;
		required?: boolean;
		disabled?: boolean;
		className?: string;
		inputClassName?: string;
		compact?: boolean;
		onvaluechange?: () => void;
	} = $props();

    const priceCurrencies = [
        "RSD",
        "EUR"
    ];

	// The field keeps its own copy of the amount instead of deriving it from the bound value: an
	// emptied input binds as null, and a derived amount would paint the committed 0 straight back
	// into the field, so the last digit could never be erased. lastCommittedAmount is what this
	// field wrote most recently, so a value replaced from the outside still refreshes the input.
	let amount = $state<number | string>('');
	let lastCommittedAmount: number | undefined;

	$effect(() => {
		if (value.amount === lastCommittedAmount) return;

		amount = unscaleAmount(value.amount, PriceMultiplier);
		lastCommittedAmount = value.amount;
	});

	// Replace the bound object rather than mutating it in place. Callers seed their draft with the
	// item's existing value, so an in-place edit would also rewrite the original they diff against
	// and the change would look like no change at all.
	function commit(changed: Partial<PTPrice>) {
		value = { ...value, ...changed };
		onvaluechange?.();
	}

	function onAmountInput() {
		lastCommittedAmount = scaleAmount(amount, PriceMultiplier);
		commit({ amount: lastCommittedAmount });
	}

	function onCurrencyChange(currency: string) {
		commit({ currency });
	}
</script>

<div class={`flex gap-2 ${className}`}>
	<NumberInput
		{id}
		bind:value={amount}
		ariaLabel="Vrednost za {ariaLabel}"
		placeholder="Unesite vrednost"
		className="min-w-0 flex-1"
		{inputClassName}
		min={0}
		max={1000000000000}
		step={0.0001}
		stepIncrement={1}
		{required}
		{disabled}
		{compact}
		onvaluechange={onAmountInput}
	/>
	<Select.Root
		type="single"
		value={value.currency ?? ''}
		onValueChange={onCurrencyChange}
		{required}
		{disabled}
	>
		<Select.Trigger
			class={`flex ${compact ? 'h-8' : 'h-10'} w-40 items-center justify-between rounded-md border border-line bg-surface px-3 text-sm text-ink transition-colors hover:border-brand/40 ${inputClassName}`}
			aria-label="Valuta za {ariaLabel}"
		>
			<Select.Value />
		</Select.Trigger>
		<Select.Portal>
			<Select.Content
				class="z-30 w-44 rounded-md border border-line bg-surface p-1 shadow-lg shadow-black/15"
				sideOffset={4}
			>
				<Select.Viewport>
                    {#each priceCurrencies as option}
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
