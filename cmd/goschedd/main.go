// Command goschedd is the scheduler daemon. It hosts the engine, persistence,
// and executor and serves the local API over IPC. In this foundational form it
// loads config, opens the store, and serves health; the scheduling engine is
// wired in during User Story 1.
package main

import (
	"context"
	"errors"
	"flag"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/clock"
	"github.com/shruggietech/go-schedule/internal/config"
	"github.com/shruggietech/go-schedule/internal/engine"
	"github.com/shruggietech/go-schedule/internal/events"
	"github.com/shruggietech/go-schedule/internal/executor"
	"github.com/shruggietech/go-schedule/internal/ipc"
	"github.com/shruggietech/go-schedule/internal/lock"
	"github.com/shruggietech/go-schedule/internal/logbus"
	"github.com/shruggietech/go-schedule/internal/service"
	"github.com/shruggietech/go-schedule/internal/store"
)

func main() {
	configPath := flag.String("config", "", "path to config file (optional)")
	flag.Parse()

	if err := mainErr(*configPath); err != nil {
		os.Stderr.WriteString("goschedd: " + err.Error() + "\n")
		os.Exit(1)
	}
}

func mainErr(configPath string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}
	// Best-effort one-time move of a pre-rebrand (goscheduler) data directory onto
	// the new (goschedule) location. Runs before the data dir is created so a fresh
	// install is detected correctly. Non-fatal by design.
	config.MigrateLegacyPaths(cfg, config.NewLogger(cfg, os.Stdout))
	if err := os.MkdirAll(cfg.DataDir, 0o755); err != nil {
		return err
	}

	// Single-instance guard, acquired before the service machinery so a second
	// daemon fails fast (a second scheduler would double-execute every task).
	lk, err := lock.Acquire(filepath.Join(cfg.DataDir, "goschedd.lock"))
	if err != nil {
		return err
	}
	defer func() { _ = lk.Release() }()

	// Run under the service manager when launched as a service; otherwise this
	// runs in the foreground until interrupted.
	return service.Run(func(ctx context.Context) error {
		return runDaemon(ctx, cfg)
	})
}

func runDaemon(ctx context.Context, cfg config.Config) error {
	// Live-event broker doubles as the log publisher, so it is created first.
	broker := events.NewBroker()

	// Log pipeline: a teeing slog handler writes every record to a rotating
	// on-disk JSONL file, a bounded in-memory ring (served by GET /v1/logs), and
	// the broker (live stream). It also echoes to stdout for foreground/dev runs.
	ring := logbus.NewRing(cfg.LogRingSize)
	rw, err := logbus.NewRotatingWriter(cfg.LogPath(), int64(cfg.LogMaxSizeBytes), cfg.LogMaxFiles)
	if err != nil {
		return err
	}
	defer rw.Close()
	log := slog.New(logbus.NewHandler(cfg.SlogLevel(), ring, io.MultiWriter(os.Stdout, rw), broker))

	st, err := store.Open(cfg.DBPath())
	if err != nil {
		return err
	}
	defer st.Close()

	endpoint := ipc.Endpoint(cfg)
	ln, err := ipc.Listen(endpoint)
	if err != nil {
		return err
	}
	defer ln.Close()

	// Scheduling engine wired to the broker for run/alert streaming.
	eng := engine.New(st, clock.NewReal(), executor.New(cfg.OutputCapBytes), log, cfg.WorkerPoolSize)
	eng.SetOnRun(broker.PublishRun)
	eng.SetOnAlert(broker.PublishAlert)
	engErr := make(chan error, 1)
	go func() { engErr <- eng.Start(ctx) }()

	srv := &http.Server{
		Handler:           server.New(st, eng, broker, ring, log).Handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	serveErr := make(chan error, 1)
	go func() {
		log.Info("daemon listening", "endpoint", endpoint, "db", cfg.DBPath())
		if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serveErr <- err
			return
		}
		serveErr <- nil
	}()

	select {
	case <-ctx.Done():
		log.Info("shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := srv.Shutdown(shutdownCtx)
		<-engErr // wait for engine to drain in-flight runs
		return err
	case err := <-serveErr:
		return err
	case err := <-engErr:
		return err
	}
}
