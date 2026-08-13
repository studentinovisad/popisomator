<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { resolve } from '$app/paths';
	import favicon from '$lib/assets/favicon.svg';
	import { api } from '$lib/api';
	import AccountLink from '$lib/components/AccountLink.svelte';
	import NavigationLinks from '$lib/components/NavigationLinks.svelte';
	import UserAvatar from '$lib/components/UserAvatar.svelte';
	import { getPageMetadata, primaryNavigation, secondaryNavigation } from '$lib/navigation';
	import { session } from '$lib/session.svelte';
	import { theme } from '$lib/theme.svelte';
	import { Button, Popover } from 'bits-ui';
	import '../app.css';

	let { children } = $props();
	let popoverSide = $state<'bottom' | 'right'>('bottom');
	let tabHelpOpen = $state(false);
	let popoverCollisionPadding = $state(24);
	let activePage = $derived(getPageMetadata(page.url.pathname));

	onMount(() => {
		void session.refresh();
		theme.initialize();

		const desktopMedia = window.matchMedia('(min-width: 48rem)');
		const updatePopoverSide = () => {
			popoverSide = desktopMedia.matches ? 'right' : 'bottom';
			popoverCollisionPadding = desktopMedia.matches ? 0 : 24;
		};

		updatePopoverSide();
		desktopMedia.addEventListener('change', updatePopoverSide);

		return () => desktopMedia.removeEventListener('change', updatePopoverSide);
	});

	$effect(() => {
		if (popoverSide === 'bottom') {
			void page.url.pathname;
			tabHelpOpen = false;
		}
	});

	$effect(() => {
		const pathname = page.url.pathname;
		if (
			session.ready &&
			!session.user &&
			pathname !== '/login' &&
			pathname !== '/register' &&
			pathname !== '/settings'
		) {
			void goto(resolve('/login'));
		}
	});

	async function logout() {
		session.clear();

		try {
			await api.logout();
		} finally {
			await goto(resolve('/login'));
		}
	}
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<div
	class="flex min-h-svh pb-16 text-ink max-sm:pb-[calc(4rem+env(safe-area-inset-bottom))] md:pb-0"
