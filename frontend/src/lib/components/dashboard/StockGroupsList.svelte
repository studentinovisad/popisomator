<script lang="ts">
	import { resolve } from '$app/paths';
	import TriangleAlert from '@lucide/svelte/icons/triangle-alert';
	import type { DashboardStockGroup } from '$lib/api';

	let {
		groups,
		showTypeName = true
	}: {
		groups: DashboardStockGroup[];
		showTypeName?: boolean;
	} = $props();

	function stockClass(group: DashboardStockGroup) {
		if (group.in_stock_count === 0) return 'text-danger';
		return group.low ? 'text-warning' : 'text-ink';
	}
</script>

<ul class="divide-y divide-line">
	{#each groups as group (`${group.type_id}-${group.name}`)}
		<li>
			<a
				href={resolve(`/?type_id=${group.type_id}&search=${encodeURIComponent(group.name)}`)}
				class="flex items-center gap-3 px-1 py-2 transition-colors hover:bg-soft"
			>
				{#if group.low}
					<TriangleAlert
						class={`size-4 shrink-0 ${group.in_stock_count === 0 ? 'text-danger' : 'text-warning'}`}
						aria-hidden="true"
					/>
				{/if}
				<div class="min-w-0 flex-1">
					<p class="truncate text-sm text-ink">{group.name}</p>
					{#if showTypeName}
						<p class="truncate text-xs text-muted">{group.type_name}</p>
					{/if}
				</div>
				<!-- The threshold is shown only where the count is short of it; against a healthy group it
				     reads as the wrong half of a fraction. -->
				<p class={`shrink-0 text-sm font-semibold tabular-nums ${stockClass(group)}`}>
					{group.in_stock_count}{#if group.low && group.low_stock_count !== null}<span
							class="font-normal text-muted">/{group.low_stock_count}</span
						>{/if}
				</p>
			</a>
		</li>
	{/each}
</ul>
