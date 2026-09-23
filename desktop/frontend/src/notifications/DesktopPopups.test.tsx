import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'
import { DesktopPopups } from './DesktopPopups'
import type { DesktopBridge } from '../connection/model'
import type { PopupPreferences, SettingsBridge, SettingsWorkspace } from '../settings/model'

const initial: PopupPreferences = { enabled: false, conditions: ['failure', 'alert'], severities: ['warning', 'error'], daemonIds: [] }
const workspace: SettingsWorkspace = { preferences: { version: 1, appearance: 'system', transition: { status: 'not_found', retired: [] }, popups: initial }, preferencePath: '', storage: [], product: { name: '', version: '', publisher: '', links: [] }, daemonAvailable: false, loadedAt: '' }
const desktop = { allSystems: async () => ({ generation: 1, startedAt: '', completedAt: '', complete: true, observations: [{ registration: { key: 'local', kind: 'local' as const, label: 'This computer', daemonId: 'daemon-local' }, state: 'connected' as const, stale: false }, { registration: { key: 'profile-a', kind: 'remote' as const, label: 'Office', daemonId: 'daemon-office' }, state: 'connected' as const, stale: false }] }) } as DesktopBridge

function settings(savePopups: SettingsBridge['savePopups'] = vi.fn(async () => ({ action: 'save_popups', outcome: 'accepted' as const, message: 'Desktop popup preferences saved.', workspace }))): SettingsBridge {
  return { workspace: async () => ({ action: 'load_settings', outcome: 'accepted', message: '', workspace }), saveAppearance: async () => ({ action: '', outcome: 'accepted', message: '' }), restore: async () => ({ action: '', outcome: 'accepted', message: '' }), copyStoragePath: async () => ({ action: '', outcome: 'accepted', message: '' }), openProductLink: async () => ({ action: '', outcome: 'accepted', message: '' }), popupStatus: async () => ({ available: true, authorized: true, message: 'Available while running.' }), savePopups }
}

describe('DesktopPopups', () => {
  it('starts muted, states the closed-app limit, and saves selected filters', async () => {
    const save = vi.fn(async () => ({ action: 'save_popups', outcome: 'accepted' as const, message: 'Desktop popup preferences saved.', workspace }))
    render(<DesktopPopups settings={settings(save)} desktop={desktop} />)
    const enable = await screen.findByRole('checkbox', { name: 'Enable desktop popups' })
    expect(enable).not.toBeChecked()
    expect(screen.getByText(/No popups are sent while the desktop app is closed/)).toBeInTheDocument()
    await userEvent.click(enable)
    await userEvent.click(await screen.findByRole('checkbox', { name: 'Office' }))
    await userEvent.click(screen.getByRole('button', { name: 'Save popup choices' }))
    await waitFor(() => expect(save).toHaveBeenCalledWith({ enabled: true, conditions: ['failure', 'alert'], severities: ['warning', 'error'], daemonIds: ['daemon-office'] }))
  })

  it('does not offer enabling when the OS facility is absent', async () => {
    const bridge = settings()
    bridge.popupStatus = async () => ({ available: false, authorized: false, message: 'No notification service.' })
    render(<DesktopPopups settings={bridge} desktop={desktop} />)
    expect(await screen.findByRole('checkbox', { name: 'Enable desktop popups' })).toBeDisabled()
  })

  it('still permits silencing a saved preference if OS support disappears', async () => {
    const bridge = settings()
    bridge.workspace = async () => ({ action: 'load_settings', outcome: 'accepted', message: '', workspace: { ...workspace, preferences: { ...workspace.preferences, popups: { ...initial, enabled: true } } } })
    bridge.popupStatus = async () => ({ available: false, authorized: false, message: 'No notification service.' })
    const save = vi.fn(async () => ({ action: 'save_popups', outcome: 'accepted' as const, message: 'Desktop popup preferences saved.' }))
    bridge.savePopups = save
    render(<DesktopPopups settings={bridge} desktop={desktop} />)
    const enable = await screen.findByRole('checkbox', { name: 'Enable desktop popups' })
    await waitFor(() => expect(enable).toBeChecked())
    expect(enable).not.toBeDisabled()
    await userEvent.click(enable)
    await userEvent.click(screen.getByRole('button', { name: 'Save popup choices' }))
    await waitFor(() => expect(save).toHaveBeenCalledWith(expect.objectContaining({ enabled: false })))
  })
})
