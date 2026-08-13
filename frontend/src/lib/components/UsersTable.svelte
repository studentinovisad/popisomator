<script lang="ts">
	import type { User, UserRole } from '$lib/api';
	import RoleSelect from '$lib/components/RoleSelect.svelte';
	import RegistrationApproval from './RegistrationApproval.svelte';

	let {
		users,
		currentUserID,
		onrolechange,
		declineuser,
		approveuser
	}: {
		users: User[];
		currentUserID: number;
		onrolechange: (user: User, role: UserRole) => void;
		declineuser: (user: User) => void;
		approveuser: (user: User) => void;
	} = $props();
</script>

<table class="hidden min-w-full table-fixed text-left text-sm md:table">
	<colgroup>
		<col class="w-[35%]" />
		<col />
		<col class="w-44" />
		<col class="w-0" />
	</colgroup>
	<thead class="border-b border-line bg-soft text-muted">
		<tr class="h-12">
			<th class="px-4 py-3 font-medium">Ime</th>
			<th class="px-4 py-3 font-medium">Email</th>
			<th class="px-4 py-3 font-medium">Uloga</th>
			<th></th>
		</tr>
	</thead>
	<tbody class="text-ink">
		{#each users as user (user.id)}
			<tr class="h-16 transition-colors hover:bg-soft/35">
				<td class="px-4 py-3 align-middle">
					<span class="block truncate" title={user.full_name}>{user.full_name}</span>
					{#if user.status !== 'active'}
						<p class="truncate font-medium text-success">Zahtev za registraciju</p>
					{/if}
				</td>
				<td class="px-4 py-3 align-middle"
					><span class="block truncate" title={user.email}>{user.email}</span></td
				>
				<td class="px-4 py-3 align-middle">
					{#if user.status === 'active'}
						<RoleSelect
							value={user.role}
							ariaLabel={`Uloga za ${user.full_name}`}
							disabled={user.id === currentUserID}
							onvaluechange={(role) => onrolechange(user, role)}
						/>
					{/if}
				</td>
				<td class="px-4 align-middle">
					{#if user.id !== currentUserID}
						{#if user.status !== 'active'}
							<RegistrationApproval
								onclick={(approve) => (approve ? approveuser(user) : declineuser(user))}
							/>
						{/if}
					{/if}
				</td>
			</tr>
		{/each}
		{#if users.length === 0}
			<tr class="h-16">
				<td class="px-4 py-3 align-middle text-muted" colspan="3">Nema korisnika.</td>
			</tr>
		{/if}
	</tbody>
</table>
