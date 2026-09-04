<script lang="ts">
	import BackendStatus from '$lib/components/app/BackendStatus.svelte';
	import ThemeToggle from '$lib/components/app/ThemeToggle.svelte';
	import { pagination } from '$lib/state/pagination.svelte';
	import { Button } from 'bits-ui';
	import { toast } from 'svelte-sonner';

	function setRowsPerPage(event: Event) {
		const input = event.currentTarget as HTMLInputElement;
		pagination.set(Number(input.value));
		input.value = String(pagination.perPage);
	}
</script>

<svelte:head>
	<title>Podešavanja | Popisomator</title>
</svelte:head>

<main class="pb-8">
	<section class="w-full" aria-label="Podešavanja">
		<div class="border-y border-line bg-surface">
			<div class="grid grid-cols-[minmax(0,1fr)_auto] items-center gap-4 p-4 sm:gap-6 sm:px-6">
				<div>
					<h2 class="font-medium text-ink">Tema aplikacije</h2>
					<p class="mt-1 text-sm text-muted">Izaberite svetlu ili tamnu temu</p>
				</div>
				<div class="justify-self-end"><ThemeToggle /></div>
			</div>

			<div
				class="grid grid-cols-[minmax(0,1fr)_auto] items-center gap-4 border-t border-line p-4 sm:gap-6 sm:px-6"
			>
				<div>
					<h2 class="font-medium text-ink">Rezultata po stranici</h2>
					<p class="mt-1 text-sm text-muted">Odredite koliko se rezultata prikazuje u tabelama</p>
				</div>
				<input
					class="pagination-page-size size-10 justify-self-end rounded-md border border-line bg-surface px-1 text-center font-mono text-sm font-medium text-ink shadow-sm transition-colors hover:border-brand/40"
					type="number"
					min="5"
					max="50"
					value={pagination.perPage}
					aria-label="Rezultata po stranici"
					onchange={setRowsPerPage}
				/>
			</div>

			<div
				class="grid grid-cols-[minmax(0,1fr)_auto] items-center gap-4 border-t border-line p-4 sm:gap-6 sm:px-6"
			>
				<div>
					<h2 class="font-medium text-ink">Status backend servisa</h2>
					<p class="mt-1 text-sm text-muted">Trenutno stanje veze sa backend servisom</p>
				</div>
				<div class="justify-self-end"><BackendStatus /></div>
			</div>

			<div
				class="grid gap-3 border-t border-line p-4 sm:grid-cols-[minmax(0,1fr)_auto] sm:items-center sm:gap-6 sm:px-6"
			>
				<div>
					<h2 class="font-medium text-ink">Originalni autori projekta</h2>
					<p class="mt-1 text-sm text-muted">Klikom na imena možete videti izvorni kôd</p>
				</div>
				<a
					class="text-sm font-medium text-balance text-brand hover:text-brand-strong sm:w-96 sm:text-right"
					href="https://github.com/studentinovisad/popisomator"
					target="_blank"
					rel="noreferrer"
				>
					Đorđe Mančić · Matija Kljajić · Miša Stefanović
				</a>
			</div>

			{#if import.meta.env.DEV}
				<div class="border-t border-line p-4 sm:px-6">
					<h2 class="font-medium text-ink">Pregled obaveštenja</h2>
					<p class="mt-1 text-sm text-muted">Dostupno samo tokom razvoja.</p>
					<div class="mt-3 flex flex-wrap gap-2">
						<Button.Root
							onclick={() => toast.info('Ovo je sistemsko obaveštenje.')}
							class="h-9 cursor-pointer rounded-md border border-brand/40 bg-brand-soft px-3 text-sm font-medium text-ink transition-colors hover:border-brand"
						>
							Obaveštenje
						</Button.Root>
						<Button.Root
							onclick={() => toast.warning('Rok trajanja stavke se približava.')}
							class="h-9 cursor-pointer rounded-md border border-warning/40 bg-warning-soft px-3 text-sm font-medium text-ink transition-colors hover:border-warning"
						>
							Upozorenje
						</Button.Root>
						<Button.Root
							onclick={() => toast.error('Radnja nije uspela.')}
							class="h-9 cursor-pointer rounded-md border border-danger/40 bg-danger-soft px-3 text-sm font-medium text-ink transition-colors hover:border-danger"
						>
							Greška
						</Button.Root>
						<Button.Root
							onclick={() => toast.success('Stavka je uspešno sačuvana.')}
							class="h-9 cursor-pointer rounded-md border border-success/40 bg-success-soft px-3 text-sm font-medium text-ink transition-colors hover:border-success"
						>
							Uspeh
						</Button.Root>
					</div>
				</div>
			{/if}
		</div>
	</section>
</main>
