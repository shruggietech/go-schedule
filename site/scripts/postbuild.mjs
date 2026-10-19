import fs from 'node:fs/promises';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import { legacyHTML } from './content.mjs';

const site = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..');
const out = path.join(site, 'out');
const routes = JSON.parse(await fs.readFile(path.join(site, '.source/routes.json'), 'utf8'));
for (const page of routes) {
  if (page.source === 'README.md') continue;
  await fs.writeFile(path.join(out, page.source.replace(/\.md$/, '.html')), legacyHTML(page.route), 'utf8');
}
await fs.writeFile(path.join(out, '.nojekyll'), '', 'utf8');
await fs.writeFile(path.join(out, 'robots.txt'), 'User-agent: *\nAllow: /\nSitemap: https://go-schedule.io/sitemap.xml\n');
const urls = ['/', '/compare/', ...routes.map(p => p.route)];
await fs.writeFile(path.join(out, 'sitemap.xml'), `<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">${urls.map(url => `<url><loc>https://go-schedule.io${url}</loc></url>`).join('')}</urlset>`);
const failures = [];
const documentCache = new Map();
async function document(file) {
  if (!documentCache.has(file)) documentCache.set(file, await fs.readFile(file, 'utf8'));
  return documentCache.get(file);
}
async function walk(dir) {
  for (const entry of await fs.readdir(dir, { withFileTypes: true })) {
    const file = path.join(dir, entry.name);
    if (entry.isDirectory()) { await walk(file); continue; }
    if (!entry.name.endsWith('.html')) continue;
    const html = await document(file);
    for (const match of html.matchAll(/(?:href|src)="([^"]+)"/g)) {
      const value = match[1].replaceAll('&amp;', '&');
      const url = new URL(value, `https://go-schedule.io/${path.relative(out, file).replaceAll(path.sep, '/')}`);
      if (url.origin !== 'https://go-schedule.io') continue;
      let target = path.join(out, decodeURIComponent(url.pathname));
      try {
        const stat = await fs.stat(target);
        if (stat.isDirectory()) target = path.join(target, 'index.html');
        await fs.access(target);
        if (url.hash && target.endsWith('.html')) {
          const fragment = decodeURIComponent(url.hash.slice(1));
          const ids = new Set([...((await document(target)).matchAll(/(?:id|name)="([^"]+)"/g))].map(m => m[1]));
          if (!ids.has(fragment)) failures.push(`${path.relative(out, file)}: missing fragment ${url.pathname}${url.hash}`);
        }
      } catch { failures.push(`${path.relative(out, file)}: missing exported target ${url.pathname}`); }
    }
  }
}
await walk(out);
async function auditAssets(dir, relative = '') {
  for (const entry of await fs.readdir(dir, {withFileTypes:true})) {
    const rel = path.join(relative, entry.name);
    if (entry.isDirectory()) { await auditAssets(path.join(dir,entry.name),rel); continue; }
    const original = await fs.readFile(path.join(dir,entry.name));
    const exported = await fs.readFile(path.join(out,'assets',rel));
    if (!original.equals(exported)) failures.push(`Changed exported canonical asset: ${rel}`);
  }
}
await auditAssets(path.resolve(site,'../docs/assets'));
if (failures.length) throw new Error(`Export link audit failed:\n${[...new Set(failures)].join('\n')}`);
console.log(`Exported ${routes.length} canonical docs, legacy aliases and public metadata; internal paths audited`);
