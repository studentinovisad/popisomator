<script lang="ts">
	import { Slider } from 'bits-ui';

	// A threshold an item type can be configured with, where zero means the warning is off. The
	// backend stores that zero as null for every such column, so the slider bottoming out and the
	// setting being absent are the same state and need no separate control to express.
	let {
		id,
		value = $bindable<number>(),
		max,
		ariaLabel,
		format,
		disabled = false
	}: {
		id: string;
		value: number;
		max: number;
		ariaLabel: string;
		format: (value: number) => string;
		disabled?: boolean;
	} = $props();
</script>

<div class="mt-3 flex items-center gap-4">
	<Slider.Root
		type="single"
		bind:value
		min={0}
		{max}
		step={1}
		{disabled}
		class="relative flex w-full grow touch-none items-center select-none"
	>
		<span class="relative h-1 w-full grow rounded-full bg-soft">
			<Slider.Range class="absolute h-full rounded-full bg-brand" />
		</span>
		<Slider.Thumb
			index={0}
			{id}
			aria-label={ariaLabel}
			class="block size-4 cursor-grab rounded-full border border-brand bg-surface shadow-sm transition-colors hover:bg-brand-soft focus-visible:ring-2 focus-visible:ring-brand focus-visible:outline-none active:cursor-grabbing disabled:cursor-not-allowed"
		/>
	</Slider.Root>
	<!-- Fixed width and lining figures, so the track does not shift as the number gains a digit. -->
	<output
		for={id}
		class="w-28 shrink-0 text-right text-sm tabular-nums {value === 0 ? 'text-muted' : 'text-ink'}"
	>
		{format(value)}
	</output>
</div>
