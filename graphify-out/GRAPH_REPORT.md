# Graph Report - .  (2026-05-11)

## Corpus Check
- 10 files · ~106,174 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 29 nodes · 24 edges · 12 communities detected
- Extraction: 100% EXTRACTED · 0% INFERRED · 0% AMBIGUOUS
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_Community 0|Community 0]]
- [[_COMMUNITY_Community 1|Community 1]]
- [[_COMMUNITY_Community 2|Community 2]]
- [[_COMMUNITY_Community 3|Community 3]]
- [[_COMMUNITY_Community 4|Community 4]]
- [[_COMMUNITY_Community 5|Community 5]]
- [[_COMMUNITY_Community 6|Community 6]]
- [[_COMMUNITY_Community 7|Community 7]]
- [[_COMMUNITY_Community 8|Community 8]]
- [[_COMMUNITY_Community 9|Community 9]]
- [[_COMMUNITY_Community 10|Community 10]]
- [[_COMMUNITY_Community 11|Community 11]]

## God Nodes (most connected - your core abstractions)
1. `generateDevLog()` - 4 edges
2. `getInitialTheme()` - 2 edges
3. `createThemeStore()` - 2 edges
4. `handleStatus()` - 2 edges
5. `handleAdminContact()` - 2 edges
6. `htmlEscape()` - 2 edges
7. `fetchGitHubEvents()` - 2 edges
8. `callGemini()` - 2 edges
9. `DevLogResponse` - 1 edges
10. `CachedData` - 1 edges

## Surprising Connections (you probably didn't know these)
- None detected - all connections are within the same source files.

## Communities

### Community 0 - "Community 0"
Cohesion: 0.17
Nodes (5): CachedData, CircuitBreaker, ContactSubmission, DevLogResponse, Metrics

### Community 1 - "Community 1"
Cohesion: 0.5
Nodes (4): callGemini(), fetchGitHubEvents(), generateDevLog(), handleStatus()

### Community 2 - "Community 2"
Cohesion: 1.0
Nodes (2): createThemeStore(), getInitialTheme()

### Community 3 - "Community 3"
Cohesion: 1.0
Nodes (2): handleAdminContact(), htmlEscape()

### Community 4 - "Community 4"
Cohesion: 1.0
Nodes (0): 

### Community 5 - "Community 5"
Cohesion: 1.0
Nodes (0): 

### Community 6 - "Community 6"
Cohesion: 1.0
Nodes (0): 

### Community 7 - "Community 7"
Cohesion: 1.0
Nodes (0): 

### Community 8 - "Community 8"
Cohesion: 1.0
Nodes (0): 

### Community 9 - "Community 9"
Cohesion: 1.0
Nodes (0): 

### Community 10 - "Community 10"
Cohesion: 1.0
Nodes (0): 

### Community 11 - "Community 11"
Cohesion: 1.0
Nodes (0): 

## Knowledge Gaps
- **5 isolated node(s):** `DevLogResponse`, `CachedData`, `CircuitBreaker`, `Metrics`, `ContactSubmission`
  These have ≤1 connection - possible missing edges or undocumented components.
- **Thin community `Community 3`** (2 nodes): `handleAdminContact()`, `htmlEscape()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 4`** (1 nodes): `svelte.config.js`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 5`** (1 nodes): `vite.config.ts`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 6`** (1 nodes): `app.d.ts`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 7`** (1 nodes): `data.ts`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 8`** (1 nodes): `index.ts`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 9`** (1 nodes): `types.ts`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 10`** (1 nodes): `index.ts`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 11`** (1 nodes): `+layout.ts`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `generateDevLog()` connect `Community 1` to `Community 0`?**
  _High betweenness centrality (0.004) - this node is a cross-community bridge._
- **What connects `DevLogResponse`, `CachedData`, `CircuitBreaker` to the rest of the system?**
  _5 weakly-connected nodes found - possible documentation gaps or missing edges._