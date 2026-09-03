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
	canonicalIconSource = "brand/platform/windows/go-schedule.ico"
	adminGroupID        = "GoScheduleAdminGroup"
	adminGroupName      = "goschedadmin"
)

type xmlElement struct {
	name   string
	attrs  map[string]string
	parent int
}

func parseInstallerElements(data []byte) ([]xmlElement, error) {
	decoder := xml.NewDecoder(strings.NewReader(string(data)))
	var elements []xmlElement
	var stack []int
	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				return elements, nil
			}
			return nil, fmt.Errorf("decode WiX XML: %w", err)
		}
		switch node := token.(type) {
		case xml.StartElement:
			attrs := make(map[string]string, len(node.Attr))
			for _, attr := range node.Attr {
				attrs[attr.Name.Local] = attr.Value
			}
			parent := -1
			if len(stack) != 0 {
				parent = stack[len(stack)-1]
			}
			elements = append(elements, xmlElement{name: node.Name.Local, attrs: attrs, parent: parent})
			stack = append(stack, len(elements)-1)
		case xml.EndElement:
			if len(stack) != 0 {
				stack = stack[:len(stack)-1]
			}
		}
	}
}

func findInstallerElementIndex(elements []xmlElement, name, id string) int {
	for index, element := range elements {
		if element.name == name && element.attrs["Id"] == id {
			return index
		}
	}
	return -1
}

func findInstallerElement(elements []xmlElement, name, id string) (xmlElement, bool) {
	index := findInstallerElementIndex(elements, name, id)
	if index >= 0 {
		return elements[index], true
	}
	return xmlElement{}, false
}

func installerElementHasAncestor(elements []xmlElement, index int, name, id string) bool {
	for parent := elements[index].parent; parent >= 0; parent = elements[parent].parent {
		if elements[parent].name == name && elements[parent].attrs["Id"] == id {
			return true
		}
	}
	return false
}

func findInstallerElementWithin(elements []xmlElement, ancestor, ancestorID, name, id string) (xmlElement, bool) {
	for index, element := range elements {
		if element.name == name && element.attrs["Id"] == id &&
			installerElementHasAncestor(elements, index, ancestor, ancestorID) {
			return element, true
		}
	}
	return xmlElement{}, false
}

func findInstallerPublish(elements []xmlElement, dialogID string, want map[string]string) (xmlElement, bool) {
	for index, element := range elements {
		if element.name != "Publish" || !installerElementHasAncestor(elements, index, "Dialog", dialogID) {
			continue
		}
		matches := true
		for name, value := range want {
			if element.attrs[name] != value {
				matches = false
				break
			}
		}
		if matches {
			return element, true
		}
	}
	return xmlElement{}, false
}

func countInstallerElements(elements []xmlElement, name string, want map[string]string) int {
	count := 0
	for _, element := range elements {
		if element.name != name {
			continue
		}
		matches := true
		for attr, value := range want {
			if element.attrs[attr] != value {
				matches = false
				break
			}
		}
		if matches {
			count++
		}
	}
	return count
}

const (
	completionActionGuard = `NOT Installed AND NOT WIX_UPGRADE_DETECTED AND NOT PATCH AND NOT RESUME AND ACTION = "INSTALL" AND UILevel = 5`
	wipeActionCondition   = `Installed AND REMOVE~="ALL" AND NOT UPGRADINGPRODUCTCODE AND NOT REINSTALL AND GOSCHEDULE_REMOVE_DATA="1"`
)

func validateShortcutFeatureAuthoring(data []byte) []string {
	elements, err := parseInstallerElements(data)
	if err != nil {
		return []string{err.Error()}
	}
	var failures []string
	mainIndex := findInstallerElementIndex(elements, "Feature", "Main")
	if mainIndex < 0 {
		return []string{"Main feature is missing"}
	}
	for _, want := range []struct {
		id, level, component string
	}{
		{"StartMenuShortcut", "1", "AppShortcut"},
		{"DesktopShortcut", "2", "DesktopShortcutComponent"},
	} {
		featureIndex := findInstallerElementIndex(elements, "Feature", want.id)
		if featureIndex < 0 {
			failures = append(failures, want.id+" feature is missing")
			continue
		}
		feature := elements[featureIndex]
		if feature.parent != mainIndex {
			failures = append(failures, want.id+" must be a direct Main child feature")
		}
		if feature.attrs["Level"] != want.level {
			failures = append(failures, want.id+" has the wrong default level")
		}
		if feature.attrs["AllowAdvertise"] != "no" || feature.attrs["InstallDefault"] != "local" {
			failures = append(failures, want.id+" must be a local non-advertised feature")
		}
		ref, ok := findInstallerElementWithin(elements, "Feature", want.id, "ComponentRef", want.component)
		if !ok || ref.parent != featureIndex {
			failures = append(failures, want.id+" has the wrong component ownership")
		}
	}

	if elements[mainIndex].attrs["AllowAbsent"] != "no" || elements[mainIndex].attrs["Display"] != "expand" {
		failures = append(failures, "Main must remain required and initially expanded")
	}
	for _, element := range elements {
		if element.name == "ComponentRef" && element.attrs["Id"] == "AppShortcut" && element.parent == mainIndex {
			failures = append(failures, "AppShortcut must not remain directly owned by Main")
		}
	}

	for _, want := range []struct {
		component, directory, shortcut string
	}{
		{"AppShortcut", "AppMenuFolder", "GuiShortcut"},
		{"DesktopShortcutComponent", "DesktopFolder", "DesktopShortcut"},
	} {
		component, ok := findInstallerElement(elements, "Component", want.component)
		if !ok {
			failures = append(failures, want.component+" component is missing")
			continue
		}
		if component.attrs["Directory"] != want.directory {
			failures = append(failures, want.component+" has the wrong directory")
		}
		shortcut, ok := findInstallerElementWithin(elements, "Component", want.component, "Shortcut", want.shortcut)
		if !ok {
			failures = append(failures, want.shortcut+" shortcut is missing")
			continue
		}
		for attr, value := range map[string]string{
			"Name": "go-schedule", "Description": "Open the go-schedule desktop app",
			"Target": "[INSTALLFOLDER]gosched-gui.exe", "Icon": canonicalIconID,
			"WorkingDirectory": "INSTALLFOLDER",
		} {
			if shortcut.attrs[attr] != value {
				failures = append(failures, fmt.Sprintf("%s %s must be %q", want.shortcut, attr, value))
			}
		}
	}
	if _, ok := findInstallerElement(elements, "StandardDirectory", "DesktopFolder"); !ok {
		failures = append(failures, "DesktopFolder standard directory is missing")
	}
	return failures
}

