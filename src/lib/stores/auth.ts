import { writable, derived } from 'svelte/store';
import { browser } from '$app/environment';
import type { User } from '$lib/types';

const TOKEN_KEY = 'portfolio_jwt';
const BASE = (import.meta.env.VITE_API_URL || '').replace(/\/$/, '');

export const token = writable<string | null>(null);
export const user = writable<User | null>(null);
export const loading = writable(true);

export const isLoggedIn = derived(user, $u => $u !== null);
export const isAdmin = derived(user, $u => $u?.role === 'admin');
export const isAuthor = derived(user, $u => $u?.role === 'author' || $u?.role === 'admin');

function getStoredToken(): string | null {
  if (!browser) return null;
  try {
    return localStorage.getItem(TOKEN_KEY);
  } catch {
    return null;
  }
}

function setStoredToken(t: string | null) {
  if (!browser) return;
  try {
    if (t) localStorage.setItem(TOKEN_KEY, t);
    else localStorage.removeItem(TOKEN_KEY);
  } catch {}
}

async function fetchUser(t: string): Promise<User | null> {
  if (!BASE) return null;
  try {
    const r = await fetch(`${BASE}/api/auth/me`, {
      headers: { Authorization: `Bearer ${t}` },
    });
    if (!r.ok) return null;
    const data = await r.json();
    return data.user ?? null;
  } catch {
    return null;
  }
}

export async function init() {
  loading.set(true);
  const t = getStoredToken();
  if (t) {
    token.set(t);
    const u = await fetchUser(t);
    if (u) {
      user.set(u);
    } else {
      setStoredToken(null);
      token.set(null);
    }
  }
  loading.set(false);
}

export async function loginWithGoogle() {
  if (!BASE) return;
  try {
    localStorage.setItem('redirect_after_login', window.location.pathname + window.location.search);
    const r = await fetch(`${BASE}/api/auth/google-url`);
    if (!r.ok) return;
    const data = await r.json();
    if (data.url) {
      window.location.href = data.url;
    }
  } catch {}
}

export async function handleCallback(tokenStr: string) {
  setStoredToken(tokenStr);
  token.set(tokenStr);
  const u = await fetchUser(tokenStr);
  if (u) {
    user.set(u);
  }
}

export function logout() {
  setStoredToken(null);
  token.set(null);
  user.set(null);
}
