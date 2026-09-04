<script lang="ts">
	import Eye from '@lucide/svelte/icons/eye';
	import EyeOff from '@lucide/svelte/icons/eye-off';
	import { Toggle } from 'bits-ui';

	let {
		id,
		value = $bindable(),
		autocomplete,
		required = false,
		minlength,
		className = '',
		invalid = false,
		describedBy,
		oninput
	}: {
		id: string;
		value: string;
		autocomplete?: 'current-password' | 'new-password';
		required?: boolean;
		minlength?: number;
		className?: string;
		invalid?: boolean;
		describedBy?: string;
		oninput?: (event: Event) => void;
	} = $props();

	let visible = $state(false);
</script>

<div class={`relative ${className}`}>
	<input
		{id}
		class={`block w-full pr-10 ${invalid ? 'field-invalid' : ''}`}
		type={visible ? 'text' : 'password'}
		bind:value
		{autocomplete}
		{required}
		{minlength}
		aria-invalid={invalid}
		aria-describedby={describedBy}
		{oninput}
	/>
	<Toggle.Root
		class="absolute top-1/2 right-1 inline-grid size-8 -translate-y-1/2 place-items-center rounded text-muted transition-colors hover:bg-soft hover:text-ink"
		type="button"
		aria-label={visible ? 'Sakrij lozinku' : 'Prikaži lozinku'}
		bind:pressed={visible}
	>
		{#if visible}
			<EyeOff class="size-4" aria-hidden="true" />
		{:else}
			<Eye class="size-4" aria-hidden="true" />
		{/if}
	</Toggle.Root>
</div>
