import { useEffect, useState } from 'react'
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
import { ConnectionsPage } from './settings/ConnectionsPage'
import { SettingsPage } from './settings/SettingsPage'
import { settingsBridge as nativeSettingsBridge } from './settings/bridge'
import type { SettingsBridge } from './settings/model'
import { useSettings } from './settings/store'
import { NotificationsPage } from './notifications/NotificationsPage'
import { notificationBridge as nativeNotificationBridge } from './notifications/bridge'
import type { NotificationBridge } from './notifications/model'
import { AgentAccessPage } from './agentaccess/AgentAccessPage'
import { agentAccessBridge as nativeAgentAccessBridge } from './agentaccess/bridge'
import type { AgentAccessBridge } from './agentaccess/model'

const copy: Record<Route, { title: string; detail: string }> = {
  tasks: { title: 'Tasks', detail: 'Create, schedule, organize, and run work on the selected scheduler.' },
  automation: { title: 'Automation Sources', detail: 'Connect events and external sources to local work.' },
  schedule: { title: 'Schedule', detail: 'Predicted work and recorded runs.' },
  activity: { title: 'Activity', detail: 'Recent runs, daemon logs, and alerts.' },
  notifications: { title: 'Notifications', detail: 'Webhook channels, outcome policies, and delivery evidence.' },
  agentAccess: { title: 'Agent Access', detail: 'Local MCP availability, authority, and revocation.' },
  connections: { title: 'Connections', detail: 'Local scheduler diagnosis and recovery.' },
  settings: { title: 'Settings', detail: 'Desktop preferences, storage ownership, and product information.' },
}

export function App({ bridge = desktopBridge, tasks = nativeTaskBridge, automation = nativeAutomationBridge, operations = nativeOperationsBridge, notifications = nativeNotificationBridge, settings = nativeSettingsBridge, agentAccess = nativeAgentAccessBridge }: { bridge?: DesktopBridge; tasks?: TaskBridge; automation?: AutomationBridge; operations?: OperationsBridge; notifications?: NotificationBridge; settings?: SettingsBridge; agentAccess?: AgentAccessBridge }) {
  const [route, setRoute] = useState<Route>('tasks')
  const [appearance, setAppearance] = useState<Appearance>('system')
  const [toast, setToast] = useState('')
  const { snapshot, announcement, retryPending, retry } = useConnection(bridge)
  const desktopSettings = useSettings(settings, `${snapshot.generation}:${snapshot.state}`)
  useEffect(() => { if (desktopSettings.workspace) setAppearance(desktopSettings.workspace.preferences.appearance) }, [desktopSettings.workspace])
  useEffect(() => { if (announcement) setToast(announcement) }, [announcement])
  useEffect(() => { if (desktopSettings.message) setToast(desktopSettings.message) }, [desktopSettings.message])
  const page = copy[route]
  const remote = snapshot.target.kind === 'remote'
  const unsupportedRemote = remote && (route === 'automation' || route === 'notifications' || route === 'agentAccess')
  const canOperate = snapshot.state === 'connected' && (!remote || snapshot.target.permissions.includes('operate'))
  const canManage = snapshot.state === 'connected' && (!remote || snapshot.target.permissions.includes('manage'))
  const targetContext = remote ? `${snapshot.target.displayName} (${snapshot.target.endpoint ?? 'remote'}, ${snapshot.target.id.slice(0, 8)})` : 'This computer'
  const saveAppearance = (value: Appearance) => { void desktopSettings.saveAppearance(value) }
  return <Shell route={route} onRoute={setRoute} appearance={appearance} onAppearance={saveAppearance} appearancePending={desktopSettings.pending} connection={snapshot} announcement={toast} onRetry={() => void retry()} onQuit={() => void bridge.quit()}>
    {desktopSettings.loading ? <StatePanel title="Loading desktop preferences" detail="Preparing your local desktop settings." busy /> : unsupportedRemote ? <><header className="page-header"><div><p className="eyebrow">{snapshot.target.displayName}</p><h1>{page.title}</h1></div></header><StatePanel title={`${page.title} is unavailable for remote targets`} detail="This feature is outside the authenticated remote operation allowlist. Select This computer to use it." /></> : route === 'tasks' ? <>{snapshot.state !== 'connected' && <Notice title={`${snapshot.target.displayName} needs attention`} tone="warning">{snapshot.message} {snapshot.action && <Button variant="secondary" onClick={() => void retry()}>Try again</Button>}</Notice>}{remote && !canOperate && <Notice title="Observe-only connection" tone="info">Task operations are disabled because this credential does not grant operate authority.</Notice>}{remote && canOperate && !canManage && <Notice title="Operate-only connection" tone="info">Task creation, editing, enablement, and deletion are disabled because this credential does not grant manage authority.</Notice>}<TasksPage bridge={tasks} platform={snapshot.target.platform} available={canOperate} workspaceAvailable={snapshot.state === 'connected'} manageAvailable={canManage} editAvailable={!remote && canManage} groupAvailable={!remote && canManage} previewAvailable={!remote && canManage} targetName={targetContext} refreshToken={snapshot.state === 'connected' ? snapshot.generation : 0} onActivity={() => setRoute('activity')} /></> : route === 'automation' ? <AutomationPage bridge={automation} available={snapshot.state === 'connected'} refreshToken={snapshot.state === 'connected' ? snapshot.generation : 0} /> : route === 'schedule' ? <SchedulePage bridge={operations} available={snapshot.state === 'connected'} targetName={targetContext} refreshToken={snapshot.state === 'connected' ? snapshot.generation : 0} /> : route === 'activity' ? <ActivityPage bridge={operations} available={snapshot.state === 'connected'} targetName={targetContext} canMutate={canOperate} refreshToken={snapshot.state === 'connected' ? snapshot.generation : 0} /> : route === 'notifications' ? <NotificationsPage bridge={notifications} available={snapshot.state === 'connected'} refreshToken={snapshot.state === 'connected' ? snapshot.generation : 0} /> : route === 'agentAccess' ? <AgentAccessPage bridge={agentAccess} available={snapshot.state === 'connected'} refreshToken={snapshot.state === 'connected' ? snapshot.generation : 0} /> : route === 'connections' ? <ConnectionsPage snapshot={snapshot} retryPending={retryPending} onRetry={() => void retry()} bridge={bridge} /> : route === 'settings' ? <SettingsPage workspace={desktopSettings.workspace} message={desktopSettings.message} pending={desktopSettings.pending} onAppearance={saveAppearance} onRestore={() => void desktopSettings.restore()} onCopy={(id) => void desktopSettings.copyStoragePath(id)} onOpen={(key) => void desktopSettings.openProductLink(key)} onConnections={() => setRoute('connections')} /> : <><header className="page-header"><div><p className="eyebrow">{snapshot.target.displayName}</p><h1>{page.title}</h1><p>{page.detail}</p></div></header><StatePanel title="Workspace unavailable" detail={page.detail} /></>}
  </Shell>
}
