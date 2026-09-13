<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import ArrowLeft from '@lucide/svelte/icons/arrow-left';
	import ArrowRight from '@lucide/svelte/icons/arrow-right';
	import History from '@lucide/svelte/icons/history';
	import SquareArrowOutUpRight from '@lucide/svelte/icons/square-arrow-out-up-right';
	import { Portal } from 'bits-ui';
	import { onMount } from 'svelte';
	import { api, ApiError, type AuditEntry } from '$lib/api';
	import AuditEntryIcon from '$lib/components/admin/AuditEntryIcon.svelte';
	import ProtectedPageState from '$lib/components/shared/ProtectedPageState.svelte';
	import UserAvatar from '$lib/components/app/UserAvatar.svelte';
	import {
		auditActionClass,
		auditActionLabel,
		auditActorName,
		auditChangeLabel,
		auditChangeValue,
		auditEntryIcon,
		auditEntryTitle,
		auditHistoryLink,
		auditTargetLink,
		auditTargetName,
		auditTargetTypeAccusative,
		auditTargetTypeLabel
	} from '$lib/domain/audit-log';
	import { formatRequestDate } from '$lib/domain/item-requests';
	import { createAuthPage } from '$lib/state/auth-page.svelte';

	const authPage = createAuthPage({
		unavailableMessage: 'Dnevnik izmena trenutno nije dostupan.',
		requiredRole: 'admin'
	});

	let entry = $state<AuditEntry | null>(null);
	let entryError = $state('');

	const entryID = $derived(Number(page.params.id));

	onMount(() => void authPage.load());

	$effect(() => {
		if (!authPage.state.authorized) return;
		void loadEntry(entryID);
	});

	async function loadEntry(id: number) {
		entryError = '';
		try {
			entry = await api.getAuditLogEntry(id);
		} catch (reason) {
			entry = null;
			entryError = reason instanceof ApiError ? reason.message : 'Izmena nije učitana.';
		}
	}

	// A deleted target has nothing left to open, so that link is dropped - but its history outlives it
	// and stays reachable, which is the whole point of keeping the entry.
	const targetLink = $derived(entry ? auditTargetLink(entry) : null);
	const historyLink = $derived(entry ? auditHistoryLink(entry.target_type, entry.target_id) : null);
</script>

<svelte:head>
	<title>Izmena | Popisomator</title>
</svelte:head>

