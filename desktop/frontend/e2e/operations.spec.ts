import AxeBuilder from '@axe-core/playwright'
import { expect, test } from '@playwright/test'

test.beforeEach(async ({ page }) => {
  await page.addInitScript(() => {
    const time = (index: number) => new Date(Date.UTC(2026, 8, 7, 0, index)).toISOString()
    const occurrences = Array.from({ length: 105 }, (_, index) => ({ id: `occurrence-${index}`, taskId: `task-${index}`, taskName: `Task ${index}`, runId: index % 2 ? `run-${index}` : undefined, time: time(index), kind: index % 2 ? 'recorded' : 'prediction', state: index % 2 ? 'success' : 'upcoming', outcome: index % 2 ? 'success' : undefined }))
    const runs = Array.from({ length: 105 }, (_, index) => ({ id: `run-${index}`, taskId: `task-${index}`, scheduledFor: time(index), startedAt: time(index), endedAt: time(index + 1), state: index === 104 ? 'failure' : 'success', outcome: index === 104 ? 'failure' : 'success', exitCode: index === 104 ? 7 : 0, output: index === 104 ? 'failure output' : '', outputTruncated: false, trigger: 'schedule' }))
    const activity = { runs, logs: [{ id: 'log-1', time: time(120), severity: 'warning', source: 'engine', message: 'slow' }], alerts: [{ id: 'alert-1', time: time(121), severity: 'error', kind: 'run_failed', message: 'failed', acknowledged: false }], logPath: '/var/log/goschedule.log', loadedAt: time(122) }
    const accepted = (action: string, extra = {}) => Promise.resolve({ action, outcome: 'accepted', message: 'Accepted.', ...extra })
    Object.assign(window, { go: { main: { App: { Snapshot: () => accepted('snapshot', { generation: 1, revision: 1, state: 'connected', target: { id: 'local', displayName: 'This computer', platform: 'linux', version: '1.2.0', capabilities: ['schedule', 'activity'], permissions: ['manage'] }, message: 'Connected.' }), ScheduleWindow: () => accepted('load_schedule', { schedule: { from: time(0), to: time(200), loadedAt: time(122), occurrences } }), ActivityWorkspace: () => accepted('load_activity', { activity }), AcknowledgeAlert: () => accepted('acknowledge_alerts', { activity }), AcknowledgeAlerts: () => accepted('acknowledge_alerts', { activity }), RetryConnection: () => accepted('retry'), Quit: () => accepted('quit') } } }, runtime: { EventsOn: () => () => undefined } })
  })
  await page.goto('/')
})

test('keeps large Schedule and Activity workspaces accessible and focused', async ({ page }) => {
  await page.getByRole('button', { name: 'Schedule', exact: true }).click(); await expect(page.getByText('105 occurrences', { exact: false })).toBeVisible(); const first = page.getByRole('button', { name: /2026/ }).first(); await first.focus(); await first.press('Enter'); await expect(page.getByRole('heading', { name: 'Task 0' })).toBeVisible()
  await page.getByLabel('View').selectOption('calendar'); await expect(page.getByRole('button', { name: /105 occurrences/ })).toBeVisible()
  await page.getByRole('button', { name: 'Activity', exact: true }).click(); await expect(page.getByText('107 matching records')).toBeVisible(); const search = page.getByLabel('Search'); await search.fill('task-104'); await expect(page.getByText('1 matching records')).toBeVisible(); await expect(search).toBeFocused()
  const results = await new AxeBuilder({ page }).withTags(['wcag2a', 'wcag2aa', 'wcag22aa']).analyze(); expect(results.violations.filter((item) => item.impact === 'serious' || item.impact === 'critical')).toEqual([])
  await page.setViewportSize({ width: 900, height: 650 }); await page.evaluate(() => { document.documentElement.style.zoom = '2' }); expect(await page.evaluate(() => document.documentElement.scrollWidth - document.documentElement.clientWidth)).toBeLessThanOrEqual(1)
})
