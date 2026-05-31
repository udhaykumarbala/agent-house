// Minimal, safe Markdown → HTML for Brain replies.
// HTML is escaped first, so output is XSS-safe even on untrusted content.
// Supports: headings, hr, blockquote, GFM tables, ordered/unordered lists,
// paragraphs, soft line breaks, **bold**, *italic*, `code`.

function esc(s: string): string {
  return s
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;");
}

function inline(s: string): string {
  let t = esc(s);
  // RFC-style `Name <email@host>` — render the address as a clean token
  // instead of bare escaped angle brackets (which read as stray markup).
  t = t.replace(
    /&lt;([^\s<>&]+@[^\s<>&]+\.[a-z]{2,})&gt;/gi,
    '<span class="eml">$1</span>'
  );
  // `code` first so ** inside code is left alone
  t = t.replace(/`([^`]+)`/g, "<code>$1</code>");
  t = t.replace(/\*\*([^*]+)\*\*/g, "<strong>$1</strong>");
  t = t.replace(/(^|[^*])\*([^*\n]+)\*/g, "$1<em>$2</em>");
  return t;
}

// "| a | b |" -> ["a", "b"] (tolerates missing edge pipes)
function cells(row: string): string[] {
  return row
    .trim()
    .replace(/^\|/, "")
    .replace(/\|$/, "")
    .split("|")
    .map((c) => c.trim());
}

const isTableSep = (l: string): boolean =>
  /^\|?\s*:?-{1,}:?\s*(\|\s*:?-{1,}:?\s*)*\|?$/.test(l.trim()) &&
  l.includes("-") &&
  l.includes("|"); // a real separator has column pipes — excludes "---" hr

const looksLikeRow = (l: string): boolean => l.trim().includes("|");

export function renderMarkdown(src: string): string {
  if (!src) return "";
  const lines = src.replace(/\r\n/g, "\n").split("\n");
  const out: string[] = [];
  let para: string[] = [];
  let list: { ordered: boolean; items: string[] } | null = null;

  const flushPara = () => {
    if (para.length) {
      out.push(`<p>${para.map(inline).join("<br>")}</p>`);
      para = [];
    }
  };
  const flushList = () => {
    if (list) {
      const tag = list.ordered ? "ol" : "ul";
      out.push(
        `<${tag}>${list.items.map((i) => `<li>${inline(i)}</li>`).join("")}</${tag}>`
      );
      list = null;
    }
  };
  const flush = () => {
    flushPara();
    flushList();
  };

  for (let i = 0; i < lines.length; i++) {
    const raw = lines[i];
    const line = raw.trim();

    if (line === "") {
      flush();
      continue;
    }

    // GFM table: a row line followed by a |---|---| separator
    if (
      looksLikeRow(line) &&
      i + 1 < lines.length &&
      isTableSep(lines[i + 1])
    ) {
      flush();
      const header = cells(line);
      i += 2; // skip header + separator
      const body: string[][] = [];
      while (i < lines.length && lines[i].trim() !== "" && looksLikeRow(lines[i])) {
        body.push(cells(lines[i]));
        i++;
      }
      i--; // for-loop will ++ back
      out.push(
        `<table><thead><tr>${header
          .map((h) => `<th>${inline(h)}</th>`)
          .join("")}</tr></thead><tbody>${body
          .map(
            (r) =>
              `<tr>${r.map((c) => `<td>${inline(c)}</td>`).join("")}</tr>`
          )
          .join("")}</tbody></table>`
      );
      continue;
    }

    if (/^(-{3,}|\*{3,}|_{3,})$/.test(line)) {
      flush();
      out.push("<hr>");
      continue;
    }
    const h = /^(#{1,3})\s+(.*)$/.exec(line);
    if (h) {
      flush();
      const lvl = h[1].length + 2; // # -> h3, ## -> h4, ### -> h5
      out.push(`<h${lvl}>${inline(h[2])}</h${lvl}>`);
      continue;
    }
    if (/^>\s+/.test(line)) {
      flushPara();
      flushList();
      out.push(`<blockquote>${inline(line.replace(/^>\s+/, ""))}</blockquote>`);
      continue;
    }
    const ol = /^(\d+)\.\s+(.*)$/.exec(line);
    const ul = /^[-*]\s+(.*)$/.exec(line);
    if (ol || ul) {
      flushPara();
      const ordered = !!ol;
      if (!list || list.ordered !== ordered) {
        flushList();
        list = { ordered, items: [] };
      }
      list.items.push((ol ? ol[2] : ul![1]).trim());
      continue;
    }
    // plain text line — part of current paragraph
    flushList();
    para.push(line);
  }
  flush();
  return out.join("");
}
