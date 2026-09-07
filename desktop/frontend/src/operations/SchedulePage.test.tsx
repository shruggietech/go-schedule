import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'
import type { OperationsBridge, ScheduleOccurrence } from './model'
import { SchedulePage } from './SchedulePage'

const at = (hour: number) => `2026-09-07T${String(hour).padStart(2, '0')}:00:00Z`
const occurrence = (index: number, state = 'upcoming'): ScheduleOccurrence => ({ id: `item-${index}`, taskId: `task-${index}`, taskName: `Task ${index}`, time: at(index % 24), kind: state === 'upcoming' ? 'prediction' : 'recorded', state, ...(state === 'upcoming' ? {} : { runId: `run-${index}`, outcome: state }) })
const result = (occurrences: ScheduleOccurrence[]) => ({ action: 'load_schedule', outcome: 'accepted' as const, message: 'Loaded.', schedule: { from: at(0), to: '2026-09-14T00:00:00Z', loadedAt: at(1), occurrences } })
const bridgeFor = (occurrences: ScheduleOccurrence[]): OperationsBridge => ({ scheduleWindow: vi.fn().mockResolvedValue(result(occurrences)), activityWorkspace: vi.fn(), acknowledgeAlert: vi.fn(), acknowledgeAlerts: vi.fn(), subscribe: () => () => undefined })

describe('SchedulePage', () => {
  it('distinguishes predictions, recorded states, ranges, calendar counts, and stable selection', async () => {
    const user = userEvent.setup(); const rows = [occurrence(9), occurrence(10, 'success'), occurrence(11, 'failure'), occurrence(12, 'skipped'), occurrence(13, 'caught_up'), occurrence(14, 'queued'), occurrence(15, 'unavailable')]; const bridge = bridgeFor(rows)
    render(<SchedulePage bridge={bridge} available refreshToken={1} />)
    expect(await screen.findByText('Prediction')).toBeVisible(); expect(screen.getAllByText('Recorded run')).toHaveLength(6)
    for (const state of ['Success', 'Failed', 'Skipped', 'Caught up', 'Queued', 'Unavailable']) expect(screen.getByText(state)).toBeVisible()
    await user.click(screen.getAllByRole('button', { name: /2026/ })[0]); expect(screen.getByRole('heading', { name: 'Task 9' })).toBeVisible()
    await user.selectOptions(screen.getByLabelText('Window'), '30'); await waitFor(() => expect(bridge.scheduleWindow).toHaveBeenLastCalledWith(30)); expect(screen.getByRole('heading', { name: 'Task 9' })).toBeVisible()
    await user.selectOptions(screen.getByLabelText('View'), 'calendar'); const day = screen.getByRole('button', { name: /September 7, 2026, 7 occurrences/ }); await user.click(day); expect(screen.getByRole('heading', { name: 'September 7, 2026' })).toBeVisible(); await user.click(screen.getByRole('button', { name: 'Next month' })); expect(screen.getByRole('heading', { name: 'October 2026' })).toBeVisible(); await user.click(screen.getByRole('button', { name: 'Previous month' })); expect(screen.getByRole('heading', { name: 'September 2026' })).toBeVisible()
  })

  it('renders and selects more than 100 rows', async () => {
    const user = userEvent.setup(); render(<SchedulePage bridge={bridgeFor(Array.from({ length: 105 }, (_, index) => occurrence(index)))} available refreshToken={1} />)
    expect(await screen.findByText('105 occurrences', { exact: false })).toBeVisible(); await user.click(screen.getAllByRole('button', { name: /2026/ })[100]); expect(screen.getByRole('heading', { name: 'Task 100' })).toBeVisible()
  })
})
