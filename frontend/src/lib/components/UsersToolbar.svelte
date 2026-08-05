<script lang="ts">
	import { Button, Portal, Select } from 'bits-ui';
	import { roleFilterOptions, type UserRoleFilter } from '$lib/users';

	let {
		total,
		role,
		search = $bindable(),
		loading,
		onrolechange,
		onsearch
	}: {
		total: number;
		role: UserRoleFilter;
		search: string;
		loading: boolean;
		onrolechange: (role: UserRoleFilter) => void;
		onsearch: (search: string) => void;
	} = $props();

	function submit(event: SubmitEvent) {
		event.preventDefault();
		onsearch(search.trim());
	}
</script>

<Portal to="#users-summary">
	<div class="flex items-center justify-between gap-4">
		<p class="font-mono text-xs leading-none font-medium tracking-wide text-muted">
			UKUPNO: {total}
		</p>
		<div class="flex items-center gap-2">
			<Select.Root
				type="single"
				value={role}
				items={roleFilterOptions}
				onValueChange={(value) => onrolechange(value as UserRoleFilter)}
			>
				<Select.Trigger
					class="flex h-9 w-32 items-center justify-between rounded-md border border-chrome-line bg-transparent px-3 text-sm text-on-chrome transition-colors hover:border-brand"
					aria-label="Filtriraj korisnike po ulozi"
				>
					<Select.Value />
				</Select.Trigger>
				<Select.Portal>
					<Select.Content
						class="z-10 w-36 rounded-md border border-line bg-surface p-1 shadow-lg shadow-ink/10"
						sideOffset={4}
					>
						<Select.Viewport>
							{#each roleFilterOptions as option (option.value)}
								<Select.Item
									value={option.value}
									label={option.label}
									class="cursor-pointer rounded px-3 py-2 outline-none data-highlighted:bg-brand-soft"
								>
									{option.label}
								</Select.Item>
							{/each}
						</Select.Viewport>
					</Select.Content>
				</Select.Portal>
			</Select.Root>
		</div>
	</div>
</Portal>

<form class="flex flex-col gap-2 sm:flex-row sm:items-end" onsubmit={submit}>
	<div class="min-w-0 flex-1">
		<div class="relative mt-3">
			<svg
				aria-hidden="true"
				viewBox="0 0 24 24"
				fill="none"
				stroke="currentColor"
				stroke-width="2"
				class="pointer-events-none absolute top-1/2 left-3 size-4 -translate-y-1/2 text-faint"
			>
				<circle cx="11" cy="11" r="6" />
				<path d="m16 16 4 4" />
			</svg>
			<input
				id="user-name-search"
				class="h-10 w-full pl-9"
				bind:value={search}
				placeholder="Pretraži po imenu"
			/>
		</div>
	</div>
	<Button.Root
		class="h-10 rounded-md bg-brand px-4 text-sm font-medium text-on-brand transition-colors hover:bg-brand-strong disabled:cursor-not-allowed disabled:opacity-50"
		disabled={loading}
		type="submit"
	>
		Pretraži
	</Button.Root>
</form>
