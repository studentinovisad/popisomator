<script lang="ts">
	import { resolve } from '$app/paths';
	import type { NavigationItem } from '$lib/navigation';
	import NavigationIcon from '$lib/components/NavigationIcon.svelte';

	let {
		items,
		pathname,
		role,
		iconOnly = false,
		iconOnlyOnSmall = false,
		class: className
	}: {
		items: NavigationItem[];
		pathname: string;
		role?: import('$lib/api').UserRole;
		iconOnly?: boolean;
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
		return `flex h-9 w-full items-center rounded-md px-3 text-sm font-medium transition-colors ${
			iconOnly
				? 'justify-center'
				: iconOnlyOnSmall
					? 'gap-3 max-sm:justify-center max-sm:gap-0'
					: 'gap-3'
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
					title={iconOnly ? item.label : undefined}
				>
					<span
						class={`min-w-0 overflow-hidden whitespace-nowrap transition-[max-width,opacity] duration-200 ease-out ${
							iconOnly
								? 'max-w-0 flex-none opacity-0'
								: iconOnlyOnSmall
									? 'max-w-40 flex-1 opacity-100 max-sm:max-w-0 max-sm:flex-none max-sm:opacity-0'
									: 'max-w-40 flex-1 opacity-100'
						}`}
					>
						{item.label}
					</span>
					<span class={iconOnly ? undefined : iconOnlyOnSmall ? 'ml-auto max-sm:ml-0' : 'ml-auto'}>
						<NavigationIcon name={item.icon} />
					</span>
				</a>
			</li>
		{/if}
	{/each}
</ul>
