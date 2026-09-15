<script lang="ts">
	import ArrowDownNarrowWide from '@lucide/svelte/icons/arrow-down-narrow-wide';
	import ArrowUpNarrowWide from '@lucide/svelte/icons/arrow-up-narrow-wide';
	import type { DashboardStockGroup } from '$lib/api';
	import DashboardCard from '$lib/components/dashboard/DashboardCard.svelte';
	import StockGroupsList from '$lib/components/dashboard/StockGroupsList.svelte';

	let {
		groups,
		showTypeName = true
	}: {
		groups: DashboardStockGroup[];
		showTypeName?: boolean;
	} = $props();

	// Scarcest first, because a shortage is the thing worth acting on and this way it leads. Flipping
	// the order turns the same list into a plain view of what is actually on the shelf; nothing is
	// filtered out either way.
	let scarcestFirst = $state(true);

	const lowCount = $derived(groups.filter((group) => group.low).length);
	const sortedGroups = $derived(
		[...groups].sort((left, right) => {
			const byCount = scarcestFirst
				? left.in_stock_count - right.in_stock_count
				: right.in_stock_count - left.in_stock_count;
			return byCount !== 0 ? byCount : left.name.localeCompare(right.name, 'sr-RS');
		})
	);

	const subtitle = $derived(
		lowCount > 0
			? `${lowCount} od ${groups.length} grupa je na pragu ili ispod njega.`
			: `${groups.length} grupa, nijedna nije pri kraju zaliha.`
	);
</script>

<DashboardCard
	title="Zalihe"
	{subtitle}
	scroll
	empty={groups.length === 0}
	emptyMessage="Nema nijedne stavke."
>
	{#snippet actions()}
		<button
			type="button"
			onclick={() => (scarcestFirst = !scarcestFirst)}
			aria-label={scarcestFirst ? 'Prikaži najveće zalihe prvo' : 'Prikaži najmanje zalihe prvo'}
			class="inline-flex items-center gap-1.5 rounded-md border border-line px-2 py-1 text-xs whitespace-nowrap text-muted transition-colors hover:border-brand/40 hover:text-ink"
		>
			{#if scarcestFirst}
				<ArrowUpNarrowWide class="size-3.5 shrink-0" aria-hidden="true" />
				Najmanje prvo
			{:else}
				<ArrowDownNarrowWide class="size-3.5 shrink-0" aria-hidden="true" />
				Najviše prvo
			{/if}
		</button>
	{/snippet}

	<StockGroupsList groups={sortedGroups} {showTypeName} />
</DashboardCard>
