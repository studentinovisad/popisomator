<script lang="ts">
	import ThresholdSlider from '$lib/components/catalog/ThresholdSlider.svelte';
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

<ThresholdSlider
	{id}
	bind:value={days}
	max={maximumExpiringSoonDays}
	ariaLabel="Broj dana pre isteka roka"
	format={(value) => (value === 0 ? 'Bez upozorenja' : `${value} dana`)}
	{disabled}
/>
