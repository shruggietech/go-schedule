import test from 'node:test';
import assert from 'node:assert/strict';
import { routeFor, rewriteMarkdown, legacyHTML } from '../scripts/content.mjs';

const routes = new Map([['README.md', '/docs/'], ['cli.md', '/docs/cli/'], ['INSTALL-windows.md', '/docs/install/windows/']]);

test('source routes retain explicit installation and root mapping', () => {
  assert.equal(routeFor('README.md'), '/docs/');
  assert.equal(routeFor('INSTALL-windows.md'), '/docs/install/windows/');
  assert.equal(routeFor('api.md'), '/docs/api/');
});
test('document links retain anchors and asset paths resolve', () => {
  assert.equal(rewriteMarkdown('[CLI](cli.md#flags) ![Mark](assets/brand/mark.svg)', 'README.md', routes), '[CLI](/docs/cli/#flags) ![Mark](/assets/brand/mark.svg)');
});
test('fenced and inline command examples remain literal', () => {
  const input = '`[example](cli.md)`\n\n```text\n[x](cli.md) <id> {"x":1}\n```\n[CLI](cli.md)';
  assert.equal(rewriteMarkdown(input, 'README.md', routes), input.replace(/\[CLI\]\(cli.md\)$/, '[CLI](/docs/cli/)'));
});
test('unknown and escaping local links fail with a correction', () => {
  assert.throws(() => rewriteMarkdown('[x](missing.md)', 'cli.md', routes), /missing.md/);
  assert.throws(() => rewriteMarkdown('[x](../outside.md)', 'cli.md', routes), /outside/);
});
test('legacy redirect retains query and fragment and provides a visible fallback', () => {
  const html = legacyHTML('/docs/cli/');
  assert.match(html, /location.search/);
  assert.match(html, /location.hash/);
  assert.match(html, /Continue to the documentation/);
  assert.match(html, /rel="canonical"/);
});
