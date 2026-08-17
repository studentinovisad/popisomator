<script lang="ts">
	import { resolve } from '$app/paths';
	import type { User } from '$lib/api';
	import UserAvatar from '$lib/components/app/UserAvatar.svelte';

	let {
		user,
		fullWidth = false,
		iconOnly = false,
		active = false
	}: { user: User; fullWidth?: boolean; iconOnly?: boolean; active?: boolean } = $props();
</script>

<a
	class={`inline-flex h-9 items-center text-sm font-medium transition-colors ${
		active
			? 'bg-brand-soft text-brand'
			: 'bg-chrome text-chrome-muted hover:bg-on-chrome/10 hover:text-on-chrome'
	} ${
		iconOnly
			? 'size-9 justify-center gap-0 rounded-md p-1'
			: fullWidth
				? 'min-w-0 flex-1 gap-2 rounded-l-md py-1 pr-2 pl-1'
				: 'max-w-64 gap-2 rounded-md pr-3 pl-1'
	}`}
	href={resolve('/account')}
	aria-label={`Moj nalog: ${user.full_name}`}
	aria-current={active ? 'page' : undefined}
	title={iconOnly ? `Moj nalog: ${user.full_name}` : undefined}
>
	<UserAvatar name={user.full_name} class="inline-flex size-7 shrink-0" />
	<span
		class={`min-w-0 overflow-hidden text-nowrap transition-[max-width,opacity] duration-200 ease-out ${
			iconOnly ? 'max-w-0 opacity-0' : 'max-w-48 opacity-100'
		}`}
	>
		{user.full_name}
	</span>
</a>
