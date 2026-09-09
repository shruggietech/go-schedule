package cli

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/shruggietech/go-schedule/internal/clientprofile"
	"github.com/shruggietech/go-schedule/internal/clientsecret"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/remoteenroll"
)

type safeProfile struct {
	ID                     string    `json:"id"`
	Label                  string    `json:"label"`
	Endpoint               string    `json:"endpoint"`
	DaemonID               string    `json:"daemon_id"`
	CredentialID           string    `json:"credential_id"`
	CertificateFingerprint string    `json:"certificate_fingerprint"`
	ClientKind             string    `json:"client_kind"`
	Capability             string    `json:"capability"`
	DaemonDisplayName      string    `json:"daemon_display_name"`
	Platform               string    `json:"platform,omitempty"`
	Architecture           string    `json:"architecture,omitempty"`
	ProductVersion         string    `json:"product_version,omitempty"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
	LastSuccessfulAt       time.Time `json:"last_successful_at,omitempty"`
}

func safeProfileOf(p clientprofile.Profile) safeProfile {
	return safeProfile{ID: p.ID, Label: p.Label, Endpoint: p.Endpoint, DaemonID: p.DaemonID, CredentialID: p.CredentialID, CertificateFingerprint: p.CertificateFingerprint, ClientKind: p.ClientKind, Capability: p.Capability, DaemonDisplayName: p.DaemonDisplayName, Platform: p.Platform, Architecture: p.Architecture, ProductVersion: p.ProductVersion, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt, LastSuccessfulAt: p.LastSuccessfulAt}
}

func newProfileCmd() *cobra.Command {
	command := &cobra.Command{Use: "profile", Short: "Pair and administer remote connection profiles"}
	command.AddCommand(newProfilePairCmd(), newProfileListCmd(), newProfileShowCmd(), newProfileRenameCmd(), newProfileRemoveCmd())
	return command
}

func newProfileListCmd() *cobra.Command {
	return &cobra.Command{Use: "list", Short: "List saved remote profiles", RunE: func(cmd *cobra.Command, _ []string) error {
		collection, err := clientprofile.NewStore("").Load()
		if err != nil {
			return err
		}
		items := make([]safeProfile, 0, len(collection.Profiles))
		for _, p := range collection.Profiles {
			items = append(items, safeProfileOf(p))
		}
		if jsonOut {
			return printJSONTo(cmd.OutOrStdout(), map[string]any{"active_desktop_profile_id": collection.ActiveDesktopProfileID, "profiles": items})
		}
		for _, p := range items {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\t%s\t%s\t%s\n", p.ID, p.Label, p.Endpoint, shortID(p.DaemonID), p.Capability)
		}
		return nil
	}}
}

func newProfileShowCmd() *cobra.Command {
	return &cobra.Command{Use: "show PROFILE", Short: "Show one saved remote profile", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		collection, err := clientprofile.NewStore("").Load()
		if err != nil {
			return err
		}
		profile, err := collection.Find(args[0])
		if err != nil {
			return fmtUsage("profile does not exist or its label is ambiguous")
		}
		safe := safeProfileOf(profile)
		if jsonOut {
			return printJSONTo(cmd.OutOrStdout(), safe)
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "Profile: %s\nID: %s\nEndpoint: %s\nDaemon ID: %s\nCredential ID: %s\nCertificate SHA-256: %s\nCapability: %s\n", safe.Label, safe.ID, safe.Endpoint, safe.DaemonID, safe.CredentialID, safe.CertificateFingerprint, safe.Capability)
		return err
	}}
}

func newProfileRenameCmd() *cobra.Command {
	return &cobra.Command{Use: "rename PROFILE LABEL", Short: "Rename one saved profile locally", Args: cobra.ExactArgs(2), RunE: func(cmd *cobra.Command, args []string) error {
		profile, err := clientprofile.NewStore("").Rename(args[0], args[1])
		if err != nil {
			return fmtUsage("profile could not be renamed")
		}
		safe := safeProfileOf(profile)
		if jsonOut {
			return printJSONTo(cmd.OutOrStdout(), safe)
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "Profile %s renamed to %s.\n", safe.ID, safe.Label)
		return err
	}}
}

func newProfileRemoveCmd() *cobra.Command {
	var confirm string
	command := &cobra.Command{Use: "remove PROFILE", Short: "Delete one profile and its native credential", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		store := clientprofile.NewStore("")
		collection, err := store.Load()
		if err != nil {
			return err
		}
		profile, err := collection.Find(args[0])
		if err != nil {
			return fmtUsage("profile does not exist or its label is ambiguous")
		}
		if confirm != profile.ID {
			return fmtUsage("--confirm must exactly match the profile ID")
		}
		if err := clientsecret.New().Delete(profile.DaemonID, profile.CredentialID); err != nil {
			return fmt.Errorf("delete native credential: %w", err)
		}
		if _, err := store.Remove(profile.ID); err != nil {
			return fmt.Errorf("remove profile metadata after credential deletion: %w", err)
		}
		if jsonOut {
			return printJSONTo(cmd.OutOrStdout(), map[string]string{"id": profile.ID, "outcome": "removed"})
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "Removed profile %s (%s).\n", profile.Label, profile.ID)
		return err
	}}
	command.Flags().StringVar(&confirm, "confirm", "", "exact profile ID confirmation")
	return command
}

func newProfilePairCmd() *cobra.Command {
	var address, daemonID, pairingID, certificateFile, clientName, capability string
	command := &cobra.Command{Use: "pair LABEL", Short: "Pair and save one remote daemon", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		certificate, err := os.ReadFile(certificateFile)
		if err != nil {
			return fmt.Errorf("read trusted certificate: %w", err)
		}
		phraseBytes, err := io.ReadAll(io.LimitReader(cmd.InOrStdin(), 4097))
		if err != nil {
			return fmt.Errorf("read pairing phrase: %w", err)
		}
		phrase := strings.TrimSpace(string(phraseBytes))
		if phrase == "" || len(phraseBytes) > 4096 {
			return fmtUsage("provide the one-time phrase on standard input")
		}
		grant := domain.Capability(capability)
		if !grant.Valid() {
			return fmtUsage("capability must be observe, operate, manage, or enroll")
		}
		if err := clientsecret.New().Probe(); err != nil {
			return fmt.Errorf("native credential storage is unavailable: %w", err)
		}
		ctx, cancel := reqCtx()
		defer cancel()
		issued, manifest, err := remoteenroll.Exchange(ctx, remoteenroll.Draft{Address: address, DaemonID: daemonID, PairingID: pairingID, Phrase: phrase, CertificatePEM: string(certificate), DisplayName: clientName, Kind: domain.ActorKindCLI, Capability: grant})
		if err != nil {
			return fmt.Errorf("pair remote profile: %w", err)
		}
		id, err := clientprofile.NewID()
		if err != nil {
			return err
		}
		now := time.Now().UTC()
		profile, err := clientprofile.Normalize(clientprofile.Profile{ID: id, Label: args[0], Endpoint: address, DaemonID: issued.DaemonID, CredentialID: issued.ID, CertificatePEM: string(certificate), ClientKind: string(issued.Actor.Kind), Capability: string(issued.Actor.Capability), DaemonDisplayName: manifest.DisplayName, Platform: manifest.Platform.OS, Architecture: manifest.Platform.Architecture, ProductVersion: manifest.ProductVersion, CreatedAt: now, UpdatedAt: now, LastSuccessfulAt: now})
		if err != nil {
			return err
		}
		secrets := clientsecret.New()
		if err := secrets.Save(profile.DaemonID, profile.CredentialID, issued.Token); err != nil {
			return fmt.Errorf("save native credential: %w", err)
		}
		if _, err := clientprofile.NewStore("").Add(profile); err != nil {
			_ = secrets.Delete(profile.DaemonID, profile.CredentialID)
			return fmt.Errorf("save connection profile: %w", err)
		}
		issued.Token = ""
		safe := safeProfileOf(profile)
		if jsonOut {
			return printJSONTo(cmd.OutOrStdout(), safe)
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "Paired profile %s (%s) with %s.\n", safe.Label, safe.ID, safe.Endpoint)
		return err
	}}
	command.Flags().StringVar(&address, "address", "", "remote HTTPS origin")
	command.Flags().StringVar(&daemonID, "expected-daemon-id", "", "expected daemon installation ID")
	command.Flags().StringVar(&pairingID, "pairing-id", "", "one-time pairing session ID")
	command.Flags().StringVar(&certificateFile, "trusted-certificate", "", "trusted PEM certificate file")
	command.Flags().StringVar(&clientName, "client-name", "", "administrator-visible client name")
	command.Flags().StringVar(&capability, "capability", "observe", "requested capability")
	_ = command.MarkFlagRequired("address")
	_ = command.MarkFlagRequired("expected-daemon-id")
	_ = command.MarkFlagRequired("pairing-id")
	_ = command.MarkFlagRequired("trusted-certificate")
	_ = command.MarkFlagRequired("client-name")
	return command
}

func shortID(value string) string {
	if len(value) > 8 {
		return value[:8]
	}
	return value
}
