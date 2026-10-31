package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/shruggietech/go-schedule/v2/internal/api/client"
	"github.com/shruggietech/go-schedule/v2/internal/api/server"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

type accessClient interface {
	ListActors(context.Context) ([]domain.Actor, error)
	CreateActor(context.Context, server.ActorCreateRequest) (domain.Actor, error)
	UpdateActor(context.Context, string, server.ActorUpdateRequest) (domain.Actor, error)
	RevokeActor(context.Context, string) (domain.Actor, error)
	ListAudit(context.Context, domain.AuditQuery) ([]domain.AuditEvent, error)
	ExportAudit(context.Context, domain.AuditQuery) ([]byte, error)
}

func newPairingCmdWithClient(api *client.Client) *cobra.Command {
	command := &cobra.Command{Use: "pairing", Short: "Create and manage one-time remote pairing phrases"}
	var kind, capability, grantDuration string
	create := &cobra.Command{Use: "create <display-name>", Short: "Create a phrase displayed once", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		grantExpiresAt, err := grantExpiration(stateFor(cmd).deps.Now().UTC(), grantDuration)
		if err != nil {
			return err
		}
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		item, err := api.CreatePairing(ctx, server.PairingCreateRequest{DisplayName: args[0], Kind: domain.ActorKind(kind), Capability: domain.Capability(capability), ExpiresAt: grantExpiresAt})
		if err != nil {
			return err
		}
		if jsonOutput(cmd) {
			return printJSONTo(cmd.OutOrStdout(), item)
		}
		grantExpiry := "non-expiring"
		if item.GrantExpiresAt != nil {
			grantExpiry = item.GrantExpiresAt.Format(time.RFC3339)
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "Pairing ID: %s\nDaemon ID: %s\nEnrollment expires: %s\nGrant expires: %s\nPhrase (shown once): %s\n", item.ID, item.DaemonID, item.ExpiresAt.Format(time.RFC3339), grantExpiry, item.Phrase)
		return err
	}}
	create.Flags().StringVar(&kind, "kind", string(domain.ActorKindDesktop), "client kind: desktop, cli, json, or mcp")
	create.Flags().StringVar(&capability, "capability", string(domain.CapabilityObserve), "capability: observe, operate, manage, or enroll")
	create.Flags().StringVar(&grantDuration, "grant-duration", "non-expiring", "access duration: 1h, 24h, 7d, 30d, or non-expiring")
	list := &cobra.Command{Use: "list", Short: "List pairing sessions and remaining enrollment attempts", Long: "List one-time remote pairing sessions on the selected daemon, including each pairing ID, state, capability, client name and remaining attempts. This does not reveal pairing phrases or create credentials.", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		items, err := api.ListPairings(ctx)
		if err != nil {
			return err
		}
		if jsonOutput(cmd) {
			return printJSONTo(cmd.OutOrStdout(), items)
		}
		return writeHumanTable(cmd.OutOrStdout(), "PAIRING ID\tSTATE\tCAPABILITY\tCLIENT NAME\tATTEMPTS REMAINING", func(out io.Writer) error {
			for _, item := range items {
				if _, err := fmt.Fprintf(out, "%s\t%s\t%s\t%s\t%d\n", item.ID, item.State, item.Capability, item.DisplayName, item.AttemptsRemaining); err != nil {
					return err
				}
			}
			return nil
		})
	}}
	cancelCmd := &cobra.Command{Use: "cancel <pairing-id>", Short: "Cancel a pairing session so its phrase cannot enroll a client", Long: "Cancel the specified pairing ID on the selected daemon. Its one-time phrase can no longer enroll a client; already issued credentials are unaffected. Use pairing list to find the ID.", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		item, err := api.CancelPairing(ctx, args[0])
		if err != nil {
			return err
		}
		return printJSONTo(cmd.OutOrStdout(), item)
	}}
	command.AddCommand(create, list, cancelCmd)
	return command
}

func grantExpiration(now time.Time, value string) (*time.Time, error) {
	durations := map[string]time.Duration{"1h": time.Hour, "24h": 24 * time.Hour, "7d": 7 * 24 * time.Hour, "30d": 30 * 24 * time.Hour}
	if value == "non-expiring" {
		return nil, nil
	}
	duration, ok := durations[value]
	if !ok {
		return nil, fmtUsage("--grant-duration must be 1h, 24h, 7d, 30d, or non-expiring")
	}
	expires := now.Add(duration)
	return &expires, nil
}

