import { useEffect, useState } from 'react'
import { Button, CardSection, Notice, StatusLabel } from '../components'
import type { DesktopBridge, SystemRegistration } from '../connection/model'
import type { PopupPreferences, PopupStatus, SettingsBridge } from '../settings/model'

const defaults: PopupPreferences = { enabled: false, conditions: ['failure', 'success', 'alert'], severities: ['info', 'warning', 'error'], daemonIds: [] }
const choices = { conditions: [['failure', 'Failed task'], ['success', 'Successful task'], ['alert', 'Scheduler alert']], severities: [['info', 'Info'], ['warning', 'Warning'], ['error', 'Error']] } as const
const toggle = (values: string[], value: string) => values.includes(value) ? values.filter((item) => item !== value) : [...values, value]

export function DesktopPopups({ settings, desktop }: { settings: SettingsBridge; desktop: DesktopBridge }) {
  const [prefs, setPrefs] = useState<PopupPreferences>(defaults)
  const [status, setStatus] = useState<PopupStatus>()
  const [registrations, setRegistrations] = useState<SystemRegistration[]>([])
  const [message, setMessage] = useState('')
  const [pending, setPending] = useState(false)
  const missingDaemonIds = prefs.daemonIds.filter((id) => !registrations.some((registration) => registration.daemonId === id))
  useEffect(() => {
    let active = true
    void settings.workspace().then((result) => { if (active && result.workspace) setPrefs(result.workspace.preferences.popups ?? defaults) }).catch(() => { if (active) setMessage('Desktop preferences could not be loaded.') })
    void settings.popupStatus?.().then((value) => { if (active) setStatus(value) }).catch(() => undefined)
    void desktop.allSystems?.().then((value) => { if (active) setRegistrations(value.observations.map((item) => ({ ...item.registration, daemonId: item.registration.kind === 'local' && item.state !== 'connected' ? undefined : item.registration.daemonId }))) }).catch(() => undefined)
    return () => { active = false }
  }, [settings, desktop])
  const save = async () => {
    if (!settings.savePopups) return
    setPending(true); setMessage('')
    try {
      const result = await settings.savePopups(prefs)
      setMessage(result.message)
      if (result.workspace?.preferences.popups) setPrefs(result.workspace.preferences.popups)
      const latest = await settings.popupStatus?.()
      if (latest) setStatus(latest)
    } catch { setMessage('Desktop popup preferences could not be saved. Try again.') }
    finally { setPending(false) }
  }
  return <CardSection eyebrow="This desktop" title="Desktop popups" status={<StatusLabel tone={prefs.enabled && status?.authorized ? 'positive' : 'neutral'}>{prefs.enabled && status?.authorized ? 'On' : 'Off'}</StatusLabel>}>
    <p>Optional operating-system popups for new activity while this app is open. They do not change scheduler results or webhook delivery. No popups are sent while the desktop app is closed.</p>
    <p>{status?.message ?? 'Checking native notification support.'} A supported system may require notification permission, and click-to-open depends on its notification service.</p>
    <label className="popup-choice"><input type="checkbox" checked={prefs.enabled} disabled={status?.available === false && !prefs.enabled} onChange={(event) => setPrefs({ ...prefs, enabled: event.target.checked })} /> Enable desktop popups</label>
    <div className="popup-filter-grid">
      <fieldset><legend>Conditions</legend>{choices.conditions.map(([value, label]) => <label className="popup-choice" key={value}><input type="checkbox" checked={prefs.conditions.includes(value)} onChange={() => setPrefs({ ...prefs, conditions: toggle(prefs.conditions, value) })} /> {label}</label>)}</fieldset>
      <fieldset><legend>Alert severity</legend>{choices.severities.map(([value, label]) => <label className="popup-choice" key={value}><input type="checkbox" checked={prefs.severities.includes(value)} onChange={() => setPrefs({ ...prefs, severities: toggle(prefs.severities, value) })} /> {label}</label>)}</fieldset>
      <fieldset><legend>Schedulers</legend><p>Leave all unchecked to include every registered scheduler.</p>{registrations.map((registration) => <label className="popup-choice" key={registration.key}><input type="checkbox" checked={prefs.daemonIds.includes(registration.daemonId ?? '')} disabled={!registration.daemonId} onChange={() => setPrefs({ ...prefs, daemonIds: toggle(prefs.daemonIds, registration.daemonId ?? '') })} /> {registration.label}</label>)}{missingDaemonIds.map((id) => <label className="popup-choice" key={id}><input type="checkbox" checked onChange={() => setPrefs({ ...prefs, daemonIds: toggle(prefs.daemonIds, id) })} /> Unavailable scheduler <code>{id.slice(0, 12)}</code></label>)}</fieldset>
    </div>
    {missingDaemonIds.length > 0 && <Notice title="Saved scheduler no longer available" tone="warning">A saved scheduler filter cannot be matched. Uncheck its unavailable entry and save to remove it, or leave it selected until that scheduler is registered again.</Notice>}
    {message && <Notice title="Desktop popup settings" tone={message.includes('saved') ? 'success' : 'warning'}>{message}</Notice>}
    <div className="actions"><Button pending={pending} disabled={pending || prefs.conditions.length === 0 || prefs.severities.length === 0} onClick={() => void save()}>Save popup choices</Button></div>
  </CardSection>
}
