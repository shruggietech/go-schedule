import AxeBuilder from '@axe-core/playwright'
import { expect, test } from '@playwright/test'

test.beforeEach(async ({ page }) => {
  await page.addInitScript(() => {
    const state = {
      appearance: 'dark',
      copied: [] as string[],
      opened: [] as string[],
      workspace: {
        preferences: { version: 1, appearance: 'dark', transition: { status: 'migrated', retired: ['appearance.font', 'appearance.scroll_sensitivity'] } },
        preferencePath: 'C:\\Users\\Ada\\AppData\\Roaming\\go-schedule\\desktop\\preferences.json',
        storage: [
          { id: 'machine-data', label: 'Machine data', path: 'C:\\ProgramData\\goschedule', owner: 'go-schedule', scope: 'machine', existence: 'present', normalRemoval: 'Preserved by normal software removal', explicitWipe: 'Removed by an explicit data wipe', copyable: true },
          { id: 'logs', label: 'Logs', owner: 'go-schedule', scope: 'machine', existence: 'unavailable', normalRemoval: 'Preserved by normal software removal', explicitWipe: 'Removed by an explicit data wipe', copyable: false },
        ],
        product: { name: 'go-schedule', version: '1.2.0', publisher: 'ShruggieTech', links: [{ key: 'source', label: 'Source repository', destination: 'https://github.com/shruggietech/go-schedule' }, { key: 'documentation', label: 'Documentation', destination: 'https://shruggietech.github.io/go-schedule/' }] },
        daemonAvailable: false,
        loadedAt: '2026-09-07T12:00:00Z',
      },
    }
    const accepted = (action: string, message: string) => ({ action, outcome: 'accepted', message, workspace: structuredClone(state.workspace) })
    const app = {
      Snapshot: async () => ({ generation: 1, revision: 1, state: 'unavailable', target: { id: 'local', displayName: 'This computer', platform: 'windows', capabilities: [], permissions: [] }, message: 'The local scheduler service is unavailable.', action: 'Start the service, then try again.' }),
      RetryConnection: async () => ({ action: 'retry', outcome: 'accepted', message: 'Trying the local scheduler service again.' }),
      Quit: async () => ({ action: 'quit', outcome: 'accepted', message: 'Closing go-schedule.' }),
      SettingsWorkspace: async () => accepted('load_settings', 'Desktop settings loaded.'),
      SaveAppearance: async (appearance: string) => { state.appearance = appearance; state.workspace.preferences.appearance = appearance; return accepted('save_appearance', 'Appearance saved.') },
      RestoreDesktopPreferences: async () => { state.appearance = 'system'; state.workspace.preferences.appearance = 'system'; return accepted('restore_preferences', 'Desktop preferences restored.') },
      CopyStoragePath: async (id: string) => { state.copied.push(id); return accepted('copy_storage_path', 'Machine data path copied.') },
      OpenProductLink: async (key: string) => { state.opened.push(key); return { action: 'open_product_link', outcome: 'accepted', message: 'Product link opened.' } },
    }
    Object.assign(window, { go: { main: { App: app } }, runtime: { EventsOn: () => () => undefined }, settingsTestState: state })
  })
})

for (const zoom of [0.8, 1, 1.5, 2]) {
  test(`Settings remains keyboard accessible at ${zoom * 100} percent zoom`, async ({ page }) => {
    await page.setViewportSize({ width: 900, height: 650 })
    await page.goto('/')
    await page.evaluate((value) => { document.documentElement.style.zoom = String(value) }, zoom)
    await page.getByRole('button', { name: 'Settings', exact: true }).click()
    await expect(page.getByRole('heading', { level: 1, name: 'Settings' })).toBeVisible()
    await expect(page.getByRole('heading', { name: 'Appearance migrated' })).toBeVisible()
    expect(await page.evaluate(() => document.documentElement.scrollWidth - document.documentElement.clientWidth)).toBeLessThanOrEqual(1)
    const results = await new AxeBuilder({ page }).withTags(['wcag2a', 'wcag2aa', 'wcag22aa']).analyze()
    expect(results.violations.filter((item) => item.impact === 'serious' || item.impact === 'critical')).toEqual([])
  })
}

test('persists appearance and routes native actions through stable identifiers', async ({ page }) => {
  await page.goto('/')
  await page.getByRole('button', { name: 'Settings', exact: true }).click()
  await page.getByRole('combobox', { name: 'Appearance' }).first().selectOption('light')
  await expect(page.locator('.app')).toHaveAttribute('data-appearance', 'light')
  await page.getByRole('button', { name: 'Copy path' }).click()
  await page.getByRole('button', { name: 'Source repository' }).click()
  const actions = await page.evaluate(() => (window as unknown as { settingsTestState: { copied: string[]; opened: string[] } }).settingsTestState)
  expect(actions.copied).toEqual(['machine-data'])
  expect(actions.opened).toEqual(['source'])
})

test('connection recovery remains inline and duplicate-safe', async ({ page }) => {
  await page.goto('/')
  await page.getByRole('button', { name: 'Connections', exact: true }).click()
  await expect(page.getByRole('heading', { level: 1, name: 'Connections' })).toBeVisible()
  await expect(page.getByText('Start or repair the local scheduler service, then try again.')).toBeVisible()
  await page.getByRole('button', { name: 'Try again' }).last().click()
  await expect(page.getByText(/Trying the local scheduler service again/).first()).toBeVisible()
  await expect(page.getByRole('dialog')).toHaveCount(0)
})
