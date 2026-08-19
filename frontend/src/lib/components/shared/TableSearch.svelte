<script lang="ts">
	import Search from '@lucide/svelte/icons/search';
	import { Button } from 'bits-ui';

	let {
		id,
		placeholder,
		search = $bindable(),
		loading = false,
		onsearch
	}: {
		id: string;
		placeholder: string;
		search: string;
		loading?: boolean;
		onsearch: (search: string) => void;
	} = $props();

	function submit(event: SubmitEvent) {
		event.preventDefault();
		onsearch(search.trim());
	}
</script>

<form class="flex flex-col gap-2 sm:flex-row sm:items-end" onsubmit={submit}>
	<div class="min-w-0 flex-1">
		<div class="relative mt-3">
			<Search
				class="pointer-events-none absolute top-1/2 left-3 size-4 -translate-y-1/2 text-faint"
				aria-hidden="true"
			/>
			<input {id} class="h-10 w-full pl-9" bind:value={search} {placeholder} />
		</div>
	</div>
	<Button.Root
		class="h-10 rounded-md bg-brand px-4 text-sm font-medium text-on-brand transition-colors hover:bg-brand-strong disabled:cursor-not-allowed disabled:opacity-50"
		disabled={loading}
		type="submit"
	>
		Pretraži
	</Button.Root>
</form>
