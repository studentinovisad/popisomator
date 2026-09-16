<script lang="ts">
	import type { LocationOption } from '$lib/api';
	import LocationBlock from '$lib/components/catalog/LocationBlock.svelte';
	import { Ellipsis, Eye, Pencil, Trash2 } from '@lucide/svelte';
	import { DropdownMenu } from 'bits-ui';

	let {
		location,
		alternateBg,
		deletelocation
	}: {
		location: LocationOption;
		alternateBg: boolean;
		deletelocation: (location: LocationOption) => void;
	} = $props();

	type DropdownOption = {
		value: string;
		label: string;
	};

	function onDelete(option: DropdownOption) {
		switch (option.value) {
			case 'edit':

			default:
				return;
		}
	}
</script>

<div class="p-4 {alternateBg ? 'bg-surface' : 'bg-chrome'} rounded-md border border-chrome-line">
	<div class="flex items-center justify-between gap-2">
		<p>{location.name}</p>
		<DropdownMenu.Root>
			<DropdownMenu.Trigger
				class="inline-grid size-8 place-items-center rounded text-muted hover:bg-soft hover:text-ink"
			>
				<Ellipsis />
			</DropdownMenu.Trigger>
			<DropdownMenu.Portal>
				<DropdownMenu.Content
					class="z-40 max-h-64 w-(--bits-combobox-anchor-width) overflow-y-auto rounded-md border border-line bg-surface p-1 shadow-lg shadow-black/15"
				>
					<a href={`/?location_id=${location.id}`}>
						<DropdownMenu.Item
							class="flex cursor-pointer items-center rounded px-3 py-2 text-sm text-ink outline-none data-highlighted:bg-brand-soft"
						>
							<Eye class="mr-2 size-4" />
							Prikaži stavke unutra
						</DropdownMenu.Item>
					</a>
					<a href={`/locations/${location.id}`}>
						<DropdownMenu.Item
							class="flex cursor-pointer items-center rounded px-3 py-2 text-sm text-ink outline-none data-highlighted:bg-brand-soft"
						>
							<Pencil class="mr-2 size-4" />
							Izmeni
						</DropdownMenu.Item>
					</a>
					<DropdownMenu.Item
						onclick={() => deletelocation(location)}
						class="flex cursor-pointer items-center rounded px-3 py-2 text-sm text-danger outline-none data-highlighted:bg-danger-soft"
					>
						<Trash2 class="mr-2 size-4" />
						Obriši
					</DropdownMenu.Item>
				</DropdownMenu.Content>
			</DropdownMenu.Portal>
		</DropdownMenu.Root>
	</div>

	{#if location.children != undefined}
		<div class="mt-4 flex flex-wrap gap-4">
			{#each location.children as childLocation}
				<LocationBlock location={childLocation} alternateBg={!alternateBg} {deletelocation} />
			{/each}
		</div>
	{/if}
</div>
