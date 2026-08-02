<svelte:head>
	<title>Popisomator</title>
	<meta name="description" content="Popisomator" />
</svelte:head>

<main class="mx-auto max-w-4xl px-6 py-16">
	<h1 class="text-3xl font-semibold text-slate-900">Popisomator</h1>

	<div class="mt-8 grid items-center gap-8 md:grid-cols-2">
		<blockquote
			class="text-3xl leading-tight font-semibold text-balance text-slate-800 sm:text-4xl"
		>
			Gde je PRIS projekat? Ja ovde još ništa ne vidim?
		</blockquote>
		<img class="w-full rounded-lg" src="/sasa.jpg" alt="Saša Matić" />
	</div>

	<ul class="mt-8 flex flex-wrap text-sm text-slate-500" aria-label="Original authors">
		{#each ['Đorđe Mančić', 'Matija Kljajić', 'Miša Stefanović'] as tim8_c (tim8_c)}
			<li class="after:mx-1.5 after:content-['·'] last:after:content-none">{tim8_c}</li>
		{/each}
	</ul>
</main>
