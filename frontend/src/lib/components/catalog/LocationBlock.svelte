<script lang="ts">
	import type { LocationOption } from '$lib/api';
    import LocationBlock from '$lib/components/catalog/LocationBlock.svelte';
	import { Ellipsis } from '@lucide/svelte';
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
        value: string,
        label: string
    }

    function onDelete(option: DropdownOption) {
        switch (option.value) {
            case "edit":

            default:
                return
        }
    }
</script>

<div class="p-4 {alternateBg ? 'bg-surface' : 'bg-chrome'} border border-chrome-line rounded-md">
    <div class="flex justify-between items-center gap-2">
        <p>{location.name}</p>
        <DropdownMenu.Root>
            <DropdownMenu.Trigger class="inline-grid size-8 place-items-center rounded text-muted hover:bg-soft hover:text-ink">
                <Ellipsis/>
            </DropdownMenu.Trigger>
            <DropdownMenu.Portal>
                <DropdownMenu.Content class="z-40 max-h-64 w-(--bits-combobox-anchor-width) overflow-y-auto rounded-md border border-line bg-surface p-1 shadow-lg shadow-black/15">
                    <a href={`/locations/${location.id}`}>
                        <DropdownMenu.Item class="cursor-pointer rounded px-3 py-2 text-sm text-ink outline-none data-highlighted:bg-brand-soft">
                            Izmeni
                        </DropdownMenu.Item>
                    </a>
                    <DropdownMenu.Item onclick={() => deletelocation(location)} class="cursor-pointer rounded px-3 py-2 text-sm text-ink outline-none data-highlighted:bg-brand-soft">
                        Obriši
                    </DropdownMenu.Item>
                </DropdownMenu.Content>
            </DropdownMenu.Portal>
        </DropdownMenu.Root>
    </div>
   
    {#if location.children != undefined}
        <div class="flex flex-wrap gap-4 mt-4">
            {#each location.children as childLocation}
                <LocationBlock location={childLocation} alternateBg={!alternateBg} deletelocation={deletelocation}/>
            {/each}
        </div>
        
    {/if}
</div>