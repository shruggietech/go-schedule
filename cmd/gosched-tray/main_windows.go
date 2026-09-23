//go:build windows

// gosched-tray is the per-logon, windowless Windows notification-area
// companion. Its elevated helper mode performs only one validated SCM action.
package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"

	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/config"
	"github.com/shruggietech/go-schedule/internal/desktopcontrol"
	"github.com/shruggietech/go-schedule/internal/ipc"
	"github.com/shruggietech/go-schedule/internal/service"
)

const (
	wmDestroy      = 0x0002
	wmNull         = 0x0000
	wmLButtonUp    = 0x0202
	wmRButtonUp    = 0x0205
	wmContextMenu  = 0x007b
	wmApp          = 0x8000
	wmTray         = wmApp + 1
	wmRefresh      = wmApp + 2
	ninSelect      = 0x0400
	ninKeySelect   = 0x0401
	nimAdd         = 0
	nimModify      = 1
	nimDelete      = 2
	nimSetVersion  = 4
	nifMessage     = 1
	nifIcon        = 2
	nifTip         = 4
	imageIcon      = 1
	lrLoadFromFile = 0x10
	mfGrayed       = 1
	mfSeparator    = 0x800
	tpmRightButton = 2
	tpmReturnCmd   = 0x100
	mbYesNo        = 4
	mbIconWarning  = 0x30
	mbDefaultNo    = 0x100
	idYes          = 6
)

type point struct{ x, y int32 }
type msg struct {
	hwnd    uintptr
	message uint32
	wParam  uintptr
	lParam  uintptr
	time    uint32
	pt      point
	private uint32
}
type wndClass struct {
	size       uint32
	style      uint32
	wndProc    uintptr
	clsExtra   int32
	wndExtra   int32
	instance   uintptr
	icon       uintptr
	cursor     uintptr
	background uintptr
	menuName   *uint16
	className  *uint16
	smallIcon  uintptr
}
type notifyIconData struct {
	size        uint32
	hwnd        uintptr
	id          uint32
	flags       uint32
	callback    uint32
	icon        uintptr
	tip         [128]uint16
	state       uint32
	stateMask   uint32
	info        [256]uint16
	version     uint32
	infoTitle   [64]uint16
	infoFlags   uint32
	guid        [16]byte
	balloonIcon uintptr
}

var (
	user32               = windows.NewLazySystemDLL("user32.dll")
	shell32              = windows.NewLazySystemDLL("shell32.dll")
	procRegisterClass    = user32.NewProc("RegisterClassExW")
	procCreateWindow     = user32.NewProc("CreateWindowExW")
	procDefWindow        = user32.NewProc("DefWindowProcW")
	procDestroyWindow    = user32.NewProc("DestroyWindow")
	procGetMessage       = user32.NewProc("GetMessageW")
	procTranslateMessage = user32.NewProc("TranslateMessage")
	procDispatchMessage  = user32.NewProc("DispatchMessageW")
	procRegisterMessage  = user32.NewProc("RegisterWindowMessageW")
	procPostMessage      = user32.NewProc("PostMessageW")
	procPostQuit         = user32.NewProc("PostQuitMessage")
	procLoadImage        = user32.NewProc("LoadImageW")
	procDestroyIcon      = user32.NewProc("DestroyIcon")
	procCreateMenu       = user32.NewProc("CreatePopupMenu")
	procAppendMenu       = user32.NewProc("AppendMenuW")
	procSetMenuDefault   = user32.NewProc("SetMenuDefaultItem")
	procDestroyMenu      = user32.NewProc("DestroyMenu")
	procTrackMenu        = user32.NewProc("TrackPopupMenu")
	procGetCursor        = user32.NewProc("GetCursorPos")
	procSetForeground    = user32.NewProc("SetForegroundWindow")
	procMessageBox       = user32.NewProc("MessageBoxW")
	procSystemParameters = user32.NewProc("SystemParametersInfoW")
	procGetSysColor      = user32.NewProc("GetSysColor")
	procNotifyIcon       = shell32.NewProc("Shell_NotifyIconW")
	currentTray          *tray
)

