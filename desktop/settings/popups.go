package settings

import "github.com/shruggietech/go-schedule/v2/internal/userprefs"

func validatePopupPreferences(value PopupPreferences) error { return userprefs.ValidatePopups(value) }
