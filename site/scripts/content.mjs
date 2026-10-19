import path from 'node:path';

export function routeFor(file) {
  if (file === 'README.md') return '/docs/';
  const platform = /^INSTALL-(windows|linux|macos)\.md$/.exec(file);
  if (platform) return `/docs/install/${platform[1]}/`;
  return `/docs/${file.replace(/\.md$/, '').toLowerCase()}/`;
}

function destination(url, file, routes) {
  if (/^(?:[a-z][a-z\d+.-]*:|\/|#)/i.test(url)) return url;
  const [resource, fragment] = url.split(/#(.*)/s);
  const resolved = path.posix.normalize(path.posix.join(path.posix.dirname(file), resource));
  if (resolved.startsWith('../')) throw new Error(`${file}: link escapes canonical docs: ${url}`);
  const suffix = fragment === undefined ? '' : `#${fragment}`;
  if (resolved.startsWith('assets/')) return `/${resolved}${suffix}`;
  if (!routes.has(resolved)) throw new Error(`${file}: unknown document ${url}; add it to the route inventory or use an absolute repository URL`);
  return `${routes.get(resolved)}${suffix}`;
}

export function rewriteMarkdown(input, file, routes) {
  let fence;
  return input.split('\n').map(line => {
    const marker = /^\s*(`{3,}|~{3,})/.exec(line);
    if (marker) {
      if (!fence) fence = marker[1][0];
      else if (fence === marker[1][0]) fence = undefined;
      return line;
    }
    if (fence) return line;
    // Inline code remains byte-for-byte literal, including command examples.
    return line.split(/(`+[^`]*`+)/g).map((part, index) => index % 2 ? part : part
      .replace(/\]\(([^\s)]+)([^)]*)\)/g, (_, url, title) => `](${destination(url, file, routes)}${title})`)
      .replace(/((?:href|src)=["'])([^"']+)(["'])/g, (_, start, url, end) => `${start}${destination(url, file, routes)}${end}`))
      .join('');
  }).join('\n');
}

export function legacyHTML(target) {
  const safe = JSON.stringify(target).replace(/</g, '\\u003c');
  return `<!doctype html><html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width"><title>Documentation moved | go-schedule</title><link rel="canonical" href="https://go-schedule.io${target}"><meta http-equiv="refresh" content="0;url=${target}"></head><body><a href="${target}">Continue to the documentation</a><script>location.replace(${safe}+location.search+location.hash)</script></body></html>`;
}
