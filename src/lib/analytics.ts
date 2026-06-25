const SESSION_KEY = "portfolio_analytics_session";
const BASE = (import.meta.env.VITE_API_URL || "").replace(/\/$/, "");
const ANALYTICS_BASE = BASE + "/api/v1/analytics";

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
  return sessionId;
}

export async function initAnalytics(): Promise<void> {
  if (!BASE || sessionId) return;

  sessionId = getOrCreateSessionId();

  const payload = {
    id: sessionId,
    country: "",
    city: "",
    referrer: document.referrer || "",
    source: getSource(),
    device: getDeviceType(),
    os: getOS(),
    browser: getBrowser(),
    theme: getTheme(),
  };

  try {
    const resp = await fetch(`${ANALYTICS_BASE}/session`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });
    if (!resp.ok) {
      const text = await resp.text().catch(() => "");
      console.warn("analytics: session POST failed", resp.status, text);
      sessionId = null;
    }
  } catch (e) {
    console.warn("analytics: session POST error", e);
    sessionId = null;
  }
}

export function trackEvent(type: string, target: string, value?: string): void {
  if (!BASE || !sessionId) return;

  const payload = [{
    session_id: sessionId,
    type,
    target,
    value: value || "",
    ts: Date.now(),
  }];

  fetch(`${ANALYTICS_BASE}/e`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  }).then(r => {
    if (!r.ok) r.text().then(t => console.warn("analytics: events POST failed", r.status, t));
  }).catch(() => {});
}
