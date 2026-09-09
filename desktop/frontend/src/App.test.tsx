import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'
import { App } from './App'
import type { ConnectionSnapshot, DesktopBridge } from './connection/model'
import type { TaskBridge } from './tasks/model'
import type { OperationsBridge } from './operations/model'
import type { SettingsBridge, SettingsWorkspace } from './settings/model'
import type { NotificationBridge } from './notifications/model'

const connected: ConnectionSnapshot = { generation: 1, revision: 2, state: 'connected', target: { id: 'local', displayName: 'This computer', platform: 'linux', version: '1.2.0', capabilities: ['tasks'], permissions: ['read'] }, message: 'Scheduler service is available.' }
const bridge: DesktopBridge = { snapshot: vi.fn().mockResolvedValue(connected), retry: vi.fn().mockResolvedValue({ action: 'retry', outcome: 'accepted', message: 'Trying again.' }), quit: vi.fn().mockResolvedValue({ action: 'quit', outcome: 'accepted', message: 'Closing.' }), subscribe: () => () => undefined }
const taskResult = { action: 'load', outcome: 'accepted' as const, message: 'Loaded.', workspace: { tasks: [], groups: [], loadedAt: '2026-09-07T00:00:00Z' } }
const tasks: TaskBridge = { workspace: vi.fn().mockResolvedValue(taskResult), task: vi.fn(), previewTask: vi.fn(), saveTask: vi.fn(), runTask: vi.fn(), setTaskEnabled: vi.fn(), deleteTask: vi.fn(), saveGroup: vi.fn(), setGroupEnabled: vi.fn(), deleteGroup: vi.fn() }
const operations: OperationsBridge = { scheduleWindow: vi.fn().mockResolvedValue({ action: 'load_schedule', outcome: 'accepted', message: 'Loaded.', schedule: { from: '2026-09-06T00:00:00Z', to: '2026-09-14T00:00:00Z', loadedAt: '2026-09-07T00:00:00Z', occurrences: [] } }), activityWorkspace: vi.fn().mockResolvedValue({ action: 'load_activity', outcome: 'accepted', message: 'Loaded.', activity: { runs: [], logs: [], alerts: [], logPath: '', loadedAt: '2026-09-07T00:00:00Z' } }), acknowledgeAlert: vi.fn(), acknowledgeAlerts: vi.fn() }
const settingsWorkspace: SettingsWorkspace = { preferences: { version: 1, appearance: 'system', transition: { status: 'not_found', retired: ['appearance.font', 'appearance.scroll_sensitivity'] } }, preferencePath: '/home/ada/.config/go-schedule/desktop/preferences.json', storage: [], product: { name: 'go-schedule', version: '1.2.0', publisher: 'ShruggieTech', links: [] }, daemonAvailable: true, loadedAt: '2026-09-07T00:00:00Z' }
const settings: SettingsBridge = { workspace: vi.fn().mockResolvedValue({ action: 'load_settings', outcome: 'accepted', message: 'Loaded.', workspace: settingsWorkspace }), saveAppearance: vi.fn().mockImplementation(async (appearance) => ({ action: 'save_appearance', outcome: 'accepted', message: 'Saved.', workspace: { ...settingsWorkspace, preferences: { ...settingsWorkspace.preferences, appearance } } })), restore: vi.fn(), copyStoragePath: vi.fn(), openProductLink: vi.fn() }
const notifications: NotificationBridge = { workspace: vi.fn().mockResolvedValue({ action: 'load_notifications', outcome: 'accepted', message: 'Loaded.', workspace: { channels: [], tasks: [], groups: [], deliveries: [], loadedAt: '2026-09-07T00:00:00Z' } }), saveChannel: vi.fn(), setChannelEnabled: vi.fn(), testChannel: vi.fn(), deleteChannel: vi.fn(), policy: vi.fn(), savePolicy: vi.fn() }

describe('production shell', () => {
  it('keeps remote identity visible and fails closed for unsupported or observe-only actions', async () => {
    const user = userEvent.setup()
    const remote: ConnectionSnapshot = { ...connected, target: { id: 'daemon-identity', profileId: 'profile-id', kind: 'remote', displayName: 'Production', endpoint: 'https://example.test', platform: 'linux', capabilities: ['tasks', 'schedule', 'activity'], permissions: ['read'] } }
    const remoteBridge: DesktopBridge = { ...bridge, snapshot: vi.fn().mockResolvedValue(remote) }
    render(<App bridge={remoteBridge} tasks={tasks} settings={settings} />)
    expect(await screen.findByText('Observe-only connection')).toBeVisible()
    expect(await screen.findByRole('button', { name: 'Create task' })).toBeDisabled()
    await user.click(screen.getByRole('button', { name: 'Automation Sources' }))
    expect(screen.getByRole('heading', { name: 'Automation Sources is unavailable for remote targets' })).toBeVisible()
    await user.click(screen.getByRole('button', { name: 'Connection details' }))
    expect(screen.getByRole('dialog', { name: 'Production connection' })).toHaveTextContent('https://example.test')
  })

  it('keeps target and page identity visible across operational routes', async () => {
    const user = userEvent.setup(); render(<App bridge={bridge} tasks={tasks} operations={operations} settings={settings} />)
    await waitFor(() => expect(screen.getByRole('heading', { level: 2, name: 'Tasks' })).toBeVisible())
    await user.click(screen.getByRole('button', { name: 'Activity' }))
    expect(screen.getByRole('heading', { level: 1, name: 'Activity' })).toBeVisible()
    expect(screen.getAllByText('This computer').length).toBeGreaterThan(0)
    expect(await screen.findByRole('heading', { level: 2, name: 'No matching activity' })).toBeVisible()
  })

  it('controls appearance, connection details, retry, and exit', async () => {
    const user = userEvent.setup(); render(<App bridge={bridge} tasks={tasks} settings={settings} />)
    await user.selectOptions(screen.getByLabelText('Appearance'), 'dark')
    await waitFor(() => expect(document.querySelector('.app')).toHaveAttribute('data-appearance', 'dark'))
    await user.click(screen.getByRole('button', { name: 'Connection details' }))
    expect(screen.getByRole('dialog', { name: 'This computer connection' })).toBeVisible()
    await user.click(screen.getByRole('button', { name: 'Exit' }))
    expect(bridge.quit).toHaveBeenCalled()
  })

  it('opens the complete Notifications workspace from primary navigation', async () => {
    const user = userEvent.setup(); render(<App bridge={bridge} tasks={tasks} notifications={notifications} settings={settings} />)
    await user.click(screen.getByRole('button', { name: 'Notifications' }))
    expect(screen.getByRole('heading', { level: 1, name: 'Notifications' })).toBeVisible()
    expect(await screen.findByRole('heading', { level: 2, name: 'No webhook channels' })).toBeVisible()
    expect(screen.getByRole('heading', { level: 2, name: 'No matching deliveries' })).toBeVisible()
  })
})
