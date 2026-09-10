import AxeBuilder from '@axe-core/playwright'
import { expect, test } from '@playwright/test'

const snapshot = {
  generation: 4,
  revision: 8,
  state: 'recovering',
  target: { id: 'daemon-identity', profileId: 'profile-id', kind: 'remote', displayName: 'Production', endpoint: 'https://example.test', platform: 'linux', capabilities: ['tasks'], permissions: ['read'] },
  message: 'The remote scheduler is temporarily unreachable.',
  lastSuccessfulAt: '2026-09-07T12:00:00Z',
  stale: true,
  retryAttempt: 2,
  nextRetryAt: '2026-09-07T12:00:10Z',
  recovery: 'automatic',
}

for (const viewport of [{ width: 1440, height: 900 }, { width: 900, height: 650 }]) {
  test(`presents stale recovery at ${viewport.width}x${viewport.height}`, async ({ page }) => {
    await page.setViewportSize(viewport)
    await page.addInitScript((value) => { Object.assign(window, { go: { main: { App: { Snapshot: async () => value } } } }) }, snapshot)
    await page.goto('/')
    await expect(page.getByText('Data may be stale')).toBeVisible()
    await page.getByRole('button', { name: 'Connection details' }).click()
    await expect(page.getByRole('dialog')).toContainText('Attempt 2 scheduled for 2026-09-07T12:00:10Z')
    expect(await page.evaluate(() => document.documentElement.scrollWidth - document.documentElement.clientWidth)).toBeLessThanOrEqual(1)
    const results = await new AxeBuilder({ page }).withTags(['wcag2a', 'wcag2aa', 'wcag22aa']).analyze()
    expect(results.violations.filter((item) => item.impact === 'serious' || item.impact === 'critical')).toEqual([])
  })
}
