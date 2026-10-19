package settings

import "github.com/shruggietech/go-schedule/internal/userprefs"

func validatePopupPreferences(value PopupPreferences) error { return userprefs.ValidatePopups(value) }
