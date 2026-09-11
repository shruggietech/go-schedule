import AxeBuilder from '@axe-core/playwright'
import { expect, test } from '@playwright/test'

test.beforeEach(async ({ page }) => {
  await page.addInitScript(() => {
    const time = (index: number) => new Date(Date.UTC(2026, 8, 7, 0, index)).toISOString()
    const channels = [{ id: 'c1', name: 'Ops hook', kind: 'webhook', endpointSummary: 'https://example.test/...', hasAuthorization: true, enabled: true, updatedAt: time(0) }, { id: 'c2', name: 'Disabled hook', kind: 'webhook', endpointSummary: 'https://disabled.test/...', hasAuthorization: false, enabled: false, updatedAt: time(0) }]
    const deliveries = Array.from({ length: 200 }, (_, index) => ({ id: `delivery-${index}`, channelId: 'c1', channelName: 'Ops hook', destinationSummary: 'https://example.test/...', kind: index === 0 ? 'test' : 'task_outcome', taskId: index === 0 ? undefined : `task-${index}`, taskName: index === 0 ? undefined : `Task ${index}`, runId: index === 0 ? undefined : `run-${index}`, state: ['queued', 'retrying', 'sending', 'successful', 'failed'][index % 5], attempts: index % 3, createdAt: time(index), lastStatus: index % 5 === 4 ? 503 : undefined, lastError: index % 5 === 4 ? 'receiver unavailable' : undefined }))
    const workspace = { channels, tasks: [{ type: 'task', id: 't1', name: 'Backup', context: 'Operations' }], groups: [{ type: 'group', id: 'g1', name: 'Operations', context: 'Operations' }], coverage: [{ type: 'task', id: 't1', name: 'Backup', context: 'Operations', sourceType: 'group', sourceName: 'Operations', onSuccess: false, onFailure: true, destinationCount: 1, enabledDestinationCount: 1, enabledSuccessDestinationCount: 0, enabledFailureDestinationCount: 1 }, { type: 'group', id: 'g1', name: 'Operations', context: 'Operations', sourceType: 'group', sourceName: 'Operations', onSuccess: false, onFailure: true, destinationCount: 1, enabledDestinationCount: 1, enabledSuccessDestinationCount: 0, enabledFailureDestinationCount: 1 }], coverageComplete: true, deliveries, loadedAt: time(200) }
    let currentWorkspace = workspace
    const policy = { scope: workspace.tasks[0], directAssignments: [], effectiveSourceType: 'group', effectiveSourceId: 'g1', effectiveSourceName: 'Operations', effectiveAssignments: [{ channelId: 'c1', onFailure: true, onSuccess: false }] }
    const accepted = (action: string, extra = {}) => Promise.resolve({ action, outcome: 'accepted', message: 'Accepted.', ...extra })
    Object.assign(window, { setNotificationWorkspace: (changes: Partial<typeof workspace>) => { currentWorkspace = { ...workspace, ...changes } }, go: { main: { App: { Snapshot: () => accepted('snapshot', { generation: 1, revision: 1, state: 'connected', target: { id: 'local', displayName: 'This computer', platform: 'linux', version: '1.3.0', capabilities: ['notifications'], permissions: ['manage'] }, message: 'Connected.' }), NotificationWorkspace: () => accepted('load_notifications', { workspace: currentWorkspace }), SaveNotificationChannel: () => accepted('save_notification_channel', { workspace: currentWorkspace }), SetNotificationChannelEnabled: () => accepted('toggle_notification_channel', { workspace: currentWorkspace }), TestNotificationChannel: () => accepted('test_notification_channel', { workspace: currentWorkspace }), DeleteNotificationChannel: () => accepted('delete_notification_channel', { workspace: currentWorkspace }), NotificationPolicy: () => accepted('load_notification_policy', { policy }), SaveNotificationPolicy: () => accepted('save_notification_policy', { policy }), RetryConnection: () => accepted('retry'), Quit: () => accepted('quit') } } }, runtime: { EventsOn: () => () => undefined } })
  })
  await page.goto('/')
})

