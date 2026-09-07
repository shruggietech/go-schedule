import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it } from 'vitest'
import { useState } from 'react'
import { CommandField, suggestions } from './CommandField'

function Harness({ platform, isNew = true }: { platform: string; isNew?: boolean }) { const [value, setValue] = useState(''); return <><CommandField platform={platform} value={value} isNew={isNew} onChange={setValue} /><button>Next</button></> }
describe('safe command suggestion', () => {
  it.each(Object.entries(suggestions))('inserts the exact %s mapping once, then traverses', async (platform, expected) => { const user = userEvent.setup(); render(<Harness platform={platform} />); const command = screen.getByLabelText('Command line') as HTMLInputElement; command.focus(); await user.tab(); expect(command).toHaveValue(expected); expect(command.selectionStart).toBe(expected.length); await user.tab(); expect(screen.getByRole('button', { name: 'Insert example' })).toHaveFocus() })
  it('keeps Shift+Tab as normal reverse traversal', async () => { const user = userEvent.setup(); render(<Harness platform="linux" />); const command = screen.getByLabelText('Command line'); command.focus(); await user.tab({ shift: true }); expect(command).toHaveValue('') })
  it('does not expose or insert suggestions while editing', async () => { const user = userEvent.setup(); render(<Harness platform="linux" isNew={false} />); const command = screen.getByLabelText('Command line'); command.focus(); await user.tab(); expect(command).toHaveValue(''); expect(screen.getByRole('button', { name: 'Next' })).toHaveFocus(); expect(screen.queryByRole('button', { name: 'Insert example' })).not.toBeInTheDocument() })
})
