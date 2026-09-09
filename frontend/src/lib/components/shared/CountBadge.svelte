<script lang="ts">
	// A small count bubble for sitting on top of an icon. Positioning is left to the caller, which
	// passes it in through `class` - the badge knows what it looks like, the caller knows where it
	// goes - so the same component works pinned to a bell or inline next to a label.
	//
	// Nothing renders at zero, and counts past the cap collapse to "99+" rather than widening the
	// bubble until it unbalances whatever it is pinned to.
	let { count, class: className }: { count: number; class?: string } = $props();

	const maximumDisplayedCount = 99;

	let label = $derived(count > maximumDisplayedCount ? `${maximumDisplayedCount}+` : String(count));
</script>

{#if count > 0}
	<span
		class={`inline-flex h-4 min-w-4 items-center justify-center rounded-full bg-brand px-1 font-mono text-[0.625rem] leading-none font-medium text-on-brand ${className ?? ''}`}
		aria-hidden="true"
	>
		{label}
	</span>
{/if}
