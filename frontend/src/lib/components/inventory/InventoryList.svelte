<script lang="ts">
	import ArrowDown from '@lucide/svelte/icons/arrow-down';
	import ArrowUp from '@lucide/svelte/icons/arrow-up';
	import ArrowUpDown from '@lucide/svelte/icons/arrow-up-down';
	import Eye from '@lucide/svelte/icons/eye';
	import { resolve } from '$app/paths';
	import type {
		ConsumptionStatus,
		Item,
		ItemProperty,
		ItemTypeOption,
		PropertyOption,
		SortOrder
	} from '$lib/api';
	import { displayJson } from '$lib/domain/items';
	import ItemConsumptionControl from '$lib/components/inventory/ItemConsumptionControl.svelte';

	let {
		items,
		itemTypes,
		properties,
		canManage,
		sortPropertyID,
		sortOrder,
		onconsumptionchange,
		onrequest,
		onsortopen
	}: {
		items: Item[];
		itemTypes: ItemTypeOption[];
		properties: PropertyOption[];
		canManage: boolean;
		// The property the list is sorted by, if any; undefined means the default newest-first order.
		sortPropertyID: number | undefined;
		sortOrder: SortOrder;
		onconsumptionchange: (item: Item, status: ConsumptionStatus) => void;
		onrequest: (itemID: number, reason: string) => Promise<void>;
		onsortopen: () => void;
	} = $props();

	let typeNames = $derived(new Map(itemTypes.map((itemType) => [itemType.id, itemType.name])));
	let propertyNames = $derived(new Map(properties.map((property) => [property.id, property.name])));
	let sortedPropertyName = $derived(
		sortPropertyID === undefined ? '' : (propertyNames.get(sortPropertyID) ?? '')
	);

	function typeName(item: Item) {
		return typeNames.get(item.type_id) ?? 'Nepoznat tip';
	}

	function bgPropertyClass(property: ItemProperty) {
		if (property.value_type === 'expiry') {
			if (property.smart_data === 'expired') return 'bg-danger-soft text-danger';
			else if (property.smart_data === 'expiring_soon') return 'bg-warning-soft text-warning';
		}

		return 'bg-soft text-muted';
	}
</script>

<div class="-mx-4 mt-4 border-y border-line bg-surface sm:-mx-6">
	<table class="hidden min-w-full table-fixed text-left text-sm lg:table">
		<colgroup>
			<col class="w-64" />
			<col />
			<col class="w-48" />
			<col class="w-24" />
		</colgroup>
		<thead class="border-b border-line bg-soft text-muted">
			<tr class="h-12">
				<th class="px-4 py-3 font-medium">Stavka</th>
				<!-- The properties of an item are one cell rather than a column each, so the header can't
				     sort on its own: it opens the dialog that asks which property to sort by. The arrow
				     only turns directional once a property sort is actually what's applied. -->
				<th class="p-0 font-medium">
					<button
						class="flex h-12 w-full items-center gap-1.5 px-4 text-left transition-colors hover:text-ink"
						type="button"
						title={sortedPropertyName ? `Sortirano po: ${sortedPropertyName}` : 'Sortiraj stavke'}
						onclick={onsortopen}
					>
						Svojstva
						{#if sortPropertyID === undefined}
							<ArrowUpDown class="size-3.5" aria-hidden="true" />
						{:else if sortOrder === 'asc'}
							<ArrowUp class="size-3.5 text-brand" aria-hidden="true" />
						{:else}
							<ArrowDown class="size-3.5 text-brand" aria-hidden="true" />
						{/if}
					</button>
				</th>
				<th class="px-4 py-3 font-medium">Stanje</th>
				<th class="px-4 py-3 text-right font-medium"><span class="sr-only">Detalji</span></th>
			</tr>
		</thead>
		<tbody class="divide-y divide-line text-ink">
			{#each items as item (item.id)}
				<tr class="h-16">
					<td class="px-4 py-3 align-middle">
						<div class="block min-w-0">
							<p class="truncate text-xs text-muted">{typeName(item)}</p>
							<p class="mt-0.5 truncate font-medium">{item.derived_name}</p>
						</div>
					</td>
					<td class="px-4 py-3 align-middle">
						<div class="flex flex-wrap gap-1.5">
							{#each item.properties as property (property.id)}
								{#if property.visibility === 'overview'}
									<span class="rounded px-2 py-1 text-xs {bgPropertyClass(property)}">
										{propertyNames.get(property.id) ?? `Svojstvo #${property.id}`}: {displayJson(
											property.value_type,
											property.value
										)}
									</span>
								{/if}
							{/each}
							{#if item.properties.length === 0}<span class="text-muted">—</span>{/if}
						</div>
					</td>
					<td class="px-4 py-3 align-middle">
						<ItemConsumptionControl
							{item}
							{canManage}
							class="h-8 w-44 px-2.5 text-xs font-medium"
							{onconsumptionchange}
							{onrequest}
						/>
					</td>
					<td class="px-4 py-3 text-right align-middle">
						<div class="flex justify-end gap-1">
							<a
								class="inline-grid size-8 place-items-center rounded text-muted hover:bg-soft hover:text-ink"
								href={resolve(`/items/${item.id}`)}
								aria-label={`Detalji stavke ${item.id}`}
								title="Detalji"
							>
								<Eye class="size-4" aria-hidden="true" />
							</a>
						</div>
					</td>
				</tr>
			{/each}
			{#if items.length === 0}
				<tr class="h-16">
					<td class="px-4 py-3 align-middle text-muted" colspan={4}> Nema stavki. </td>
				</tr>
			{/if}
		</tbody>
	</table>

	<ul class="divide-y divide-line lg:hidden" aria-label="Stavke">
		{#each items as item (item.id)}
			{@const overviewProperties = item.properties.filter(
				(property) => property.visibility === 'overview'
			)}
			<li class="px-4 py-3">
				<div class="flex items-start justify-between gap-3">
					<a class="min-w-0 flex-1" href={resolve(`/items/${item.id}`)}>
						{#if item.derived_name}
							<p class="truncate text-xs text-muted">{typeName(item)}</p>
							<p class="mt-0.5 truncate text-sm font-medium text-ink hover:text-brand">
								{item.derived_name}
							</p>
						{:else}
							<p class="truncate text-sm font-medium text-ink hover:text-brand">{typeName(item)}</p>
						{/if}
						{#if overviewProperties.length}
							<p class="mt-2 text-xs leading-relaxed text-muted">
								{overviewProperties
									.map(
										(property) =>
											`${propertyNames.get(property.id) ?? `Svojstvo #${property.id}`}: ${displayJson(property.value_type, property.value)}`
									)
									.join(' · ')}
							</p>
						{/if}
					</a>
					<div class="flex shrink-0 gap-1">
						<a
							class="inline-grid size-8 place-items-center rounded text-muted hover:bg-soft hover:text-ink"
							href={resolve(`/items/${item.id}`)}
							aria-label={`Detalji stavke ${item.id}`}
							title="Detalji"
						>
							<Eye class="size-4" aria-hidden="true" />
						</a>
					</div>
				</div>
				<div class="mt-3">
					<ItemConsumptionControl
						{item}
						{canManage}
						class="h-8 w-full px-2.5 text-xs font-medium"
						{onconsumptionchange}
						{onrequest}
					/>
				</div>
			</li>
		{/each}
		{#if items.length === 0}<li class="px-4 py-3 text-sm text-muted">Nema stavki.</li>{/if}
	</ul>
</div>
