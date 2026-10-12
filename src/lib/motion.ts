/**
 * Shared motion loader — single gsap/ScrollTrigger registration point.
 * F-P0-1: dynamic import() only when motion is allowed, unified
 * `gsap/ScrollTrigger` path (never `gsap/dist/ScrollTrigger`).
 */

let registered = false;
let cached: { gsap: any; ScrollTrigger: any } | null = null;
let inflight: Promise<{ gsap: any; ScrollTrigger: any } | null> | null = null;

export function prefersReducedMotion(): boolean {
	if (typeof window === "undefined") return false;
	return window.matchMedia("(prefers-reduced-motion: reduce)").matches;
}

export function isTouchDevice(): boolean {
	if (typeof window === "undefined") return false;
	return "ontouchstart" in window || navigator.maxTouchPoints > 0;
}

/** Motion (GSAP + Lenis reveals) is only allowed on non-touch, full-motion devices. */
export function motionAllowed(): boolean {
	if (typeof window === "undefined") return false;
	return !prefersReducedMotion() && !isTouchDevice();
}

/**
 * Dynamically import gsap + ScrollTrigger, register once, apply defaults once.
 * Returns null when motion is not allowed (caller must leave content visible).
 */
export function ensureGsap(): Promise<{ gsap: any; ScrollTrigger: any } | null> {
	if (typeof window === "undefined") return Promise.resolve(null);
	if (cached) return Promise.resolve(cached);
	if (inflight) return inflight;
	inflight = (async (): Promise<{ gsap: any; ScrollTrigger: any } | null> => {
		if (!motionAllowed()) return null;
		const [{ default: gsap }, { ScrollTrigger }] = await Promise.all([
			import("gsap"),
			import("gsap/ScrollTrigger")
		]);
		if (!registered) {
			gsap.registerPlugin(ScrollTrigger);
			gsap.defaults({ overwrite: "auto" });
			registered = true;
		}
		cached = { gsap, ScrollTrigger };
		return cached;
	})();
	void inflight.then(
		() => {
			inflight = null;
		},
		() => {
			inflight = null;
		}
	);
	return inflight;
}

/** Lightbox images present on this page? Gate for lightbox3 dynamic import. */
export function hasLightboxTargets(): boolean {
	if (typeof document === "undefined") return false;
	return document.querySelector("[data-lightbox]") !== null;
}
