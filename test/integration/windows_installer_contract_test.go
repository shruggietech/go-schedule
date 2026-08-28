package integration

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	canonicalIconID     = "GoSchedule.ico"
	canonicalIconSource = "cmd/gosched-gui/icon.ico"
)

type xmlElement struct {
	name  string
	attrs map[string]string
}

func parseInstallerElements(data []byte) ([]xmlElement, error) {
	decoder := xml.NewDecoder(strings.NewReader(string(data)))
	var elements []xmlElement
	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				return elements, nil
			}
			return nil, fmt.Errorf("decode WiX XML: %w", err)
		}
		start, ok := token.(xml.StartElement)
		if !ok {
			continue
		}
		attrs := make(map[string]string, len(start.Attr))
		for _, attr := range start.Attr {
			attrs[attr.Name.Local] = attr.Value
		}
		elements = append(elements, xmlElement{name: start.Name.Local, attrs: attrs})
	}
}

func findInstallerElement(elements []xmlElement, name, id string) (xmlElement, bool) {
	for _, element := range elements {
		if element.name == name && element.attrs["Id"] == id {
			return element, true
		}
	}
	return xmlElement{}, false
}

func validateInstallerIdentity(data []byte) []string {
	elements, err := parseInstallerElements(data)
	if err != nil {
		return []string{err.Error()}
	}

	var failures []string
	icon, ok := findInstallerElement(elements, "Icon", canonicalIconID)
	if !ok {
		failures = append(failures, "canonical Icon declaration is missing")
	} else if icon.attrs["SourceFile"] != canonicalIconSource {
		failures = append(failures, "canonical Icon SourceFile is incorrect")
	}

	arp, ok := findInstallerElement(elements, "Property", "ARPPRODUCTICON")
	if !ok {
		failures = append(failures, "ARPPRODUCTICON property is missing")
	} else if arp.attrs["Value"] != canonicalIconID {
		failures = append(failures, "ARPPRODUCTICON does not reference the canonical Icon")
	}

	shortcut, ok := findInstallerElement(elements, "Shortcut", "GuiShortcut")
	if !ok {
		failures = append(failures, "GuiShortcut is missing")
	} else if shortcut.attrs["Icon"] != canonicalIconID {
		failures = append(failures, "GuiShortcut does not reference the canonical Icon")
	}

	return failures
}

func readRepositoryFile(t *testing.T, parts ...string) []byte {
	t.Helper()
	pathParts := append([]string{"..", ".."}, parts...)
	data, err := os.ReadFile(filepath.Join(pathParts...))
	if err != nil {
		t.Fatalf("read repository file: %v", err)
	}
	return data
}

func TestWindowsInstallerIdentityContract(t *testing.T) {
	wxs := readRepositoryFile(t, "build", "windows", "goschedule.wxs")
	if failures := validateInstallerIdentity(wxs); len(failures) != 0 {
		t.Fatalf("Windows installer identity contract failed:\n - %s", strings.Join(failures, "\n - "))
	}
}

func TestWindowsInstallerIdentityRejectsBrokenRelationships(t *testing.T) {
	valid := `<Wix><Package>
<Icon Id="GoSchedule.ico" SourceFile="cmd/gosched-gui/icon.ico" />
<Property Id="ARPPRODUCTICON" Value="GoSchedule.ico" />
<Shortcut Id="GuiShortcut" Icon="GoSchedule.ico" />
</Package></Wix>`

	tests := []struct {
		name        string
		old         string
		replacement string
		want        string
	}{
		{"missing icon", `<Icon Id="GoSchedule.ico" SourceFile="cmd/gosched-gui/icon.ico" />`, "", "canonical Icon declaration is missing"},
		{"wrong icon source", canonicalIconSource, "gui/assets/icon.ico", "canonical Icon SourceFile is incorrect"},
		{"missing ARP property", `<Property Id="ARPPRODUCTICON" Value="GoSchedule.ico" />`, "", "ARPPRODUCTICON property is missing"},
		{"wrong ARP reference", `<Property Id="ARPPRODUCTICON" Value="GoSchedule.ico" />`, `<Property Id="ARPPRODUCTICON" Value="Other.ico" />`, "ARPPRODUCTICON does not reference the canonical Icon"},
		{"missing shortcut", `<Shortcut Id="GuiShortcut" Icon="GoSchedule.ico" />`, "", "GuiShortcut is missing"},
		{"wrong shortcut reference", `<Shortcut Id="GuiShortcut" Icon="GoSchedule.ico" />`, `<Shortcut Id="GuiShortcut" Icon="Other.ico" />`, "GuiShortcut does not reference the canonical Icon"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mutated := strings.Replace(valid, test.old, test.replacement, 1)
			failures := strings.Join(validateInstallerIdentity([]byte(mutated)), "\n")
			if !strings.Contains(failures, test.want) {
				t.Fatalf("failures %q do not contain %q", failures, test.want)
			}
		})
	}
}

func TestWindowsInstallerGUIResourceContract(t *testing.T) {
	workflow := string(readRepositoryFile(t, ".github", "workflows", "release.yml"))
	required := []string{
		"go run github.com/josephspurrier/goversioninfo/cmd/goversioninfo@v1.4.1",
		"-64",
		"-icon=cmd/gosched-gui/icon.ico",
		"-o=cmd/gosched-gui/resource_windows_amd64.syso",
		"./cmd/gosched-gui",
	}
	for _, fragment := range required {
		if !strings.Contains(workflow, fragment) {
			t.Errorf("release workflow is missing GUI resource contract fragment %q", fragment)
		}
	}
}
