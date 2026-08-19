<script lang="ts">
	import { Select } from 'bits-ui';
	import type { ItemTypeOption } from '$lib/api';
	import TableSearch from '$lib/components/shared/TableSearch.svelte';

	let {
		total,
		itemTypes,
		typeFilter,
		search = $bindable(),
		loading,
		onitemtypechange,
		onsearch
	}: {
		total: number;
		itemTypes: ItemTypeOption[];
		typeFilter: string;
		search: string;
		loading: boolean;
		onitemtypechange: (itemTypeID: number | undefined) => void;
		onsearch: (search: string) => void;
	} = $props();

	let itemTypeOptions = $derived([
		{ value: 'all', label: 'Svi tipovi' },
		...itemTypes.map((itemType) => ({ value: String(itemType.id), label: itemType.name }))
	]);
</script>

<div class="flex items-center justify-between gap-4">
	<p class="font-mono text-xs leading-none font-medium tracking-wide text-muted">UKUPNO: {total}</p>
	{#if itemTypes.length > 1}
		<Select.Root
			type="single"
			value={typeFilter}
			items={itemTypeOptions}
			onValueChange={(value) => onitemtypechange(value === 'all' ? undefined : Number(value))}
		>
			<Select.Trigger
				class="flex h-9 w-40 items-center justify-between rounded-md border border-chrome-line bg-transparent px-3 text-sm text-on-chrome transition-colors hover:border-brand"
				aria-label="Filtriraj stavke po tipu"
			>
				<Select.Value />
			</Select.Trigger>
			<Select.Portal>
				<Select.Content
					class="z-10 w-44 rounded-md border border-line bg-surface p-1 shadow-lg shadow-black/15"
					sideOffset={4}
				>
					<Select.Viewport>
						{#each itemTypeOptions as option (option.value)}
							<Select.Item
								value={option.value}
								label={option.label}
								class="cursor-pointer rounded px-3 py-2 outline-none data-highlighted:bg-brand-soft"
							>
								{option.label}
							</Select.Item>
						{/each}
					</Select.Viewport>
				</Select.Content>
			</Select.Portal>
		</Select.Root>
	{/if}
</div>

<TableSearch
	id="item-derived-name-search"
	placeholder="Pretraži po nazivu stavke"
	bind:search
	{loading}
	{onsearch}
/>
