<script lang="ts">
	let status = $state<'checking' | 'available' | 'unavailable'>('checking');

	$effect(() => {
		let active = true;
		let controller: AbortController | undefined;

		const checkBackend = () => {
			controller?.abort();
			controller = new AbortController();
			const currentController = controller;
			const timeout = window.setTimeout(() => currentController.abort(), 5_000);

			void fetch('/api/health', { signal: currentController.signal })
				.then((response) => {
					if (active && controller === currentController) {
						status = response.ok ? 'available' : 'unavailable';
					}
				})
				.catch(() => {
					if (active && controller === currentController) {
						status = 'unavailable';
					}
				})
				.finally(() => window.clearTimeout(timeout));
		};

		checkBackend();
		const interval = window.setInterval(checkBackend, 5_000);

		return () => {
			active = false;
			window.clearInterval(interval);
			controller?.abort();
		};
	});
</script>

<p class="font-mono text-xs tracking-wide sm:text-right" aria-live="polite">
	{#if status === 'checking'}
		Saša vizuelizuje backend…
	{:else if status === 'available'}
		Saša ne vidi problem sa backend-om.
	{:else}
		Saša ne može da pronađe backend.
	{/if}
</p>
