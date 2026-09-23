//go:build linux

// gosched-indicator is the per-session Linux StatusNotifier companion.
// It never owns or starts the system daemon implicitly.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/gogpu/systray"

	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/config"
	"github.com/shruggietech/go-schedule/internal/desktopcontrol"
	"github.com/shruggietech/go-schedule/internal/ipc"
	"github.com/shruggietech/go-schedule/internal/service"
)

type indicator struct {
	tray         *systray.SystemTray
	monitor      desktopcontrol.Monitor
	status       *systray.MenuItem
	feedback     *systray.MenuItem
	start        *systray.MenuItem
	stop         *systray.MenuItem
	restart      *systray.MenuItem
	mu           sync.Mutex
	actions      sync.WaitGroup
	pending      bool
	closing      bool
	message      string
	messageUntil time.Time
	ctx          context.Context
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "gosched-indicator:", err)
		os.Exit(1)
	}
}

func run() error {
	instance, owner, err := desktopcontrol.ClaimTray()
	if err != nil {
		return fmt.Errorf("claim session indicator: %w", err)
	}
	if !owner {
		return nil
	}
	defer instance.Close() //nolint:errcheck // process-lifetime session lock
	bus, err := dbus.ConnectSessionBus()
	if err != nil {
		return fmt.Errorf("no graphical D-Bus session bus; use the GUI Connections page or 'gosched service status': %w", err)
	}
	if err := bus.Close(); err != nil {
		return fmt.Errorf("close D-Bus preflight connection: %w", err)
	}
	configPath, err := service.InstalledConfigPath()
	if err != nil {
		return err
	}
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("load local daemon configuration: %w", err)
	}
	local := client.New(ipc.Endpoint(cfg))
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	icon, err := loadIcon()
	if err != nil {
		return err
	}
	i := &indicator{tray: systray.New(), ctx: ctx}
	i.monitor = desktopcontrol.NewMonitor(func(checkCtx context.Context) error {
		_, healthErr := local.Health(checkCtx)
		return healthErr
	})
	menu := systray.NewMenu()
	i.status = menu.Add("This computer: Checking", nil)
	i.status.SetDisabled(true)
	i.feedback = menu.Add("Local service controls", nil)
	i.feedback.SetDisabled(true)
	menu.AddSeparator()
	menu.Add("Open go-schedule", func() { go i.openGUI() })
	i.start = menu.Add("Start service", func() { i.startAction("start") })
	i.stop = menu.Add("Stop service", func() { i.startAction("stop") })
	i.restart = menu.Add("Restart service", func() { i.startAction("restart") })
	menu.AddSeparator()
	menu.Add("Quit indicator", cancel)
	i.tray.SetIcon(icon).SetTooltip("go-schedule: Checking local service").SetMenu(menu).OnClick(func() { go i.openGUI() }).Show()
	i.refresh()
	done := make(chan struct{})
	go func() {
		defer close(done)
		ticker := time.NewTicker(4 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				i.refresh()
			}
		}
	}()
	go func() {
		<-ctx.Done()
		<-done
		i.mu.Lock()
		i.closing = true
		i.mu.Unlock()
		i.actions.Wait()
		i.tray.Remove()
	}()
	err = i.tray.Run()
	cancel()
	<-done
	return err
}

func loadIcon() ([]byte, error) {
	exe, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("locate indicator executable: %w", err)
	}
	iconRel := filepath.Join("share", "icons", "hicolor", "32x32", "apps", "go-schedule.png")
	for _, candidate := range []string{
		filepath.Join(filepath.Dir(exe), iconRel),
		filepath.Join("/usr/local", iconRel),
		filepath.Join("/usr", iconRel),
		filepath.Join("brand", "platform", "linux", "hicolor", "32x32", "apps", "go-schedule.png"),
	} {
		icon, readErr := os.ReadFile(candidate)
		if readErr == nil {
			return icon, nil
		}
	}
	return nil, fmt.Errorf("compact go-schedule icon is unavailable; install the desktop archive's share/icons tree")
}

func stateLabel(state string) string {
	switch state {
	case "not_installed":
		return "Not installed"
	case "unreachable":
		return "Unreachable"
	case "unknown":
		return "Unknown"
	default:
		if state == "" {
			return "Unknown"
		}
		return strings.ToUpper(state[:1]) + state[1:]
	}
}

