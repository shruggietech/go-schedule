//go:build windows

package winuninstall

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

const (
	profileListKey = `SOFTWARE\Microsoft\Windows NT\CurrentVersion\ProfileList`
	profileLeaf    = `AppData\Roaming\fyne\tech.shruggie.goschedule`
	resultKey      = `Software\ShruggieTech\go-schedule-uninstall`
	resultFolder   = `ShruggieTech\go-schedule-uninstall\b6f3c2e1-7a4d-4c9e-9b2a-1f6d8e5a0c34`
)

type windowsBackend struct {
	programData string
	resultPath  string
}

// Wipe derives the fixed Windows-owned roots and executes the bounded cleanup contract.
func Wipe() Result {
	programData, err := windows.KnownFolderPath(windows.FOLDERID_ProgramData, 0)
	if err != nil {
		return Result{Schema: resultSchema, State: StateInternalError, Error: fmt.Sprintf("resolve ProgramData: %v", err)}
	}
	backend := &windowsBackend{
		programData: filepath.Clean(programData),
		resultPath:  filepath.Join(programData, resultFolder, "cleanup-result.json"),
	}
	return Run(backend)
}

func (b *windowsBackend) Discover() (targets []Target, resultErr error) {
	targets = []Target{{
		Kind: TargetMachine, Path: filepath.Join(b.programData, "goschedule"),
		base: b.programData, relative: "goschedule",
	}}
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, profileListKey, registry.ENUMERATE_SUB_KEYS|registry.WOW64_64KEY)
	if err != nil {
		return nil, fmt.Errorf("open registered profile list: %w", err)
	}
	defer func() {
		if err := key.Close(); err != nil {
			resultErr = errors.Join(resultErr, fmt.Errorf("close registered profile list: %w", err))
		}
	}()
	sids, err := key.ReadSubKeyNames(-1)
	if err != nil {
		return nil, fmt.Errorf("enumerate registered profiles: %w", err)
	}
	sort.Strings(sids)
	seen := map[string]bool{strings.ToLower(filepath.Clean(targets[0].Path)): true}
	for _, sid := range sids {
		profileKey, err := registry.OpenKey(key, sid, registry.QUERY_VALUE)
		if err != nil {
			return nil, fmt.Errorf("open registered profile %s: %w", sid, err)
		}
		rawPath, _, valueErr := profileKey.GetStringValue("ProfileImagePath")
		profileCloseErr := profileKey.Close()
		if profileCloseErr != nil {
			return nil, errors.Join(valueErr, fmt.Errorf("close registered profile %s: %w", sid, profileCloseErr))
		}
		if errors.Is(valueErr, registry.ErrNotExist) || strings.TrimSpace(rawPath) == "" {
			continue
		}
		if valueErr != nil {
			return nil, fmt.Errorf("read registered profile %s path: %w", sid, valueErr)
		}
		profilePath, err := expandWindowsEnvironment(rawPath)
		if err != nil {
			return nil, fmt.Errorf("expand registered profile %s path: %w", sid, err)
		}
		profilePath = filepath.Clean(profilePath)
		candidate := filepath.Join(profilePath, profileLeaf)
		identity := strings.ToLower(filepath.Clean(candidate))
		if seen[identity] {
			continue
		}
		seen[identity] = true
		targets = append(targets, Target{
			Kind: TargetProfile, SID: sid, Path: candidate,
			base: profilePath, relative: profileLeaf,
		})
	}
	return targets, nil
}

