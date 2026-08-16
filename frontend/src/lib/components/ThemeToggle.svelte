<script lang="ts">
	import Moon from '@lucide/svelte/icons/moon';
	import Sun from '@lucide/svelte/icons/sun';
	import { Button } from 'bits-ui';
	import { theme } from '$lib/theme.svelte';
</script>

<Button.Root
	aria-label={theme.current === 'dark' ? 'Uključi svetlu temu' : 'Uključi tamnu temu'}
	class="inline-flex size-10 items-center justify-center rounded-md border border-line bg-surface text-muted shadow-sm transition-colors hover:border-brand/40 hover:bg-brand-soft hover:text-brand"
	type="button"
	onclick={() => theme.set(theme.current === 'dark' ? 'light' : 'dark')}
>
	{#if theme.current === 'dark'}
		<Sun class="size-4" aria-hidden="true" />
	{:else}
		<Moon class="size-4" aria-hidden="true" />
	{/if}
</Button.Root>
