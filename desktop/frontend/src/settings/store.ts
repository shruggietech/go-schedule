import { useCallback, useEffect, useRef, useState } from "react";
import type { Appearance } from "../connection/model";
import { settingsBridge } from "./bridge";
import type {
  SettingsBridge,
  SettingsResult,
  SettingsWorkspace,
} from "./model";

export type SettingsStatus = Pick<
  SettingsResult,
  "action" | "outcome" | "message"
> & { id: number };

export function useSettings(
  bridge: SettingsBridge = settingsBridge,
  refreshToken: string | number = 0,
) {
  const [workspace, setWorkspace] = useState<SettingsWorkspace>();
  const [status, setStatus] = useState<SettingsStatus>();
  const [pendingActions, setPendingActions] = useState<ReadonlySet<string>>(
    new Set(),
  );
  const [copyResults, setCopyResults] = useState<
    Record<string, "copied" | "failed">
  >({});
  const [loading, setLoading] = useState(true);
  const pendingRef = useRef(new Set<string>());
  const statusSequence = useRef(0);
  const publish = useCallback(
    (result: Pick<SettingsResult, "action" | "outcome" | "message">) => {
      statusSequence.current += 1;
      setStatus({ id: statusSequence.current, ...result });
    },
    [],
  );

  useEffect(() => {
    let active = true;
    void bridge.workspace().then(
      (result) => {
        if (!active) return;
        if (result.outcome === "accepted" && result.workspace)
          setWorkspace(result.workspace);
        setLoading(false);
      },
      () => {
        if (!active) return;
        publish({
          action: "load_settings",
          outcome: "unavailable",
          message:
            "Desktop settings are unavailable. Check access to the user configuration directory, then try again.",
        });
        setLoading(false);
      },
    );
    return () => {
      active = false;
    };
  }, [bridge, publish, refreshToken]);

  const run = useCallback(
    async (
      key: string,
      operation: () => ReturnType<SettingsBridge["workspace"]>,
    ) => {
      if (pendingRef.current.has(key)) return undefined;
      pendingRef.current.add(key);
      setPendingActions(new Set(pendingRef.current));
      try {
        const result = await operation();
        publish(result);
        if (result.outcome === "accepted" && result.workspace)
          setWorkspace(result.workspace);
        return result;
      } catch {
        publish({
          action: "desktop_action",
          outcome: "unavailable",
          message: "The desktop action could not be completed. Try again.",
        });
        return undefined;
      } finally {
        pendingRef.current.delete(key);
        setPendingActions(new Set(pendingRef.current));
      }
    },
    [publish],
  );

  const saveAppearance = useCallback(
    (value: Appearance) =>
      run("appearance", () => bridge.saveAppearance(value)),
    [bridge, run],
  );
  const restore = useCallback(
    () => run("restore", () => bridge.restore()),
    [bridge, run],
  );
  const copyStoragePath = useCallback(
    async (id: string) => {
      setCopyResults((current) => {
        const next = { ...current };
        delete next[id];
        return next;
      });
      const result = await run(`copy:${id}`, () => bridge.copyStoragePath(id));
      if (result)
        setCopyResults((current) => ({
          ...current,
          [id]: result.outcome === "accepted" ? "copied" : "failed",
        }));
      return result;
    },
    [bridge, run],
  );
  const openProductLink = useCallback(
    (key: string) => run(`open:${key}`, () => bridge.openProductLink(key)),
    [bridge, run],
  );

  return {
    workspace,
    message: status?.message ?? "",
    status,
    loading,
    pending: pendingActions.size > 0,
    pendingActions,
    copyResults,
    saveAppearance,
    restore,
    copyStoragePath,
    openProductLink,
  };
}
