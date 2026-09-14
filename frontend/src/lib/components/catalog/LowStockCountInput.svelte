<script lang="ts">
	import ThresholdSlider from '$lib/components/catalog/ThresholdSlider.svelte';
	import { maximumLowStockCount } from '$lib/domain/items';

	let {
		id,
		count = $bindable<number>(),
		disabled = false
	}: {
		id: string;
		count: number;
		disabled?: boolean;
	} = $props();
</script>

<ThresholdSlider
	{id}
	bind:value={count}
	max={maximumLowStockCount}
	ariaLabel="Najmanji broj dostupnih stavki"
	format={(value) => {
		if (value === 0) return 'Bez upozorenja';
		// 1 komad, but 2 and up komada - and 11 goes with the many form despite ending in one, which
		// is why the last two digits decide rather than the last.
		const singular = value % 10 === 1 && value % 100 !== 11;
		return `${value} ${singular ? 'komad' : 'komada'}`;
	}}
	{disabled}
/>
