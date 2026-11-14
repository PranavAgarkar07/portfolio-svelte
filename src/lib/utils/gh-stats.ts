/**
 * Live GitHub engagement stats (stars + forks) for the project cards.
 * Replaces the removed Go-API like counter: the value shown is synced
 * from the GitHub REST API, cached per browser session for one hour so
 * a visit stays well inside the unauthenticated rate budget (60 req/h).
 */

export interface RepoStats {
  stars: number;
  forks: number;
}

/** Extract "owner/repo" from a github.com URL (tolerates .git and extra path segments). */
export function repoFullNameFromUrl(url: string): string | null {
  const m = /(?:^|\/\/)github\.com\/([^/\s]+)\/([^/\s?#]+)/.exec(url);
  if (!m) return null;
  return `${m[1]}/${m[2].replace(/\.git$/, "")}`;
}

const CACHE_KEY = "gh-repo-stats-v1";
const CACHE_TTL = 60 * 60 * 1000; // 1 hour

interface CacheEntry {
  stats: RepoStats;
  cachedAt: number;
}

function readCache(): Record<string, CacheEntry> {
  if (typeof sessionStorage === "undefined") return {};
  try {
    return JSON.parse(sessionStorage.getItem(CACHE_KEY) || "{}") as Record<
      string,
      CacheEntry
    >;
  } catch {
    return {};
  }
}

function writeCache(cache: Record<string, CacheEntry>): void {
  try {
    sessionStorage.setItem(CACHE_KEY, JSON.stringify(cache));
  } catch {
    /* Quota or blocked storage — stats are cosmetic, drop them. */
  }
}

async function fetchOne(repo: string): Promise<RepoStats | null> {
  try {
    const res = await fetch(`https://api.github.com/repos/${repo}`, {
      headers: { Accept: "application/vnd.github+json" },
    });
    if (!res.ok) return null;
    const data = await res.json();
    if (typeof data.stargazers_count !== "number") return null;
    return {
      stars: data.stargazers_count,
      forks: typeof data.forks_count === "number" ? data.forks_count : 0,
    };
  } catch {
    /* Offline or rate limited — resolve null, caller keeps the fallback. */
    return null;
  }
}

let inflight: Promise<Record<string, RepoStats | null>> | null = null;

/**
 * Fetch stats for a list of "owner/repo" names. Individual values may be
 * null when a repo could not be loaded (offline, rate limited, renamed).
 */
export function fetchRepoStats(
  repos: string[]
): Promise<Record<string, RepoStats | null>> {
  const cache = readCache();
  const now = Date.now();
  const resolved: Record<string, RepoStats | null> = {};
  const missing: string[] = [];
  for (const repo of new Set(repos)) {
    const entry = cache[repo];
    if (entry && now - entry.cachedAt < CACHE_TTL) {
      resolved[repo] = entry.stats;
    } else {
      resolved[repo] = null;
      missing.push(repo);
    }
  }
  if (missing.length === 0) return Promise.resolve(resolved);

  // Shared promise — concurrent callers must not duplicate requests.
  if (!inflight) {
    inflight = (async () => {
      const fresh: Record<string, RepoStats | null> = {};
      await Promise.all(
        missing.map(async (repo) => {
          const stats = await fetchOne(repo);
          fresh[repo] = stats;
          if (stats) {
            const updated = readCache();
            updated[repo] = { stats, cachedAt: Date.now() };
            writeCache(updated);
          }
        })
      );
      return fresh;
    })().finally(() => {
      inflight = null;
    });
  }
  return inflight.then((fresh) => ({ ...resolved, ...fresh }));
}
