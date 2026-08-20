<script lang="ts">
	import ChevronLeft from '@lucide/svelte/icons/chevron-left';
	import ChevronRight from '@lucide/svelte/icons/chevron-right';
	import { Pagination, Portal } from 'bits-ui';
	import { onMount } from 'svelte';

	let {
		total,
		perPage,
		page,
		hasPreviousPage,
		hasNextPage,
		loading,
		onpagechange
	}: {
		total: number;
		perPage: number;
		page: number;
		hasPreviousPage: boolean;
		hasNextPage: boolean;
		loading: boolean;
		onpagechange: (page: number) => void;
	} = $props();

	let siblingCount = $state(1);

	onMount(() => {
		const desktopQuery = window.matchMedia('(min-width: 64rem)');
		const updateSiblingCount = () => {
			siblingCount = desktopQuery.matches ? 4 : 1;
		};

		updateSiblingCount();
		desktopQuery.addEventListener('change', updateSiblingCount);

		return () => desktopQuery.removeEventListener('change', updateSiblingCount);
	});
</script>

<Portal to="#page-footer-actions">
	<Pagination.Root
		class="grid grid-cols-[minmax(0,1fr)_auto_minmax(0,1fr)] items-center gap-4 text-sm"
		count={total}
		{perPage}
		{page}
		{siblingCount}
		onPageChange={onpagechange}
	>
		{#snippet children({ pages })}
			<div>
				<Pagination.PrevButton
					class="inline-flex size-8 items-center justify-center rounded-md border border-line text-ink hover:bg-surface disabled:opacity-50 sm:size-9"
					disabled={!hasPreviousPage || loading}
					aria-label="Prethodna stranica"
				>
					<ChevronLeft class="size-4" aria-hidden="true" />
				</Pagination.PrevButton>
			</div>
			<div class="flex items-center gap-1">
				{#each pages as page (page.key)}
					{#if page.type === 'ellipsis'}
						<span class="px-1 text-muted">…</span>
					{:else}
						<Pagination.Page
							{page}
							class="rounded-md px-2 py-1 text-xs hover:bg-surface data-selected:bg-brand data-selected:text-on-brand sm:px-3 sm:py-1.5 sm:text-sm"
						>
							{page.value}
						</Pagination.Page>
					{/if}
				{/each}
			</div>
			<div class="flex justify-end">
				<Pagination.NextButton
					class="inline-flex size-8 items-center justify-center rounded-md border border-line text-ink hover:bg-surface disabled:opacity-50 sm:size-9"
					disabled={!hasNextPage || loading}
					aria-label="Sledeća stranica"
				>
					<ChevronRight class="size-4" aria-hidden="true" />
				</Pagination.NextButton>
			</div>
		{/snippet}
	</Pagination.Root>
</Portal>