test('manages notification setup and large redacted history accessibly', async ({ page }) => {
  await page.getByRole('button', { name: 'Notifications', exact: true }).click(); await expect(page.getByRole('heading', { name: 'Notifications need attention' })).toBeVisible()
  await expect(page.getByText('1 active destination', { exact: true })).toBeVisible(); await expect(page.getByText('1 configured task', { exact: true })).toBeVisible(); await expect(page.getByText('199 recent task outcomes', { exact: true })).toBeVisible(); await expect(page.getByTestId('recent-notification-result')).toHaveCount(5)
  for (const summary of ['Manage destinations', 'Manage assignment rules', 'Inspect delivery diagnostics']) await expect(page.getByText(summary).locator('..')).not.toHaveAttribute('open', '')
  await page.getByText('Manage destinations').click()
  await page.getByRole('button', { name: 'Edit' }).first().click(); await expect(page.getByLabel('New HTTPS endpoint')).toHaveCount(0); await page.getByLabel('Replace stored endpoint').check(); await expect(page.getByLabel('New HTTPS endpoint')).toHaveValue('')
  await expect(page.getByText('Bearer secret')).toHaveCount(0); await expect(page.getByText('/hook', { exact: true })).toHaveCount(0)
  await page.getByText('Manage assignment rules').click()
  await page.getByLabel('Task or group').selectOption('task:t1'); await expect(page.getByText(/inherits from Operations/)).toBeVisible(); await page.getByLabel('Success').first().check(); await expect(page.getByText('Success notifications can be noisy')).toBeVisible()
  await page.getByText('Inspect delivery diagnostics').click(); await expect(page.getByText('200 matching deliveries')).toBeVisible()
  await page.getByLabel('State').selectOption('failed'); await expect(page.getByText('40 matching deliveries')).toBeVisible(); const delivery = page.locator('.row-select').first(); await delivery.focus(); await delivery.press('Enter'); await expect(page.getByText('receiver unavailable')).toBeVisible()
  const results = await new AxeBuilder({ page }).withTags(['wcag2a', 'wcag2aa', 'wcag22aa']).analyze(); expect(results.violations.filter((item) => item.impact === 'serious' || item.impact === 'critical')).toEqual([])
  await page.setViewportSize({ width: 800, height: 600 }); for (const zoom of ['0.8', '1', '1.5', '2']) { await page.evaluate((value) => { document.documentElement.style.zoom = value }, zoom); expect(await page.evaluate(() => document.documentElement.scrollWidth - document.documentElement.clientWidth)).toBeLessThanOrEqual(1) }
})

test('explains every overview state without serious accessibility violations', async ({ page }) => {
  await page.getByRole('button', { name: 'Notifications', exact: true }).click()
  const scenarios = [
    { title: 'Set up notifications', changes: { channels: [], coverage: [], deliveries: [] } },
    { title: 'Notifications are paused', changes: { channels: [{ id: 'c1', name: 'Ops hook', kind: 'webhook', endpointSummary: 'https://example.test/...', hasAuthorization: true, enabled: false, updatedAt: '2026-09-07T00:00:00.000Z' }], coverage: [], deliveries: [] } },
    { title: 'Choose what should notify', changes: { coverage: [], deliveries: [] } },
    { title: 'Assigned notifications are paused', changes: { coverage: [{ type: 'task', id: 't1', name: 'Backup', context: 'Operations', sourceType: 'group', sourceName: 'Operations', onSuccess: false, onFailure: true, destinationCount: 1, enabledDestinationCount: 0, enabledSuccessDestinationCount: 0, enabledFailureDestinationCount: 0 }], deliveries: [] } },
    { title: 'Notification coverage is incomplete', changes: { coverage: [], coverageComplete: false, deliveries: [] } },
    { title: 'Some notification outcomes are paused', changes: { coverage: [{ type: 'task', id: 't1', name: 'Backup', context: 'Operations', sourceType: 'group', sourceName: 'Operations', onSuccess: true, onFailure: true, destinationCount: 2, enabledDestinationCount: 1, enabledSuccessDestinationCount: 1, enabledFailureDestinationCount: 0 }], deliveries: [] } },
    { title: 'Notifications are active', changes: { deliveries: [{ id: 'healthy', channelId: 'c1', channelName: 'Ops hook', destinationSummary: 'https://example.test/...', kind: 'task_outcome', taskId: 't1', taskName: 'Backup', runId: 'r1', state: 'successful', attempts: 1, createdAt: '2026-09-07T00:01:00.000Z' }] } },
    { title: 'Notifications are working', changes: { deliveries: [{ id: 'queued', channelId: 'c1', channelName: 'Ops hook', destinationSummary: 'https://example.test/...', kind: 'task_outcome', taskId: 't1', taskName: 'Backup', runId: 'r1', state: 'queued', attempts: 0, createdAt: '2026-09-07T00:01:00.000Z' }] } },
    { title: 'A notification is retrying', changes: { deliveries: [{ id: 'retrying', channelId: 'c1', channelName: 'Ops hook', destinationSummary: 'https://example.test/...', kind: 'task_outcome', taskId: 't1', taskName: 'Backup', runId: 'r1', state: 'retrying', attempts: 1, createdAt: '2026-09-07T00:01:00.000Z' }] } },
    { title: 'Notifications need attention', changes: { deliveries: [{ id: 'failed', channelId: 'c1', channelName: 'Ops hook', destinationSummary: 'https://example.test/...', kind: 'task_outcome', taskId: 't1', taskName: 'Backup', runId: 'r1', state: 'failed', attempts: 3, createdAt: '2026-09-07T00:01:00.000Z' }] } },
  ]
  for (const scenario of scenarios) {
    await page.evaluate((changes) => (window as typeof window & { setNotificationWorkspace: (value: typeof changes) => void }).setNotificationWorkspace(changes), scenario.changes)
    await page.getByRole('button', { name: 'Refresh', exact: true }).click()
    await expect(page.getByRole('heading', { name: scenario.title })).toBeVisible()
    const results = await new AxeBuilder({ page }).withTags(['wcag2a', 'wcag2aa', 'wcag22aa']).analyze()
    expect(results.violations.filter((item) => item.impact === 'serious' || item.impact === 'critical')).toEqual([])
  }
})
