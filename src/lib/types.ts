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

export interface User {
  id: string;
  email: string;
  name: string;
  avatar_url: string;
  role: 'user' | 'author' | 'admin';
  created_at: string;
}

export interface BlogImage {
  url: string;
  alt?: string;
  caption?: string;
}

export interface BlogPost {
  id: string;
  slug: string;
  title: string;
  content_md: string;
  excerpt: string;
  cover_image: string;
  images?: BlogImage[];
  tags: string[];
  published: boolean;
  published_at?: string;
  author_id?: string;
  author_name?: string;
  created_at: string;
  updated_at: string;
}

export interface ProjectReview {
  id?: number;
  user_id: string;
  user_name: string;
  avatar_url?: string;
  project_name: string;
  rating: number;
  comment?: string;
  created_at: string;
}

export interface BlogComment {
  id: string;
  post_id: string;
  parent_id?: string;
  user_id: string;
  user_name: string;
  avatar_url: string;
  content: string;
  created_at: string;
  updated_at: string;
  replies: BlogComment[];
}

export interface MarqueeItem {
  type: string;
  user_name: string;
  avatar_url?: string;
  project_name: string;
  rating: number;
  comment?: string;
  created_at: string;
}
