<script lang="ts">
	let backendStatus = $state<'checking' | 'available' | 'unavailable'>('checking');

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
						backendStatus = response.ok ? 'available' : 'unavailable';
					}
				})
				.catch(() => {
					if (active && controller === currentController) {
						backendStatus = 'unavailable';
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

<svelte:head>
	<title>Popisomator</title>
	<meta name="description" content="Popisomator" />
</svelte:head>

<main class="mx-auto max-w-4xl px-6 py-16">
	<h1 class="text-3xl font-semibold text-slate-900">Popisomator</h1>

	<div class="mt-8 grid items-center gap-8 md:grid-cols-2">
		<blockquote
			class="text-3xl leading-tight font-semibold text-balance text-slate-800 sm:text-4xl"
		>
			Gde je PRIS projekat? Ja ovde još ništa ne vidim?
		</blockquote>
		<img class="w-full rounded-lg" src="/sasa.jpg" alt="Saša Matić" />
	</div>

	<footer
		class="mt-8 flex flex-col items-start gap-3 text-sm text-slate-500 sm:flex-row sm:items-center sm:justify-between"
	>
		<ul class="flex flex-wrap" aria-label="Original authors">
			{#each ['Đorđe Mančić', 'Matija Kljajić', 'Miša Stefanović'] as tim8_c (tim8_c)}
				<li class="after:mx-1.5 after:content-['·'] last:after:content-none">{tim8_c}</li>
			{/each}
		</ul>

		<p class="sm:text-right" aria-live="polite">
			Status:
			{#if backendStatus === 'checking'}
				Saša vizuelizuje backend…
			{:else if backendStatus === 'available'}
				Saša ne vidi problem sa backend-om.
			{:else}
				Saša ne može da pronađe backend.
			{/if}
		</p>
	</footer>
</main>
