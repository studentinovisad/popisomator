<script lang="ts">
	import ArrowUpDown from '@lucide/svelte/icons/arrow-up-down';
	import { Button, Select } from 'bits-ui';
	import type { ItemPropertyTotal, ItemTypeOption, PropertyOption } from '$lib/api';
	import TableSearch from '$lib/components/shared/TableSearch.svelte';
	import { displayJson } from '$lib/domain/items';

	let {
		total,
		propertyTotals = [],
		properties = [],
		itemTypes,
		typeFilter,
		search = $bindable(),
		loading,
		onitemtypechange,
		onsearch,
		onsortopen
	}: {
		total: number;
		propertyTotals?: ItemPropertyTotal[];
		properties?: PropertyOption[];
		itemTypes: ItemTypeOption[];
		typeFilter: string;
		search: string;
		loading: boolean;
		onitemtypechange: (itemTypeID: number | undefined) => void;
		onsearch: (search: string) => void;
		onsortopen: () => void;
	} = $props();

	// The totals cover every item matching the current filters, not just this page, so they belong
	// beside the item count rather than in the list below.
	let propertyNamesByID = $derived(
		new Map(properties.map((property) => [property.id, property.name]))
	);
	// Names are uppercased to sit consistently beside the UKUPNO label; the amounts are left alone,
	// since unit symbols are case sensitive and 'mL' must not become 'ML'.
	let summedProperties = $derived(
		propertyTotals.map((totalled) => ({
			name: (propertyNamesByID.get(totalled.property_id) ?? '').toLocaleUpperCase('sr-Latn-RS'),
			amount: displayJson(totalled.value_type, totalled.value)
		}))
	);

	let itemTypeOptions = $derived([
		{ value: 'all', label: 'Svi tipovi' },
		...itemTypes.map((itemType) => ({ value: String(itemType.id), label: itemType.name }))
	]);
</script>

<div class="flex items-center justify-between gap-4">
	<p
		class="min-w-0 font-mono text-xs leading-relaxed font-medium tracking-wide text-balance text-muted"
	>
		UKUPNO: {total}{#each summedProperties as summed (summed.name + summed.amount)}
			<span class="px-1.5 text-line">·</span>{summed.name}: {summed.amount}{/each}
	</p>
	<div class="flex shrink-0 items-center gap-2">
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
		<!-- Kept outside the type filter's {#if}: below lg the item table's header row is replaced by
		     a card list, so this button is the only way into the sort dialog there. -->
		<Button.Root
			class="flex size-9 items-center justify-center rounded-md border border-chrome-line bg-transparent text-on-chrome transition-colors hover:border-brand"
			onclick={onsortopen}
			aria-label="Sortiraj stavke"
			title="Sortiraj stavke"
		>
			<ArrowUpDown class="size-4" aria-hidden="true" />
		</Button.Root>
	</div>
</div>

<TableSearch
	id="item-derived-name-search"
	placeholder="Pretraži po nazivu stavke"
	bind:search
	{loading}
	{onsearch}
/>
