<script lang="ts">
	import Eye from '@lucide/svelte/icons/eye';
	import { Select } from 'bits-ui';
	import { resolve } from '$app/paths';
	import type { ConsumptionStatus, Item, ItemTypeOption, PropertyOption } from '$lib/api';
	import {
		consumptionClass,
		consumptionLabel,
		consumptionOptions,
		displayJson
	} from '$lib/domain/items';

	let {
		items,
		itemTypes,
		properties,
		onconsumptionchange
	}: {
		items: Item[];
		itemTypes: ItemTypeOption[];
		properties: PropertyOption[];
		onconsumptionchange: (item: Item, status: ConsumptionStatus) => void;
	} = $props();

	let typeNames = $derived(new Map(itemTypes.map((itemType) => [itemType.id, itemType.name])));
	let propertyNames = $derived(new Map(properties.map((property) => [property.id, property.name])));

	function typeName(item: Item) {
		return typeNames.get(item.type_id) ?? 'Nepoznat tip';
	}
</script>

<div class="-mx-4 mt-4 border-y border-line bg-surface sm:-mx-6">
	<table class="hidden min-w-full table-fixed text-left text-sm md:table">
		<colgroup>
			<col class="w-64" />
			<col />
			<col class="w-48" />
			<col class="w-24" />
		</colgroup>
		<thead class="border-b border-line bg-soft text-muted">
			<tr class="h-12">
				<th class="px-4 py-3 font-medium">Stavka</th>
				<th class="px-4 py-3 font-medium">Svojstva</th>
				<th class="px-4 py-3 font-medium">Stanje</th>
				<th class="px-4 py-3 text-right font-medium">Detalji</th>
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
									<span class="rounded bg-soft px-2 py-1 text-xs text-muted">
										{propertyNames.get(property.id) ?? `Svojstvo #${property.id}`}: {displayJson(
											property.value
										)}
									</span>
								{/if}
							{/each}
							{#if item.properties.length === 0}<span class="text-muted">—</span>{/if}
						</div>
					</td>
					<td class="px-4 py-3 align-middle">
						<Select.Root
							type="single"
							value={item.consumption}
							items={consumptionOptions}
							onValueChange={(value) => onconsumptionchange(item, value as ConsumptionStatus)}
						>
							<Select.Trigger
								class={`flex h-8 w-44 items-center justify-between rounded-md px-2.5 text-xs font-medium ${consumptionClass(item.consumption)}`}
								aria-label={`Stanje stavke ${item.id}`}
							>
								<Select.Value />
							</Select.Trigger>
							<Select.Portal>
								<Select.Content
									class="z-30 w-48 rounded-md border border-line bg-surface p-1 shadow-lg shadow-black/15"
									sideOffset={4}
								>
									<Select.Viewport>
										{#each consumptionOptions as option (option.value)}
											<Select.Item
												value={option.value}
												label={option.label}
												class="cursor-pointer rounded px-3 py-2 text-sm outline-none data-highlighted:bg-brand-soft"
											>
												{option.label}
											</Select.Item>
										{/each}
									</Select.Viewport>
								</Select.Content>
							</Select.Portal>
						</Select.Root>
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

	<ul class="divide-y divide-line md:hidden" aria-label="Stavke">
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
											`${propertyNames.get(property.id) ?? `Svojstvo #${property.id}`}: ${displayJson(property.value)}`
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
					<Select.Root
						type="single"
						value={item.consumption}
						items={consumptionOptions}
						onValueChange={(value) => onconsumptionchange(item, value as ConsumptionStatus)}
					>
						<Select.Trigger
							class={`flex h-8 w-full items-center justify-between rounded-md px-2.5 text-xs font-medium ${consumptionClass(item.consumption)}`}
						>
							<Select.Value>{consumptionLabel(item.consumption)}</Select.Value>
						</Select.Trigger>
						<Select.Portal>
							<Select.Content
								class="z-30 w-(--bits-select-anchor-width) rounded-md border border-line bg-surface p-1 shadow-lg shadow-black/15"
								sideOffset={4}
							>
								<Select.Viewport>
									{#each consumptionOptions as option (option.value)}
										<Select.Item
											value={option.value}
											label={option.label}
											class="cursor-pointer rounded px-3 py-2 text-sm outline-none data-highlighted:bg-brand-soft"
										>
											{option.label}
										</Select.Item>
									{/each}
								</Select.Viewport>
							</Select.Content>
						</Select.Portal>
					</Select.Root>
				</div>
			</li>
		{/each}
		{#if items.length === 0}<li class="px-4 py-3 text-sm text-muted">Nema stavki.</li>{/if}
	</ul>
</div>
