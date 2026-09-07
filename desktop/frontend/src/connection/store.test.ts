import { act, renderHook, waitFor } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import { unavailableSnapshot } from './bridge'
import { acceptSnapshot, useConnection } from './store'
import type { ConnectionSnapshot, DesktopBridge, DesktopEvent } from './model'

const connected: ConnectionSnapshot = { generation: 2, revision: 4, state: 'connected', target: { id: 'local', displayName: 'This computer', platform: 'windows', version: '1.2.0', capabilities: ['tasks'], permissions: ['read'] }, message: 'Scheduler service is available.' }

function fakeBridge(initial = connected) {
  let listener: ((event: DesktopEvent) => void) | undefined
  const bridge: DesktopBridge = {
    snapshot: vi.fn().mockResolvedValue(initial), retry: vi.fn().mockResolvedValue({ action: 'retry', outcome: 'accepted', message: 'Trying again.' }), quit: vi.fn().mockResolvedValue({ action: 'quit', outcome: 'accepted', message: 'Closing.' }),
    subscribe: (next) => { listener = next; return () => { listener = undefined } },
  }
  return { bridge, emit: (event: DesktopEvent) => listener?.(event) }
}

describe('connection store', () => {
  it('rejects stale snapshots', () => { expect(acceptSnapshot(connected, { ...unavailableSnapshot, generation: 1 })).toBe(connected) })

  it('rejects an older snapshot from the same generation', () => {
    expect(acceptSnapshot(connected, { ...connected, revision: 3, state: 'connecting' })).toBe(connected)
  })

  it('loads native state, accepts current events, and announces recovery', async () => {
    const fake = fakeBridge()
    const { result } = renderHook(() => useConnection(fake.bridge))
    await waitFor(() => expect(result.current.snapshot.state).toBe('connected'))
    act(() => fake.emit({ id: '3', kind: 'connection.changed', message: 'Recovering.', generation: 3, occurredAt: 'now', snapshot: { ...connected, generation: 3, revision: 5, state: 'recovering' } }))
    expect(result.current.snapshot.state).toBe('recovering')
    expect(result.current.announcement).toBe('Recovering.')
  })

  it('keeps a truthful browser fallback and exposes manual retry', async () => {
    const fake = fakeBridge(unavailableSnapshot)
    const { result } = renderHook(() => useConnection(fake.bridge))
    await waitFor(() => expect(result.current.snapshot.state).toBe('unavailable'))
    await act(() => result.current.retry())
    expect(fake.bridge.retry).toHaveBeenCalledOnce()
    expect(result.current.announcement).toBe('Trying again.')
  })
})
