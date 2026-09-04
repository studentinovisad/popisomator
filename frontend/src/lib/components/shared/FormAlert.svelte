<script lang="ts">
	import CircleAlert from '@lucide/svelte/icons/circle-alert';

	let { message }: { message: string } = $props();
</script>

<div
	class="flex items-start gap-2 rounded-md border border-danger/40 bg-danger-soft px-3 py-2 text-sm text-danger"
	role="alert"
>
	<CircleAlert class="mt-0.5 size-4 shrink-0" aria-hidden="true" />
	<p>{message}</p>
</div>
