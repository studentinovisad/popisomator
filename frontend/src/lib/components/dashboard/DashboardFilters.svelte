<script lang="ts">
	import { Select } from 'bits-ui';
	import ChevronDown from '@lucide/svelte/icons/chevron-down';
	import type { DashboardMonths, ItemTypeOption } from '$lib/api';
	import { dashboardMonthRanges } from '$lib/domain/dashboard';

	let {
		months,
		typeID,
		itemTypes,
		disabled = false,
		onmonthschange,
		ontypechange
	}: {
		months: DashboardMonths;
		typeID: string;
		itemTypes: ItemTypeOption[];
		disabled?: boolean;
		onmonthschange: (months: DashboardMonths) => void;
		ontypechange: (typeID: string) => void;
	} = $props();

	const typeOptions = $derived([
		{ value: 'all', label: 'Svi tipovi' },
		...itemTypes.map((itemType) => ({ value: String(itemType.id), label: itemType.name }))
	]);
	const selectedTypeLabel = $derived(
		typeOptions.find((option) => option.value === (typeID || 'all'))?.label ?? 'Svi tipovi'
	);
</script>

<div class="flex flex-wrap items-center gap-3">
	<div
		class="inline-flex rounded-md border border-line bg-surface p-0.5"
		role="group"
		aria-label="Period prikaza"
	>
		{#each dashboardMonthRanges as range (range.value)}
			<button
				type="button"
				{disabled}
				aria-pressed={months === range.value}
				class={`rounded px-3 py-1.5 text-sm transition-colors ${
					months === range.value
						? 'bg-brand text-on-brand'
						: 'text-muted hover:bg-soft hover:text-ink'
				}`}
				onclick={() => onmonthschange(range.value)}
			>
				{range.label}
			</button>
		{/each}
	</div>

	<Select.Root
		type="single"
		value={typeID || 'all'}
		items={typeOptions}
		{disabled}
		onValueChange={ontypechange}
	>
		<Select.Trigger
			aria-label="Tip stavke"
			class="flex h-9 w-52 items-center justify-between rounded-md border border-line bg-surface px-3 text-sm text-ink transition-colors hover:border-brand/40"
		>
			<span class="truncate">{selectedTypeLabel}</span>
			<ChevronDown class="size-4 shrink-0 text-muted" aria-hidden="true" />
		</Select.Trigger>
		<Select.Portal>
			<Select.Content
				class="z-30 max-h-72 w-52 overflow-y-auto rounded-md border border-line bg-surface p-1 shadow-lg shadow-black/15"
			>
				{#each typeOptions as option (option.value)}
					<Select.Item
						value={option.value}
						label={option.label}
						class="cursor-pointer rounded px-2 py-1.5 text-sm text-ink data-highlighted:bg-soft"
					>
						{option.label}
					</Select.Item>
				{/each}
			</Select.Content>
		</Select.Portal>
	</Select.Root>
</div>
