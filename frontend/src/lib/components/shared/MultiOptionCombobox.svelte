<script lang="ts">
	import X from '@lucide/svelte/icons/x';
	import { Combobox } from 'bits-ui';

	type Option = {
		id: number;
		name: string;
	};

	let {
		options,
		values = $bindable([]),
		placeholder = 'Pretražite opcije',
		emptyMessage = 'Nema odgovarajućih opcija.',
		disabled = false,
		onvaluechange
	}: {
		options: Option[];
		values?: string[];
		placeholder?: string;
		emptyMessage?: string;
		disabled?: boolean;
		onvaluechange?: (values: string[]) => void;
	} = $props();

	let query = $state('');
	let filteredOptions = $derived(
		options.filter((option) => option.name.toLocaleLowerCase().includes(query.toLocaleLowerCase()))
	);
	let selectedOptions = $derived(options.filter((option) => values.includes(String(option.id))));
	let items = $derived(options.map((option) => ({ value: String(option.id), label: option.name })));

	function handleInput(event: Event) {
		query = (event.currentTarget as HTMLInputElement).value;
	}

	function handleValueChange(nextValues: string[]) {
		onvaluechange?.(nextValues);
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
	{disabled}
	onValueChange={handleValueChange}
>
	<div
		class="flex min-h-10 flex-wrap items-center gap-1.5 rounded-md border border-line bg-surface px-2 py-1.5 focus-within:border-brand hover:border-brand"
	>
		{#each selectedOptions as option (option.id)}
			<span class="inline-flex items-center gap-1 rounded bg-brand-soft px-2 py-1 text-xs text-ink">
				{option.name}
				<button
					class="rounded text-muted hover:text-ink"
					type="button"
					onclick={() => removeOption(option.id)}
					aria-label={`Ukloni ${option.name}`}
				>
					<X class="size-3" aria-hidden="true" />
				</button>
			</span>
		{/each}
		<Combobox.Input
			class="min-w-32 grow bg-transparent px-1 py-0.5 text-sm text-ink outline-none placeholder:text-muted"
			{placeholder}
			oninput={handleInput}
		/>
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
