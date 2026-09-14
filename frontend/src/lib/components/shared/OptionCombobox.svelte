<script lang="ts">
	import { X } from '@lucide/svelte';
	import ChevronDown from '@lucide/svelte/icons/chevron-down';
	import { Combobox } from 'bits-ui';

	type Option = {
		id: number;
		name: string;
	};

	let {
		id,
		options,
		value = $bindable(''),
		placeholder = 'Pretražite opcije',
		emptyMessage = 'Nema odgovarajućih opcija.',
		disabled = false,
		invalid = false,
		clearable = false,
		height = "10",
		bgColor = "surface",
		describedBy,
		onvaluechange
	}: {
		id?: string;
		options: Option[];
		value?: string;
		placeholder?: string;
		emptyMessage?: string;
		disabled?: boolean;
		invalid?: boolean;
		clearable?: boolean;
		height?: string;
		bgColor?: string;
		describedBy?: string;
		onvaluechange?: (value: string) => void;
	} = $props();

	let items = $derived(options.map((option) => ({ value: String(option.id), label: option.name })));
	let query = $derived(options.find((option) => String(option.id) === value)?.name || '');
	let filteredOptions = $derived(
		options.filter((option) => option.name.toLocaleLowerCase().includes(query.toLocaleLowerCase()))
	);
	

	function handleInput(event: Event) {
		query = (event.currentTarget as HTMLInputElement).value;
	}

	function handleValueChange(nextValue: string) {
		query = options.find((option) => String(option.id) === nextValue)?.name ?? '';
		onvaluechange?.(nextValue);
	}

	function clear() {
		query = '';
		value = '';
		onvaluechange?.(value);
	}
</script>

<Combobox.Root
	type="single"
	bind:value
	inputValue={query}
	{items}
	{disabled}
	allowDeselect={false}
	onValueChange={handleValueChange}
>
	<div class="relative">
		<Combobox.Input
			{id}
			class={`block h-${height} w-full rounded-md border border-line bg-${bgColor} py-0 pr-10 pl-3 text-sm text-ink placeholder:text-muted hover:border-brand focus-visible:border-brand ${invalid ? 'field-invalid' : ''}`}
			{placeholder}
			aria-invalid={invalid}
			aria-describedby={describedBy}
			oninput={handleInput}
		/>
		{#if value && clearable}
			<button
				type="button"
				class="absolute top-1/2 right-8 grid size-5 -translate-y-1/2 cursor-pointer place-items-center rounded text-muted hover:text-ink focus-visible:text-brand focus-visible:outline-none"
				onclick={clear}
				aria-label={`Poništi`}
			>
				<X class="size-3.5" aria-hidden="true" />
			</button>
		{/if}
		<Combobox.Trigger
			class="group absolute inset-y-0 right-0 grid w-10 cursor-pointer place-items-center rounded-r-md text-muted outline-none hover:text-ink focus-visible:ring-1 focus-visible:ring-brand disabled:cursor-not-allowed"
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
						class="cursor-pointer rounded px-3 py-2 text-sm text-ink outline-none data-highlighted:bg-brand-soft"
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
