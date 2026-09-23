<script lang="ts">
	import ChevronDown from '@lucide/svelte/icons/chevron-down';
	import { Collapsible } from 'bits-ui';
	import type { DashboardTypeConsumptionQuantity, PropertyValueType } from '$lib/api';
	import {
		consumptionGroupRows,
		formatMonthLong,
		formatPropertyTotal,
		sortGroupsByTotal,
		sumRowsFooters
	} from '$lib/domain/dashboard';

	let {
		types,
		valueTypes
	}: {
		types: DashboardTypeConsumptionQuantity[];
		valueTypes: PropertyValueType[];
	} = $props();

	const visibleTypes = $derived(
		types
			.map((type) => {
				const groups = sortGroupsByTotal(
					type.groups.filter((group) => consumptionGroupRows(group, valueTypes).length > 0),
					valueTypes
				);
				const rows = groups.flatMap((group) => consumptionGroupRows(group, valueTypes));
				return { ...type, groups, rows, footers: sumRowsFooters(rows) };
			})
			.filter((type) => type.groups.length > 0)
	);
</script>

<div class="space-y-2">
	{#each visibleTypes as type (type.type_id)}
		<Collapsible.Root open={visibleTypes.length === 1} class="rounded-md border border-line">
			<Collapsible.Trigger
				class="flex w-full items-center justify-between gap-2 px-3 py-2 text-left text-sm font-medium text-ink transition-colors hover:bg-soft"
			>
				<span>{type.type_name}</span>
				<ChevronDown class="size-4 shrink-0 text-muted" aria-hidden="true" />
			</Collapsible.Trigger>
			<Collapsible.Content class="overflow-x-auto border-t border-line">
				<table class="w-full text-sm">
					<thead>
						<tr class="border-b border-line text-xs text-muted">
							<th class="px-3 py-2 text-left font-medium">Naziv</th>
							{#each type.groups[0].buckets as bucket (bucket.month)}
								<th class="px-3 py-2 text-right font-medium">{formatMonthLong(bucket.month)}</th>
							{/each}
							<th class="px-3 py-2 text-right font-medium">Ukupno</th>
						</tr>
					</thead>
					<tbody>
						{#each type.rows as row (row.key)}
							<tr class="border-b border-line last:border-0">
								<td class="px-3 py-2 text-ink">{row.label}</td>
								{#each row.cells as cell, index (type.groups[0].buckets[index].month)}
									<td class="px-3 py-2 text-right text-ink tabular-nums">
										{cell ? formatPropertyTotal(cell) : '—'}
									</td>
								{/each}
								<td class="px-3 py-2 text-right font-semibold text-ink tabular-nums">
									{formatPropertyTotal(row.periodTotal)}
								</td>
							</tr>
						{/each}
						{#each type.footers as group (group.label)}
							<tr class="border-t-2 border-line font-semibold text-ink">
								<td class="px-3 py-2">{group.label}</td>
								{#each group.footer.cells as cell, index (type.groups[0].buckets[index].month)}
									<td class="px-3 py-2 text-right tabular-nums">
										{cell.value_count > 0 ? formatPropertyTotal(cell) : '—'}
									</td>
								{/each}
								<td class="px-3 py-2 text-right tabular-nums"
									>{formatPropertyTotal(group.footer.total)}</td
								>
							</tr>
						{/each}
					</tbody>
				</table>
			</Collapsible.Content>
		</Collapsible.Root>
	{/each}
</div>
