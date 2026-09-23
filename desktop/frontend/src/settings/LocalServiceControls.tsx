import { Button, CardSection, Notice, StatusLabel } from "../components";
import type { LocalServiceSnapshot } from "../connection/model";

export function LocalServiceControls({
  snapshot,
  pending,
  message,
  onAction,
}: {
  snapshot?: LocalServiceSnapshot;
  pending: boolean;
  message: string;
  onAction(action: "start" | "stop" | "restart"): void;
}) {
  if (!snapshot || snapshot.state === "unsupported") return null;
  const canStart = snapshot.state === "stopped";
  const canStop = snapshot.scmState === "running";
  return (
    <CardSection
      eyebrow="This computer"
      title="Local service"
      status={
        <StatusLabel tone={snapshot.state === "running" ? "positive" : "warning"}>
          {snapshot.state.replaceAll("_", " ")}
        </StatusLabel>
      }
    >
      <p>{snapshot.detail}</p>
      {message && <Notice title="Service action" tone="info" dismissible={false}>{message}</Notice>}
      <div className="actions">
        {canStart && <Button pending={pending} onClick={() => onAction("start")}>Start service</Button>}
        {canStop && <>
          <Button variant="secondary" pending={pending} onClick={() => onAction("restart")}>Restart service</Button>
          <Button variant="danger" pending={pending} onClick={() => onAction("stop")}>Stop service</Button>
        </>}
      </div>
      {snapshot.state === "not_installed" && <p>Install the system service to manage it here.</p>}
      {snapshot.state === "unreachable" && <p>Inspect the running service in your system service manager, then retry its local connection.</p>}
    </CardSection>
  );
}
