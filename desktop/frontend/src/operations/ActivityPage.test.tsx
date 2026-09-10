import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'
import type { ActivityWorkspace, OperationsBridge } from './model'
import { ActivityPage, activityItems } from './ActivityPage'

const workspace: ActivityWorkspace = { loadedAt: '2026-09-07T12:00:00Z', logPath: 'C:\\ProgramData\\goschedule\\logs\\goschedule.log', runs: [{ id: 'run-1', taskId: 'task-1', scheduledFor: '2026-09-07T11:00:00Z', startedAt: '2026-09-07T11:00:01Z', endedAt: '2026-09-07T11:00:02Z', state: 'failure', outcome: 'failure', exitCode: 7, output: 'broken', outputTruncated: true, trigger: 'watcher', sourceWatcherId: 'watcher-1' }], logs: [{ id: 'log-1', time: '2026-09-07T11:01:00Z', severity: 'warning', source: 'engine', message: 'slow operation', detail: 'elapsed: 3s' }], alerts: [{ id: 'alert-1', time: '2026-09-07T11:02:00Z', severity: 'error', kind: 'run_failed', message: 'Task failed', runId: 'run-1', acknowledged: false }] }
const accepted = { action: 'load_activity', outcome: 'accepted' as const, message: 'Loaded.', activity: workspace }
const bridge = (): OperationsBridge => ({ scheduleWindow: vi.fn(), activityWorkspace: vi.fn().mockResolvedValue(accepted), acknowledgeAlert: vi.fn().mockResolvedValue({ ...accepted, action: 'acknowledge_alerts' }), acknowledgeAlerts: vi.fn().mockResolvedValue({ ...accepted, action: 'acknowledge_alerts' }), subscribe: () => () => undefined })

describe('ActivityPage', () => {
  it('orders and clears runs by execution time instead of scheduled time', () => {
    const items = activityItems({ ...workspace, runs: [{ ...workspace.runs[0], scheduledFor: '2026-09-01T00:00:00Z', startedAt: '2026-09-07T11:59:00Z', endedAt: '2026-09-07T12:00:00Z' }] })
    expect(items[0].type).toBe('run'); expect(items[0].time).toBe('2026-09-07T12:00:00Z')
  })

  it('keeps record types distinct, filters, and exposes full run diagnostics', async () => {
    const user = userEvent.setup(); render(<ActivityPage bridge={bridge()} available refreshToken={1} />)
    expect(await screen.findByText(workspace.logPath)).toBeVisible(); expect(screen.getByText('3 matching records')).toBeVisible()
    await user.selectOptions(screen.getByLabelText('Record type'), 'run'); expect(screen.getByText('1 matching records')).toBeVisible(); expect(screen.queryByText('slow operation')).not.toBeInTheDocument()
    await user.click(screen.getByRole('button', { name: /9\/7\/2026/ })); expect(screen.getByRole('heading', { name: 'Run details' })).toBeVisible(); expect(screen.getByText('broken')).toBeVisible(); expect(screen.getByText('watcher-1')).toBeVisible(); expect(screen.getByText('Yes')).toBeVisible()
  })

  it('acknowledges one alert and clears only the visible view', async () => {
    const user = userEvent.setup(); const api = bridge(); render(<ActivityPage bridge={api} available refreshToken={1} />); await screen.findByText('3 matching records')
    await user.selectOptions(screen.getByLabelText('Record type'), 'alert'); await user.click(screen.getByRole('button', { name: /9\/7\/2026/ })); await user.click(screen.getByRole('button', { name: 'Acknowledge alert' })); await waitFor(() => expect(api.acknowledgeAlert).toHaveBeenCalledWith('alert-1'))
    await user.click(screen.getByRole('button', { name: 'Clear View' })); expect(screen.getByRole('heading', { name: 'No matching activity' })).toBeVisible(); expect(screen.getByText(/Persisted records are not deleted/)).toBeVisible()
  })

  it('filters more than 100 activity rows without losing search focus', async () => {
    const many = { ...workspace, runs: Array.from({ length: 105 }, (_, index) => ({ ...workspace.runs[0], id: `run-${index}`, taskId: `task-${index}`, state: index === 104 ? 'success' : 'failure', output: `output ${index}` })) }; const api = bridge(); api.activityWorkspace = vi.fn().mockResolvedValue({ ...accepted, activity: many }); const user = userEvent.setup()
    render(<ActivityPage bridge={api} available refreshToken={1} />); expect(await screen.findByText('107 matching records')).toBeVisible(); const search = screen.getByLabelText('Search'); await user.type(search, 'task-104'); expect(screen.getByText('1 matching records')).toBeVisible(); expect(search).toHaveFocus()
  })

  it('blocks alert mutation after an uncertain outcome until authoritative refresh', async () => {
    const api = bridge(); api.acknowledgeAlert = vi.fn().mockResolvedValue({ action: 'acknowledge_alerts', outcome: 'uncertain', message: 'The remote request may have completed. Refresh before trying again.' })
    const user = userEvent.setup(); render(<ActivityPage bridge={api} available refreshToken={1} />); await screen.findByText('3 matching records')
    await user.selectOptions(screen.getByLabelText('Record type'), 'alert'); await user.click(screen.getByRole('button', { name: /9\/7\/2026/ })); await user.click(screen.getByRole('button', { name: 'Acknowledge alert' }))
    await waitFor(() => expect(screen.getByText(/may have completed/)).toBeVisible())
    expect(screen.getByRole('button', { name: 'Acknowledge alert' })).toBeDisabled()
  })
})
