package settings

import (
	"fmt"
	"strings"
)

func validatePopupPreferences(value PopupPreferences) error {
	if len(value.Conditions) == 0 || len(value.Severities) == 0 {
		return fmt.Errorf("choose at least one popup condition and severity")
	}
	if len(value.Conditions) > 3 || len(value.Severities) > 3 || len(value.DaemonIDs) > 64 {
		return fmt.Errorf("too many popup filters")
	}
	if !uniqueAllowed(value.Conditions, map[string]bool{"failure": true, "success": true, "alert": true}) || !uniqueAllowed(value.Severities, map[string]bool{"info": true, "warning": true, "error": true}) {
		return fmt.Errorf("unsupported popup condition or severity")
	}
	seen := make(map[string]bool, len(value.DaemonIDs))
	for _, id := range value.DaemonIDs {
		if strings.TrimSpace(id) != id || id == "" || len(id) > 128 || seen[id] {
			return fmt.Errorf("invalid or duplicate popup daemon")
		}
		seen[id] = true
	}
	return nil
}

func uniqueAllowed(values []string, allowed map[string]bool) bool {
	seen := make(map[string]bool, len(values))
	for _, value := range values {
		if !allowed[value] || seen[value] {
			return false
		}
		seen[value] = true
	}
	return true
}
