<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import Eye from '@lucide/svelte/icons/eye';
	import X from '@lucide/svelte/icons/x';
	import { Button, Select } from 'bits-ui';
	import { onMount } from 'svelte';
	import { api, type AuditActorOption, type AuditEntry, type AuditTargetType } from '$lib/api';
	import AuditEntryIcon from '$lib/components/admin/AuditEntryIcon.svelte';
	import PaginationFooter from '$lib/components/shared/PaginationFooter.svelte';
	import {
		auditActionAllowedFor,
		auditActionClass,
		auditActionFilterOptions,
		auditActionFilterOptionsFor,
		auditActionLabel,
		auditEntryDetail,
		auditEntryIcon,
		auditEntryTitle,
		auditTargetTypeAccusative,
		auditTargetTypeFilterOptions,
		type AuditActionFilter,
		type AuditTargetTypeFilter
	} from '$lib/domain/audit-log';
	import { formatRequestDate } from '$lib/domain/item-requests';
	import { createServerPagination } from '$lib/state/server-pagination.svelte';
	import { getTableFilter, getTablePage, updateTableQuery } from '$lib/state/table-query';

	let { onloaderror }: { onloaderror?: (message: string) => void } = $props();

	const entriesPage = createServerPagination<
		AuditEntry,
		{
			action: AuditActionFilter;
			targetType: AuditTargetTypeFilter;
			targetID: string;
			actorID: string;
		}
	>({
		initialFilters: { action: 'all', targetType: 'all', targetID: '', actorID: 'all' },
		loadPage: ({ limit, offset, action, targetType, targetID, actorID }) =>
			api.listAuditLog({
				limit,
				offset,
				action: action === 'all' ? undefined : action,
				targetType: targetType === 'all' ? undefined : targetType,
				targetID: targetID ? Number(targetID) : undefined,
				actorID: actorID === 'all' ? undefined : Number(actorID)
			}),
		unavailableMessage: 'Dnevnik izmena nije učitan.'
	});

	const actionOptions = $derived(auditActionFilterOptionsFor(entriesPage.filters.targetType));

	let actors = $state<AuditActorOption[]>([]);
	const actorFilterOptions = $derived([
		{ value: 'all', label: 'Svi korisnici' },
		...actors.map((actor) => ({ value: String(actor.id), label: actor.name }))
	]);

	// A specific target is not something you pick from a dropdown - there are far too many items for
	// that. It arrives on the URL, from the history button on an item or from an entry's own page, and
	// shows as a removable chip instead.
	let pinnedTarget = $derived(
		entriesPage.filters.targetID
			? {
					type: entriesPage.filters.targetType as AuditTargetType,
					id: Number(entriesPage.filters.targetID),
					label: entriesPage.items[0]?.target_label || ''
				}
			: null
	);

	onMount(() => {
		void loadActors();
	});

	$effect(() => {
		const url = page.url;

		const requestedTargetType = getTableFilter(url, 'target_type');
		const targetType = auditTargetTypeFilterOptions.some(
			(option) => option.value === requestedTargetType
		)
			? (requestedTargetType as AuditTargetTypeFilter)
			: 'all';

		const requestedAction = getTableFilter(url, 'action');
		const namedAction = auditActionFilterOptions.some((option) => option.value === requestedAction)
			? (requestedAction as AuditActionFilter)
			: 'all';
		const action = auditActionAllowedFor(namedAction, targetType) ? namedAction : 'all';

		// An id only means anything alongside a type, so it is dropped without one - the same rule the
		// backend enforces, applied here so it never sends a request it knows will be rejected.
		const requestedTargetID = Number(getTableFilter(url, 'target_id'));
		const targetID =
			targetType !== 'all' && Number.isSafeInteger(requestedTargetID) && requestedTargetID > 0
				? String(requestedTargetID)
				: '';

		const requestedActorID = Number(getTableFilter(url, 'actor_id'));
		const actorID =
			Number.isSafeInteger(requestedActorID) && requestedActorID > 0
				? String(requestedActorID)
				: 'all';

		entriesPage.sync({
			page: getTablePage(url),
			filters: { action, targetType, targetID, actorID }
		});
	});

	$effect(() => {
		if (entriesPage.error) onloaderror?.(entriesPage.error);
	});

	async function loadActors() {
		try {
			actors = await api.listAuditLogActors();
		} catch {
			actors = [];
		}
	}

	function filterByAction(action: AuditActionFilter) {
		updateTableQuery({ action: action === 'all' ? null : action, page: 1 });
	}

	function filterByTargetType(targetType: AuditTargetTypeFilter) {
		// Changing the kind of entity abandons any pinned id: it belonged to the previous kind.
		const next: Record<string, string | number | null> = {
			target_type: targetType === 'all' ? null : targetType,
			target_id: null,
			page: 1
		};
		// An action the new kind never produces would filter the list down to nothing, so it goes too.
		// The key is only set when it has to change - updateTableQuery drops whatever it is handed.
		if (!auditActionAllowedFor(entriesPage.filters.action, targetType)) next.action = null;
		updateTableQuery(next);
	}

	function filterByActor(actorID: string) {
		updateTableQuery({ actor_id: actorID === 'all' ? null : actorID, page: 1 });
	}

	function clearPinnedTarget() {
		updateTableQuery({ target_type: null, target_id: null, page: 1 });
	}
