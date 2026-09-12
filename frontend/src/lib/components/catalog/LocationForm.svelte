<script lang="ts">
	import {
		api,
		ApiError,
		type CreateLocationRequest,
		type Location,

		type LocationOption

	} from '$lib/api';
	import { requiredTextError } from '$lib/domain/form-validation';
	import { Button, Label, Portal, Select, Separator } from 'bits-ui';
	import Check from '@lucide/svelte/icons/check';
	import Plus from '@lucide/svelte/icons/plus';
	import X from '@lucide/svelte/icons/x';
	import { toast } from 'svelte-sonner';
	import OptionCombobox from '../shared/OptionCombobox.svelte';

	let {
		location,
		locationOptions,
		onsaved,
		oncancel
	}: {
		location?: Location;
		locationOptions: LocationOption[];
		onsaved: () => void;
		oncancel?: () => void;
	} = $props();

	let filteredLocationOptions = $derived(locationOptions.filter((option) => option.id != location?.id))
	let name = $state('');
	let description = $state('');
	let hasParentLocation = $state(false);
	let parentLocationID: string = $state('');
	let saving = $state(false);
	let nameError = $state('');

	$effect(() => {
		name = location?.name ?? '';
		description = location?.description ?? '';
		if (location?.parent_id != undefined) {
			hasParentLocation = true;
			parentLocationID = String(location?.parent_id) ?? '';
		}
	});

	async function save(event: SubmitEvent) {
		event.preventDefault();
		nameError = requiredTextError(name, 'naziv lokacije') ?? '';
		if (nameError) return;
		saving = true;

		try {
			const savedParentID = hasParentLocation ? Number(parentLocationID) : null;
			if (location) {
				await api.updateLocation(location.id, {
					name,
					description,
					parent_id: savedParentID
				});
			} else {
				const payload: CreateLocationRequest = {
					name,
					description,
					parent_id: savedParentID
				};
				await api.createLocation(payload);
				name = '';
				description = '';
			}
			toast.success(location ? 'Lokacija je izmenjena.' : 'Lokacija je dodata.');
			onsaved();
		} catch (reason) {
			toast.error(reason instanceof ApiError ? reason.message : 'Lokacija nije sačuvana.');
		} finally {
			saving = false;
		}
	}
</script>

<Portal to="#page-header-actions">
	<Button.Root
		class="inline-grid size-10 place-items-center rounded-md bg-brand text-on-brand hover:bg-brand-strong disabled:opacity-60"
		disabled={saving}
		form="location-form"
		type="submit"
		aria-label={location ? 'Sačuvaj izmene' : 'Dodaj lokaciju'}
		title={location ? 'Sačuvaj izmene' : 'Dodaj lokaciju'}
	>
		{#if location}
			<Check class="size-4" aria-hidden="true" />
		{:else}
			<Plus class="size-4" aria-hidden="true" />
		{/if}
	</Button.Root>
	{#if oncancel}
		<Button.Root
			class="inline-grid size-10 place-items-center rounded-md border border-line bg-surface text-ink hover:border-brand/40 hover:bg-brand-soft hover:text-brand disabled:opacity-60"
			disabled={saving}
			type="button"
			onclick={oncancel}
			aria-label={location ? 'Odbaci izmene' : 'Otkaži dodavanje'}
			title={location ? 'Odbaci izmene' : 'Otkaži dodavanje'}
		>
			<X class="size-4" aria-hidden="true" />
		</Button.Root>
	{/if}
</Portal>

<form id="location-form" class="grid w-full gap-5" novalidate onsubmit={save}>
	<div class="grid gap-4">
		<div class="grid gap-4 sm:grid-cols-2">
			<div>
				<Label.Root class="text-sm font-medium text-ink" for="location-name">Naziv</Label.Root>
				<input
					id="location-name"
					class={`mt-1 block w-full ${nameError ? 'field-invalid' : ''}`}
					bind:value={name}
					aria-invalid={Boolean(nameError)}
					aria-describedby={nameError ? 'location-name-error' : undefined}
					oninput={() => (nameError = '')}
				/>
				<p id="location-name-error" class="mt-1 min-h-4 text-xs text-danger" aria-live="polite">
					{nameError}
				</p>
			</div>
			
		</div>
		<div>
			<Label.Root class="text-sm font-medium text-ink" for="location-description">Opis</Label.Root>
			<textarea
				id="location-description"
				class="mt-1 block min-h-24 w-full"
				bind:value={description}></textarea>
		</div>
	</div>
	<Separator.Root class="h-px bg-line" decorative />
	<fieldset>
		<label class="flex cursor-pointer items-center gap-2 text-sm font-medium text-ink">
			<input type="checkbox" bind:checked={hasParentLocation} />
			Postavi nadlokaciju
		</label>
		{#if hasParentLocation}
			<div class="mt-4">
				<OptionCombobox
					id="parent-location"
					options={filteredLocationOptions}
					bind:value={parentLocationID}
					placeholder="Odaberite lokaciju"
				/>
			</div>
		{/if}
	</fieldset>
</form>