func newCredentialCmdWithClient(api *client.Client) *cobra.Command {
	command := &cobra.Command{Use: "credential", Short: "Inspect, rotate, and revoke remote credentials"}
	command.AddCommand(&cobra.Command{Use: "list", Short: "List remote credential identities and their state", Long: "List remote credential IDs, associated actor IDs, expiry and state on the selected daemon. Tokens are never included. Use these IDs with credential rotate or credential revoke.", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		items, err := api.ListCredentials(ctx)
		if err != nil {
			return err
		}
		return printJSONTo(cmd.OutOrStdout(), items)
	}})
	command.AddCommand(&cobra.Command{Use: "rotate <credential-id>", Short: "Replace a credential token and display the replacement once", Long: "Rotate the specified remote credential ID on the selected daemon. The old token stops working and the replacement token is returned once with its actor and daemon identities. Store it before leaving this response; this does not update saved client profiles automatically.", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		item, err := api.RotateCredential(ctx, args[0])
		if err != nil {
			return err
		}
		return printJSONTo(cmd.OutOrStdout(), item)
	}})
	command.AddCommand(&cobra.Command{Use: "revoke <credential-id>", Short: "Revoke a remote credential so its token no longer works", Long: "Revoke the specified credential ID on the selected daemon. The credential's token can no longer authenticate; its associated actor record and other credentials remain. Use credential list to find the ID.", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		item, err := api.RevokeCredential(ctx, args[0])
		if err != nil {
			return err
		}
		return printJSONTo(cmd.OutOrStdout(), item)
	}})
	return command
}

func newActorCmdWithClient(api accessClient) *cobra.Command {
	command := &cobra.Command{Use: "actor", Short: "Manage client actors"}
	command.AddCommand(newActorListCmd(api), newActorCreateCmd(api), newActorUpdateCmd(api), newActorRevokeCmd(api))
	return command
}

func newActorListCmd(api accessClient) *cobra.Command {
	return &cobra.Command{Use: "list", Short: "List actors", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		actors, err := api.ListActors(ctx)
		if err != nil {
			return err
		}
		if jsonOutput(cmd) {
			return printJSONTo(cmd.OutOrStdout(), actors)
		}
		return writeHumanTable(cmd.OutOrStdout(), "ACTOR ID\tKIND\tCAPABILITY\tSTATE\tNAME", func(out io.Writer) error {
			for _, actor := range actors {
				if _, err := fmt.Fprintf(out, "%s\t%s\t%s\t%s\t%s\n", actor.ID, actor.Kind, actor.Capability, actor.State, actor.DisplayName); err != nil {
					return err
				}
			}
			return nil
		})
	}}
}

func newActorCreateCmd(api accessClient) *cobra.Command {
	var kind, capability, expiration string
	command := &cobra.Command{Use: "create <display-name>", Short: "Create an actor without a credential", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		request := server.ActorCreateRequest{Kind: domain.ActorKind(kind), DisplayName: args[0], Capability: domain.Capability(capability)}
		if expiration != "" {
			parsed, err := time.Parse(time.RFC3339, expiration)
			if err != nil {
				return fmtUsage("--expires must be an RFC3339 timestamp")
			}
			request.ExpiresAt = &parsed
		}
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		actor, err := api.CreateActor(ctx, request)
		if err != nil {
			return err
		}
		return printActor(cmd, actor, "Created")
	}}
	command.Flags().StringVar(&kind, "kind", string(domain.ActorKindCLI), "actor kind: desktop, cli, json, or mcp")
	command.Flags().StringVar(&capability, "capability", string(domain.CapabilityObserve), "capability: observe, operate, manage, or enroll")
	command.Flags().StringVar(&expiration, "expires", "", "optional RFC3339 expiration")
	return command
}

func newActorUpdateCmd(api accessClient) *cobra.Command {
	var name, capability, state, expiration string
	var clearExpiration bool
	command := &cobra.Command{Use: "update <actor-id>", Short: "Update an actor", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		request := server.ActorUpdateRequest{ClearExpiration: clearExpiration}
		if cmd.Flags().Changed("name") {
			request.DisplayName = &name
		}
		if cmd.Flags().Changed("capability") {
			value := domain.Capability(capability)
			request.Capability = &value
		}
		if cmd.Flags().Changed("state") {
			value := domain.ActorState(state)
			request.State = &value
		}
		if expiration != "" {
			parsed, err := time.Parse(time.RFC3339, expiration)
			if err != nil {
				return fmtUsage("--expires must be an RFC3339 timestamp")
			}
			request.ExpiresAt = &parsed
		}
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		actor, err := api.UpdateActor(ctx, args[0], request)
		if err != nil {
			return err
		}
		return printActor(cmd, actor, "Updated")
	}}
	command.Flags().StringVar(&name, "name", "", "new display name")
	command.Flags().StringVar(&capability, "capability", "", "new capability")
	command.Flags().StringVar(&state, "state", "", "new state: active, expired, or revoked")
	command.Flags().StringVar(&expiration, "expires", "", "new RFC3339 expiration")
	command.Flags().BoolVar(&clearExpiration, "clear-expiration", false, "remove expiration")
	return command
}

