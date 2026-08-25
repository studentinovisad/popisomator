<script lang="ts">
	import { ScrollArea } from 'bits-ui';
	import PropertyFilterCombobox from '$lib/components/inventory/PropertyFilterCombobox.svelte';

	export type PropertyFilter = {
		id: number;
		name: string;
		value: string;
	};

	let {
		filters,
		filterOptions,
		onfilterchange,
		onfiltervaluesearch
	}: {
		filters: PropertyFilter[];
		filterOptions: Record<number, string[]>;
		onfilterchange: (propertyID: number, value: string) => void;
		onfiltervaluesearch: (propertyID: number, search: string) => void;
	} = $props();
</script>

{#if filters.length > 0}
	<ScrollArea.Root class="mt-3 overflow-hidden" type="auto">
		<ScrollArea.Viewport class="w-full">
			<div class="flex w-max gap-2 pb-2">
				{#each filters as filter (filter.id)}
					<PropertyFilterCombobox
						label={filter.name}
						options={filterOptions[filter.id] ?? []}
						value={filter.value}
						onvaluechange={(value) => onfilterchange(filter.id, value)}
						onsearchchange={(search) => onfiltervaluesearch(filter.id, search)}
					/>
				{/each}
			</div>
		</ScrollArea.Viewport>
		<ScrollArea.Scrollbar
			class="mt-0.5 flex h-1 touch-none rounded-full bg-transparent"
			orientation="horizontal"
		>
			<ScrollArea.Thumb class="rounded-full bg-line/60" />
		</ScrollArea.Scrollbar>
	</ScrollArea.Root>
{/if}
