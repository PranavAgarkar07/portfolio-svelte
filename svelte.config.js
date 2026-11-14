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
			base: process.env.NODE_ENV === 'production' ? '/portfolio-svelte' : '',
			// Emit base-absolute URLs (/portfolio-svelte/blog…) for every link —
			// the default per-page-relative mode resolves wrongly when a page is
			// served at an unexpected URL (e.g. 404.html fallback, /blog w/o slash).
			relative: false
		},
		prerender: {
			handleMissingId: 'ignore'
		},
		alias: {
			'$data': 'src/data'
		}
	}
};

export default config;
