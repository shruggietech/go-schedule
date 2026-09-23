import { useEffect, useRef, useState } from "react";
import { Button, Notice, StatePanel } from "./components";
import { Shell, type ShellFeedback } from "./components/Shell";
import { desktopBridge } from "./connection/bridge";
import { useConnection } from "./connection/store";
import type { Appearance, ConnectionSnapshot, DesktopBridge, LocalServiceSnapshot, Route } from "./connection/model";
import { taskBridge as nativeTaskBridge } from "./tasks/bridge";
import type { TaskBridge } from "./tasks/model";
import { TasksPage } from "./tasks/TasksPage";
import { AutomationPage } from "./automation/AutomationPage";
import { automationBridge as nativeAutomationBridge } from "./automation/bridge";
import type { AutomationBridge } from "./automation/model";
import { ActivityPage } from "./operations/ActivityPage";
import { operationsBridge as nativeOperationsBridge } from "./operations/bridge";
import type { OperationsBridge } from "./operations/model";
import { SchedulePage } from "./operations/SchedulePage";
import { ConnectionsPage } from "./settings/ConnectionsPage";
import { SettingsPage } from "./settings/SettingsPage";
import { settingsBridge as nativeSettingsBridge } from "./settings/bridge";
import type { SettingsBridge } from "./settings/model";
import { useSettings } from "./settings/store";
import { NotificationsPage } from "./notifications/NotificationsPage";
import { notificationBridge as nativeNotificationBridge } from "./notifications/bridge";
import type { NotificationBridge } from "./notifications/model";
import { AgentAccessPage } from "./agentaccess/AgentAccessPage";
import { agentAccessBridge as nativeAgentAccessBridge } from "./agentaccess/bridge";
import type { AgentAccessBridge } from "./agentaccess/model";
import { SystemsPage, type Drilldown } from "./systems/SystemsPage";
import { SearchPage } from "./search/SearchPage";
import { searchBridge as nativeSearchBridge } from "./search/bridge";
import type { SearchBridge } from "./search/model";
import { BundlesPage } from "./bundles/BundlesPage";
import { bundleBridge as nativeBundleBridge } from "./bundles/bridge";
import type { BundleBridge } from "./bundles/model";

const copy: Record<Route, { title: string; detail: string }> = {
  systems: {
    title: "All Systems",
    detail: "Bounded operational observations from every registered scheduler.",
  },
  search: {
    title: "Search",
    detail: "Find automation and operational evidence across registered schedulers.",
  },
  tasks: {
    title: "Tasks",
    detail:
      "Create, schedule, organize, and run work on the selected scheduler.",
  },
  automation: {
    title: "Automation Sources",
    detail: "Connect events and external sources to local work.",
  },
	 bundles: { title: "Portable bundles", detail: "Export, compare, and apply safe automation intent on one selected scheduler." },
  schedule: { title: "Schedule", detail: "Predicted work and recorded runs." },
  activity: {
    title: "Activity",
    detail: "Recent runs, daemon logs, and alerts.",
  },
  notifications: {
    title: "Notifications",
    detail: "Webhook channels, outcome policies, and delivery evidence.",
  },
  agentAccess: {
    title: "Agent Access",
    detail: "Local MCP availability, authority, and revocation.",
  },
  connections: {
    title: "Connections",
    detail: "Local scheduler diagnosis and recovery.",
  },
  settings: {
    title: "Settings",
    detail: "Desktop preferences, storage ownership, and product information.",
  },
};

const reconnectTimeout = 6500;

function selectedTarget(snapshot: ConnectionSnapshot, registrationKey: string) {
  return registrationKey === "local"
    ? snapshot.target.kind !== "remote"
    : snapshot.target.profileId === registrationKey;
}

