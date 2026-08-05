<script lang="ts">
	import type { NavigationIconName } from '$lib/navigation';

	let { name }: { name: NavigationIconName } = $props();
</script>

<svg
	aria-hidden="true"
	viewBox="0 0 24 24"
	fill="none"
	stroke="currentColor"
	stroke-width="2"
	class="size-4 shrink-0"
>
	{#if name === 'overview'}
		<rect x="3" y="3" width="7" height="7" rx="1" />
		<rect x="14" y="3" width="7" height="7" rx="1" />
		<rect x="3" y="14" width="7" height="7" rx="1" />
		<rect x="14" y="14" width="7" height="7" rx="1" />
	{:else if name === 'users'}
		<circle cx="9" cy="8" r="3" />
		<path d="M3 21v-2a4 4 0 0 1 4-4h4a4 4 0 0 1 4 4v2" />
		<path d="M16 5a3 3 0 0 1 0 6M21 21v-2a4 4 0 0 0-3-3.87" />
	{:else}
		<circle cx="12" cy="12" r="3" />
		<path
			d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06-2.12 2.12-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V20h-3v-.08a1.65 1.65 0 0 0-1-1.51 1.65 1.65 0 0 0-1.82.33l-.06.06-2.12-2.12.06-.06A1.65 1.65 0 0 0 7.18 15a1.65 1.65 0 0 0-1.51-1H5.6v-3h.08a1.65 1.65 0 0 0 1.51-1 1.65 1.65 0 0 0-.33-1.82L6.8 8.12 8.92 6l.06.06a1.65 1.65 0 0 0 1.82.33h.01a1.65 1.65 0 0 0 1-1.51V4.8h3v.08a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33L17.7 6l2.12 2.12-.06.06a1.65 1.65 0 0 0-.33 1.82v.01a1.65 1.65 0 0 0 1.51 1h.08v3h-.08a1.65 1.65 0 0 0-1.54 1Z"
		/>
	{/if}
</svg>
