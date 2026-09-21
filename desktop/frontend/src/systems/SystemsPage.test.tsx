import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'
import type { DesktopBridge, SystemsSnapshot } from '../connection/model'
import { SystemsPage } from './SystemsPage'

const snapshot: SystemsSnapshot = {
  generation: 3,
  startedAt: '2026-09-21T15:00:00Z',
  completedAt: '2026-09-21T15:00:01Z',
  complete: true,
  observations: [
    { registration: { key: 'local', kind: 'local', label: 'Workshop', daemonId: 'local-daemon', shortDaemonId: 'local-da', platform: 'windows', version: '1.5.0' }, state: 'connected', observedAt: '2026-09-21T15:00:00Z', stale: false, summary: { schema: 'go-schedule.system-summary.v1', observed_at: '2026-09-21T15:00:00Z', active_task_count: 2, next_occurrence: { task_id: 'task-1', task_name: 'Daily export', scheduled_for: '2026-09-21T16:00:00Z' }, recent_failure_count: 0, unacknowledged_alert_count: 0, notification_problem_count: 0 } },
    { registration: { key: 'remote-1', profileId: 'remote-1', kind: 'remote', label: 'Workshop', endpoint: 'https://remote.test', daemonId: 'remote-daemon', shortDaemonId: 'remote-d', platform: 'linux', version: '1.5.0' }, state: 'timed_out', observedAt: '2026-09-21T14:00:00Z', stale: true, failure: { state: 'timed_out', message: 'The scheduler timed out.', action: 'Check the network.' }, summary: { schema: 'go-schedule.system-summary.v1', observed_at: '2026-09-21T14:00:00Z', active_task_count: 1, recent_failure_count: 1, recent_failure: { run_id: 'run-7', task_id: 'task-7', task_name: 'Archive', ended_at: '2026-09-21T13:00:00Z' }, unacknowledged_alert_count: 1, unacknowledged_alert: { alert_id: 'alert-8', task_id: 'task-8', severity: 'error', kind: 'run_failed', created_at: '2026-09-21T13:30:00Z' }, notification_problem_count: 1 } },
  ],
}

function bridge(): DesktopBridge {
  return { snapshot: vi.fn(), retry: vi.fn(), quit: vi.fn(), subscribe: () => () => undefined, allSystems: vi.fn().mockResolvedValue(snapshot) }
}

describe('SystemsPage', () => {
  it('keeps same-named registrations distinguishable and filters attention', async () => {
    const user = userEvent.setup()
    render(<SystemsPage bridge={bridge()} onOpen={vi.fn()} />)
    expect(await screen.findByText('This computer · local-da')).toBeVisible()
    expect(screen.getByText('https://remote.test · remote-d')).toBeVisible()
    expect(screen.getAllByRole('heading', { name: 'Workshop' })).toHaveLength(2)
    await user.click(screen.getByLabelText('Needs attention only'))
    expect(screen.queryByText('This computer · local-da')).not.toBeInTheDocument()
    expect(screen.getByText(/Stale since/)).toBeVisible()
  })

  it('carries exact registration and representative context into drill-down', async () => {
    const user = userEvent.setup()
    const onOpen = vi.fn().mockResolvedValue(undefined)
    render(<SystemsPage bridge={bridge()} onOpen={onOpen} />)
    await screen.findByText('https://remote.test · remote-d')
    await user.click(screen.getByRole('button', { name: 'Failed run' }))
    expect(onOpen).toHaveBeenCalledWith(expect.objectContaining({ registrationKey: 'remote-1', destination: 'activity', taskId: 'task-7', recordId: 'run-7' }))
    await user.click(screen.getByRole('button', { name: 'Alert' }))
    expect(onOpen).toHaveBeenLastCalledWith(expect.objectContaining({ registrationKey: 'remote-1', destination: 'activity', taskId: 'task-8', recordId: 'alert-8' }))
  })

  it('supports state filtering, search, and stable label sorting', async () => {
    const user = userEvent.setup()
    render(<SystemsPage bridge={bridge()} onOpen={vi.fn()} />)
    await screen.findByText('2 registered systems refreshed; 2 shown.')
    await user.selectOptions(screen.getByLabelText('Connection state'), 'connected')
    expect(screen.getByText('2 registered systems refreshed; 1 shown.')).toBeVisible()
    await user.clear(screen.getByLabelText('Search'))
    await user.type(screen.getByLabelText('Search'), 'remote.test')
    expect(screen.getByRole('heading', { name: 'No matching systems' })).toBeVisible()
  })

  it('publishes completed target observations while the aggregate refresh is still running', async () => {
    let publish!: (value: SystemsSnapshot) => void
    let complete!: (value: SystemsSnapshot) => void
    const pending = new Promise<SystemsSnapshot>((resolve) => { complete = resolve })
    const progressiveBridge: DesktopBridge = { ...bridge(), allSystems: vi.fn(() => pending), subscribeSystems: (listener) => { publish = listener; return () => undefined } }
    render(<SystemsPage bridge={progressiveBridge} onOpen={vi.fn()} />)
    publish({ ...snapshot, generation: 4, completedAt: '', complete: false, observations: [snapshot.observations[0]] })
    expect(await screen.findByText('This computer · local-da')).toBeVisible()
    expect(screen.getByText('Refreshing registered systems.')).toBeVisible()
    complete({ ...snapshot, generation: 4 })
    expect(await screen.findByText('2 registered systems refreshed; 2 shown.')).toBeVisible()
  })

  it('routes remote notification context to the supported activity destination', async () => {
    const user = userEvent.setup()
    const onOpen = vi.fn().mockResolvedValue(undefined)
    render(<SystemsPage bridge={bridge()} onOpen={onOpen} />)
    await screen.findByText('https://remote.test · remote-d')
    await user.click(screen.getByRole('button', { name: 'Notification issue' }))
    expect(onOpen).toHaveBeenCalledWith(expect.objectContaining({ registrationKey: 'remote-1', destination: 'activity' }))
  })
})
