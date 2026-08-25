<script lang="ts">
	import ChevronDown from '@lucide/svelte/icons/chevron-down';
	import X from '@lucide/svelte/icons/x';
	import { Combobox } from 'bits-ui';

	let {
		label,
		options,
		value,
		onvaluechange,
		onsearchchange
	}: {
		label: string;
		options: string[];
		value: string;
		onvaluechange: (value: string) => void;
		onsearchchange: (search: string) => void;
	} = $props();

	let query = $state('');
	let open = $state(false);
	let inputValue = $state('');
	let searchTimeout: ReturnType<typeof setTimeout> | undefined;
	let filteredOptions = $derived(
		options.filter((option) => option.toLocaleLowerCase().includes(query.toLocaleLowerCase()))
	);
	let items = $derived(options.map((option) => ({ value: option, label: option })));

	$effect(() => {
		inputValue = value;
	});

	function handleInput(event: Event) {
		query = (event.currentTarget as HTMLInputElement).value;
		queueSearch();
	}

	function handleOpenChange(nextOpen: boolean) {
		open = nextOpen;
		if (nextOpen) onsearchchange(query);
	}

	function queueSearch() {
		if (searchTimeout) clearTimeout(searchTimeout);
		searchTimeout = setTimeout(() => onsearchchange(query), 200);
	}

	function handleValueChange(nextValue: string) {
		query = nextValue;
		inputValue = nextValue;
		onvaluechange(nextValue);
	}

	function clear() {
		query = '';
		inputValue = '';
		onvaluechange('');
	}
</script>

<Combobox.Root
	type="single"
	{items}
	{value}
	{inputValue}
	{open}
	allowDeselect={false}
	onOpenChange={handleOpenChange}
	onValueChange={handleValueChange}
>
	<div class="relative max-w-64 min-w-44">
		<Combobox.Input
			class="property-filter-input h-10 w-full py-0 pr-13 pl-3 text-sm"
			placeholder={label}
			aria-label={`Filtriraj po svojstvu: ${label}`}
			oninput={handleInput}
		/>
		{#if value}
			<button
				type="button"
				class="absolute top-1/2 right-8 grid size-5 -translate-y-1/2 cursor-pointer place-items-center rounded text-muted hover:text-ink focus-visible:text-brand focus-visible:outline-none"
				onclick={clear}
				aria-label={`Poništi filter: ${label}`}
			>
				<X class="size-3.5" aria-hidden="true" />
			</button>
		{/if}
		<Combobox.Trigger
			class="group absolute top-1/2 right-2 grid size-5 -translate-y-1/2 cursor-pointer place-items-center rounded text-muted hover:text-ink focus-visible:text-brand focus-visible:outline-none"
			aria-label={`Filtriraj po svojstvu: ${label}`}
		>
			<ChevronDown
				class="size-4 transition-transform duration-150 group-data-[state=open]:rotate-180"
				aria-hidden="true"
			/>
		</Combobox.Trigger>
	</div>
	<Combobox.Portal>
		<Combobox.Content
			class="z-40 max-h-64 w-(--bits-combobox-anchor-width) overflow-y-auto rounded-md border border-line bg-surface p-1 shadow-lg shadow-black/15"
			sideOffset={4}
		>
			<Combobox.Viewport>
				{#each filteredOptions as option (option)}
					<Combobox.Item
						value={option}
						label={option}
						class="cursor-pointer rounded px-3 py-2 text-sm text-ink outline-none data-highlighted:bg-brand-soft"
					>
						{option}
					</Combobox.Item>
				{/each}
				{#if filteredOptions.length === 0}
					<p class="px-3 py-2 text-sm text-muted">Nema odgovarajućih vrednosti.</p>
				{/if}
			</Combobox.Viewport>
		</Combobox.Content>
	</Combobox.Portal>
</Combobox.Root>

<style>
	:global(.property-filter-input:focus-visible) {
		outline: none !important;
		box-shadow: none !important;
	}
</style>
