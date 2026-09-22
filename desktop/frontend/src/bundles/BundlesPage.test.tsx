import { fireEvent, render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'
import { BundlesPage } from './BundlesPage'
import type { BundleBridge } from './model'

describe('BundlesPage', () => {
  it('sends watcher paths as target-local preview bindings', async () => {
    const preview = vi.fn().mockResolvedValue({ action: 'preview_bundle', outcome: 'accepted', message: 'Review the plan.', plan: { id: 'plan', bundle_digest: 'digest', target_daemon_id: 'target', target_fingerprint: 'fingerprint', items: [] } })
    const bridge: BundleBridge = { exportBundle: vi.fn(), validate: vi.fn(), compare: vi.fn(), preview, apply: vi.fn() }
    const user = userEvent.setup()
    render(<BundlesPage bridge={bridge} targetName="Target" manageAvailable />)
    const document = '{"schema":"go-schedule.bundle/v2","watchers":[{"portable_id":"watcher-1","name":"Files"}]}'
    fireEvent.change(screen.getByLabelText('Bundle JSON'), { target: { value: document } })
    await user.type(screen.getByLabelText('Files'), 'C:\\Data\\Files')
    await user.click(screen.getByRole('button', { name: 'Preview changes' }))
    expect(preview).toHaveBeenCalledWith(document, { 'watcher-1': 'C:\\Data\\Files' })
  })
})
