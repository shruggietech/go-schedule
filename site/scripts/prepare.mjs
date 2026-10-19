import fs from 'node:fs/promises';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import { parse as parseYAML } from 'yaml';
import { routeFor, rewriteMarkdown } from './content.mjs';

const site = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..');
const root = path.dirname(site);
const docsDir = path.join(root, 'docs');
const files = (await fs.readdir(docsDir)).filter(f => f.endsWith('.md')).sort();
const routes = new Map(files.map(f => [f, routeFor(f)]));
if (new Set(routes.values()).size !== files.length) throw new Error('Duplicate documentation route; rename the conflicting canonical source');
const generated = path.join(site, 'content/docs');
await fs.mkdir(generated, { recursive: true });
if ((await fs.lstat(generated)).isSymbolicLink() || !(await fs.realpath(generated)).startsWith(`${await fs.realpath(site)}${path.sep}`)) throw new Error('Generated content must stay inside the site workspace');
// Remove only known generated documentation under the site-owned directory.
for (const entry of await fs.readdir(generated)) await fs.rm(path.join(generated, entry), { recursive: true, force: true });
const pages = [];
for (const file of files) {
  const raw = await fs.readFile(path.join(docsDir, file), 'utf8');
  const front = /^---\r?\n([\s\S]*?)\r?\n---\r?\n/.exec(raw);
  if (!front) throw new Error(`${file}: missing document metadata`);
  const parsed = { data: parseYAML(front[1]), content: raw.slice(front[0].length) };
  if (typeof parsed.data.title !== 'string' || !Number.isFinite(parsed.data.nav_order)) throw new Error(`${file}: title and numeric nav_order are required for navigation`);
  const route = routes.get(file);
  const slug = route.slice('/docs/'.length).replace(/\/$/, '') || 'index';
  const output = path.join(generated, `${slug}.md`);
  await fs.mkdir(path.dirname(output), { recursive: true });
  const body = rewriteMarkdown(parsed.content, file, routes);
  await fs.writeFile(output, `---\ntitle: ${JSON.stringify(parsed.data.title)}\n---\n${body}`, 'utf8');
  pages.push({ source: file, route, slug, title: parsed.data.title, order: parsed.data.nav_order });
}
pages.sort((a, b) => a.order - b.order || a.source.localeCompare(b.source));
await fs.writeFile(path.join(generated, 'meta.json'), JSON.stringify({ title: 'Documentation', pages: pages.filter(p => !p.slug.includes('/')).map(p => p.slug) }, null, 2));
await fs.mkdir(path.join(generated, 'install'), { recursive: true });
await fs.writeFile(path.join(generated, 'install/meta.json'), JSON.stringify({ title: 'Installation', pages: ['windows', 'linux', 'macos'] }, null, 2));
// Installation's index is the canonical install.md, moved into its folder.
await fs.rename(path.join(generated, 'install.md'), path.join(generated, 'install/index.md'));
await fs.mkdir(path.join(site, 'public/assets'), { recursive: true });
const assets = path.join(site, 'public/assets');
if ((await fs.lstat(assets)).isSymbolicLink() || !(await fs.realpath(assets)).startsWith(`${await fs.realpath(site)}${path.sep}`)) throw new Error('Generated assets must stay inside the site workspace');
for (const entry of await fs.readdir(assets)) await fs.rm(path.join(assets,entry), {recursive:true,force:true});
await fs.cp(path.join(docsDir, 'assets'), path.join(site, 'public/assets'), { recursive: true });
await fs.mkdir(path.join(site, '.source'), { recursive: true });
await fs.writeFile(path.join(site, '.source/routes.json'), JSON.stringify(pages, null, 2));
const published = JSON.parse(await fs.readFile(path.join(site, 'data/published.json'), 'utf8'));
const selectedTag = process.env.SITE_RELEASE_TAG || published.tag;
if (!/^v\d+\.\d+\.\d+$/.test(selectedTag)) throw new Error('SITE_RELEASE_TAG must identify a public vMAJOR.MINOR.PATCH release');
const comparison = JSON.parse(await fs.readFile(path.join(site, 'data/comparisons', `${selectedTag}.json`), 'utf8').catch(() => { throw new Error(`Public release ${selectedTag} has no comparison snapshot; add it before deploying the site`); }));
if (comparison.tag !== selectedTag) throw new Error('Selected comparison snapshot tag does not match the public release');
await fs.writeFile(path.join(site, '.source/comparison.json'), JSON.stringify(comparison, null, 2));
console.log(`Prepared ${pages.length} canonical documents and their assets`);
