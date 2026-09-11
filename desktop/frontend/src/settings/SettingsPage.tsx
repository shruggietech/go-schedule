import {
  Button,
  CardSection,
  DescriptionList,
  Notice,
  PathDisplay,
  StatePanel,
  StatusLabel,
} from "../components";
import type { Appearance } from "../connection/model";
import type { SettingsWorkspace } from "./model";

const transitionCopy: Record<string, { title: string; detail: string }> = {
  migrated: {
    title: "Appearance migrated",
    detail:
      "Your previous light, dark, or system choice was carried into the new desktop.",
  },
  not_found: {
    title: "System appearance selected",
    detail:
      "No earlier desktop appearance was found, so go-schedule follows your operating system.",
  },
  invalid: {
    title: "Earlier preference could not be used",
    detail:
      "The previous appearance value was invalid. go-schedule uses system appearance and remains usable.",
  },
  unreadable: {
    title: "Earlier preference could not be read",
    detail:
      "The previous preference file was unavailable. go-schedule uses system appearance and remains usable.",
  },
  not_required: {
    title: "Desktop defaults restored",
    detail:
      "The current desktop preference file is authoritative and follows your operating system.",
  },
};
type CopyResult = "copied" | "failed";

export function SettingsPage({
  workspace,
  message,
  pending = false,
  pendingActions = new Set(),
  copyResults = {},
  onAppearance,
  onRestore,
  onCopy,
  onOpen,
  onConnections,
}: {
  workspace?: SettingsWorkspace;
  message: string;
  pending?: boolean;
  pendingActions?: ReadonlySet<string>;
  copyResults?: Record<string, CopyResult>;
  onAppearance(value: Appearance): void;
  onRestore(): void;
  onCopy(id: string): void;
  onOpen(key: string): void;
  onConnections(): void;
}) {
  const isPending = (key: string) =>
    pendingActions.size ? pendingActions.has(key) : pending;
  if (!workspace)
    return (
      <>
        <PageHeader />
        <StatePanel
          title="Desktop settings are unavailable"
          detail={
            message ||
            "Check access to the user configuration directory, then try again."
          }
          busy={pending}
          action="Restore desktop defaults"
          onAction={onRestore}
        />
      </>
    );
  const transition =
    transitionCopy[workspace.preferences.transition.status] ??
    transitionCopy.not_found;
  return (
    <>
      <PageHeader />
      <div className="settings-grid">
        <CardSection eyebrow="Desktop preference" title="Appearance">
          <p>Choose a durable color mode for this desktop.</p>
          <label className="settings-control">
            Appearance
            <select
              value={workspace.preferences.appearance}
              disabled={isPending("appearance")}
              onChange={(event) =>
                onAppearance(event.target.value as Appearance)
              }
            >
              <option value="system">System</option>
              <option value="light">Light</option>
              <option value="dark">Dark</option>
            </select>
          </label>
          <Button
            variant="secondary"
            pending={isPending("restore")}
            onClick={onRestore}
          >
            Restore desktop defaults
          </Button>
          <PathDisplay
            label="Preference file path"
            value={workspace.preferencePath}
          />
        </CardSection>
        <CardSection eyebrow="Upgrade" title={transition.title}>
          <p>{transition.detail}</p>
          <p>
            Legacy font selection and scroll sensitivity were retired because
            browser typography and native scrolling now provide those behaviors.
          </p>
        </CardSection>
      </div>
      {!workspace.daemonAvailable && (
        <Notice title="Daemon storage details are unavailable" tone="warning">
          Local preferences and product information are still available.
          Daemon-owned paths are not guessed.{" "}
          <Button variant="secondary" onClick={onConnections}>
            Open Connections
          </Button>
        </Notice>
      )}
      <CardSection eyebrow="Local files" title="Application storage">
        <p>
          Paths are read-only. Ownership and removal behavior remain explicit,
          including locations configured outside application-owned data.
        </p>
        <div className="storage-list">
          {workspace.storage.map((record) => {
            const copyPending = isPending(`copy:${record.id}`);
            const copyResult = copyResults[record.id];
            return (
              <article
                className={`storage-record storage-${record.existence}`}
                key={record.id}
              >
                <div className="storage-heading">
                  <h3>{record.label}</h3>
                  <StatusLabel
                    tone={
                      record.existence === "present"
                        ? "positive"
                        : record.existence === "absent"
                          ? "neutral"
                          : "warning"
                    }
                  >
                    {record.existence === "present"
                      ? "Present"
                      : record.existence === "absent"
                        ? "Not present"
                        : "Unavailable"}
                  </StatusLabel>
                </div>
                <PathDisplay
                  label={`${record.label} path`}
                  value={record.path}
                />
                <DescriptionList>
                  <dt>Owner</dt>
                  <dd>{record.owner}</dd>
                  <dt>Scope</dt>
                  <dd>{record.scope}</dd>
                  <dt>Normal removal</dt>
                  <dd>{record.normalRemoval}</dd>
                  <dt>Explicit wipe</dt>
                  <dd>{record.explicitWipe}</dd>
                </DescriptionList>
                {record.copyable && (
                  <Button
                    aria-label={`Copy ${record.label} path`}
                    variant="secondary"
                    pending={copyPending}
                    onClick={() => onCopy(record.id)}
                  >
                    {copyPending
                      ? "Copying"
                      : copyResult === "copied"
                        ? "Copied"
                        : copyResult === "failed"
                          ? "Copy failed"
                          : "Copy path"}
                  </Button>
                )}
              </article>
            );
          })}
        </div>
      </CardSection>
      <CardSection
        eyebrow="Product information"
        title={workspace.product.name}
        className="about-section"
      >
        <img src="/go-schedule-mark.svg" alt="" />
        <div>
          <p>Version {workspace.product.version}</p>
          <p>Built and maintained by {workspace.product.publisher}.</p>
          <div className="product-links">
            {workspace.product.links.map((link) => (
              <Button
                key={link.key}
                variant="secondary"
                pending={isPending(`open:${link.key}`)}
                onClick={() => onOpen(link.key)}
              >
                {link.label}
              </Button>
            ))}
          </div>
        </div>
      </CardSection>
    </>
  );
}

function PageHeader() {
  return (
    <header className="page-header">
      <div>
        <p className="eyebrow">This computer</p>
        <h1>Settings</h1>
        <p>Desktop preferences, storage ownership, and product information.</p>
      </div>
    </header>
  );
}
