// Command brand-import installs a pinned official brand archive and its declared consumers.
package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/shruggietech/go-schedule/internal/brandrelease"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "usage: go run ./scripts/brand-import <release-archive.zip> <SHA256SUMS> <new-consumer-map.json>")
		os.Exit(2)
	}
	root, err := os.Getwd()
	if err == nil {
		err = run(root, os.Args[1], os.Args[2], os.Args[3])
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "brand-import: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("brand-import: installed pinned archive and synchronized consumers")
}

func run(root, archive, checksums, mapFile string) error {
	pin, err := brandrelease.ReadPin(filepath.Join(root, "brand", "source.json"))
	if err != nil {
		return err
	}
	if err := brandrelease.VerifyChecksums(pin, checksums); err != nil {
		return err
	}
	kit, err := brandrelease.Open(pin, archive)
	if err != nil {
		return err
	}
	newMap, err := brandrelease.ReadConsumers(mapFile)
	if err != nil {
		return err
	}
	if newMap.Version != 2 {
		return fmt.Errorf("new consumer map must be version 2")
	}
	if err := brandrelease.ValidateConsumers(root, kit, newMap); err != nil {
		return err
	}
	oldMap, err := brandrelease.ReadConsumers(filepath.Join(root, "brand", "repository-consumers.json"))
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	newTargets := map[string]bool{}
	type output struct {
		path string
		data []byte
	}
	var outputs []output
	for _, mapping := range newMap.Mappings {
		data, err := brandrelease.Render(kit, mapping)
		if err != nil {
			return err
		}
		for _, target := range mapping.Targets {
			newTargets[target] = true
			outputs = append(outputs, output{path: target, data: data})
		}
	}
	for _, mapping := range oldMap.Mappings {
		for _, target := range mapping.Targets {
			if err := brandrelease.SafePath(target); err != nil || !oldTargetAllowed(target) {
				return fmt.Errorf("old consumer map has unsafe target %q", target)
			}
			if err := brandrelease.CheckNoSymlink(root, target); err != nil {
				return err
			}
		}
	}
	archiveData, err := os.ReadFile(archive)
	if err != nil {
		return fmt.Errorf("reread brand archive: %w", err)
	}
	mapData, err := os.ReadFile(mapFile)
	if err != nil {
		return fmt.Errorf("reread new consumer map: %w", err)
	}
	if err := write(root, "brand/"+pin.Package+".zip", archiveData); err != nil {
		return err
	}
	for _, item := range outputs {
		if err := write(root, item.path, item.data); err != nil {
			return err
		}
	}
	if err := write(root, "brand/repository-consumers.json", mapData); err != nil {
		return err
	}
	for _, mapping := range oldMap.Mappings {
		for _, target := range mapping.Targets {
			if !newTargets[target] {
				if err := os.Remove(filepath.Join(root, filepath.FromSlash(target))); err != nil && !os.IsNotExist(err) {
					return fmt.Errorf("remove retired consumer %q: %w", target, err)
				}
			}
		}
	}
	allowedBrand := map[string]bool{
		"brand/source.json": true, "brand/repository-consumers.json": true, "brand/README.md": true, "brand/REPOSITORY.md": true,
		"brand/" + pin.Package + ".zip": true, "brand/platform/linux/go-schedule.desktop": true, "brand/platform/linux/go-schedule-indicator.desktop": true,
	}
	for target := range newTargets {
		if strings.HasPrefix(target, "brand/") {
			allowedBrand[target] = true
		}
	}
	if err := filepath.WalkDir(filepath.Join(root, "brand"), func(filename string, entry fs.DirEntry, walkErr error) error {
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
		if allowedBrand[relative] {
			return nil
		}
		if err := os.Remove(filename); err != nil {
			return fmt.Errorf("remove retired brand file %q: %w", relative, err)
		}
		return nil
	}); err != nil {
		return err
	}
	if _, _, err := brandrelease.CheckInstalled(root); err != nil {
		return fmt.Errorf("verify installed brand: %w", err)
	}
	return nil
}

func oldTargetAllowed(name string) bool {
	return strings.HasPrefix(name, "brand/") || strings.HasPrefix(name, "docs/assets/brand/") || strings.HasPrefix(name, "docs/assets/favicons/") || strings.HasPrefix(name, "docs/assets/fonts/") || strings.HasPrefix(name, "desktop/build/") || strings.HasPrefix(name, "desktop/frontend/public/")
}

func write(root, relative string, data []byte) error {
	if err := brandrelease.CheckNoSymlink(root, relative); err != nil {
		return err
	}
	filename := filepath.Join(root, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(filename), 0o755); err != nil {
		return fmt.Errorf("create parent of %q: %w", relative, err)
	}
	if err := os.WriteFile(filename, data, 0o644); err != nil {
		return fmt.Errorf("write %q: %w", relative, err)
	}
	return nil
}
