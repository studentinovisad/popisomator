<script lang="ts">
	import { page } from '$app/state';
	import { getContext, onDestroy, onMount } from 'svelte';
	import Printer from '@lucide/svelte/icons/printer';
	import { Button } from 'bits-ui';
	import { api, type Dashboard, type ItemTypeOption } from '$lib/api';
	import ConsumptionByMonthChart from '$lib/components/dashboard/ConsumptionByMonthChart.svelte';
	import ConsumptionQuantityReport from '$lib/components/dashboard/ConsumptionQuantityReport.svelte';
	import DashboardCard from '$lib/components/dashboard/DashboardCard.svelte';
	import DashboardFilters from '$lib/components/dashboard/DashboardFilters.svelte';
	import ExpiringItemsList from '$lib/components/dashboard/ExpiringItemsList.svelte';
	import ExpiryByMonthChart from '$lib/components/dashboard/ExpiryByMonthChart.svelte';
	import MostConsumedTable from '$lib/components/dashboard/MostConsumedTable.svelte';
	import StockGroupsCard from '$lib/components/dashboard/StockGroupsCard.svelte';
	import ProtectedPageState from '$lib/components/shared/ProtectedPageState.svelte';
	import { itemCountLabel, monthRangeLabel, parseDashboardMonths } from '$lib/domain/dashboard';
	import { createAuthPage } from '$lib/state/auth-page.svelte';
	import {
		dashboardReportPrintContextKey,
		type DashboardReportPrintContext
	} from '$lib/state/dashboard-report-print-context';
	import { getTableFilter, updateTableQuery } from '$lib/state/table-query';

	const printContext = getContext<DashboardReportPrintContext>(dashboardReportPrintContextKey);

	const authPage = createAuthPage({
		unavailableMessage: 'Kontrolna tabla trenutno nije dostupna.',
		requiredRoles: ['manager', 'admin']
	});

	let dashboard = $state<Dashboard | null>(null);
	let itemTypes = $state<ItemTypeOption[]>([]);
	let dashboardError = $state('');
	let loading = $state(false);
	let loadedQueryKey = '';
	let loadVersion = 0;

	const months = $derived(parseDashboardMonths(getTableFilter(page.url, 'months')));
	const typeID = $derived(getTableFilter(page.url, 'type_id'));
	const rangeLabel = $derived(monthRangeLabel(months));

	// One type's groups all carry its name, which is worth saying once in the filter instead of on
	// every row.
	const showTypeName = $derived(!typeID);

	onMount(() => {
		void authPage.load().then(() => {
			if (authPage.state.authorized) void loadItemTypes();
		});
	});

	onDestroy(() => {
		printContext.setDashboardReport(null);
	});

	$effect(() => {
		if (!authPage.state.authorized) return;

		const queryKey = `${months}|${typeID}`;
		if (queryKey === loadedQueryKey) return;
		loadedQueryKey = queryKey;

		void loadDashboard(months, typeID);
	});

	async function loadDashboard(selectedMonths: typeof months, selectedTypeID: string) {
		const version = ++loadVersion;
		loading = true;
		dashboardError = '';

		try {
			const next = await api.getDashboard({
				months: selectedMonths,
				typeID: selectedTypeID ? Number(selectedTypeID) : undefined
			});
			if (version !== loadVersion) return;
			dashboard = next;
			printContext.setDashboardReport(next);
		} catch {
			if (version !== loadVersion) return;
			dashboardError = 'Kontrolna tabla nije učitana.';
		} finally {
			if (version === loadVersion) loading = false;
		}
	}

	async function loadItemTypes() {
		try {
			itemTypes = await api.getItemTypeOptions();
		} catch {
			// The filter is an extra, not the page: without the list it stays on "Svi tipovi" and every
			// figure below is still correct.
			itemTypes = [];
		}
	}
</script>

<svelte:head>
	<title>Kontrolna tabla | Popisomator</title>
</svelte:head>

