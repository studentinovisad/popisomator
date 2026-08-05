<script lang="ts">
	import { resolve } from '$app/paths';
	import type { User } from '$lib/api';
	import UserAvatar from '$lib/components/UserAvatar.svelte';

	let {
		user,
		fullWidth = false,
		active = false
	}: { user: User; fullWidth?: boolean; active?: boolean } = $props();
</script>

<a
	class={`inline-flex h-9 items-center gap-2 text-sm font-medium transition-colors ${
		active
			? 'bg-brand-soft text-brand'
			: 'bg-chrome text-chrome-muted hover:bg-on-chrome/10 hover:text-on-chrome'
	} ${fullWidth ? 'min-w-0 flex-1 rounded-l-md py-1 pr-2 pl-1' : 'max-w-64 rounded-md pr-3 pl-1'}`}
	href={resolve('/account')}
	aria-label={`Moj nalog: ${user.full_name}`}
	aria-current={active ? 'page' : undefined}
>
	<UserAvatar name={user.full_name} class="inline-flex size-7 shrink-0" />
	<span class="truncate">{user.full_name}</span>
</a>
