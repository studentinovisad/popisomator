<script lang="ts">
	import BookOpen from '@lucide/svelte/icons/book-open';
	import ClipboardCheck from '@lucide/svelte/icons/clipboard-check';
	import FlaskConical from '@lucide/svelte/icons/flask-conical';
	import Pencil from '@lucide/svelte/icons/pencil';
	import Plus from '@lucide/svelte/icons/plus';
	import Trash2 from '@lucide/svelte/icons/trash-2';
	import type { AuditIconName } from '$lib/domain/audit-log';

	let { name, size = 'size-4' }: { name: AuditIconName; size?: string } = $props();
</script>

{#if name === 'create'}
	<Plus class={`${size} shrink-0`} aria-hidden="true" />
{:else if name === 'delete'}
	<Trash2 class={`${size} shrink-0`} aria-hidden="true" />
{:else if name === 'consume'}
	<FlaskConical class={`${size} shrink-0`} aria-hidden="true" />
{:else if name === 'request'}
	<ClipboardCheck class={`${size} shrink-0`} aria-hidden="true" />
{:else if name === 'catalog'}
	<BookOpen class={`${size} shrink-0`} aria-hidden="true" />
{:else}
	<Pencil class={`${size} shrink-0`} aria-hidden="true" />
{/if}