func (i *indicator) refresh() {
	if i.ctx.Err() != nil {
		return
	}
	snapshot := i.monitor.Observe(i.ctx)
	if i.ctx.Err() != nil {
		return
	}
	i.mu.Lock()
	pending, message := i.pending, i.message
	if !pending && time.Now().After(i.messageUntil) {
		message = ""
	}
	i.mu.Unlock()
	i.status.SetLabel("This computer: " + stateLabel(snapshot.State))
	if message == "" {
		message = snapshot.Detail
	}
	i.feedback.SetLabel(message)
	i.tray.SetTooltip("go-schedule: " + stateLabel(snapshot.State) + " - " + snapshot.Detail)
	canStart, canStop := availableActions(snapshot)
	i.start.SetDisabled(pending || !canStart)
	i.stop.SetDisabled(pending || !canStop)
	i.restart.SetDisabled(pending || !canStop)
}

func availableActions(snapshot desktopcontrol.Snapshot) (bool, bool) {
	return snapshot.State == "stopped", snapshot.SCMState == "running"
}

func (i *indicator) startAction(action string) {
	i.mu.Lock()
	if i.ctx.Err() != nil || i.closing {
		i.mu.Unlock()
		return
	}
	i.actions.Add(1)
	i.mu.Unlock()
	go func() {
		defer i.actions.Done()
		i.perform(action)
	}()
}

func (i *indicator) perform(action string) {
	i.mu.Lock()
	if i.pending {
		i.mu.Unlock()
		return
	}
	i.pending = true
	i.message = "Requesting " + action + "..."
	i.messageUntil = time.Time{}
	i.mu.Unlock()
	defer func() {
		i.mu.Lock()
		i.pending = false
		i.mu.Unlock()
		i.refresh()
	}()
	i.refresh()
	if action == "stop" || action == "restart" {
		confirmed, err := confirmImpact(i.ctx, action)
		if err != nil {
			i.setMessage(err.Error())
			return
		}
		if !confirmed {
			i.setMessage("Service action cancelled. No change was requested.")
			return
		}
	}
	result := i.monitor.RequestAction(i.ctx, action)
	i.setMessage(result.Message)
}

func (i *indicator) setMessage(message string) {
	i.mu.Lock()
	i.message = message
	i.messageUntil = time.Now().Add(10 * time.Second)
	i.mu.Unlock()
}

func confirmImpact(ctx context.Context, action string) (bool, error) {
	text := "Stop the local service? Scheduled tasks on this computer will cease until it is started again."
	if action == "restart" {
		text = "Restart the local service? Active tasks may be interrupted."
	}
	if path, err := exec.LookPath("zenity"); err == nil {
		return runConfirmation(ctx, path, "--question", "--title=go-schedule", "--text="+text)
	}
	if path, err := exec.LookPath("kdialog"); err == nil {
		return runConfirmation(ctx, path, "--title", "go-schedule", "--yesno", text)
	}
	return false, errors.New("no graphical confirmation dialog is available; use the GUI Connections page to stop or restart")
}

func runConfirmation(ctx context.Context, path string, args ...string) (bool, error) {
	confirmCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()
	err := exec.CommandContext(confirmCtx, path, args...).Run()
	if confirmCtx.Err() != nil {
		return false, fmt.Errorf("confirmation dialog timed out or was interrupted: %w", confirmCtx.Err())
	}
	if err == nil {
		return true, nil
	}
	var exit *exec.ExitError
	if errors.As(err, &exit) && exit.ExitCode() == 1 {
		return false, nil
	}
	return false, fmt.Errorf("could not show service impact confirmation: %w", err)
}

func (i *indicator) openGUI() {
	for attempt := 0; attempt < 3; attempt++ {
		found, err := desktopcontrol.SignalGUI()
		if err == nil && found {
			return
		}
		if err != nil {
			i.setMessage(fmt.Sprintf("Could not focus GUI: %v", err))
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	exe, err := os.Executable()
	if err != nil {
		i.setMessage(fmt.Sprintf("Could not locate GUI: %v", err))
		return
	}
	gui := filepath.Join(filepath.Dir(exe), "gosched-gui")
	if _, err := os.Stat(gui); err != nil {
		gui, err = exec.LookPath("gosched-gui")
		if err != nil {
			i.setMessage("Could not find gosched-gui; install the Linux desktop bundle.")
			return
		}
	}
	cmd := exec.Command(gui)
	if err := cmd.Start(); err != nil {
		i.setMessage(fmt.Sprintf("Could not open GUI: %v", err))
		return
	}
	if err := cmd.Process.Release(); err != nil {
		i.setMessage(fmt.Sprintf("Could not release GUI launcher: %v", err))
	}
}
