<script lang="ts">
	import type { Snippet } from 'svelte';

	let {
		title,
		subtitle,
		empty = false,
		emptyMessage = 'Nema podataka za izabrani period.',
		// Lists show every row rather than a shortlist, so they scroll inside the card instead of
		// stretching it past whatever sits beside it.
		scroll = false,
		actions,
		children
	}: {
		title: string;
		subtitle?: string;
		empty?: boolean;
		emptyMessage?: string;
		scroll?: boolean;
		actions?: Snippet;
		children: Snippet;
	} = $props();
</script>

<section class="flex min-w-0 flex-col rounded-md border border-line bg-surface">
	<header class="flex items-start justify-between gap-3 border-b border-line px-4 py-3">
		<div class="min-w-0">
			<h2 class="text-sm font-semibold text-ink">{title}</h2>
			{#if subtitle}
				<p class="mt-0.5 text-xs text-muted">{subtitle}</p>
			{/if}
		</div>
		{#if actions}
			<div class="shrink-0">{@render actions()}</div>
		{/if}
	</header>

	<div
		class={`flex min-h-0 flex-1 flex-col p-4 ${scroll && !empty ? 'max-h-80 overflow-y-auto' : ''}`}
	>
		{#if empty}
			<p class="grid flex-1 place-items-center py-8 text-center text-sm text-muted">
				{emptyMessage}
			</p>
		{:else}
			{@render children()}
		{/if}
	</div>
</section>