type tray struct {
	hwnd           uintptr
	icon           uintptr
	lightTheme     bool
	taskbarCreated uint32
	monitor        desktopcontrol.Monitor
	cancel         context.CancelFunc
	mu             sync.RWMutex
	snapshot       desktopcontrol.Snapshot
	action         string
	lastMessage    string
}

func main() {
	if len(os.Args) == 3 && os.Args[1] == "--elevated-service-action" {
		action := os.Args[2]
		if action != "start" && action != "stop" && action != "restart" {
			os.Exit(2)
		}
		if _, err := service.Control(action, "", nil); err != nil {
			os.Exit(3)
		}
		return
	}
	if len(os.Args) != 1 {
		os.Exit(2)
	}
	instance, owner, err := desktopcontrol.ClaimTray()
	if err != nil {
		os.Exit(1)
	}
	if !owner {
		return
	}
	defer instance.Close() //nolint:errcheck // process-lifetime mutex
	cfg, err := config.Load("")
	if err != nil {
		os.Exit(1)
	}
	local := client.New(ipc.Endpoint(cfg))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	t := &tray{monitor: desktopcontrol.NewMonitor(func(ctx context.Context) error {
		_, err := local.Health(ctx)
		return err
	}), cancel: cancel}
	if err := t.run(ctx); err != nil {
		os.Exit(1)
	}
}

func (t *tray) run(ctx context.Context) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	currentTray = t
	className := windows.StringToUTF16Ptr("GoScheduleTrayWindow")
	class := wndClass{size: uint32(unsafe.Sizeof(wndClass{})), wndProc: syscall.NewCallback(wndProc), className: className}
	if r, _, err := procRegisterClass.Call(uintptr(unsafe.Pointer(&class))); r == 0 {
		return fmt.Errorf("register tray window: %w", err)
	}
	hwnd, _, err := procCreateWindow.Call(0, uintptr(unsafe.Pointer(className)), uintptr(unsafe.Pointer(className)), 0, 0, 0, 0, 0, 0, 0, 0, 0)
	if hwnd == 0 {
		return fmt.Errorf("create tray window: %w", err)
	}
	t.hwnd = hwnd
	t.taskbarCreated = uint32(registerMessage("TaskbarCreated"))
	t.loadIcon()
	t.refresh(ctx)
	t.updateIcon(nimAdd)
	go t.poll(ctx)
	var message msg
	for {
		result, _, callErr := procGetMessage.Call(uintptr(unsafe.Pointer(&message)), 0, 0, 0)
		if int32(result) == -1 {
			return fmt.Errorf("read tray message: %w", callErr)
		}
		if result == 0 {
			return nil
		}
		_, _, _ = procTranslateMessage.Call(uintptr(unsafe.Pointer(&message))) // message translation has no error contract
		_, _, _ = procDispatchMessage.Call(uintptr(unsafe.Pointer(&message)))  // dispatch result is window-procedure defined
	}
}

func registerMessage(value string) uintptr {
	r, _, _ := procRegisterMessage.Call(uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(value))))
	return r
}

func (t *tray) poll(ctx context.Context) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			t.refresh(ctx)
			_, _, _ = procPostMessage.Call(t.hwnd, wmRefresh, 0, 0) // next polling tick also refreshes
		}
	}
}

func (t *tray) refresh(ctx context.Context) {
	snapshot := t.monitor.Observe(ctx)
	t.mu.Lock()
	t.snapshot = snapshot
	t.mu.Unlock()
}

func (t *tray) status() (desktopcontrol.Snapshot, string) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.snapshot, t.action
}

func (t *tray) loadIcon() {
	exe, err := os.Executable()
	if err != nil {
		return
	}
	t.lightTheme = lightTaskbar()
	name := "go-schedule-light.ico"
	if t.lightTheme {
		name = "go-schedule-dark.ico"
	}
	path := filepath.Join(filepath.Dir(exe), name)
	if _, err := os.Stat(path); err != nil {
		path = filepath.Join(filepath.Dir(exe), "go-schedule.ico")
	}
	ptr, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return
	}
	t.icon, _, _ = procLoadImage.Call(0, uintptr(unsafe.Pointer(ptr)), imageIcon, 0, 0, lrLoadFromFile)
}