func newActorRevokeCmd(api accessClient) *cobra.Command {
	return &cobra.Command{Use: "revoke <actor-id>", Short: "Revoke an actor", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		actor, err := api.RevokeActor(ctx, args[0])
		if err != nil {
			return err
		}
		return printActor(cmd, actor, "Revoked")
	}}
}

func printActor(cmd *cobra.Command, actor domain.Actor, verb string) error {
	if jsonOutput(cmd) {
		return printJSONTo(cmd.OutOrStdout(), actor)
	}
	_, err := fmt.Fprintf(cmd.OutOrStdout(), "%s actor %s (%s).\n", verb, actor.DisplayName, actor.ID)
	return err
}

func newAuditCmdWithClient(api accessClient) *cobra.Command {
	command := &cobra.Command{Use: "audit", Short: "Inspect management audit history"}
	command.AddCommand(newAuditListCmd(api), newAuditExportCmd(api))
	return command
}

type auditFlags struct {
	actor, operation, result, since, until string
	limit                                  int
}

func (f *auditFlags) bind(command *cobra.Command) {
	command.Flags().StringVar(&f.actor, "actor", "", "filter by actor ID")
	command.Flags().StringVar(&f.operation, "operation", "", "filter by operation ID")
	command.Flags().StringVar(&f.result, "result", "", "filter by result")
	command.Flags().StringVar(&f.since, "since", "", "include events at or after RFC3339 time")
	command.Flags().StringVar(&f.until, "until", "", "include events at or before RFC3339 time")
	command.Flags().IntVar(&f.limit, "limit", 100, "maximum events (1 through 1000)")
}

func (f auditFlags) query() (domain.AuditQuery, error) {
	query := domain.AuditQuery{ActorID: f.actor, Operation: f.operation, Result: domain.AuditResult(f.result), Limit: f.limit}
	if f.since != "" {
		value, err := time.Parse(time.RFC3339, f.since)
		if err != nil {
			return query, fmtUsage("--since must be RFC3339")
		}
		query.Since = &value
	}
	if f.until != "" {
		value, err := time.Parse(time.RFC3339, f.until)
		if err != nil {
			return query, fmtUsage("--until must be RFC3339")
		}
		query.Until = &value
	}
	if f.limit < 1 || f.limit > 1000 {
		return query, fmtUsage("--limit must be between 1 and 1000 audit events")
	}
	if query.Result != "" && !query.Result.Valid() {
		return query, fmtUsage("--result must be succeeded, denied, failed, or uncertain; omit it to include all audit results")
	}
	return query, nil
}

func newAuditListCmd(api accessClient) *cobra.Command {
	var flags auditFlags
	command := &cobra.Command{Use: "list", Short: "List audit events", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		query, err := flags.query()
		if err != nil {
			return err
		}
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		events, err := api.ListAudit(ctx, query)
		if err != nil {
			return err
		}
		if jsonOutput(cmd) {
			return printJSONTo(cmd.OutOrStdout(), events)
		}
		return writeHumanTable(cmd.OutOrStdout(), "TIME\tRESULT\tACTOR ID\tOPERATION\tTARGET ID", func(out io.Writer) error {
			for _, event := range events {
				if _, err := fmt.Fprintf(out, "%s\t%s\t%s\t%s\t%s\n", event.OccurredAt.Format(time.RFC3339), event.Result, event.ActorID, event.Operation, event.TargetID); err != nil {
					return err
				}
			}
			return nil
		})
	}}
	flags.bind(command)
	return command
}

func newAuditExportCmd(api accessClient) *cobra.Command {
	var flags auditFlags
	var output string
	command := &cobra.Command{Use: "export", Short: "Export audit events as newline-delimited JSON", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		query, err := flags.query()
		if err != nil {
			return err
		}
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		data, err := api.ExportAudit(ctx, query)
		if err != nil {
			return err
		}
		if output == "" {
			_, err = cmd.OutOrStdout().Write(data)
			return err
		}
		if strings.ContainsRune(output, '\x00') {
			return fmtUsage("--output is invalid")
		}
		file, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
		if err != nil {
			return err
		}
		if err := file.Chmod(0o600); err != nil {
			_ = file.Close()
			return err
		}
		if _, err := file.Write(data); err != nil {
			_ = file.Close()
			return err
		}
		return file.Close()
	}}
	flags.bind(command)
	command.Flags().StringVarP(&output, "output", "o", "", "write export to a file")
	return command
}
