<script lang="ts">
	import { Button } from 'bits-ui';
	import type { UserRole } from '$lib/api';

	let {
		onclick
	}: {
		onclick?: (approved: boolean) => void;
	} = $props();
</script>

<div class="flex justify-end gap-2">
	<Button.Root
		class="h-10 w-10 rounded-md bg-brand text-sm font-medium text-on-brand transition-colors hover:bg-brand-strong disabled:cursor-not-allowed disabled:opacity-50"
		onclick={() => onclick?.(true)}
	>
		Y
	</Button.Root>
	<Button.Root
		class="h-10 w-10 rounded-md bg-brand text-sm font-medium text-on-brand transition-colors hover:bg-brand-strong disabled:cursor-not-allowed disabled:opacity-50"
		onclick={() => onclick?.(false)}
	>
		N
	</Button.Root>
</div>
