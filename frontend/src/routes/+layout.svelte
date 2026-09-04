<script lang="ts">
	import { onMount, setContext } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { resolve } from '$app/paths';
	import favicon from '$lib/assets/favicon.svg';
	import ChevronLeft from '@lucide/svelte/icons/chevron-left';
	import ChevronRight from '@lucide/svelte/icons/chevron-right';
	import CircleCheck from '@lucide/svelte/icons/circle-check';
	import CircleX from '@lucide/svelte/icons/circle-x';
	import Info from '@lucide/svelte/icons/info';
	import LogIn from '@lucide/svelte/icons/log-in';
	import LogOut from '@lucide/svelte/icons/log-out';
	import TriangleAlert from '@lucide/svelte/icons/triangle-alert';
	import X from '@lucide/svelte/icons/x';
	import { api, type ItemRequestPreparationReport } from '$lib/api';
	import AccountLink from '$lib/components/app/AccountLink.svelte';
	import NavigationLinks from '$lib/components/app/NavigationLinks.svelte';
	import UserAvatar from '$lib/components/app/UserAvatar.svelte';
	import PreparationReport from '$lib/components/admin/ItemRequestPreparationReport.svelte';
	import { getPageMetadata, primaryNavigation, secondaryNavigation } from '$lib/domain/navigation';
	import { session } from '$lib/state/session.svelte';
	import {
		preparationReportPrintContextKey,
		type PreparationReportPrintContext
	} from '$lib/state/preparation-report-print-context';
	import { theme } from '$lib/state/theme.svelte';
	import { Button, Collapsible, Popover, ScrollArea } from 'bits-ui';
	import { Toaster } from 'svelte-sonner';
	import '../app.css';

	let { children, data } = $props();
	let popoverSide = $state<'bottom' | 'right'>('bottom');
	let tabHelpOpen = $state(false);
	let popoverCollisionPadding = $state(24);
	let sidebarExpanded = $state(true);
	let sessionHydrated = $state(false);
	let currentUser = $derived(sessionHydrated ? session.user : data.currentUser);
	let preparationReport = $state<ItemRequestPreparationReport | null>(null);

	// svelte-ignore state_referenced_locally
	if (!data.sidebarExpanded) {
		sidebarExpanded = false;
	}

	let activePage = $derived(getPageMetadata(page.url.pathname));

	setContext<PreparationReportPrintContext>(preparationReportPrintContextKey, {
		setPreparationReport: (report) => {
			preparationReport = report;
		},
		print: () => window.print()
	});

	onMount(() => {
		if (data.currentUser) {
			session.setUser(data.currentUser);
		} else {
			void session.refresh();
		}
		sessionHydrated = true;
		theme.initialize();
		localStorage.setItem('popisomator-sidebar', sidebarExpanded ? 'expanded' : 'collapsed');
		const desktopMedia = window.matchMedia('(min-width: 48rem)');
		const updatePopoverSide = () => {
			popoverSide = desktopMedia.matches ? 'right' : 'bottom';
			popoverCollisionPadding = desktopMedia.matches ? 0 : 24;
		};

		updatePopoverSide();
		desktopMedia.addEventListener('change', updatePopoverSide);

		return () => {
			desktopMedia.removeEventListener('change', updatePopoverSide);
		};
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

	function setSidebarExpanded(expanded: boolean) {
		sidebarExpanded = expanded;
		localStorage.setItem('popisomator-sidebar', expanded ? 'expanded' : 'collapsed');
		document.cookie = `popisomator-sidebar=${expanded ? 'expanded' : 'collapsed'}; Path=/; Max-Age=31536000; SameSite=Lax`;
	}
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<div
	class="app-shell flex h-svh overflow-hidden pb-16 text-ink max-sm:pb-[calc(4rem+env(safe-area-inset-bottom))] lg:pb-0"
>
	<Collapsible.Root
		open={sidebarExpanded}
		onOpenChange={setSidebarExpanded}
		class={`app-sidebar hidden h-full shrink-0 transition-[width] duration-200 ease-out lg:block ${
			sidebarExpanded ? 'w-60' : 'w-16'
		}`}
	>
		<aside class="flex h-full w-full flex-col border-r border-chrome-line bg-chrome">
			<div class="px-3">
				<div class="flex h-16 items-center">
					<div
						class={`flex h-9 w-full items-center ${sidebarExpanded ? 'justify-between gap-3' : 'justify-center'}`}
					>
						<a
							class={`min-w-0 overflow-hidden font-semibold tracking-tight whitespace-nowrap text-on-chrome transition-[max-width,opacity] duration-200 ease-out ${
								sidebarExpanded ? 'max-w-32 opacity-100' : 'max-w-0 opacity-0'
							}`}
							href={resolve('/')}
						>
							Popisomator
						</a>
						<Collapsible.Trigger
							class="inline-grid size-9 shrink-0 place-items-center rounded-md text-chrome-muted transition-colors hover:bg-on-chrome/10 hover:text-on-chrome"
							aria-label={sidebarExpanded ? 'Skupi bočnu navigaciju' : 'Raširi bočnu navigaciju'}
							title={sidebarExpanded ? 'Skupi bočnu navigaciju' : 'Raširi bočnu navigaciju'}
						>
							{#if sidebarExpanded}
								<ChevronLeft class="size-4" aria-hidden="true" />
							{:else}
								<ChevronRight class="size-4" aria-hidden="true" />
							{/if}
						</Collapsible.Trigger>
					</div>
				</div>

				<nav class="mt-10" aria-label="Glavna navigacija">
					<NavigationLinks
						items={currentUser ? primaryNavigation : []}
						pathname={page.url.pathname}
						role={currentUser?.role}
						iconOnly={!sidebarExpanded}
						class="mt-2 space-y-1"
					/>
				</nav>
			</div>

			<div class={`mt-auto px-3 ${sidebarExpanded ? 'py-4' : 'pt-3 pb-4'}`}>
				<NavigationLinks
					items={secondaryNavigation}
					pathname={page.url.pathname}
					role={currentUser?.role}
					iconOnly={!sidebarExpanded}
					class="space-y-1"
				/>
				{#if currentUser}
					<div class={`mt-4 flex w-full ${sidebarExpanded ? '' : 'flex-col items-center gap-1'}`}>
						<AccountLink
							user={currentUser}
							fullWidth={sidebarExpanded}
							iconOnly={!sidebarExpanded}
							active={page.url.pathname === '/account'}
						/>
						<Button.Root
							class={`inline-flex h-9 shrink-0 items-center justify-center bg-chrome text-chrome-muted transition-colors hover:bg-on-chrome/10 hover:text-on-chrome ${
								sidebarExpanded ? 'rounded-r-md px-3' : 'w-9 rounded-md'
							}`}
							aria-label="Odjavi se"
							title={!sidebarExpanded ? 'Odjavi se' : undefined}
							onclick={() => void logout()}
						>
							<LogOut class="size-4" aria-hidden="true" />
						</Button.Root>
					</div>
				{:else}
					<a
						class={`mt-4 flex h-9 items-center rounded-md text-sm font-medium transition-colors ${
							sidebarExpanded ? 'w-full justify-between px-3' : 'w-full justify-center px-3'
						} ${
							page.url.pathname === '/login'
								? 'bg-brand-soft text-brand'
								: 'bg-chrome text-chrome-muted hover:bg-on-chrome/10 hover:text-on-chrome'
						}`}
						href={resolve('/login')}
						aria-current={page.url.pathname === '/login' ? 'page' : undefined}
						title={!sidebarExpanded ? 'Prijava' : undefined}
					>
						<span
							class={`overflow-hidden whitespace-nowrap transition-[max-width,opacity] duration-200 ease-out ${
								sidebarExpanded ? 'max-w-24 opacity-100' : 'max-w-0 opacity-0'
							}`}
						>
							Prijava
						</span>
						<LogIn class="size-4" aria-hidden="true" />
					</a>
				{/if}
			</div>
		</aside>
	</Collapsible.Root>

	<div class="flex min-h-0 min-w-0 flex-1 flex-col">
		<header class="border-b border-chrome-line bg-chrome text-on-chrome lg:hidden">
			<div class="flex h-16 items-center justify-between px-4">
				<a class="font-semibold tracking-tight" href={resolve('/')}>Popisomator</a>
				<div class="flex items-center gap-2">
					{#if currentUser}
						<a
							class={`inline-flex size-9 items-center justify-center rounded-md transition-colors ${
								page.url.pathname === '/account'
									? 'bg-brand-soft'
									: 'bg-chrome hover:bg-on-chrome/10'
							}`}
							href={resolve('/account')}
							aria-label={`Moj nalog: ${currentUser.full_name}`}
							aria-current={page.url.pathname === '/account' ? 'page' : undefined}
						>
							<UserAvatar name={currentUser.full_name} class="size-7" />
						</a>
						<Button.Root
							class="inline-flex h-9 w-9 items-center justify-center rounded-md bg-chrome pr-3 pl-2 text-chrome-muted transition-colors hover:bg-on-chrome/10 hover:text-on-chrome"
							aria-label="Odjavi se"
							onclick={() => void logout()}
						>
							<LogOut class="size-4" aria-hidden="true" />
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
							<LogIn class="size-4" aria-hidden="true" />
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
							<Info class="mt-0.5 size-4" strokeWidth={2} aria-hidden="true" />
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

		<ScrollArea.Root class="app-content min-h-0 flex-1 overflow-hidden" type="auto">
			<ScrollArea.Viewport class="h-full w-full">
				<div class="min-h-full">{@render children()}</div>
			</ScrollArea.Viewport>
			<ScrollArea.Scrollbar class="flex w-2.5 touch-none bg-soft p-0.5" orientation="vertical">
				<ScrollArea.Thumb class="flex-1 rounded-full bg-line" />
			</ScrollArea.Scrollbar>
		</ScrollArea.Root>

		<footer
			class="page-footer h-16 shrink-0 border-t border-chrome-line bg-chrome text-chrome-muted"
		>
			<div id="page-footer-actions" class="page-footer-actions"></div>
		</footer>
	</div>

	<nav
		class="fixed inset-x-0 bottom-0 z-10 border-t border-chrome-line bg-chrome pb-[env(safe-area-inset-bottom)] lg:hidden"
		aria-label="Glavna navigacija"
	>
		<NavigationLinks
			items={currentUser ? [...primaryNavigation, ...secondaryNavigation] : secondaryNavigation}
			pathname={page.url.pathname}
			role={currentUser?.role}
			iconOnlyOnSmall
			class="flex h-16 items-center justify-center gap-4 overflow-x-auto px-4 text-sm max-sm:justify-around max-sm:gap-2"
		/>
	</nav>
</div>

{#if preparationReport}
	<PreparationReport report={preparationReport} />
{/if}

{#snippet sonnerInfoIcon()}
	<Info class="size-4" aria-hidden="true" />
{/snippet}

{#snippet sonnerSuccessIcon()}
	<CircleCheck class="size-4" aria-hidden="true" />
{/snippet}

{#snippet sonnerWarningIcon()}
	<TriangleAlert class="size-4" aria-hidden="true" />
{/snippet}

{#snippet sonnerErrorIcon()}
	<CircleX class="size-4" aria-hidden="true" />
{/snippet}

{#snippet sonnerCloseIcon()}
	<X class="size-4" aria-hidden="true" />
{/snippet}

<Toaster
	position="top-center"
	theme={theme.current}
	closeButton
	infoIcon={sonnerInfoIcon}
	successIcon={sonnerSuccessIcon}
	warningIcon={sonnerWarningIcon}
	errorIcon={sonnerErrorIcon}
	closeIcon={sonnerCloseIcon}
	duration={5000}
	gap={8}
	mobileOffset="12px"
	toastOptions={{
		unstyled: true,
		class:
			'flex w-[calc(100vw-1.5rem)] max-w-md items-center gap-3 rounded-md border border-l-[3px] border-line border-l-brand bg-surface px-3.5 py-3 text-sm text-ink shadow-lg shadow-black/15',
		classes: {
			content: 'min-w-0 flex-1',
			title: 'whitespace-pre-line font-medium leading-5',
			icon: 'shrink-0 text-brand',
			default: 'border-brand/45 border-l-brand bg-brand-soft',
			info: 'border-brand/45 border-l-brand bg-brand-soft',
			error: 'border-danger/45 border-l-danger !bg-danger-soft [&_[data-icon]]:text-danger',
			success: 'border-success/45 border-l-success !bg-success-soft [&_[data-icon]]:text-success',
			warning: 'border-warning/45 border-l-warning !bg-warning-soft [&_[data-icon]]:text-warning',
			closeButton:
				'order-last ml-auto inline-grid size-6 shrink-0 cursor-pointer place-items-center self-center rounded text-muted transition-colors hover:bg-soft hover:text-ink focus-visible:outline focus-visible:outline-1 focus-visible:outline-brand'
		}
	}}
/>
