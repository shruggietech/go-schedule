import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { createRef } from 'react'
import { describe, expect, it, vi } from 'vitest'
import { Button, DataTable, Dialog, Disclosure, Field, Notice, StatePanel, StatusBadge, ToastRegion } from '.'

describe('shared component catalog', () => {
  it('provides status, notice, state, field, table, disclosure, and notification semantics', () => {
    render(<><StatusBadge state="degraded" /><Notice title="Delayed">Updates are delayed.</Notice><StatePanel title="No tasks" detail="Create one later." busy /><Field label="Name" help="Required" error="Enter a name"><input /></Field><DataTable caption="Tasks" headings={['Name']} rows={[[<>Backup</>]]} /><Disclosure summary="Details">More information</Disclosure><ToastRegion message="Saved" /></>)
    expect(screen.getByText('Degraded')).toBeVisible()
    expect(screen.getByText('Updates are delayed.')).toBeVisible()
    expect(screen.getByRole('table', { name: 'Tasks' })).toBeVisible()
    expect(screen.getByText('Details')).toBeVisible()
    expect(document.querySelector('[aria-live="polite"]')).not.toBeNull()
  })

  it('honors disabled buttons', async () => {
    const click = vi.fn(); const user = userEvent.setup()
    render(<Button disabled onClick={click}>Save</Button>)
    await user.click(screen.getByRole('button', { name: 'Save' }))
    expect(click).not.toHaveBeenCalled()
  })

  it('traps dialog focus and returns it to the exact invoker', async () => {
    const user = userEvent.setup(); const close = vi.fn(); const ref = createRef<HTMLButtonElement>()
    const { rerender } = render(<><button ref={ref}>Open</button><Dialog open title="Confirm" invoker={ref.current} onClose={close}><button>Confirm</button></Dialog></>)
    rerender(<><button ref={ref}>Open</button><Dialog open title="Confirm" invoker={ref.current} onClose={close}><button>Confirm</button></Dialog></>)
    expect(screen.getByRole('button', { name: 'Confirm' })).toHaveFocus()
    await user.keyboard('{Escape}')
    expect(close).toHaveBeenCalledOnce()
    expect(ref.current).toHaveFocus()
  })
})