>
	<aside class="hidden w-60 shrink-0 flex-col border-r border-chrome-line bg-chrome md:flex">
		<div class="p-5">
			<a class="block text-on-chrome" href={resolve('/')}>
				<span class="font-semibold tracking-tight">Popisomator</span>
			</a>

			<nav class="mt-10" aria-label="Glavna navigacija">
				<NavigationLinks
					items={session.user ? primaryNavigation : []}
					pathname={page.url.pathname}
					role={session.user?.role}
					class="mt-2 space-y-1"
				/>
			</nav>
		</div>

		<div class="mt-auto p-4">
			<NavigationLinks
				items={secondaryNavigation}
				pathname={page.url.pathname}
				role={session.user?.role}
				class="space-y-1"
			/>
			{#if session.user}
				<div class="mt-4 flex w-full">
					<AccountLink user={session.user} fullWidth active={page.url.pathname === '/account'} />
					<Button.Root
						class="inline-flex h-9 shrink-0 items-center justify-center rounded-r-md bg-chrome pr-3 pl-2.5 text-chrome-muted transition-colors hover:bg-on-chrome/10 hover:text-on-chrome"
						aria-label="Odjavi se"
						onclick={() => void logout()}
					>
						<svg
							aria-hidden="true"
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							stroke-width="2"
							class="size-4"
						>
							<path d="M21 19V5a2 2 0 0 0-2-2h-6" />
							<path d="M10 17l-5-5 5-5" />
							<path d="M5 12h11" />
						</svg>
					</Button.Root>
				</div>
			{:else}
				<a
					class={`mt-4 flex h-9 w-full items-center justify-between rounded-md px-3 text-sm font-medium transition-colors ${
						page.url.pathname === '/login'
							? 'bg-brand-soft text-brand'
							: 'bg-chrome text-chrome-muted hover:bg-on-chrome/10 hover:text-on-chrome'
					}`}
					href={resolve('/login')}
					aria-current={page.url.pathname === '/login' ? 'page' : undefined}
				>
					<span>Prijava</span>
					<svg
						aria-hidden="true"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						class="size-4"
					>
						<path d="M10 17l5-5-5-5" />
						<path d="M15 12H3" />
						<path d="M21 19V5a2 2 0 0 0-2-2h-6" />
					</svg>
				</a>
			{/if}
		</div>
	</aside>

	<div class="flex min-w-0 flex-1 flex-col">
		<header class="border-b border-chrome-line bg-chrome text-on-chrome md:hidden">
			<div class="flex h-16 items-center justify-between px-4">
				<a class="font-semibold tracking-tight" href={resolve('/')}>Popisomator</a>
				<div class="flex items-center gap-2">
					{#if session.user}
						<a
							class={`inline-flex size-9 items-center justify-center rounded-md transition-colors ${
								page.url.pathname === '/account'
									? 'bg-brand-soft'
									: 'bg-chrome hover:bg-on-chrome/10'
							}`}
							href={resolve('/account')}
							aria-label={`Moj nalog: ${session.user.full_name}`}
							aria-current={page.url.pathname === '/account' ? 'page' : undefined}
						>
							<UserAvatar name={session.user.full_name} class="size-7" />
						</a>
						<Button.Root
							class="inline-flex h-9 w-9 items-center justify-center rounded-md bg-chrome pr-3 pl-2 text-chrome-muted transition-colors hover:bg-on-chrome/10 hover:text-on-chrome"
							aria-label="Odjavi se"
							onclick={() => void logout()}
						>
							<svg
								aria-hidden="true"
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="2"
								class="size-4"
							>
								<path d="M21 19V5a2 2 0 0 0-2-2h-6" />
								<path d="M10 17l-5-5 5-5" />
								<path d="M5 12h11" />
							</svg>
						</Button.Root>
					{:else}
						<a
							class={`inline-flex h-9 items-center gap-2 rounded-md px-3 text-sm font-medium transition-colors ${
								page.url.pathname === '/login'
									? 'bg-brand-soft text-brand'
									: 'bg-chrome text-chrome-muted hover:bg-on-chrome/10 hover:text-on-chrome'
							}`}
							href={resolve('/login')}
							aria-current={page.url.pathname === '/login' ? 'page' : undefined}
						>
							<span>Prijava</span>
							<svg
								aria-hidden="true"
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="2"
								class="size-4"
							>
								<path d="M10 17l5-5-5-5" />
								<path d="M15 12H3" />
								<path d="M21 19V5a2 2 0 0 0-2-2h-6" />
							</svg>
						</a>
					{/if}
				</div>
			</div>
		</header>

		<header class="h-16 shrink-0 border-b border-chrome-line bg-chrome text-brand">
			<div class="flex h-full items-center justify-between gap-4 px-4 sm:px-6">
				<div class="flex items-center gap-2">
					<h1 class="flex items-center gap-3 text-base leading-none font-semibold tracking-tight">
						{activePage.title}
					</h1>
					<Popover.Root bind:open={tabHelpOpen}>
						<Popover.Trigger
							class="relative -top-px grid size-5 place-items-center text-base leading-none text-chrome-muted hover:text-on-chrome"
							aria-label={`Pomoć za ${activePage.title}`}
						>
							🛈
						</Popover.Trigger>
						<Popover.Portal>
							<Popover.Content
								side={popoverSide}
								sideOffset={8}
								align={popoverSide === 'right' ? 'start' : 'center'}
								alignOffset={popoverSide === 'right' ? -8 : 0}
								avoidCollisions={true}
								collisionPadding={popoverCollisionPadding}
								sticky="always"
								strategy="fixed"
								class="z-50 w-[min(20rem,calc(100dvw-3rem))] rounded-md border border-line bg-surface px-3 py-2 text-sm leading-snug text-ink shadow-lg shadow-black/15"
							>
								{activePage.description}
							</Popover.Content>
						</Popover.Portal>
					</Popover.Root>
				</div>
				<div id="page-header-actions" class="page-header-actions"></div>
			</div>
		</header>

		<div class="app-content flex-1">
			{@render children()}
		</div>

		<footer
			class="page-footer h-16 shrink-0 border-t border-chrome-line bg-chrome text-chrome-muted"
		>
			<div id="page-footer-actions" class="page-footer-actions"></div>
		</footer>
	</div>

	<nav
		class="fixed inset-x-0 bottom-0 z-10 border-t border-chrome-line bg-chrome pb-[env(safe-area-inset-bottom)] md:hidden"
		aria-label="Glavna navigacija"
	>
		<NavigationLinks
			items={session.user ? [...primaryNavigation, ...secondaryNavigation] : secondaryNavigation}
			pathname={page.url.pathname}
			role={session.user?.role}
			iconOnlyOnSmall
			class="flex h-16 items-center justify-center gap-4 overflow-x-auto px-4 text-sm max-sm:justify-around max-sm:gap-2"
		/>
	</nav>
</div>
