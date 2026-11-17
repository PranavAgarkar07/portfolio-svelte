import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [sveltekit()],
	define: {
		// Build-scoped cache buster for the static stylesheets in /static
		// (styles.css, styles.critical.css, light-mode.css). Those files are
		// copied verbatim, so their URLs would otherwise be cached by the
		// browser across deploys and serve stale CSS.
		__CSS_BUILD__: JSON.stringify(Date.now().toString(36))
	},
	build: {
		rollupOptions: {
			output: {
				manualChunks(id: string) {
					if (id.includes('node_modules/gsap')) return 'vendor-gsap';
					if (id.includes('node_modules/lenis')) return 'vendor-lenis';
					if (id.includes('node_modules/lightbox3')) return 'vendor-lightbox';
					if (id.includes('node_modules/@humanspeak/svelte-motion')) return 'vendor-motion';
				}
			}
		}
	}
});
