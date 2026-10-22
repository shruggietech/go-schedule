import { useEffect, useState } from 'react'
import { Field } from '../components'
import { splitSMTPRecipients } from './recipients'
import type { ChannelDraft, SMTPConfig } from './model'

export const blankSMTP = (): SMTPConfig => ({ host: '', port: 587, tls_mode: 'starttls', from: '', recipients: [], timeout_seconds: 10, html: false })

export function SMTPFields({ draft, change }: { draft: ChannelDraft; change: (draft: ChannelDraft) => void }) {
  const config = draft.smtp ?? blankSMTP()
  const [recipients, setRecipients] = useState(config.recipients.join(', '))
  useEffect(() => { setRecipients((draft.smtp?.recipients ?? []).join(', ')) }, [draft.id, draft.isNew])
  const update = (values: Partial<SMTPConfig>) => change({ ...draft, smtp: { ...config, ...values } })
  return <fieldset><legend>SMTP email settings</legend>
    <Field label="SMTP server" help="Hostname or IP address, without a URL scheme or port."><input value={config.host} autoComplete="off" onChange={(event) => update({ host: event.target.value })} /></Field>
    <Field label="SMTP port"><input type="number" min="1" max="65535" value={config.port} onChange={(event) => update({ port: Number(event.target.value) })} /></Field>
    <Field label="Encryption" help="STARTTLS is required when selected. TLS certificates and hostnames are verified. Plain relay sends unencrypted traffic."><select value={config.tls_mode} onChange={(event) => update({ tls_mode: event.target.value as SMTPConfig['tls_mode'] })}><option value="starttls">Required STARTTLS</option><option value="tls">Implicit TLS</option><option value="plain">Plain relay (unencrypted)</option></select></Field>
    <Field label="Sender email"><input type="text" inputMode="email" value={config.from} onChange={(event) => update({ from: event.target.value })} /></Field>
    <Field label="Recipient emails" help="One to 100 addresses separated by commas or new lines."><textarea value={recipients} onChange={(event) => { const value = event.target.value; setRecipients(value); update({ recipients: splitSMTPRecipients(value) }) }} /></Field>
    <Field label="SMTP timeout (seconds)" help="One shared deadline for the complete attempt, from 1 to 60 seconds."><input type="number" min="1" max="60" value={config.timeout_seconds} onChange={(event) => update({ timeout_seconds: Number(event.target.value) })} /></Field>
    <label className="check-field"><input type="checkbox" checked={config.html} onChange={(event) => update({ html: event.target.checked })} /> Include simple HTML alongside plain text</label>
    {!draft.isNew && <label className="check-field"><input type="checkbox" checked={draft.replaceAuthorization} onChange={(event) => change({ ...draft, replaceAuthorization: event.target.checked, smtpUsername: '', smtpPassword: '' })} /> Replace or remove stored SMTP credentials</label>}
    {(draft.isNew || draft.replaceAuthorization) && <>
      <Field label="SMTP username (optional)" help="Stored credentials are never redisplayed. Leave both fields blank for an unauthenticated relay or to clear existing credentials."><input value={draft.smtpUsername ?? ''} autoComplete="off" onChange={(event) => change({ ...draft, smtpUsername: event.target.value })} /></Field>
      <Field label="SMTP password"><input type="password" value={draft.smtpPassword ?? ''} autoComplete="new-password" onChange={(event) => change({ ...draft, smtpPassword: event.target.value })} /></Field>
    </>}
  </fieldset>
}
