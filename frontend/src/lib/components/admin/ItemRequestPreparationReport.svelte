<script lang="ts">
	import type { ItemRequestPreparationReport } from '$lib/api';
	import { consumptionLabel, displayJson } from '$lib/domain/items';
	import { SvelteMap } from 'svelte/reactivity';

	let { report }: { report: ItemRequestPreparationReport } = $props();
	type PreparationItem = ItemRequestPreparationReport['items'][number];

	type StorageGroup = {
		cabinet: string;
		box: string;
		items: PreparationItem[];
	};

	type LocationGroup = {
		name: string;
		groups: StorageGroup[];
	};

	// TODO: Replace these temporary property-name conventions once locations and containers are entities.
	const storagePropertyNames = new Set(['Lokacija', 'Ormar', 'Mesto/kutija']);
	let groupedItems = $derived(groupItems(report.items));

	function formatDate(value: string) {
		return new Intl.DateTimeFormat('sr-Latn-RS', {
			dateStyle: 'long',
			timeStyle: 'short'
		}).format(new Date(value));
	}

	function derivedNamePropertyNames(item: PreparationItem) {
		return new Set(
			[...item.derived_name_format.matchAll(/\{([^{}]+)\}/g)].map((match) => match[1].trim())
		);
	}

	function propertySections(item: PreparationItem) {
		const names = derivedNamePropertyNames(item);
		const derivedName = [] as PreparationItem['properties'];
		const overview = [] as PreparationItem['properties'];
		const details = [] as PreparationItem['properties'];

		for (const property of orderedProperties(item)) {
			if (storagePropertyNames.has(property.name)) continue;
			if (names.has(property.name)) derivedName.push(property);
			else if (property.visibility === 'overview') overview.push(property);
			else details.push(property);
		}

		return { derivedName, overview, details };
	}

	function orderedProperties(item: PreparationItem) {
		return [...item.properties].sort((left, right) => left.position - right.position);
	}

	function storageValue(item: PreparationItem, propertyName: string) {
		const property = item.properties.find((entry) => entry.name === propertyName);
		return property ? displayJson(property.value_type, property.value) : '';
	}

	function groupItems(items: PreparationItem[]) {
		const locations = new SvelteMap<string, LocationGroup>();

		for (const item of items) {
			const location = storageValue(item, 'Lokacija') || 'Neraspoređeno';
			const cabinet = storageValue(item, 'Ormar');
			const box = storageValue(item, 'Mesto/kutija');
			const locationGroup = locations.get(location) ?? { name: location, groups: [] };
			const key = [cabinet, box].join('\u0000');
			const group = locationGroup.groups.find(
				(candidate) => [candidate.cabinet, candidate.box].join('\u0000') === key
			) ?? { cabinet, box, items: [] };

			group.items.push(item);
			if (!locationGroup.groups.includes(group)) locationGroup.groups.push(group);
			locations.set(location, locationGroup);
		}

		return [...locations.values()]
			.sort((left, right) => left.name.localeCompare(right.name, 'sr'))
			.map((location) => ({
				...location,
				groups: location.groups
					.map((group) => ({
						...group,
						items: group.items.sort((left, right) => left.name.localeCompare(right.name, 'sr'))
					}))
					.sort((left, right) =>
						[left.cabinet, left.box]
							.join('\u0000')
							.localeCompare([right.cabinet, right.box].join('\u0000'), 'sr')
					)
			}));
	}
</script>

<article
	class="preparation-report mx-auto max-w-3xl px-6 py-10 text-ink sm:px-10"
	aria-hidden="true"
