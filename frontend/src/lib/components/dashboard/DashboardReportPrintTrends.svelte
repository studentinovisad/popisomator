<script lang="ts">
	import type { Dashboard, DashboardTypeConsumptionQuantity, PropertyValueType } from '$lib/api';
	import {
		consumptionGroupRows,
		formatMonthLong,
		formatPropertyTotal,
		formatPropertyTotals,
		monthRangeLabel,
		sortGroupsByTotal,
		sumRowsFooters
	} from '$lib/domain/dashboard';

	let { report }: { report: Dashboard } = $props();

	function visibleTypes(
		types: DashboardTypeConsumptionQuantity[],
		valueTypes: PropertyValueType[]
	) {
		return types
			.map((type) => {
				const groups = sortGroupsByTotal(
					type.groups.filter((group) => consumptionGroupRows(group, valueTypes).length > 0),
					valueTypes
				);
				const rows = groups.flatMap((group) => consumptionGroupRows(group, valueTypes));
				return { ...type, groups, rows, footers: sumRowsFooters(rows) };
			})
			.filter((type) => type.groups.length > 0);
	}
</script>

<article
	class="dashboard-report-trends mx-auto max-w-6xl px-6 py-10 text-ink sm:px-10"
	aria-hidden="true"
>
	<header class="border-b-2 border-ink pb-6">
		<h1 class="mt-1 text-2xl font-semibold tracking-tight">Izveštaj o potrošnji i troškovima</h1>
		<p class="mt-3 text-sm text-muted">
			Period: <span class="font-medium text-ink">poslednjih {monthRangeLabel(report.months)}</span>
		</p>
	</header>

	{#each [{ heading: 'Potrošene količine po mesecima', valueTypes: ['mass', 'volume'] as PropertyValueType[] }, { heading: 'Troškovi po mesecima', valueTypes: ['price'] as PropertyValueType[] }] as section (section.heading)}
		{@const types = visibleTypes(report.consumption_quantity, section.valueTypes)}
		{#if types.length > 0}
			<section class="mt-8">
				<h2 class="border-b border-ink pb-1 text-base font-semibold">{section.heading}</h2>
				{#each types as type (type.type_id)}
					<div class="mt-4">
						<h3 class="text-sm font-semibold">{type.type_name}</h3>
						<table class="mt-2 w-full text-sm">
							<thead>
								<tr class="border-b border-ink text-left text-xs text-muted">
									<th class="py-1.5 pr-3 font-medium">Naziv</th>
									{#each type.groups[0].buckets as bucket (bucket.month)}
										<th class="py-1.5 pr-3 text-right font-medium"
											>{formatMonthLong(bucket.month)}</th
										>
									{/each}
									<th class="py-1.5 text-right font-medium">Ukupno</th>
								</tr>
							</thead>
							<tbody>
								{#each type.rows as row (row.key)}
									<tr class="border-b border-line last:border-0">
										<td class="py-1.5 pr-3">{row.label}</td>
										{#each row.cells as cell, index (type.groups[0].buckets[index].month)}
											<td class="py-1.5 pr-3 text-right tabular-nums">
												{cell ? formatPropertyTotal(cell) : '—'}
											</td>
										{/each}
										<td class="py-1.5 text-right font-semibold tabular-nums">
											{formatPropertyTotal(row.periodTotal)}
										</td>
									</tr>
								{/each}
								{#each type.footers as group (group.label)}
									<tr class="border-t-2 border-ink font-semibold">
										<td class="py-1.5 pr-3">{group.label}</td>
										{#each group.footer.cells as cell, index (type.groups[0].buckets[index].month)}
											<td class="py-1.5 pr-3 text-right tabular-nums">
												{cell.value_count > 0 ? formatPropertyTotal(cell) : '—'}
											</td>
										{/each}
										<td class="py-1.5 text-right tabular-nums"
											>{formatPropertyTotal(group.footer.total)}</td
										>
									</tr>
								{/each}
							</tbody>
						</table>
					</div>
				{/each}
			</section>
		{/if}
	{/each}

	<section class="mt-8">
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
	.dashboard-report-trends {
		display: none;
	}

	@media print {
		:global(body) {
			background: white;
		}

		:global(.app-shell) {
			display: none;
		}

		.dashboard-report-trends {
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

		.dashboard-report-trends tr {
			break-inside: avoid;
		}

		.dashboard-report-trends h2,
		.dashboard-report-trends h3 {
			break-after: avoid;
		}
	}
</style>
