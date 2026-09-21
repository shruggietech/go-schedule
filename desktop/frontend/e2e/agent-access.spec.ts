import AxeBuilder from '@axe-core/playwright'
import { expect, test } from '@playwright/test'

test.beforeEach(async ({ page }) => {
  await page.addInitScript(() => {
    const workspace = { daemon: { id: 'daemon-local', name: 'This computer' }, mcpState: 'off', transports: [{ id: 'stdio', name: 'Stdio', state: 'available_on_demand', description: 'Available only while a host launches it.' }, { id: 'localhost_http', name: 'Localhost HTTP', state: 'off', description: 'Optional loopback listener.' }, { id: 'remote_https', name: 'Remote HTTPS', state: 'off', description: 'Separately configured authenticated listener.' }], grants: [], stdioDescription: 'Available on demand when an MCP host launches gosched mcp serve. Stdio opens no network listener.', http: { enabled: false, allowedOrigins: [], requestCount: 0 }, authorities: [{ name: 'Observe', status: 'available', description: 'Read bounded scheduler data.' }, { name: 'Operate', status: 'available', description: 'Run, enable, and disable existing tasks.' }, { name: 'Manage', status: 'available', description: 'Create, update, and delete bounded automation definitions.' }] }
    const settings = { preferences: { version: 1, appearance: 'system', transition: { status: 'not_required', retired: [] } }, preferencePath: '', storage: [], product: { name: 'go-schedule', version: 'test', publisher: 'ShruggieTech', links: [] }, daemonAvailable: true, loadedAt: '2026-09-09T12:00:00Z' }
    const app = {
      Snapshot: async () => ({ generation: 1, revision: 1, state: 'connected', target: { id: 'local', displayName: 'This computer', platform: 'windows', capabilities: ['agent-access'], permissions: ['observe'] }, message: 'Connected.' }),
      RetryConnection: async () => ({ action: 'retry', outcome: 'accepted', message: 'Retrying.' }), Quit: async () => ({ action: 'quit', outcome: 'accepted', message: 'Closing.' }),
      SettingsWorkspace: async () => ({ action: 'load_settings', outcome: 'accepted', message: '', workspace: settings }),
      AgentAccessWorkspace: async () => ({ action: 'load_agent_access', outcome: 'accepted', message: '', workspace }),
      EnableAgentAccess: async () => ({ action: 'enable_agent_access', outcome: 'accepted', message: 'Credential copied once.', workspace }),
      CreateAgentGrant: async (draft: { clientName: string; capability: string }) => ({ action: 'create_agent_grant', outcome: 'accepted', message: 'Enrollment copied once.', workspace: { ...workspace, mcpState: 'active', transports: workspace.transports.map((transport) => transport.id === 'remote_https' ? { ...transport, state: 'active' } : transport), grants: [{ id: 'actor-1', clientName: draft.clientName, daemonId: 'daemon-local', daemonName: 'This computer', capability: draft.capability, capabilityDescription: 'Read bounded scheduler data.', transport: 'remote_https', createdAt: '2026-09-20T12:00:00Z', expiresAt: '2026-09-21T12:00:00Z', state: 'active' }] } }),
      EditAgentGrant: async () => ({ action: 'edit_agent_grant', outcome: 'accepted', message: 'Grant narrowed.', workspace }),
      RevokeAgentGrant: async () => ({ action: 'revoke_agent_grant', outcome: 'accepted', message: 'Grant revoked.', workspace }),
      AgentGrantActions: async () => ({ action: 'load_agent_actions', outcome: 'accepted', message: '', actions: [] }),
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
    await page.getByText('Configure localhost HTTP').click()
    await page.getByLabel('Client name').focus()
    await page.keyboard.press('Tab')
    expect(await page.evaluate(() => document.documentElement.scrollWidth - document.documentElement.clientWidth)).toBeLessThanOrEqual(1)
    const results = await new AxeBuilder({ page }).withTags(['wcag2a', 'wcag2aa', 'wcag22aa']).analyze()
    expect(results.violations.filter((item) => item.impact === 'serious' || item.impact === 'critical')).toEqual([])
  })
}

test('Agent Access supports a keyboard-only bounded grant workflow', async ({ page }) => {
  await page.setViewportSize({ width: 1100, height: 760 })
  await page.goto('/')
  await page.getByRole('button', { name: 'Agent Access', exact: true }).click()
  const create = page.getByRole('button', { name: 'Grant remote access' })
  await create.focus()
  await page.keyboard.press('Enter')
  const dialog = page.getByRole('dialog', { name: 'Grant remote MCP access' })
  await dialog.getByLabel('Client name').fill('Build agent')
  await dialog.getByLabel('Authority').selectOption('observe')
  await dialog.getByLabel('Access duration').selectOption('24h')
  await dialog.getByRole('button', { name: 'Create and copy enrollment' }).click()
  await expect(page.getByRole('heading', { name: 'Build agent' })).toBeVisible()
  await expect(page.getByText('MCP On')).toBeVisible()
  await page.getByRole('button', { name: 'Recent actions' }).click()
  await expect(page.getByRole('dialog', { name: 'Build agent recent actions' })).toBeVisible()
  await page.keyboard.press('Escape')
  await expect(page.getByRole('button', { name: 'Recent actions' })).toBeFocused()
})
