import AxeBuilder from '@axe-core/playwright'
import { expect, test } from '@playwright/test'

test.beforeEach(async ({ page }) => {
  await page.addInitScript(() => {
    const groups = Array.from({ length: 20 }, (_, index) => ({ id: `group-${index}`, name: `Level ${index + 1}`, parentId: index ? `group-${index - 1}` : '', path: Array.from({ length: index + 1 }, (_, part) => `Level ${part + 1}`).join(' / '), depth: index, declaredEnabled: true, effectiveEnabled: true, effectiveReason: '', childCount: index === 19 ? 0 : 1, taskCount: 5, descendantCount: 19 - index, updatedAt: '2026-09-07T00:00:00Z' }))
    const tasks = Array.from({ length: 100 }, (_, index) => ({ id: `task-${index}`, name: `Task ${String(index + 1).padStart(3, '0')}`, groupId: `group-${index % 20}`, groupPath: groups[index % 20].path, commandConfigured: true, declaredEnabled: false, effectiveState: 'task_disabled', effectiveReason: 'Automatic activation is disabled.', lifecycle: 'active', timezone: 'UTC', scheduleSummary: 'Every day at 09:00', policySummary: 'Default calendar policies', nextRuns: [], updatedAt: '2026-09-07T00:00:00Z' }))
    const operation = (action: string) => Promise.resolve({ action, outcome: 'accepted', message: 'Accepted.', workspace: { tasks, groups, loadedAt: '2026-09-07T00:00:00Z' } })
    Object.assign(window, { go: { main: { App: { Snapshot: () => Promise.resolve({ generation: 1, revision: 1, state: 'connected', target: { id: 'local', displayName: 'This computer', platform: 'linux', version: '1.0.0', capabilities: ['tasks'], permissions: ['manage'] }, message: 'Connected.' }), RetryConnection: () => operation('retry'), Quit: () => operation('quit'), Workspace: () => operation('load'), Task: (id: string) => Promise.resolve({ action: 'load', outcome: 'accepted', message: 'Loaded.', task: { id, name: 'Task 001', groupId: 'group-0', commandLine: '', workingDir: '', environment: [], stdin: '', runAs: '', enabled: false, timezone: 'UTC', mode: 'recurring', schedule: 'every day at 09:00', scheduleSyntax: 'human', at: '', overlapPolicy: 'queue_one', catchupPolicy: 'one', missingDatePolicy: 'skip', timeBasis: 'wall_clock', dstGapPolicy: 'next_valid', dstOverlapPolicy: 'first', scheduleSummary: 'Every day', policySummary: '', readiness: 'disabled', updatedAt: '2026-09-07T00:00:00Z', nextRuns: [] } }), PreviewTask: () => operation('preview'), SaveTask: () => operation('save_task'), RunTask: () => operation('run_task'), SetTaskEnabled: () => operation('toggle_task'), DeleteTask: () => operation('delete_task'), SaveGroup: () => operation('save_group'), SetGroupEnabled: () => operation('toggle_group'), DeleteGroup: () => operation('delete_group') } } }, runtime: { EventsOn: () => () => undefined } })
  })
  await page.goto('/')
  await expect(page.getByRole('heading', { level: 1, name: 'Tasks' })).toBeVisible()
})

test('keeps one hundred tasks and twenty group levels searchable and accessible', async ({ page }) => {
  await expect(page.getByText('100 tasks')).toBeVisible()
  await page.getByRole('searchbox', { name: 'Search' }).fill('Task 100')
  await expect(page.getByRole('button', { name: 'Task 100' })).toBeVisible()
  await expect(page.getByRole('button', { name: 'Task 001' })).toHaveCount(0)
  const results = await new AxeBuilder({ page }).withTags(['wcag2a', 'wcag2aa', 'wcag22aa']).analyze()
  expect(results.violations.filter((item) => item.impact === 'serious' || item.impact === 'critical')).toEqual([])
})

test('inserts the Linux example once and remains usable at narrow width and zoom', async ({ page }) => {
  await page.setViewportSize({ width: 900, height: 650 })
  const initialHeight = await page.locator('#main-content').evaluate((element) => element.scrollHeight)
  await page.getByRole('button', { name: 'Create task' }).click()
  await expect(page.getByRole('dialog', { name: 'Create task' })).toBeVisible()
  expect(await page.locator('#main-content').evaluate((element) => element.scrollHeight)).toBe(initialHeight)
  await expect(page.getByText('Advanced settings').locator('..')).not.toHaveAttribute('open')
  const command = page.getByLabel('Command line')
  await command.focus(); await page.keyboard.press('Tab'); await expect(command).toHaveValue('uname -a')
  await page.keyboard.press('Tab'); await expect(page.getByRole('button', { name: 'Insert example' })).toBeFocused()
  await page.evaluate(() => { document.documentElement.style.zoom = '2' })
  expect(await page.evaluate(() => document.documentElement.scrollWidth - document.documentElement.clientWidth)).toBeLessThanOrEqual(1)
  expect(await page.getByRole('dialog').evaluate((element) => element.scrollWidth - element.clientWidth)).toBeLessThanOrEqual(1)
})
