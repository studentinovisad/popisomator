<script lang="ts">
	import { api, ApiError, type Property } from '$lib/api';
	import { displayJson, propertyValueTypeLabel } from '$lib/items';
	import { Button, Label, ScrollArea } from 'bits-ui';

	let { properties, oncreated }: { properties: Property[]; oncreated: () => void } = $props();

	let name = $state('');
	let description = $state('');
	let selectedPropertyIDs = $state<number[]>([]);
	let creating = $state(false);
	let error = $state('');

	async function createItemType(event: SubmitEvent) {
		event.preventDefault();
		creating = true;
		error = '';

		try {
			await api.createItemType({
				name,
				description,
				properties: selectedPropertyIDs.map((id) => ({
					id,
					default_value: properties.find((property) => property.id === id)?.default_value ?? null
				}))
			});
			name = '';
			description = '';
			selectedPropertyIDs = [];
			oncreated();
		} catch (reason) {
			error = reason instanceof ApiError ? reason.message : 'Tip stavke nije sačuvan.';
		} finally {
			creating = false;
		}
	}
</script>

<form class="grid gap-4" onsubmit={createItemType}>
	<div>
		<Label.Root class="text-sm font-medium text-ink" for="item-type-name">Naziv</Label.Root>
		<input id="item-type-name" class="mt-1 block w-full" bind:value={name} required />
	</div>
	<div>
		<Label.Root class="text-sm font-medium text-ink" for="item-type-description">Opis</Label.Root>
		<textarea id="item-type-description" class="mt-1 block min-h-20 w-full" bind:value={description}
		></textarea>
	</div>
	<fieldset class="border-t border-line pt-4">
		<legend class="text-sm font-medium text-ink">Svojstva tipa</legend>
		<p class="mt-1 text-xs text-muted">Vrednosti se podrazumevano preuzimaju sa svojstva.</p>
		<ScrollArea.Root class="mt-3 h-64 overflow-hidden rounded-md border border-line" type="auto">
			<ScrollArea.Viewport class="h-full w-full">
				<div class="divide-y divide-line">
					{#each properties as property (property.id)}
						<label
							class="flex cursor-pointer items-center justify-between gap-4 px-3 py-2 hover:bg-soft"
						>
							<span>
								<span class="block text-sm font-medium text-ink">{property.name}</span>
								<span class="block text-xs text-muted">
									{propertyValueTypeLabel(property.value_type)} · {displayJson(
										property.default_value
									)}
								</span>
							</span>
							<input type="checkbox" bind:group={selectedPropertyIDs} value={property.id} />
						</label>
					{/each}
					{#if properties.length === 0}<p class="px-3 py-3 text-sm text-muted">
							Najpre dodajte svojstvo.
						</p>{/if}
				</div>
			</ScrollArea.Viewport>
			<ScrollArea.Scrollbar class="flex w-2.5 touch-none bg-soft p-0.5" orientation="vertical">
				<ScrollArea.Thumb class="flex-1 rounded-full bg-line" />
			</ScrollArea.Scrollbar>
		</ScrollArea.Root>
	</fieldset>
	<div class="border-t border-line pt-4">
		<Button.Root
			class="rounded-md bg-brand px-4 py-2 text-sm font-medium text-on-brand hover:bg-brand-strong disabled:opacity-60"
			disabled={creating}
			type="submit"
		>
			{creating ? 'Čuvanje…' : 'Dodaj tip'}
		</Button.Root>
		{#if error}<p class="mt-3 text-sm text-danger" role="alert">{error}</p>{/if}
	</div>
</form>
