import { Button, StatusBadge } from '../components'
import type { ConnectionSnapshot } from '../connection/model'

const guidance: Record<ConnectionSnapshot['state'], string> = {
  connecting: 'Wait while go-schedule checks the protected local endpoint.',
  connected: 'The local scheduler service is ready. No recovery action is needed.',
  degraded: 'The last connection was interrupted. Current data remains visible while recovery continues.',
  recovering: 'A bounded retry is in progress. You can keep using desktop-local Settings.',
  unavailable: 'Start or repair the local scheduler service, then try again.',
  access_denied: 'Check that this account has permission to use the installed scheduler service.',
  incompatible: 'Update the desktop application and scheduler service to compatible versions.',
  timed_out: 'Check whether the scheduler service is responsive, then try again.',
}

export function ConnectionsPage({ snapshot, retryPending, onRetry }: { snapshot: ConnectionSnapshot; retryPending: boolean; onRetry(): void }) {
  const busy = snapshot.state === 'connecting' || snapshot.state === 'recovering' || retryPending
  return <><header className="page-header"><div><p className="eyebrow">This computer</p><h1>Connections</h1><p>Local scheduler diagnosis and recovery.</p></div><StatusBadge state={snapshot.state} /></header><section className="panel connection-workspace" aria-busy={busy}><div><p className="eyebrow">Current diagnosis</p><h2>{snapshot.target.displayName}</h2><p>{snapshot.message}</p><p>{guidance[snapshot.state]}</p></div><dl><dt>Platform</dt><dd>{snapshot.target.platform}</dd><dt>Version</dt><dd>{snapshot.target.version || 'Unavailable'}</dd><dt>Capabilities</dt><dd>{snapshot.target.capabilities.length ? snapshot.target.capabilities.join(', ') : 'Unavailable'}</dd><dt>Last successful connection</dt><dd>{snapshot.lastSuccessfulAt || 'Not available in this session'}</dd></dl><Button disabled={busy || !snapshot.action} onClick={onRetry}>{busy ? 'Trying again' : snapshot.action ? 'Try again' : 'No retry needed'}</Button></section></>
}
