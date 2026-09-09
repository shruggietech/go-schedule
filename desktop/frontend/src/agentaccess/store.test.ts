import { act, renderHook, waitFor } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import type { AgentAccessBridge, AgentAccessWorkspace } from './model'
import { useAgentAccess } from './store'

const workspace: AgentAccessWorkspace = { stdioDescription: 'On demand.', http: { enabled: false, allowedOrigins: [], requestCount: 0 }, authorities: [] }
const bridge = (): AgentAccessBridge => ({ workspace: vi.fn().mockResolvedValue({ action: 'load_agent_access', outcome: 'accepted', message: '', workspace }), enable: vi.fn(), rotate: vi.fn(), revoke: vi.fn(), openGuide: vi.fn() })

describe('useAgentAccess', () => {
  it('rejects stale loads when connection identity changes', async () => {
    let resolveFirst!: (value: unknown) => void
    const api = bridge()
    api.workspace = vi.fn().mockImplementationOnce(() => new Promise((resolve) => { resolveFirst = resolve })).mockResolvedValueOnce({ action: 'load_agent_access', outcome: 'accepted', message: 'current', workspace })
    const { result, rerender } = renderHook(({ token }) => useAgentAccess(api, true, token), { initialProps: { token: 1 } })
    await waitFor(() => expect(api.workspace).toHaveBeenCalledOnce())
    rerender({ token: 2 })
    await waitFor(() => expect(api.workspace).toHaveBeenCalledTimes(2))
    act(() => resolveFirst({ action: 'load_agent_access', outcome: 'accepted', message: 'stale', workspace: { ...workspace, stdioDescription: 'stale' } }))
    await waitFor(() => expect(result.current.message).toBe('current'))
    expect(result.current.workspace?.stdioDescription).toBe('On demand.')
  })

  it('suppresses duplicate lifecycle actions', async () => {
    let release!: () => void
    const api = bridge()
    api.rotate = vi.fn().mockImplementation(() => new Promise((resolve) => { release = () => resolve({ action: 'rotate_agent_access', outcome: 'accepted', message: 'rotated', workspace }) }))
    const { result } = renderHook(() => useAgentAccess(api, true, 1))
    await waitFor(() => expect(result.current.workspace).toEqual(workspace))
    act(() => { void result.current.rotate(); void result.current.rotate() })
    expect(api.rotate).toHaveBeenCalledOnce()
    act(() => release())
    await waitFor(() => expect(result.current.pending).toBe(false))
  })

  it('turns a rejected initial call into a bounded unavailable result', async () => {
    const api = bridge()
    api.workspace = vi.fn().mockRejectedValue(new Error('native detail'))
    const { result } = renderHook(() => useAgentAccess(api, true, 1))
    await waitFor(() => expect(result.current.message).toMatch(/unavailable/i))
    expect(result.current.message).not.toContain('native detail')
  })

  it('refreshes evidence while localhost access is active', async () => {
    vi.useFakeTimers()
    const active = { ...workspace, http: { ...workspace.http, enabled: true } }
    const api = bridge()
    api.workspace = vi.fn().mockResolvedValue({ action: 'load_agent_access', outcome: 'accepted', message: '', workspace: active })
    const { result } = renderHook(() => useAgentAccess(api, true, 1))
    await act(async () => { await vi.runOnlyPendingTimersAsync() })
    expect(result.current.workspace?.http.enabled).toBe(true)
    const calls = vi.mocked(api.workspace).mock.calls.length
    await act(async () => { await vi.advanceTimersByTimeAsync(5000) })
    expect(api.workspace).toHaveBeenCalledTimes(calls + 1)
    vi.useRealTimers()
  })
})
