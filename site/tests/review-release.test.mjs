import test from 'node:test';
import assert from 'node:assert/strict';
import fs from 'node:fs/promises';
import path from 'node:path';
import os from 'node:os';
import { fileURLToPath } from 'node:url';
import { execFileSync } from 'node:child_process';
import { createHash } from 'node:crypto';

test('review helper creates exact content digests and cannot overwrite a review', async t => {
  const root = await fs.mkdtemp(path.join(os.tmpdir(), 'goschedule-comparison-review-'));
  t.after(async () => {
    assert.ok((await fs.realpath(root)).startsWith(`${await fs.realpath(os.tmpdir())}${path.sep}`));
    await fs.rm(root, {recursive:true,force:true});
  });
  const repository = fileURLToPath(new URL('../../', import.meta.url));
  const helper = 'site/scripts/review-release.mjs';
  const snapshot = await fs.readFile(path.join(repository,'site/data/comparisons/v1.5.1.json'));
  const section = '## [1.5.1] - 2026-09-26\n\nReleased capability.\n\n';
  const notes = '# Release\n\nReleased capability.\n';
  const inputs = {
    [helper]: await fs.readFile(path.join(repository,helper)),
    'site/data/comparisons/v1.5.1.json': snapshot,
    'CHANGELOG.md': '# Changelog\n\n## [Unreleased]\n\n'+section+'## [1.5.0] - 2026-09-23\n\nOlder.\n',
    '.github/release-notes/v1.5.1.md': notes,
  };
  for (const [relative,data] of Object.entries(inputs)) {
    const file=path.join(root,relative); await fs.mkdir(path.dirname(file),{recursive:true}); await fs.writeFile(file,data);
  }
  const run=tag=>execFileSync(process.execPath,[path.join(root,helper),tag],{windowsHide:true,stdio:'pipe'});
  run('v1.5.1');
  const target=path.join(root,'site/data/reviews/v1.5.1.json');
  const raw=await fs.readFile(target,'utf8'); const review=JSON.parse(raw);
  const digest=value=>createHash('sha256').update(value).digest('hex');
  assert.equal(review.source_ref,'refs/tags/v1.5.1');
  assert.equal(review.dataset_sha256,digest(snapshot));
  assert.equal(review.changelog_sha256,digest(section));
  assert.equal(review.notes_sha256,digest(notes));
  assert.deepEqual(Object.keys(review.decisions),JSON.parse(snapshot).capabilities.map(row=>row.id));
  assert.ok(Object.values(review.decisions).every(value=>value==='review-required'));
  assert.throws(()=>run('v1.5.1'),/EEXIST/);
  assert.equal(await fs.readFile(target,'utf8'),raw);
  assert.throws(()=>run('../../escape'),/exact vMAJOR.MINOR.PATCH/);
});