</script>

<section aria-labelledby="audit-log-heading">
	<h2 id="audit-log-heading" class="sr-only">Dnevnik izmena</h2>

	<div class="flex flex-wrap items-center justify-between gap-4">
		<p class="font-mono text-xs font-medium tracking-wide text-muted">
			UKUPNO: {entriesPage.total}
		</p>
		<div class="flex flex-wrap items-center gap-2">
			<Select.Root
				type="single"
				value={entriesPage.filters.actorID}
				items={actorFilterOptions}
				onValueChange={(value) => filterByActor(value)}
			>
				<Select.Trigger
					class="flex h-9 w-44 items-center justify-between rounded-md border border-line bg-surface px-3 text-sm text-ink transition-colors hover:border-brand/40"
					aria-label="Filtriraj izmene po korisniku"
				>
					<Select.Value />
				</Select.Trigger>
				<Select.Portal>
					<Select.Content
						class="z-30 max-h-72 w-52 overflow-y-auto rounded-md border border-line bg-surface p-1 shadow-lg shadow-black/15"
						sideOffset={4}
					>
						<Select.Viewport>
							{#each actorFilterOptions as option (option.value)}
								<Select.Item
									value={option.value}
									label={option.label}
									class="cursor-pointer rounded px-3 py-2 text-sm outline-none data-highlighted:bg-brand-soft"
								>
									{option.label}
								</Select.Item>
							{/each}
						</Select.Viewport>
					</Select.Content>
				</Select.Portal>
			</Select.Root>

			<Select.Root
				type="single"
				value={entriesPage.filters.targetType}
				items={auditTargetTypeFilterOptions}
				onValueChange={(value) => filterByTargetType(value as AuditTargetTypeFilter)}
			>
				<Select.Trigger
					class="flex h-9 w-36 items-center justify-between rounded-md border border-line bg-surface px-3 text-sm text-ink transition-colors hover:border-brand/40"
					aria-label="Filtriraj izmene po vrsti"
				>
					<Select.Value />
				</Select.Trigger>
				<Select.Portal>
					<Select.Content
						class="z-30 w-44 rounded-md border border-line bg-surface p-1 shadow-lg shadow-black/15"
						sideOffset={4}
					>
						<Select.Viewport>
							{#each auditTargetTypeFilterOptions as option (option.value)}
								<Select.Item
									value={option.value}
									label={option.label}
									class="cursor-pointer rounded px-3 py-2 text-sm outline-none data-highlighted:bg-brand-soft"
								>
									{option.label}
								</Select.Item>
							{/each}
						</Select.Viewport>
					</Select.Content>
				</Select.Portal>
			</Select.Root>

			<Select.Root
				type="single"
				value={entriesPage.filters.action}
				items={actionOptions}
				onValueChange={(value) => filterByAction(value as AuditActionFilter)}
			>
				<Select.Trigger
					class="flex h-9 w-52 items-center justify-between rounded-md border border-line bg-surface px-3 text-sm text-ink transition-colors hover:border-brand/40"
					aria-label="Filtriraj izmene po radnji"
				>
					<Select.Value />
				</Select.Trigger>
				<Select.Portal>
					<Select.Content
						class="z-30 max-h-72 w-64 overflow-y-auto rounded-md border border-line bg-surface p-1 shadow-lg shadow-black/15"
						sideOffset={4}
					>
						<Select.Viewport>
							{#each actionOptions as option (option.value)}
								<Select.Item
									value={option.value}
									label={option.label}
									class="cursor-pointer rounded px-3 py-2 text-sm outline-none data-highlighted:bg-brand-soft"
								>
									{option.label}
								</Select.Item>
							{/each}
						</Select.Viewport>
					</Select.Content>
				</Select.Portal>
			</Select.Root>
		</div>
	</div>

	{#if pinnedTarget}
		<div
			class="mt-3 flex items-center gap-2 rounded-md border border-brand/30 bg-brand-soft px-3 py-2 text-sm text-ink"
		>
			<span class="min-w-0 truncate">
				Istorija za {auditTargetTypeAccusative(pinnedTarget.type)}
				<span class="font-medium">{pinnedTarget.label || `#${pinnedTarget.id}`}</span>
			</span>
			<Button.Root
				class="ml-auto inline-grid size-7 shrink-0 place-items-center rounded text-muted transition-colors hover:bg-surface hover:text-ink"
				type="button"
				onclick={clearPinnedTarget}
				aria-label="Prikaži sve izmene"
				title="Prikaži sve izmene"
			>
				<X class="size-4" aria-hidden="true" />
			</Button.Root>
		</div>
	{/if}

	<div class="-mx-4 mt-4 border-y border-line bg-surface sm:-mx-6">
		<table class="hidden w-full table-fixed text-left text-sm lg:table">
			<colgroup>
				<col />
				<!-- The longest label, "Promena redosleda svojstava", needs 209px once the cell and badge
				padding are counted; w-52 gives 208 and wraps the pill onto two lines. -->
				<col class="w-56" />
				<col class="w-44" />
				<!-- Wide enough for the 32px button plus the cell's own padding, the way the inventory
				table sizes its actions column. Narrower and the button gets squeezed. -->
				<col class="w-24" />
			</colgroup>
			<thead class="border-b border-line bg-soft text-muted">
				<tr class="h-12">
					<th class="px-4 py-3 font-medium">Izmena</th>
					<th class="px-4 py-3 font-medium">Radnja</th>
					<th class="px-4 py-3 font-medium">Vreme</th>
					<th class="px-4 py-3"><span class="sr-only">Detalji</span></th>
				</tr>
			</thead>
			<tbody class="divide-y divide-line text-ink">
				{#each entriesPage.items as entry (entry.id)}
					{@const detail = auditEntryDetail(entry)}
					<tr class="h-16 transition-colors hover:bg-soft/35">
						<td class="px-4 py-3 align-middle">
							<a
								class="flex min-w-0 items-start gap-3"
								href={resolve(`/admin/audit-log/${entry.id}`)}
							>
								<span class="mt-0.5 shrink-0 text-muted" aria-hidden="true">
									<AuditEntryIcon name={auditEntryIcon(entry)} />
								</span>
								<span class="min-w-0">
									<span class="block truncate">{auditEntryTitle(entry)}</span>
									{#if detail}
										<span class="mt-0.5 block truncate text-xs text-muted" title={detail}>
											{detail}
										</span>
									{/if}
								</span>
							</a>
						</td>
						<td class="px-4 py-3 align-middle">
							<span
								class={`inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium whitespace-nowrap ${auditActionClass(entry.action)}`}
							>
								{auditActionLabel(entry.action)}
							</span>
						</td>
						<td class="px-4 py-3 align-middle text-muted">{formatRequestDate(entry.created_at)}</td>
						<td class="px-4 py-3 text-right align-middle">
							<div class="flex justify-end gap-1">
								<a
									class="inline-grid size-8 place-items-center rounded text-muted hover:bg-soft hover:text-ink"
									href={resolve(`/admin/audit-log/${entry.id}`)}
									aria-label={`Detalji izmene ${entry.id}`}
									title="Detalji"
								>
									<Eye class="size-4" aria-hidden="true" />
								</a>
							</div>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>

		<ul class="divide-y divide-line lg:hidden">
			{#each entriesPage.items as entry (entry.id)}
				{@const detail = auditEntryDetail(entry)}
				<li>
					<a
						class="flex items-start gap-3 px-4 py-3"
						href={resolve(`/admin/audit-log/${entry.id}`)}
					>
						<span class="mt-0.5 shrink-0 text-muted" aria-hidden="true">
							<AuditEntryIcon name={auditEntryIcon(entry)} />
						</span>
						<span class="min-w-0 flex-1">
							<span class="block text-sm text-ink">{auditEntryTitle(entry)}</span>
							{#if detail}
								<span class="mt-0.5 block truncate text-xs text-muted">{detail}</span>
							{/if}
							<span class="mt-1 flex items-center gap-2">
								<span
									class={`inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium whitespace-nowrap ${auditActionClass(entry.action)}`}
								>
									{auditActionLabel(entry.action)}
								</span>
								<span class="text-xs text-muted">{formatRequestDate(entry.created_at)}</span>
							</span>
						</span>
						<span
							class="mt-0.5 inline-grid size-8 shrink-0 place-items-center rounded text-muted"
							aria-hidden="true"
						>
							<Eye class="size-4" />
						</span>
					</a>
				</li>
			{/each}
		</ul>

		{#if entriesPage.items.length === 0 && entriesPage.loaded}
			<p class="px-4 py-3 text-sm text-muted">Nema zabeleženih izmena.</p>
		{/if}
	</div>

	<PaginationFooter
		total={entriesPage.total}
		perPage={entriesPage.perPage}
		page={entriesPage.currentPage}
		hasPreviousPage={entriesPage.hasPreviousPage}
		hasNextPage={entriesPage.hasNextPage}
		loading={entriesPage.loading}
		onpagechange={(nextPage) => updateTableQuery({ page: nextPage })}
	/>
</section>