func validateOwnedInstallerUI(data []byte) []string {
	elements, err := parseInstallerElements(data)
	if err != nil {
		return []string{err.Error()}
	}
	var failures []string
	removeARP, ok := findInstallerElement(elements, "Property", "ARPNOREMOVE")
	if !ok || removeARP.attrs["Value"] != "1" {
		failures = append(failures, "ARPNOREMOVE must disable the reduced-interface direct removal entry")
	}
	if _, ok := findInstallerElement(elements, "Property", "ARPNOMODIFY"); ok {
		failures = append(failures, "ARPNOMODIFY must remain absent so maintenance opens the guided removal flow")
	}
	modifyPath, ok := findInstallerElement(elements, "RegistryValue", "ApplicationManagementModifyPath")
	if !ok || modifyPath.attrs["Root"] != "HKLM" ||
		modifyPath.attrs["Key"] != `Software\Microsoft\Windows\CurrentVersion\Uninstall\[ProductCode]` ||
		modifyPath.attrs["Name"] != "ModifyPath" || modifyPath.attrs["Type"] != "expandable" ||
		modifyPath.attrs["Value"] != `MsiExec.exe /I[ProductCode]` || modifyPath.attrs["KeyPath"] != "yes" {
		failures = append(failures, "ModifyPath must be an owned expandable maintenance command for the current ProductCode")
	}
	if _, ok := findInstallerElementWithin(elements, "Feature", "Main", "ComponentRef", "ApplicationManagementRegistration"); !ok {
		failures = append(failures, "Main feature must install the application-management registration component")
	}
	for _, element := range elements {
		if element.name == "WixUI" {
			failures = append(failures, "stock WixUI must not be composed with the package-owned UI")
			break
		}
	}
	for _, dialog := range []string{"GoScheduleMaintenanceTypeDlg", "GoScheduleUninstallDlg", "GoScheduleWipeConfirmDlg", "GoScheduleExitDlg"} {
		if _, ok := findInstallerElement(elements, "Dialog", dialog); !ok {
			failures = append(failures, dialog+" is missing")
		}
	}
	maintenanceRemove, ok := findInstallerElementWithin(elements, "Dialog", "GoScheduleMaintenanceTypeDlg", "Control", "RemoveButton")
	if !ok || strings.Contains(maintenanceRemove.attrs["DisableCondition"], "ARPNOREMOVE") {
		failures = append(failures, "package-owned maintenance Remove must remain enabled when direct ARP removal is suppressed")
	}
	for _, want := range []struct {
		id, value string
	}{
		{"LAUNCH_GOSCHEDULE", "1"},
		{"GOSCHEDULE_REMOVE_DATA", "0"},
	} {
		property, ok := findInstallerElement(elements, "Property", want.id)
		if !ok || property.attrs["Value"] != want.value {
			failures = append(failures, want.id+" has the wrong default")
		}
	}
	removeProperty, ok := findInstallerElement(elements, "Property", "GOSCHEDULE_REMOVE_DATA")
	if !ok || removeProperty.attrs["Secure"] != "yes" {
		failures = append(failures, "GOSCHEDULE_REMOVE_DATA must be secure")
	}

	for _, want := range []struct {
		id, property string
	}{
		{"LaunchGuiCheckBox", "LAUNCH_GOSCHEDULE"},
		{"OpenDocsCheckBox", "OPEN_GOSCHEDULE_DOCS"},
	} {
		control, ok := findInstallerElementWithin(elements, "Dialog", "GoScheduleExitDlg", "Control", want.id)
		if !ok || control.attrs["Type"] != "CheckBox" || control.attrs["Property"] != want.property || control.attrs["CheckBoxValue"] != "1" {
			failures = append(failures, want.id+" is not an independent completion checkbox")
		}
	}

	for _, want := range []struct {
		property, value, order, selection string
	}{
		{"WixUnelevatedShellExecTarget", "[#gosched_gui.exe]", "1", `LAUNCH_GOSCHEDULE = "1"`},
		{"WixUnelevatedShellExecTarget", "https://shruggietech.github.io/go-schedule/", "3", `OPEN_GOSCHEDULE_DOCS = "1"`},
	} {
		publish, ok := findInstallerPublish(elements, "GoScheduleExitDlg", map[string]string{
			"Property": want.property, "Value": want.value, "Order": want.order,
		})
		if !ok || publish.attrs["Condition"] != completionActionGuard+" AND "+want.selection {
			failures = append(failures, "completion target publish is missing or insufficiently guarded")
		}
	}
	for _, want := range []struct {
		action, order, selection string
	}{
		{"LaunchGui", "2", `LAUNCH_GOSCHEDULE = "1"`}, {"OpenDocs", "4", `OPEN_GOSCHEDULE_DOCS = "1"`},
	} {
		publish, ok := findInstallerPublish(elements, "GoScheduleExitDlg", map[string]string{
			"Event": "DoAction", "Value": want.action, "Order": want.order,
		})
		if !ok || publish.attrs["Condition"] != completionActionGuard+" AND "+want.selection {
			failures = append(failures, want.action+" publish is missing or insufficiently guarded")
		}
		action, ok := findInstallerElement(elements, "CustomAction", want.action)
		if !ok || action.attrs["BinaryRef"] != "Wix4UtilCA_$(sys.BUILDARCHSHORT)" || action.attrs["DllEntry"] != "WixUnelevatedShellExec" || action.attrs["Execute"] != "immediate" || action.attrs["Return"] != "ignore" {
			failures = append(failures, want.action+" must use ignored immediate WixUnelevatedShellExec")
		}
	}

	for _, sequence := range []string{"InstallUISequence", "AdminUISequence"} {
		successRows := 0
		for index, element := range elements {
			if element.name == "Show" && element.attrs["Dialog"] == "GoScheduleExitDlg" && element.attrs["OnExit"] == "success" &&
				installerElementHasAncestor(elements, index, sequence, "") {
				successRows++
			}
		}
		if successRows != 1 {
			failures = append(failures, sequence+" must schedule exactly one package-owned success dialog")
		}
	}
	if _, ok := findInstallerPublish(elements, "GoScheduleMaintenanceTypeDlg", map[string]string{
		"Property": "WixUI_InstallMode", "Value": "Remove", "Order": "1",
	}); !ok {
		failures = append(failures, "maintenance Remove must set the MSI maintenance mode first")
	}
	for _, want := range []map[string]string{
		{"Dialog": "GoScheduleMaintenanceTypeDlg", "Control": "RemoveButton", "Property": "GOSCHEDULE_REMOVE_DATA", "Value": "0", "Order": "2"},
		{"Dialog": "GoScheduleMaintenanceTypeDlg", "Control": "RemoveButton", "Property": "GoScheduleRemoveChoice", "Value": "preserve", "Order": "3"},
	} {
		if countInstallerElements(elements, "Publish", want) != 1 {
			failures = append(failures, "maintenance Remove must reset the guided choice to preserve before navigation")
		}
	}
	if countInstallerElements(elements, "Publish", map[string]string{"Dialog": "GoScheduleMaintenanceTypeDlg", "Control": "RemoveButton", "Event": "NewDialog", "Value": "GoScheduleUninstallDlg", "Order": "4"}) != 1 {
		failures = append(failures, "maintenance Remove must route through the uninstall inventory")
	}
	foundDirectRemove := false
	for index, element := range elements {
		if element.name == "Show" && element.attrs["Dialog"] == "GoScheduleUninstallDlg" &&
			installerElementHasAncestor(elements, index, "InstallUISequence", "") &&
			strings.Contains(element.attrs["Condition"], `REMOVE~="ALL"`) &&
			strings.Contains(element.attrs["Condition"], "Preselected") && element.attrs["Before"] == "ProgressDlg" {
			foundDirectRemove = true
		}
	}
	if !foundDirectRemove {
		failures = append(failures, "direct full-UI uninstall interception is missing")
	}
	for _, controlID := range []string{"AlwaysRemovedText", "PreservedDataText", "WipedDataText", "SecurityStateText", "RemoveChoice"} {
		if _, ok := findInstallerElementWithin(elements, "Dialog", "GoScheduleUninstallDlg", "Control", controlID); !ok {
			failures = append(failures, "uninstall inventory control "+controlID+" is missing")
		}
	}
	confirm, ok := findInstallerElementWithin(elements, "Dialog", "GoScheduleWipeConfirmDlg", "Control", "ConfirmWipe")
	if !ok || confirm.attrs["Default"] == "yes" {
		failures = append(failures, "destructive confirmation must be explicit and non-default")
	}
	return failures
}

