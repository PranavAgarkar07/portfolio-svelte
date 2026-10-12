/**
 * Shared backend warmup — single deduped health-check + timeout fetch.
 * F-P0-4: shared between DevLog.svelte + ContactForm.svelte.
 * - Dedupes concurrent warmups (singleflight).
 * - Cooldown: reuses a recent success for 30s without refetching.
 * - Every fetch has an AbortController + timeout (default 5s).
 */

const WARM_COOLDOWN_MS = 30_000;
const DEFAULT_TIMEOUT_MS = 5_000;

let warmInflight: Promise<boolean> | null = null;
let lastWarmSuccess = 0;

/** fetch() with AbortController timeout. Throws on abort/network. */
export async function fetchWithTimeout(
	url: string,
	init: RequestInit = {},
	timeoutMs: number = DEFAULT_TIMEOUT_MS
): Promise<Response> {
	const ctrl = new AbortController();
	const timer = setTimeout(() => ctrl.abort(), timeoutMs);
	try {
		return await fetch(url, { ...init, signal: ctrl.signal });
	} finally {
		clearTimeout(timer);
	}
}

/**
 * Deduped backend wake-up ping. Resolves true when backend answers ok,
 * false otherwise. Concurrent callers share one request; a success within
 * the cooldown window resolves true without any request.
 */
export function warmBackend(healthUrl: string, timeoutMs: number = DEFAULT_TIMEOUT_MS): Promise<boolean> {
	if (typeof fetch === "undefined") return Promise.resolve(false);
	if (Date.now() - lastWarmSuccess < WARM_COOLDOWN_MS) return Promise.resolve(true);
	if (warmInflight) return warmInflight;
	warmInflight = fetchWithTimeout(healthUrl, {}, timeoutMs)
		.then((res) => {
			if (res.ok) lastWarmSuccess = Date.now();
			return res.ok;
		})
		.catch(() => false)
		.finally(() => {
			warmInflight = null;
		}) as Promise<boolean>;
	return warmInflight;
}

/** Test/SSR helper — clears dedup state. */
export function resetWarmupCache(): void {
	warmInflight = null;
	lastWarmSuccess = 0;
}
