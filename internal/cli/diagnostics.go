package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/clientprofile"
)

type commandContextError struct {
	message string
	cause   error
}

func (e *commandContextError) Error() string { return e.message }
func (e *commandContextError) Unwrap() error { return e.cause }

// contextualCommandError adds actionable CLI context without replaying request
// payloads, credentials, trigger keys, or arbitrary server error messages.
func contextualCommandError(cmd *cobra.Command, args []string, err error) error {
	if err == nil || cmd == nil || isReported(err) {
		return err
	}
	var existing *commandContextError
	if errors.As(err, &existing) {
		return err
	}
	var uncertain *client.MutationUncertainError
	if errors.As(err, &uncertain) {
		return err
	}
	var status *client.StatusError
	if !errors.As(err, &status) {
		if cmd.Parent() != nil && cmd.Parent().Name() == "profile" && len(args) > 0 && (errors.Is(err, clientprofile.ErrNotFound) || errors.Is(err, clientprofile.ErrAmbiguous)) {
			return &commandContextError{message: fmt.Sprintf("%s: profile %s could not be selected; run gosched profile list and use an exact profile ID or an unambiguous label", cmd.CommandPath(), diagnosticIdentity(args[0])), cause: err}
		}
		return err
	}
	path := cmd.CommandPath()
	var message string
	switch status.Code {
	case server.CodeNotFound:
		kind, identity := requestedObject(cmd, args, status.Field)
		message = fmt.Sprintf("%s: requested %s was not found; check its identity and selected daemon, then run %s --help", path, kind, path)
		if identity != "" {
			message = fmt.Sprintf("%s: %s %s was not found; check its identity and selected daemon, then run %s --help", path, kind, diagnosticIdentity(identity), path)
		}
	case server.CodeValidation:
		field := diagnosticField(status.Field)
		if field == "watcher_paths" {
			message = watcherBindingDiagnostic(cmd, status)
		} else if field == "bundle" && len(args) > 0 && strings.Contains(cmd.Use, "<bundle-file>") {
			message = fmt.Sprintf("%s: invalid bundle file %s; run gosched bundle validate %s to list the document and compatibility issues, then correct the bundle and preview again", path, diagnosticIdentity(args[0]), diagnosticIdentity(args[0]))
		} else if label, guidance := scheduleValidationGuidance(field); label != "" {
			message = fmt.Sprintf("%s: invalid %s; %s", path, label, guidance)
		} else {
			label := fmt.Sprintf("field %q", field)
			flagName := strings.ReplaceAll(field, "_", "-")
			if cmd.Flag(flagName) != nil {
				label = "--" + flagName
			}
			message = fmt.Sprintf("%s: invalid %s; correct this input and run %s --help for accepted options", path, label, path)
		}
	default:
		return err
	}
	return &commandContextError{message: message, cause: err}
}

func watcherBindingDiagnostic(cmd *cobra.Command, status *client.StatusError) string {
	guidance := "use a watcher portable_id declared in the bundle exactly once as --watcher-path portable-id=absolute-path, or remove the binding"
	if flag := cmd.Flag("watcher-path"); flag != nil {
		if values, ok := flag.Value.(interface{ GetSlice() []string }); ok {
			for _, binding := range values.GetSlice() {
				id, _, found := strings.Cut(binding, "=")
				if !found {
					continue
				}
				// Match only server-owned diagnostic templates against identities
				// already supplied by this invocation. Bound paths are never read.
				if status.Message == fmt.Sprintf("watcher path binding %q is not a watcher in this bundle; use a watcher portable_id from the bundle or remove the binding", id) {
					return fmt.Sprintf("%s: --watcher-path identity %s is not a watcher in this bundle; %s", cmd.CommandPath(), diagnosticIdentity(id), guidance)
				}
				if status.Message == fmt.Sprintf("duplicate watcher path binding %q; supply each watcher portable_id exactly once", id) {
					return fmt.Sprintf("%s: duplicate --watcher-path identity %s; %s", cmd.CommandPath(), diagnosticIdentity(id), guidance)
				}
			}
		}
	}
	return fmt.Sprintf("%s: invalid watcher_paths binding; %s", cmd.CommandPath(), guidance)
}

func requestedObject(cmd *cobra.Command, args []string, field string) (string, string) {
	selectors := map[string]struct{ kind, flag string }{
		"group_id": {"group", "group"}, "target_task_id": {"task", "task"},
		"source_task_id": {"task", "source"}, "channel_id": {"notification channel", "channel"},
	}
	if selector, ok := selectors[field]; ok {
		if flag := cmd.Flag(selector.flag); flag != nil && flag.Changed {
			return selector.kind, flag.Value.String()
		}
	}
	nouns := map[string]string{"task": "task", "group": "group", "chain": "chain", "watcher": "watcher", "trigger": "trigger", "runs": "run", "actor": "actor", "pairing": "pairing", "credential": "credential", "channel": "notification channel", "set": "trigger set", "alerts": "alert"}
	kind := "object"
	for parent := cmd.Parent(); parent != nil; parent = parent.Parent() {
		if noun, ok := nouns[parent.Name()]; ok {
			kind = noun
			break
		}
	}
	// Only ID operands are echoed. In particular, trigger fire <key>, names,
	// schedule expressions, and document contents are never treated as IDs.
	use := strings.Fields(cmd.Use)
	if len(args) > 0 && len(use) > 1 && strings.Contains(strings.ToLower(use[1]), "id") {
		return kind, args[0]
	}
	return kind, ""
}

func diagnosticIdentity(value string) string {
	runes := []rune(value)
	if len(runes) > 160 {
		value = string(runes[:160]) + "..."
	}
	return fmt.Sprintf("%q", value)
}

func diagnosticField(field string) string {
	if field == "" || len(field) > 80 {
		return "request"
	}
	for _, value := range field {
		if (value < 'a' || value > 'z') && (value < 'A' || value > 'Z') && (value < '0' || value > '9') && !strings.ContainsRune("_.[]", value) {
			return "request"
		}
	}
	return field
}
