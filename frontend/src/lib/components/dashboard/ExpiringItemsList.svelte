<script lang="ts">
	import { resolve } from '$app/paths';
	import CalendarClock from '@lucide/svelte/icons/calendar-clock';
	import CalendarX from '@lucide/svelte/icons/calendar-x';
	import type { DashboardExpiringItem } from '$lib/api';
	import { expiryDelayLabel, formatExpiryDate } from '$lib/domain/dashboard';

	let { items }: { items: DashboardExpiringItem[] } = $props();
</script>

<ul class="divide-y divide-line">
	{#each items as item (item.id)}
		{@const expired = item.days_remaining < 0}
		<li>
			<a
				href={resolve(`/items/${item.id}`)}
				class="flex items-center gap-3 px-1 py-2 transition-colors hover:bg-soft"
			>
				{#if expired}
					<CalendarX class="size-4 shrink-0 text-danger" aria-hidden="true" />
				{:else}
					<CalendarClock class="size-4 shrink-0 text-warning" aria-hidden="true" />
				{/if}
				<div class="min-w-0 flex-1">
					<p class="truncate text-sm text-ink">{item.name}</p>
					<p class="truncate text-xs text-muted">{formatExpiryDate(item.expires_on)}</p>
				</div>
				<p
					class={`shrink-0 text-xs font-medium whitespace-nowrap ${expired ? 'text-danger' : 'text-warning'}`}
				>
					{expiryDelayLabel(item.days_remaining)}
				</p>
			</a>
		</li>
	{/each}
</ul>
