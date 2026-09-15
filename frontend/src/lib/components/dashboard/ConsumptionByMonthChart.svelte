<script lang="ts">
	import { BarChart } from 'layerchart';
	import type { DashboardConsumptionBucket } from '$lib/api';
	import {
		consumptionSeries,
		consumptionSeriesLabel,
		formatMonthLong,
		formatMonthTick,
		itemCountLabel
	} from '$lib/domain/dashboard';

	let { buckets }: { buckets: DashboardConsumptionBucket[] } = $props();

	const ticks = $derived(
		new Map(
			buckets.map((bucket, index) => [
				bucket.month,
				formatMonthTick(bucket.month, buckets[index - 1]?.month)
			])
		)
	);

	const series = consumptionSeries.map(({ key, color }) => ({
		key,
		label: consumptionSeriesLabel(key),
		value: (bucket: DashboardConsumptionBucket) => bucket[key],
		color
	}));
</script>

<div class="h-64">
	<BarChart
		data={buckets}
		x="month"
		{series}
		seriesLayout="stack"
		bandPadding={0.3}
		legend
		props={{
			xAxis: { format: (month: string) => ticks.get(month) ?? month },
			yAxis: { format: (count: number) => String(count) },
			tooltip: {
				header: { format: (month: string) => formatMonthLong(month) },
				item: { format: (count: number) => itemCountLabel(count) }
			}
		}}
	/>
</div>
