<script lang="ts">
	import ErrorState from '$lib/components/ErrorState.svelte';
</script>

<svelte:head>
	<title>Popisomator</title>
	<meta name="description" content="Popisomator" />
</svelte:head>

<main class="grid min-h-[calc(100svh-14rem)] place-items-center px-4 sm:px-6">
	<ErrorState code="67" message="Gde je projekat? Ja ovde još ništa ne vidim?" />
</main>
