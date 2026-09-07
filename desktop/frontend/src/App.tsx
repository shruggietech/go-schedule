import { useState } from 'react'
import { Button, DataTable, Disclosure, Field, Notice, StatePanel } from './components'
import { Shell } from './components/Shell'
import { desktopBridge } from './connection/bridge'
import { useConnection } from './connection/store'
import type { Appearance, DesktopBridge, Route } from './connection/model'

const copy: Record<Route, { title: string; detail: string }> = {
  tasks: { title: 'Tasks', detail: 'Task workflows arrive in the next migration slice.' },
  schedule: { title: 'Schedule', detail: 'Upcoming run workflows are not migrated yet.' },
  activity: { title: 'Activity', detail: 'Run history workflows are not migrated yet.' },
  connections: { title: 'Connections', detail: 'This computer is the only connection in this release.' },
  settings: { title: 'Settings', detail: 'Preference migration is planned for a later slice.' },
}

export function App({ bridge = desktopBridge }: { bridge?: DesktopBridge }) {
  const [route, setRoute] = useState<Route>('tasks')
  const [appearance, setAppearance] = useState<Appearance>('system')
  const { snapshot, announcement, retry } = useConnection(bridge)
  const page = copy[route]
  return <Shell route={route} onRoute={setRoute} appearance={appearance} onAppearance={setAppearance} connection={snapshot} announcement={announcement} onRetry={() => void retry()} onQuit={() => void bridge.quit()}>
    <header className="page-header"><div><p className="eyebrow">This computer</p><h1>{page.title}</h1><p>{page.detail}</p></div><Button disabled>Create task</Button></header>
    {snapshot.state === 'connected' ? <section className="catalog panel" aria-labelledby="foundation-title"><h2 id="foundation-title">Desktop foundation ready</h2><Notice title="Migration boundary" tone="info">The shell is connected. Feature workflows remain in the current application until their migration slices are complete.</Notice><Field label="Example field" help="Shared field help and validation remain associated."><input disabled value="Available to migrated screens" readOnly /></Field><DataTable caption="Component contract example" headings={['Pattern', 'Status']} rows={[[<>Table, empty, and loading states</>, <>Ready</>]]} /><Disclosure summary="What is available now?">Connection identity, lifecycle recovery, navigation, appearance, notices, forms, tables, dialogs, and notifications.</Disclosure></section> : <StatePanel title={snapshot.state === 'connecting' || snapshot.state === 'recovering' ? 'Connecting to This computer' : 'This computer needs attention'} detail={snapshot.message} busy={snapshot.state === 'connecting' || snapshot.state === 'recovering'} action={snapshot.action ? 'Try again' : undefined} onAction={() => void retry()} />}
  </Shell>
}
