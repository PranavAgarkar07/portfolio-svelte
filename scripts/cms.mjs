#!/usr/bin/env node
/**
 * Local offline CMS server — zero dependencies (node built-ins only).
 *
 *   npm run cms   →  http://localhost:4323
 *
 * Drag-and-drop images/documents; files land in content/ and static/
 * with repo-convention names, and `build:content` regenerates src/data
 * automatically after every save. Binds to 127.0.0.1 only.
 */
import http from 'http';
import fs from 'fs';
import path from 'path';
import { execFileSync, spawnSync } from 'child_process';
import { fileURLToPath } from 'url';

const ROOT = path.join(path.dirname(fileURLToPath(import.meta.url)), '..');
const PORT = 4323;
const HOST = '127.0.0.1';
const MAX_UPLOAD = 25 * 1024 * 1024;

const BLOG_DIR = path.join(ROOT, 'content', 'blog');
const BLOG_IMG_DIR = path.join(ROOT, 'static', 'images', 'blog');
const CERTS_DIR = path.join(ROOT, 'content', 'certs');
const CERT_STATIC = path.join(ROOT, 'static', 'certs');
const UI = path.join(ROOT, 'scripts', 'cms-ui.html');

const IMG_EXT = ['jpg', 'jpeg', 'png', 'webp', 'gif', 'svg', 'avif'];
const DOC_EXT = [...IMG_EXT, 'pdf'];

/* ─────────────── helpers ─────────────── */

export function slugify(s) {
  return String(s || '')
    .toLowerCase()
    .replace(/['"`,]/g, '')
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '')
    .slice(0, 80);
}

function unique(dir, name) {
  let p = path.join(dir, name);
  if (!fs.existsSync(p)) return name;
  const ext = path.extname(name);
  const base = name.slice(0, -ext.length);
  for (let i = 2; ; i++) {
    if (!fs.existsSync(path.join(dir, `${base}-${i}${ext}`))) return `${base}-${i}${ext}`;
  }
}

function parseFrontmatter(content) {
  const m = content.match(/^---\s*\n([\s\S]*?)\n---\s*\n?([\s\S]*)$/);
  if (!m) return { frontmatter: {}, body: content };
  const fm = {};
  for (const line of m[1].split('\n')) {
    const i = line.indexOf(':');
    if (i === -1) continue;
    const key = line.slice(0, i).trim();
    let val = line.slice(i + 1).trim();
    if ((val.startsWith('"') && val.endsWith('"')) || (val.startsWith("'") && val.endsWith("'")))
      val = val.slice(1, -1);
    fm[key] = val;
  }
  return { frontmatter: fm, body: m[2].replace(/\s+$/, '') };
}

function quote(s) {
  return JSON.stringify(String(s ?? ''));
}

function blogToMd(fm, body) {
  return (
    `---\n` +
    `title: ${quote(fm.title)}\n` +
    `slug: ${fm.slug}\n` +
    `date: ${fm.date}\n` +
    `summary: ${quote(fm.summary)}\n` +
    `tags: ${JSON.stringify(fm.tags)}\n` +
    `published: ${fm.published !== false}\n` +
    `cover: ${quote(fm.cover || '')}\n` +
    `---\n\n` +
    body.replace(/\r\n/g, '\n').trim() +
    `\n`
  );
}

function runContentBuild() {
  const r = spawnSync('node', ['scripts/build-content.js'], { cwd: ROOT, encoding: 'utf8' });
  return { ok: r.status === 0, out: [r.stdout, r.stderr].filter(Boolean).join('\n') };
}

/** PDF → static/certs/thumbs/<name>.jpg via pdftoppm, else magick. */
function makePdfThumb(pdfPath, thumbPath) {
  try {
    execFileSync('pdftoppm', ['-jpeg', '-f', '1', '-singlefile', '-r', '110', pdfPath, thumbPath.replace(/\.jpg$/, '')], { cwd: ROOT });
    return null;
  } catch (e1) {
    try {
      execFileSync('magick', [`${pdfPath}[0]`, '-resize', '1100x1100', '-quality', '82', 'jpg:' + thumbPath], { cwd: ROOT });
      return null;
    } catch (e2) {
      return `thumbnail failed: ${e1.message} / ${e2.message}`;
    }
  }
}

/* ─────────────── multipart (no deps) ─────────────── */

function parseMultipart(buf, contentType) {
  const m = contentType.match(/boundary=(?:"([^"]+)"|([^;]+))/);
  if (!m) return null;
  const boundary = Buffer.from(`--${m[1] || m[2]}`);
  const parts = {};
  let idx = buf.indexOf(boundary);
  while (idx !== -1) {
    const start = idx + boundary.length;
    if (buf.slice(start, start + 2).toString() === '--') break;
    const partStart = buf.indexOf(Buffer.from('\r\n\r\n'), start);
    if (partStart === -1) break;
    const headers = buf.slice(start, partStart).toString('utf8');
    const next = buf.indexOf(boundary, partStart);
    // strip trailing \r\n before boundary
    const data = buf.slice(partStart + 4, next - 2);
    const dispLine = headers.match(/content-disposition:[^\n]*/i)?.[0] || '';
    const fields = {};
    for (const tok of dispLine.split(';')) {
      const kv = tok.trim().match(/^(name|filename)="([^"]*)"/i);
      if (kv) fields[kv[1].toLowerCase()] = kv[2];
    }
    if (fields.name) parts[fields.name] = { data, filename: fields.filename || '' };
    idx = next;
  }
  return parts;
}

