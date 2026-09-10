import { useCallback, useEffect, useState } from 'react'
import { Button, Notice, StatusBadge } from '../components'
import { desktopBridge } from '../connection/bridge'
import type { ConnectionResult, ConnectionSnapshot, ConnectionWorkspace, DesktopBridge } from '../connection/model'
import { PairingForm } from '../remotepairing/PairingForm'

const guidance: Record<ConnectionSnapshot['state'], string> = {
  connecting: 'Wait while go-schedule checks the selected endpoint.',
  connected: 'The selected scheduler service is ready. No recovery action is needed.',
  degraded: 'The last connection was interrupted. Current data remains visible while recovery continues.',
  recovering: 'A bounded retry is in progress. You can keep using desktop-local Settings.',
  unavailable: 'Start or repair the local scheduler service, then try again.',
  access_denied: 'Repair the connection or ask an administrator to update its grant.',
  unauthorized: 'Repair the saved credential, then reconnect.',
  revoked: 'This credential was revoked. Repair the profile with a new enrollment bundle.',
  forbidden: 'Ask an administrator for the authority required by this connection.',
  incompatible: 'Update the desktop application or selected scheduler to compatible versions.',
  trust_changed: 'Verify the endpoint certificate before repairing the saved trust record.',
  identity_changed: 'Verify that this is the intended daemon before replacing the saved profile.',
  timed_out: 'Check whether the selected scheduler service is responsive, then try again.',
}

export function ConnectionsPage({ snapshot, retryPending, onRetry, bridge = desktopBridge }: { snapshot: ConnectionSnapshot; retryPending: boolean; onRetry(): void; bridge?: DesktopBridge }) {
  const [workspace, setWorkspace] = useState<ConnectionWorkspace>()
  const [message, setMessage] = useState('')
  const [repairProfileId, setRepairProfileId] = useState('')
  const load = useCallback(async () => {
    const result = await bridge.connectionProfiles?.()
    if (result?.workspace) setWorkspace(result.workspace)
    if (result && result.outcome !== 'accepted') setMessage(result.message)
  }, [bridge])
  useEffect(() => { void load() }, [load])
  const accept = (result: ConnectionResult | undefined) => { if (!result) return; setMessage(result.message); if (result.workspace) setWorkspace(result.workspace) }
  const select = async (id: string) => accept(await bridge.selectConnection?.(id))
  const rename = async (id: string, current: string) => { const label = window.prompt('New local profile label', current); if (label !== null) accept(await bridge.renameConnection?.(id, label)) }
  const remove = async (id: string, label: string) => { if (window.confirm(`Remove ${label}? Its native credential will also be deleted.`)) accept(await bridge.removeConnection?.(id)) }
  const busy = snapshot.state === 'connecting' || snapshot.state === 'recovering' || retryPending
  const stateGuidance = snapshot.target.kind === 'remote' && snapshot.state === 'unavailable' ? 'Check the selected remote scheduler and network, then try again.' : guidance[snapshot.state]
  return <><header className="page-header"><div><p className="eyebrow">{snapshot.target.kind === 'remote' ? 'Remote daemon' : 'This computer'}</p><h1>Connections</h1><p>Select, diagnose, pair, and administer scheduler targets.</p></div><StatusBadge state={snapshot.state} /></header>{message && <Notice title="Connection update" tone="info">{message}</Notice>}<section className="panel connection-workspace" aria-busy={busy}><div><p className="eyebrow">Current diagnosis</p><h2>{snapshot.target.displayName}</h2><p>{snapshot.message}</p><p>{stateGuidance}</p>{snapshot.stale && <Notice title="Data may be stale" tone="warning">Last-known scheduler data remains visible until an authoritative refresh succeeds.</Notice>}</div><dl>{snapshot.target.endpoint && <><dt>Endpoint</dt><dd>{snapshot.target.endpoint}</dd></>}<dt>Platform</dt><dd>{snapshot.target.platform}{snapshot.target.architecture ? `/${snapshot.target.architecture}` : ''}</dd><dt>Version</dt><dd>{snapshot.target.version || 'Unavailable'}</dd><dt>Capabilities</dt><dd>{snapshot.target.capabilities.length ? snapshot.target.capabilities.join(', ') : 'Unavailable'}</dd><dt>Last successful connection</dt><dd>{snapshot.lastSuccessfulAt || 'Not available in this session'}</dd><dt>Recovery mode</dt><dd>{snapshot.recovery === 'automatic' ? 'Automatic retry' : snapshot.recovery === 'manual' ? 'Manual action required' : 'None'}</dd>{snapshot.recovery === 'automatic' && <><dt>Retry attempt</dt><dd>{snapshot.retryAttempt || 1}</dd><dt>Next retry</dt><dd>{snapshot.nextRetryAt || 'Pending'}</dd></>}</dl><Button disabled={busy || !snapshot.action} onClick={onRetry}>{busy ? 'Trying again' : snapshot.action ? 'Try again' : 'No retry needed'}</Button></section><section className="panel"><div className="section-heading"><div><p className="eyebrow">Saved targets</p><h2>Connection profiles</h2></div><Button variant="secondary" onClick={() => void load()}>Refresh profiles</Button></div><article><h3>Local connection</h3><p>Protected local IPC</p><Button disabled={!workspace?.activeProfileId} onClick={() => void select('')}>Select This computer</Button></article>{workspace?.profiles.map(profile => <article key={profile.id}><h3>{profile.label}{profile.active ? ' (selected)' : ''}</h3><p>{profile.endpoint} · {profile.shortDaemonId}</p><dl><dt>Capability</dt><dd>{profile.capability}</dd><dt>Platform</dt><dd>{profile.platform}{profile.architecture ? `/${profile.architecture}` : ''}</dd><dt>Certificate SHA-256</dt><dd>{profile.fingerprint}</dd><dt>Last successful connection</dt><dd>{profile.lastSuccessfulAt || 'Never'}</dd></dl><div className="actions"><Button disabled={profile.active} onClick={() => void select(profile.id)}>Select</Button><Button variant="secondary" onClick={() => void rename(profile.id, profile.label)}>Rename</Button><Button variant="secondary" onClick={() => setRepairProfileId(profile.id)}>Repair</Button><Button variant="danger" onClick={() => void remove(profile.id, profile.label)}>Remove</Button></div></article>)}</section><PairingForm repairProfileId={repairProfileId} onPaired={() => { setRepairProfileId(''); void load() }} /></>
}
