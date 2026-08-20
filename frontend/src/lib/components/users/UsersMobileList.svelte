<script lang="ts">
	import type { User, UserRole } from '$lib/api';
	import RoleSelect from '$lib/components/auth/RoleSelect.svelte';

	let {
		users,
		currentUserID,
		onrolechange
	}: {
		users: User[];
		currentUserID: number;
		onrolechange: (user: User, role: UserRole) => void;
	} = $props();
</script>

<ul class="divide-y divide-line lg:hidden" aria-label="Korisnici">
	{#each users as user (user.id)}
		<li class="px-4 py-3">
			<div class="grid grid-cols-[minmax(0,1fr)_8rem] items-center gap-4">
				<div class="min-w-0">
					<p class="truncate text-sm font-medium text-ink" title={user.full_name}>
						{user.full_name}
					</p>
					<p class="mt-0.5 truncate text-sm text-muted" title={user.email}>{user.email}</p>
				</div>
				<div>
					<RoleSelect
						value={user.role}
						ariaLabel={`Uloga za ${user.full_name}`}
						disabled={user.id === currentUserID}
						onvaluechange={(role) => onrolechange(user, role)}
					/>
				</div>
			</div>
		</li>
	{/each}
	{#if users.length === 0}
		<li class="px-4 py-3 text-sm text-muted">Nema korisnika.</li>
	{/if}
</ul>
