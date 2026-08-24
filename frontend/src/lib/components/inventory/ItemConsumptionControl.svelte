<script lang="ts">
	import { Button, Dialog, Label, Select } from 'bits-ui';
	import { ApiError, type ConsumptionStatus, type Item } from '$lib/api';
	import { consumptionClass, consumptionLabel, consumptionOptions } from '$lib/domain/items';

	let {
		item,
		canManage,
		disabled = false,
		class: sizeClass = 'h-9 min-w-0 flex-1 px-3 text-sm font-medium sm:w-44 sm:flex-none',
		onconsumptionchange,
		onrequest
	}: {
		item: Item;
		canManage: boolean;
		disabled?: boolean;
		class?: string;
		onconsumptionchange: (item: Item, status: ConsumptionStatus) => void;
		onrequest: (itemID: number, reason: string) => Promise<void>;
	} = $props();

	let canConsume = $derived(canManage || item.request_status === 'approved');
	let isTerminal = $derived(
		item.consumption === 'fully_consumed' || item.consumption === 'damaged'
	);

	let dialogOpen = $state(false);
	let reason = $state('');
	let submitting = $state(false);
	let error = $state('');

	async function submitRequest(event: SubmitEvent) {
		event.preventDefault();
		submitting = true;
		error = '';

		try {
			await onrequest(item.id, reason);
			reason = '';
			dialogOpen = false;
		} catch (caught) {
			error = caught instanceof ApiError ? caught.message : 'Zahtev nije poslat.';
		} finally {
			submitting = false;
		}
	}
</script>

{#if canConsume}
	<Select.Root
		type="single"
		value={item.consumption}
		items={consumptionOptions}
		{disabled}
		onValueChange={(value) => onconsumptionchange(item, value as ConsumptionStatus)}
	>
		<Select.Trigger
			class={`flex items-center justify-between rounded-md ${sizeClass} ${consumptionClass(item.consumption)}`}
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
{:else if isTerminal}
	<span
		class={`inline-flex items-center rounded-md ${sizeClass} ${consumptionClass(item.consumption)}`}
	>
		{consumptionLabel(item.consumption)}
	</span>
{:else if item.request_status === 'requested'}
	<span class={`inline-flex items-center rounded-md ${sizeClass} bg-warning-soft text-warning`}>
		Zahtev na čekanju
	</span>
{:else}
	<Dialog.Root bind:open={dialogOpen}>
		<Dialog.Trigger
			class={`inline-flex items-center justify-center rounded-md border border-line bg-surface ${sizeClass} text-ink transition-colors hover:border-brand/40 hover:bg-brand-soft hover:text-brand`}
		>
			Zatraži korišćenje
		</Dialog.Trigger>
		<Dialog.Portal>
			<Dialog.Overlay class="fixed inset-0 z-20 bg-black/35 backdrop-blur-sm" />
			<Dialog.Content
				class="fixed top-1/2 left-1/2 z-30 w-[calc(100%-2rem)] max-w-md -translate-x-1/2 -translate-y-1/2 rounded-lg border border-line bg-surface p-6 shadow-black/20"
			>
				<div class="flex items-start justify-between gap-4">
					<div>
						<Dialog.Title class="text-xl font-semibold text-ink">Zahtev za korišćenje</Dialog.Title>
						<Dialog.Description class="mt-1 text-sm text-muted">
							Objasnite administratoru zašto vam je potrebna ova stavka.
						</Dialog.Description>
					</div>
					<Dialog.Close
						class="rounded-md px-2 py-1 text-sm text-muted hover:bg-soft hover:text-ink"
					>
						Zatvori
					</Dialog.Close>
				</div>
				<form class="mt-6" onsubmit={submitRequest}>
					<Label.Root class="text-sm font-medium text-ink" for={`item-request-reason-${item.id}`}>
						Razlog <span class="font-normal text-muted">(opciono)</span>
					</Label.Root>
					<textarea
						id={`item-request-reason-${item.id}`}
						class="mt-1 block min-h-20 w-full"
						maxlength="400"
						bind:value={reason}></textarea>
					<Button.Root
						class="mt-4 rounded-md bg-brand px-4 py-2 font-medium text-on-brand hover:bg-brand-strong disabled:opacity-60"
						disabled={submitting}
						type="submit"
					>
						{submitting ? 'Slanje…' : 'Pošalji zahtev'}
					</Button.Root>
				</form>
				{#if error}<p class="mt-3 text-sm text-danger" role="alert">{error}</p>{/if}
			</Dialog.Content>
		</Dialog.Portal>
	</Dialog.Root>
{/if}
