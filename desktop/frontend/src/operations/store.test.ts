import { act, renderHook, waitFor } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import type { OperationResult, OperationsBridge } from './model'
import { useActivity, useSchedule } from './store'

const schedule = (loadedAt: string): OperationResult => ({ action: 'load_schedule', outcome: 'accepted', message: 'Loaded.', schedule: { from: '2026-09-07T00:00:00Z', to: '2026-09-14T00:00:00Z', loadedAt, occurrences: [] } })
const activity = (loadedAt: string): OperationResult => ({ action: 'load_activity', outcome: 'accepted', message: 'Loaded.', activity: { runs: [], logs: [], alerts: [], logPath: '', loadedAt } })

describe('operations stores', () => {
  it('discards an older Schedule response after a range change', async () => {
    let resolveOlder!: (value: OperationResult) => void; let resolveNewer!: (value: OperationResult) => void
    const older = new Promise<OperationResult>((done) => { resolveOlder = done }); const newer = new Promise<OperationResult>((done) => { resolveNewer = done })
    const bridge = { scheduleWindow: vi.fn().mockReturnValueOnce(older).mockReturnValueOnce(newer), subscribe: () => () => undefined } as unknown as OperationsBridge
    const { result, rerender } = renderHook(({ days }) => useSchedule(bridge, days, true, 1), { initialProps: { days: 7 } }); rerender({ days: 30 })
    act(() => resolveNewer(schedule('newer'))); await waitFor(() => expect(result.current.snapshot?.loadedAt).toBe('newer')); act(() => resolveOlder(schedule('older'))); await waitFor(() => expect(result.current.snapshot?.loadedAt).toBe('newer'))
  })

  it('preserves the last complete Activity workspace after a failed refresh', async () => {
    const bridge = { activityWorkspace: vi.fn().mockResolvedValueOnce(activity('complete')).mockResolvedValueOnce({ action: 'load_activity', outcome: 'unavailable', message: 'Unavailable.' }), subscribe: () => () => undefined } as unknown as OperationsBridge
    const { result } = renderHook(() => useActivity(bridge, true, 1)); await waitFor(() => expect(result.current.workspace?.loadedAt).toBe('complete')); await act(async () => { await result.current.load() }); expect(result.current.workspace?.loadedAt).toBe('complete'); expect(result.current.status?.outcome).toBe('unavailable')
  })

  it('retains schedule data while unavailable and refreshes on the recovered generation', async () => {
    const bridge = { scheduleWindow: vi.fn().mockResolvedValueOnce(schedule('known')).mockResolvedValueOnce(schedule('current')), subscribe: () => () => undefined } as unknown as OperationsBridge
    const { result, rerender } = renderHook(({ available, token }) => useSchedule(bridge, 7, available, token), { initialProps: { available: true, token: 1 } })
    await waitFor(() => expect(result.current.snapshot?.loadedAt).toBe('known'))
    rerender({ available: false, token: 0 })
    expect(result.current.snapshot?.loadedAt).toBe('known')
    rerender({ available: true, token: 2 })
    await waitFor(() => expect(result.current.snapshot?.loadedAt).toBe('current'))
    expect(bridge.scheduleWindow).toHaveBeenCalledTimes(2)
  })
})
