import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'
import { App } from './App'
import type { ConnectionSnapshot, DesktopBridge } from './connection/model'

const connected: ConnectionSnapshot = { generation: 1, state: 'connected', target: { id: 'local', displayName: 'This computer', platform: 'linux', version: '1.2.0', capabilities: ['tasks'], permissions: ['read'] }, message: 'Scheduler service is available.' }
const bridge: DesktopBridge = { snapshot: vi.fn().mockResolvedValue(connected), retry: vi.fn().mockResolvedValue({ action: 'retry', outcome: 'accepted', message: 'Trying again.' }), quit: vi.fn().mockResolvedValue({ action: 'quit', outcome: 'accepted', message: 'Closing.' }), subscribe: () => () => undefined }

describe('production shell', () => {
  it('keeps target and page identity visible across honest placeholders', async () => {
    const user = userEvent.setup(); render(<App bridge={bridge} />)
    await waitFor(() => expect(screen.getByText('Desktop foundation ready')).toBeVisible())
    await user.click(screen.getByRole('button', { name: 'Activity' }))
    expect(screen.getByRole('heading', { level: 1, name: 'Activity' })).toBeVisible()
    expect(screen.getAllByText('This computer').length).toBeGreaterThan(0)
    expect(screen.getByText(/not migrated yet/i)).toBeVisible()
  })

  it('controls appearance, connection details, retry, and exit', async () => {
    const user = userEvent.setup(); render(<App bridge={bridge} />)
    await user.selectOptions(screen.getByLabelText('Appearance'), 'dark')
    expect(document.querySelector('.app')).toHaveAttribute('data-appearance', 'dark')
    await user.click(screen.getByRole('button', { name: 'Connection details' }))
    expect(screen.getByRole('dialog', { name: 'This computer connection' })).toBeVisible()
    await user.click(screen.getByRole('button', { name: 'Exit' }))
    expect(bridge.quit).toHaveBeenCalled()
  })
})
