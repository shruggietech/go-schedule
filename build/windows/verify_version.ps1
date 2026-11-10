<#
.SYNOPSIS
Checks Windows desktop PE and optional MSI version metadata without executing the product.
#>
param(
  [Parameter(Mandatory)][string]$ExePath,
  [Parameter(Mandatory)][string]$Version,
  [string]$MsiPath
)
$ErrorActionPreference = 'Stop'
$numeric = $Version -replace '^v', ''
if ($numeric -notmatch '^\d+\.\d+\.\d+$') { throw 'Expected a selected numeric release version.' }
$expected = [version]"$numeric.0"
$info = [Diagnostics.FileVersionInfo]::GetVersionInfo((Resolve-Path -LiteralPath $ExePath).Path)
$product = [version]::new($info.ProductMajorPart, $info.ProductMinorPart, $info.ProductBuildPart, $info.ProductPrivatePart)
$actual = [version]::new($info.FileMajorPart, $info.FileMinorPart, $info.FileBuildPart, $info.FilePrivatePart)
if ($actual -ne $expected -or $product -ne $expected -or $info.FileVersion -ne $numeric -or $info.ProductVersion -ne $numeric) {
  throw "Desktop resource mismatch: expected $numeric, fixed $actual, product $($info.ProductVersion)"
}
if ($MsiPath) {
  $installer = New-Object -ComObject WindowsInstaller.Installer
  $db = $installer.OpenDatabase((Resolve-Path -LiteralPath $MsiPath).Path, 0)
  foreach ($query in @(
    'SELECT `Value` FROM `Property` WHERE `Property`=''ProductVersion''',
    'SELECT `Version` FROM `File` WHERE `File`=''gosched_gui.exe'''
  )) {
    $view = $db.OpenView($query)
    $view.Execute()
    $row = $view.Fetch()
    if (-not $row) { throw "Missing MSI version row: $query" }
    $value = $row.StringData(1)
    $parts = [version]$value
    $normalized = [version]::new($parts.Major, $parts.Minor, $parts.Build, [Math]::Max(0, $parts.Revision))
    if ($normalized -ne $expected) { throw "MSI version mismatch: expected $numeric, found $value" }
    $view.Close()
  }
}
Write-Output "Windows desktop identity matches $Version (metadata inspection only)."
