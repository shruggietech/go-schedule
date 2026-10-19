import test from 'node:test';
import assert from 'node:assert/strict';
import fs from 'node:fs/promises';
import { parse } from 'yaml';

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