<main class="px-4 pt-4 pb-8 sm:px-6">
	<ProtectedPageState
		loading={authPage.state.loading}
		error={authPage.state.error || entryError}
		authorized={authPage.state.authorized}
	>
		<Portal to="#page-header-actions">
			<a
				class="inline-grid size-8 place-items-center rounded text-muted transition-colors hover:bg-soft hover:text-ink"
				href={resolve('/admin/audit-log')}
				aria-label="Nazad na dnevnik izmena"
				title="Nazad"
			>
				<ArrowLeft class="size-4" aria-hidden="true" />
			</a>
		</Portal>

		{#if entry}
			<section class="mx-auto max-w-4xl" aria-labelledby="audit-entry-heading">
				<div class="border-b border-line pb-5">
					<div class="flex items-center gap-2">
						<span
							class={`inline-flex items-center gap-1.5 rounded-full px-2 py-0.5 text-xs font-medium ${auditActionClass(entry.action)}`}
						>
							<AuditEntryIcon name={auditEntryIcon(entry)} size="size-3" />
							{auditActionLabel(entry.action)}
						</span>
						<span class="text-sm text-muted">{formatRequestDate(entry.created_at)}</span>
					</div>
					<h2 id="audit-entry-heading" class="mt-2 text-xl font-semibold text-ink">
						{auditEntryTitle(entry)}
					</h2>

					<div class="mt-3 flex flex-wrap items-center gap-2">
						{#if targetLink}
							<a
								class="inline-flex h-9 items-center gap-2 rounded-md border border-line bg-surface px-3 text-sm text-ink transition-colors hover:border-brand/40 hover:bg-brand-soft"
								href={resolve(targetLink)}
							>
								<SquareArrowOutUpRight class="size-4" aria-hidden="true" />
								Otvori {auditTargetTypeAccusative(entry.target_type)}
							</a>
						{/if}
						{#if historyLink}
							<a
								class="inline-flex h-9 items-center gap-2 rounded-md border border-line bg-surface px-3 text-sm text-ink transition-colors hover:border-brand/40 hover:bg-brand-soft"
								href={resolve(historyLink)}
							>
								<History class="size-4" aria-hidden="true" />
								Istorija izmena
							</a>
						{/if}
					</div>
				</div>

				<section class="mt-5" aria-labelledby="audit-entry-who-heading">
					<h3 id="audit-entry-who-heading" class="text-base font-semibold text-ink">Ko i šta</h3>
					<dl class="mt-3 grid gap-3 sm:grid-cols-2">
						<div class="rounded-md border border-line bg-surface px-4 py-3">
							<dt class="text-xs text-muted">Korisnik</dt>
							<dd class="mt-1 flex items-center gap-2">
								<UserAvatar name={auditActorName(entry)} class="inline-flex size-7 shrink-0" />
								<span class="min-w-0 truncate font-medium text-ink">{auditActorName(entry)}</span>
							</dd>
							{#if entry.actor_id === null && entry.actor_name}
								<p class="mt-1 text-xs text-muted">Nalog je u međuvremenu obrisan.</p>
							{/if}
						</div>
						<div class="rounded-md border border-line bg-surface px-4 py-3">
							<dt class="text-xs text-muted">
								{auditTargetTypeLabel(entry.target_type)}
							</dt>
							<dd class="mt-1 min-w-0 truncate font-medium text-ink">{auditTargetName(entry)}</dd>
							{#if !targetLink}
								<p class="mt-1 text-xs text-muted">Više ne postoji.</p>
							{/if}
						</div>
					</dl>
				</section>

				{#if entry.changes.length}
					<section class="mt-5" aria-labelledby="audit-entry-changes-heading">
						<h3 id="audit-entry-changes-heading" class="text-base font-semibold text-ink">
							Promenjene vrednosti
						</h3>
						<div class="-mx-4 mt-3 border-y border-line bg-surface sm:mx-0 sm:rounded-md sm:border">
							<ul class="divide-y divide-line">
								{#each entry.changes as change, index (index)}
									<li class="px-4 py-3">
										<p class="text-xs text-muted">{auditChangeLabel(change)}</p>
										<div class="mt-1 flex flex-wrap items-center gap-2 text-sm">
											{#if change.old !== undefined && change.old !== null}
												<span class="text-muted line-through"
													>{auditChangeValue(change, 'old')}</span
												>
												<ArrowRight class="size-3.5 shrink-0 text-muted" aria-hidden="true" />
											{/if}
											{#if change.new !== undefined && change.new !== null}
												<span class="font-medium text-ink">{auditChangeValue(change, 'new')}</span>
											{:else}
												<span class="text-muted">Uklonjeno</span>
											{/if}
										</div>
									</li>
								{/each}
							</ul>
						</div>
					</section>
				{/if}

				{#if entry.context.reason || entry.context.type_name || entry.context.batch_size || entry.context.properties?.length}
					<section class="mt-5" aria-labelledby="audit-entry-context-heading">
						<h3 id="audit-entry-context-heading" class="text-base font-semibold text-ink">
							Dodatni podaci
						</h3>
						<dl class="mt-3 space-y-3">
							{#if entry.context.reason}
								<div class="rounded-md border border-line bg-surface px-4 py-3">
									<dt class="text-xs text-muted">Razlog</dt>
									<dd class="mt-1 text-sm text-ink">{entry.context.reason}</dd>
								</div>
							{/if}
							{#if entry.context.type_name}
								<div class="rounded-md border border-line bg-surface px-4 py-3">
									<dt class="text-xs text-muted">Tip stavke</dt>
									<dd class="mt-1 text-sm text-ink">{entry.context.type_name}</dd>
								</div>
							{/if}
							{#if entry.context.batch_size && entry.context.batch_size > 1}
								<div class="rounded-md border border-line bg-surface px-4 py-3">
									<dt class="text-xs text-muted">Grupno dodavanje</dt>
									<dd class="mt-1 text-sm text-ink">
										Dodato je {entry.context.batch_size} stavki odjednom, svaka sa svojim zapisom.
									</dd>
								</div>
							{/if}
							{#if entry.context.properties?.length}
								<div class="rounded-md border border-line bg-surface px-4 py-3">
									<dt class="text-xs text-muted">Početna svojstva</dt>
									<dd class="mt-1 text-sm text-ink">
										{entry.context.properties.map((property) => property.name).join(', ')}
									</dd>
								</div>
							{/if}
						</dl>
					</section>
				{/if}
			</section>
		{/if}
	</ProtectedPageState>
</main>
