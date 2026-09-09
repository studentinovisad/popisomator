<script lang="ts">
	import { resolve } from '$app/paths';
	import type { AppPath, NavigationItem } from '$lib/domain/navigation';
	import CountBadge from '$lib/components/shared/CountBadge.svelte';
	import NavigationIcon from '$lib/components/app/NavigationIcon.svelte';

	let {
		items,
		pathname,
		role,
		iconOnly = false,
		iconOnlyOnSmall = false,
		badgeCounts,
		class: className
	}: {
		items: NavigationItem[];
		pathname: string;
		role?: import('$lib/api').UserRole;
		iconOnly?: boolean;
		iconOnlyOnSmall?: boolean;
		// Counts to pin onto an item's icon, keyed by path. Only the paths present here get a badge,
		// so an ordinary link stays an ordinary link; today that is the unread count on /notifications.
		badgeCounts?: Partial<Record<AppPath, number>>;
		class?: string;
	} = $props();

	function isCurrentPath(path: string) {
		return (
			pathname === path ||
			(path === '/' && pathname.startsWith('/items/')) ||
			(path === '/admin/users' && pathname.startsWith('/admin/users/')) ||
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
			{@const badgeCount = badgeCounts?.[item.path] ?? 0}
			<li>
				<a
					class={linkClass(item.path)}
					href={resolve(item.path)}
					aria-current={isCurrentPath(item.path) ? 'page' : undefined}
					aria-label={badgeCount > 0 ? `${item.label}, ${badgeCount} nepročitanih` : undefined}
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
					<span
						class={`relative ${iconOnly ? '' : iconOnlyOnSmall ? 'ml-auto max-sm:ml-0' : 'ml-auto'}`}
					>
						<NavigationIcon name={item.icon} />
						<CountBadge count={badgeCount} class="absolute -top-1.5 -right-2" />
					</span>
				</a>
			</li>
		{/if}
	{/each}
</ul>
