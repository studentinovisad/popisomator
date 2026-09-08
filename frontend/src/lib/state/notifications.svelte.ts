import { api } from '$lib/api';

const pollIntervalMs = 60_000;
// Focus and visibility changes often fire together; this keeps them from stacking up requests.
const minimumFetchGapMs = 10_000;

class NotificationsBadge {
	unreadCount = $state(0);

	#interval: ReturnType<typeof setInterval> | null = null;
	#lastFetchedAt = 0;
	#pending: Promise<void> | null = null;

	// A stale badge is not worth a toast, so failures leave the last known count in place.
	refresh = async ({ force = false } = {}) => {
		if (!force && Date.now() - this.#lastFetchedAt < minimumFetchGapMs) return;
		if (this.#pending) return this.#pending;

		this.#lastFetchedAt = Date.now();
		this.#pending = api
			.countUnreadNotifications()
			.then((count) => {
				this.unreadCount = count;
			})
			.catch(() => {})
			.finally(() => {
				this.#pending = null;
			});

		return this.#pending;
	};

	// Starts polling and returns a cleanup function. Safe to call again; the previous poller is
	// stopped first so a re-mount never leaves two intervals running.
	start = () => {
		this.stop();

		void this.refresh({ force: true });
		this.#interval = setInterval(() => void this.refresh({ force: true }), pollIntervalMs);

		const refreshWhenVisible = () => {
			if (document.visibilityState === 'visible') void this.refresh();
		};

		document.addEventListener('visibilitychange', refreshWhenVisible);
		window.addEventListener('focus', refreshWhenVisible);

		return () => {
			this.stop();
			document.removeEventListener('visibilitychange', refreshWhenVisible);
			window.removeEventListener('focus', refreshWhenVisible);
		};
	};

	stop = () => {
		if (this.#interval === null) return;

		clearInterval(this.#interval);
		this.#interval = null;
	};

	// The notifications page marks everything read in one request, so the badge can drop to zero
	// without asking the server again.
	markAllRead = () => {
		this.unreadCount = 0;
		this.#lastFetchedAt = Date.now();
	};

	clear = () => {
		this.stop();
		this.unreadCount = 0;
		this.#lastFetchedAt = 0;
	};
}

export const notifications = new NotificationsBadge();