async function selectAndWaitForTarget(bridge: DesktopBridge, registrationKey: string, expectedDaemonId?: string): Promise<ConnectionSnapshot> {
  return new Promise((resolve, reject) => {
    let settled = false;
    let unsubscribe: () => void = () => undefined;
    let timeout = 0;
    const finish = (snapshot?: ConnectionSnapshot, message?: string) => {
      if (settled) return;
      settled = true;
      clearTimeout(timeout);
      unsubscribe();
      if (snapshot) resolve(snapshot);
      else reject(new Error(message ?? "The scheduler did not reconnect in time. Check the connection and try again."));
    };
    const inspect = (snapshot: ConnectionSnapshot) => {
      if (!selectedTarget(snapshot, registrationKey)) return;
      if (snapshot.state === "connected" && expectedDaemonId && snapshot.target.id !== expectedDaemonId) finish(undefined, "The scheduler identity changed. Search was left open and no record was opened.");
      else if (snapshot.state === "connected") finish(snapshot);
      else if (snapshot.state !== "connecting" && snapshot.state !== "recovering") finish(undefined, `${snapshot.message}${snapshot.action ? ` ${snapshot.action}` : ""}`);
    };
    unsubscribe = bridge.subscribe((event) => { if (event.snapshot) inspect(event.snapshot); });
    timeout = window.setTimeout(() => finish(), reconnectTimeout);
    const request = bridge.selectConnection?.(registrationKey === "local" ? "" : registrationKey);
    if (!request) {
      finish(undefined, "This scheduler could not be selected.");
      return;
    }
    void request.then((selected) => {
      if (!selected || selected.outcome !== "accepted") {
        finish(undefined, selected?.message ?? "This scheduler could not be selected.");
        return;
      }
      void bridge.snapshot().then(inspect).catch(() => undefined);
    }).catch(() => finish(undefined, "This scheduler could not be selected."));
  });
}