<main class="px-4 pt-4 pb-8 sm:px-6">
	<ProtectedPageState
		loading={authPage.state.loading}
		contentLoaded={dashboard !== null}
		error={authPage.state.error || dashboardError}
		authorized={authPage.state.authorized}
	>
		<div class="flex flex-wrap items-center justify-between gap-3">
			<DashboardFilters
				{months}
				{typeID}
				{itemTypes}
				disabled={loading}
				onmonthschange={(next) => updateTableQuery({ months: next })}
				ontypechange={(next) => updateTableQuery({ type_id: next === 'all' ? null : next })}
			/>
			<Button.Root
				type="button"
				class="inline-flex h-9 items-center justify-center gap-2 rounded-md border border-line bg-surface px-3 text-sm text-ink transition-colors hover:border-brand/40 hover:bg-brand-soft"
				disabled={!dashboard}
				onclick={() => printContext.print()}
			>
				<Printer class="size-4" aria-hidden="true" />
				<span>Odštampaj izveštaj</span>
			</Button.Root>
		</div>

		{#if dashboard}
			{@const summary = dashboard}
			<div class={`mt-4 grid gap-4 xl:grid-cols-2 ${loading ? 'opacity-60' : ''}`}>
				<DashboardCard
					title="Rokovi koji ističu"
					subtitle={`Stavke na stanju kojima rok ističe u narednih ${rangeLabel}.`}
					empty={summary.expiring_by_month.every((bucket) => bucket.count === 0)}
					emptyMessage="Nijednoj stavki rok ne ističe u izabranom periodu."
				>
					{#snippet actions()}
						{#if summary.expired_backlog > 0}
							<span
								class="inline-flex items-center rounded-full bg-danger-soft px-2 py-0.5 text-xs font-medium whitespace-nowrap text-danger"
							>
								{itemCountLabel(summary.expired_backlog)} sa isteklim rokom
							</span>
						{/if}
					{/snippet}
					<ExpiryByMonthChart buckets={summary.expiring_by_month} />
				</DashboardCard>

				<DashboardCard
					title="Potrošnja"
					subtitle={`Šta je skinuto sa stanja u proteklih ${rangeLabel}, po ishodu.`}
					empty={summary.consumption_by_month.every(
						(bucket) => bucket.fully_consumed + bucket.partially_consumed + bucket.damaged === 0
					)}
					emptyMessage="U izabranom periodu ništa nije evidentirano kao potrošeno."
				>
					<ConsumptionByMonthChart buckets={summary.consumption_by_month} />
				</DashboardCard>

				<DashboardCard
					title="Uskoro ističe"
					subtitle="Stavke kojima je rok istekao ili se bliži pragu svog tipa."
					scroll
					empty={summary.expiring_items.length === 0}
					emptyMessage="Nijednoj stavki rok ne ističe uskoro."
				>
					{#snippet actions()}
						<span class="text-xs whitespace-nowrap text-muted">
							{summary.expiring_items.length < summary.expiring_items_total
								? `prikazano ${summary.expiring_items.length} od ${summary.expiring_items_total}`
								: itemCountLabel(summary.expiring_items_total)}
						</span>
					{/snippet}
					<ExpiringItemsList items={summary.expiring_items} />
				</DashboardCard>

				<StockGroupsCard groups={summary.stock_groups} {showTypeName} />
			</div>

			<div class={`mt-4 grid gap-4 ${loading ? 'opacity-60' : ''}`}>
				<DashboardCard
					title="Potrošnja i troškovi po mesecima"
					subtitle={`Iskorišćene količine i vrednost, po tipu, u proteklih ${rangeLabel}.`}
					empty={summary.consumption_quantity.length === 0}
					emptyMessage="Nijedan tip nema evidentiranu potrošnju u izabranom periodu."
				>
					<ConsumptionQuantityReport types={summary.consumption_quantity} />
				</DashboardCard>

				<DashboardCard
					title="Najviše korišćeno"
					subtitle={`Nazivi sa najviše potrošnje u proteklih ${rangeLabel}.`}
					empty={summary.most_consumed.length === 0}
					emptyMessage="U izabranom periodu ništa nije evidentirano kao potrošeno."
				>
					<MostConsumedTable groups={summary.most_consumed} {showTypeName} />
				</DashboardCard>
			</div>
		{/if}
	</ProtectedPageState>
</main>
