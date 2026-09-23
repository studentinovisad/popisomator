<script lang="ts">
	import { BarChart } from 'layerchart';
	import type { DashboardMonthCount } from '$lib/api';
	import { formatMonthLong, formatMonthTick, itemCountLabel } from '$lib/domain/dashboard';

	let { buckets }: { buckets: DashboardMonthCount[] } = $props();

	// The tick carries a year only where it changes, which needs the bucket before it.
	const ticks = $derived(
		new Map(
			buckets.map((bucket, index) => [
				bucket.month,
				formatMonthTick(bucket.month, buckets[index - 1]?.month)
			])
		)
	);
</script>

<div class="h-64 overflow-hidden">
	<BarChart
		data={buckets}
		x="month"
		y="count"
		bandPadding={0.3}
		props={{
			xAxis: { format: (month: string) => ticks.get(month) ?? month },
			yAxis: { format: (count: number) => String(count) },
			tooltip: {
				header: { format: (month: string) => formatMonthLong(month) },
				item: { label: 'Ističe', format: (count: number) => itemCountLabel(count) }
			}
		}}
	/>
</div>
