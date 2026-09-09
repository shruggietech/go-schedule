import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import { PairingForm, type PairingDraft } from './PairingForm'

describe('PairingForm', () => {
  it('sends the bounded enrollment fields and clears the phrase', async () => {
    const pair = vi.fn(async (_draft: PairingDraft) => ({ action: 'pair_remote', outcome: 'accepted', message: 'Paired.' }))
    const nativeWindow = { go: { main: { App: { PairRemote: pair } } } } as unknown as Window
    render(<PairingForm nativeWindow={nativeWindow} />)
    fireEvent.change(screen.getByLabelText('HTTPS address'), { target: { value: 'https://daemon.test:8443' } })
    fireEvent.change(screen.getByLabelText('Daemon ID'), { target: { value: 'daemon-id' } })
    fireEvent.change(screen.getByLabelText('Pairing ID'), { target: { value: 'pairing-id' } })
    fireEvent.change(screen.getByLabelText('Client display name'), { target: { value: 'My desktop' } })
    fireEvent.change(screen.getByLabelText('One-time phrase'), { target: { value: 'one-time-phrase' } })
    fireEvent.change(screen.getByLabelText('Trusted certificate (PEM)'), { target: { value: 'certificate' } })
    fireEvent.click(screen.getByRole('button', { name: 'Pair remote daemon' }))
    await waitFor(() => expect(pair).toHaveBeenCalledWith(expect.objectContaining({ address: 'https://daemon.test:8443', daemon_id: 'daemon-id', pairing_id: 'pairing-id', display_name: 'My desktop', capability: 'observe', phrase: 'one-time-phrase', certificate_pem: 'certificate' })))
    expect(screen.getByLabelText('One-time phrase')).toHaveValue('')
  })
})
