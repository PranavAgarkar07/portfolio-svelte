export interface Certificate {
  id: number;
  title: string;
  issuer: string;
  date: string;
  /** Source folder, e.g. "01_Internships" (mirrors the documents directory). */
  category: string;
  /** Human-readable category name. */
  category_label: string;
  /** Folder-order rank (01 → 1) used for the default showcase sorting. */
  category_rank: number;
  /** Document type: "pdf" opens in a new tab; "image" opens in the lightbox. */
  kind: "pdf" | "image";
  /** Static path of the actual document (committed to static/certs/). */
  file_url: string;
  /** Small evidence line, e.g. "AICTE ID 2335036 · 6 weeks". */
  note: string;
  credential_url: string;
  /** Display thumbnail — generated first-page JPG for PDFs, the file itself for images. */
  image_url: string;
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
  replies?: BlogComment[] | null;
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
