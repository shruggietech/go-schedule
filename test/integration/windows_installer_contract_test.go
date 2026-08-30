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
	adminGroupID        = "GoScheduleAdminGroup"
	adminGroupName      = "goschedadmin"
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

func validateInstallerAdminGroup(data []byte) []string {
	elements, err := parseInstallerElements(data)
	if err != nil {
		return []string{err.Error()}
	}

	var failures []string
	group, ok := findInstallerElement(elements, "Group", adminGroupID)
	if !ok {
		failures = append(failures, "administrative Group declaration is missing")
	} else {
		want := map[string]string{
			"Name":              adminGroupName,
			"Domain":            "[ComputerName]",
			"CreateGroup":       "yes",
			"FailIfExists":      "no",
			"RemoveOnUninstall": "no",
			"UpdateIfExists":    "yes",
			"Vital":             "yes",
		}
		for attr, value := range want {
			if group.attrs[attr] != value {
				failures = append(failures, fmt.Sprintf("administrative Group %s must be %q", attr, value))
			}
		}
	}

	installerUser, ok := findInstallerElement(elements, "User", "InstallingUser")
	if !ok {
		failures = append(failures, "installing User declaration is missing")
	} else {
		want := map[string]string{
			"Name":              "[LogonUser]",
			"CreateUser":        "no",
			"FailIfExists":      "no",
			"RemoveOnUninstall": "no",
			"UpdateIfExists":    "yes",
			"Vital":             "yes",
		}
		for attr, value := range want {
			if installerUser.attrs[attr] != value {
				failures = append(failures, fmt.Sprintf("installing User %s must be %q", attr, value))
			}
		}
	}

	groupRef, ok := findInstallerElement(elements, "GroupRef", adminGroupID)
	if !ok {
		failures = append(failures, "installing user administrative GroupRef is missing")
	} else if groupRef.attrs["Id"] != adminGroupID {
		failures = append(failures, "installing user GroupRef is incorrect")
	}
	if _, ok := findInstallerElement(elements, "ComponentRef", "AdminAccessProvisioning"); !ok {
		failures = append(failures, "administrative access component is absent from the main feature")
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

func TestWindowsInstallerAdminGroupContract(t *testing.T) {
	wxs := readRepositoryFile(t, "build", "windows", "goschedule.wxs")
	if failures := validateInstallerAdminGroup(wxs); len(failures) != 0 {
		t.Fatalf("Windows installer administrative-group contract failed:\n - %s", strings.Join(failures, "\n - "))
	}

	workflow := string(readRepositoryFile(t, ".github", "workflows", "release.yml"))
	for _, dependency := range []string{
		"dotnet tool install --global wix --version 6.0.2",
		"WixToolset.UI.wixext/6.0.2",
		"WixToolset.Util.wixext/6.0.2",
	} {
		if !strings.Contains(workflow, dependency) {
			t.Errorf("release workflow is missing pinned WiX 6 dependency %q", dependency)
		}
	}
	if strings.Contains(workflow, "5.0.2") {
		t.Error("release workflow retains the obsolete WiX 5.0.2 pin")
	}
}

func TestWindowsInstallerAdminGroupRejectsBrokenLifecycle(t *testing.T) {
	valid := `<Wix><Package><Feature><ComponentRef Id="AdminAccessProvisioning" /></Feature><Component Id="AdminAccessProvisioning">
<Group Id="GoScheduleAdminGroup" Name="goschedadmin" Domain="[ComputerName]" CreateGroup="yes" FailIfExists="no" RemoveOnUninstall="no" UpdateIfExists="yes" Vital="yes" />
<User Id="InstallingUser" Name="[LogonUser]" CreateUser="no" FailIfExists="no" RemoveOnUninstall="no" UpdateIfExists="yes" Vital="yes"><GroupRef Id="GoScheduleAdminGroup" /></User>
</Component></Package></Wix>`
	for _, test := range []struct {
		name        string
		old         string
		replacement string
		want        string
	}{
		{name: "missing group", old: `<Group Id="GoScheduleAdminGroup" Name="goschedadmin" Domain="[ComputerName]" CreateGroup="yes" FailIfExists="no" RemoveOnUninstall="no" UpdateIfExists="yes" Vital="yes" />`, want: "administrative Group declaration is missing"},
		{name: "wrong group name", old: `Name="goschedadmin"`, replacement: `Name="other"`, want: "administrative Group Name"},
		{name: "destructive uninstall", old: `RemoveOnUninstall="no"`, replacement: `RemoveOnUninstall="yes"`, want: "administrative Group RemoveOnUninstall"},
		{name: "creates user", old: `CreateUser="no"`, replacement: `CreateUser="yes"`, want: "installing User CreateUser"},
		{name: "missing membership", old: `<GroupRef Id="GoScheduleAdminGroup" />`, want: "GroupRef is missing"},
		{name: "feature omission", old: `<ComponentRef Id="AdminAccessProvisioning" />`, want: "component is absent"},
	} {
		t.Run(test.name, func(t *testing.T) {
			mutated := strings.Replace(valid, test.old, test.replacement, 1)
			failures := strings.Join(validateInstallerAdminGroup([]byte(mutated)), "\n")
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

func TestWindowsInstallerEvidenceToolingContract(t *testing.T) {
	inspector := string(readRepositoryFile(t, "test", "windows", "inspect-installer.ps1"))
	for _, fragment := range []string{
		"[ValidateSet('candidate', 'published')]",
		"[string]$ArtifactClass",
		"[string]$ArtifactOrigin",
		`"- Evidence class: **$ArtifactClass artifact**"`,
		`"- Artifact origin: $ArtifactOrigin"`,
		`"- $ArtifactClass artifact status: **$status**"`,
	} {
		if !strings.Contains(inspector, fragment) {
			t.Errorf("MSI inspector is missing evidence-provenance fragment %q", fragment)
		}
	}
	if strings.Contains(inspector, "Candidate/published artifact status") {
		t.Error("MSI inspector combines candidate and published evidence statuses")
	}

	lifecycle := string(readRepositoryFile(t, "test", "windows", "install-lifecycle.ps1"))
	for _, fragment := range []string{
		"[ValidateSet('candidate', 'published')]",
		"[string]$ArtifactClass",
		"[string]$ArtifactOrigin",
		`"- Evidence class: **$ArtifactClass artifact**"`,
		`"- Artifact origin: $ArtifactOrigin"`,
		"function Read-NativeObservation",
		`$nativeObservations[$surface] = Read-NativeObservation -Surface $surface`,
		"Cleanup uninstall failed; the package may remain installed",
		"Final machine state: product registered=",
		"Write-Evidence -Status $lifecycleStatus -Problems $failures",
	} {
		if !strings.Contains(lifecycle, fragment) {
			t.Errorf("lifecycle verifier is missing evidence-integrity fragment %q", fragment)
		}
	}
	if strings.Contains(lifecycle, "unavailable until recorded by the operator") {
		t.Error("lifecycle verifier still hard-codes unavailable native observations")
	}
	finallyPosition := strings.LastIndex(lifecycle, "} finally {")
	writePosition := strings.LastIndex(lifecycle, "Write-Evidence -Status $lifecycleStatus")
	if finallyPosition < 0 || writePosition < finallyPosition {
		t.Error("lifecycle verifier writes final evidence before cleanup completes")
	}
}
