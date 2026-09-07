import { act, renderHook, waitFor } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import { useAutomationWorkspace } from './store'
import type { AutomationBridge, OperationResult } from './model'

describe('useAutomationWorkspace', () => {
  it('ignores an older load that resolves after a newer mutation', async () => {
    let resolve!: (result: OperationResult) => void
    const older = new Promise<OperationResult>((done) => { resolve = done })
    const bridge = { workspace: vi.fn().mockReturnValue(older), subscribe: () => () => undefined } as unknown as AutomationBridge
    const { result } = renderHook(() => useAutomationWorkspace(bridge))
    const newer: OperationResult = { action: 'save_chain', outcome: 'accepted', message: 'Saved.', workspace: { tasks: [], chains: [], triggers: [], triggerSets: [], watchers: [], loadedAt: 'newer' } }
    act(() => result.current.accept(newer)); resolve({ action: 'load', outcome: 'accepted', message: 'Old.', workspace: { ...newer.workspace!, loadedAt: 'older' } })
    await waitFor(() => expect(result.current.workspace?.loadedAt).toBe('newer'))
  })

  it('discards a save refresh when a newer load has started', async () => {
    let resolveSaveRefresh!: (result: OperationResult) => void; let resolveNewerLoad!: (result: OperationResult) => void
    const saveRefresh = new Promise<OperationResult>((done) => { resolveSaveRefresh = done }); const newerLoad = new Promise<OperationResult>((done) => { resolveNewerLoad = done })
    const initial: OperationResult = { action: 'load', outcome: 'accepted', message: 'Initial.', workspace: { tasks: [], chains: [], triggers: [], triggerSets: [], watchers: [], loadedAt: 'initial' } }
    const bridge = { workspace: vi.fn().mockResolvedValueOnce(initial).mockReturnValueOnce(saveRefresh).mockReturnValueOnce(newerLoad), subscribe: () => () => undefined } as unknown as AutomationBridge
    const { result } = renderHook(() => useAutomationWorkspace(bridge)); await waitFor(() => expect(result.current.workspace?.loadedAt).toBe('initial'))
    const saved: OperationResult = { action: 'save_trigger', outcome: 'accepted', message: 'Saved.' }; act(() => { void result.current.refreshAfter(saved) }); act(() => { void result.current.load() })
    resolveNewerLoad({ ...initial, workspace: { ...initial.workspace!, loadedAt: 'newer' } }); await waitFor(() => expect(result.current.workspace?.loadedAt).toBe('newer'))
    resolveSaveRefresh({ ...initial, workspace: { ...initial.workspace!, loadedAt: 'older-save-refresh' } }); await waitFor(() => expect(result.current.workspace?.loadedAt).toBe('newer'))
  })
})
