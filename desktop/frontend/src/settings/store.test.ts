import { act, renderHook, waitFor } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import { useSettings } from './store'
import type { SettingsBridge, SettingsWorkspace } from './model'

const workspace: SettingsWorkspace = { preferences: { version: 1, appearance: 'system', transition: { status: 'not_found', retired: [] } }, preferencePath: '/prefs.json', storage: [], product: { name: 'go-schedule', version: 'dev', publisher: 'ShruggieTech', links: [] }, daemonAvailable: false, loadedAt: '2026-09-07T00:00:00Z' }

function bridge(): SettingsBridge { return { workspace: vi.fn().mockResolvedValue({ action: 'load_settings', outcome: 'accepted', message: 'Loaded.', workspace }), saveAppearance: vi.fn().mockResolvedValue({ action: 'save_appearance', outcome: 'unavailable', message: 'Could not save.' }), restore: vi.fn(), copyStoragePath: vi.fn(), openProductLink: vi.fn() } }

describe('useSettings', () => {
  it('preserves the active workspace after a failed write', async () => {
    const api = bridge(); const { result } = renderHook(() => useSettings(api))
    await waitFor(() => expect(result.current.workspace).toEqual(workspace))
    await act(async () => { await result.current.saveAppearance('dark') })
    expect(result.current.workspace?.preferences.appearance).toBe('system')
    expect(result.current.message).toBe('Could not save.')
  })

  it('suppresses duplicate pending mutations', async () => {
    let release!: () => void
    const pending = new Promise<void>((resolve) => { release = resolve })
    const api = bridge(); api.saveAppearance = vi.fn().mockImplementation(async () => { await pending; return { action: 'save_appearance', outcome: 'accepted', message: 'Saved.', workspace } })
    const { result } = renderHook(() => useSettings(api))
    await waitFor(() => expect(result.current.workspace).toEqual(workspace))
    let first!: Promise<unknown> | undefined
    act(() => { first = result.current.saveAppearance('dark'); void result.current.saveAppearance('light') })
    expect(api.saveAppearance).toHaveBeenCalledOnce()
    release(); await act(async () => { await first })
  })

  it('isolates concurrent copy operations by storage record', async () => {
    let releaseDatabase!: () => void
    let releaseLogs!: () => void
    const databasePending = new Promise<void>((resolve) => { releaseDatabase = resolve })
    const logsPending = new Promise<void>((resolve) => { releaseLogs = resolve })
    const api = bridge(); api.copyStoragePath = vi.fn().mockImplementation(async (id) => { await (id === 'database' ? databasePending : logsPending); return { action: 'copy_storage_path', outcome: 'accepted', message: `${id} copied.`, workspace } })
    const { result } = renderHook(() => useSettings(api))
    await waitFor(() => expect(result.current.workspace).toEqual(workspace))
    let database!: Promise<unknown>; let logs!: Promise<unknown>
    act(() => { database = result.current.copyStoragePath('database'); logs = result.current.copyStoragePath('logs') })
    expect(api.copyStoragePath).toHaveBeenCalledTimes(2)
    expect(result.current.pendingActions).toEqual(new Set(['copy:database', 'copy:logs']))
    releaseLogs(); await act(async () => { await logs })
    expect(result.current.pendingActions).toEqual(new Set(['copy:database']))
    expect(result.current.copyResults.logs).toBe('copied')
    releaseDatabase(); await act(async () => { await database })
    expect(result.current.copyResults.database).toBe('copied')
  })

  it('publishes a distinct event for repeated identical outcomes', async () => {
    const api = bridge(); api.saveAppearance = vi.fn().mockResolvedValue({ action: 'save_appearance', outcome: 'accepted', message: 'Appearance saved.', workspace })
    const { result } = renderHook(() => useSettings(api))
    await waitFor(() => expect(result.current.workspace).toEqual(workspace))
    await act(async () => { await result.current.saveAppearance('dark') })
    const firstID = result.current.status?.id
    await act(async () => { await result.current.saveAppearance('light') })
    expect(result.current.status?.message).toBe('Appearance saved.')
    expect(result.current.status?.id).toBeGreaterThan(firstID ?? 0)
  })

  it('refreshes daemon-backed settings when connection identity changes', async () => {
    const api = bridge(); const { result, rerender } = renderHook(({ token }) => useSettings(api, token), { initialProps: { token: '0:unavailable' } })
    await waitFor(() => expect(api.workspace).toHaveBeenCalledOnce())
    rerender({ token: '1:connected' })
    await waitFor(() => expect(api.workspace).toHaveBeenCalledTimes(2))
    expect(result.current.workspace).toEqual(workspace)
  })
})
