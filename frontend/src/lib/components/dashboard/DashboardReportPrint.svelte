<script lang="ts">
	import type { Dashboard } from '$lib/api';
	import {
		consumptionPropertyRows,
		formatMonthLong,
		formatPropertyTotal,
		formatPropertyTotals,
		monthRangeLabel
	} from '$lib/domain/dashboard';

	let { report }: { report: Dashboard } = $props();
</script>

<article class="dashboard-report mx-auto max-w-4xl px-6 py-10 text-ink sm:px-10" aria-hidden="true">
	<header class="border-b-2 border-ink pb-6">
		<h1 class="mt-1 text-2xl font-semibold tracking-tight">Izveštaj o zalihama i potrošnji</h1>
		<p class="mt-3 text-sm text-muted">
			Period: <span class="font-medium text-ink">poslednjih {monthRangeLabel(report.months)}</span>
		</p>
	</header>

	<section class="mt-8 break-inside-avoid">
		<h2 class="border-b border-ink pb-1 text-base font-semibold">Trenutne zalihe</h2>
		<table class="mt-3 w-full text-sm">
			<thead>
				<tr class="border-b border-ink text-left text-xs text-muted">
					<th class="py-1.5 pr-3 font-medium">Tip</th>
					<th class="py-1.5 pr-3 font-medium">Naziv</th>
					<th class="py-1.5 pr-3 text-right font-medium">Na stanju</th>
					<th class="py-1.5 text-right font-medium">Količina</th>
				</tr>
			</thead>
			<tbody>
				{#each report.stock_groups as group (`${group.type_id}-${group.name}`)}
					<tr class="border-b border-line last:border-0">
						<td class="py-1.5 pr-3 text-muted">{group.type_name}</td>
						<td class="py-1.5 pr-3">{group.name}</td>
						<td class="py-1.5 pr-3 text-right tabular-nums">{group.in_stock_count}</td>
						<td class="py-1.5 text-right tabular-nums">
							{group.totals.length > 0 ? formatPropertyTotals(group.totals) : '—'}
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</section>

	<section class="mt-8">
		<h2 class="border-b border-ink pb-1 text-base font-semibold">
			Potrošnja i troškovi po mesecima
		</h2>
		{#each report.consumption_quantity as type (type.type_id)}
			<div class="mt-4 break-inside-avoid">
				<h3 class="text-sm font-semibold">{type.type_name}</h3>
				<table class="mt-2 w-full text-sm">
					<thead>
						<tr class="border-b border-ink text-left text-xs text-muted">
							<th class="py-1.5 pr-3 font-medium">Svojstvo</th>
							{#each type.buckets as bucket (bucket.month)}
								<th class="py-1.5 pr-3 text-right font-medium">{formatMonthLong(bucket.month)}</th>
							{/each}
							<th class="py-1.5 text-right font-medium">Ukupno</th>
						</tr>
					</thead>
					<tbody>
						{#each consumptionPropertyRows(type) as row (row.propertyID)}
							<tr class="border-b border-line last:border-0">
								<td class="py-1.5 pr-3">{row.periodTotal.property_name}</td>
								{#each row.cells as cell, index (type.buckets[index].month)}
									<td class="py-1.5 pr-3 text-right tabular-nums">
										{cell ? formatPropertyTotal(cell) : '—'}
									</td>
								{/each}
								<td class="py-1.5 text-right font-semibold tabular-nums">
									{formatPropertyTotal(row.periodTotal)}
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		{/each}
	</section>

	<section class="mt-8 break-inside-avoid">
		<h2 class="border-b border-ink pb-1 text-base font-semibold">Najviše korišćeno</h2>
		<table class="mt-3 w-full text-sm">
			<thead>
				<tr class="border-b border-ink text-left text-xs text-muted">
					<th class="w-8 py-1.5 pr-3 font-medium">#</th>
					<th class="py-1.5 pr-3 font-medium">Naziv</th>
					<th class="py-1.5 pr-3 font-medium">Tip</th>
					<th class="py-1.5 pr-3 text-right font-medium">Potrošeno</th>
					<th class="py-1.5 text-right font-medium">Količina</th>
				</tr>
			</thead>
			<tbody>
				{#each report.most_consumed as group, index (`${group.type_id}-${group.name}`)}
					<tr class="border-b border-line last:border-0">
						<td class="py-1.5 pr-3 tabular-nums">{index + 1}</td>
						<td class="py-1.5 pr-3">{group.name}</td>
						<td class="py-1.5 pr-3 text-muted">{group.type_name}</td>
						<td class="py-1.5 pr-3 text-right font-semibold tabular-nums">{group.consumed_count}</td
						>
						<td class="py-1.5 text-right tabular-nums">
							{group.totals.length > 0 ? formatPropertyTotals(group.totals) : '—'}
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</section>
</article>

<style>
	.dashboard-report {
		display: none;
	}

	@media print {
		:global(body) {
			background: white;
		}

		:global(.app-shell) {
			display: none;
		}

		.dashboard-report {
			display: block;
			max-width: none;
			padding: 0;
			color-scheme: light;
			--color-canvas: #f7f7f6;
			--color-surface: #ffffff;
			--color-ink: #272726;
			--color-muted: #6e6e6b;
			--color-faint: #737370;
			--color-line: #e2e2df;
			--color-soft: #f0f0ee;
		}
	}
</style>
