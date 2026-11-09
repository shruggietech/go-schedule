# Research: Windows version identity

## Decision

Use a small maintained Go build command to temporarily supply Wails info.productVersion from the selected version and restore exact configuration bytes on success or failure. Wails v2.17.1 cached source internal/project/project.go defaults this field to 1.0.0; flags/build.go exposes no product-version flag and packager.go renders the Windows resource template from the project info. Go linker flags alone cannot change it. The existing CI MSI selects 1.0.1 separately and must use the same candidate selection for GUI and installer.

## Alternatives

A static config pin duplicates release identity and will drift. Editing published v2.0.2 changes accepted immutable bytes. Copying the entire desktop into an alternate build workspace is disproportionate. A deferred exact-byte restore is bounded to one build process; simultaneous builds in one checkout are refused by an exclusive lock.

## Development disposition

No selected release means runtime 0.0.1-dev, numeric resources 0.0.1, and an explicit Development build comment. Release selection accepts canonical vMAJOR.MINOR.PATCH or MAJOR.MINOR.PATCH within Windows Installer numeric limits. Other inputs fail before changing config. Child console processes use existing platform.HideConsole and noninteractive input.

## Authority

Owner kickoff S141, issue #458, accepted main f76ff50ba31810418cbb16587fabe39e2445f36a. Published tags and assets remain immutable. No native runtime or upgrade claim follows metadata inspection.

## Source packaging discovery

The prior Windows resource template existed only in an ignored generated directory. S141 explicitly tracks the maintained info.json through a narrow ignore exception. A build with the local generated icon/manifest moved aside recreated those inputs and passed selected PE inspection. Other generated Windows assets remain ignored. Git and credential-helper interactivity are disabled in the managed child environment, alongside CI/npm noninteractive settings.
