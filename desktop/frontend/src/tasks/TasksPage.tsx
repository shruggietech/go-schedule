import { useMemo, useState } from "react";
import { Button, Notice, StatePanel } from "../components";
import { GroupsPanel } from "./GroupsPanel";
import { TaskActions } from "./TaskActions";
import { TaskEditor } from "./TaskEditor";
import {
  blankTask,
  type OperationResult,
  type TaskBridge,
  type TaskDraft,
} from "./model";
import { useTaskWorkspace } from "./store";
export function TasksPage({
  bridge,
  platform,
  available = true,
  workspaceAvailable = true,
  manageAvailable = available,
  editAvailable = manageAvailable,
  groupAvailable = manageAvailable,
  previewAvailable = manageAvailable,
  targetName = "This computer",
  refreshToken = 0,
  onActivity,
}: {
  bridge: TaskBridge;
  platform: string;
  available?: boolean;
  workspaceAvailable?: boolean;
  manageAvailable?: boolean;
  editAvailable?: boolean;
  groupAvailable?: boolean;
  previewAvailable?: boolean;
  targetName?: string;
  refreshToken?: number;
  onActivity(): void;
}) {
  const { workspace, status, selected, setSelected, load, accept, setStatus } =
    useTaskWorkspace(bridge, workspaceAvailable, refreshToken);
  const [editor, setEditor] = useState<{
    draft: TaskDraft;
    invoker: HTMLElement | null;
  }>();
  const [query, setQuery] = useState("");
  const [state, setState] = useState("all");
  const tasks = useMemo(
    () =>
      workspace?.tasks.filter(
        (task) =>
          (state === "all" || task.effectiveState === state) &&
          `${task.name} ${task.groupPath} ${task.scheduleSummary}`
            .toLowerCase()
            .includes(query.toLowerCase()),
      ) ?? [],
    [workspace, query, state],
  );
  const task = workspace?.tasks.find((value) => value.id === selected);
  const reconciling = status?.outcome === "uncertain";
  const handle = (result: OperationResult) => {
    accept(result);
    if (result.action === "run_task" && result.outcome === "accepted")
      onActivity();
  };
  const edit = async (invoker: HTMLElement) => {
    if (!task) return;
    const result = await bridge.task(task.id);
    setStatus(result);
    if (result.task)
      setEditor({
        draft: {
          ...result.task,
          isNew: false,
          originalUpdatedAt: result.task.updatedAt,
          overwriteStale: false,
        },
        invoker,
      });
  };
  if (!workspace && status?.outcome === "unavailable")
    return (
      <>
        <header className="page-header">
          <div>
            <p className="eyebrow">{targetName}</p>
            <h1>Tasks</h1>
            <p>Create, schedule, organize, and run local work.</p>
          </div>
        </header>
        <StatePanel
          title="Tasks are unavailable"
          detail={status.message}
          action="Try again"
          onAction={() => void load()}
        />
      </>
    );
  if (!workspace)
    return (
      <>
        <header className="page-header">
          <div>
            <p className="eyebrow">{targetName}</p>
            <h1>Tasks</h1>
          </div>
        </header>
        <StatePanel
          title="Loading tasks and groups"
          detail={`Reading the current scheduler state from ${targetName}.`}
          busy
        />
      </>
    );
  return (
    <>
      <header className="page-header">
        <div>
          <p className="eyebrow">{targetName}</p>
          <h1>Tasks</h1>
          <p>
            Create, schedule, organize, and run work on the selected scheduler.
          </p>
        </div>
        <Button
          disabled={!manageAvailable || reconciling}
          onClick={(event) =>
            setEditor({ draft: blankTask(platform), invoker: event.currentTarget })
          }
        >
          Create task
        </Button>
      </header>
      {status && status.action !== "load" && (
        <Notice
          title={status.outcome === "accepted" ? "Complete" : "Action needed"}
          tone={status.outcome === "accepted" ? "success" : "error"}
        >
          {status.message}
        </Notice>
      )}
      <div className="task-layout">
        <GroupsPanel
          groups={workspace.groups}
          bridge={bridge}
          available={groupAvailable && !reconciling}
          targetName={targetName}
          onResult={handle}
        />
        <section className="panel task-list" aria-labelledby="task-list-title">
          <div className="section-heading">
            <h2 id="task-list-title">Tasks</h2>
            <Button variant="quiet" onClick={() => void load()}>
              Refresh
            </Button>
          </div>
          <div className="filters">
            <label>
              Search{" "}
              <input
                type="search"
                value={query}
                onChange={(event) => setQuery(event.target.value)}
              />
            </label>
            <label>
              State{" "}
              <select
                value={state}
                onChange={(event) => setState(event.target.value)}
              >
                <option value="all">All</option>
                <option value="runnable">Runnable</option>
                <option value="manual_only">Manual only</option>
                <option value="task_disabled">Disabled</option>
                <option value="group_disabled">Group disabled</option>
                <option value="not_runnable">Needs setup</option>
              </select>
            </label>
          </div>
          {tasks.length === 0 ? (
            <StatePanel
              title="No matching tasks"
              detail="Create a task or change the current filters."
            />
          ) : (
            <div className="table-scroll" tabIndex={0}>
              <table>
                <caption>{tasks.length} tasks</caption>
                <thead>
                  <tr>
                    <th scope="col">Task</th>
                    <th scope="col">Group</th>
                    <th scope="col">State</th>
                    <th scope="col">Schedule</th>
                  </tr>
                </thead>
                <tbody>
                  {tasks.map((item) => (
                    <tr
                      key={item.id}
                      aria-selected={item.id === selected}
                      onClick={() => setSelected(item.id)}
                    >
                      <th scope="row">
                        <button
                          className="row-select"
                          onClick={() => setSelected(item.id)}
                        >
                          {item.name}
                        </button>
                      </th>
                      <td>{item.groupPath}</td>
                      <td title={item.effectiveReason}>
                        {item.effectiveState.replaceAll("_", " ")}
                      </td>
                      <td>{item.scheduleSummary}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
          {task && (
            <div className="task-detail">
              <h3>{task.name}</h3>
              <p>{task.effectiveReason}</p>
              <dl>
                <dt>Group</dt>
                <dd>{task.groupPath}</dd>
                <dt>Timezone</dt>
                <dd>{task.timezone}</dd>
                <dt>Policy</dt>
                <dd>{task.policySummary || "Default"}</dd>
              </dl>
              <div className="actions">
                <Button
                  variant="secondary"
                  disabled={!editAvailable || reconciling}
                  onClick={(event) => void edit(event.currentTarget)}
                >
                  Edit
                </Button>
                <TaskActions
                  task={task}
                  bridge={bridge}
                  available={available && !reconciling}
                  manageAvailable={manageAvailable && !reconciling}
                  targetName={targetName}
                  onResult={handle}
                />
              </div>
            </div>
          )}
        </section>
      </div>
      {editor && (
        <TaskEditor
          initial={editor.draft}
          invoker={editor.invoker}
          groups={workspace.groups}
          platform={platform}
          bridge={bridge}
          available={manageAvailable && !reconciling}
          previewAvailable={previewAvailable && !reconciling}
          targetName={targetName}
          onSaved={(result) => {
            handle(result);
            setEditor(undefined);
          }}
          onCancel={() => setEditor(undefined)}
        />
      )}
    </>
  );
}
