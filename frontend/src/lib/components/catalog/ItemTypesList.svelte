<script lang="ts">
	import Pencil from '@lucide/svelte/icons/pencil';
	import Trash2 from '@lucide/svelte/icons/trash-2';
	import { resolve } from '$app/paths';
	import type { ItemType, PropertyOption } from '$lib/api';

	let {
		itemTypes,
		propertyOptions,
		deleteitemtype
	}: {
		itemTypes: ItemType[];
		propertyOptions: PropertyOption[];
		deleteitemtype: (itemType: ItemType) => void;
	} = $props();

	function assignedProperties(itemType: ItemType) {
		return itemType.properties.length
			? itemType.properties
					.map((itemProperty) => propertyOptions.find(property => property.id === itemProperty.id)?.name ?? 'Nepoznato svojstvo')
					.join(' · ')
			: 'Bez svojstava';
	}
</script>

<div class="-mx-4 mt-4 overflow-hidden border-y border-line bg-surface sm:-mx-6">
	<table class="hidden min-w-full table-fixed text-left text-sm md:table">
		<colgroup>
			<col class="w-[28%]" />
			<col />
			<col class="w-[35%]" />
			<col class="w-24" />
		</colgroup>
		<thead class="border-b border-line bg-soft text-muted">
			<tr class="h-12">
				<th class="px-4 py-3 font-medium">Naziv</th>
				<th class="px-4 py-3 font-medium">Opis</th>
				<th class="px-4 py-3 font-medium">Svojstva</th>
				<th class="px-4 py-3 text-right font-medium">Radnje</th>
			</tr>
		</thead>
		<tbody class="text-ink">
			{#each itemTypes as itemType (itemType.id)}
				<tr class="h-16 transition-colors hover:bg-soft/35">
					<td class="px-4 py-3 align-middle">
						<span class="block truncate font-medium" title={itemType.name}>{itemType.name}</span>
					</td>
					<td class="px-4 py-3 align-middle">
						<span class="block truncate text-muted" title={itemType.description || undefined}
							>{itemType.description || '—'}</span
						>
					</td>
					<td class="px-4 py-3 align-middle">
						<span class="block truncate text-muted" title={assignedProperties(itemType)}
							>{assignedProperties(itemType)}</span
						>
					</td>
					<td class="px-4 py-3 text-right align-middle">
						<div class="flex justify-end gap-1">
							<a
								class="inline-grid size-8 place-items-center rounded text-muted hover:bg-soft hover:text-ink"
								href={resolve(`/catalog/item-types/${itemType.id}`)}
								aria-label={`Izmeni tip stavke ${itemType.name}`}
								title="Izmeni"
							>
								<Pencil class="size-4" aria-hidden="true" />
							</a>
							<button
								class="inline-grid size-8 place-items-center rounded text-danger hover:bg-danger-soft"
								type="button"
								onclick={() => deleteitemtype(itemType)}
								aria-label={`Obriši tip stavke ${itemType.name}`}
								title="Obriši"
							>
								<Trash2 class="size-4" aria-hidden="true" />
							</button>
						</div>
					</td>
				</tr>
			{/each}
			{#if itemTypes.length === 0}
				<tr class="h-16">
					<td class="px-4 py-3 align-middle text-muted" colspan="4">Nema tipova stavki.</td>
				</tr>
			{/if}
		</tbody>
	</table>

	<div
		class="grid grid-cols-[minmax(0,1fr)_auto] gap-x-3 border-b border-line bg-soft px-4 py-2 text-xs font-medium text-muted md:hidden"
	>
		<span>Naziv · opis · svojstva</span>
		<span>Radnje</span>
	</div>
	<ul class="divide-y divide-line md:hidden" aria-label="Tipovi stavki">
		{#each itemTypes as itemType (itemType.id)}
			<li class="px-4 py-3">
				<div class="grid grid-cols-[minmax(0,1fr)_auto] gap-x-3 gap-y-2">
					<div class="min-w-0">
						<p class="truncate text-sm font-medium text-ink" title={itemType.name}>
							{itemType.name}
						</p>
						{#if itemType.description}
							<p class="mt-0.5 truncate text-sm text-muted" title={itemType.description}>
								Opis: {itemType.description}
							</p>
						{/if}
					</div>
					<div class="flex shrink-0 gap-1">
						<a
							class="inline-grid size-8 place-items-center rounded text-muted hover:bg-soft hover:text-ink"
							href={resolve(`/catalog/item-types/${itemType.id}`)}
							aria-label={`Izmeni tip stavke ${itemType.name}`}
							title="Izmeni"
						>
							<Pencil class="size-4" aria-hidden="true" />
						</a>
						<button
							class="inline-grid size-8 place-items-center rounded text-danger hover:bg-danger-soft"
							type="button"
							onclick={() => deleteitemtype(itemType)}
							aria-label={`Obriši tip stavke ${itemType.name}`}
							title="Obriši"
						>
							<Trash2 class="size-4" aria-hidden="true" />
						</button>
					</div>
					<p class="col-span-2 truncate text-xs text-muted" title={assignedProperties(itemType)}>
						Svojstva: {assignedProperties(itemType)}
					</p>
				</div>
			</li>
		{/each}
		{#if itemTypes.length === 0}<li class="px-4 py-3 text-sm text-muted">
				Nema tipova stavki.
			</li>{/if}
	</ul>
</div>
