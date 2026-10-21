package clievaluate

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/shruggietech/go-schedule/internal/clifixture"
)

func (s *session) special(ctx context.Context, definition step) (bool, execution) {
	switch definition.Hook {
	case "invalid-daemon-config", "occupied-daemon-config":
		return true, s.daemonConfiguration(ctx, definition)
	case "mcp-stdio":
		return true, s.mcpProtocol(ctx, "manage", true)
	case "authority-audit":
		return true, s.authority(ctx)
	case "event-stream":
		return true, s.eventStream(ctx)
	case "notification-retry":
		return true, notificationRetry(ctx)
	case "notification-transport":
		return true, notificationTransport(ctx)
	case "watcher-recovery":
		return true, s.watcherRecovery(ctx)
	case "descendant-drain":
		return true, s.descendantDrain(ctx)
	case "canceled-secret", "partial-systems":
		return true, s.partialSystems(ctx, definition)
	case "mcp-authority":
		return true, s.mcpAuthority(ctx)
	case "dst-gap-fold", "month-boundaries":
		return true, calendarPolicy(definition.Hook)
	case "bundle-reference-integrity":
		return true, bundlePolicy()
	case "expired-pairing":
		return true, expiredPairingPolicy()
	case "popup-filtering":
		return true, popupPolicy(filepath.Join(s.w.Runtime, "client", "popup-fixture.json"))
	case "catchup-startup", "overlap-policies":
		return true, enginePolicy(ctx, definition.Hook)
	case "one-off-completion":
		return true, oneOffPolicy(ctx)
	case "socket-bound":
		parent := filepath.Join(s.w.Runtime, strings.Repeat("q", 120))
		w, err := clifixture.NewWorkspace(parent)
		_, statErr := os.Stat(parent)
		return true, policyResult(map[string]any{"workspace_absent": w == nil, "parent_absent": os.IsNotExist(statErr), "error": errorText(err)}, w == nil && err != nil && os.IsNotExist(statErr), nil)
	}
	if definition.Hook == "private-secret" {
		return true, s.secretBoundary(ctx, definition)
	}
	if strings.HasPrefix(definition.Spec.CommandPath, "@") {
		return true, execution{Blocked: "scenario_prerequisite", Detail: "scenario dispatcher is not yet bound: " + definition.Hook}
	}
	return false, execution{}
}
