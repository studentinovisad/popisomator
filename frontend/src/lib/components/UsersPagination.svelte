<script lang="ts">
	import { Pagination, Portal } from 'bits-ui';

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
</script>

<Portal to="#page-footer-actions">
	<Pagination.Root
		class="flex items-center justify-between gap-4 text-sm"
		count={total}
		{perPage}
		{page}
		onPageChange={onpagechange}
	>
		{#snippet children({ pages })}
			<Pagination.PrevButton
				class="rounded-md border border-line px-1.5 py-1.5 text-sm font-medium text-ink hover:bg-surface disabled:opacity-50"
				disabled={!hasPreviousPage || loading}
			>
				Prethodna
			</Pagination.PrevButton>
			<div class="flex items-center gap-1">
				{#each pages as page (page.key)}
					{#if page.type === 'ellipsis'}
						<span class="px-1 text-muted">…</span>
					{:else}
						<Pagination.Page
							{page}
							class="rounded-md px-3 py-1.5 text-sm hover:bg-surface data-selected:bg-brand data-selected:text-on-brand"
						>
							{page.value}
						</Pagination.Page>
					{/if}
				{/each}
			</div>
			<Pagination.NextButton
				class="rounded-md border border-line px-3 py-1.5 text-sm font-medium text-ink hover:bg-surface disabled:opacity-50"
				disabled={!hasNextPage || loading}
			>
				Sledeća
			</Pagination.NextButton>
		{/snippet}
	</Pagination.Root>
</Portal>
