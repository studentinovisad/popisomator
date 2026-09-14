<script lang="ts">
	import Eye from '@lucide/svelte/icons/eye';
	import { resolve } from '$app/paths';
	import type { User } from '$lib/api';
	import { userRoleLabel } from '$lib/domain/users';

	let { users }: { users: User[] } = $props();
</script>

<ul class="divide-y divide-line lg:hidden" aria-label="Korisnici">
	{#each users as user (user.id)}
		<li class="px-4 py-3">
			<div class="grid grid-cols-[minmax(0,1fr)_auto] gap-x-3 gap-y-2">
				<div class="min-w-0">
					<p class="truncate text-sm font-medium text-ink" title={user.full_name}>
						{user.full_name}
					</p>
					<p class="mt-0.5 truncate text-sm text-muted" title={user.email}>{user.email}</p>
				</div>
				<a
					class="inline-grid size-8 place-items-center rounded text-muted hover:bg-soft hover:text-ink"
					href={resolve(`/users/${user.id}`)}
					aria-label={`Pregledaj korisnika ${user.full_name}`}
					title="Pregledaj"
				>
					<Eye class="size-4" aria-hidden="true" />
				</a>
				<span class="col-span-2 text-sm text-muted">{userRoleLabel(user.role)}</span>
			</div>
		</li>
	{/each}
	{#if users.length === 0}
		<li class="px-4 py-3 text-sm text-muted">Nema korisnika.</li>
	{/if}
</ul>
