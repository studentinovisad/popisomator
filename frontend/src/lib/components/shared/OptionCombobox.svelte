<script lang="ts">
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
		onvaluechange
	}: {
		id?: string;
		options: Option[];
		value?: string;
		placeholder?: string;
		emptyMessage?: string;
		disabled?: boolean;
		onvaluechange?: (value: string) => void;
	} = $props();

	let query = $state('');
	let filteredOptions = $derived(
		options.filter((option) => option.name.toLocaleLowerCase().includes(query.toLocaleLowerCase()))
	);
	let items = $derived(options.map((option) => ({ value: String(option.id), label: option.name })));

	function handleInput(event: Event) {
		query = (event.currentTarget as HTMLInputElement).value;
	}

	function handleValueChange(nextValue: string) {
		query = options.find((option) => String(option.id) === nextValue)?.name ?? '';
		onvaluechange?.(nextValue);
	}
</script>

<Combobox.Root
	type="single"
	bind:value
	{items}
	{disabled}
	allowDeselect={false}
	onValueChange={handleValueChange}
>
	<div class="relative">
		<Combobox.Input
			{id}
			class="block h-10 w-full rounded-md border border-line bg-surface py-0 pr-10 pl-3 text-sm text-ink placeholder:text-muted hover:border-brand focus-visible:border-brand"
			{placeholder}
			oninput={handleInput}
		/>
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
