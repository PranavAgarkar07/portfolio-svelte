export function renderMarkdown(md: string): string {
  let html = md.replace(/[&<>]/g, escapeHtml);

  html = html.replace(/^```(\w*)\n([\s\S]*?)```$/gm, (_, lang, code) => {
    const langAttr = lang ? ` class="lang-${lang}"` : '';
    return `<pre><code${langAttr}>${code.trim()}</code></pre>`;
  });

  html = html.replace(/^> (.+)$/gm, '<blockquote><p>$1</p></blockquote>');

  html = html.replace(/^### (.+)$/gm, '<h3>$1</h3>');
  html = html.replace(/^## (.+)$/gm, '<h2>$1</h2>');
  html = html.replace(/^# (.+)$/gm, '<h1>$1</h1>');

  html = html.replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>');
  html = html.replace(/\*(.+?)\*/g, '<em>$1</em>');
  html = html.replace(/`([^`]+)`/g, '<code>$1</code>');
  html = html.replace(/~~(.+?)~~/g, '<del>$1</del>');

  html = html.replace(/!\[([^\]]*)\]\(([^)]+)\)/g, '<img src="$2" alt="$1" loading="lazy" />');
  html = html.replace(/\[([^\]]+)\]\(([^)]+)\)/g, '<a href="$2" target="_blank" rel="noopener">$1</a>');

  html = html.replace(/^---$/gm, '<hr />');

  html = html.replace(/^- (.+)$/gm, '<li>$1</li>');
  html = html.replace(/(<li>.*<\/li>\n?)+/g, '<ul>$&</ul>');

  const tableRegex = /^\|(.+)\|\n\|[-| :]+\|\n((?:\|.+\|\n?)*)/gm;
  html = html.replace(tableRegex, (_, header, body) => {
    const headers = header.split('|').map((h: string) => h.trim()).filter(Boolean);
    const rows = body.trim().split('\n').map((row: string) =>
      row.split('|').map((c: string) => c.trim()).filter(Boolean)
    );
    let tbl = '<table><thead><tr>';
    for (const h of headers) tbl += `<th>${h}</th>`;
    tbl += '</tr></thead><tbody>';
    for (const row of rows) {
      if (row.length) {
        tbl += '<tr>';
        for (const c of row) tbl += `<td>${c}</td>`;
        tbl += '</tr>';
      }
    }
    return tbl + '</tbody></table>';
  });

  html = html.replace(/\n\n/g, '</p><p>');
  html = html.replace(/^(?!<[hublpta])/gm, '<p>');
  html = html.replace(/(<\/p>)\s*<p>(?=<[hublpta])/g, '$1');
  html = html.replace(/<p><\/p>/g, '');

  html = html.replace(/<blockquote>\s*<p>/g, '<blockquote><p>');
  html = html.replace(/<\/p>\s*<\/blockquote>/g, '</p></blockquote>');

  return '<p>' + html + '</p>';
}

function escapeHtml(s: string): string {
  return s
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;');
}

export function readingTime(md: string): string {
  const words = md.split(/\s+/).length;
  const min = Math.max(1, Math.round(words / 200));
  return `${min} min read`;
}

export function extractHeadings(md: string): Array<{ level: number; text: string; id: string }> {
  const headings: Array<{ level: number; text: string; id: string }> = [];
  const regex = /^(#{1,3})\s+(.+)$/gm;
  let match;
  while ((match = regex.exec(md)) !== null) {
    const level = match[1].length;
    const text = match[2].trim();
    const id = text.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-|-$/g, '');
    headings.push({ level, text, id });
  }
  return headings;
}
