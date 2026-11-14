import { browser } from "$app/environment";

const SESSION_KEY = "portfolio_analytics_session";
const EVENTS_KEY = "portfolio_analytics_events";

let sessionId: string | null = null;

function generateUUID(): string {
  if (crypto.randomUUID) return crypto.randomUUID();
  return `${Date.now()}-${Math.random().toString(36).slice(2, 10)}`;
}

function getDeviceType(): string {
  const ua = navigator.userAgent;
  if (/tablet|ipad|playbook|silk/i.test(ua)) return "tablet";
  if (/mobile|iphone|ipod|android.*mobile|blackberry|windows phone/i.test(ua)) return "mobile";
  return "desktop";
}

function getOS(): string {
  const ua = navigator.userAgent;
  if (/windows/i.test(ua)) return "Windows";
  if (/mac os|macintosh/i.test(ua)) return "macOS";
  if (/linux/i.test(ua)) return "Linux";
  if (/android/i.test(ua)) return "Android";
  if (/iphone|ipad|ipod/i.test(ua)) return "iOS";
  return "Unknown";
}

function getBrowser(): string {
  const ua = navigator.userAgent;
  if (/edg/i.test(ua)) return "Edge";
  if (/firefox/i.test(ua)) return "Firefox";
  if (/chrome/i.test(ua)) return "Chrome";
  if (/safari/i.test(ua)) return "Safari";
  return "Unknown";
}

function getTheme(): string {
  try {
    return localStorage.getItem("theme") || "dark";
  } catch {
    return "dark";
  }
}

function getSource(): string {
  const params = new URLSearchParams(window.location.search);
  const utm = params.get("utm_source");
  if (utm) return utm;

  const ref = document.referrer || "";
  if (!ref) return "direct";
  if (ref.includes("linkedin.com")) return "linkedin";
  if (ref.includes("facebook.com") || ref.includes("fb.com") || ref.includes("fb.me")) return "facebook";
  if (ref.includes("instagram.com")) return "instagram";
  if (ref.includes("wa.me") || ref.includes("whatsapp.com")) return "whatsapp";
  if (ref.includes("x.com") || ref.includes("twitter.com") || ref.includes("t.co")) return "x";
  if (ref.includes("github.com")) return "github";
  try {
    const u = new URL(ref);
    return u.hostname;
  } catch {
    return "other";
  }
}

function getOrCreateSessionId(): string {
  try {
    let id = localStorage.getItem(SESSION_KEY);
    if (!id) {
      id = generateUUID();
      localStorage.setItem(SESSION_KEY, id);
    }
    return id;
  } catch {
    return generateUUID();
  }
}

export function getSessionId(): string | null {
  if (!browser) return null;
  if (!sessionId) sessionId = getOrCreateSessionId();
  return sessionId;
}

function hasAnalyticsConsent(): boolean {
  try {
    return localStorage.getItem("portfolio_analytics_consent") === "true";
  } catch {
    return false;
  }
}

export function initAnalytics(): void {
  if (!browser) return;
  if (!hasAnalyticsConsent()) return;
  sessionId = getOrCreateSessionId();
}

export function pageview(path: string): void {
  if (!browser || !hasAnalyticsConsent()) return;
  const session = getSessionId();
  if (!session) return;
  const event: AnalyticsEvent = {
    type: "pageview",
    path,
    session_id: session,
    ts: Date.now(),
    referrer: document.referrer,
    source: getSource(),
    device: getDeviceType(),
    os: getOS(),
    browser: getBrowser(),
    theme: getTheme(),
  };
  pushEvent(event);
}

export function event(name: string, data?: Record<string, string>): void {
  if (!browser || !hasAnalyticsConsent()) return;
  const session = getSessionId();
  if (!session) return;
  const event: AnalyticsEvent = {
    type: "event",
    name,
    data,
    session_id: session,
    ts: Date.now(),
  };
  pushEvent(event);
}

export function dumpAnalytics(): AnalyticsEvent[] {
  if (!browser) return [];
  try {
    const raw = localStorage.getItem(EVENTS_KEY);
    return raw ? JSON.parse(raw) : [];
  } catch {
    return [];
  }
}

export function clearAnalytics(): void {
  if (!browser) return;
  try {
    localStorage.removeItem(EVENTS_KEY);
  } catch {}
}

interface AnalyticsEvent {
  type: "pageview" | "event";
  session_id: string;
  ts: number;
  path?: string;
  name?: string;
  data?: Record<string, string>;
  referrer?: string;
  source?: string;
  device?: string;
  os?: string;
  browser?: string;
  theme?: string;
}

function pushEvent(event: AnalyticsEvent): void {
  try {
    const raw = localStorage.getItem(EVENTS_KEY);
    const events: AnalyticsEvent[] = raw ? JSON.parse(raw) : [];
    events.push(event);
    // Keep only last 500 events to avoid unbounded growth
    if (events.length > 500) events.splice(0, events.length - 500);
    localStorage.setItem(EVENTS_KEY, JSON.stringify(events));
  } catch {}
}
