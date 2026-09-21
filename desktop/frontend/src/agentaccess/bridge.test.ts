import { describe, expect, it, vi } from 'vitest'
import { createAgentAccessBridge } from './bridge'

describe('Agent Access bridge', () => {
  it('uses bounded fallbacks outside Wails', async () => {
    const bridge = createAgentAccessBridge({} as Window)
    expect((await bridge.workspace()).outcome).toBe('unavailable')
    expect((await bridge.actions('actor-1')).actions).toEqual([])
  })
  it('passes the enable draft through the native boundary', async () => {
    const enable = vi.fn().mockResolvedValue({ action: 'enable_agent_access', outcome: 'accepted', message: 'copied' })
    const bridge = createAgentAccessBridge({ go: { main: { App: { EnableAgentAccess: enable } } } } as unknown as Window)
    const draft = { clientName: 'Codex', port: 43123, allowedOrigins: [], permission: 'operate' as const }
    await bridge.enable(draft)
    expect(enable).toHaveBeenCalledWith(draft)
  })
  it('passes grant lifecycle requests through the native boundary', async () => {
    const create = vi.fn().mockResolvedValue({ action: 'create_agent_grant', outcome: 'accepted', message: 'copied' })
    const actions = vi.fn().mockResolvedValue({ action: 'load_agent_actions', outcome: 'accepted', message: '', actions: [] })
    const bridge = createAgentAccessBridge({ go: { main: { App: { CreateAgentGrant: create, AgentGrantActions: actions } } } } as unknown as Window)
    const draft = { clientName: 'Build agent', capability: 'observe' as const, duration: '24h' as const }
    await bridge.createGrant(draft)
    await bridge.actions('actor-1')
    expect(create).toHaveBeenCalledWith(draft)
    expect(actions).toHaveBeenCalledWith('actor-1')
  })
})
