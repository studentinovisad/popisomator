<script lang="ts">
	import { page } from '$app/state';
	import ErrorState from '$lib/components/shared/ErrorState.svelte';

	let message = $derived(page.status === 404 ? 'Ova stranica ne postoji…' : 'Nastao je problem…');
</script>

<svelte:head>
	<title>{page.status} | Popisomator</title>
</svelte:head>

<main class="grid min-h-[calc(100svh-14rem)] place-items-center px-4 sm:px-6">
	<ErrorState code={page.status} {message} />
</main>
