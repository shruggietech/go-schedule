import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { expect, it, vi } from 'vitest'
import type { TaskBridge, TaskSummary } from './model'
import { TaskActions } from './TaskActions'

it('names the target and suppresses duplicate Run now activation', async () => {
  const user = userEvent.setup(); let finish: ((value: unknown) => void) | undefined; const runTask = vi.fn(() => new Promise((resolve) => { finish = resolve })); const bridge = { runTask } as unknown as TaskBridge; const task = { id: 'one', name: 'Nightly backup', commandConfigured: true, declaredEnabled: false, effectiveState: 'manual_only' } as TaskSummary
  render(<TaskActions task={task} bridge={bridge} onResult={() => undefined} />); await user.click(screen.getByRole('button', { name: 'Run now' })); expect(screen.getByRole('dialog', { name: 'Run Nightly backup now?' })).toBeVisible(); await user.click(screen.getByRole('button', { name: 'Confirm run' })); expect(runTask).toHaveBeenCalledTimes(1); expect(screen.getByRole('button', { name: 'Confirm run' })).toBeDisabled(); finish?.({ action: 'run_task', outcome: 'accepted', message: 'Accepted.' })
})

it('separates operate permission from manage-only controls', () => {
  const task = { id: 'one', name: 'Nightly backup', commandConfigured: true, declaredEnabled: false, effectiveState: 'manual_only' } as TaskSummary
  render(<TaskActions task={task} bridge={{} as TaskBridge} available manageAvailable={false} onResult={() => undefined} />)
  expect(screen.getByRole('button', { name: 'Run now' })).toBeEnabled()
  expect(screen.getByRole('button', { name: 'Enable' })).toBeDisabled()
  expect(screen.getByRole('button', { name: 'Delete' })).toBeDisabled()
  expect(screen.getByRole('button', { name: 'Enable' })).toHaveClass('button-affirmative')
  expect(screen.getByRole('button', { name: 'Delete' })).toHaveClass('button-danger')
})
