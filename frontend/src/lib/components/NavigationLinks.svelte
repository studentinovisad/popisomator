<script lang="ts">
	import { resolve } from '$app/paths';
	import type { NavigationItem } from '$lib/navigation';
	import NavigationIcon from '$lib/components/NavigationIcon.svelte';

	let {
		items,
		pathname,
		role,
		iconOnlyOnSmall = false,
		class: className
	}: {
		items: NavigationItem[];
		pathname: string;
		role?: import('$lib/api').UserRole;
		iconOnlyOnSmall?: boolean;
		class?: string;
	} = $props();

	function isCurrentPath(path: string) {
		return (
			pathname === path ||
			(path === '/' && pathname.startsWith('/items/')) ||
			(path === '/catalog/item-types' && pathname.startsWith('/catalog/'))
		);
	}

	function linkClass(path: string) {
		return `flex items-center justify-between gap-3 rounded-md px-3 py-2 text-sm font-medium transition-colors ${
			iconOnlyOnSmall ? 'max-sm:size-9 max-sm:justify-center max-sm:p-0' : ''
		} ${
			isCurrentPath(path)
				? 'bg-brand-soft text-brand'
				: 'text-chrome-muted hover:bg-on-chrome/10 hover:text-on-chrome'
		}`;
	}
</script>

<ul class={className}>
	{#each items as item (item.path)}
		{#if !item.requiredRoles || (role && item.requiredRoles.includes(role))}
			<li>
				<a
					class={linkClass(item.path)}
					href={resolve(item.path)}
					aria-current={isCurrentPath(item.path) ? 'page' : undefined}
				>
					<span class={iconOnlyOnSmall ? 'max-sm:sr-only' : undefined}>{item.label}</span>
					<NavigationIcon name={item.icon} />
				</a>
			</li>
		{/if}
	{/each}
</ul>
