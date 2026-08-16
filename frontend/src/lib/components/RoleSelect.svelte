<script lang="ts">
	import { Select } from 'bits-ui';
	import type { UserRole } from '$lib/api';
	import { userRoleOptions } from '$lib/users';

	let {
		id,
		value = $bindable<UserRole>(),
		ariaLabel,
		disabled = false,
		onvaluechange
	}: {
		id?: string;
		value: UserRole;
		ariaLabel: string;
		disabled?: boolean;
		onvaluechange?: (value: UserRole) => void;
	} = $props();
</script>

<Select.Root
	type="single"
	bind:value={value as never}
	items={userRoleOptions}
	{disabled}
	onValueChange={(value) => onvaluechange?.(value as UserRole)}
>
	<Select.Trigger
		{id}
		class="flex w-full items-center justify-between rounded-md border border-line bg-surface px-3 py-2 text-left text-ink hover:border-muted disabled:cursor-not-allowed disabled:opacity-50"
		aria-label={ariaLabel}
		{disabled}
	>
		<Select.Value />
	</Select.Trigger>
	<Select.Portal>
		<Select.Content
			class="z-40 w-40 rounded-md border border-line bg-surface p-1 shadow-lg shadow-black/15"
			sideOffset={4}
		>
			<Select.Viewport>
				{#each userRoleOptions as role (role.value)}
					<Select.Item
						value={role.value}
						label={role.label}
						class="cursor-pointer rounded px-3 py-2 outline-none data-highlighted:bg-brand-soft"
					>
						{role.label}
					</Select.Item>
				{/each}
			</Select.Viewport>
		</Select.Content>
	</Select.Portal>
</Select.Root>