func validateWipeInstallerAuthoring(data []byte) []string {
	elements, err := parseInstallerElements(data)
	if err != nil {
		return []string{err.Error()}
	}
	var failures []string
	closeGUI, ok := findInstallerElement(elements, "CloseApplication", "CloseRunningGui")
	if !ok || closeGUI.attrs["Target"] != "gosched-gui.exe" || closeGUI.attrs["TerminateProcess"] != "1" {
		failures = append(failures, "running GUI close authoring is missing")
	}
	cleanup, ok := findInstallerElement(elements, "Binary", "GoScheduleCleanup")
	if !ok || cleanup.attrs["SourceFile"] != `$(StageDir)\gosched-cleanup.exe` {
		failures = append(failures, "installer-private cleanup Binary is missing")
	}
	wipe, ok := findInstallerElement(elements, "CustomAction", "WipeApplicationData")
	if !ok || wipe.attrs["BinaryRef"] != "GoScheduleCleanup" || wipe.attrs["ExeCommand"] != "wipe" || wipe.attrs["Execute"] != "commit" || wipe.attrs["Impersonate"] != "no" || wipe.attrs["Return"] != "ignore" {
		failures = append(failures, "WipeApplicationData must be an ignored non-impersonated commit helper action")
	}
	sequenceFound := false
	for index, element := range elements {
		if element.name == "Custom" && element.attrs["Action"] == "WipeApplicationData" &&
			installerElementHasAncestor(elements, index, "InstallExecuteSequence", "") &&
			element.attrs["Before"] == "InstallFinalize" && element.attrs["Condition"] == wipeActionCondition {
			sequenceFound = true
		}
		if element.name == "Custom" && (element.attrs["Action"] == "LaunchGui" || element.attrs["Action"] == "OpenDocs") &&
			installerElementHasAncestor(elements, index, "InstallExecuteSequence", "") {
			failures = append(failures, "completion actions must not be execute-sequence actions")
		}
	}
	if !sequenceFound {
		failures = append(failures, "WipeApplicationData has the wrong scheduling condition")
	}
	invalidGuard := false
	for _, element := range elements {
		if element.name == "Launch" && strings.Contains(element.attrs["Condition"], "GOSCHEDULE_REMOVE_DATA") &&
			strings.Contains(element.attrs["Condition"], `= "0"`) && strings.Contains(element.attrs["Condition"], `= "1"`) {
			invalidGuard = true
		}
	}
	if !invalidGuard {
		failures = append(failures, "invalid GOSCHEDULE_REMOVE_DATA values are not rejected")
	}
	return failures
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
		if group.attrs["Domain"] != "" {
			failures = append(failures, "administrative Group Domain must be empty for elevated local-group creation")
		}
	}

	installerUser, ok := findInstallerElement(elements, "User", "InstallingUser")
	if !ok {
		failures = append(failures, "installing User declaration is missing")
	} else {
		want := map[string]string{
			"Name":              "[LogonUser]",
			"Domain":            "[%USERDOMAIN]",
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
<Icon Id="GoSchedule.ico" SourceFile="brand/platform/windows/go-schedule.ico" />
<Property Id="ARPPRODUCTICON" Value="GoSchedule.ico" />
<Shortcut Id="GuiShortcut" Icon="GoSchedule.ico" />
</Package></Wix>`

	tests := []struct {
		name        string
		old         string
		replacement string
		want        string
	}{
		{"missing icon", `<Icon Id="GoSchedule.ico" SourceFile="brand/platform/windows/go-schedule.ico" />`, "", "canonical Icon declaration is missing"},
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

func TestWindowsInstallerLifecycleAuthoringContract(t *testing.T) {
	wxs := readRepositoryFile(t, "build", "windows", "goschedule.wxs")
	var failures []string
	failures = append(failures, validateShortcutFeatureAuthoring(wxs)...)
	failures = append(failures, validateOwnedInstallerUI(wxs)...)
	failures = append(failures, validateWipeInstallerAuthoring(wxs)...)
	if len(failures) != 0 {
		t.Fatalf("Windows installer lifecycle contract failed:\n - %s", strings.Join(failures, "\n - "))
	}
}

func TestWindowsInstallerCIScriptAcceptsCleanGUIProcessSnapshot(t *testing.T) {
	script := string(readRepositoryFile(t, "test", "windows", "Invoke-InstallerContractCI.ps1"))
	script = strings.ReplaceAll(script, "\r\n", "\n")
	contract := "[AllowEmptyCollection()]\n    [int[]]$GuiProcessIdsBefore"
	if !strings.Contains(script, contract) {
		t.Fatalf("installer CI harness must accept an empty pre-operation GUI process snapshot")
	}
}

func TestWindowsInstallerCIScriptComparesCanonicalShortcutPaths(t *testing.T) {
	script := string(readRepositoryFile(t, "test", "windows", "Invoke-InstallerContractCI.ps1"))
	for _, contract := range []string{
		"function Test-EquivalentWindowsPath",
		"Test-EquivalentWindowsPath -Actual $shortcut.TargetPath -Expected $expectedTarget",
		"Test-EquivalentWindowsPath -Actual $shortcut.WorkingDirectory -Expected $installDirectory",
	} {
		if !strings.Contains(script, contract) {
			t.Fatalf("installer CI harness must compare shortcut paths canonically; missing %q", contract)
		}
	}
}

func TestWindowsInstallerLifecycleContractRejectsMutations(t *testing.T) {
	shortcutFixture := `<Wix><Package>
<StandardDirectory Id="DesktopFolder" />
<Feature Id="Main" Level="1" AllowAbsent="no" Display="expand">
  <Feature Id="StartMenuShortcut" Level="1" AllowAdvertise="no" InstallDefault="local"><ComponentRef Id="AppShortcut" /></Feature>
  <Feature Id="DesktopShortcut" Level="2" AllowAdvertise="no" InstallDefault="local"><ComponentRef Id="DesktopShortcutComponent" /></Feature>
</Feature>
<Component Id="AppShortcut" Directory="AppMenuFolder"><Shortcut Id="GuiShortcut" Name="go-schedule" Description="Open the go-schedule desktop app" Target="[INSTALLFOLDER]gosched-gui.exe" Icon="GoSchedule.ico" WorkingDirectory="INSTALLFOLDER" /></Component>
<Component Id="DesktopShortcutComponent" Directory="DesktopFolder"><Shortcut Id="DesktopShortcut" Name="go-schedule" Description="Open the go-schedule desktop app" Target="[INSTALLFOLDER]gosched-gui.exe" Icon="GoSchedule.ico" WorkingDirectory="INSTALLFOLDER" /></Component>
</Package></Wix>`
	for _, test := range []struct {
		name, old, replacement, want string
	}{
		{"start menu default", `Id="StartMenuShortcut" Level="1"`, `Id="StartMenuShortcut" Level="2"`, "StartMenuShortcut has the wrong default level"},
		{"desktop default", `Id="DesktopShortcut" Level="2"`, `Id="DesktopShortcut" Level="1"`, "DesktopShortcut has the wrong default level"},
		{"desktop component ownership", `<ComponentRef Id="DesktopShortcutComponent" />`, `<ComponentRef Id="AppShortcut" />`, "DesktopShortcut has the wrong component ownership"},
		{"desktop location", `Id="DesktopShortcutComponent" Directory="DesktopFolder"`, `Id="DesktopShortcutComponent" Directory="AppMenuFolder"`, "DesktopShortcutComponent has the wrong directory"},
		{"desktop icon", `Id="DesktopShortcut" Name="go-schedule" Description="Open the go-schedule desktop app" Target="[INSTALLFOLDER]gosched-gui.exe" Icon="GoSchedule.ico"`, `Id="DesktopShortcut" Name="go-schedule" Description="Open the go-schedule desktop app" Target="[INSTALLFOLDER]gosched-gui.exe" Icon="Other.ico"`, "DesktopShortcut Icon"},
	} {
		t.Run(test.name, func(t *testing.T) {
			mutated := strings.Replace(shortcutFixture, test.old, test.replacement, 1)
			failures := strings.Join(validateShortcutFeatureAuthoring([]byte(mutated)), "\n")
			if !strings.Contains(failures, test.want) {
				t.Fatalf("failures %q do not contain %q", failures, test.want)
			}
		})
	}

	uiFixture := `<Wix><Package>
	<Property Id="ARPNOREMOVE" Value="1"/><Property Id="LAUNCH_GOSCHEDULE" Value="1"/><Property Id="GOSCHEDULE_REMOVE_DATA" Value="0" Secure="yes"/>
	<Feature Id="Main"><ComponentRef Id="ApplicationManagementRegistration"/></Feature><Component Id="ApplicationManagementRegistration"><RegistryValue Id="ApplicationManagementModifyPath" Root="HKLM" Key="Software\Microsoft\Windows\CurrentVersion\Uninstall\[ProductCode]" Name="ModifyPath" Type="expandable" Value="MsiExec.exe /I[ProductCode]" KeyPath="yes"/></Component>
<UI>
<Dialog Id="GoScheduleMaintenanceTypeDlg"><Control Id="RemoveButton" Type="PushButton" DisableCondition="BURNMSIREPAIR OR BURNMSIMODIFY"><Publish Property="WixUI_InstallMode" Value="Remove" Order="1"/></Control></Dialog>
<Dialog Id="GoScheduleUninstallDlg"><Control Id="AlwaysRemovedText"/><Control Id="PreservedDataText"/><Control Id="WipedDataText"/><Control Id="SecurityStateText"/><Control Id="RemoveChoice"/></Dialog>
<Dialog Id="GoScheduleWipeConfirmDlg"><Control Id="ConfirmWipe" Type="PushButton" Default="no"/></Dialog>
<Dialog Id="GoScheduleExitDlg"><Control Id="LaunchGuiCheckBox" Type="CheckBox" Property="LAUNCH_GOSCHEDULE" CheckBoxValue="1"/><Control Id="OpenDocsCheckBox" Type="CheckBox" Property="OPEN_GOSCHEDULE_DOCS" CheckBoxValue="1"/><Control Id="Finish"><Publish Property="WixUnelevatedShellExecTarget" Value="[#gosched_gui.exe]" Order="1" Condition='` + completionActionGuard + ` AND LAUNCH_GOSCHEDULE = "1"'/><Publish Event="DoAction" Value="LaunchGui" Order="2" Condition='` + completionActionGuard + ` AND LAUNCH_GOSCHEDULE = "1"'/><Publish Property="WixUnelevatedShellExecTarget" Value="https://shruggietech.github.io/go-schedule/" Order="3" Condition='` + completionActionGuard + ` AND OPEN_GOSCHEDULE_DOCS = "1"'/><Publish Event="DoAction" Value="OpenDocs" Order="4" Condition='` + completionActionGuard + ` AND OPEN_GOSCHEDULE_DOCS = "1"'/></Control></Dialog>
<Publish Dialog="GoScheduleMaintenanceTypeDlg" Control="RemoveButton" Property="GOSCHEDULE_REMOVE_DATA" Value="0" Order="2"/><Publish Dialog="GoScheduleMaintenanceTypeDlg" Control="RemoveButton" Property="GoScheduleRemoveChoice" Value="preserve" Order="3"/><Publish Dialog="GoScheduleMaintenanceTypeDlg" Control="RemoveButton" Event="NewDialog" Value="GoScheduleUninstallDlg" Order="4"/>
<InstallUISequence><Show Dialog="GoScheduleUninstallDlg" Before="ProgressDlg" Condition='Installed AND REMOVE~="ALL" AND Preselected AND UILevel = 5'/><Show Dialog="GoScheduleExitDlg" OnExit="success"/></InstallUISequence><AdminUISequence><Show Dialog="GoScheduleExitDlg" OnExit="success"/></AdminUISequence>
</UI>
<CustomAction Id="LaunchGui" BinaryRef="Wix4UtilCA_$(sys.BUILDARCHSHORT)" DllEntry="WixUnelevatedShellExec" Execute="immediate" Return="ignore"/><CustomAction Id="OpenDocs" BinaryRef="Wix4UtilCA_$(sys.BUILDARCHSHORT)" DllEntry="WixUnelevatedShellExec" Execute="immediate" Return="ignore"/>
</Package></Wix>`
	for _, test := range []struct {
		name, old, replacement, want string
	}{
		{"stock UI composition", `<UI>`, `<WixUI Id="WixUI_FeatureTree"/><UI>`, "stock WixUI"},
		{"direct Settings removal restored", `<Property Id="ARPNOREMOVE" Value="1"/>`, ``, "ARPNOREMOVE"},
		{"maintenance suppressed", `<Property Id="ARPNOREMOVE" Value="1"/>`, `<Property Id="ARPNOREMOVE" Value="1"/><Property Id="ARPNOMODIFY" Value="1"/>`, "ARPNOMODIFY"},
		{"maintenance command removed", `<RegistryValue Id="ApplicationManagementModifyPath" Root="HKLM" Key="Software\Microsoft\Windows\CurrentVersion\Uninstall\[ProductCode]" Name="ModifyPath" Type="expandable" Value="MsiExec.exe /I[ProductCode]" KeyPath="yes"/>`, ``, "ModifyPath"},
		{"maintenance command made destructive", `Value="MsiExec.exe /I[ProductCode]"`, `Value="MsiExec.exe /X[ProductCode]"`, "ModifyPath"},
		{"maintenance component unreferenced", `<ComponentRef Id="ApplicationManagementRegistration"/>`, ``, "Main feature"},
		{"maintenance Remove disabled", `DisableCondition="BURNMSIREPAIR OR BURNMSIMODIFY"`, `DisableCondition="ARPNOREMOVE OR BURNMSIREPAIR OR BURNMSIMODIFY"`, "maintenance Remove"},
		{"launch elevation guard", ` AND UILevel = 5`, ``, "insufficiently guarded"},
		{"second success dialog", `<Show Dialog="GoScheduleExitDlg" OnExit="success"/></InstallUISequence>`, `<Show Dialog="GoScheduleExitDlg" OnExit="success"/><Show Dialog="GoScheduleExitDlg" OnExit="success"/></InstallUISequence>`, "InstallUISequence must schedule exactly one"},
		{"remove routing", `Value="GoScheduleUninstallDlg" Order="4"/>`, `Value="VerifyReadyDlg" Order="4"/>`, "maintenance Remove"},
		{"destructive default", `Id="ConfirmWipe" Type="PushButton" Default="no"`, `Id="ConfirmWipe" Type="PushButton" Default="yes"`, "explicit and non-default"},
	} {
		t.Run(test.name, func(t *testing.T) {
			mutated := strings.Replace(uiFixture, test.old, test.replacement, 1)
			failures := strings.Join(validateOwnedInstallerUI([]byte(mutated)), "\n")
			if !strings.Contains(failures, test.want) {
				t.Fatalf("failures %q do not contain %q", failures, test.want)
			}
		})
	}

	wipeFixture := `<Wix><Package>
<Launch Condition='GOSCHEDULE_REMOVE_DATA = "0" OR GOSCHEDULE_REMOVE_DATA = "1"'/>
<CloseApplication Id="CloseRunningGui" Target="gosched-gui.exe" TerminateProcess="1"/>
<Binary Id="GoScheduleCleanup" SourceFile="$(StageDir)\gosched-cleanup.exe"/>
<CustomAction Id="WipeApplicationData" BinaryRef="GoScheduleCleanup" ExeCommand="wipe" Execute="commit" Impersonate="no" Return="ignore"/>
<InstallExecuteSequence><Custom Action="WipeApplicationData" Before="InstallFinalize" Condition='` + wipeActionCondition + `'/></InstallExecuteSequence>
</Package></Wix>`
	for _, test := range []struct {
		name, old, replacement, want string
	}{
		{"helper source", `gosched-cleanup.exe`, `gosched.exe`, "cleanup Binary"},
		{"commit timing", `Execute="commit"`, `Execute="deferred"`, "ignored non-impersonated commit"},
		{"ignored return", `Return="ignore"`, `Return="check"`, "ignored non-impersonated commit"},
		{"upgrade exclusion", ` AND NOT UPGRADINGPRODUCTCODE`, ``, "wrong scheduling condition"},
		{"invalid mode accepted", ` OR GOSCHEDULE_REMOVE_DATA = "1"`, ``, "invalid GOSCHEDULE_REMOVE_DATA"},
	} {
		t.Run(test.name, func(t *testing.T) {
			mutated := strings.Replace(wipeFixture, test.old, test.replacement, 1)
			failures := strings.Join(validateWipeInstallerAuthoring([]byte(mutated)), "\n")
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

func TestWindowsInstallGuideDirectMemberTokenRefreshContract(t *testing.T) {
	guide := string(readRepositoryFile(t, "docs", "INSTALL-windows.md"))
	for _, fragment := range []string{
		"authorizes that direct user",
		"launch immediately without",
		"signing out or running the desktop application elevated",
		"nested-group membership",
	} {
		if !strings.Contains(guide, fragment) {
			t.Errorf("Windows install guide is missing direct-member guidance %q", fragment)
		}
	}
	if strings.Contains(guide, "Sign out and back in once after the first install") {
		t.Error("Windows install guide still mandates the obsolete direct-member sign-out workaround")
	}
}

func TestWindowsInstallerAdminGroupRejectsBrokenLifecycle(t *testing.T) {
	domainQualifiedGroup := `<Wix><Package><Feature><ComponentRef Id="AdminAccessProvisioning" /></Feature><Component Id="AdminAccessProvisioning">
<Group Id="GoScheduleAdminGroup" Name="goschedadmin" Domain="[ComputerName]" CreateGroup="yes" FailIfExists="no" RemoveOnUninstall="no" UpdateIfExists="yes" Vital="yes" />
<User Id="InstallingUser" Name="[LogonUser]" Domain="[%USERDOMAIN]" CreateUser="no" FailIfExists="no" RemoveOnUninstall="no" UpdateIfExists="yes" Vital="yes"><GroupRef Id="GoScheduleAdminGroup" /></User>
</Component></Package></Wix>`
	domainFailures := strings.Join(validateInstallerAdminGroup([]byte(domainQualifiedGroup)), "\n")
	if !strings.Contains(domainFailures, "administrative Group Domain must be empty") {
		t.Fatalf("domain-qualified local group failures %q do not reject domain routing", domainFailures)
	}

	valid := `<Wix><Package><Feature><ComponentRef Id="AdminAccessProvisioning" /></Feature><Component Id="AdminAccessProvisioning">
<Group Id="GoScheduleAdminGroup" Name="goschedadmin" CreateGroup="yes" FailIfExists="no" RemoveOnUninstall="no" UpdateIfExists="yes" Vital="yes" />
<User Id="InstallingUser" Name="[LogonUser]" Domain="[%USERDOMAIN]" CreateUser="no" FailIfExists="no" RemoveOnUninstall="no" UpdateIfExists="yes" Vital="yes"><GroupRef Id="GoScheduleAdminGroup" /></User>
</Component></Package></Wix>`
	for _, test := range []struct {
		name        string
		old         string
		replacement string
		want        string
	}{
		{name: "missing group", old: `<Group Id="GoScheduleAdminGroup" Name="goschedadmin" CreateGroup="yes" FailIfExists="no" RemoveOnUninstall="no" UpdateIfExists="yes" Vital="yes" />`, want: "administrative Group declaration is missing"},
		{name: "wrong group name", old: `Name="goschedadmin"`, replacement: `Name="other"`, want: "administrative Group Name"},
		{name: "destructive uninstall", old: `RemoveOnUninstall="no"`, replacement: `RemoveOnUninstall="yes"`, want: "administrative Group RemoveOnUninstall"},
		{name: "unqualified user", old: `Domain="[%USERDOMAIN]"`, want: "installing User Domain"},
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
		"-icon=brand/platform/windows/go-schedule.ico",
		`cp brand/platform/macos/go-schedule.icns "$app/Contents/Resources/icon.icns"`,
		`cp brand/platform/linux/go-schedule.desktop "$stage/share/applications/"`,
		`cp -R brand/platform/linux/hicolor "$stage/share/icons/"`,
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
		"[ValidateSet('candidate', 'published', 'local-demo')]",
		"[string]$ArtifactClass",
		"[string]$ArtifactOrigin",
		`"- Evidence class: **$ArtifactClass artifact**"`,
		`"- Artifact origin: $ArtifactOrigin"`,
		`"- $ArtifactClass artifact status: **$status**"`,
		"FROM ``Wix4Group`` WHERE ``Group``='GoScheduleAdminGroup'",
		"Wix4Group.GoScheduleAdminGroup.Domain",
		"expected empty for elevated local-group creation",
		"- Administrative group row:",
		"Property.ARPNOREMOVE",
		"Property.ARPNOMODIFY must remain absent",
		"Registry.ApplicationManagementModifyPath",
		"GoScheduleMaintenanceTypeDlg.RemoveButton",
		"- Application-management registration:",
	} {
		if !strings.Contains(inspector, fragment) {
			t.Errorf("MSI inspector is missing evidence-provenance fragment %q", fragment)
		}
	}
	if strings.Contains(inspector, "Candidate/published artifact status") {
		t.Error("MSI inspector combines candidate and published evidence statuses")
	}

	contractCI := string(readRepositoryFile(t, "test", "windows", "Invoke-InstallerContractCI.ps1"))
	for _, fragment := range []string{
		"function Assert-ApplicationManagementRegistration",
		"NoRemove",
		"NoModify",
		"ModifyPath",
		"UninstallString",
		"application-management-registration",
	} {
		if !strings.Contains(contractCI, fragment) {
			t.Errorf("installer CI harness is missing application-management evidence fragment %q", fragment)
		}
	}

	lifecycle := string(readRepositoryFile(t, "test", "windows", "Invoke-InstallerLifecycle.ps1"))
	for _, fragment := range []string{
		"[ValidateSet('candidate', 'published')]",
		"[ValidateSet('fresh','upgrade','access-probe','installed-core-probe')]",
		"[string]$ArtifactClass",
		"[string]$ArtifactOrigin",
		`"- Evidence class: **$ArtifactClass artifact**"`,
		`"- Artifact origin: $ArtifactOrigin"`,
		"-WindowStyle Hidden",
		"candidate-repair",
		"candidate-reinstall",
		"candidate-upgrade",
		"GroupSid",
		"IntendedMemberPresent",
		"GroupBeforeService",
		"MembershipBeforeService",
		"Windows 11 client required",
		"ExpectedPriorHash",
		"Installed candidate CLI is missing",
		"CandidateProductCode",
		"Installed product identity does not match the candidate MSI",
		"Token contains goschedadmin SID",
		"Expected restricted pipe descriptor",
		"AccessProbeExitCode",
		"AccessProbeOutput",
		"Daemon admin_group is",
		"PSChildName",
		"Invoke-InstalledCoreExecutionProbe",
		"'*/5 * * * * *'",
		"'WindowsPowerShell'",
		"[IO.File]::AppendAllText",
		"Wait-TaskRun -TaskId $successTask.id",
		"-Trigger manual",
		"-Trigger schedule",
		"S038-controlled-failure",
		"missing-s038-executable.exe",
		"process start failed for",
		"ExecutionEnvironmentKeys = '<none>'",
		"Cleanup uninstall failed",
		"Write-Evidence -Status $status -Problems $problems",
		"$OperationArguments[0]",
		"$OperationArguments | Select-Object -Skip 1",
		"if (-not $PSCmdlet.ShouldProcess(",
		"Lifecycle action skipped; stopping without evidence.",
		"if (-not $actionSkipped -and",
	} {
		if !strings.Contains(lifecycle, fragment) {
			t.Errorf("lifecycle verifier is missing evidence-integrity fragment %q", fragment)
		}
	}
	if strings.Contains(lifecycle, "Start-Process -FilePath 'msiexec.exe' -ArgumentList") &&
		!strings.Contains(lifecycle, "-WindowStyle Hidden") {
		t.Error("lifecycle verifier launches msiexec without a hidden-window guarantee")
	}
	if strings.Contains(lifecycle, "$arguments = @($OperationArguments) + @(") {
		t.Error("lifecycle verifier appends the MSI after reinstall properties")
	}
	if got := strings.Count(lifecycle, "if (-not $PSCmdlet.ShouldProcess("); got < 7 {
		t.Errorf("lifecycle verifier gates only %d destructive actions; want at least 7", got)
	}
	finallyPosition := strings.LastIndex(lifecycle, "} finally {")
	writePosition := strings.LastIndex(lifecycle, "Write-Evidence -Status $status")
	if finallyPosition < 0 || writePosition < finallyPosition {
		t.Error("lifecycle verifier writes final evidence before cleanup completes")
	}

	serviceProbe := string(readRepositoryFile(t, "test", "windows", "Invoke-ServiceCoreCI.ps1"))
	for _, fragment := range []string{
		"Probe token unexpectedly contains the newly created group SID",
		"Invoke-ProbeCli -Arguments @('service', 'install')",
		"Wait-ProbeRun -TaskId $successTask.id -Trigger manual",
		"Wait-ProbeRun -TaskId $successTask.id -Trigger schedule",
		"S038-controlled-failure",
		"missing-s038-executable.exe",
		"service_account = $service.StartName",
		"environment_keys = @()",
		"} finally {",
		"Remove-LocalGroup -Name $groupName",
	} {
		if !strings.Contains(serviceProbe, fragment) {
			t.Errorf("service-core verifier is missing boundary fragment %q", fragment)
		}
	}

	ci := string(readRepositoryFile(t, ".github", "workflows", "ci.yml"))
	for _, fragment := range []string{
		"windows-service-core:",
		"runs-on: windows-latest",
		"Invoke-ServiceCoreCI.ps1",
		"windows-service-core-evidence",
	} {
		if !strings.Contains(ci, fragment) {
			t.Errorf("CI workflow is missing service-core fragment %q", fragment)
		}
	}
}
