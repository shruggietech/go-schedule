package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/clientprofile"
	"github.com/shruggietech/go-schedule/internal/clientsecret"
	"github.com/shruggietech/go-schedule/internal/config"
	"github.com/shruggietech/go-schedule/internal/ipc"
)

type remoteTargetFlags struct{ Profile, Endpoint, DaemonID, CredentialID, CertificateFile string }

func resolveTarget(ctx context.Context, flags remoteTargetFlags) (*client.Client, string, error) {
	explicit := []string{flags.Endpoint, flags.DaemonID, flags.CredentialID, flags.CertificateFile}
	count := 0
	for _, value := range explicit {
		if strings.TrimSpace(value) != "" {
			count++
		}
	}
	if flags.Profile != "" && count != 0 {
		return nil, "", fmtUsage("--profile cannot be combined with explicit remote target flags")
	}
	if flags.Profile == "" && count == 0 {
		cfg, _ := config.Load("")
		return client.New(ipc.Endpoint(cfg)), "", nil
	}
	var profile clientprofile.Profile
	if flags.Profile != "" {
		collection, err := clientprofile.NewStore("").Load()
		if err != nil {
			return nil, "", fmt.Errorf("load connection profiles: %w", err)
		}
		profile, err = collection.Find(flags.Profile)
		if err != nil {
			return nil, "", fmtUsage("remote profile does not exist or its label is ambiguous")
		}
	} else {
		if count != len(explicit) {
			return nil, "", fmtUsage("--endpoint, --daemon-id, --credential-id, and --certificate-file are required together")
		}
		certificate, err := os.ReadFile(flags.CertificateFile)
		if err != nil {
			return nil, "", fmt.Errorf("read trusted certificate: %w", err)
		}
		profile = clientprofile.Profile{Label: "explicit target", Endpoint: flags.Endpoint, DaemonID: flags.DaemonID, CredentialID: flags.CredentialID, CertificatePEM: string(certificate)}
	}
	token, err := clientsecret.New().Load(profile.DaemonID, profile.CredentialID)
	if err != nil {
		return nil, "", fmt.Errorf("load native remote credential: %w", err)
	}
	remote, err := client.NewRemote(profile.Endpoint, profile.CertificatePEM, token, profile.DaemonID)
	if err != nil {
		return nil, "", fmtUsage("remote target configuration is invalid")
	}
	if _, err := remote.VerifyIdentity(ctx); err != nil {
		return nil, "", fmt.Errorf("verify remote daemon identity: %w", err)
	}
	short := profile.DaemonID
	if len(short) > 8 {
		short = short[:8]
	}
	return remote, fmt.Sprintf("target: %s (%s, %s)", profile.Label, profile.Endpoint, short), nil
}
