<script lang="ts">
	import { onMount } from 'svelte';
	import { api, ApiError, type User, type UserStatus } from '$lib/api';
	import RegistrationApproval from '$lib/components/auth/RegistrationApproval.svelte';
	import PaginationFooter from '$lib/components/shared/PaginationFooter.svelte';
	import { createServerPagination } from '$lib/state/server-pagination.svelte';

	let { onempty }: { onempty: () => void } = $props();

	const registrationsPage = createServerPagination<User, { status: UserStatus }>({
		initialFilters: { status: 'requested' },
		loadPage: ({ limit, offset, status }) => api.listUsers({ limit, offset, status }),
		unavailableMessage: 'Zahtevi nisu učitani.'
	});

	onMount(() => void registrationsPage.load());

	$effect(() => {
		if (registrationsPage.loaded && registrationsPage.total === 0) onempty();
	});

	async function decideRegistration(user: User, approve: boolean) {
		registrationsPage.error = '';

		try {
			if (approve) {
				await api.updateUser(user.id, { status: 'active' });
			} else {
				await api.deleteUser(user.id);
			}

			registrationsPage.reloadAfterDelete();
		} catch (reason) {
			registrationsPage.error =
				reason instanceof ApiError ? reason.message : 'Zahtev nije obrađen.';
		}
	}
</script>

<section aria-labelledby="pending-registrations-heading">
	<h2 id="pending-registrations-heading" class="sr-only">Zahtevi za registraciju</h2>
	<p class="font-mono text-xs font-medium tracking-wide text-muted">
		UKUPNO: {registrationsPage.total}
	</p>
	{#if registrationsPage.error}
		<p class="mt-3 text-sm text-danger" role="alert">{registrationsPage.error}</p>
	{/if}
	<div class="-mx-4 mt-4 border-y border-line bg-surface sm:-mx-6">
		<table class="hidden w-full table-fixed text-left text-sm lg:table">
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
				{#each registrationsPage.items as user (user.id)}
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
				{#if registrationsPage.items.length === 0}
					<tr class="h-16"><td class="px-4 py-3 text-muted" colspan="3">Nema zahteva.</td></tr>
				{/if}
			</tbody>
		</table>
		<ul class="divide-y divide-line lg:hidden" aria-label="Zahtevi za registraciju">
			{#each registrationsPage.items as user (user.id)}
				<li class="px-4 py-3">
					<p class="truncate text-sm font-medium text-ink">{user.full_name}</p>
					<p class="mt-0.5 truncate text-sm text-muted">{user.email}</p>
					<div class="mt-3">
						<RegistrationApproval onclick={(approved) => void decideRegistration(user, approved)} />
					</div>
				</li>
			{/each}
			{#if registrationsPage.items.length === 0}
				<li class="px-4 py-3 text-sm text-muted">Nema zahteva.</li>
			{/if}
		</ul>
	</div>
	<PaginationFooter
		total={registrationsPage.total}
		perPage={registrationsPage.perPage}
		page={registrationsPage.currentPage}
		hasPreviousPage={registrationsPage.hasPreviousPage}
		hasNextPage={registrationsPage.hasNextPage}
		loading={registrationsPage.loading}
		onpagechange={registrationsPage.goToPage}
	/>
</section>
