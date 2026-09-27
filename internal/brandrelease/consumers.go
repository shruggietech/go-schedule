package brandrelease

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// ReadConsumers loads a repository mapping. Version 2 supports controlled transforms.
func ReadConsumers(filename string) (ConsumerMap, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return ConsumerMap{}, fmt.Errorf("read consumer map: %w", err)
	}
	return ParseConsumers(data)
}

// ParseConsumers parses the exact map bytes that an importer will install.
func ParseConsumers(data []byte) (ConsumerMap, error) {
	var result ConsumerMap
	if err := json.Unmarshal(data, &result); err != nil {
		return result, fmt.Errorf("parse consumer map: %w", err)
	}
	if result.Version != 1 && result.Version != 2 {
		return result, fmt.Errorf("unsupported consumer map version %d", result.Version)
	}
	return result, nil
}

// ValidateConsumers checks mapping identity and destination boundaries before writes.
func ValidateConsumers(root string, kit *Kit, consumers ConsumerMap) error {
	seenSource := map[string]bool{}
	seenTarget := map[string]bool{}
	targetSource := map[string]string{}
	var docsManifest *Mapping
	for _, mapping := range consumers.Mappings {
		if err := SafePath(mapping.Source); err != nil {
			return fmt.Errorf("consumer source: %w", err)
		}
		if _, ok := kit.Files[mapping.Source]; !ok {
			return fmt.Errorf("consumer source %q is missing from brand kit", mapping.Source)
		}
		if seenSource[strings.ToLower(mapping.Source)] || len(mapping.Targets) == 0 || strings.TrimSpace(mapping.Purpose) == "" {
			return fmt.Errorf("consumer source %q is duplicated or incomplete", mapping.Source)
		}
		seenSource[strings.ToLower(mapping.Source)] = true
		if mapping.Transform != "" && mapping.Transform != "docs-site-manifest" {
			return fmt.Errorf("unsupported consumer transform %q", mapping.Transform)
		}
		if mapping.Transform == "docs-site-manifest" && mapping.Source != "favicons/site.webmanifest" {
			return fmt.Errorf("docs-site-manifest transform has wrong source")
		}
		if mapping.Transform == "docs-site-manifest" {
			docsManifest = &mapping
		}
		for _, target := range mapping.Targets {
			if err := SafePath(target); err != nil {
				return fmt.Errorf("consumer target: %w", err)
			}
			if !allowedTarget(target) || reservedTarget(target, kit.Pin.Package+".zip") {
				return fmt.Errorf("consumer target %q is outside approved brand destinations", target)
			}
			if seenTarget[strings.ToLower(target)] {
				return fmt.Errorf("duplicate consumer target %q", target)
			}
			seenTarget[strings.ToLower(target)] = true
			targetSource[target] = mapping.Source
			if mapping.Transform == "docs-site-manifest" && target != "docs/assets/favicons/site.webmanifest" {
				return fmt.Errorf("docs-site-manifest transform has wrong target")
			}
			if err := CheckNoSymlink(root, target); err != nil {
				return err
			}
		}
	}
	if docsManifest != nil {
		manifest, err := Render(kit, *docsManifest)
		if err != nil {
			return err
		}
		var doc struct {
			Icons []struct {
				Src string `json:"src"`
			} `json:"icons"`
		}
		if err := json.Unmarshal(manifest, &doc); err != nil {
			return fmt.Errorf("parse docs web manifest: %w", err)
		}
		for _, icon := range doc.Icons {
			source := "favicons/" + icon.Src
			if targetSource["docs/assets/favicons/"+icon.Src] != source {
				return fmt.Errorf("docs web icon %q has no matching consumer", icon.Src)
			}
		}
	}
	return nil
}

func allowedTarget(name string) bool {
	return strings.HasPrefix(name, "brand/") || strings.HasPrefix(name, "docs/assets/brand/") || strings.HasPrefix(name, "docs/assets/favicons/") || strings.HasPrefix(name, "docs/assets/fonts/") || strings.HasPrefix(name, "desktop/build/") || strings.HasPrefix(name, "desktop/frontend/public/")
}

func reservedTarget(name, archive string) bool {
	switch name {
	case "brand/source.json", "brand/repository-consumers.json", "brand/REPOSITORY.md", "brand/README.md", "brand/" + archive, "brand/platform/linux/go-schedule.desktop", "brand/platform/linux/go-schedule-indicator.desktop":
		return true
	}
	return false
}

