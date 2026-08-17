<script lang="ts">
	import { Avatar } from 'bits-ui';
	import type { HTMLAttributes } from 'svelte/elements';

	let {
		name,
		class: className,
		fallbackClass = ''
	}: {
		name: string;
		class: string;
		fallbackClass?: string;
	} = $props();

	let initials = $derived(
		name
			.trim()
			.split(/\s+/)
			.slice(0, 2)
			.map((part) => part.at(0)?.toLocaleUpperCase('sr-Latn-RS'))
			.join('')
	);
</script>

{#snippet initialsAvatar({ props }: { props: Avatar.RootProps })}
	{@const avatarProps = props as unknown as HTMLAttributes<HTMLSpanElement>}
	<span {...avatarProps}>
		<Avatar.Fallback
			class={`flex size-full items-center justify-center rounded-full bg-brand-soft font-mono text-xs font-semibold text-brand-strong ${fallbackClass}`}
		>
			{initials}
		</Avatar.Fallback>
	</span>
{/snippet}

<Avatar.Root class={className} child={initialsAvatar} aria-hidden="true" />
