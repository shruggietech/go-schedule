package clievaluate

import (
	"context"
	"path/filepath"
	"testing"
	"time"
)

func TestExecutablePolicyScenarios(t *testing.T) {
	for _, name := range []string{"dst-gap-fold", "month-boundaries", "bundle-reference-integrity", "expired-pairing", "notification-retry", "notification-transport", "catchup-startup", "one-off-completion", "overlap-policies", "popup-filtering"} {
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(DefaultBudget().CaseTimeoutMS)*time.Millisecond)
			defer cancel()
			var result execution
			switch name {
			case "dst-gap-fold", "month-boundaries":
				result = calendarPolicy(name)
			case "bundle-reference-integrity":
				result = bundlePolicy()
			case "expired-pairing":
				result = expiredPairingPolicy()
			case "notification-retry":
				result = notificationRetry(ctx)
			case "notification-transport":
				result = notificationTransport(ctx)
			case "one-off-completion":
				result = oneOffPolicy(ctx)
			case "catchup-startup", "overlap-policies":
				result = enginePolicy(ctx, name)
			case "popup-filtering":
				result = popupPolicy(filepath.Join(t.TempDir(), "preferences.json"))
			}
			for _, passed := range result.Checks {
				if !passed {
					t.Fatalf("scenario did not satisfy its contract: %s", result.State)
				}
			}
		})
	}
}