func lightTaskbar() bool {
	// In high-contrast mode use the Windows text color as the actual surface
	// contrast reference instead of the ordinary taskbar theme preference.
	type highContrast struct {
		size          uint32
		flags         uint32
		defaultScheme *uint16
	}
	hc := highContrast{size: uint32(unsafe.Sizeof(highContrast{}))}
	if ok, _, _ := procSystemParameters.Call(0x42, uintptr(hc.size), uintptr(unsafe.Pointer(&hc)), 0); ok != 0 && hc.flags&1 != 0 {
		color, _, _ := procGetSysColor.Call(8) // COLOR_WINDOWTEXT
		red, green, blue := color&0xff, (color>>8)&0xff, (color>>16)&0xff
		return red+green+blue < 384
	}
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Themes\Personalize`, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer key.Close() //nolint:errcheck // read-only key
	value, _, err := key.GetIntegerValue("SystemUsesLightTheme")
	return err == nil && value != 0
}

func (t *tray) updateIcon(op uintptr) {
	snapshot, action := t.status()
	data := notifyIconData{size: uint32(unsafe.Sizeof(notifyIconData{})), hwnd: t.hwnd, id: 1, flags: nifMessage | nifIcon | nifTip, callback: wmTray, icon: t.icon}
	copy(data.tip[:], windows.StringToUTF16(tooltipText(snapshot, action)))
	result, _, _ := procNotifyIcon.Call(op, uintptr(unsafe.Pointer(&data)))
	if result == 0 && op == nimModify {
		// Explorer can be unavailable at logon or recreate its notification area
		// without a timely TaskbarCreated broadcast. The next refresh retries.
		result, _, _ = procNotifyIcon.Call(nimAdd, uintptr(unsafe.Pointer(&data)))
		op = nimAdd
	}
	if result != 0 && op == nimAdd {
		data.version = 4
		_, _, _ = procNotifyIcon.Call(nimSetVersion, uintptr(unsafe.Pointer(&data))) // default callback version remains usable
	}
}

func tooltipText(snapshot desktopcontrol.Snapshot, action string) string {
	label := strings.ReplaceAll(snapshot.State, "_", " ")
	if action != "" {
		label = action + " in progress"
	}
	return "go-schedule: This computer - " + label
}

func availableActions(snapshot desktopcontrol.Snapshot, pending string) []string {
	if pending != "" {
		return nil
	}
	if snapshot.State == "stopped" {
		return []string{"start"}
	}
	if snapshot.SCMState == "running" {
		return []string{"stop", "restart"}
	}
	return nil
}

func wndProc(hwnd uintptr, message uint32, wParam, lParam uintptr) uintptr {
	t := currentTray
	if t == nil {
		r, _, _ := procDefWindow.Call(hwnd, uintptr(message), wParam, lParam)
		return r
	}
	switch message {
	case wmTray:
		switch uint32(lParam) & 0xffff {
		case wmLButtonUp, ninSelect, ninKeySelect:
			go t.openGUI()
		case wmRButtonUp, wmContextMenu:
			t.menu()
		}
		return 0
	case wmRefresh:
		if t.lightTheme != lightTaskbar() {
			if t.icon != 0 {
				_, _, _ = procDestroyIcon.Call(t.icon) // best-effort GDI cleanup
			}
			t.loadIcon()
		}
		t.updateIcon(nimModify)
		return 0
	case wmDestroy:
		t.updateIcon(nimDelete)
		t.cancel()
		if t.icon != 0 {
			_, _, _ = procDestroyIcon.Call(t.icon) // process exit also releases the icon handle
		}
		_, _, _ = procPostQuit.Call(0) // no failure return
		return 0
	default:
		if message == t.taskbarCreated {
			t.updateIcon(nimAdd)
			return 0
		}
	}
	r, _, _ := procDefWindow.Call(hwnd, uintptr(message), wParam, lParam)
	return r
}

func (t *tray) menu() {
	menu, _, _ := procCreateMenu.Call()
	if menu == 0 {
		return
	}
	defer func() { _, _, _ = procDestroyMenu.Call(menu) }() // native menu cleanup
	snapshot, action := t.status()
	appendMenu(menu, 0, 1, "Open go-schedule")
	_, _, _ = procSetMenuDefault.Call(menu, 1, 0) // Open is the one native default command
	appendMenu(menu, mfGrayed, 2, "This computer: "+strings.ReplaceAll(snapshot.State, "_", " "))
	appendMenu(menu, mfSeparator, 0, "")
	for _, available := range availableActions(snapshot, action) {
		switch available {
		case "start":
			appendMenu(menu, 0, 3, "Start service")
		case "stop":
			appendMenu(menu, 0, 4, "Stop service")
		case "restart":
			appendMenu(menu, 0, 5, "Restart service")
		}
	}
	appendMenu(menu, mfSeparator, 0, "")
	appendMenu(menu, 0, 6, "Quit companion")
	var pos point
	if ok, _, _ := procGetCursor.Call(uintptr(unsafe.Pointer(&pos))); ok == 0 {
		return
	}
	_, _, _ = procSetForeground.Call(t.hwnd) // Windows may deny focus transfer; menu still opens
	selected, _, _ := procTrackMenu.Call(menu, tpmRightButton|tpmReturnCmd, uintptr(pos.x), uintptr(pos.y), 0, t.hwnd, 0)
	_, _, _ = procPostMessage.Call(t.hwnd, wmNull, 0, 0) // documented TrackPopupMenu dismissal sequence
	switch selected {
	case 1:
		go t.openGUI()
	case 3:
		go t.perform("start")
	case 4:
		if t.confirm("Stop the local service? Scheduled tasks on this computer will cease until it is started again.") {
			go t.perform("stop")
		}
	case 5:
		if t.confirm("Restart the local service? Active tasks may be interrupted.") {
			go t.perform("restart")
		}
	case 6:
		_, _, _ = procDestroyWindow.Call(t.hwnd) // WM_DESTROY performs cleanup
	}
}

func appendMenu(menu, flags, id uintptr, label string) {
	var ptr *uint16
	if label != "" {
		ptr = windows.StringToUTF16Ptr(label)
	}
	_, _, _ = procAppendMenu.Call(menu, flags, id, uintptr(unsafe.Pointer(ptr))) // a failed item is omitted from the menu
}

func (t *tray) confirm(message string) bool {
	text := windows.StringToUTF16Ptr(message)
	title := windows.StringToUTF16Ptr("go-schedule service")
	result, _, _ := procMessageBox.Call(t.hwnd, uintptr(unsafe.Pointer(text)), uintptr(unsafe.Pointer(title)), mbYesNo|mbIconWarning|mbDefaultNo)
	return result == idYes
}

func (t *tray) perform(action string) {
	t.mu.Lock()
	if t.action != "" {
		t.mu.Unlock()
		return
	}
	t.action = action
	t.mu.Unlock()
	_, _, _ = procPostMessage.Call(t.hwnd, wmRefresh, 0, 0) // next polling tick also refreshes
	result := t.monitor.RequestAction(context.Background(), action)
	t.mu.Lock()
	t.action = ""
	t.snapshot = result.Snapshot
	t.lastMessage = result.Message
	t.mu.Unlock()
	_, _, _ = procPostMessage.Call(t.hwnd, wmRefresh, 0, 0) // next polling tick also refreshes
	if result.Outcome != "accepted" {
		text := windows.StringToUTF16Ptr(result.Message)
		title := windows.StringToUTF16Ptr("go-schedule service")
		_, _, _ = procMessageBox.Call(t.hwnd, uintptr(unsafe.Pointer(text)), uintptr(unsafe.Pointer(title)), mbIconWarning) // failure remains in observed state
	}
}

func (t *tray) openGUI() {
	found, err := desktopcontrol.SignalGUI()
	if err == nil && found {
		return
	}
	exe, err := os.Executable()
	if err != nil {
		return
	}
	gui := filepath.Join(filepath.Dir(exe), "gosched-gui.exe")
	cmd := exec.Command(gui)
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: windows.CREATE_NO_WINDOW}
	if err := cmd.Start(); err == nil {
		_ = cmd.Process.Release()
	}
}
