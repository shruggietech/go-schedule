import { useState } from 'react'
import { Button, Notice } from '../components'

export interface PairingDraft { address: string; daemon_id: string; pairing_id: string; phrase: string; certificate_pem: string; display_name: string; capability: string }
export interface PairingResult { action: string; outcome: string; message: string; credential_id?: string }
type NativeWindow = Window & { go?: { main?: { App?: { PairRemote?(draft: PairingDraft): Promise<PairingResult> } } } }

export function PairingForm({ nativeWindow = window }: { nativeWindow?: NativeWindow }) {
  const [draft, setDraft] = useState<PairingDraft>({ address: '', daemon_id: '', pairing_id: '', phrase: '', certificate_pem: '', display_name: '', capability: 'observe' })
  const [pending, setPending] = useState(false)
  const [message, setMessage] = useState('')
  const update = (key: keyof PairingDraft, value: string) => setDraft(current => ({ ...current, [key]: value }))
  const pair = async () => {
    setPending(true)
    const result = await nativeWindow.go?.main?.App?.PairRemote?.(draft) ?? { action: 'pair_remote', outcome: 'unavailable', message: 'Remote pairing is available in the installed desktop application.' }
    setMessage(result.message)
    setDraft(current => ({ ...current, phrase: '' }))
    setPending(false)
  }
  return <section className="panel"><div><p className="eyebrow">Remote daemon</p><h2>Pair this desktop</h2><p>Paste the values shown by the administrator exactly. The phrase is cleared after the attempt and the durable credential is stored by the operating system.</p></div>{message && <Notice title="Pairing result" tone="info">{message}</Notice>}<label>HTTPS address<input value={draft.address} onChange={event => update('address', event.target.value)} placeholder="https://192.0.2.10:8443" /></label><label>Daemon ID<input value={draft.daemon_id} onChange={event => update('daemon_id', event.target.value)} /></label><label>Pairing ID<input value={draft.pairing_id} onChange={event => update('pairing_id', event.target.value)} /></label><label>Client display name<input value={draft.display_name} onChange={event => update('display_name', event.target.value)} /></label><label>Requested capability<select value={draft.capability} onChange={event => update('capability', event.target.value)}><option value="observe">Observe</option><option value="operate">Operate</option><option value="manage">Manage</option><option value="enroll">Enroll</option></select></label><label>One-time phrase<input type="password" autoComplete="off" value={draft.phrase} onChange={event => update('phrase', event.target.value)} /></label><label>Trusted certificate (PEM)<textarea value={draft.certificate_pem} onChange={event => update('certificate_pem', event.target.value)} /></label><Button disabled={pending} onClick={() => void pair()}>{pending ? 'Pairing' : 'Pair remote daemon'}</Button></section>
}
