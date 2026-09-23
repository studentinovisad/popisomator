<script lang="ts">
	import ChevronDown from '@lucide/svelte/icons/chevron-down';
	import { Collapsible } from 'bits-ui';
	import type { DashboardTypeConsumptionQuantity } from '$lib/api';
	import {
		consumptionPropertyRows,
		formatMonthLong,
		formatPropertyTotal
	} from '$lib/domain/dashboard';

	let { types }: { types: DashboardTypeConsumptionQuantity[] } = $props();
</script>

<div class="space-y-2">
	{#each types as type (type.type_id)}
		<Collapsible.Root open={types.length === 1} class="rounded-md border border-line">
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
							<th class="px-3 py-2 text-left font-medium">Svojstvo</th>
							{#each type.buckets as bucket (bucket.month)}
								<th class="px-3 py-2 text-right font-medium">{formatMonthLong(bucket.month)}</th>
							{/each}
							<th class="px-3 py-2 text-right font-medium">Ukupno</th>
						</tr>
					</thead>
					<tbody>
						{#each consumptionPropertyRows(type) as row (row.propertyID)}
							<tr class="border-b border-line last:border-0">
								<td class="px-3 py-2 text-ink">{row.periodTotal.property_name}</td>
								{#each row.cells as cell, index (type.buckets[index].month)}
									<td class="px-3 py-2 text-right text-ink tabular-nums">
										{cell ? formatPropertyTotal(cell) : '—'}
									</td>
								{/each}
								<td class="px-3 py-2 text-right font-semibold text-ink tabular-nums">
									{formatPropertyTotal(row.periodTotal)}
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</Collapsible.Content>
		</Collapsible.Root>
	{/each}
</div>
