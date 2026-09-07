import { useState } from 'react'
import { Button, Notice, StatePanel } from './components'
import { Shell } from './components/Shell'
import { desktopBridge } from './connection/bridge'
import { useConnection } from './connection/store'
import type { Appearance, DesktopBridge, Route } from './connection/model'
import { taskBridge as nativeTaskBridge } from './tasks/bridge'
import type { TaskBridge } from './tasks/model'
import { TasksPage } from './tasks/TasksPage'

const copy: Record<Route, { title: string; detail: string }> = {
  tasks: { title: 'Tasks', detail: 'Create, schedule, organize, and run local work.' },
  schedule: { title: 'Schedule', detail: 'Upcoming run workflows are not migrated yet.' },
  activity: { title: 'Activity', detail: 'Run history workflows are not migrated yet.' },
  connections: { title: 'Connections', detail: 'This computer is the only connection in this release.' },
  settings: { title: 'Settings', detail: 'Preference migration is planned for a later slice.' },
}

export function App({ bridge = desktopBridge, tasks = nativeTaskBridge }: { bridge?: DesktopBridge; tasks?: TaskBridge }) {
  const [route, setRoute] = useState<Route>('tasks')
  const [appearance, setAppearance] = useState<Appearance>('system')
  const { snapshot, announcement, retry } = useConnection(bridge)
  const page = copy[route]
  return <Shell route={route} onRoute={setRoute} appearance={appearance} onAppearance={setAppearance} connection={snapshot} announcement={announcement} onRetry={() => void retry()} onQuit={() => void bridge.quit()}>
    {route === 'tasks' ? <>{snapshot.state !== 'connected' && <Notice title="This computer needs attention" tone="warning">{snapshot.message} {snapshot.action && <Button variant="secondary" onClick={() => void retry()}>Try again</Button>}</Notice>}<TasksPage bridge={tasks} platform={snapshot.target.platform} available={snapshot.state === 'connected'} refreshToken={snapshot.state === 'connected' ? snapshot.generation : 0} onActivity={() => setRoute('activity')} /></> : snapshot.state === 'connected' ? <><header className="page-header"><div><p className="eyebrow">This computer</p><h1>{page.title}</h1><p>{page.detail}</p></div></header><StatePanel title={`${page.title} is not migrated yet`} detail={page.detail} /></> : <><header className="page-header"><div><p className="eyebrow">This computer</p><h1>{page.title}</h1><p>{page.detail}</p></div></header><StatePanel title={snapshot.state === 'connecting' || snapshot.state === 'recovering' ? 'Connecting to This computer' : 'This computer needs attention'} detail={snapshot.message} busy={snapshot.state === 'connecting' || snapshot.state === 'recovering'} action={snapshot.action ? 'Try again' : undefined} onAction={() => void retry()} /></>}
  </Shell>
}
