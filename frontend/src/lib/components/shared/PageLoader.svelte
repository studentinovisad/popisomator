<script lang="ts">
	import { Progress } from 'bits-ui';
</script>

<div class="grid min-h-[calc(100svh-14rem)] place-items-center" aria-busy="true">
	<Progress.Root
		value={null}
		class="size-10 animate-spin rounded-full border-4 border-line border-t-brand"
		aria-label="Učitavanje"
	/>
</div>
