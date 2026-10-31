package cli

import (
	"fmt"
	"os/user"
	"strings"

	"github.com/spf13/cobra"

	"github.com/shruggietech/go-schedule/internal/userprefs"
)

type preferenceOutput struct {
	User        string                `json:"user"`
	Path        string                `json:"path"`
	Preferences userprefs.Preferences `json:"preferences"`
}

func currentPreferenceUser() (string, error) {
	current, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("identify the current preference owner: %w", err)
	}
	return current.Username, nil
}

func preferenceArgs(count int) cobra.PositionalArgs {
	return func(_ *cobra.Command, args []string) error {
		if len(args) != count {
			return fmtUsage(fmt.Sprintf("expected %d positional arguments, received %d", count, len(args)))
		}
		return nil
	}
}

func preferencesFor(cmd *cobra.Command) *userprefs.Store {
	deps := stateFor(cmd).deps
	return userprefs.NewStore(userprefs.Options{Path: deps.PreferencePath})
}

func writePreferencesResult(cmd *cobra.Command, owner string, store *userprefs.Store, prefs userprefs.Preferences) error {
	result := preferenceOutput{User: owner, Path: store.Path(), Preferences: prefs}
	if jsonOutput(cmd) {
		return printJSONTo(cmd.OutOrStdout(), result)
	}
	_, err := fmt.Fprintf(cmd.OutOrStdout(), "User: %s\nPreference path: %s\nAppearance: %s\nDesktop popups enabled: %t\nPopup conditions: %s\nPopup severities: %s\nPopup source daemon IDs: %s\n", owner, store.Path(), prefs.Appearance, prefs.Popups.Enabled, strings.Join(prefs.Popups.Conditions, ", "), strings.Join(prefs.Popups.Severities, ", "), popupSources(prefs.Popups.DaemonIDs))
	return err
}

func popupSources(values []string) string {
	if len(values) == 0 {
		return "all registered schedulers"
	}
	return strings.Join(values, ", ")
}

func newPreferencesCmd() *cobra.Command {
	command := &cobra.Command{Use: "preferences", Short: "Read and change this user's desktop preferences without opening the desktop"}
	MarkLocal(command)
	show := &cobra.Command{Use: "show", Short: "Show all durable desktop preferences and their user-owned path", Args: preferenceArgs(0), RunE: func(cmd *cobra.Command, _ []string) error {
		owner, err := currentPreferenceUser()
		if err != nil {
			return err
		}
		store := preferencesFor(cmd)
		prefs, err := store.Load()
		if err != nil {
			return err
		}
		return writePreferencesResult(cmd, owner, store, prefs)
	}}
	var appearance string
	var enabled bool
	var conditions, severities, daemons []string
	set := &cobra.Command{Use: "set", Short: "Change only explicitly supplied appearance or desktop popup choices", Args: preferenceArgs(0), RunE: func(cmd *cobra.Command, _ []string) error {
		var patch userprefs.Patch
		mode := userprefs.Appearance(appearance)
		if cmd.Flags().Changed("appearance") {
			patch.Appearance = &mode
		}
		if cmd.Flags().Changed("popup-enabled") {
			patch.PopupEnabled = &enabled
		}
		if cmd.Flags().Changed("popup-conditions") {
			patch.PopupConditions = &conditions
		}
		if cmd.Flags().Changed("popup-severities") {
			patch.PopupSeverities = &severities
		}
		if cmd.Flags().Changed("popup-daemons") {
			patch.PopupDaemonIDs = &daemons
		}
		if err := userprefs.ValidatePatch(patch); err != nil {
			return fmtUsage(err.Error())
		}
		owner, err := currentPreferenceUser()
		if err != nil {
			return err
		}
		store := preferencesFor(cmd)
		prefs, err := store.Patch(patch)
		if err != nil {
			return err
		}
		return writePreferencesResult(cmd, owner, store, prefs)
	}}
	set.Flags().StringVar(&appearance, "appearance", "", "system, light, or dark")
	set.Flags().BoolVar(&enabled, "popup-enabled", false, "enable desktop popups while the desktop is open; use =false to disable")
	set.Flags().StringSliceVar(&conditions, "popup-conditions", nil, "replace popup conditions: failure, success, alert (comma separated)")
	set.Flags().StringSliceVar(&severities, "popup-severities", nil, "replace popup severities: info, warning, error (comma separated)")
	set.Flags().StringSliceVar(&daemons, "popup-daemons", nil, "replace source daemon IDs (comma separated); = clears the filter to include all")
	var scope string
	reset := &cobra.Command{Use: "reset --scope appearance|popups|all", Short: "Restore explicit preference scope; all also repairs an unreadable preference document", Args: preferenceArgs(0), RunE: func(cmd *cobra.Command, _ []string) error {
		if scope != "appearance" && scope != "popups" && scope != "all" {
			return fmtUsage("--scope must be appearance, popups, or all")
		}
		owner, err := currentPreferenceUser()
		if err != nil {
			return err
		}
		store := preferencesFor(cmd)
		prefs, err := store.Reset(scope)
		if err != nil {
			return err
		}
		return writePreferencesResult(cmd, owner, store, prefs)
	}}
	reset.Flags().StringVar(&scope, "scope", "", "appearance, popups, or all; saved connections and daemon notifications are separate")
	command.AddCommand(show, set, reset)
	return command
}
