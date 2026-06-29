import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	preprocess: vitePreprocess(),

	kit: {
		adapter: adapter({
			pages: 'build',
			assets: 'build',
			fallback: '404.html',
			precompress: true,
			strict: true
		}),
		paths: {
			base: process.env.NODE_ENV === 'production' ? '/portfolio-svelte' : ''
		},
		prerender: {
			handleMissingId: 'ignore'
		}
	}
};

export default config;
