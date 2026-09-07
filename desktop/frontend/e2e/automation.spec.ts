import AxeBuilder from '@axe-core/playwright'
import { expect, test } from '@playwright/test'

test.beforeEach(async ({ page }) => {
  await page.addInitScript(() => {
    const tasks = [{ id: 'task-1', name: 'Build', readiness: 'ready', reason: 'Ready.' }]
    const triggers = Array.from({ length: 100 }, (_, index) => ({ id: `trigger-${index}`, name: `Hook ${String(index + 1).padStart(3, '0')}`, targetTaskId: 'task-1', targetTaskName: 'Build', enabled: true, readiness: index === 99 ? 'target_missing' : 'ready', reason: index === 99 ? 'Target task is missing.' : 'Ready.', updatedAt: '2026-09-07T00:00:00Z' }))
    const automation = { tasks, chains: [], triggers, triggerSets: [], watchers: [{ id: 'watcher-1', name: 'Incoming', kind: 'directory', path: '/a/very/long/path/that/remains/visible/to/assist/recovery', pattern: '*.csv', recursive: true, debounce: '500ms', stability: '1s', targetTaskId: 'task-1', targetTaskName: 'Build', enabled: true, health: 'degraded', healthReason: 'Path is unavailable.', readiness: 'ready', reason: '', updatedAt: '2026-09-07T00:00:00Z' }], loadedAt: '2026-09-07T00:00:00Z' }
    const result = (action: string) => Promise.resolve({ action, outcome: 'accepted', message: 'Accepted.', workspace: automation })
    Object.assign(window, { go: { main: { App: { Snapshot: () => Promise.resolve({ generation: 1, revision: 1, state: 'connected', target: { id: 'local', displayName: 'This computer', platform: 'linux', version: '1.0.0', capabilities: ['automation'], permissions: ['manage'] }, message: 'Connected.' }), AutomationWorkspace: () => result('load'), RetryConnection: () => result('retry'), Quit: () => result('quit') } } }, runtime: { EventsOn: () => () => undefined } })
  })
  await page.goto('/')
  await page.getByRole('button', { name: 'Automation Sources', exact: true }).click()
})

test('keeps one hundred sources searchable, honest, accessible, and responsive', async ({ page }) => {
  await expect(page.getByText('100 matching sources')).toBeVisible()
  const search = page.getByRole('searchbox', { name: 'Search' }); await search.fill('Hook 100'); await expect(page.getByRole('heading', { name: 'Hook 100' })).toBeVisible(); await expect(search).toBeFocused()
  await page.getByRole('checkbox', { name: 'Needs attention only' }).check(); await expect(page.getByText('Target task is missing.')).toBeVisible()
  await search.fill(''); await page.getByLabel('Source type').selectOption('watchers'); await expect(page.getByText('Path is unavailable.')).toBeVisible()
  const results = await new AxeBuilder({ page }).withTags(['wcag2a', 'wcag2aa', 'wcag22aa']).analyze()
  expect(results.violations.filter((item) => item.impact === 'serious' || item.impact === 'critical')).toEqual([])
  await page.setViewportSize({ width: 900, height: 650 }); await page.evaluate(() => { document.documentElement.style.zoom = '2' }); expect(await page.evaluate(() => document.documentElement.scrollWidth - document.documentElement.clientWidth)).toBeLessThanOrEqual(1)
})
