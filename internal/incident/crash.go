package incident

import (
	"path/filepath"
	"strings"
)

// Environment replaces only explicitly named child variables, case-insensitively.
func Environment(base []string, values map[string]string) []string {
	result := make([]string, 0, len(base)+len(values))
	for _, entry := range base {
		key, _, ok := strings.Cut(entry, "=")
		replaced := false
		if ok {
			for name := range values {
				if strings.EqualFold(key, name) {
					replaced = true
					break
				}
			}
		}
		if !replaced {
			result = append(result, entry)
		}
	}
	for key, value := range values {
		result = append(result, key+"="+value)
	}
	return result
}

// CoreCLREnvironment enables only this owned runtime's Mini dump collection.
// Retained size is checked separately; this is not an operating-system disk quota.
func (c *Capture) CoreCLREnvironment(base []string) ([]string, error) {
	if _, err := c.path("runtime.json"); err != nil {
		return base, err
	}
	return Environment(base, map[string]string{
		"DOTNET_DbgEnableMiniDump":     "1",
		"DOTNET_DbgMiniDumpType":       "1",
		"DOTNET_DbgMiniDumpName":       filepath.Join(strings.ReplaceAll(c.Dir, "%", "%%"), "crash.%p.dmp"),
		"DOTNET_CreateDumpDiagnostics": "1",
		"DOTNET_CreateDumpLogToFile":   filepath.Join(c.Dir, "createdump.log"),
		"GOSCHEDULE_INCIDENT_IDENTITY": filepath.Join(c.Dir, "runtime.json"),
	}), nil
}
