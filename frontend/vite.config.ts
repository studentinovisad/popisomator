import adapter from '@sveltejs/adapter-node';
import { sveltekit } from '@sveltejs/kit/vite';
import tailwindcss from '@tailwindcss/vite';
import { fileURLToPath } from 'node:url';
import { defineConfig, loadEnv } from 'vite';

const projectRoot = fileURLToPath(new URL('..', import.meta.url));

export default defineConfig(({ mode }) => {
	const environment = loadEnv(mode, projectRoot, 'POPISOMATOR_BACKEND_URL');
	const backendUrl = environment.POPISOMATOR_BACKEND_URL ?? 'http://localhost:8080';

	return {
		plugins: [
			tailwindcss(),
			sveltekit({
				compilerOptions: {
					runes: ({ filename }) =>
						filename.split(/[/\\]/).includes('node_modules') ? undefined : true
				},
				adapter: adapter()
			})
		],
		server: {
			proxy: {
				'/api': {
					target: backendUrl,
					rewrite: (path) => path.replace(/^\/api/, '')
				}
			}
		}
	};
});
