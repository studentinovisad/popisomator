export type Theme = 'light' | 'dark';

const storageKey = 'popisomator-theme';

class ThemePreference {
	current = $state<Theme>('light');

	initialize() {
		const storedTheme = localStorage.getItem(storageKey);
		const nextTheme =
			storedTheme === 'light' || storedTheme === 'dark'
				? storedTheme
				: window.matchMedia('(prefers-color-scheme: dark)').matches
					? 'dark'
					: 'light';

		this.apply(nextTheme);
	}

	set(nextTheme: Theme) {
		localStorage.setItem(storageKey, nextTheme);
		this.apply(nextTheme);
	}

	private apply(nextTheme: Theme) {
		document.documentElement.classList.add('theme-changing');
		this.current = nextTheme;
		document.documentElement.classList.toggle('dark', nextTheme === 'dark');
		document.documentElement.style.colorScheme = nextTheme;

		requestAnimationFrame(() => {
			document.documentElement.classList.remove('theme-changing');
		});
	}
}

export const theme = new ThemePreference();
