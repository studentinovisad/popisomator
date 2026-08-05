<script lang="ts">
	import { Toggle } from 'bits-ui';
	import { theme } from '$lib/theme.svelte';
</script>

<Toggle.Root
	pressed={theme.current === 'dark'}
	onPressedChange={(pressed) => theme.set(pressed ? 'dark' : 'light')}
	aria-label={theme.current === 'dark' ? 'Uključi svetlu temu' : 'Uključi tamnu temu'}
	class="inline-flex size-9 items-center justify-center rounded-md border border-on-chrome/30 text-on-chrome transition-colors hover:bg-on-chrome/10 data-[state=on]:bg-on-chrome/15"
>
	{#if theme.current === 'dark'}
		<svg
			aria-hidden="true"
			viewBox="0 0 24 24"
			fill="none"
			stroke="currentColor"
			stroke-width="2"
			class="size-4"
		>
			<circle cx="12" cy="12" r="4" />
			<path
				d="M12 2v2M12 20v2M4.93 4.93l1.41 1.41M17.66 17.66l1.41 1.41M2 12h2M20 12h2M4.93 19.07l1.41-1.41M17.66 6.34l1.41-1.41"
			/>
		</svg>
	{:else}
		<svg
			aria-hidden="true"
			viewBox="0 0 24 24"
			fill="none"
			stroke="currentColor"
			stroke-width="2"
			class="size-4"
		>
			<path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79Z" />
		</svg>
	{/if}
</Toggle.Root>
