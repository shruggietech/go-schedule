import AxeBuilder from '@axe-core/playwright'
import { expect, test } from '@playwright/test'

const pages = ['Tasks', 'Task editor', 'Schedule', 'Activity', 'Targets']

for (const viewport of [{ width: 1440, height: 900 }, { width: 900, height: 650 }]) {
  test.describe(`${viewport.width}x${viewport.height}`, () => {
    test.use({ viewport })

    test('representative views retain hierarchy, target context, and horizontal fit', async ({ page }) => {
      await page.goto('/')
      for (const destination of pages) {
        await page.getByRole('button', { name: destination, exact: true }).click()
        await expect(page.getByRole('heading', { level: 1, name: destination })).toBeVisible()
        await expect(page.getByText('This computer').first()).toBeVisible()
        const overflow = await page.evaluate(() => document.documentElement.scrollWidth - document.documentElement.clientWidth)
        expect(overflow).toBeLessThanOrEqual(1)
        const screenshot = await page.screenshot({ fullPage: true })
        expect(screenshot.byteLength).toBeGreaterThan(0)
      }
    })
  })
}

test('required state and appearance samples have no serious or critical axe violations', async ({ page }) => {
  await page.goto('/')
  for (const appearance of ['light', 'dark']) {
    await page.getByLabel('Appearance').selectOption(appearance)
    for (const condition of ['connected', 'disconnected', 'degraded', 'destructive']) {
      await page.getByLabel('State').selectOption(condition)
      const results = await new AxeBuilder({ page }).withTags(['wcag2a', 'wcag2aa', 'wcag22aa']).analyze()
      expect(results.violations.filter((violation) => violation.impact === 'serious' || violation.impact === 'critical')).toEqual([])
    }
  }
})

test('keyboard navigation reaches content and changes destinations without a trap', async ({ page }) => {
  await page.goto('/')
  await page.keyboard.press('Tab')
  await expect(page.getByRole('link', { name: 'Skip to main content' })).toBeFocused()
  await page.keyboard.press('Enter')
  await expect(page.locator('#main-content')).toBeFocused()
  await page.getByRole('button', { name: 'Activity', exact: true }).focus()
  await page.keyboard.press('Enter')
  await expect(page.getByRole('heading', { level: 1, name: 'Activity' })).toBeVisible()
  await page.keyboard.press('Tab')
  await expect(page.locator(':focus')).toBeVisible()
})

test('primary workflow remains usable at 200 percent zoom', async ({ page }) => {
  await page.setViewportSize({ width: 900, height: 650 })
  await page.goto('/')
  await page.evaluate(() => { document.documentElement.style.zoom = '2' })
  await expect(page.getByText('This computer').first()).toBeVisible()
  await expect(page.getByRole('button', { name: 'Create task' })).toBeVisible()
  const overflow = await page.evaluate(() => document.documentElement.scrollWidth - document.documentElement.clientWidth)
  expect(overflow).toBeLessThanOrEqual(1)
})
