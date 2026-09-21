import { useState } from 'react'
import { Button, CardSection, DataTable, DescriptionList, Dialog, Disclosure, Field, FormGrid, Notice, StatePanel, StatusLabel } from '../components'
import type { AgentAccessBridge, AgentGrant, AgentGrantDraft, AgentGrantEditDraft } from './model'
import { useAgentAccess } from './store'

const durationOptions = [
  ['1h', '1 hour'],
  ['24h', '24 hours'],
  ['7d', '7 days'],
  ['30d', '30 days'],
] as const

const capabilityLabel = (value: string) => value === 'manage' ? 'Manage' : value === 'operate' ? 'Operate' : 'Observe'
const transportLabel = (value: string) => value === 'remote_https' ? 'Remote HTTPS' : value === 'localhost_http' ? 'Localhost HTTP' : 'Stdio'
const timeLabel = (value?: string) => value ? new Date(value).toLocaleString() : 'Not recorded'
const stateTone = (state: string) => state === 'active' ? 'positive' : state === 'revoked' ? 'danger' : 'warning'

export function AgentAccessPage({ bridge, available, refreshToken }: { bridge: AgentAccessBridge; available: boolean; refreshToken: number }) {
  const access = useAgentAccess(bridge, available, refreshToken)
  const [name, setName] = useState('Codex')
  const [port, setPort] = useState('43123')
  const [origins, setOrigins] = useState('')
  const [permission, setPermission] = useState<'observe' | 'operate' | 'manage'>('observe')
  const [requireConfirmation, setRequireConfirmation] = useState(false)
  const [createOpen, setCreateOpen] = useState(false)
  const [createInvoker, setCreateInvoker] = useState<HTMLElement | null>(null)
  const [grantName, setGrantName] = useState('')
  const [grantCapability, setGrantCapability] = useState<AgentGrantDraft['capability']>('observe')
  const [grantDuration, setGrantDuration] = useState<AgentGrantDraft['duration']>('24h')
  const [nonExpiringConfirmed, setNonExpiringConfirmed] = useState(false)
  const [editGrant, setEditGrant] = useState<AgentGrant>()
  const [editInvoker, setEditInvoker] = useState<HTMLElement | null>(null)
  const [editCapability, setEditCapability] = useState<AgentGrantEditDraft['capability']>()
  const [editDuration, setEditDuration] = useState<AgentGrantEditDraft['duration']>()
  const [revokeGrant, setRevokeGrant] = useState<AgentGrant>()
  const [revokeInvoker, setRevokeInvoker] = useState<HTMLElement | null>(null)
  const [activityGrant, setActivityGrant] = useState<AgentGrant>()
  const [activityInvoker, setActivityInvoker] = useState<HTMLElement | null>(null)

  if (!available) return <StatePanel title="Agent Access unavailable" detail="Reconnect to the local scheduler service to inspect or control agent access." />
  if (!access.workspace && access.message) return <StatePanel title="Agent Access unavailable" detail={access.message} action="Try again" onAction={() => void access.load()} />
  if (!access.workspace) return <StatePanel title="Loading Agent Access" detail="Reading the local MCP access boundary." busy />

  const http = access.workspace.http
  const grants = access.workspace.grants ?? []
  const transports = access.workspace.transports ?? [
    { id: 'stdio' as const, name: 'Stdio', state: 'available_on_demand', description: access.workspace.stdioDescription },
    { id: 'localhost_http' as const, name: 'Localhost HTTP', state: http.enabled ? 'active' : 'off', description: 'Optional loopback listener.' },
  ]
  const daemon = access.workspace.daemon ?? { id: 'local', name: 'This computer' }

  const enableForm = <form onSubmit={(event) => { event.preventDefault(); void access.enable({ clientName: name, port: Number(port), allowedOrigins: origins.split(/\r?\n/).map((value) => value.trim()).filter(Boolean), permission, requireConfirmation: permission === 'manage' && requireConfirmation }) }}><FormGrid><Field label="Client name"><input value={name} maxLength={64} required onChange={(event) => setName(event.target.value)} /></Field><Field label="Permission"><select value={permission} onChange={(event) => setPermission(event.target.value as 'observe' | 'operate' | 'manage')}><option value="observe">Observe (read-only)</option><option value="operate">Operate (run, enable, and disable tasks)</option><option value="manage">Manage (change automation definitions)</option></select></Field><Field label="Loopback port"><input value={port} type="number" min="1" max="65535" required onChange={(event) => setPort(event.target.value)} /></Field><Field label="Allowed browser origins (one per line)"><textarea value={origins} placeholder="http://127.0.0.1:3000" onChange={(event) => setOrigins(event.target.value)} /></Field>{permission === 'manage' && <Field label="Manage safety"><label><input type="checkbox" checked={requireConfirmation} onChange={(event) => setRequireConfirmation(event.target.checked)} /> Require confirmed calls</label></Field>}</FormGrid><Button pending={access.pending} type="submit">Enable and copy credential</Button></form>

  const openEdit = (grant: AgentGrant, invoker: HTMLElement) => {
    setEditGrant(grant)
    setEditInvoker(invoker)
    setEditCapability(grant.capability)
    setEditDuration(undefined)
  }

  return <section aria-labelledby="agent-access-title">
    <header className="page-header"><div><p className="eyebrow">{daemon.name}</p><h1 id="agent-access-title">Agent Access</h1><p>See who can access this daemon, at what authority, through which transport, and until when.</p></div><div className="actions"><StatusLabel tone={access.workspace.mcpState === 'active' ? 'positive' : 'neutral'}>MCP {access.workspace.mcpState === 'active' ? 'On' : 'Off'}</StatusLabel><Button variant="secondary" onClick={() => void access.openGuide()}>Setup guide</Button><Button onClick={(event) => { setCreateInvoker(event.currentTarget); setCreateOpen(true) }}>Grant remote access</Button></div></header>
    {access.message && <Notice title="Agent Access status" tone={access.outcome === 'accepted' ? 'success' : 'warning'} identity={access.message}>{access.message}</Notice>}

    <div className="admin-stack">
      <div className="transport-grid" aria-label="MCP transports">{transports.map((transport) => <CardSection key={transport.id} eyebrow="Transport" title={transport.name} status={<StatusLabel tone={transport.state === 'active' || transport.state === 'available_on_demand' ? 'positive' : 'neutral'}>{transport.state === 'available_on_demand' ? 'Available on demand' : transport.state === 'active' ? 'Active' : 'Off'}</StatusLabel>}><p>{transport.description}</p></CardSection>)}</div>

      <CardSection eyebrow="Permission grants" title="Agents with access" status={<StatusLabel tone={grants.some((grant) => grant.state === 'active') ? 'positive' : 'neutral'}>{grants.filter((grant) => grant.state === 'active').length} active</StatusLabel>}>
        {grants.length === 0 ? <StatePanel title="No agent grants" detail="Create a remote grant or launch a local MCP client when access is needed." /> : <div className="grant-list">{grants.map((grant) => <article className={`grant-card grant-${grant.capability}`} key={grant.id}><header><div><p className="eyebrow">{transportLabel(grant.transport)}</p><h3>{grant.clientName}</h3></div><div className="grant-badges"><StatusLabel tone={grant.capability === 'manage' ? 'danger' : grant.capability === 'operate' ? 'warning' : 'neutral'}>{capabilityLabel(grant.capability)}</StatusLabel><StatusLabel tone={stateTone(grant.state)}>{grant.state === 'active' ? 'Active' : grant.state === 'expired' ? 'Expired' : 'Revoked'}</StatusLabel></div></header><p>{grant.capabilityDescription}</p><DescriptionList><dt>Daemon</dt><dd>{grant.daemonName} <code className="operational-value">{grant.daemonId}</code></dd><dt>Created</dt><dd>{timeLabel(grant.createdAt)}</dd><dt>Last use</dt><dd>{timeLabel(grant.lastUsedAt)}</dd><dt>Expires</dt><dd>{grant.expiresAt ? timeLabel(grant.expiresAt) : 'Non-expiring'}</dd>{grant.credentialFingerprint && <><dt>Credential fingerprint</dt><dd><code className="operational-value">{grant.credentialFingerprint}</code></dd></>}</DescriptionList><div className="actions"><Button variant="secondary" onClick={(event) => { setActivityGrant(grant); setActivityInvoker(event.currentTarget); access.clearActions(); void access.loadActions(grant.id) }}>Recent actions</Button>{grant.state === 'active' && <><Button variant="secondary" onClick={(event) => openEdit(grant, event.currentTarget)}>Narrow access</Button><Button variant="danger" onClick={(event) => { setRevokeGrant(grant); setRevokeInvoker(event.currentTarget) }}>Revoke</Button></>}</div></article>)}</div>}
      </CardSection>

      <div className="settings-grid"><CardSection eyebrow="Authority" title="What agents can do"><div className="authority-list">{access.workspace.authorities.map((authority) => <section className={`authority-${authority.name.toLowerCase()}`} key={authority.name}><h3>{authority.name}</h3><p>{authority.description}</p></section>)}</div></CardSection><CardSection eyebrow="Optional network access" title="Localhost HTTP" status={<StatusLabel tone={http.enabled ? 'positive' : 'neutral'}>{http.enabled ? 'Active' : 'Off'}</StatusLabel>}>{http.enabled ? <><DescriptionList><dt>Client</dt><dd>{http.clientName}</dd><dt>Permission</dt><dd>{capabilityLabel(http.permission ?? 'observe')}</dd>{http.permission === 'manage' && <><dt>Confirmation</dt><dd>{http.requireConfirmation ? 'Required on every mutation' : 'Not required'}</dd></>}<dt>Endpoint</dt><dd><code className="operational-value">{http.endpoint}</code></dd><dt>Allowed browser origins</dt><dd className="operational-value">{http.allowedOrigins.length ? http.allowedOrigins.join(', ') : 'None'}</dd><dt>Credential fingerprint</dt><dd><code className="operational-value">{http.credentialFingerprint}</code></dd><dt>Enabled</dt><dd>{timeLabel(http.enabledAt)}</dd><dt>Last successful access</dt><dd>{timeLabel(http.lastAccessedAt)}</dd><dt>Successful requests</dt><dd>{http.requestCount}</dd></DescriptionList><div className="actions"><Button disabled={access.pending} onClick={() => void access.rotate()}>Rotate credential</Button><Button variant="danger" disabled={access.pending} onClick={() => void access.revoke()}>Revoke localhost access</Button></div></> : <Disclosure summary="Configure localhost HTTP">{enableForm}</Disclosure>}</CardSection></div>
    </div>

    <Dialog open={createOpen} title="Grant remote MCP access" invoker={createInvoker} dismissDisabled={access.pending} onClose={() => setCreateOpen(false)} actions={<Button pending={access.pending} disabled={!grantName.trim() || (grantDuration === 'non-expiring' && !nonExpiringConfirmed)} onClick={() => { void access.createGrant({ clientName: grantName.trim(), capability: grantCapability, duration: grantDuration }); setCreateOpen(false) }}>Create and copy enrollment</Button>}><p>The listener remains a separate setting. This creates one named permission grant for {daemon.name}.</p><FormGrid><Field label="Client name"><input required maxLength={80} value={grantName} onChange={(event) => setGrantName(event.target.value)} /></Field><Field label="Authority"><select value={grantCapability} onChange={(event) => setGrantCapability(event.target.value as AgentGrantDraft['capability'])}><option value="observe">Observe, read-only context</option><option value="operate">Operate, run and toggle existing tasks</option><option value="manage">Manage, change automation definitions</option></select></Field><Field label="Access duration"><select value={grantDuration} onChange={(event) => { setGrantDuration(event.target.value as AgentGrantDraft['duration']); setNonExpiringConfirmed(false) }}>{durationOptions.map(([value, label]) => <option key={value} value={value}>{label}</option>)}<option value="non-expiring">Non-expiring (deliberate)</option></select></Field>{grantDuration === 'non-expiring' && <Field label="Non-expiring acknowledgement"><label><input type="checkbox" checked={nonExpiringConfirmed} onChange={(event) => setNonExpiringConfirmed(event.target.checked)} /> I intend this grant to remain active until revoked.</label></Field>}</FormGrid><p className="muted">The one-time enrollment bundle is copied natively and expires in 10 minutes. It is never displayed here.</p></Dialog>

    <Dialog open={Boolean(editGrant)} title={`Narrow ${editGrant?.clientName ?? 'agent'} access`} invoker={editInvoker} dismissDisabled={access.pending} onClose={() => setEditGrant(undefined)} actions={<Button pending={access.pending} onClick={() => { if (editGrant) void access.editGrant({ actorId: editGrant.id, capability: editCapability, duration: editDuration }); setEditGrant(undefined) }}>Apply narrower access</Button>}><p>Existing connections use the lower boundary on their next request. Higher or longer access requires a new grant.</p><FormGrid><Field label="Authority"><select value={editCapability} onChange={(event) => setEditCapability(event.target.value as AgentGrantEditDraft['capability'])}><option value="observe">Observe</option>{editGrant?.capability !== 'observe' && <option value="operate">Operate</option>}{editGrant?.capability === 'manage' && <option value="manage">Manage (unchanged)</option>}</select></Field><Field label="Earlier expiry"><select value={editDuration ?? ''} onChange={(event) => setEditDuration((event.target.value || undefined) as AgentGrantEditDraft['duration'])}><option value="">Keep current expiry</option>{durationOptions.map(([value, label]) => <option key={value} value={value}>{label} from now</option>)}</select></Field></FormGrid></Dialog>

    <Dialog open={Boolean(revokeGrant)} title={`Revoke ${revokeGrant?.clientName ?? 'agent'}?`} invoker={revokeInvoker} dismissDisabled={access.pending} onClose={() => setRevokeGrant(undefined)} actions={<Button variant="danger" pending={access.pending} onClick={() => { if (revokeGrant) void access.revokeGrant(revokeGrant.id); setRevokeGrant(undefined) }}>Confirm revoke</Button>}><p>This permanently denies the grant. Existing connections are rejected on their next request.</p></Dialog>

    <Dialog open={Boolean(activityGrant)} title={`${activityGrant?.clientName ?? 'Agent'} recent actions`} invoker={activityInvoker} dismissDisabled={access.pending} onClose={() => { setActivityGrant(undefined); access.clearActions() }}>{access.pending && !access.actions ? <StatePanel title="Loading recent actions" detail="Reading shared audit evidence." busy /> : access.actions?.outcome !== 'accepted' ? <StatePanel title="Recent actions unavailable" detail={access.actions?.message ?? 'Recent actions could not be loaded.'} /> : <DataTable caption={`Newest actions for ${activityGrant?.clientName ?? 'agent'}`} headings={['Time', 'Action', 'Target', 'Result', 'Daemon']} rows={access.actions.actions.map((action) => [timeLabel(action.occurredAt), <code>{action.operation}</code>, `${action.targetKind}${action.targetId ? ` ${action.targetId}` : ''}`, action.result, <code>{action.daemonId}</code>])} />}</Dialog>
  </section>
}
