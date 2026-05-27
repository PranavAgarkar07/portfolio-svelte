import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [sveltekit()],
	build: {
		rollupOptions: {
			output: {
				manualChunks(id: string) {
					if (id.includes('node_modules/gsap')) return 'vendor-gsap';
					if (id.includes('node_modules/lenis')) return 'vendor-lenis';
				}
			}
		}
	}
});
