import AxeBuilder from '@axe-core/playwright'
import { expect, test } from '@playwright/test'

for (const viewport of [{ width: 1440, height: 900 }, { width: 900, height: 650 }]) {
  test.describe(`${viewport.width}x${viewport.height}`, () => {
    test.use({ viewport })
    test('keeps long page and target identity reachable without page overflow', async ({ page }) => {
      await page.goto('/')
      for (const destination of ['Tasks', 'Schedule', 'Activity', 'Connections', 'Settings']) {
        await page.getByRole('button', { name: destination, exact: true }).click()
        await expect(page.getByRole('heading', { level: 1, name: destination })).toBeVisible()
        await expect(page.getByText('This computer').first()).toBeVisible()
        expect(await page.evaluate(() => document.documentElement.scrollWidth - document.documentElement.clientWidth)).toBeLessThanOrEqual(1)
      }
    })
  })
}

test('appearance and browser-unavailable state have no serious or critical axe findings', async ({ page }) => {
  await page.goto('/')
  for (const appearance of ['light', 'dark', 'system']) {
    await page.getByLabel('Appearance').selectOption(appearance)
    const results = await new AxeBuilder({ page }).withTags(['wcag2a', 'wcag2aa', 'wcag22aa']).analyze()
    expect(results.violations.filter((item) => item.impact === 'serious' || item.impact === 'critical')).toEqual([])
  }
})

test('keyboard and 200 percent zoom retain primary controls', async ({ page }) => {
  await page.setViewportSize({ width: 900, height: 650 }); await page.goto('/')
  await page.keyboard.press('Tab'); await expect(page.getByRole('link', { name: 'Skip to main content' })).toBeFocused(); await page.keyboard.press('Enter'); await expect(page.locator('#main-content')).toBeFocused()
  await page.evaluate(() => { document.documentElement.style.zoom = '2' })
  await expect(page.getByText('This computer').first()).toBeVisible()
  expect(await page.evaluate(() => document.documentElement.scrollWidth - document.documentElement.clientWidth)).toBeLessThanOrEqual(1)
})

test('uses only local assets', async ({ page }) => {
  const remote: string[] = []
  page.on('request', (request) => { if (!request.url().startsWith('http://127.0.0.1:4173')) remote.push(request.url()) })
  await page.goto('/'); await page.waitForLoadState('networkidle'); expect(remote).toEqual([])
})

test('respects reduced motion preference', async ({ page }) => {
  await page.emulateMedia({ reducedMotion: 'reduce' })
  await page.goto('/')
  const durationMilliseconds = await page.locator('.button').first().evaluate((element) => {
    const duration = getComputedStyle(element).transitionDuration
    return duration.endsWith('ms') ? Number.parseFloat(duration) : Number.parseFloat(duration) * 1000
  })
  expect(durationMilliseconds).toBeLessThanOrEqual(0.01)
})
