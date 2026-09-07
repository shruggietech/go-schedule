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
})
