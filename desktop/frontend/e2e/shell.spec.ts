import AxeBuilder from '@axe-core/playwright'
import { expect, test } from '@playwright/test'

function contrast(foreground: string, background: string) {
  const channel = (value: number) => { const normalized = value / 255; return normalized <= 0.04045 ? normalized / 12.92 : ((normalized + 0.055) / 1.055) ** 2.4 }
  const values = (color: string) => color.match(/[\d.]+/g)!.slice(0, 3).map(Number)
  const luminance = (color: string) => { const [red, green, blue] = values(color).map(channel); return 0.2126 * red + 0.7152 * green + 0.0722 * blue }
  const [lighter, darker] = [luminance(foreground), luminance(background)].sort((a, b) => b - a)
  return (lighter + 0.05) / (darker + 0.05)
}

for (const viewport of [{ width: 1440, height: 900 }, { width: 900, height: 650 }, { width: 800, height: 600 }]) {
  test.describe(`${viewport.width}x${viewport.height}`, () => {
    test.use({ viewport })
    test('keeps long page and target identity reachable without page overflow', async ({ page }) => {
      await page.goto('/')
      for (const destination of ['Tasks', 'Automation Sources', 'Schedule', 'Activity', 'Connections', 'Settings']) {
        await page.getByRole('button', { name: destination, exact: true }).click()
        await expect(page.getByRole('heading', { level: 1, name: destination })).toBeVisible()
        await expect(page.getByText('This computer').first()).toBeVisible()
        expect(await page.evaluate(() => document.documentElement.scrollWidth - document.documentElement.clientWidth)).toBeLessThanOrEqual(1)
        expect(await page.evaluate(() => document.documentElement.scrollHeight - document.documentElement.clientHeight)).toBeLessThanOrEqual(1)
        await expect(page.getByRole('button', { name: 'Exit' })).toBeVisible()
        await expect(page.getByLabel('Appearance')).toBeVisible()
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

test('keeps every semantic action compact, readable, and visibly interactive in both palettes', async ({ page }) => {
  await page.goto('/')
  for (const appearance of ['light', 'dark']) {
    await page.getByLabel('Appearance').selectOption(appearance)
    for (const variant of ['primary', 'secondary', 'subtle', 'affirmative', 'danger']) {
      const button = page.locator('body').locator(`.test-${variant}`)
      await page.evaluate((name) => { const element = document.createElement('button'); element.className = `button button-${name} test-${name}`; element.textContent = name; document.querySelector('main')!.append(element) }, variant)
      const rest = await button.evaluate((element) => { const style = getComputedStyle(element); const background = style.backgroundColor === 'rgba(0, 0, 0, 0)' ? getComputedStyle(document.querySelector('.app')!).backgroundColor : style.backgroundColor; return { color: style.color, background, rawBackground: style.backgroundColor, border: style.borderColor, height: element.getBoundingClientRect().height, padding: Number.parseFloat(style.paddingLeft) } })
      expect(rest.height).toBeGreaterThanOrEqual(24); expect(rest.height).toBeLessThanOrEqual(34); expect(rest.padding).toBeLessThan(16); expect(contrast(rest.color, rest.background)).toBeGreaterThanOrEqual(4.5)
      await button.hover(); await expect.poll(() => button.evaluate((element) => getComputedStyle(element).backgroundColor)).not.toBe(rest.rawBackground)
      await button.focus(); const focus = await button.evaluate((element) => ({ color: getComputedStyle(element).outlineColor, background: getComputedStyle(document.querySelector('.app')!).backgroundColor })); expect(await button.evaluate((element) => getComputedStyle(element).outlineStyle)).not.toBe('none'); expect(contrast(focus.color, focus.background)).toBeGreaterThanOrEqual(3)
      await button.evaluate((element) => element.remove())
    }
  }
})

test('uses one page scroller while keeping shell controls fixed', async ({ page }) => {
  await page.setViewportSize({ width: 800, height: 600 }); await page.goto('/'); await page.getByRole('button', { name: 'Settings' }).click()
  const shellBefore = await page.locator('.rail').boundingBox(); await page.locator('#main-content').evaluate((element) => { element.scrollTop = element.scrollHeight })
  expect((await page.locator('.rail').boundingBox())?.y).toBe(shellBefore?.y)
  expect(await page.locator('#main-content').evaluate((element) => element.scrollTop)).toBeGreaterThan(0)
  await expect(page.getByRole('button', { name: 'Exit' })).toBeVisible()
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
