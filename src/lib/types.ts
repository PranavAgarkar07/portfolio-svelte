export interface DevLogResponse {
  summary: string;
  last_update: string;
  source: "cache" | "live" | "stale-cache";
}

export interface Certificate {
  id: number;
  title: string;
  issuer: string;
  date: string;
  credential_url: string;
  image_url: string;
  thumb_url?: string;
  tags: string[];
  is_verified: boolean;
  display_order: number;
  created_at: string;
}

export interface Badge {
  id: number;
  name: string;
  image_url: string;
  credential_url: string;
  rarity: "common" | "uncommon" | "rare";
  category: string;
  important: boolean;
  display_order: number;
  created_at: string;
}

export interface SkillItem {
  name: string;
  icon: string;
  level: string;
  projects?: string[];
}
