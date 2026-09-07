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
})
