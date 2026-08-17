<script lang="ts">
	import PageLoader from '$lib/components/shared/PageLoader.svelte';
	import type { Snippet } from 'svelte';

	let {
		loading,
		error = '',
		authorized,
		children
	}: {
		loading: boolean;
		error?: string;
		authorized: boolean;
		children: Snippet;
	} = $props();
</script>

{#if loading}
	<PageLoader />
{:else if error}
	<div class="grid min-h-[calc(100svh-14rem)] place-items-center">
		<p class="text-danger" role="alert">{error}</p>
	</div>
{:else if authorized}
	{@render children()}
{/if}
