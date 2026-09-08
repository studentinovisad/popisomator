<script lang="ts">
	import ChevronDown from '@lucide/svelte/icons/chevron-down';
	import X from '@lucide/svelte/icons/x';
	import { tick } from 'svelte';
	import { Combobox } from 'bits-ui';

	type Option = {
		id: number;
		name: string;
	};

	let {
		id,
		options,
		values = $bindable([]),
		placeholder = 'Pretražite opcije',
		emptyMessage = 'Nema odgovarajućih opcija.',
		disabled = false,
		showSelected = true,
		onvaluechange
	}: {
		id?: string;
		options: Option[];
		values?: string[];
		placeholder?: string;
		emptyMessage?: string;
		disabled?: boolean;
		showSelected?: boolean;
		onvaluechange?: (values: string[]) => void;
	} = $props();

	let query = $state('');
	let input = $state<HTMLInputElement | null>(null);
	let filteredOptions = $derived(
		options.filter(
			(option) =>
				!values.includes(String(option.id)) &&
				option.name.toLocaleLowerCase().includes(query.toLocaleLowerCase())
		)
	);
	let selectedOptions = $derived(
		values.flatMap((value) => {
			const option = options.find((candidate) => candidate.id === Number(value));
			return option ? [option] : [];
		})
	);
	let items = $derived(options.map((option) => ({ value: String(option.id), label: option.name })));

	function handleInput(event: Event) {
		query = (event.currentTarget as HTMLInputElement).value;
	}

	async function handleValueChange(nextValues: string[]) {
		onvaluechange?.(nextValues);
		await tick();
		query = '';
		if (input) {
			input.value = '';
			input.dispatchEvent(new Event('input', { bubbles: true }));
		}
	}

	function removeOption(id: number) {
		values = values.filter((value) => value !== String(id));
		onvaluechange?.(values);
	}
</script>

<Combobox.Root
	type="multiple"
	bind:value={values}
	{items}
	inputValue={query}
	{disabled}
	onValueChange={handleValueChange}
>
	{#if showSelected && selectedOptions.length}
		<ul class="mb-3 flex flex-wrap gap-2" aria-label="Izabrana svojstva">
			{#each selectedOptions as option (option.id)}
				<li
					class="inline-flex items-center gap-1.5 rounded-md border border-line bg-soft py-1.5 pr-1.5 pl-2.5 text-sm text-ink"
				>
					{option.name}
					<button
						class="grid size-5 cursor-pointer place-items-center rounded text-muted hover:bg-surface hover:text-ink"
						type="button"
						onclick={() => removeOption(option.id)}
						aria-label={`Ukloni ${option.name}`}
					>
						<X class="size-3.5" aria-hidden="true" />
					</button>
				</li>
			{/each}
		</ul>
	{/if}
	<div class="relative">
		<Combobox.Input
			{id}
			bind:ref={input}
			class="block h-10 w-full rounded-md border border-line bg-surface py-0 pr-10 pl-3 text-sm text-ink placeholder:text-muted hover:border-brand focus-visible:border-brand"
			{placeholder}
			oninput={handleInput}
		/>
		<Combobox.Trigger
			class="group absolute top-1/2 right-1 grid size-8 -translate-y-1/2 cursor-pointer place-items-center rounded text-muted outline-none hover:bg-soft hover:text-ink focus-visible:ring-1 focus-visible:ring-brand disabled:cursor-not-allowed"
			aria-label="Prikaži opcije"
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
				{#each filteredOptions as option (option.id)}
					<Combobox.Item
						value={String(option.id)}
						label={option.name}
						class="flex cursor-pointer items-center justify-between rounded px-3 py-2 text-sm text-ink outline-none data-highlighted:bg-brand-soft data-selected:text-brand"
					>
						{option.name}
					</Combobox.Item>
				{/each}
				{#if filteredOptions.length === 0}
					<p class="px-3 py-2 text-sm text-muted">{emptyMessage}</p>
				{/if}
			</Combobox.Viewport>
		</Combobox.Content>
	</Combobox.Portal>
</Combobox.Root>
