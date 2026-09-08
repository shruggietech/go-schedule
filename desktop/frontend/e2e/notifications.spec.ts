import AxeBuilder from '@axe-core/playwright'
import { expect, test } from '@playwright/test'

test.beforeEach(async ({ page }) => {
  await page.addInitScript(() => {
    const time = (index: number) => new Date(Date.UTC(2026, 8, 7, 0, index)).toISOString()
    const channels = [{ id: 'c1', name: 'Ops hook', kind: 'webhook', endpointSummary: 'https://example.test/...', hasAuthorization: true, enabled: true, updatedAt: time(0) }, { id: 'c2', name: 'Disabled hook', kind: 'webhook', endpointSummary: 'https://disabled.test/...', hasAuthorization: false, enabled: false, updatedAt: time(0) }]
    const deliveries = Array.from({ length: 200 }, (_, index) => ({ id: `delivery-${index}`, channelId: 'c1', channelName: 'Ops hook', destinationSummary: 'https://example.test/...', kind: index === 0 ? 'test' : 'task_outcome', taskId: index === 0 ? undefined : `task-${index}`, taskName: index === 0 ? undefined : `Task ${index}`, runId: index === 0 ? undefined : `run-${index}`, state: ['queued', 'retrying', 'sending', 'successful', 'failed'][index % 5], attempts: index % 3, createdAt: time(index), lastStatus: index % 5 === 4 ? 503 : undefined, lastError: index % 5 === 4 ? 'receiver unavailable' : undefined }))
    const workspace = { channels, tasks: [{ type: 'task', id: 't1', name: 'Backup', context: 'Operations' }], groups: [{ type: 'group', id: 'g1', name: 'Operations', context: 'Operations' }], deliveries, loadedAt: time(200) }
    const policy = { scope: workspace.tasks[0], directAssignments: [], effectiveSourceType: 'group', effectiveSourceId: 'g1', effectiveSourceName: 'Operations', effectiveAssignments: [{ channelId: 'c1', onFailure: true, onSuccess: false }] }
    const accepted = (action: string, extra = {}) => Promise.resolve({ action, outcome: 'accepted', message: 'Accepted.', ...extra })
    Object.assign(window, { go: { main: { App: { Snapshot: () => accepted('snapshot', { generation: 1, revision: 1, state: 'connected', target: { id: 'local', displayName: 'This computer', platform: 'linux', version: '1.3.0', capabilities: ['notifications'], permissions: ['manage'] }, message: 'Connected.' }), NotificationWorkspace: () => accepted('load_notifications', { workspace }), SaveNotificationChannel: () => accepted('save_notification_channel', { workspace }), SetNotificationChannelEnabled: () => accepted('toggle_notification_channel', { workspace }), TestNotificationChannel: () => accepted('test_notification_channel', { workspace }), DeleteNotificationChannel: () => accepted('delete_notification_channel', { workspace }), NotificationPolicy: () => accepted('load_notification_policy', { policy }), SaveNotificationPolicy: () => accepted('save_notification_policy', { policy }), RetryConnection: () => accepted('retry'), Quit: () => accepted('quit') } } }, runtime: { EventsOn: () => () => undefined } })
  })
  await page.goto('/')
})

test('manages notification setup and large redacted history accessibly', async ({ page }) => {
  await page.getByRole('button', { name: 'Notifications', exact: true }).click(); await expect(page.getByText('200 matching deliveries')).toBeVisible()
  await page.getByRole('button', { name: 'Edit' }).first().click(); await expect(page.getByLabel('New HTTPS endpoint')).toHaveCount(0); await page.getByLabel('Replace stored endpoint').check(); await expect(page.getByLabel('New HTTPS endpoint')).toHaveValue('')
  await expect(page.getByText('Bearer secret')).toHaveCount(0); await expect(page.getByText('/hook', { exact: true })).toHaveCount(0)
  await page.getByLabel('Task or group').selectOption('task:t1'); await expect(page.getByText(/inherits from Operations/)).toBeVisible(); await page.getByLabel('Success').first().check(); await expect(page.getByText('Success notifications can be noisy')).toBeVisible()
  await page.getByLabel('State').selectOption('failed'); await expect(page.getByText('40 matching deliveries')).toBeVisible(); const delivery = page.locator('.row-select').first(); await delivery.focus(); await delivery.press('Enter'); await expect(page.getByText('receiver unavailable')).toBeVisible()
  const results = await new AxeBuilder({ page }).withTags(['wcag2a', 'wcag2aa', 'wcag22aa']).analyze(); expect(results.violations.filter((item) => item.impact === 'serious' || item.impact === 'critical')).toEqual([])
  await page.setViewportSize({ width: 900, height: 650 }); for (const zoom of ['0.8', '1', '1.5', '2']) { await page.evaluate((value) => { document.documentElement.style.zoom = value }, zoom); expect(await page.evaluate(() => document.documentElement.scrollWidth - document.documentElement.clientWidth)).toBeLessThanOrEqual(1) }
})