>
	<header class="report-header border-b-2 border-ink pb-6">
		<h1 class="mt-1 text-2xl font-semibold tracking-tight">Priprema stavki</h1>
		<p class="mt-3 text-sm text-muted">
			Korisnik: <span class="font-medium text-ink">{report.user.name}</span>
			<span aria-hidden="true"> · </span> Broj potrebnih stavki:
			<span class="font-medium text-ink">{report.items.length}</span>
		</p>
	</header>

	{#if groupedItems.length > 0}
		<div class="mt-6 space-y-7">
			{#each groupedItems as location (location.name)}
				<section>
					<h2 class="border-b border-ink pb-1 text-base font-semibold">{location.name}</h2>
					{#each location.groups as group (`${group.cabinet}:${group.box}`)}
						{#if group.cabinet || group.box}
							<p class="mt-4 pl-3 text-sm text-muted">
								{#if group.cabinet}Ormar: {group.cabinet}{/if}
								{#if group.cabinet && group.box}
									·
								{/if}
								{#if group.box}Kutija: {group.box}{/if}
							</p>
						{/if}
						<ol class="mt-2 divide-y divide-line">
							{#each group.items as item (item.id)}
								{@const sections = propertySections(item)}
								<li class="break-inside-avoid py-4 first:pt-3">
									<div class="grid grid-cols-[1.5rem_minmax(0,1fr)_8rem] gap-x-3">
										<span class="row-start-1 block size-6 self-center rounded-sm border border-ink"
										></span>
										<div class="row-start-1 min-w-0">
											<p class="text-xs text-muted">{item.type_name}</p>
											<p class="mt-0.5 font-medium">{item.name}</p>
										</div>
										<time
											class="row-start-1 self-center text-right text-xs leading-relaxed text-muted"
											>{formatDate(item.requested_at)}</time
										>
										<div class="col-span-2 col-start-2 min-w-0">
											{#if item.reason}
												<p class="mt-3 border-l-2 border-line pl-3 text-sm text-muted">
													{item.reason}
												</p>
											{/if}
											{#if sections.derivedName.length > 0}
												<ul class="mt-3 flex flex-wrap gap-1.5" aria-label="Naziv stavke">
													{#each sections.derivedName as property (property.name)}
														<li
															class="max-w-full rounded-full border border-line bg-soft px-2 py-1 text-xs text-ink"
														>
															<span class="font-medium">{property.name}:</span>
															{displayJson(property.value_type, property.value)}
														</li>
													{/each}
												</ul>
											{/if}
											<ul class="mt-2 flex flex-wrap gap-1.5" aria-label="Pregled stavke">
												<li
													class="rounded-full border border-line bg-soft px-2 py-1 text-xs text-ink"
												>
													<span class="font-medium">Stanje:</span>
													{consumptionLabel(item.consumption)}
												</li>
												{#each sections.overview as property (property.name)}
													<li
														class="max-w-full rounded-full border border-line bg-soft px-2 py-1 text-xs text-ink"
													>
														<span class="font-medium">{property.name}:</span>
														{displayJson(property.value_type, property.value)}
													</li>
												{/each}
											</ul>
											{#if sections.details.length > 0}
												<ul class="mt-2 flex flex-wrap gap-1.5" aria-label="Detalji stavke">
													{#each sections.details as property (property.name)}
														<li
															class="max-w-full rounded-full border border-line px-2 py-1 text-xs text-ink"
														>
															<span class="font-medium">{property.name}:</span>
															{displayJson(property.value_type, property.value)}
														</li>
													{/each}
												</ul>
											{/if}
										</div>
									</div>
								</li>
							{/each}
						</ol>
					{/each}
				</section>
			{/each}
		</div>
	{:else}
		<p class="mt-6 text-sm text-muted">Ovaj korisnik trenutno nema zahteva za pripremu.</p>
	{/if}
</article>

<style>
	.preparation-report {
		display: none;
	}

	@media print {
		:global(body) {
			background: white;
		}

		:global(.app-shell) {
			display: none;
		}

		.preparation-report {
			display: block;
			max-width: none;
			padding: 0;
			color-scheme: light;
			--color-canvas: #f7f7f6;
			--color-surface: #ffffff;
			--color-ink: #272726;
			--color-muted: #6e6e6b;
			--color-faint: #737370;
			--color-line: #e2e2df;
			--color-soft: #f0f0ee;
		}
	}
</style>
