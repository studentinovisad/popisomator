<script lang="ts">
	import { Slider } from 'bits-ui';
	import { maximumExpiringSoonDays } from '$lib/domain/items';

	let {
		id,
		days = $bindable<number>(),
		disabled = false
	}: {
		id: string;
		days: number;
		disabled?: boolean;
	} = $props();
</script>

<div class="mt-3 flex items-center gap-4">
	<Slider.Root
		type="single"
		bind:value={days}
		min={0}
		max={maximumExpiringSoonDays}
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
			aria-label="Broj dana pre isteka roka"
			class="block size-4 cursor-grab rounded-full border border-brand bg-surface shadow-sm transition-colors hover:bg-brand-soft focus-visible:ring-2 focus-visible:ring-brand focus-visible:outline-none active:cursor-grabbing disabled:cursor-not-allowed"
		/>
	</Slider.Root>
	<!-- Fixed width and lining figures, so the track does not shift as the number gains a digit. -->
	<output
		for={id}
		class="w-28 shrink-0 text-right text-sm tabular-nums {days === 0 ? 'text-muted' : 'text-ink'}"
	>
		{days === 0 ? 'Bez upozorenja' : `${days} dana`}
	</output>
</div>
