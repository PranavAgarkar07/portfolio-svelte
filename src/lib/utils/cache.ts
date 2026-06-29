import { browser } from '$app/environment';

const PREFIX = 'pcache_v2_';

interface Entry<T> {
  data: T;
  ts: number;
}

export function cacheGet<T>(key: string): T | null {
  if (!browser) return null;
  try {
    const raw = localStorage.getItem(PREFIX + key);
    if (!raw) return null;
    return (JSON.parse(raw) as Entry<T>).data;
  } catch { return null; }
}

export function cacheSet<T>(key: string, data: T): void {
  if (!browser) return;
  try {
    localStorage.setItem(PREFIX + key, JSON.stringify({ data, ts: Date.now() } as Entry<T>));
  } catch {}
}

export function cacheAge(key: string): number | null {
  if (!browser) return null;
  try {
    const raw = localStorage.getItem(PREFIX + key);
    if (!raw) return null;
    const entry = JSON.parse(raw) as Entry<unknown>;
    return Date.now() - entry.ts;
  } catch { return null; }
}

export function cacheDel(key: string): void {
  if (!browser) return;
  try { localStorage.removeItem(PREFIX + key); } catch {}
}

export function cacheClear(pattern?: string): void {
  if (!browser) return;
  try {
    const fullPrefix = PREFIX + (pattern || '');
    const toRemove: string[] = [];
    for (let i = 0; i < localStorage.length; i++) {
      const k = localStorage.key(i);
      if (k?.startsWith(fullPrefix)) toRemove.push(k);
    }
    toRemove.forEach(k => localStorage.removeItem(k));
  } catch {}
}

export async function staleWhileRevalidate<T>(
  key: string,
  url: string,
  ttlMs: number,
  onData: (data: T) => void,
  parse?: (raw: unknown) => T,
  init?: RequestInit,
): Promise<boolean> {
  const cached = cacheGet<T>(key);
  const age = cacheAge(key);
  const isStale = age !== null && age > ttlMs;

  if (cached && !isStale) {
    onData(cached);
    return true;
  }

  if (cached && isStale) {
    onData(cached);
  }

  try {
    const r = await fetch(url, init);
    if (!r.ok) return false;
    const json = await r.json();
    const data: T = parse ? parse(json) : (json as T);
    cacheSet(key, data);
    onData(data);
    return true;
  } catch { return false; }
}
