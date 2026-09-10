import { act, renderHook, waitFor } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import type { OperationResult, TaskBridge } from './model'
import { useTaskWorkspace } from './store'

const result = (ids: string[]): OperationResult => ({ action: 'load', outcome: 'accepted', message: 'Loaded.', workspace: { loadedAt: 'now', groups: [], tasks: ids.map((id) => ({ id, name: id, groupId: '', groupPath: 'Not assigned', commandConfigured: true, declaredEnabled: false, effectiveState: 'manual_only', effectiveReason: '', lifecycle: 'active', timezone: 'UTC', scheduleSummary: 'Manual only', policySummary: '', nextRuns: [], updatedAt: 'now' })) } })

describe('task workspace authority', () => {
  it('coalesces relevant events and preserves a still-valid stable selection', async () => {
    let listener: ((event: { kind: string }) => void) | undefined
    const workspace = vi.fn().mockResolvedValueOnce(result(['one', 'two'])).mockResolvedValue(result(['one', 'two', 'three']))
    const bridge = { workspace, subscribe: (next: typeof listener) => { listener = next; return () => undefined } } as unknown as TaskBridge
    const { result: hook } = renderHook(() => useTaskWorkspace(bridge))
    await act(async () => { await Promise.resolve() })
    await waitFor(() => expect(hook.current.selected).toBe('one'))
    act(() => hook.current.setSelected('two'))
    act(() => { listener?.({ kind: 'task.updated' }); listener?.({ kind: 'group.updated' }) })
    await act(async () => { await new Promise((resolve) => setTimeout(resolve, 100)) })
    expect(workspace).toHaveBeenCalledTimes(2); expect(hook.current.selected).toBe('two')
  })

  it('ignores an older workspace response that finishes last', async () => {
    let resolveFirst: ((value: OperationResult) => void) | undefined; let resolveSecond: ((value: OperationResult) => void) | undefined
    const workspace = vi.fn().mockImplementationOnce(() => new Promise<OperationResult>((resolve) => { resolveFirst = resolve })).mockImplementationOnce(() => new Promise<OperationResult>((resolve) => { resolveSecond = resolve }))
    const bridge = { workspace } as unknown as TaskBridge
    const { result: hook } = renderHook(() => useTaskWorkspace(bridge))
    await waitFor(() => expect(workspace).toHaveBeenCalledTimes(1)); act(() => { void hook.current.load() }); await waitFor(() => expect(workspace).toHaveBeenCalledTimes(2))
    await act(async () => { resolveSecond?.(result(['new'])); await Promise.resolve() }); expect(hook.current.selected).toBe('new')
    await act(async () => { resolveFirst?.(result(['old'])); await Promise.resolve() }); expect(hook.current.selected).toBe('new')
  })

  it('retains the last complete workspace while unavailable and refreshes after recovery', async () => {
    const workspace = vi.fn().mockResolvedValueOnce(result(['known'])).mockResolvedValueOnce(result(['current']))
    const bridge = { workspace } as unknown as TaskBridge
    const { result: hook, rerender } = renderHook(({ available, token }) => useTaskWorkspace(bridge, available, token), { initialProps: { available: true, token: 1 } })
    await waitFor(() => expect(hook.current.selected).toBe('known'))
    rerender({ available: false, token: 0 })
    expect(hook.current.workspace?.tasks[0].id).toBe('known')
    rerender({ available: true, token: 2 })
    await waitFor(() => expect(hook.current.selected).toBe('current'))
    expect(workspace).toHaveBeenCalledTimes(2)
  })

  it('reports unavailable instead of loading forever when the initial read is blocked', async () => {
    const bridge = { workspace: vi.fn() } as unknown as TaskBridge
    const { result: hook } = renderHook(() => useTaskWorkspace(bridge, false, 0))
    await waitFor(() => expect(hook.current.status?.outcome).toBe('unavailable'))
    expect(hook.current.workspace).toBeUndefined()
    expect(bridge.workspace).not.toHaveBeenCalled()
  })
})