// CheckNoSymlink rejects a destination whose existing path or ancestor is a symlink.
func CheckNoSymlink(root, relative string) error {
	parts := strings.Split(relative, "/")
	current := root
	for _, part := range parts {
		current = filepath.Join(current, part)
		info, err := os.Lstat(current)
		if os.IsNotExist(err) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("inspect destination %q: %w", relative, err)
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("destination %q crosses a symlink", relative)
		}
	}
	return nil
}

// Render returns the exact official bytes or the one documented docs manifest adaptation.
func Render(kit *Kit, mapping Mapping) ([]byte, error) {
	data, ok := kit.Files[mapping.Source]
	if !ok {
		return nil, fmt.Errorf("missing consumer source %q", mapping.Source)
	}
	if mapping.Transform == "" {
		return data, nil
	}
	var doc map[string]json.RawMessage
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("parse official web manifest: %w", err)
	}
	var icons []map[string]json.RawMessage
	if err := json.Unmarshal(doc["icons"], &icons); err != nil || len(icons) == 0 {
		return nil, fmt.Errorf("official web manifest has no usable icons: %v", err)
	}
	for i := range icons {
		var src string
		if err := json.Unmarshal(icons[i]["src"], &src); err != nil {
			return nil, fmt.Errorf("parse official web icon URL: %w", err)
		}
		if !strings.HasPrefix(src, "/") || strings.Contains(strings.TrimPrefix(src, "/"), "/") {
			return nil, fmt.Errorf("unexpected official web icon URL %q", src)
		}
		if _, ok := kit.Files["favicons/"+strings.TrimPrefix(src, "/")]; !ok {
			return nil, fmt.Errorf("web icon %q is missing from kit", src)
		}
		encodedSrc, err := json.Marshal(strings.TrimPrefix(src, "/"))
		if err != nil {
			return nil, fmt.Errorf("encode docs web icon URL: %w", err)
		}
		icons[i]["src"] = encodedSrc
	}
	encodedIcons, err := json.Marshal(icons)
	if err != nil {
		return nil, fmt.Errorf("encode docs web icons: %w", err)
	}
	doc["icons"] = encodedIcons
	encoded, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("encode docs web manifest: %w", err)
	}
	return append(encoded, '\n'), nil
}

// CheckInstalled validates every mapped copy and detects stale files under brand/.
func CheckInstalled(root string) (int, int, error) {
	pin, err := ReadPin(filepath.Join(root, "brand", "source.json"))
	if err != nil {
		return 0, 0, err
	}
	kit, err := Open(pin, filepath.Join(root, "brand", pin.Package+".zip"))
	if err != nil {
		return 0, 0, err
	}
	consumers, err := ReadConsumers(filepath.Join(root, "brand", "repository-consumers.json"))
	if err != nil {
		return 0, 0, err
	}
	if err := ValidateConsumers(root, kit, consumers); err != nil {
		return 0, 0, err
	}
	known := map[string]bool{
		"brand/source.json": true, "brand/repository-consumers.json": true, "brand/README.md": true, "brand/REPOSITORY.md": true,
		"brand/" + pin.Package + ".zip": true, "brand/platform/linux/go-schedule.desktop": true, "brand/platform/linux/go-schedule-indicator.desktop": true,
	}
	var count int
	for _, mapping := range consumers.Mappings {
		want, err := Render(kit, mapping)
		if err != nil {
			return 0, 0, err
		}
		for _, target := range mapping.Targets {
			known[target] = true
			actual, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(target)))
			if err != nil {
				return 0, 0, fmt.Errorf("read consumer %q: %w", target, err)
			}
			if !bytes.Equal(actual, want) {
				return 0, 0, fmt.Errorf("consumer %q differs from official source %q", target, mapping.Source)
			}
			count++
		}
	}
	err = filepath.WalkDir(filepath.Join(root, "brand"), func(filename string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		relative, err := filepath.Rel(root, filename)
		if err != nil {
			return err
		}
		relative = filepath.ToSlash(relative)
		if !known[relative] {
			return fmt.Errorf("brand file %q is not a current mapped asset", relative)
		}
		return nil
	})
	if err != nil {
		return 0, 0, err
	}
	return len(kit.Manifest.Files), count, nil
}
