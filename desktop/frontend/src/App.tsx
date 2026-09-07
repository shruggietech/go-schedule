import { useState } from 'react'
import { Button, Notice, StatePanel } from './components'
import { Shell } from './components/Shell'
import { desktopBridge } from './connection/bridge'
import { useConnection } from './connection/store'
import type { Appearance, DesktopBridge, Route } from './connection/model'
import { taskBridge as nativeTaskBridge } from './tasks/bridge'
import type { TaskBridge } from './tasks/model'
import { TasksPage } from './tasks/TasksPage'
import { AutomationPage } from './automation/AutomationPage'
import { automationBridge as nativeAutomationBridge } from './automation/bridge'
import type { AutomationBridge } from './automation/model'
import { ActivityPage } from './operations/ActivityPage'
import { operationsBridge as nativeOperationsBridge } from './operations/bridge'
import type { OperationsBridge } from './operations/model'
import { SchedulePage } from './operations/SchedulePage'

const copy: Record<Route, { title: string; detail: string }> = {
  tasks: { title: 'Tasks', detail: 'Create, schedule, organize, and run local work.' },
  automation: { title: 'Automation Sources', detail: 'Connect events and external sources to local work.' },
  schedule: { title: 'Schedule', detail: 'Predicted work and recorded runs.' },
  activity: { title: 'Activity', detail: 'Recent runs, daemon logs, and alerts.' },
  connections: { title: 'Connections', detail: 'This computer is the only connection in this release.' },
  settings: { title: 'Settings', detail: 'Preference migration is planned for a later slice.' },
}

export function App({ bridge = desktopBridge, tasks = nativeTaskBridge, automation = nativeAutomationBridge, operations = nativeOperationsBridge }: { bridge?: DesktopBridge; tasks?: TaskBridge; automation?: AutomationBridge; operations?: OperationsBridge }) {
  const [route, setRoute] = useState<Route>('tasks')
  const [appearance, setAppearance] = useState<Appearance>('system')
  const { snapshot, announcement, retry } = useConnection(bridge)
  const page = copy[route]
  return <Shell route={route} onRoute={setRoute} appearance={appearance} onAppearance={setAppearance} connection={snapshot} announcement={announcement} onRetry={() => void retry()} onQuit={() => void bridge.quit()}>
    {route === 'tasks' ? <>{snapshot.state !== 'connected' && <Notice title="This computer needs attention" tone="warning">{snapshot.message} {snapshot.action && <Button variant="secondary" onClick={() => void retry()}>Try again</Button>}</Notice>}<TasksPage bridge={tasks} platform={snapshot.target.platform} available={snapshot.state === 'connected'} refreshToken={snapshot.state === 'connected' ? snapshot.generation : 0} onActivity={() => setRoute('activity')} /></> : route === 'automation' ? <AutomationPage bridge={automation} available={snapshot.state === 'connected'} refreshToken={snapshot.state === 'connected' ? snapshot.generation : 0} /> : route === 'schedule' ? <SchedulePage bridge={operations} available={snapshot.state === 'connected'} refreshToken={snapshot.state === 'connected' ? snapshot.generation : 0} /> : route === 'activity' ? <ActivityPage bridge={operations} available={snapshot.state === 'connected'} refreshToken={snapshot.state === 'connected' ? snapshot.generation : 0} /> : snapshot.state === 'connected' ? <><header className="page-header"><div><p className="eyebrow">This computer</p><h1>{page.title}</h1><p>{page.detail}</p></div></header><StatePanel title={`${page.title} is not migrated yet`} detail={page.detail} /></> : <><header className="page-header"><div><p className="eyebrow">This computer</p><h1>{page.title}</h1><p>{page.detail}</p></div></header><StatePanel title={snapshot.state === 'connecting' || snapshot.state === 'recovering' ? 'Connecting to This computer' : 'This computer needs attention'} detail={snapshot.message} busy={snapshot.state === 'connecting' || snapshot.state === 'recovering'} action={snapshot.action ? 'Try again' : undefined} onAction={() => void retry()} /></>}
  </Shell>
}