func (b *windowsBackend) Preflight(target Target) (bool, error) {
	if err := validateLexicalPath(target.Path); err != nil {
		return false, err
	}
	expected := filepath.Clean(filepath.Join(target.base, target.relative))
	if !strings.EqualFold(filepath.Clean(target.Path), expected) {
		return false, fmt.Errorf("candidate is outside its declared owned root")
	}
	volume := filepath.VolumeName(target.Path)
	rootPointer, err := windows.UTF16PtrFromString(volume + `\`)
	if err != nil {
		return false, fmt.Errorf("encode volume root: %w", err)
	}
	if driveType := windows.GetDriveType(rootPointer); driveType != windows.DRIVE_FIXED {
		return false, fmt.Errorf("candidate volume is not a fixed local drive")
	}
	if err := rejectReparseAncestors(target.Path); err != nil {
		return false, err
	}
	if _, err := os.Lstat(target.Path); errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("inspect owned root: %w", err)
	}
	canonicalBase, err := canonicalPath(target.base)
	if err != nil {
		return false, fmt.Errorf("canonicalize trusted base: %w", err)
	}
	canonicalTarget, err := canonicalPath(target.Path)
	if err != nil {
		return false, fmt.Errorf("canonicalize owned root: %w", err)
	}
	canonicalExpected := filepath.Clean(filepath.Join(canonicalBase, target.relative))
	if !strings.EqualFold(canonicalTarget, canonicalExpected) {
		return false, fmt.Errorf("canonical owned root escaped its trusted base")
	}
	return true, nil
}

func (b *windowsBackend) Remove(target Target) (resultErr error) {
	if err := validateLexicalPath(target.Path); err != nil {
		return err
	}
	volumeRoot := filepath.VolumeName(target.Path) + `\`
	relative, err := filepath.Rel(volumeRoot, target.Path)
	if err != nil || relative == "." || relative == ".." || strings.HasPrefix(relative, `..\`) {
		return fmt.Errorf("derive volume-relative owned root")
	}
	parent, err := os.OpenRoot(volumeRoot)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("open owned-root parent: %w", err)
	}
	defer func() {
		if err := parent.Close(); err != nil {
			resultErr = errors.Join(resultErr, fmt.Errorf("close owned-root parent: %w", err))
		}
	}()
	if err := parent.RemoveAll(relative); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("remove without following reparse entries: %w", err)
	}
	return nil
}

func (b *windowsBackend) WriteResult(result Result) (resultErr error) {
	if err := secureResultDirectory(filepath.Dir(b.resultPath)); err != nil {
		return err
	}
	if err := writeLedger(b.resultPath, result); err != nil {
		return err
	}
	key, _, err := registry.CreateKey(registry.LOCAL_MACHINE, resultKey, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("create cleanup status key: %w", err)
	}
	defer func() {
		if err := key.Close(); err != nil {
			resultErr = errors.Join(resultErr, fmt.Errorf("close cleanup status key: %w", err))
		}
	}()
	if err := key.SetStringValue("State", string(result.State)); err != nil {
		return fmt.Errorf("write cleanup state: %w", err)
	}
	if err := key.SetDWordValue("RemainingCount", uint32(result.Remaining)); err != nil {
		return fmt.Errorf("write cleanup remaining count: %w", err)
	}
	if err := key.SetStringValue("ReportPath", b.resultPath); err != nil {
		return fmt.Errorf("write cleanup report path: %w", err)
	}
	return nil
}

func (b *windowsBackend) ClearResult() (resultErr error) {
	if err := registry.DeleteKey(registry.LOCAL_MACHINE, resultKey); err != nil && !errors.Is(err, registry.ErrNotExist) {
		return fmt.Errorf("delete cleanup status key: %w", err)
	}
	volumeRoot := filepath.VolumeName(b.programData) + `\`
	relative, err := filepath.Rel(volumeRoot, filepath.Join(b.programData, resultFolder))
	if err != nil || relative == "." || relative == ".." || strings.HasPrefix(relative, `..\`) {
		return fmt.Errorf("derive volume-relative cleanup result path")
	}
	root, err := os.OpenRoot(volumeRoot)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("open ProgramData for result cleanup: %w", err)
	}
	defer func() {
		if err := root.Close(); err != nil {
			resultErr = errors.Join(resultErr, fmt.Errorf("close ProgramData cleanup root: %w", err))
		}
	}()
	if err := root.RemoveAll(relative); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("delete completed cleanup ledger: %w", err)
	}
	return nil
}

func validateLexicalPath(path string) error {
	if strings.TrimSpace(path) == "" {
		return fmt.Errorf("candidate path is empty")
	}
	if strings.HasPrefix(path, `\\`) || strings.HasPrefix(path, `\\?\`) || strings.HasPrefix(path, `\\.\`) {
		return fmt.Errorf("candidate uses an UNC or device path")
	}
	if !filepath.IsAbs(path) {
		return fmt.Errorf("candidate path is not absolute")
	}
	if strings.Count(path, ":") != 1 || len(path) < 3 || path[1] != ':' {
		return fmt.Errorf("candidate contains a device or alternate data stream")
	}
	cleaned := filepath.Clean(path)
	if strings.EqualFold(cleaned, filepath.VolumeName(cleaned)+`\`) {
		return fmt.Errorf("candidate resolves to a volume root")
	}
	for _, part := range strings.FieldsFunc(path[2:], func(r rune) bool { return r == '\\' || r == '/' }) {
		if part == "." || part == ".." {
			return fmt.Errorf("candidate contains a dot segment")
		}
		if strings.TrimRight(part, ". ") != part || strings.ContainsRune(part, 0) {
			return fmt.Errorf("candidate contains a malformed path component")
		}
		base := strings.ToUpper(strings.TrimSuffix(part, filepath.Ext(part)))
		if base == "CON" || base == "PRN" || base == "AUX" || base == "NUL" ||
			(len(base) == 4 && (strings.HasPrefix(base, "COM") || strings.HasPrefix(base, "LPT")) && base[3] >= '1' && base[3] <= '9') {
			return fmt.Errorf("candidate contains a reserved device component")
		}
	}
	return nil
}

func rejectReparseAncestors(path string) error {
	cleaned := filepath.Clean(path)
	volume := filepath.VolumeName(cleaned)
	remainder := strings.TrimPrefix(cleaned, volume)
	current := volume + `\`
	for _, part := range strings.FieldsFunc(remainder, func(r rune) bool { return r == '\\' || r == '/' }) {
		current = filepath.Join(current, part)
		pointer, err := windows.UTF16PtrFromString(current)
		if err != nil {
			return fmt.Errorf("encode candidate ancestor: %w", err)
		}
		attributes, err := windows.GetFileAttributes(pointer)
		if errors.Is(err, windows.ERROR_FILE_NOT_FOUND) || errors.Is(err, windows.ERROR_PATH_NOT_FOUND) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("inspect candidate ancestor %q: %w", current, err)
		}
		if attributes&windows.FILE_ATTRIBUTE_REPARSE_POINT != 0 {
			return fmt.Errorf("candidate root or ancestor is a reparse point: %s", current)
		}
	}
	return nil
}

func canonicalPath(path string) (resolvedPath string, resultErr error) {
	pointer, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return "", err
	}
	handle, err := windows.CreateFile(pointer, windows.FILE_READ_ATTRIBUTES,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE|windows.FILE_SHARE_DELETE,
		nil, windows.OPEN_EXISTING,
		windows.FILE_FLAG_BACKUP_SEMANTICS|windows.FILE_FLAG_OPEN_REPARSE_POINT, 0)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := windows.CloseHandle(handle); err != nil {
			resultErr = errors.Join(resultErr, fmt.Errorf("close canonical path handle: %w", err))
		}
	}()
	buffer := make([]uint16, 512)
	for {
		length, err := windows.GetFinalPathNameByHandle(handle, &buffer[0], uint32(len(buffer)), 0)
		if err != nil {
			return "", err
		}
		if int(length) < len(buffer) {
			resolved := windows.UTF16ToString(buffer[:length])
			resolved = strings.TrimPrefix(resolved, `\\?\`)
			return filepath.Clean(resolved), nil
		}
		buffer = make([]uint16, length+1)
	}
}

func expandWindowsEnvironment(value string) (string, error) {
	source, err := windows.UTF16PtrFromString(value)
	if err != nil {
		return "", err
	}
	length, err := windows.ExpandEnvironmentStrings(source, nil, 0)
	if err != nil {
		return "", err
	}
	buffer := make([]uint16, length)
	if _, err := windows.ExpandEnvironmentStrings(source, &buffer[0], uint32(len(buffer))); err != nil {
		return "", err
	}
	return windows.UTF16ToString(buffer), nil
}

func secureResultDirectory(path string) (resultErr error) {
	if err := validateLexicalPath(path); err != nil {
		return fmt.Errorf("validate cleanup result directory: %w", err)
	}
	if err := rejectReparseAncestors(path); err != nil {
		return fmt.Errorf("validate cleanup result directory: %w", err)
	}
	volumeRoot := filepath.VolumeName(path) + `\`
	relative, err := filepath.Rel(volumeRoot, path)
	if err != nil || relative == "." || relative == ".." || strings.HasPrefix(relative, `..\`) {
		return fmt.Errorf("derive volume-relative cleanup result directory")
	}
	root, err := os.OpenRoot(volumeRoot)
	if err != nil {
		return fmt.Errorf("open cleanup result volume: %w", err)
	}
	defer func() {
		if err := root.Close(); err != nil {
			resultErr = errors.Join(resultErr, fmt.Errorf("close cleanup result volume: %w", err))
		}
	}()
	if err := root.MkdirAll(relative, 0o700); err != nil {
		return fmt.Errorf("create cleanup result directory: %w", err)
	}
	if err := rejectReparseAncestors(path); err != nil {
		return fmt.Errorf("validate created cleanup result directory: %w", err)
	}
	descriptor, err := windows.SecurityDescriptorFromString("D:P(A;OICI;FA;;;SY)(A;OICI;FA;;;BA)")
	if err != nil {
		return fmt.Errorf("create cleanup result security descriptor: %w", err)
	}
	dacl, _, err := descriptor.DACL()
	if err != nil {
		return fmt.Errorf("read cleanup result DACL: %w", err)
	}
	information := windows.SECURITY_INFORMATION(windows.DACL_SECURITY_INFORMATION | windows.PROTECTED_DACL_SECURITY_INFORMATION)
	if err := windows.SetNamedSecurityInfo(path, windows.SE_FILE_OBJECT, information, nil, nil, dacl, nil); err != nil {
		return fmt.Errorf("protect cleanup result directory: %w", err)
	}
	return nil
}
