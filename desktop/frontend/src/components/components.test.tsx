import { act, fireEvent, render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { createRef } from 'react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { Button, DataTable, Dialog, Disclosure, Field, Notice, StatePanel, StatusBadge, ToastRegion } from '.'

describe('shared component catalog', () => {
  afterEach(() => vi.useRealTimers())

  it('provides status, notice, state, field, table, disclosure, and notification semantics', () => {
    render(<><StatusBadge state="degraded" /><Notice title="Delayed">Updates are delayed.</Notice><StatePanel title="No tasks" detail="Create one later." busy /><Field label="Name" help="Required" error="Enter a name"><input /></Field><DataTable caption="Tasks" headings={['Name']} rows={[[<>Backup</>]]} /><Disclosure summary="Details">More information</Disclosure><ToastRegion message="Saved" /></>)
    expect(screen.getByText('Degraded')).toBeVisible()
    expect(screen.getByText('Updates are delayed.')).toBeVisible()
    expect(screen.getByRole('table', { name: 'Tasks' })).toBeVisible()
    expect(screen.getByText('Details')).toBeVisible()
    expect(document.querySelector('[aria-live="polite"]')).not.toBeNull()
    const field = screen.getByRole('textbox', { name: 'Name' })
    expect(field).toHaveAttribute('aria-invalid', 'true')
    expect(field).toHaveAccessibleDescription('Required Enter a name')
  })

  it('honors disabled buttons', async () => {
    const click = vi.fn(); const user = userEvent.setup()
    render(<Button disabled onClick={click}>Save</Button>)
    await user.click(screen.getByRole('button', { name: 'Save' }))
    expect(click).not.toHaveBeenCalled()
  })

  it('exposes compact semantic button variants and pending state', () => {
    render(<><Button variant="affirmative">Enable</Button><Button variant="danger">Delete</Button><Button variant="subtle">Cancel</Button><Button pending>Save</Button></>)
    expect(screen.getByRole('button', { name: 'Enable' })).toHaveClass('button-affirmative')
    expect(screen.getByRole('button', { name: 'Delete' })).toHaveClass('button-danger')
    expect(screen.getByRole('button', { name: 'Cancel' })).toHaveClass('button-subtle')
    expect(screen.getByRole('button', { name: 'Save' })).toBeDisabled()
    expect(screen.getByRole('button', { name: 'Save' })).toHaveAttribute('aria-busy', 'true')
  })

  it('traps dialog focus and returns it to the exact invoker', async () => {
    const user = userEvent.setup(); const close = vi.fn(); const ref = createRef<HTMLButtonElement>()
    const { rerender } = render(<><button ref={ref}>Open</button><Dialog open title="Confirm" invoker={ref.current} onClose={close} actions={<button>Confirm</button>}>Content</Dialog></>)
    rerender(<><button ref={ref}>Open</button><Dialog open title="Confirm" invoker={ref.current} onClose={close} actions={<button>Confirm</button>}>Content</Dialog></>)
    expect(screen.getByRole('button', { name: 'Confirm' })).toHaveFocus()
    await user.keyboard('{Escape}')
    expect(close).toHaveBeenCalledOnce()
    expect(ref.current).toHaveFocus()
    expect(screen.getByRole('button', { name: 'Close' }).parentElement).toHaveClass('dialog-actions')
  })

  it('dismisses persistent notices explicitly', async () => {
    const user = userEvent.setup()
    render(<Notice title="Action needed" tone="error">Try again.</Notice>)
    await user.click(screen.getByRole('button', { name: 'Dismiss Action needed' }))
    expect(screen.queryByText('Try again.')).not.toBeInTheDocument()
  })

  it('restores a dismissed JSX notice when its event identity changes', async () => {
    const user = userEvent.setup()
    const { rerender } = render(<Notice title="Could not save group" tone="error" identity="first"><span>Enter a group name.</span></Notice>)
    await user.click(screen.getByRole('button', { name: 'Dismiss Could not save group' }))
    rerender(<Notice title="Could not save group" tone="error" identity="second"><span>The group name is already in use.</span></Notice>)
    expect(screen.getByText('The group name is already in use.')).toBeVisible()
  })

  it('replaces and dismisses transient feedback after five seconds', async () => {
    vi.useFakeTimers()
    const { rerender } = render(<ToastRegion message="First saved" />)
    expect(screen.getByRole('status')).toHaveTextContent('First saved')
    rerender(<ToastRegion message="Second saved" />)
    expect(screen.getByRole('status')).toHaveTextContent('Second saved')
    expect(screen.queryByText('First saved')).not.toBeInTheDocument()
    await act(async () => vi.advanceTimersByTimeAsync(5_000))
    expect(screen.queryByText('Second saved')).not.toBeInTheDocument()
  })

  it('restarts identical transient feedback when its event identity changes', async () => {
    const user = userEvent.setup()
    const { rerender } = render(<ToastRegion message="Appearance saved." identity={1} />)
    await user.click(screen.getByRole('button', { name: 'Dismiss notification' }))
    rerender(<ToastRegion message="Appearance saved." identity={2} />)
    expect(screen.getByRole('status')).toHaveTextContent('Appearance saved.')
  })

  it('restarts the full timeout for an identical event that arrives while visible', async () => {
    vi.useFakeTimers()
    const { rerender } = render(<ToastRegion message="Appearance saved." identity={1} />)
    await act(async () => vi.advanceTimersByTimeAsync(4_900))
    rerender(<ToastRegion message="Appearance saved." identity={2} />)
    await act(async () => vi.advanceTimersByTimeAsync(200))
    expect(screen.getByText('Appearance saved.')).toBeVisible()
    await act(async () => vi.advanceTimersByTimeAsync(4_800))
    expect(screen.queryByText('Appearance saved.')).not.toBeInTheDocument()
  })

  it('preserves compound success content in transient feedback', () => {
    render(<Notice title="Ready" tone="success"><><span>Preview generated.</span><strong> Review the command.</strong></></Notice>)
    expect(screen.getByRole('status')).toHaveTextContent('Ready: Preview generated. Review the command.')
  })

  it('pauses transient dismissal while hovered', async () => {
    vi.useFakeTimers()
    render(<ToastRegion message="Saved" />)
    const toast = screen.getByRole('status')
    fireEvent.mouseEnter(toast)
    await act(async () => vi.advanceTimersByTimeAsync(6_000))
    expect(screen.getByText('Saved')).toBeVisible()
    fireEvent.mouseLeave(toast)
    await act(async () => vi.advanceTimersByTimeAsync(5_000))
    expect(screen.queryByText('Saved')).not.toBeInTheDocument()
  })
})
