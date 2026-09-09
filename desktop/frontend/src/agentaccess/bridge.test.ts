import { describe, expect, it, vi } from 'vitest'
import { createAgentAccessBridge } from './bridge'

describe('Agent Access bridge', () => {
  it('uses bounded fallbacks outside Wails', async () => { expect((await createAgentAccessBridge({} as Window).workspace()).outcome).toBe('unavailable') })
  it('passes the enable draft through the native boundary', async () => {
    const enable = vi.fn().mockResolvedValue({ action: 'enable_agent_access', outcome: 'accepted', message: 'copied' })
    const bridge = createAgentAccessBridge({ go: { main: { App: { EnableAgentAccess: enable } } } } as unknown as Window)
    const draft = { clientName: 'Codex', port: 43123, allowedOrigins: [] }
    await bridge.enable(draft)
    expect(enable).toHaveBeenCalledWith(draft)
  })
})
