import { readFileSync, readdirSync, writeFileSync, existsSync, mkdirSync } from 'fs';
import { join } from 'path';

const ROOT = process.cwd();
const BLOG_DIR = join(ROOT, 'content', 'blog');
const CERTS_DIR = join(ROOT, 'content', 'certs');
const DATA_DIR = join(ROOT, 'src', 'data');

function parseFrontmatter(content) {
  const match = content.match(/^---\s*\n([\s\S]*?)\n---\s*\n([\s\S]*)$/);
  if (!match) return { frontmatter: {}, body: content };

  const fmText = match[1];
  const body = match[2];
  const frontmatter = {};

  for (const line of fmText.split('\n')) {
    const colonIdx = line.indexOf(':');
    if (colonIdx === -1) continue;
    const key = line.slice(0, colonIdx).trim();
    let value = line.slice(colonIdx + 1).trim();

    if ((value.startsWith('"') && value.endsWith('"')) || (value.startsWith("'") && value.endsWith("'"))) {
      value = value.slice(1, -1);
    }

    frontmatter[key] = value;
  }

  if (frontmatter.tags) {
    const t = frontmatter.tags.trim();
    if (t.startsWith('[') && t.endsWith(']')) {
      frontmatter.tags = JSON.parse(t);
    } else {
      frontmatter.tags = t.split(',').map(tag => tag.trim()).filter(Boolean).join(',');
    }
  }

  return { frontmatter, body };
}

function buildBlog() {
  const posts = [];

  if (!existsSync(BLOG_DIR)) {
    writeFileSync(join(DATA_DIR, 'blog.ts'), generateBlogTs(posts));
    return posts;
  }

  const files = readdirSync(BLOG_DIR).filter(f => f.endsWith('.md')).sort();
  for (const file of files) {
    const content = readFileSync(join(BLOG_DIR, file), 'utf-8');
    const { frontmatter, body } = parseFrontmatter(content);
    const date = frontmatter.date || '';

    posts.push({
      id: frontmatter.slug,
      slug: frontmatter.slug,
      title: frontmatter.title || '',
      content_md: body.trim(),
      excerpt: frontmatter.summary || '',
      cover_image: frontmatter.cover || '',
      tags: Array.isArray(frontmatter.tags) ? frontmatter.tags : (frontmatter.tags ? frontmatter.tags.split(',').map(t => t.trim()).filter(Boolean) : []),
      published: frontmatter.published !== 'false',
      published_at: date,
      created_at: date,
      updated_at: date,
    });
  }

  posts.sort((a, b) => new Date(b.published_at || 0) - new Date(a.published_at || 0));
  writeFileSync(join(DATA_DIR, 'blog.ts'), generateBlogTs(posts));
  return posts;
}

function generateBlogTs(posts) {
  const items = posts.map(p => `  {
    id: ${JSON.stringify(p.id)},
    slug: ${JSON.stringify(p.slug)},
    title: ${JSON.stringify(p.title)},
    content_md: ${JSON.stringify(p.content_md)},
    excerpt: ${JSON.stringify(p.excerpt)},
    cover_image: ${JSON.stringify(p.cover_image)},
    tags: ${JSON.stringify(p.tags)},
    published: ${p.published},
    published_at: ${JSON.stringify(p.published_at)},
    created_at: ${JSON.stringify(p.created_at)},
    updated_at: ${JSON.stringify(p.updated_at)},
  }`).join(',\n');

  return `import type { BlogPost } from '$lib/types';

export const blogData: BlogPost[] = [
${items}
];

export default blogData;
`;
}

const CERT_CATEGORY_LABELS = {
  '01_Internships': 'Internships',
  '02_Certifications': 'Certifications',
  '03_Trainings_and_Courses': 'Trainings & Courses',
  '04_Events_and_Participation': 'Events & Participation',
};

function buildCerts() {
  const certs = [];

  if (!existsSync(CERTS_DIR)) {
    writeFileSync(join(DATA_DIR, 'certs.ts'), generateCertsTs(certs));
    return certs;
  }

  // Certificates mirror the documents folder structure:
  // content/certs/<01..NN_category>/<id>.json — folder prefix = showcase rank.
  const categories = readdirSync(CERTS_DIR, { withFileTypes: true })
    .filter((d) => d.isDirectory())
    .map((d) => d.name)
    .sort();

  let id = 0;
  for (const category of categories) {
    const dir = join(CERTS_DIR, category);
    const files = readdirSync(dir).filter((f) => f.endsWith('.json')).sort();

    for (const f of files) {
      const cert = JSON.parse(readFileSync(join(dir, f), 'utf-8'));
      id += 1;

      const file = cert.file || '';
      const dotIdx = file.lastIndexOf('.');
      const ext = dotIdx === -1 ? '' : file.slice(dotIdx + 1).toLowerCase();
      const kind = ext === 'pdf' ? 'pdf' : 'image';

      const fileUrl = file ? `certs/${category}/${file}` : '';
      let imageUrl = fileUrl;
      if (kind === 'pdf' && fileUrl) {
        // PDFs can't render as <img> — use the pre-generated first-page thumbnail.
        imageUrl = `certs/thumbs/${file.slice(0, -ext.length - 1)}.jpg`;
      }

      const date = cert.date || '';
      certs.push({
        id,
        title: cert.title || '',
        issuer: cert.issuer || '',
        date,
        category,
        category_label: CERT_CATEGORY_LABELS[category] || category,
        category_rank: parseInt(category.slice(0, 2), 10) || 99,
        kind,
        file_url: fileUrl,
        note: cert.note || '',
        credential_url: cert.credential_url || '',
        image_url: imageUrl,
        tags: Array.isArray(cert.tags) ? cert.tags : [],
        is_verified: !!cert.is_verified,
        display_order: id,
        created_at: date ? `${date}-01T00:00:00Z` : '2026-01-01T00:00:00Z',
      });
    }
  }

  certs.sort((a, b) => a.category_rank - b.category_rank || a.display_order - b.display_order);
  writeFileSync(join(DATA_DIR, 'certs.ts'), generateCertsTs(certs));
  return certs;
}

function generateCertsTs(certs) {
  const items = certs.map(c => `  {
    id: ${c.id},
    title: ${JSON.stringify(c.title)},
    issuer: ${JSON.stringify(c.issuer)},
    date: ${JSON.stringify(c.date)},
    category: ${JSON.stringify(c.category)},
    category_label: ${JSON.stringify(c.category_label)},
    category_rank: ${c.category_rank},
    kind: ${JSON.stringify(c.kind)},
    file_url: ${JSON.stringify(c.file_url)},
    note: ${JSON.stringify(c.note)},
    credential_url: ${JSON.stringify(c.credential_url)},
    image_url: ${JSON.stringify(c.image_url)},
    tags: ${JSON.stringify(c.tags)},
    is_verified: ${c.is_verified},
    display_order: ${c.display_order},
    created_at: ${JSON.stringify(c.created_at)},
  }`).join(',\n');

  return `import type { Certificate } from '$lib/types';

export const certsData: Certificate[] = [
${items}
];

export default certsData;
`;
}

mkdirSync(DATA_DIR, { recursive: true });
const posts = buildBlog();
const certs = buildCerts();
console.log(`Build complete: ${posts.length} blog post(s), ${certs.length} certificate(s)`);
