<script lang="ts">
	import { onMount } from 'svelte';
	import { api, ApiError, type User } from '$lib/api';
	import RegistrationApproval from '$lib/components/RegistrationApproval.svelte';
	import UsersPagination from '$lib/components/UsersPagination.svelte';
	import { pagination } from '$lib/pagination.svelte';

	let { onempty }: { onempty: () => void } = $props();

	let usersPerPage = $derived(pagination.perPage);

	let users = $state<User[]>([]);
	let userOffset = $state(0);
	let usersTotal = $state(0);
	let loading = $state(false);
	let error = $state('');
	let loadVersion = 0;
	let currentPage = $derived(Math.floor(userOffset / usersPerPage) + 1);
	let hasPreviousPage = $derived(userOffset > 0);
	let hasNextPage = $derived(userOffset + users.length < usersTotal);

	onMount(() => void loadUsers(0));

	async function loadUsers(offset: number) {
		const version = ++loadVersion;
		loading = true;
		error = '';

		try {
			const page = await api.listUsers({ limit: usersPerPage, offset, status: 'requested' });
			if (version !== loadVersion) return;

			users = page.items;
			userOffset = page.offset;
			usersTotal = page.total;

			if (page.total === 0) {
				onempty();
			}
		} catch (reason) {
			if (version !== loadVersion) return;

			error = reason instanceof ApiError ? reason.message : 'Zahtevi nisu učitani.';
		} finally {
			if (version === loadVersion) loading = false;
		}
	}

	async function decideRegistration(user: User, approve: boolean) {
		error = '';

		try {
			if (approve) {
				await api.approveRegistration(user.id);
			} else {
				await api.declineRegistration(user.id);
			}

			const nextOffset =
				users.length === 1 && userOffset > 0 ? userOffset - usersPerPage : userOffset;
			await loadUsers(nextOffset);
		} catch (reason) {
			error = reason instanceof ApiError ? reason.message : 'Zahtev nije obrađen.';
		}
	}

	function goToPage(page: number) {
		void loadUsers((page - 1) * usersPerPage);
	}
</script>

<section aria-labelledby="pending-registrations-heading">
	<h2 id="pending-registrations-heading" class="sr-only">Zahtevi za registraciju</h2>
	<p class="font-mono text-xs font-medium tracking-wide text-muted">UKUPNO: {usersTotal}</p>
	{#if error}
		<p class="mt-3 text-sm text-danger" role="alert">{error}</p>
	{/if}
	<div class="-mx-4 mt-4 border-y border-line bg-surface sm:-mx-6">
		<table class="hidden min-w-full table-fixed text-left text-sm md:table">
			<colgroup>
				<col class="w-[35%]" />
				<col />
				<col class="w-52" />
			</colgroup>
			<thead class="border-b border-line bg-soft text-muted">
				<tr class="h-12">
					<th class="px-4 py-3 font-medium">Ime</th>
					<th class="px-4 py-3 font-medium">Email</th>
					<th class="px-4 py-3 text-right font-medium">Radnje</th>
				</tr>
			</thead>
			<tbody class="text-ink">
				{#each users as user (user.id)}
					<tr class="h-16 transition-colors hover:bg-soft/35">
						<td class="px-4 py-3 align-middle"
							><span class="block truncate">{user.full_name}</span></td
						>
						<td class="px-4 py-3 align-middle"><span class="block truncate">{user.email}</span></td>
						<td class="px-4 py-3 align-middle">
							<RegistrationApproval
								onclick={(approved) => void decideRegistration(user, approved)}
							/>
						</td>
					</tr>
				{/each}
				{#if users.length === 0}
					<tr class="h-16"><td class="px-4 py-3 text-muted" colspan="3">Nema zahteva.</td></tr>
				{/if}
			</tbody>
		</table>
		<ul class="divide-y divide-line md:hidden" aria-label="Zahtevi za registraciju">
			{#each users as user (user.id)}
				<li class="px-4 py-3">
					<p class="truncate text-sm font-medium text-ink">{user.full_name}</p>
					<p class="mt-0.5 truncate text-sm text-muted">{user.email}</p>
					<div class="mt-3">
						<RegistrationApproval onclick={(approved) => void decideRegistration(user, approved)} />
					</div>
				</li>
			{/each}
			{#if users.length === 0}<li class="px-4 py-3 text-sm text-muted">Nema zahteva.</li>{/if}
		</ul>
	</div>
	<UsersPagination
		total={usersTotal}
		perPage={usersPerPage}
		page={currentPage}
		{hasPreviousPage}
		{hasNextPage}
		{loading}
		onpagechange={goToPage}
	/>
</section>
