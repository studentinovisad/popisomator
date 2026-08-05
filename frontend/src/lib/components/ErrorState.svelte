<script lang="ts">
	let {
		code,
		message
	}: {
		code?: number | string;
		message: string;
	} = $props();
</script>

<section
	class="grid w-full max-w-4xl items-center gap-10 md:grid-cols-[18rem_minmax(0,1fr)]"
	aria-label={code ? `Greška ${code}` : 'Greška'}
>
	<div
		class="relative aspect-square w-full max-w-72 justify-self-center overflow-hidden rounded-md bg-soft md:order-1 md:justify-self-start"
	>
		{#if code}
			<span
				class="pointer-events-none absolute top-4 left-5 z-10 font-mono text-6xl font-bold tracking-tighter text-brand/30"
				aria-hidden="true"
			>
				{code}
			</span>
		{/if}
		<img class="h-full w-full object-cover" src="/sasa.jpg" alt="Saša Matić" />
	</div>

	<div class="md:order-2">
		<h1 class="sr-only">{message}</h1>
		<blockquote
			class="border-l-4 border-brand pl-5 text-3xl leading-tight font-semibold text-balance text-ink sm:text-4xl"
		>
			{message}
		</blockquote>
	</div>
</section>
