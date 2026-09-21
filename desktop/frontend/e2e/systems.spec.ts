import AxeBuilder from '@axe-core/playwright'
import { expect, test } from '@playwright/test'

test('keeps All Systems operable at the supported compact viewport and zoom', async ({ page }) => {
  await page.setViewportSize({ width: 800, height: 600 })
  await page.goto('/')
  await page.getByRole('button', { name: 'All Systems' }).click()
  await expect(page.getByRole('heading', { level: 1, name: 'All Systems' })).toBeVisible()
  await expect(page.getByRole('button', { name: 'Refresh all' })).toBeVisible()
  await page.getByLabel('Search').focus()
  await expect(page.getByLabel('Search')).toBeFocused()
  await page.evaluate(() => { document.documentElement.style.zoom = '2' })
  await expect(page.getByLabel('Connection state')).toBeVisible()
  await expect(page.getByLabel('Sort by')).toBeVisible()
  expect(await page.evaluate(() => document.documentElement.scrollWidth - document.documentElement.clientWidth)).toBeLessThanOrEqual(1)
})

test('has no serious or critical accessibility findings', async ({ page }) => {
  await page.goto('/')
  await page.getByRole('button', { name: 'All Systems' }).click()
  const results = await new AxeBuilder({ page }).withTags(['wcag2a', 'wcag2aa', 'wcag22aa']).analyze()
  expect(results.violations.filter((item) => item.impact === 'serious' || item.impact === 'critical')).toEqual([])
})
