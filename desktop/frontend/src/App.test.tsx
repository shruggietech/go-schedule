import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'
import { App } from './App'
import type { ConnectionSnapshot, DesktopBridge } from './connection/model'
import type { TaskBridge } from './tasks/model'
import type { OperationsBridge } from './operations/model'

const connected: ConnectionSnapshot = { generation: 1, revision: 2, state: 'connected', target: { id: 'local', displayName: 'This computer', platform: 'linux', version: '1.2.0', capabilities: ['tasks'], permissions: ['read'] }, message: 'Scheduler service is available.' }
const bridge: DesktopBridge = { snapshot: vi.fn().mockResolvedValue(connected), retry: vi.fn().mockResolvedValue({ action: 'retry', outcome: 'accepted', message: 'Trying again.' }), quit: vi.fn().mockResolvedValue({ action: 'quit', outcome: 'accepted', message: 'Closing.' }), subscribe: () => () => undefined }
const taskResult = { action: 'load', outcome: 'accepted' as const, message: 'Loaded.', workspace: { tasks: [], groups: [], loadedAt: '2026-09-07T00:00:00Z' } }
const tasks: TaskBridge = { workspace: vi.fn().mockResolvedValue(taskResult), task: vi.fn(), previewTask: vi.fn(), saveTask: vi.fn(), runTask: vi.fn(), setTaskEnabled: vi.fn(), deleteTask: vi.fn(), saveGroup: vi.fn(), setGroupEnabled: vi.fn(), deleteGroup: vi.fn() }
const operations: OperationsBridge = { scheduleWindow: vi.fn().mockResolvedValue({ action: 'load_schedule', outcome: 'accepted', message: 'Loaded.', schedule: { from: '2026-09-06T00:00:00Z', to: '2026-09-14T00:00:00Z', loadedAt: '2026-09-07T00:00:00Z', occurrences: [] } }), activityWorkspace: vi.fn().mockResolvedValue({ action: 'load_activity', outcome: 'accepted', message: 'Loaded.', activity: { runs: [], logs: [], alerts: [], logPath: '', loadedAt: '2026-09-07T00:00:00Z' } }), acknowledgeAlert: vi.fn(), acknowledgeAlerts: vi.fn() }

describe('production shell', () => {
  it('keeps target and page identity visible across operational routes', async () => {
    const user = userEvent.setup(); render(<App bridge={bridge} tasks={tasks} operations={operations} />)
    await waitFor(() => expect(screen.getByRole('heading', { level: 2, name: 'Tasks' })).toBeVisible())
    await user.click(screen.getByRole('button', { name: 'Activity' }))
    expect(screen.getByRole('heading', { level: 1, name: 'Activity' })).toBeVisible()
    expect(screen.getAllByText('This computer').length).toBeGreaterThan(0)
    expect(await screen.findByRole('heading', { level: 2, name: 'No matching activity' })).toBeVisible()
  })

  it('controls appearance, connection details, retry, and exit', async () => {
    const user = userEvent.setup(); render(<App bridge={bridge} tasks={tasks} />)
    await user.selectOptions(screen.getByLabelText('Appearance'), 'dark')
    expect(document.querySelector('.app')).toHaveAttribute('data-appearance', 'dark')
    await user.click(screen.getByRole('button', { name: 'Connection details' }))
    expect(screen.getByRole('dialog', { name: 'This computer connection' })).toBeVisible()
    await user.click(screen.getByRole('button', { name: 'Exit' }))
    expect(bridge.quit).toHaveBeenCalled()
  })
})
