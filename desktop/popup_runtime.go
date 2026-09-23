package main

import (
	"context"

	"github.com/shruggietech/go-schedule/desktop/popups"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type popupPlatform interface {
	Initialize(context.Context, func(string)) error
	Available() bool
	Authorized() bool
	RequestAuthorization() bool
	Cleanup()
	popups.Presenter
}

type wailsPopup struct{ ctx context.Context }

func (p *wailsPopup) Initialize(ctx context.Context, onResponse func(string)) error {
	if err := wailsRuntime.InitializeNotifications(ctx); err != nil {
		return err
	}
	p.ctx = ctx
	wailsRuntime.OnNotificationResponse(ctx, func(result wailsRuntime.NotificationResult) {
		if result.Error == nil && result.Response.ID != "" {
			onResponse(result.Response.ID)
		}
	})
	return nil
}
func (p *wailsPopup) Available() bool {
	return p.ctx != nil && wailsRuntime.IsNotificationAvailable(p.ctx)
}
func (p *wailsPopup) Authorized() bool {
	if !p.Available() {
		return false
	}
	ok, err := wailsRuntime.CheckNotificationAuthorization(p.ctx)
	return err == nil && ok
}
func (p *wailsPopup) RequestAuthorization() bool {
	if !p.Available() {
		return false
	}
	ok, err := wailsRuntime.RequestNotificationAuthorization(p.ctx)
	return err == nil && ok
}
func (p *wailsPopup) Present(value popups.Notice) error {
	return wailsRuntime.SendNotification(p.ctx, wailsRuntime.NotificationOptions{ID: value.ID, Title: value.Title, Body: value.Body})
}
func (p *wailsPopup) Cleanup() {
	if p.ctx != nil {
		wailsRuntime.CleanupNotifications(p.ctx)
	}
}

func (a *App) domReady(ctx context.Context) {
	if a.popupRuntime == nil {
		return
	}
	if err := a.popupRuntime.Initialize(ctx, func(id string) {
		if a.popups == nil {
			return
		}
		intent, ok := a.popups.Intent(id)
		if !ok {
			return
		}
		wailsShowWindow(ctx)
		if a.emitter != nil {
			a.emitter.Emit(ctx, popupEventName, intent)
		}
	}); err != nil {
		return
	}
	if a.popups != nil {
		a.popups.Start(ctx)
	}
}
