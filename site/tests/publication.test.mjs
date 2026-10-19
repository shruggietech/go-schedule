import test from 'node:test';
import assert from 'node:assert/strict';
import fs from 'node:fs/promises';
import { parse } from 'yaml';

// Remove this transition regression only after the replacement site is live.
test('pre-cutover readers retain live README links and the legacy Pages publishing inputs', async () => {
  const read = name => fs.readFile(new URL(`../../${name}`,import.meta.url),'utf8');
  const readme = await read('README.md');
  assert.ok(readme.includes('https://shruggietech.github.io/go-schedule/'));
  assert.doesNotMatch(readme,/https:\/\/go-schedule\.io(?:\/|["\)])/);
  assert.ok(readme.includes('[Comparison data](site/data/comparisons/v1.5.1.json)'));
  const legacy = parse(await read('docs/_config.yml'));
  assert.equal(legacy.url,'https://shruggietech.github.io');
  assert.equal(legacy.baseurl,'/go-schedule');
  assert.equal(legacy.remote_theme,'just-the-docs/just-the-docs@v0.4.2');
  assert.equal(legacy.search_enabled,true);
  assert.equal(legacy.relative_links.enabled,true);
  for (const file of ['docs/_includes/head_custom.html','docs/_includes/nav_footer_custom.html','docs/_sass/color_schemes/goschedule.scss','docs/_sass/custom/custom.scss']) assert.ok((await read(file)).trim(),`${file} must remain usable until cutover`);
});

test('website refresh follows both successful promotion workflows despite token event suppression', async () => {
  const workflow = async name => parse(await fs.readFile(new URL(`../../.github/workflows/${name}`,import.meta.url),'utf8'));
  const pages = await workflow('pages.yml');
  assert.deepEqual(pages.on.workflow_run.types,['completed']);
  for (const file of ['promote-release.yml','promote-release-waived.yml']) {
    const promotion = await workflow(file);
    assert.ok(pages.on.workflow_run.workflows.includes(promotion.name),`${file} must refresh the public comparison`);
  }
  assert.equal(pages.jobs.build.if,"github.event_name != 'workflow_run' || github.event.workflow_run.conclusion == 'success'");
  assert.equal(pages.jobs.build.steps[0].with.ref,'main');
  assert.equal(pages.permissions.contents,'read');
  assert.equal(pages.jobs.deploy.needs,'build');
  const selection=pages.jobs.build.steps.find(step=>step.name==='Select the latest public release comparison');
  assert.match(selection.run,/releases\/latest/);
  assert.match(selection.run,/SITE_RELEASE_TAG=/);
});
