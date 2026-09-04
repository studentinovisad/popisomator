<script lang="ts">
	import ChevronDown from '@lucide/svelte/icons/chevron-down';
	import X from '@lucide/svelte/icons/x';
	import { onDestroy } from 'svelte';
	import { Combobox } from 'bits-ui';
	import type { PropertyValue, PropertyValueType } from '$lib/api';
	import { displayJson } from '$lib/domain/items';

	let {
		label,
		options,
		value,
		value_type,
		onvaluechange,
		onsearchchange
	}: {
		label: string;
		options: PropertyValue[];
		value: PropertyValue;
		value_type: PropertyValueType;
		onvaluechange: (value: PropertyValue | undefined) => void;
		onsearchchange: (search: string) => void;
	} = $props();

	let open = $state(false);
	let inputText = $derived(value != undefined ? displayJson(value_type, value) : '');
	let searchTimeout: ReturnType<typeof setTimeout> | undefined;
	let control: HTMLDivElement | null = null;
	let filteredOptions = $derived(
		options.filter((option) =>
			displayJson(value_type, option).toLocaleLowerCase().includes(inputText.toLocaleLowerCase())
		)
	);
	let items = $derived(
		options.map((option) => ({
			value: JSON.stringify(option),
			label: displayJson(value_type, option)
		}))
	);

	onDestroy(() => {
		if (searchTimeout) clearTimeout(searchTimeout);
	});

	function handleInput(event: Event) {
		inputText = (event.currentTarget as HTMLInputElement).value;
		queueSearch();
	}

	function handleOpenChange(nextOpen: boolean) {
		open = nextOpen;
		if (nextOpen) onsearchchange(inputText);
	}

	function queueSearch() {
		if (searchTimeout) clearTimeout(searchTimeout);
		searchTimeout = setTimeout(() => onsearchchange(inputText), 200);
	}

	function handleValueChange(nextValue: string) {
		inputText = nextValue;
		onvaluechange(JSON.parse(nextValue));
	}

	function restoreActiveValue(event: FocusEvent) {
		if (event.relatedTarget instanceof Node && control?.contains(event.relatedTarget)) return;

		const input = event.currentTarget as HTMLInputElement;
		if (input.value === '') {
			clear();
			return;
		}

		inputText = value != undefined ? displayJson(value_type, value) : '';
	}

	function clear() {
		inputText = '';
		onvaluechange(undefined);
	}
</script>

<Combobox.Root
	type="single"
	{items}
	value={JSON.stringify(value)}
	inputValue={inputText}
	{open}
	allowDeselect={false}
	onOpenChange={handleOpenChange}
	onValueChange={handleValueChange}
>
	<div bind:this={control} class="relative max-w-64 min-w-44">
		<Combobox.Input
			class="property-filter-input h-10 w-full py-0 pr-13 pl-3 text-sm"
			placeholder={label}
			aria-label={`Filtriraj po svojstvu: ${label}`}
			oninput={handleInput}
			onblur={restoreActiveValue}
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
						value={JSON.stringify(option)}
						label={displayJson(value_type, option)}
						class="cursor-pointer rounded px-3 py-2 text-sm text-ink outline-none data-highlighted:bg-brand-soft"
					>
						{displayJson(value_type, option)}
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
