import AxeBuilder from '@axe-core/playwright'
import { expect, test } from '@playwright/test'

test('keeps search keyboard-operable at compact viewport and 200 percent zoom', async ({ page }) => {
  await page.setViewportSize({ width: 800, height: 600 })
  await page.goto('/')
  await page.getByRole('button', { name: 'Search', exact: true }).click()
  const query = page.getByLabel('Search all systems')
  await query.focus()
  await expect(query).toBeFocused()
  await query.fill('archive')
  await page.evaluate(() => { document.documentElement.style.zoom = '2' })
  await expect(page.getByRole('button', { name: 'Search', exact: true }).last()).toBeVisible()
  await expect(page.getByLabel('task')).toBeVisible()
  expect(await page.evaluate(() => document.documentElement.scrollWidth - document.documentElement.clientWidth)).toBeLessThanOrEqual(1)
})

test('has no serious or critical accessibility findings', async ({ page }) => {
  await page.goto('/')
  await page.getByRole('button', { name: 'Search', exact: true }).click()
  const results = await new AxeBuilder({ page }).withTags(['wcag2a', 'wcag2aa', 'wcag22aa']).analyze()
  expect(results.violations.filter((item) => item.impact === 'serious' || item.impact === 'critical')).toEqual([])
})

test('keeps 100 duplicate-labeled profiles and one partial failure distinguishable', async ({ page }) => {
  await page.addInitScript(() => {
    const observations = Array.from({ length: 100 }, (_, index) => {
      const id = `daemon-${String(index).padStart(3, '0')}`
      if (index === 99) return { registration: { key: `profile-${index}`, kind: 'remote', label: 'Production', daemonId: id, shortDaemonId: id }, state: 'timed_out', truncated: false, matches: [], failure: { state: 'timed_out', message: 'Timed out.', action: 'Try again.' } }
      return { registration: { key: `profile-${index}`, kind: 'remote', label: 'Production', daemonId: id, shortDaemonId: id }, state: 'connected', observedAt: '2026-09-21T12:00:00Z', truncated: false, matches: [{ registrationKey: `profile-${index}`, expectedDaemonId: id, sourceLabel: 'Production', sourceShortId: id, result: { kind: 'task', object_id: `task-${index}`, task_id: `task-${index}`, name: 'Archive', action_hints: ['run_now'] }, availableActions: ['open', 'run_now'] }] }
    })
    Object.assign(window, { go: { main: { App: { Snapshot: async () => ({ generation: 1, revision: 1, state: 'connected', target: { id: 'local-daemon', kind: 'local', displayName: 'This computer', platform: 'windows', capabilities: ['search'], permissions: ['read', 'manage'] }, message: 'Connected.' }), RetryConnection: async () => ({ action: 'retry', outcome: 'accepted', message: 'Retrying.' }), Quit: async () => ({ action: 'quit', outcome: 'accepted', message: 'Closing.' }), SearchAcrossSystems: async (request: { query: string }) => ({ generation: 1, query: request.query, startedAt: '2026-09-21T12:00:00Z', completedAt: '2026-09-21T12:00:01Z', complete: true, observations }) } } }, runtime: { EventsOn: () => () => undefined } })
  })
  await page.goto('/')
  await page.getByRole('button', { name: 'Search', exact: true }).click()
  await page.getByLabel('Search all systems').fill('archive')
  await page.getByRole('button', { name: 'Search', exact: true }).last().click()
  await expect(page.getByText('100 systems searched; 99 results.')).toBeVisible()
  await expect(page.getByText('Timed out. Try again.')).toBeVisible()
  await expect(page.getByText('daemon-000', { exact: true })).toBeVisible()
  await expect(page.getByText('daemon-098', { exact: true })).toBeVisible()
})
