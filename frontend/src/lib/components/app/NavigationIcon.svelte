<script lang="ts">
	import BookOpen from '@lucide/svelte/icons/book-open';
	import Boxes from '@lucide/svelte/icons/boxes';
	import ClipboardCheck from '@lucide/svelte/icons/clipboard-check';
	import Settings from '@lucide/svelte/icons/settings';
	import UsersRound from '@lucide/svelte/icons/users-round';
	import type { NavigationIconName } from '$lib/domain/navigation';

	let { name }: { name: NavigationIconName } = $props();
</script>

{#if name === 'inventory'}
	<Boxes class="size-4 shrink-0" aria-hidden="true" />
{:else if name === 'catalog'}
	<BookOpen class="size-4 shrink-0" aria-hidden="true" />
{:else if name === 'users'}
	<UsersRound class="size-4 shrink-0" aria-hidden="true" />
{:else if name === 'requests'}
	<ClipboardCheck class="size-4 shrink-0" aria-hidden="true" />
{:else}
	<Settings class="size-4 shrink-0" aria-hidden="true" />
{/if}
