<script lang="ts">
	import { untrack } from 'svelte';
	import { api, ApiError, type User, type UserRole } from '$lib/api';
	import UsersMobileList from '$lib/components/UsersMobileList.svelte';
	import PaginationFooter from '$lib/components/PaginationFooter.svelte';
	import UsersTable from '$lib/components/UsersTable.svelte';
	import UsersToolbar from '$lib/components/UsersToolbar.svelte';
	import { pagination } from '$lib/pagination.svelte';
	import type { UserRoleFilter } from '$lib/users';

	let { refreshKey, currentUserID }: { refreshKey: number; currentUserID: number } = $props();

	let usersPerPage = $derived(pagination.perPage);

	let users = $state<User[]>([]);
	let userOffset = $state(0);
	let usersTotal = $state(0);
	let pendingUsersTotal = $state(0);
	let usersLoading = $state(false);
	let error = $state('');
	let nameSearch = $state('');
	let activeNameSearch = $state('');
	let roleFilter = $state<UserRoleFilter>('all');
	let loadVersion = 0;
	let currentPage = $derived(Math.floor(userOffset / usersPerPage) + 1);
	let hasPreviousPage = $derived(userOffset > 0);
	let hasNextPage = $derived(userOffset + users.length < usersTotal);

	$effect(() => {
		void refreshKey;
		untrack(() => {
			void loadUsers(0, activeNameSearch, roleFilter);
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

	async function loadUsers(offset: number, search: string, role: UserRoleFilter) {
		const version = ++loadVersion;
		usersLoading = true;
		error = '';

		try {
			const page = await api.listUsers({
				limit: usersPerPage,
				offset,
				search,
				role: role === 'all' ? undefined : role,
				status: 'active'
			});
			if (version !== loadVersion) return;

			users = page.items;
			userOffset = page.offset;
			usersTotal = page.total;
		} catch (reason) {
			if (version !== loadVersion) return;

			error = reason instanceof ApiError ? reason.message : 'Korisnici nisu učitani.';
		} finally {
			if (version === loadVersion) {
				usersLoading = false;
			}
		}
	}

	async function updateRole(user: User, nextRole: UserRole) {
		error = '';

		try {
			const updatedUser = await api.updateUserRole(user.id, nextRole);
			users = users.map((listedUser) => (listedUser.id === user.id ? updatedUser : listedUser));
		} catch (reason) {
			error = reason instanceof ApiError ? reason.message : 'Uloga nije promenjena.';
		}
	}

	function goToPage(page: number) {
		void loadUsers((page - 1) * usersPerPage, activeNameSearch, roleFilter);
	}

	function searchUsers(search: string) {
		activeNameSearch = search;
		void loadUsers(0, search, roleFilter);
	}

	function filterByRole(role: UserRoleFilter) {
		roleFilter = role;
		void loadUsers(0, activeNameSearch, role);
	}
</script>

<section aria-labelledby="users-heading">
	<h2 id="users-heading" class="sr-only">Tabela korisnika</h2>
	<UsersToolbar
		total={usersTotal}
		hasPendingUsers={pendingUsersTotal > 0}
		role={roleFilter}
		bind:search={nameSearch}
		loading={usersLoading}
		onrolechange={filterByRole}
		onsearch={searchUsers}
	/>
	{#if error}
		<p class="mt-3 text-sm text-danger" role="alert">{error}</p>
	{/if}
	<div class="-mx-4 mt-4 border-y border-line bg-surface sm:-mx-6">
		<UsersMobileList {users} {currentUserID} onrolechange={updateRole} />
		<UsersTable {users} {currentUserID} onrolechange={updateRole} />
	</div>
	<PaginationFooter
		total={usersTotal}
		perPage={usersPerPage}
		page={currentPage}
		{hasPreviousPage}
		{hasNextPage}
		loading={usersLoading}
		onpagechange={goToPage}
	/>
</section>
