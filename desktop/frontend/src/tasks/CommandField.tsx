import { useRef } from "react";
import { Button, Field } from "../components";

const suggestions: Record<string, string> = {
  windows: "cmd.exe /d /c ver",
  macos: "/usr/bin/sw_vers",
  linux: "uname -a",
};
export function CommandField({
  platform,
  value,
  isNew,
  onChange,
  error,
}: {
  platform: string;
  value: string;
  isNew: boolean;
  onChange(value: string): void;
  error?: string;
}) {
  const inserted = useRef(false);
  const suggestion = isNew ? suggestions[platform] : undefined;
  const insert = () => {
    if (!value && suggestion) {
      inserted.current = true;
      onChange(suggestion);
    }
  };
  return (
    <div className="command-field">
      <Field
        label="Command line"
        help="Runs the program directly. Shell operators and substitutions are not interpreted."
        error={error}
      >
        <input
          name="command"
          value={value}
          onChange={(event) => {
            inserted.current = false;
            onChange(event.target.value);
          }}
          onKeyDown={(event) => {
            if (
              event.key === "Tab" &&
              !event.shiftKey &&
              !value &&
              suggestion &&
              !inserted.current
            ) {
              event.preventDefault();
              insert();
            }
          }}
        />
      </Field>
      {suggestion && (
        <p className="field-suggestion">
          Safe example for this computer: <code>{suggestion}</code>. It prints
          platform information and exits. Save the inactive task, choose Run
          now, then inspect captured output in Activity.{" "}
          <Button type="button" variant="secondary" onClick={insert}>
            Insert example
          </Button>
        </p>
      )}
    </div>
  );
}
export { suggestions };
