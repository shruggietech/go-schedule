package clievaluate

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/api/client"
	"github.com/shruggietech/go-schedule/v2/internal/api/server"
	"github.com/shruggietech/go-schedule/v2/internal/cli"
	"github.com/shruggietech/go-schedule/v2/internal/clientsecret"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
	"github.com/shruggietech/go-schedule/v2/internal/remoteenroll"
)

func (s *session) authority(ctx context.Context) execution {
	observations := []any{}
	good := true
	for _, capability := range []domain.Capability{domain.CapabilityObserve, domain.CapabilityOperate, domain.CapabilityManage} {
		pair, err := s.api.CreatePairing(ctx, server.PairingCreateRequest{DisplayName: "QA authority " + string(capability), Kind: domain.ActorKindCLI, Capability: capability})
		if err != nil {
			return policyResult(observations, false, err)
		}
		s.redactor.Register(pair.Phrase)
		credential, _, err := remoteenroll.Exchange(ctx, remoteenroll.Draft{Address: s.vars["remote"], DaemonID: s.w.DaemonID, PairingID: pair.ID, Phrase: pair.Phrase, CertificatePEM: s.vars["certificate_pem"], DisplayName: pair.DisplayName, Kind: pair.Kind, Capability: pair.Capability})
		if err != nil {
			return policyResult(observations, false, err)
		}
		s.redactor.Register(credential.Token)
		remote, err := client.NewRemote(s.vars["remote"], s.vars["certificate_pem"], credential.Token, s.w.DaemonID)
		if err != nil {
			return policyResult(observations, false, err)
		}
		if _, err := remote.VerifyIdentity(ctx); err != nil {
			return policyResult(observations, false, err)
		}
		_, healthErr := remote.Health(ctx)
		good = good && healthErr == nil
		disabled := false
		created, createErr := remote.CreateTask(ctx, server.TaskCreateRequest{Name: "QA remote " + string(capability), Command: s.binaries["runner"], Timezone: "UTC", Enabled: &disabled})
		if capability == domain.CapabilityManage {
			good = good && createErr == nil && created.Task.ID != ""
			if created.Task.ID != "" {
				if err := remote.DeleteTask(ctx, created.Task.ID); err != nil {
					return policyResult(observations, false, err)
				}
			}
		} else {
			good = good && createErr != nil
		}
		runErr := remote.RunNow(ctx, s.vars["base_task"])
		good = good && ((capability == domain.CapabilityObserve && runErr != nil) || (capability != domain.CapabilityObserve && runErr == nil))
		audit, auditErr := s.api.ListAudit(ctx, domain.AuditQuery{ActorID: credential.Actor.ID, Limit: 20})
		good = good && auditErr == nil && len(audit) > 0
		observations = append(observations, map[string]any{"capability": capability, "actor_id": credential.Actor.ID, "create_error": errorText(createErr), "run_error": errorText(runErr), "audit": audit})
		if _, err := s.api.RevokeCredential(ctx, credential.ID); err != nil {
			return policyResult(observations, false, err)
		}
		_, revokedErr := remote.ListTasks(ctx, "", "")
		good = good && revokedErr != nil
	}
	return policyResult(observations, good, nil)
}

func (s *session) eventStream(ctx context.Context) execution {
	observations := []any{}
	for attempt := range 2 {
		streamCtx, cancel := context.WithCancel(ctx)
		seen := make(chan client.RemoteEvent, 8)
		done := make(chan error, 1)
		go func() {
			done <- s.api.StreamRemoteEvents(streamCtx, func(event client.RemoteEvent) {
				select {
				case seen <- event:
				default:
				}
			})
		}()
		ticker := time.NewTicker(40 * time.Millisecond)
		var received client.RemoteEvent
		var streamErr error
	loop:
		for {
			select {
			case received = <-seen:
				break loop
			case err := <-done:
				cancel()
				ticker.Stop()
				return policyResult(observations, false, fmt.Errorf("stream ended before observation: %w", err))
			case <-ticker.C:
				if err := s.w.Guard(ctx, s.api); err != nil {
					streamErr = err
					break loop
				}
				_, streamErr = s.api.UpdateGroup(ctx, s.vars["base_group"], server.GroupUpdateRequest{Name: fmt.Sprintf("QA stream %d", attempt)})
				if streamErr != nil {
					break loop
				}
			case <-ctx.Done():
				streamErr = ctx.Err()
				break loop
			}
		}
		ticker.Stop()
		cancel()
		stopErr := <-done
		observations = append(observations, map[string]any{"connection": attempt + 1, "event": received, "canceled_stream_error": errorText(stopErr)})
		if streamErr != nil || received.Kind == "" || stopErr == nil {
			return policyResult(observations, false, streamErr)
		}
	}
	return policyResult(observations, true, nil)
}

type interruptedSecrets struct {
	cli.SecretStore
	loads atomic.Int32
}

func (s *interruptedSecrets) LoadContext(context.Context, string, string) (string, error) {
	s.loads.Add(1)
	return "", context.Canceled
}

func (s *session) partialSystems(ctx context.Context, definition step) execution {
	args, err := s.arguments(definition.Spec.Argv)
	if err != nil {
		return execution{Blocked: "fixture_prerequisite", Detail: err.Error()}
	}
	var out, diagnostic bytes.Buffer
	deps := s.dependencies(strings.NewReader(""), &out, &diagnostic)
	secrets := &interruptedSecrets{SecretStore: clientsecret.NewWithBackend(s.recorder)}
	deps.Secrets = secrets
	code := cli.ExecuteWithDependencies(ctx, args, deps)
	observed := map[string]any{"secret_loads": secrets.loads.Load(), "diagnostic": diagnostic.String(), "partial_local_result": strings.Contains(out.String(), s.w.DaemonID)}
	return execution{Code: code, Args: args, Out: out.Bytes(), Err: diagnostic.Bytes(), State: mustObservation(observed), Checks: []bool{secrets.loads.Load() > 0, strings.Contains(out.String(), s.w.DaemonID) && code != 0}}
}
