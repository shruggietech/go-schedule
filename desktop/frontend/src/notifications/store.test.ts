import { act, renderHook, waitFor } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import { useNotifications } from './store'
import type { NotificationBridge, NotificationResult, NotificationWorkspace } from './model'

const snapshot = (name: string): NotificationWorkspace => ({ channels: [{ id: name, name, kind: 'webhook', endpointSummary: 'https://safe/...', hasAuthorization: false, enabled: true, updatedAt: '' }], tasks: [], groups: [], deliveries: [], loadedAt: '' })
const result = (workspace?: NotificationWorkspace, outcome: NotificationResult['outcome'] = 'accepted'): NotificationResult => ({ action: 'load_notifications', outcome, message: outcome === 'accepted' ? 'ok' : 'safe error', workspace })
const deferred = <T,>() => { let resolve!: (value: T) => void; const promise = new Promise<T>((done) => { resolve = done }); return { promise, resolve } }
const base = (workspace = snapshot('initial')): NotificationBridge => ({ workspace: () => Promise.resolve(result(workspace)), saveChannel: () => Promise.resolve(result(workspace)), setChannelEnabled: () => Promise.resolve(result(workspace)), testChannel: () => Promise.resolve(result(workspace)), deleteChannel: () => Promise.resolve(result(workspace)), policy: () => Promise.resolve({ action: 'load_notification_policy', outcome: 'accepted', message: 'ok' }), savePolicy: () => Promise.resolve({ action: 'save_notification_policy', outcome: 'accepted', message: 'ok' }) })

describe('useNotifications', () => {
  it('discards stale workspace completions and preserves a complete snapshot on failure', async () => {
    const old = deferred<NotificationResult>(); const fresh = deferred<NotificationResult>(); const workspace = vi.fn().mockReturnValueOnce(Promise.resolve(result(snapshot('initial')))).mockReturnValueOnce(old.promise).mockReturnValueOnce(fresh.promise).mockResolvedValue(result(undefined, 'unavailable')); const testBridge = { ...base(), workspace }
    const { result: hook } = renderHook(() => useNotifications(testBridge, true, 1)); await waitFor(() => expect(hook.current.workspace?.channels[0].name).toBe('initial'))
    const first = hook.current.load(); const second = hook.current.load(); fresh.resolve(result(snapshot('fresh'))); await act(async () => { await second }); old.resolve(result(snapshot('old'))); await act(async () => { await first })
    expect(hook.current.workspace?.channels[0].name).toBe('fresh'); await act(async () => { await hook.current.load() }); expect(hook.current.workspace?.channels[0].name).toBe('fresh'); expect(hook.current.status?.outcome).toBe('unavailable')
  })

  it('suppresses duplicate mutations before React can rerender', async () => {
    const pending = deferred<NotificationResult>(); const testChannel = vi.fn(() => pending.promise); const testBridge = { ...base(), testChannel }; const { result: hook } = renderHook(() => useNotifications(testBridge, true, 1)); await waitFor(() => expect(hook.current.workspace).toBeDefined())
    act(() => { void hook.current.testChannel('c1'); void hook.current.testChannel('c1') }); expect(testChannel).toHaveBeenCalledTimes(1); pending.resolve(result(snapshot('done'))); await waitFor(() => expect(hook.current.pending).toBe(false))
  })
})
