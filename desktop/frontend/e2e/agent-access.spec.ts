import AxeBuilder from '@axe-core/playwright'
import { expect, test } from '@playwright/test'

test.beforeEach(async ({ page }) => {
  await page.addInitScript(() => {
    const workspace = { stdioDescription: 'Available on demand when an MCP host launches gosched mcp serve. Stdio opens no network listener.', http: { enabled: false, allowedOrigins: [], requestCount: 0 }, authorities: [{ name: 'Observe', status: 'available', description: 'Read bounded scheduler data.' }, { name: 'Operate', status: 'future', description: 'Unavailable. Agents cannot run or change scheduled work.' }, { name: 'Manage', status: 'future', description: 'Unavailable. Agents cannot change configuration or access controls.' }] }
    const settings = { preferences: { version: 1, appearance: 'system', transition: { status: 'not_required', retired: [] } }, preferencePath: '', storage: [], product: { name: 'go-schedule', version: 'test', publisher: 'ShruggieTech', links: [] }, daemonAvailable: true, loadedAt: '2026-09-09T12:00:00Z' }
    const app = {
      Snapshot: async () => ({ generation: 1, revision: 1, state: 'connected', target: { id: 'local', displayName: 'This computer', platform: 'windows', capabilities: ['agent-access'], permissions: ['observe'] }, message: 'Connected.' }),
      RetryConnection: async () => ({ action: 'retry', outcome: 'accepted', message: 'Retrying.' }), Quit: async () => ({ action: 'quit', outcome: 'accepted', message: 'Closing.' }),
      SettingsWorkspace: async () => ({ action: 'load_settings', outcome: 'accepted', message: '', workspace: settings }),
      AgentAccessWorkspace: async () => ({ action: 'load_agent_access', outcome: 'accepted', message: '', workspace }),
      EnableAgentAccess: async () => ({ action: 'enable_agent_access', outcome: 'accepted', message: 'Credential copied once.', workspace }),
      OpenAgentAccessGuide: async () => ({ action: 'open_agent_access_guide', outcome: 'accepted', message: 'Guide opened.' }),
    }
    Object.assign(window, { go: { main: { App: app } }, runtime: { EventsOn: () => () => undefined } })
  })
})

for (const zoom of [0.8, 1, 1.5, 2]) {
  test(`Agent Access is accessible at ${zoom * 100} percent zoom`, async ({ page }) => {
    await page.setViewportSize({ width: 900, height: 650 })
    await page.goto('/')
    await page.evaluate((value) => { document.documentElement.style.zoom = String(value) }, zoom)
    await page.getByRole('button', { name: 'Agent Access', exact: true }).click()
    await expect(page.getByRole('heading', { level: 1, name: 'Agent Access' })).toBeVisible()
    await page.getByLabel('Client name').focus()
    await page.keyboard.press('Tab')
    expect(await page.evaluate(() => document.documentElement.scrollWidth - document.documentElement.clientWidth)).toBeLessThanOrEqual(1)
    const results = await new AxeBuilder({ page }).withTags(['wcag2a', 'wcag2aa', 'wcag22aa']).analyze()
    expect(results.violations.filter((item) => item.impact === 'serious' || item.impact === 'critical')).toEqual([])
  })
}
