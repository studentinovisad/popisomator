<script lang="ts">
	import { resolve } from '$app/paths';
	import Bell from '@lucide/svelte/icons/bell';
	import { notifications } from '$lib/state/notifications.svelte';

	let {
		variant,
		iconOnly = false,
		active = false
	}: { variant: 'sidebar' | 'header'; iconOnly?: boolean; active?: boolean } = $props();

	const label = 'Obaveštenja';

	let unreadCount = $derived(notifications.unreadCount);
	let badgeLabel = $derived(unreadCount > 99 ? ':D' : String(unreadCount));
	let ariaLabel = $derived(unreadCount > 0 ? `${label}, ${unreadCount} nepročitanih` : label);

	// The sidebar bell sits among the navigation links, so it mirrors their sizing and the
	// label collapse animation; the header bell matches the mobile account button instead.
	let linkClass = $derived(
		variant === 'sidebar'
			? `flex h-9 w-full items-center rounded-md px-3 text-sm font-medium transition-colors ${
					iconOnly ? 'justify-center' : 'gap-3'
				} ${
					active
						? 'bg-brand-soft text-brand'
						: 'text-chrome-muted hover:bg-on-chrome/10 hover:text-on-chrome'
				}`
			: `inline-flex size-9 items-center justify-center rounded-md transition-colors ${
					active
						? 'bg-brand-soft text-brand'
						: 'bg-chrome text-chrome-muted hover:bg-on-chrome/10 hover:text-on-chrome'
				}`
	);
</script>

<a
	class={linkClass}
	href={resolve('/notifications')}
	aria-label={ariaLabel}
	aria-current={active ? 'page' : undefined}
	title={variant === 'header' || iconOnly ? label : undefined}
>
	{#if variant === 'sidebar'}
		<span
			class={`min-w-0 overflow-hidden whitespace-nowrap transition-[max-width,opacity] duration-200 ease-out ${
				iconOnly ? 'max-w-0 flex-none opacity-0' : 'max-w-40 flex-1 opacity-100'
			}`}
		>
			{label}
		</span>
	{/if}
	<span class={`relative ${variant === 'sidebar' && !iconOnly ? 'ml-auto' : ''}`}>
		<Bell class="size-4 shrink-0" aria-hidden="true" />
		{#if unreadCount > 0}
			<span
				class="absolute -top-1.5 -right-2 inline-flex h-4 min-w-4 items-center justify-center rounded-full bg-brand px-1 font-mono text-[0.625rem] leading-none font-medium text-on-brand"
				aria-hidden="true"
			>
				{badgeLabel}
			</span>
		{/if}
	</span>
</a>
