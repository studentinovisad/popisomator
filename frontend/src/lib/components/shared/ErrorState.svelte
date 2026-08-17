<script lang="ts">
	import FallenClipboard from '$lib/components/shared/FallenClipboard.svelte';

	let {
		code,
		message
	}: {
		code?: number | string;
		message: string;
	} = $props();
</script>

<section
	class="grid w-full max-w-2xl items-center gap-3"
	aria-label={code ? `Greška ${code}` : 'Greška'}
>
	<div class="relative aspect-square w-full justify-self-center min-[380px]:max-w-72">
		<FallenClipboard {code} />
	</div>

	<div class="flex justify-center">
		<h1 class="sr-only">{message}</h1>
		<blockquote
			class="w-fit max-w-full text-center text-3xl leading-tight font-semibold text-balance text-ink sm:text-4xl"
		>
			{message}
		</blockquote>
	</div>
</section>
