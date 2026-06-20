/// <reference types="@sveltejs/kit" />
/// <reference no-default-lib="true"/>
/// <reference lib="esnext" />
/// <reference lib="webworker" />

const sw = self as unknown as ServiceWorkerGlobalScope;

const CACHE = "portfolio-v1";

const ASSETS = [
	"/",
	"/showcase",
];

const IGNORED_PATTERNS = [
	/^https?:\/\/.*(?:s3\.amazonaws\.com|execute-api\.)/,
];

sw.addEventListener("install", (event) => {
	event.waitUntil(
		caches.open(CACHE).then((cache) => cache.addAll(ASSETS)),
	);
	sw.skipWaiting();
});

sw.addEventListener("activate", (event) => {
	event.waitUntil(
		caches.keys().then((keys) =>
			Promise.all(keys.filter((k) => k !== CACHE).map((k) => caches.delete(k))),
		),
	);
	sw.clients.claim();
});

sw.addEventListener("fetch", (event) => {
	if (IGNORED_PATTERNS.some((p) => p.test(event.request.url))) return;

	if (event.request.method !== "GET") return;

	if (event.request.url.includes("/api/")) {
		event.respondWith(networkFirst(event.request));
		return;
	}

	if (
		event.request.url.includes("/_app/immutable/") ||
		event.request.url.match(/\.(js|css|woff2?|png|webp|svg|jpg|jpeg)$/)
	) {
		event.respondWith(cacheFirst(event.request));
		return;
	}

	event.respondWith(networkFirst(event.request));
});

async function cacheFirst(request: Request): Promise<Response> {
	const cached = await caches.match(request);
	if (cached) return cached;
	const response = await fetch(request);
	if (response.ok) {
		const cache = await caches.open(CACHE);
		cache.put(request, response.clone());
	}
	return response;
}

async function networkFirst(request: Request): Promise<Response> {
	try {
		const response = await fetch(request);
		if (response.ok) {
			const cache = await caches.open(CACHE);
			cache.put(request, response.clone());
		}
		return response;
	} catch {
		const cached = await caches.match(request);
		if (cached) return cached;
		return new Response("Offline", { status: 503 });
	}
}
