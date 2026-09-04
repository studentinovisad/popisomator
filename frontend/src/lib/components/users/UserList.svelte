<script lang="ts">
	import { page } from '$app/state';
	import { api, ApiError, type User, type UserRole } from '$lib/api';
	import UsersMobileList from '$lib/components/users/UsersMobileList.svelte';
	import PaginationFooter from '$lib/components/shared/PaginationFooter.svelte';
	import UsersTable from '$lib/components/users/UsersTable.svelte';
	import UsersToolbar from '$lib/components/users/UsersToolbar.svelte';
	import { createServerPagination } from '$lib/state/server-pagination.svelte';
	import { toast } from 'svelte-sonner';
	import { roleFilterOptions, type UserRoleFilter } from '$lib/domain/users';
	import {
		getTableFilter,
		getTablePage,
		getTableSearch,
		updateTableQuery
	} from '$lib/state/table-query';

	let {
		refreshKey,
		currentUserID,
		onloaderror
	}: {
		refreshKey: number;
		currentUserID: number;
		onloaderror?: (message: string) => void;
	} = $props();

	let pendingUsersTotal = $state(0);
	let search = $state('');
	let previousRefreshKey = $state<number | undefined>();
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
		const url = page.url;
		const nextSearch = getTableSearch(url);
		const requestedRole = getTableFilter(url, 'role');
		const role = roleFilterOptions.some((option) => option.value === requestedRole)
			? (requestedRole as UserRoleFilter)
			: 'all';

		search = nextSearch;
		usersPage.sync({ page: getTablePage(url), search: nextSearch, filters: { role } });
	});

	$effect(() => {
		if (previousRefreshKey === refreshKey) return;

		const shouldReloadUsers = previousRefreshKey !== undefined;
		previousRefreshKey = refreshKey;
		void loadPendingUsersTotal();
		if (shouldReloadUsers) void usersPage.load(usersPage.offset);
	});

	$effect(() => {
		if (usersPage.error) onloaderror?.(usersPage.error);
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
		try {
			const updatedUser = await api.updateUser(user.id, { role: nextRole });
			usersPage.items = usersPage.items.map((listedUser) =>
				listedUser.id === user.id ? updatedUser : listedUser
			);
			toast.success('Uloga korisnika je promenjena.');
		} catch (reason) {
			toast.error(reason instanceof ApiError ? reason.message : 'Uloga nije promenjena.');
		}
	}

	function filterByRole(role: UserRoleFilter) {
		updateTableQuery({ role: role === 'all' ? undefined : role, page: 1 });
	}

	function searchUsers(nextSearch: string) {
		updateTableQuery({ search: nextSearch, page: 1 });
	}

	function goToPage(nextPage: number) {
		updateTableQuery({ page: nextPage });
	}
</script>

<section aria-labelledby="users-heading">
	<h2 id="users-heading" class="sr-only">Tabela korisnika</h2>
	<UsersToolbar
		total={usersPage.total}
		hasPendingUsers={pendingUsersTotal > 0}
		role={usersPage.filters.role}
		bind:search
		loading={usersPage.loading}
		onrolechange={filterByRole}
		onsearch={searchUsers}
	/>
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
		onpagechange={goToPage}
	/>
</section>
