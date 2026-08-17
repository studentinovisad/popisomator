<script lang="ts">
	let { code }: { code?: number | string } = $props();
</script>

<svg
	class="h-full w-full"
	viewBox="35 35 250 240"
	fill="none"
	role="img"
	aria-label="Ispuštena kontrolna tabla"
	xmlns="http://www.w3.org/2000/svg"
>
	<path d="M42 256H278" class="stroke-line" stroke-width="2" stroke-linecap="round" />
	<ellipse cx="163" cy="262" rx="97" ry="14" class="fill-black/10" />
	<g transform="rotate(-11 160 163)">
		<rect
			x="88"
			y="58"
			width="148"
			height="198"
			rx="14"
			class="fill-surface stroke-line"
			stroke-width="3"
		/>
		<path
			d="M134 70V61C134 53.268 140.268 47 148 47H176C183.732 47 190 53.268 190 61V70"
			class="stroke-muted"
			stroke-width="7"
			stroke-linecap="round"
		/>
		<rect x="108" y="89" width="108" height="144" rx="6" class="fill-soft" />
		<path d="M126 118H192" class="stroke-line" stroke-width="8" stroke-linecap="round" />
		<path d="M126 148H178" class="stroke-line" stroke-width="8" stroke-linecap="round" />
		<path d="M126 178H156" class="stroke-line" stroke-width="8" stroke-linecap="round" />
		<text
			x="188"
			y="184"
			class="fill-brand"
			font-size="22"
			font-family="var(--font-mono)"
			font-weight="600"
			text-anchor="middle">{code ?? '!'}</text
		>
		<path d="M126 208H181" class="stroke-line" stroke-width="8" stroke-linecap="round" />
	</g>
</svg>
