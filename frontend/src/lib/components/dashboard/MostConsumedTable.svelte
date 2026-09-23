<script lang="ts">
	import type { DashboardMostConsumedGroup } from '$lib/api';
	import { formatPropertyTotals } from '$lib/domain/dashboard';

	let {
		groups,
		showTypeName = true
	}: {
		groups: DashboardMostConsumedGroup[];
		showTypeName?: boolean;
	} = $props();
</script>

<div class="overflow-x-auto">
	<table class="w-full text-sm">
		<thead>
			<tr class="border-b border-line text-xs text-muted">
				<th class="w-8 px-3 py-2 text-left font-medium">#</th>
				<th class="px-3 py-2 text-left font-medium">Naziv</th>
				<th class="px-3 py-2 text-right font-medium">Potrošeno</th>
				<th class="px-3 py-2 text-right font-medium">Količina</th>
			</tr>
		</thead>
		<tbody>
			{#each groups as group, index (`${group.type_id}-${group.name}`)}
				<tr class="border-b border-line last:border-0">
					<td class="px-3 py-2 text-muted tabular-nums">{index + 1}</td>
					<td class="px-3 py-2 text-ink">
						<p class="truncate">{group.name}</p>
						{#if showTypeName}
							<p class="truncate text-xs text-muted">{group.type_name}</p>
						{/if}
					</td>
					<td class="px-3 py-2 text-right font-semibold text-ink tabular-nums"
						>{group.consumed_count}</td
					>
					<td class="px-3 py-2 text-right text-ink tabular-nums">
						{group.totals.length > 0 ? formatPropertyTotals(group.totals) : '—'}
					</td>
				</tr>
			{/each}
		</tbody>
	</table>
</div>
