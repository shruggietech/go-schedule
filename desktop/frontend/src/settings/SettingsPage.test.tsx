import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'
import { SettingsPage } from './SettingsPage'
import type { SettingsWorkspace } from './model'

const workspace: SettingsWorkspace = {
  preferences: { version: 1, appearance: 'dark', transition: { status: 'migrated', retired: ['appearance.font', 'appearance.scroll_sensitivity'] } },
  preferencePath: '/home/ada/.config/go-schedule/desktop/preferences.json',
  storage: [
    { id: 'database', label: 'Task database', path: '/var/lib/goschedule/tasks.db', owner: 'go-schedule', scope: 'machine', existence: 'present', normalRemoval: 'Preserved by normal software removal', explicitWipe: 'No built-in data wipe on this platform', copyable: true },
    { id: 'logs', label: 'Logs', owner: 'external', scope: 'external', existence: 'unavailable', normalRemoval: 'Preserved by normal software removal', explicitWipe: 'Preserved because it is outside the application-owned data root', copyable: false },
  ],
  product: { name: 'go-schedule', version: '1.2.0', publisher: 'ShruggieTech', links: [{ key: 'source', label: 'Source repository', destination: 'https://github.com/shruggietech/go-schedule' }] },
  daemonAvailable: false,
  loadedAt: '2026-09-07T12:00:00Z',
}

function renderPage(overrides: Partial<Parameters<typeof SettingsPage>[0]> = {}) {
  const props = { workspace, message: '', pending: false, onAppearance: vi.fn(), onRestore: vi.fn(), onCopy: vi.fn(), onOpen: vi.fn(), onConnections: vi.fn(), ...overrides }
  render(<SettingsPage {...props} />)
  return props
}

describe('SettingsPage', () => {
  it('explains migration, retired controls, storage truth, and offline state', () => {
    renderPage()
    expect(screen.getByRole('heading', { name: 'Appearance migrated' })).toBeVisible()
    expect(screen.getByText(/font selection and scroll sensitivity were retired/i)).toBeVisible()
    expect(screen.getByText('/var/lib/goschedule/tasks.db')).toBeVisible()
    expect(screen.getByText('Preserved because it is outside the application-owned data root')).toBeVisible()
    expect(screen.getByText('Daemon storage details are unavailable')).toBeVisible()
    expect(screen.getByRole('heading', { name: 'go-schedule' })).toBeVisible()
  })

  it('routes appearance, restore, copy, link, and recovery actions by identifiers', async () => {
    const user = userEvent.setup(); const props = renderPage()
    await user.selectOptions(screen.getByRole('combobox', { name: 'Appearance' }), 'light')
    await user.click(screen.getByRole('button', { name: 'Restore desktop defaults' }))
    await user.click(screen.getByRole('button', { name: 'Copy path' }))
    await user.click(screen.getByRole('button', { name: 'Source repository' }))
    await user.click(screen.getByRole('button', { name: 'Open Connections' }))
    expect(props.onAppearance).toHaveBeenCalledWith('light')
    expect(props.onRestore).toHaveBeenCalledOnce()
    expect(props.onCopy).toHaveBeenCalledWith('database')
    expect(props.onOpen).toHaveBeenCalledWith('source')
    expect(props.onConnections).toHaveBeenCalledOnce()
  })

  it('shows a useful local failure without inventing content', () => {
    renderPage({ workspace: undefined, message: 'Desktop settings are unavailable.' })
    expect(screen.getByRole('heading', { name: 'Desktop settings are unavailable' })).toBeVisible()
    expect(screen.queryByRole('button', { name: 'Copy path' })).not.toBeInTheDocument()
  })
})
