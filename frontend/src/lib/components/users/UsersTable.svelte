<script lang="ts">
	import Eye from '@lucide/svelte/icons/eye';
	import { resolve } from '$app/paths';
	import type { User } from '$lib/api';
	import { userRoleLabel } from '$lib/domain/users';

	let { users }: { users: User[] } = $props();
</script>

<table class="hidden w-full table-fixed text-left text-sm lg:table">
	<colgroup>
		<col class="w-[35%]" />
		<col />
		<col class="w-44" />
		<col class="w-24" />
	</colgroup>
	<thead class="border-b border-line bg-soft text-muted">
		<tr class="h-12">
			<th class="px-4 py-3 font-medium">Ime</th>
			<th class="px-4 py-3 font-medium">Email</th>
			<th class="px-4 py-3 font-medium">Uloga</th>
			<th class="sr-only">Radnje</th>
		</tr>
	</thead>
	<tbody class="text-ink">
		{#each users as user (user.id)}
			<tr class="h-16 transition-colors hover:bg-soft/35">
				<td class="px-4 py-3 align-middle">
					<span class="block truncate" title={user.full_name}>{user.full_name}</span>
				</td>
				<td class="px-4 py-3 align-middle">
					<span class="block truncate" title={user.email}>{user.email}</span>
				</td>
				<td class="px-4 py-3 align-middle">
					<span class="text-muted">{userRoleLabel(user.role)}</span>
				</td>
				<td class="px-4 py-3 text-right align-middle">
					<a
						class="inline-grid size-8 place-items-center rounded text-muted hover:bg-soft hover:text-ink"
						href={resolve(`/users/${user.id}`)}
						aria-label={`Pregledaj korisnika ${user.full_name}`}
						title="Pregledaj"
					>
						<Eye class="size-4" aria-hidden="true" />
					</a>
				</td>
			</tr>
		{/each}
		{#if users.length === 0}
			<tr class="h-16">
				<td class="px-4 py-3 align-middle text-muted" colspan="4">Nema korisnika.</td>
			</tr>
		{/if}
	</tbody>
</table>
