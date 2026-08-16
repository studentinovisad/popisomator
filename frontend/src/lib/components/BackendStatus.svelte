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

<p class="inline-flex" role="status">
	<span
		class="inline-flex size-10 items-center justify-center rounded-md border border-line bg-surface shadow-sm"
		aria-hidden="true"
	>
		<span
			class={`size-2.5 rounded-full ${
				status === 'checking'
					? 'animate-pulse bg-warning'
					: status === 'available'
						? 'bg-success'
						: 'bg-danger'
			}`}
		></span>
	</span>
	<span class="sr-only">
		{#if status === 'checking'}
			Backend se proverava.
		{:else if status === 'available'}
			Backend je dostupan.
		{:else}
			Backend nije dostupan.
		{/if}
		></span
	>
</p>
