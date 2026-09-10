import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'
import { ConnectionsPage } from './ConnectionsPage'
import type { ConnectionSnapshot, ConnectionState, DesktopBridge } from '../connection/model'

const base: ConnectionSnapshot = { generation: 1, revision: 1, state: 'connected', target: { id: 'local', displayName: 'This computer', platform: 'windows', version: '1.2.0', capabilities: ['tasks'], permissions: ['read'] }, message: 'Scheduler service is available.', lastSuccessfulAt: '2026-09-07T12:00:00Z' }

describe('ConnectionsPage', () => {
  it('lists and selects same-named profiles by stable identity', async () => {
    const user = userEvent.setup()
    const selectConnection = vi.fn().mockResolvedValue({ action: 'select_connection', outcome: 'accepted' as const, message: 'Selected.', workspace: { activeProfileId: 'two', profiles: [] } })
    const bridge = { snapshot: vi.fn(), retry: vi.fn(), quit: vi.fn(), subscribe: () => () => undefined, connectionProfiles: vi.fn().mockResolvedValue({ action: 'load_connections', outcome: 'accepted' as const, message: 'Loaded.', workspace: { profiles: [{ id: 'one', label: 'Workshop', endpoint: 'https://one.test', daemonId: 'daemon-one', shortDaemonId: 'daemon-o', fingerprint: 'aa', capability: 'observe', platform: 'linux', active: false }, { id: 'two', label: 'Workshop', endpoint: 'https://two.test', daemonId: 'daemon-two', shortDaemonId: 'daemon-t', fingerprint: 'bb', capability: 'operate', platform: 'linux', active: false }] } }), selectConnection } satisfies DesktopBridge
    render(<ConnectionsPage snapshot={base} retryPending={false} onRetry={vi.fn()} bridge={bridge} />)
    expect(await screen.findByText('https://one.test · daemon-o')).toBeVisible()
    expect(screen.getByText('https://two.test · daemon-t')).toBeVisible()
    await user.click(screen.getAllByRole('button', { name: 'Select' })[1])
    expect(selectConnection).toHaveBeenCalledWith('two')
  })

  for (const state of ['connected', 'unavailable', 'timed_out', 'access_denied', 'unauthorized', 'revoked', 'forbidden', 'incompatible', 'trust_changed', 'identity_changed', 'recovering'] satisfies ConnectionState[]) {
    it(`renders explicit ${state} diagnosis`, () => {
      render(<ConnectionsPage snapshot={{ ...base, state, action: state === 'connected' ? undefined : 'Try again.' }} retryPending={false} onRetry={vi.fn()} />)
      const labels: Partial<Record<ConnectionState, string>> = { access_denied: 'Access denied', timed_out: 'Timed out', unauthorized: 'Unauthorized', revoked: 'Credential revoked', forbidden: 'Forbidden', trust_changed: 'Trust changed', identity_changed: 'Identity changed' }
      expect(screen.getByText(labels[state] ?? state[0].toUpperCase() + state.slice(1))).toBeVisible()
      expect(screen.getByRole('heading', { name: 'This computer' })).toBeVisible()
    })
  }

  it('shows stale data and bounded automatic retry metadata without relying on color', () => {
    render(<ConnectionsPage snapshot={{ ...base, state: 'recovering', stale: true, recovery: 'automatic', retryAttempt: 3, nextRetryAt: '2026-09-07T12:00:10Z' }} retryPending={false} onRetry={vi.fn()} />)
    expect(screen.getByText('Data may be stale')).toBeVisible()
    expect(screen.getByText('Automatic retry')).toBeVisible()
    expect(screen.getByText('2026-09-07T12:00:10Z')).toBeVisible()
  })

  it('offers one disabled retry while pending', async () => {
    const user = userEvent.setup(); const retry = vi.fn()
    const { rerender } = render(<ConnectionsPage snapshot={{ ...base, state: 'unavailable', action: 'Start the service.' }} retryPending={false} onRetry={retry} />)
    await user.click(screen.getByRole('button', { name: 'Try again' }))
    expect(retry).toHaveBeenCalledOnce()
    rerender(<ConnectionsPage snapshot={{ ...base, state: 'recovering', action: 'Try again.' }} retryPending onRetry={retry} />)
    expect(screen.getByRole('button', { name: 'Trying again' })).toBeDisabled()
  })

  it('keeps the focused retry control mounted across recovery and success', async () => {
    const user = userEvent.setup(); const retry = vi.fn()
    const { rerender } = render(<ConnectionsPage snapshot={{ ...base, state: 'unavailable', action: 'Start the service.' }} retryPending={false} onRetry={retry} />)
    const button = screen.getByRole('button', { name: 'Try again' })
    await user.click(button)
    rerender(<ConnectionsPage snapshot={{ ...base, state: 'recovering', action: undefined }} retryPending={false} onRetry={retry} />)
    expect(screen.getByRole('button', { name: 'Trying again' })).toBe(button)
    expect(button).toHaveFocus()
    rerender(<ConnectionsPage snapshot={base} retryPending={false} onRetry={retry} />)
    expect(screen.getByRole('button', { name: 'No retry needed' })).toBe(button)
    expect(button).toHaveFocus()
  })
})
