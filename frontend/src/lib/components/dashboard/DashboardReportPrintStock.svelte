<script lang="ts">
	import type { Dashboard } from '$lib/api';
	import { formatPropertyTotals } from '$lib/domain/dashboard';

	let { report }: { report: Dashboard } = $props();
</script>

<article
	class="dashboard-report-stock mx-auto max-w-3xl px-6 py-10 text-ink sm:px-10"
	aria-hidden="true"
>
	<header class="border-b-2 border-ink pb-6">
		<h1 class="mt-1 text-2xl font-semibold tracking-tight">Izveštaj o trenutnom stanju zaliha</h1>
	</header>

	<section class="mt-8">
		<table class="mt-3 w-full text-sm">
			<thead>
				<tr class="border-b border-ink text-left text-xs text-muted">
					<th class="py-1.5 pr-3 font-medium">Tip</th>
					<th class="py-1.5 pr-3 font-medium">Naziv</th>
					<th class="py-1.5 pr-3 text-right font-medium">Na stanju</th>
					<th class="py-1.5 text-right font-medium">Količina</th>
				</tr>
			</thead>
			<tbody>
				{#each report.stock_groups as group (`${group.type_id}-${group.name}`)}
					<tr class="border-b border-line last:border-0">
						<td class="py-1.5 pr-3 text-muted">{group.type_name}</td>
						<td class="py-1.5 pr-3">{group.name}</td>
						<td class="py-1.5 pr-3 text-right tabular-nums">{group.in_stock_count}</td>
						<td class="py-1.5 text-right tabular-nums">
							{group.totals.length > 0 ? formatPropertyTotals(group.totals) : '—'}
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</section>
</article>

<style>
	.dashboard-report-stock {
		display: none;
	}

	@media print {
		:global(body) {
			background: white;
		}

		:global(.app-shell) {
			display: none;
		}

		.dashboard-report-stock {
			display: block;
			max-width: none;
			padding: 0;
			color-scheme: light;
			--color-canvas: #f7f7f6;
			--color-surface: #ffffff;
			--color-ink: #272726;
			--color-muted: #6e6e6b;
			--color-faint: #737370;
			--color-line: #e2e2df;
			--color-soft: #f0f0ee;
		}

		.dashboard-report-stock tr {
			break-inside: avoid;
		}
	}
</style>