function readBody(req) {
  return new Promise((resolve, reject) => {
    const chunks = [];
    let size = 0;
    req.on('data', (c) => {
      size += c.length;
      if (size > MAX_UPLOAD) { reject(new Error('upload exceeds 25MB')); req.destroy(); return; }
      chunks.push(c);
    });
    req.on('end', () => resolve(Buffer.concat(chunks)));
    req.on('error', reject);
  });
}

/* ─────────────── API handlers ─────────────── */

function listBlog() {
  const files = fs.existsSync(BLOG_DIR) ? fs.readdirSync(BLOG_DIR).filter((f) => f.endsWith('.md')) : [];
  return files
    .map((f) => {
      const { frontmatter } = parseFrontmatter(fs.readFileSync(path.join(BLOG_DIR, f), 'utf8'));
      return {
        file: f,
        title: frontmatter.title || '',
        date: frontmatter.date || '',
        cover: frontmatter.cover || '',
        tags: (frontmatter.tags || '').toString().replace(/^\[/, '').replace(/\]$/, '').replace(/"/g, ''),
        published: frontmatter.published !== 'false',
      };
    })
    .sort((a, b) => (b.date || '').localeCompare(a.date || ''));
}

function listCerts() {
  const out = {};
  const cats = fs.existsSync(CERTS_DIR)
    ? fs.readdirSync(CERTS_DIR, { withFileTypes: true }).filter((d) => d.isDirectory()).map((d) => d.name).sort()
    : [];
  for (const cat of cats) {
    out[cat] = fs
      .readdirSync(path.join(CERTS_DIR, cat))
      .filter((f) => f.endsWith('.json'))
      .map((f) => JSON.parse(fs.readFileSync(path.join(CERTS_DIR, cat, f), 'utf8')));
  }
  return out;
}

function apiBlogFile(query) {
  const file = path.basename(String(query.get('file') || ''));
  if (!file.endsWith('.md')) throw new Error('bad file');
  const full = path.join(BLOG_DIR, file);
  if (!fs.existsSync(full)) throw new Error('not found');
  const { frontmatter, body } = parseFrontmatter(fs.readFileSync(full, 'utf8'));
  let tags = frontmatter.tags || '';
  if (tags.startsWith('[')) tags = JSON.parse(tags);
  else tags = tags ? tags.split(',').map((t) => t.trim()).filter(Boolean) : [];
  return {
    slug: frontmatter.slug || file.replace(/\.md$/, ''),
    title: frontmatter.title || '',
    date: frontmatter.date || '',
    summary: frontmatter.summary || '',
    cover: frontmatter.cover || '',
    published: frontmatter.published !== 'false',
    tags,
    body,
  };
}

async function apiImage(body, ct) {
  const parts = parseMultipart(body, ct);
  if (!parts?.file) throw new Error('missing file');
  const mode = parts.mode?.data.toString() === 'cover' ? 'cover' : 'inline';
  const slug = slugify(parts.slug?.data.toString());
  if (!slug) throw new Error('missing slug');
  const dest = parts.dest?.data.toString() === 'cert' ? 'cert' : 'blog';

  let dir, urlDir, nameBase;
  if (dest === 'blog') {
    dir = BLOG_IMG_DIR; urlDir = '/images/blog';
    nameBase = mode === 'cover' ? slug : `${slug}-img`;
  } else {
    throw new Error('use cert-doc endpoint for certificate documents');
  }
  fs.mkdirSync(dir, { recursive: true });

  const raw = parts.file.filename || '';
  const ext = path.extname(raw).slice(1).toLowerCase();
  if (!IMG_EXT.includes(ext)) throw new Error(`extension .${ext || '?'} not in ${IMG_EXT.join('/')}`);

  let name;
  if (mode === 'cover') {
    name = `${slug}.${ext}`;
  } else {
    // numbered inline: <slug>-img1.jpg, -img2.jpg … even across drops
    let n = 1;
    while (fs.existsSync(path.join(dir, `${slug}-img${n}.${ext}`)) ||
           fs.readdirSync(dir).some((f) => f.startsWith(`${slug}-img${n}.`))) n++;
    name = `${slug}-img${n}.${ext}`;
  }
  name = unique(dir, name);
  fs.writeFileSync(path.join(dir, name), parts.file.data);
  return { url: `${urlDir}/${name}`, name, mode };
}

async function apiCertDoc(body, ct) {
  const parts = parseMultipart(body, ct);
  if (!parts?.file) throw new Error('missing file');
  const category = String(parts.category?.data || '');
  const cats = fs.readdirSync(CERTS_DIR).filter((c) => fs.statSync(path.join(CERTS_DIR, c)).isDirectory());
  if (!cats.includes(category)) throw new Error(`unknown category (expect ${cats.join(', ')})`);
  const title = parts.title?.data.toString() || '';
  const issuer = parts.issuer?.data.toString() || '';
  const date = parts.date?.data.toString() || '';

  const ext = path.extname(parts.file.filename || '').slice(1).toLowerCase();
  if (!DOC_EXT.includes(ext)) throw new Error(`extension .${ext || '?'} not in ${DOC_EXT.join('/')}`);

  const dir = path.join(CERT_STATIC, category);
  fs.mkdirSync(dir, { recursive: true });
  let name = `${slugify(title)}${issuer ? `_${slugify(issuer)}` : ''}${date ? `_${date}` : ''}`;
  name = unique(dir, `${name}.${ext}`);
  fs.writeFileSync(path.join(dir, name), parts.file.data);

  let thumbError = null;
  if (ext === 'pdf') {
    const thumbDir = path.join(CERT_STATIC, 'thumbs');
    fs.mkdirSync(thumbDir, { recursive: true });
    const thumbPath = path.join(thumbDir, `${name.slice(0, -ext.length - 1)}.jpg`);
    if (!fs.existsSync(thumbPath)) thumbError = makePdfThumb(path.join(dir, name), thumbPath);
  }
  return { file: `certs/${category}/${name}`, category, name, thumbError };
}

async function routeApi(req, res, url) {
  const send = (code, obj) => {
    const body = JSON.stringify(obj);
    res.writeHead(code, { 'Content-Type': 'application/json' });
    res.end(body);
  };
  try {
    if (req.method === 'GET') {
      if (url.pathname === '/api/blog') return send(200, listBlog());
      if (url.pathname === '/api/certs') return send(200, listCerts());
      if (url.pathname === '/api/blog-file') return send(200, apiBlogFile(url.searchParams));
      return send(404, { error: 'no such endpoint' });
    }
    // POST endpoints
    const body = await readBody(req);
    if (url.pathname === '/api/image') {
      const ct = req.headers['content-type'] || '';
      if (!ct.includes('multipart')) throw new Error('multipart required');
      const r = await apiImage(body, ct);
      const build = runContentBuild();
      return send(200, { ...r, build });
    }
    if (url.pathname === '/api/cert-doc') {
      const r = await apiCertDoc(body, req.headers['content-type'] || '');
      return send(200, r);
    }
    if (url.pathname === '/api/blog-save') {
      const fm = JSON.parse(body.toString('utf8'));
      fm.slug = slugify(fm.slug);
      if (!fm.slug) throw new Error('slug is required');
      if (!fm.date) throw new Error('date is required (YYYY-MM-DD)');
      fs.mkdirSync(BLOG_DIR, { recursive: true });
      fs.writeFileSync(path.join(BLOG_DIR, `${fm.slug}.md`), blogToMd(fm, fm.body || ''));
      const build = runContentBuild();
      return send(200, { ok: true, file: `content/blog/${fm.slug}.md`, build });
    }
    if (url.pathname === '/api/cert-save') {
      const j = JSON.parse(body.toString('utf8'));
      const category = String(j.category || '');
      const cats = fs.readdirSync(CERTS_DIR).filter((c) => fs.statSync(path.join(CERTS_DIR, c)).isDirectory());
      if (!cats.includes(category)) throw new Error(`unknown category`);
      if (!j.file) throw new Error('drop the certificate document first');
      const doc = {
        title: j.title || '',
        issuer: j.issuer || '',
        date: j.date || '',
        note: j.note || '',
        credential_url: j.credential_url || '',
        tags: j.tags || [],
        is_verified: !!j.is_verified,
        file: String(j.file).split('/').pop(),
      };
      fs.writeFileSync(path.join(CERTS_DIR, category, `${slugify(j.title)}.json`), JSON.stringify(doc, null, 2) + '\n');
      const build = runContentBuild();
      return send(200, { ok: true, file: `${slugify(j.title)}.json`, build });
    }
    if (url.pathname === '/api/export-check') {
      // final sanity: run svelte-check quietly? keep light – content build only
      return send(200, runContentBuild());
    }
    return send(404, { error: 'no such endpoint' });
  } catch (err) {
    return send(400, { error: err.message || String(err) });
  }
}

/* ─────────────── server ─────────────── */

const MIME = {
  '.jpg': 'image/jpeg', '.jpeg': 'image/jpeg', '.png': 'image/png', '.webp': 'image/webp',
  '.gif': 'image/gif', '.svg': 'image/svg+xml', '.avif': 'image/avif', '.pdf': 'application/pdf',
};

function serveStatic(url, res) {
  if (!url.pathname.startsWith('/static/')) return false;
  const rel = decodeURIComponent(url.pathname.slice('/static/'.length));
  const full = path.join(ROOT, 'static', path.normalize(rel).replace(/^(\.\.[/\\])+/, ''));
  if (!full.startsWith(path.join(ROOT, 'static') + path.sep) || !fs.existsSync(full) || !fs.statSync(full).isFile()) {
    res.writeHead(404); res.end('no such file'); return true;
  }
  res.writeHead(200, { 'Content-Type': MIME[path.extname(full).toLowerCase()] || 'application/octet-stream' });
  res.end(fs.readFileSync(full));
  return true;
}

const server = http.createServer(async (req, res) => {
  const url = new URL(req.url, `http://${req.headers.host || 'localhost'}`);
  if (url.pathname.startsWith('/api/')) return routeApi(req, res, url);
  if (serveStatic(url, res)) return;
  if (url.pathname === '/' || url.pathname === '/index.html') {
    const html = fs.readFileSync(UI);
    res.writeHead(200, { 'Content-Type': 'text/html; charset=utf-8', 'Cache-Control': 'no-store' });
    return res.end(html);
  }
  res.writeHead(404, { 'Content-Type': 'text/plain' });
  res.end('not found');
});

server.listen(PORT, HOST, () => {
  console.log(`offline CMS  →  http://${HOST}:${PORT}`);
  console.log('drag image/posts anywhere — files land in content/ + static/ + auto build:content');
});
