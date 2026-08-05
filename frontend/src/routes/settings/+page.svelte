<script lang="ts">
	import BackendStatus from '$lib/components/BackendStatus.svelte';
	import ThemeToggle from '$lib/components/ThemeToggle.svelte';
</script>

<svelte:head>
	<title>Podešavanja | Popisomator</title>
</svelte:head>

<main class="pb-8">
	<section class="w-full" aria-label="Podešavanja">
		<div class="border-y border-line bg-surface">
			<div
				class="grid grid-cols-[minmax(0,1fr)_auto] gap-4 p-5 sm:grid-cols-[minmax(0,1fr)_24rem] sm:items-center sm:gap-6 sm:px-6"
			>
				<div>
					<h2 class="font-medium text-ink">Tema</h2>
					<p class="mt-1 text-sm text-muted">Izaberite svetlu ili tamnu temu.</p>
				</div>
				<div class="self-center justify-self-end"><ThemeToggle /></div>
			</div>

			<div
				class="grid gap-4 border-t border-line p-5 sm:grid-cols-[minmax(0,1fr)_24rem] sm:items-center sm:gap-6 sm:px-6"
			>
				<div>
					<h2 class="font-medium text-ink">Backend</h2>
					<p class="mt-1 text-sm text-muted">Trenutno stanje veze sa backend servisom.</p>
				</div>
				<div class="sm:justify-self-end"><BackendStatus /></div>
			</div>

			<div
				class="grid gap-4 border-t border-line p-5 sm:grid-cols-[minmax(0,1fr)_24rem] sm:items-center sm:gap-6 sm:px-6"
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
		</div>
	</section>
</main>
