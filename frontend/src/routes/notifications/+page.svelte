<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { resolve } from '$app/paths';
	import Bell from '@lucide/svelte/icons/bell';
	import CalendarClock from '@lucide/svelte/icons/calendar-clock';
	import CalendarX from '@lucide/svelte/icons/calendar-x';
	import ClipboardCheck from '@lucide/svelte/icons/clipboard-check';
	import X from '@lucide/svelte/icons/x';
	import { api, ApiError, type Notification } from '$lib/api';
	import PaginationFooter from '$lib/components/shared/PaginationFooter.svelte';
	import ProtectedPageState from '$lib/components/shared/ProtectedPageState.svelte';
	import { formatRequestDate } from '$lib/domain/item-requests';
	import {
		notificationDetail,
		notificationHref,
		notificationIcon,
		notificationTitle
	} from '$lib/domain/notifications';
	import { createAuthPage } from '$lib/state/auth-page.svelte';
	import { notifications } from '$lib/state/notifications.svelte';
	import { createServerPagination } from '$lib/state/server-pagination.svelte';
	import { getTablePage, updateTableQuery } from '$lib/state/table-query';
	import { Button } from 'bits-ui';
	import { toast } from 'svelte-sonner';

	const authPage = createAuthPage({ unavailableMessage: 'Obaveštenja trenutno nisu dostupna.' });

	// ServerPagination only tracks items/offset/total, so the unread total is captured on the way
	// through the loader.
	let totalUnread = $state(0);

	const notificationsPage = createServerPagination<Notification>({
		loadPage: async ({ limit, offset }) => {
			const loaded = await api.listNotifications({ limit, offset });
			totalUnread = loaded.total_unread;
			return loaded;
		},
		unavailableMessage: 'Obaveštenja nisu učitana.'
	});

	let markReadTimer: ReturnType<typeof setTimeout> | null = null;
	let markedRead = false;

	onMount(() => {
		void authPage.load();

		return () => {
			if (markReadTimer !== null) clearTimeout(markReadTimer);
		};
	});

	$effect(() => {
		if (!authPage.state.authorized) return;
		notificationsPage.sync({ page: getTablePage(page.url) });
	});

	// Give the user a moment to actually look at what is new before clearing the badge, scaled a
	// little by how much there is to take in. The timer is dropped on unmount, so bouncing straight
	// off the page leaves everything unread.
	$effect(() => {
		if (markedRead || !notificationsPage.loaded || totalUnread === 0) return;

		const unreadOnPage = notificationsPage.items.filter((item) => !item.read).length;
		markedRead = true;
		markReadTimer = setTimeout(
			() => {
				markReadTimer = null;
				void api
					.readNotifications()
					.then(() => notifications.markAllRead())
					.catch(() => {
						markedRead = false;
					});
			},
			Math.min(5000, 1200 + 150 * unreadOnPage)
		);
	});

	async function dismiss(notification: Notification) {
		try {
			await api.deleteNotification(notification.id);
			notificationsPage.reloadAfterDelete();
			if (!notification.read) void notifications.refresh({ force: true });
		} catch (reason) {
			toast.error(reason instanceof ApiError ? reason.message : 'Obaveštenje nije obrisano.');
		}
	}

	function goToPage(nextPage: number) {
		updateTableQuery({ page: nextPage });
	}
</script>

<svelte:head>
	<title>Obaveštenja | Popisomator</title>
</svelte:head>

{#snippet body(notification: Notification, detail: string)}
	<p class={`text-sm text-ink ${notification.read ? '' : 'font-medium'}`}>
		{notificationTitle(notification)}
	</p>
	{#if detail}
		<p class="mt-0.5 truncate text-sm text-muted" title={detail}>{detail}</p>
	{/if}
	<p class="mt-1 text-xs text-muted">{formatRequestDate(notification.created_at)}</p>
{/snippet}

<main class="px-4 pt-4 pb-8 sm:px-6">
	<ProtectedPageState
		loading={authPage.state.loading}
		error={authPage.state.error || notificationsPage.error}
		authorized={authPage.state.authorized}
	>
		<p class="font-mono text-xs font-medium tracking-wide text-muted">
			UKUPNO: {notificationsPage.total}{totalUnread > 0 ? ` · NEPROČITANIH: ${totalUnread}` : ''}
		</p>
		<div class="-mx-4 mt-4 border-y border-line bg-surface sm:-mx-6">
			<ul class="divide-y divide-line" aria-label="Obaveštenja">
				{#each notificationsPage.items as notification (notification.id)}
					{@const href = notificationHref(notification, authPage.state.user?.role)}
					{@const detail = notificationDetail(notification)}
					{@const icon = notificationIcon(notification)}
					<li
						class={`flex items-start gap-3 px-4 py-3 ${
							notification.read ? '' : 'border-l-2 border-l-brand bg-brand-soft'
						}`}
					>
						<span
							class={`mt-0.5 shrink-0 ${notification.read ? 'text-muted' : 'text-brand'}`}
							aria-hidden="true"
						>
							{#if icon === 'request'}
								<ClipboardCheck class="size-4" />
							{:else if icon === 'expiring'}
								<CalendarClock class="size-4" />
							{:else if icon === 'expired'}
								<CalendarX class="size-4" />
							{:else}
								<Bell class="size-4" />
							{/if}
						</span>
						{#if href}
							<a class="min-w-0 flex-1" href={resolve(href)}>
								{@render body(notification, detail)}
							</a>
						{:else}
							<div class="min-w-0 flex-1">{@render body(notification, detail)}</div>
						{/if}
						<Button.Root
							class="inline-grid size-8 shrink-0 place-items-center rounded-md text-muted transition-colors hover:bg-soft hover:text-ink"
							aria-label="Obriši obaveštenje"
							title="Obriši obaveštenje"
							onclick={() => void dismiss(notification)}
						>
							<X class="size-4" aria-hidden="true" />
						</Button.Root>
					</li>
				{/each}
				{#if notificationsPage.items.length === 0}
					<li class="px-4 py-3 text-sm text-muted">Nema obaveštenja.</li>
				{/if}
			</ul>
		</div>
		<PaginationFooter
			total={notificationsPage.total}
			perPage={notificationsPage.perPage}
			page={notificationsPage.currentPage}
			hasPreviousPage={notificationsPage.hasPreviousPage}
			hasNextPage={notificationsPage.hasNextPage}
			loading={notificationsPage.loading}
			onpagechange={goToPage}
		/>
	</ProtectedPageState>
</main>
