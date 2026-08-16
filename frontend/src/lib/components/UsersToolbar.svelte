<script lang="ts">
	import { resolve } from '$app/paths';
	import Search from '@lucide/svelte/icons/search';
	import { Button, Portal, Select } from 'bits-ui';
	import { roleFilterOptions, type UserRoleFilter } from '$lib/users';

	let {
		total,
		hasPendingUsers,
		role,
		search = $bindable(),
		loading,
		onrolechange,
		onsearch
	}: {
		total: number;
		hasPendingUsers: boolean;
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
			{#if hasPendingUsers}
				<a
					class="pending-requests-link inline-flex h-9 items-center rounded-md border border-brand bg-brand-soft px-3 text-sm font-medium text-brand transition-colors hover:border-brand hover:bg-brand-soft"
					href={resolve('/admin/users/pending')}
				>
					Zahtevi
				</a>
			{/if}
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
						class="z-10 w-36 rounded-md border border-line bg-surface p-1 shadow-lg shadow-black/15"
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
			<Search
				class="pointer-events-none absolute top-1/2 left-3 size-4 -translate-y-1/2 text-faint"
				aria-hidden="true"
			/>
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
