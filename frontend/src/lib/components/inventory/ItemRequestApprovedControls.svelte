<script lang="ts">
	import { RotateCcw, Trash2 } from '@lucide/svelte';
	import { Button } from 'bits-ui';

	let {
		onclick
	}: {
		onclick?: (isDelete: boolean) => void;
	} = $props();
</script>

<div class="flex justify-end gap-2">
    <Button.Root
        class="inline-flex size-9 items-center justify-center rounded-md border border-brand bg-surface text-brand transition-colors hover:bg-brand-soft"
        onclick={() => onclick?.(false)}
        aria-label="Vrati na čekanje"
        title="Vrati na čekanje"
    >
        <RotateCcw class="size-4" aria-hidden="true" />
    </Button.Root>
    <Button.Root
        class="inline-flex size-9 items-center justify-center rounded-md border border-danger bg-surface text-danger transition-colors hover:bg-danger-soft"
        onclick={() => onclick?.(true)}
        aria-label="Ukloni odobrenje"
        title="Ukloni odobrenje"
    >
        <Trash2 class="size-4" aria-hidden="true" />
    </Button.Root>
</div>