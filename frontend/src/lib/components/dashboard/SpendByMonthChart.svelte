<script lang="ts">
	import { BarChart } from 'layerchart';
	import type { MonthlySpend } from '$lib/domain/dashboard';
	import { formatMonthLong, formatMonthTick } from '$lib/domain/dashboard';

	let { spend }: { spend: MonthlySpend[] } = $props();

	const currency = $derived(spend.find((entry) => entry.currency)?.currency ?? 'RSD');
	const amountFormatter = $derived(
		new Intl.NumberFormat('sr-Latn-RS', {
			style: 'currency',
			currency,
			maximumFractionDigits: 0
		})
	);
	const ticks = $derived(
		new Map(
			spend.map((entry, index) => [
				entry.month,
				formatMonthTick(entry.month, spend[index - 1]?.month)
			])
		)
	);
</script>

<div class="h-64 overflow-hidden">
	<BarChart
		data={spend}
		x="month"
		y="amount"
		bandPadding={0.3}
		props={{
			xAxis: { format: (month: string) => ticks.get(month) ?? month },
			yAxis: { format: (amount: number) => amountFormatter.format(amount) },
			tooltip: {
				header: { format: (month: string) => formatMonthLong(month) },
				item: { label: 'Trošak', format: (amount: number) => amountFormatter.format(amount) }
			}
		}}
	/>
</div>
