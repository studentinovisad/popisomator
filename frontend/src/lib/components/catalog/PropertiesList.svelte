<script lang="ts">
	import ChevronDown from '@lucide/svelte/icons/chevron-down';
	import Pencil from '@lucide/svelte/icons/pencil';
	import Trash2 from '@lucide/svelte/icons/trash-2';
	import { resolve } from '$app/paths';
	import { Collapsible } from 'bits-ui';
	import type { Property } from '$lib/api';
	import { propertyValueTypeLabel } from '$lib/domain/items';

	let {
		properties,
		deleteproperty
	}: {
		properties: Property[];
		deleteproperty: (property: Property) => void;
	} = $props();

	function hasDefaultValue(property: Property) {
		return property.default_value !== null && property.default_value !== undefined;
	}
</script>

<div class="-mx-4 mt-4 overflow-hidden border-y border-line bg-surface sm:-mx-6">
	<table class="hidden w-full table-fixed text-left text-sm md:table">
		<colgroup>
			<col class="w-[30%]" />
			<col class="w-36" />
			<col class="w-56" />
			<col class="w-28" />
		</colgroup>
		<thead class="border-b border-line bg-soft text-muted">
			<tr class="h-12">
				<th class="px-4 py-3 font-medium">Naziv</th>
				<th class="px-4 py-3 font-medium">Tip</th>
				<th class="px-4 py-3 font-medium">Podrazumevana vrednost</th>
				<th class="px-4 py-3 text-right font-medium">Radnje</th>
			</tr>
		</thead>
		<tbody class="text-ink">
			{#each properties as property (property.id)}
				<Collapsible.Root>
					{#snippet child({ props: rootProps })}
						<tr {...rootProps} class="h-16 transition-colors hover:bg-soft/35">
							<td class="px-4 py-3 align-middle">
								<span class="block truncate font-medium" title={property.name}>{property.name}</span
								>
							</td>
							<td class="px-4 py-3 align-middle">
								<span class="rounded bg-soft px-2 py-0.5 text-sm font-medium text-muted"
									>{propertyValueTypeLabel(property.value_type)}</span
								>
							</td>
							<td class="px-4 py-3 align-middle">
								<span
									class={`rounded px-2 py-0.5 text-sm font-medium ${hasDefaultValue(property) ? 'bg-success-soft text-success' : 'bg-danger-soft text-danger'}`}
								>
									{hasDefaultValue(property) ? 'da' : 'ne'}
								</span>
							</td>
							<td class="px-4 py-3 text-right align-middle">
								<div class="flex justify-end gap-1">
									{#if property.description}
										<Collapsible.Trigger
											class="description-toggle inline-grid size-8 place-items-center rounded text-muted hover:bg-soft hover:text-ink"
											aria-label="Prikaži ili sakrij opis svojstva"
											title="Prikaži opis"
										>
											<ChevronDown class="size-4" aria-hidden="true" />
										</Collapsible.Trigger>
									{/if}
									<a
										class="inline-grid size-8 place-items-center rounded text-muted hover:bg-soft hover:text-ink"
										href={resolve(`/catalog/properties/${property.id}`)}
										aria-label={`Izmeni svojstvo ${property.name}`}
										title="Izmeni"
									>
										<Pencil class="size-4" aria-hidden="true" />
									</a>
									<button
										class="inline-grid size-8 place-items-center rounded text-danger hover:bg-danger-soft"
										type="button"
										onclick={() => deleteproperty(property)}
										aria-label={`Obriši svojstvo ${property.name}`}
										title="Obriši"
									>
										<Trash2 class="size-4" aria-hidden="true" />
									</button>
								</div>
							</td>
						</tr>
						{#if property.description}
							<Collapsible.Content>
								{#snippet child({ props: contentProps })}
									<tr {...contentProps} class="property-description-content">
										<td colspan="4" class="px-4 pb-3">
											<p class="border-t border-line pt-3 text-sm leading-6 text-muted">
												{property.description}
											</p>
										</td>
									</tr>
								{/snippet}
							</Collapsible.Content>
						{/if}
					{/snippet}
				</Collapsible.Root>
			{/each}
			{#if properties.length === 0}
				<tr class="h-16">
					<td class="px-4 py-3 align-middle text-muted" colspan="4">Nema svojstava.</td>
				</tr>
			{/if}
		</tbody>
	</table>

	<div
		class="grid grid-cols-[minmax(0,1fr)_auto] gap-x-3 border-b border-line bg-soft px-4 py-2 text-xs font-medium text-muted md:hidden"
	>
		<span>Naziv · tip · ima podrazumevanu vrednost</span>
		<span>Radnje</span>
	</div>
	<ul class="divide-y divide-line md:hidden" aria-label="Svojstva">
		{#each properties as property (property.id)}
			<li class="px-4 py-3">
				<Collapsible.Root>
					<div class="grid grid-cols-[minmax(0,1fr)_auto] gap-x-3 gap-y-2">
						<p class="min-w-0 truncate text-sm font-medium text-ink" title={property.name}>
							{property.name}
						</p>
						<div class="row-span-2 flex shrink-0 items-center gap-1">
							{#if property.description}
								<Collapsible.Trigger
									class="description-toggle inline-grid size-8 place-items-center rounded text-muted hover:bg-soft hover:text-ink"
									aria-label="Prikaži ili sakrij opis svojstva"
									title="Prikaži opis"
								>
									<ChevronDown class="size-4" aria-hidden="true" />
								</Collapsible.Trigger>
							{/if}
							<a
								class="inline-grid size-8 place-items-center rounded text-muted hover:bg-soft hover:text-ink"
								href={resolve(`/catalog/properties/${property.id}`)}
								aria-label={`Izmeni svojstvo ${property.name}`}
								title="Izmeni"
							>
								<Pencil class="size-4" aria-hidden="true" />
							</a>
							<button
								class="inline-grid size-8 place-items-center rounded text-danger hover:bg-danger-soft"
								type="button"
								onclick={() => deleteproperty(property)}
								aria-label={`Obriši svojstvo ${property.name}`}
								title="Obriši"
							>
								<Trash2 class="size-4" aria-hidden="true" />
							</button>
						</div>
						<div class="order-3 flex min-w-0 flex-wrap items-center gap-1.5 text-sm">
							<span class="text-muted">Tip:</span>
							<span class="rounded bg-soft px-2 py-0.5 font-medium text-muted"
								>{propertyValueTypeLabel(property.value_type)}</span
							>
							<span
								class={`rounded px-2 py-0.5 font-medium ${hasDefaultValue(property) ? 'bg-success-soft text-success' : 'bg-danger-soft text-danger'}`}
							>
								{hasDefaultValue(property) ? 'da' : 'ne'}
							</span>
						</div>
					</div>
					{#if property.description}
						<Collapsible.Content class="property-description-content">
							<p class="mt-3 border-t border-line pt-3 text-sm leading-6 text-muted">
								{property.description}
							</p>
						</Collapsible.Content>
					{/if}
				</Collapsible.Root>
			</li>
		{/each}
		{#if properties.length === 0}<li class="px-4 py-3 text-sm text-muted">Nema svojstava.</li>{/if}
	</ul>
</div>

<style>
	:global(.description-toggle svg) {
		transition: transform 160ms ease;
	}

	:global(.description-toggle[data-state='open'] svg) {
		transform: rotate(180deg);
	}

	:global(.property-description-content[data-state='open']) {
		overflow: hidden;
		animation: property-description-open 180ms ease-out;
	}

	:global(.property-description-content[data-state='closed']) {
		overflow: hidden;
		animation: property-description-close 140ms ease-in;
	}

	@keyframes property-description-open {
		from {
			height: 0;
			opacity: 0;
			transform: translateY(-0.25rem);
		}

		to {
			height: var(--bits-collapsible-content-height);
		}
	}

	@keyframes property-description-close {
		from {
			height: var(--bits-collapsible-content-height);
		}

		to {
			height: 0;
			opacity: 0;
			transform: translateY(-0.25rem);
		}
	}

	@media (prefers-reduced-motion: reduce) {
		:global(.description-toggle svg),
		:global(.property-description-content[data-state]) {
			transition: none;
			animation: none;
		}
	}
</style>