export function App({
  bridge = desktopBridge,
  tasks = nativeTaskBridge,
  automation = nativeAutomationBridge,
  operations = nativeOperationsBridge,
  notifications = nativeNotificationBridge,
  settings = nativeSettingsBridge,
  agentAccess = nativeAgentAccessBridge,
  search = nativeSearchBridge,
	 bundles = nativeBundleBridge,
}: {
  bridge?: DesktopBridge;
  tasks?: TaskBridge;
  automation?: AutomationBridge;
  operations?: OperationsBridge;
  notifications?: NotificationBridge;
  settings?: SettingsBridge;
  agentAccess?: AgentAccessBridge;
  search?: SearchBridge;
	 bundles?: BundleBridge;
}) {
  const [route, setRoute] = useState<Route>("tasks");
  const [appearance, setAppearance] = useState<Appearance>("system");
  const [announcementFeedback, setAnnouncementFeedback] =
    useState<ShellFeedback>();
  const [settingsError, setSettingsError] = useState<ShellFeedback>();
  const [drilldown, setDrilldown] = useState<Drilldown>();
  const [localService, setLocalService] = useState<LocalServiceSnapshot>();
  const [localServicePending, setLocalServicePending] = useState(false);
  const [localServiceMessage, setLocalServiceMessage] = useState("");
  const feedbackSequence = useRef(0);
  const { snapshot, announcement, retryPending, retry } = useConnection(bridge);
  useEffect(() => {
    let active = true;
    const refresh = () => {
      void bridge.localServiceSnapshot?.().then((value) => {
        if (active) setLocalService(value);
      }).catch(() => {
        if (active) setLocalService({ state: "unknown", scmState: "unknown", detail: "Could not read local Windows service status.", observedAt: "" });
      });
    };
    refresh();
    const timer = window.setInterval(refresh, 4000);
    return () => { active = false; window.clearInterval(timer); };
  }, [bridge]);
  const localServiceAction = async (action: "start" | "stop" | "restart") => {
    if (action === "stop" && !window.confirm("Stop the local service? Scheduled tasks on this computer will cease until it is started again.")) return;
    if (action === "restart" && !window.confirm("Restart the local service? Active tasks may be interrupted.")) return;
    setLocalServicePending(true);
    setLocalServiceMessage(`Requesting ${action} through Windows service control...`);
    try {
      const result = await bridge.controlLocalService?.(action, true);
      if (result) {
        setLocalService(result.snapshot);
        setLocalServiceMessage(result.message);
      }
    } catch {
      setLocalServiceMessage("The service action could not be completed. Check Windows Services and try again.");
    } finally {
      setLocalServicePending(false);
      void bridge.localServiceSnapshot?.().then(setLocalService);
    }
  };
  const desktopSettings = useSettings(
    settings,
    `${snapshot.generation}:${snapshot.state}`,
  );
  useEffect(() => {
    if (desktopSettings.workspace)
      setAppearance(desktopSettings.workspace.preferences.appearance);
  }, [desktopSettings.workspace]);
  useEffect(() => {
    if (!announcement) return;
    feedbackSequence.current += 1;
    setAnnouncementFeedback({
      identity: `connection:${feedbackSequence.current}`,
      message: announcement,
      tone: "success",
    });
  }, [announcement]);
  useEffect(() => {
    if (!desktopSettings.status) return;
    const next: ShellFeedback = {
      identity: `settings:${desktopSettings.status.id}`,
      message: desktopSettings.status.message,
      tone: desktopSettings.status.outcome === "accepted" ? "success" : "error",
    };
    if (next.tone === "error") setSettingsError(next);
    else {
      setSettingsError(undefined);
      setAnnouncementFeedback(next);
    }
  }, [desktopSettings.status]);
  const page = copy[route];
  const remote = snapshot.target.kind === "remote";
  const unsupportedRemote =
    remote &&
    (route === "automation" ||
      route === "notifications" ||
      route === "agentAccess");
  const canOperate =
    snapshot.state === "connected" &&
    (!remote || snapshot.target.permissions.includes("operate"));
  const canManage =
    snapshot.state === "connected" &&
    (!remote || snapshot.target.permissions.includes("manage"));
  const targetContext = remote
    ? `${snapshot.target.displayName} (${snapshot.target.endpoint ?? "remote"}, ${snapshot.target.id.slice(0, 8)})`
    : "This computer";
  const saveAppearance = (value: Appearance) => {
    void desktopSettings.saveAppearance(value);
  };
  const openSystem = async (intent: Drilldown) => {
    feedbackSequence.current += 1;
    try {
      await selectAndWaitForTarget(bridge, intent.registrationKey, intent.expectedDaemonId);
    } catch (error) {
      setAnnouncementFeedback({ identity: `system:${feedbackSequence.current}`, message: error instanceof Error ? error.message : "This scheduler could not be reconnected.", tone: "error" });
      return;
    }
    setDrilldown(intent);
    setAnnouncementFeedback({ identity: `system:${feedbackSequence.current}`, message: `${intent.label} selected. Opened ${intent.context}.`, tone: "success" });
    setRoute(intent.destination);
  };
  return (
    <Shell
      route={route}
      onRoute={setRoute}
      appearance={appearance}
      onAppearance={saveAppearance}
      appearancePending={desktopSettings.pendingActions.has("appearance")}
      connection={snapshot}
      localService={localService}
      announcement={announcementFeedback}
      persistentFeedback={settingsError}
      onRetry={() => void retry()}
      onQuit={() => void bridge.quit()}
    >
      {localService?.state === "stopped" && (
        <Notice title="Local daemon stopped" tone="warning" dismissible={false}>
          Scheduled tasks on this computer are not running.{" "}
          <Button pending={localServicePending} onClick={() => void localServiceAction("start")}>Start service</Button>
        </Notice>
      )}
      {drilldown && route === drilldown.destination && <Notice title={`Opened from ${drilldown.label}`} tone="info" dismissible={false}>{drilldown.context}{drilldown.taskId ? ` · Task ${drilldown.taskId}` : ""}{drilldown.recordId ? ` · Record ${drilldown.recordId}` : ""}</Notice>}
      {desktopSettings.loading ? (
        <StatePanel
          title="Loading desktop preferences"
          detail="Preparing your local desktop settings."
          busy
        />
      ) : unsupportedRemote ? (
        <>
          <header className="page-header">
            <div>
              <p className="eyebrow">{snapshot.target.displayName}</p>
              <h1>{page.title}</h1>
            </div>
          </header>
          <StatePanel
            title={`${page.title} is unavailable for remote targets`}
            detail="This feature is outside the authenticated remote operation allowlist. Select This computer to use it."
          />
        </>
      ) : route === "systems" ? (
        <SystemsPage bridge={bridge} onOpen={openSystem} />
      ) : route === "search" ? (
        <SearchPage bridge={search} onOpen={openSystem} />
      ) : route === "tasks" ? (
        <>
          {snapshot.state !== "connected" && (
            <Notice
              title={`${snapshot.target.displayName} needs attention`}
              tone="warning"
              identity={`${snapshot.generation}:${snapshot.revision}:${snapshot.state}:${snapshot.message}`}
            >
              {snapshot.message}{" "}
              {snapshot.action && (
                <Button variant="secondary" onClick={() => void retry()}>
                  Try again
                </Button>
              )}
            </Notice>
          )}
          {remote && !canOperate && (
            <Notice title="Observe-only connection" tone="info">
              Task operations are disabled because this credential does not
              grant operate authority.
            </Notice>
          )}
          {remote && canOperate && !canManage && (
            <Notice title="Operate-only connection" tone="info">
              Task creation, editing, enablement, and deletion are disabled
              because this credential does not grant manage authority.
            </Notice>
          )}
          <TasksPage
            bridge={tasks}
            platform={snapshot.target.platform}
            available={canOperate}
            workspaceAvailable={snapshot.state === "connected"}
            manageAvailable={canManage}
            editAvailable={!remote && canManage}
            groupAvailable={!remote && canManage}
            previewAvailable={!remote && canManage}
            targetName={targetContext}
            initialTaskId={drilldown?.destination === "tasks" ? drilldown.taskId : undefined}
            refreshToken={
              snapshot.state === "connected" ? snapshot.generation : 0
            }
            onActivity={() => setRoute("activity")}
          />
        </>
      ) : route === "automation" ? (
        <AutomationPage
          bridge={automation}
          available={snapshot.state === "connected"}
          refreshToken={
            snapshot.state === "connected" ? snapshot.generation : 0
          }
        />
	  ) : route === "bundles" ? (
		<BundlesPage bridge={bundles} targetName={targetContext} manageAvailable={canManage} />
      ) : route === "schedule" ? (
        <SchedulePage
          bridge={operations}
          available={snapshot.state === "connected"}
          targetName={targetContext}
          initialRecordId={drilldown?.destination === "schedule" ? drilldown.recordId : undefined}
		  initialTaskId={drilldown?.destination === "schedule" ? drilldown.taskId : undefined}
		  initialOccurredAt={drilldown?.destination === "schedule" ? drilldown.occurredAt : undefined}
          refreshToken={
            snapshot.state === "connected" ? snapshot.generation : 0
          }
        />
      ) : route === "activity" ? (
        <ActivityPage
          bridge={operations}
          available={snapshot.state === "connected"}
          targetName={targetContext}
          canMutate={canOperate}
          initialRecordId={drilldown?.destination === "activity" ? drilldown.recordId : undefined}
		  initialRecordKind={drilldown?.destination === "activity" ? drilldown.recordKind : undefined}
          refreshToken={
            snapshot.state === "connected" ? snapshot.generation : 0
          }
        />
      ) : route === "notifications" ? (
        <NotificationsPage
          bridge={notifications}
          available={snapshot.state === "connected"}
          refreshToken={
            snapshot.state === "connected" ? snapshot.generation : 0
          }
        />
      ) : route === "agentAccess" ? (
        <AgentAccessPage
          bridge={agentAccess}
          available={snapshot.state === "connected"}
          refreshToken={
            snapshot.state === "connected" ? snapshot.generation : 0
          }
        />
      ) : route === "connections" ? (
        <ConnectionsPage
          snapshot={snapshot}
          retryPending={retryPending}
          onRetry={() => void retry()}
          bridge={bridge}
          localService={localService}
          localServicePending={localServicePending}
          localServiceMessage={localServiceMessage}
          onLocalServiceAction={(action) => void localServiceAction(action)}
        />
      ) : route === "settings" ? (
        <SettingsPage
          workspace={desktopSettings.workspace}
          message={desktopSettings.message}
          pending={desktopSettings.pending}
          pendingActions={desktopSettings.pendingActions}
          copyResults={desktopSettings.copyResults}
          onAppearance={saveAppearance}
          onRestore={() => void desktopSettings.restore()}
          onCopy={(id) => void desktopSettings.copyStoragePath(id)}
          onOpen={(key) => void desktopSettings.openProductLink(key)}
          onConnections={() => setRoute("connections")}
        />
      ) : (
        <>
          <header className="page-header">
            <div>
              <p className="eyebrow">{snapshot.target.displayName}</p>
              <h1>{page.title}</h1>
              <p>{page.detail}</p>
            </div>
          </header>
          <StatePanel title="Workspace unavailable" detail={page.detail} />
        </>
      )}
    </Shell>
  );
}
