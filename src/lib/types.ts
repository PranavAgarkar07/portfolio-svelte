export interface DevLogResponse {
  summary: string;
  last_update: string;
  source: "cache" | "live" | "stale-cache";
}
