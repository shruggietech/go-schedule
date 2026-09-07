import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'
import { ConnectionsPage } from './ConnectionsPage'
import type { ConnectionSnapshot, ConnectionState } from '../connection/model'

const base: ConnectionSnapshot = { generation: 1, revision: 1, state: 'connected', target: { id: 'local', displayName: 'This computer', platform: 'windows', version: '1.2.0', capabilities: ['tasks'], permissions: ['read'] }, message: 'Scheduler service is available.', lastSuccessfulAt: '2026-09-07T12:00:00Z' }

describe('ConnectionsPage', () => {
  for (const state of ['connected', 'unavailable', 'timed_out', 'access_denied', 'incompatible', 'recovering'] satisfies ConnectionState[]) {
    it(`renders explicit ${state} diagnosis`, () => {
      render(<ConnectionsPage snapshot={{ ...base, state, action: state === 'connected' ? undefined : 'Try again.' }} retryPending={false} onRetry={vi.fn()} />)
      expect(screen.getByText(state === 'access_denied' ? 'Access denied' : state === 'timed_out' ? 'Timed out' : state[0].toUpperCase() + state.slice(1))).toBeVisible()
      expect(screen.getByRole('heading', { name: 'This computer' })).toBeVisible()
    })
  }

  it('offers one disabled retry while pending', async () => {
    const user = userEvent.setup(); const retry = vi.fn()
    const { rerender } = render(<ConnectionsPage snapshot={{ ...base, state: 'unavailable', action: 'Start the service.' }} retryPending={false} onRetry={retry} />)
    await user.click(screen.getByRole('button', { name: 'Try again' }))
    expect(retry).toHaveBeenCalledOnce()
    rerender(<ConnectionsPage snapshot={{ ...base, state: 'recovering', action: 'Try again.' }} retryPending onRetry={retry} />)
    expect(screen.getByRole('button', { name: 'Trying again' })).toBeDisabled()
  })
})
