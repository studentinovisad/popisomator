<script lang="ts">
	import { untrack } from 'svelte';
	import { api, ApiError, type User, type UserRole } from '$lib/api';
	import UsersMobileList from '$lib/components/users/UsersMobileList.svelte';
	import PaginationFooter from '$lib/components/shared/PaginationFooter.svelte';
	import UsersTable from '$lib/components/users/UsersTable.svelte';
	import UsersToolbar from '$lib/components/users/UsersToolbar.svelte';
	import { createServerPagination } from '$lib/state/server-pagination.svelte';
	import type { UserRoleFilter } from '$lib/domain/users';

	let { refreshKey, currentUserID }: { refreshKey: number; currentUserID: number } = $props();

	let pendingUsersTotal = $state(0);
	const usersPage = createServerPagination<User, { role: UserRoleFilter }>({
		initialFilters: { role: 'all' },
		loadPage: ({ limit, offset, search, role }) =>
			api.listUsers({
				limit,
				offset,
				search,
				role: role === 'all' ? undefined : role,
				status: 'active'
			}),
		unavailableMessage: 'Korisnici nisu učitani.'
	});

	$effect(() => {
		void refreshKey;
		untrack(() => {
			void usersPage.load();
			void loadPendingUsersTotal();
		});
	});

	async function loadPendingUsersTotal() {
		try {
			const page = await api.listUsers({ limit: 1, status: 'requested' });
			pendingUsersTotal = page.total;
		} catch {
			pendingUsersTotal = 0;
		}
	}

	async function updateRole(user: User, nextRole: UserRole) {
		usersPage.error = '';

		try {
			const updatedUser = await api.updateUser(user.id, { role: nextRole });
			usersPage.items = usersPage.items.map((listedUser) =>
				listedUser.id === user.id ? updatedUser : listedUser
			);
		} catch (reason) {
			usersPage.error = reason instanceof ApiError ? reason.message : 'Uloga nije promenjena.';
		}
	}

	function filterByRole(role: UserRoleFilter) {
		usersPage.setFilters({ role });
	}
</script>

<section aria-labelledby="users-heading">
	<h2 id="users-heading" class="sr-only">Tabela korisnika</h2>
	<UsersToolbar
		total={usersPage.total}
		hasPendingUsers={pendingUsersTotal > 0}
		role={usersPage.filters.role}
		bind:search={usersPage.search}
		loading={usersPage.loading}
		onrolechange={filterByRole}
		onsearch={usersPage.searchBy}
	/>
	{#if usersPage.error}
		<p class="mt-3 text-sm text-danger" role="alert">{usersPage.error}</p>
	{/if}
	<div class="-mx-4 mt-4 border-y border-line bg-surface sm:-mx-6">
		<UsersMobileList users={usersPage.items} {currentUserID} onrolechange={updateRole} />
		<UsersTable users={usersPage.items} {currentUserID} onrolechange={updateRole} />
	</div>
	<PaginationFooter
		total={usersPage.total}
		perPage={usersPage.perPage}
		page={usersPage.currentPage}
		hasPreviousPage={usersPage.hasPreviousPage}
		hasNextPage={usersPage.hasNextPage}
		loading={usersPage.loading}
		onpagechange={usersPage.goToPage}
	/>
</section>
