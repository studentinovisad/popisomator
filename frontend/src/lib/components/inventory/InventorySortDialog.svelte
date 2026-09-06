<script lang="ts">
	import ArrowDown from '@lucide/svelte/icons/arrow-down';
	import ArrowUp from '@lucide/svelte/icons/arrow-up';
	import Check from '@lucide/svelte/icons/check';
	import { Button, Dialog, RadioGroup } from 'bits-ui';
	import type { SortOrder } from '$lib/api';

	// Items carry no sortable column of their own beyond their creation time, so the field list is
	// this one entry followed by the properties. It doubles as the RadioGroup value for that entry,
	// which is why it is a string rather than a property id.
	const createdAtField = 'created_at';

	let {
		open = $bindable(false),
		options,
		sortPropertyID,
		sortOrder,
		onsortchange
	}: {
		open?: boolean;
		options: { id: number; name: string }[];
		// The property currently sorted by; undefined means the items are ordered by creation time.
		sortPropertyID: number | undefined;
		sortOrder: SortOrder;
		onsortchange: (propertyID: number | undefined, order: SortOrder) => void;
	} = $props();

	// The choice is edited as a draft and only handed over on submit, so picking a field and a
	// direction reloads the list once instead of twice. Reopening starts from the applied sort.
	let draftField = $state(createdAtField);
	let draftOrder = $state<SortOrder>('desc');

	let fields = $derived([
		{ value: createdAtField, label: 'Datum unosa' },
		...options.map((option) => ({ value: String(option.id), label: option.name }))
	]);

	// Seeding on the closed -> open edge rather than in Dialog.Root's onOpenChange, because that only
	// fires for opens the dialog performs itself; here both triggers live outside it and set `open`
	// through the binding, which the callback never sees. wasOpen is deliberately not $state: it
	// exists to detect the edge, and making it reactive would re-run this on its own write.
	let wasOpen = false;
	$effect(() => {
		if (open === wasOpen) return;

		wasOpen = open;
		if (!open) return;
		draftField = sortPropertyID === undefined ? createdAtField : String(sortPropertyID);
		draftOrder = sortOrder;
	});

	function applySort(event: SubmitEvent) {
		event.preventDefault();
		onsortchange(draftField === createdAtField ? undefined : Number(draftField), draftOrder);
		open = false;
	}

	function directionClass(direction: SortOrder) {
		const selected = draftOrder === direction;
		return `inline-flex h-10 items-center justify-center gap-2 rounded-md border text-sm transition-colors ${
			selected
				? 'border-brand bg-brand-soft font-medium text-brand'
				: 'border-line bg-surface text-ink hover:border-brand/40 hover:bg-soft'
		}`;
	}
</script>

<Dialog.Root bind:open>
	<Dialog.Portal>
		<Dialog.Overlay class="fixed inset-0 z-20 bg-black/35 backdrop-blur-sm" />
		<Dialog.Content
			class="fixed top-1/2 left-1/2 z-30 w-[calc(100%-2rem)] max-w-md -translate-x-1/2 -translate-y-1/2 rounded-lg border border-line bg-surface p-6 shadow-black/20"
		>
			<div class="flex items-start justify-between gap-4">
				<div>
					<Dialog.Title class="text-xl font-semibold text-ink">Sortiranje</Dialog.Title>
					<Dialog.Description class="mt-1 text-sm text-muted">
						Izaberite po čemu se stavke ređaju.
					</Dialog.Description>
				</div>
				<Dialog.Close class="rounded-md px-2 py-1 text-sm text-muted hover:bg-soft hover:text-ink">
					Zatvori
				</Dialog.Close>
			</div>
			<form class="mt-6" onsubmit={applySort}>
				<p class="text-sm font-medium text-ink" id="inventory-sort-field-label">Svojstvo</p>
				<RadioGroup.Root
					class="mt-1 max-h-56 space-y-0.5 overflow-y-auto"
					bind:value={draftField}
					aria-labelledby="inventory-sort-field-label"
				>
					{#each fields as field (field.value)}
						<RadioGroup.Item
							value={field.value}
							class="flex w-full cursor-pointer items-center justify-between gap-2 rounded px-3 py-2 text-left text-sm text-ink outline-none hover:bg-soft focus-visible:bg-soft data-[state=checked]:bg-brand-soft data-[state=checked]:font-medium data-[state=checked]:text-brand"
						>
							{#snippet children({ checked })}
								<span class="truncate">{field.label}</span>
								{#if checked}<Check class="size-4 shrink-0" aria-hidden="true" />{/if}
							{/snippet}
						</RadioGroup.Item>
					{/each}
				</RadioGroup.Root>

				<p class="mt-4 text-sm font-medium text-ink" id="inventory-sort-order-label">Redosled</p>
				<div
					class="mt-1 grid grid-cols-2 gap-2"
					role="group"
					aria-labelledby="inventory-sort-order-label"
				>
					<Button.Root
						class={directionClass('asc')}
						type="button"
						aria-pressed={draftOrder === 'asc'}
						onclick={() => (draftOrder = 'asc')}
					>
						<ArrowUp class="size-4" aria-hidden="true" />
						Rastuće
					</Button.Root>
					<Button.Root
						class={directionClass('desc')}
						type="button"
						aria-pressed={draftOrder === 'desc'}
						onclick={() => (draftOrder = 'desc')}
					>
						<ArrowDown class="size-4" aria-hidden="true" />
						Opadajuće
					</Button.Root>
				</div>

				<Button.Root
					class="mt-6 w-full rounded-md bg-brand px-4 py-2 font-medium text-on-brand hover:bg-brand-strong"
					type="submit"
				>
					Primeni
				</Button.Root>
			</form>
		</Dialog.Content>
	</Dialog.Portal>
</Dialog.Root>
