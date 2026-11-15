// Prerender the listing and every post: the content is committed static data
// (src/data/blog.ts), and GitHub Pages has no SPA fallback — a prerender=false
// blog means /blog and all deep links 404 in production.
export const prerender = true;
