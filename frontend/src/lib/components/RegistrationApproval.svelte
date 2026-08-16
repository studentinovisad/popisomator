<script lang="ts">
	import Check from '@lucide/svelte/icons/check';
	import X from '@lucide/svelte/icons/x';
	import { Button } from 'bits-ui';

	let {
		onclick
	}: {
		onclick?: (approved: boolean) => void;
	} = $props();
</script>

<div class="flex justify-end gap-2">
	<Button.Root
		class="inline-flex size-9 items-center justify-center rounded-md bg-brand text-on-brand transition-colors hover:bg-brand-strong disabled:cursor-not-allowed disabled:opacity-50"
		onclick={() => onclick?.(true)}
		aria-label="Prihvati zahtev"
		title="Prihvati zahtev"
	>
		<Check class="size-4" aria-hidden="true" />
	</Button.Root>
	<Button.Root
		class="inline-flex size-9 items-center justify-center rounded-md border border-danger bg-surface text-danger transition-colors hover:bg-danger-soft disabled:cursor-not-allowed disabled:opacity-50"
		onclick={() => onclick?.(false)}
		aria-label="Odbij zahtev"
		title="Odbij zahtev"
	>
		<X class="size-4" aria-hidden="true" />
	</Button.Root>
</div>
