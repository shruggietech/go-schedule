package cli

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
)

type accessClient interface {
	ListActors(context.Context) ([]domain.Actor, error)
	CreateActor(context.Context, server.ActorCreateRequest) (domain.Actor, error)
	UpdateActor(context.Context, string, server.ActorUpdateRequest) (domain.Actor, error)
	RevokeActor(context.Context, string) (domain.Actor, error)
	ListAudit(context.Context, domain.AuditQuery) ([]domain.AuditEvent, error)
	ExportAudit(context.Context, domain.AuditQuery) ([]byte, error)
}

func newActorCmd() *cobra.Command { return newActorCmdWithClient(newClient()) }

func newActorCmdWithClient(api accessClient) *cobra.Command {
	command := &cobra.Command{Use: "actor", Short: "Manage client actors"}
	command.AddCommand(newActorListCmd(api), newActorCreateCmd(api), newActorUpdateCmd(api), newActorRevokeCmd(api))
	return command
}

func newActorListCmd(api accessClient) *cobra.Command {
	return &cobra.Command{Use: "list", Short: "List actors", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		ctx, cancel := reqCtx()
		defer cancel()
		actors, err := api.ListActors(ctx)
		if err != nil {
			return err
		}
		if jsonOut {
			return printJSONTo(cmd.OutOrStdout(), actors)
		}
		for _, actor := range actors {
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\t%s\t%s\t%s\n", actor.ID, actor.Kind, actor.Capability, actor.State, actor.DisplayName); err != nil {
				return err
			}
		}
		return nil
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
		ctx, cancel := reqCtx()
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
		ctx, cancel := reqCtx()
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
		ctx, cancel := reqCtx()
		defer cancel()
		actor, err := api.RevokeActor(ctx, args[0])
		if err != nil {
			return err
		}
		return printActor(cmd, actor, "Revoked")
	}}
}

func printActor(cmd *cobra.Command, actor domain.Actor, verb string) error {
	if jsonOut {
		return printJSONTo(cmd.OutOrStdout(), actor)
	}
	_, err := fmt.Fprintf(cmd.OutOrStdout(), "%s actor %s (%s).\n", verb, actor.DisplayName, actor.ID)
	return err
}

func newAuditCmd() *cobra.Command { return newAuditCmdWithClient(newClient()) }

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
	if f.limit < 1 || f.limit > 1000 || (query.Result != "" && !query.Result.Valid()) {
		return query, fmtUsage("audit filters are invalid")
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
		ctx, cancel := reqCtx()
		defer cancel()
		events, err := api.ListAudit(ctx, query)
		if err != nil {
			return err
		}
		if jsonOut {
			return printJSONTo(cmd.OutOrStdout(), events)
		}
		for _, event := range events {
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\t%s\t%s\t%s\n", event.OccurredAt.Format(time.RFC3339), event.Result, event.ActorID, event.Operation, event.TargetID); err != nil {
				return err
			}
		}
		return nil
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
		ctx, cancel := reqCtx()
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
