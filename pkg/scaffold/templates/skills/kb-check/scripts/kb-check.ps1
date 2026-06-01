param(
  [switch]$List,
  [switch]$All
)

$ErrorActionPreference = "Stop"

function Add-Check {
  param([string]$Name, [string]$Command, [string]$Reason)
  [PSCustomObject]@{ Name = $Name; Command = $Command; Reason = $Reason }
}

function Get-PowerShellCommand {
  $pwsh = Get-Command pwsh -ErrorAction SilentlyContinue | Select-Object -First 1
  if ($pwsh) {
    return "pwsh -NoProfile -File"
  }
  $windowsPowerShell = Get-Command powershell -ErrorAction SilentlyContinue | Select-Object -First 1
  if ($windowsPowerShell) {
    return "powershell -NoProfile -ExecutionPolicy Bypass -File"
  }
  throw "Neither pwsh nor powershell is available."
}

$psFile = Get-PowerShellCommand
$checks = New-Object System.Collections.Generic.List[object]

if (Test-Path "package.json") {
  $pkg = Get-Content "package.json" -Raw | ConvertFrom-Json
  $scripts = $pkg.scripts
  foreach ($name in @("lint", "typecheck", "test", "test:unit", "test:integration", "test:e2e", "build")) {
    if ($scripts -and $scripts.PSObject.Properties.Name -contains $name) {
      $runner = if (Test-Path "pnpm-lock.yaml") { "pnpm" } elseif (Test-Path "yarn.lock") { "yarn" } else { "npm run" }
      $cmd = if ($runner -eq "npm run") { "npm run $name" } else { "$runner $name" }
      $checks.Add((Add-Check $name $cmd "package.json script"))
    }
  }
}

if ((Test-Path "pyproject.toml") -or (Test-Path "pytest.ini")) {
  $checks.Add((Add-Check "pytest" "python -m pytest" "Python test config detected"))
}

if (Test-Path "go.mod") {
  $checks.Add((Add-Check "go-test" "go test ./..." "Go module detected"))
}

if (Get-ChildItem -Filter "*.sln" -ErrorAction SilentlyContinue) {
  $sln = (Get-ChildItem -Filter "*.sln" | Select-Object -First 1).Name
  $checks.Add((Add-Check "dotnet-test" "dotnet test `"$sln`"" ".NET solution detected"))
  $checks.Add((Add-Check "dotnet-build" "dotnet build `"$sln`" --no-restore" ".NET solution detected"))
} elseif (Get-ChildItem -Filter "*.csproj" -Recurse -ErrorAction SilentlyContinue | Select-Object -First 1) {
  $checks.Add((Add-Check "dotnet-test" "dotnet test" ".NET project detected"))
}

if (Test-Path "Makefile") {
  $checks.Add((Add-Check "make-test" "make test" "Makefile detected"))
}

if ((Test-Path ".github/skills") -and (Test-Path "config/skill-quality.json")) {
  if (Test-Path "scripts/skill-lint.ps1") {
    $checks.Add((Add-Check "skill-lint" "$psFile scripts\skill-lint.ps1" "skill quality config detected"))
  }
  if (Test-Path "scripts/route-complexity-eval.ps1") {
    $checks.Add((Add-Check "route-complexity-eval" "$psFile scripts\route-complexity-eval.ps1" "route complexity eval fixtures detected"))
  }
  if (Test-Path "scripts/skill-eval.ps1") {
    $checks.Add((Add-Check "skill-eval" "$psFile scripts\skill-eval.ps1" "skill eval selftest fixtures detected"))
  }
  if (Test-Path "scripts/skill-eval-manifest-selftest.ps1") {
    $checks.Add((Add-Check "skill-eval-manifest-selftest" "$psFile scripts\skill-eval-manifest-selftest.ps1" "skill eval protected-file hash selftest detected"))
  }
  if (Test-Path "scripts/skill-eval-baseline-selftest.ps1") {
    $checks.Add((Add-Check "skill-eval-baseline-selftest" "$psFile scripts\skill-eval-baseline-selftest.ps1" "skill eval baseline regression selftest detected"))
  }
  if (Test-Path "scripts/skill-eval-run-codex.ps1") {
    $checks.Add((Add-Check "skill-eval-codex-dry-run" "$psFile scripts\skill-eval-run-codex.ps1 -FixtureId tiny-typo-fix -DryRun" "Codex skill eval adapter detected"))
  }
  if (Test-Path "scripts/skill-eval-run-ghcp.ps1") {
    $checks.Add((Add-Check "skill-eval-ghcp-dry-run" "$psFile scripts\skill-eval-run-ghcp.ps1 -FixtureId tiny-typo-fix -DryRun" "GHCP skill eval adapter detected"))
  }
  if ((Test-Path "scripts/skill-eval-wrap.ps1") -and (Test-Path "scripts/skill-eval-run-ghcp.ps1")) {
    $checks.Add((Add-Check "skill-eval-observed-trace-dry-run" "$psFile scripts\skill-eval-wrap.ps1 -Runner scripts\skill-eval-run-ghcp.ps1 -FixtureId tiny-typo-fix -DryRun -Sealed" "observed trace eval wrapper detected"))
  }
  if (Test-Path "scripts/skill-eval-quality.ps1") {
    $checks.Add((Add-Check "skill-eval-quality" "$psFile scripts\skill-eval-quality.ps1" "skill output quality rubric fixtures detected"))
  }
  if (Test-Path "scripts/kb-pipeline-selftest.ps1") {
    $checks.Add((Add-Check "kb-pipeline-selftest" "$psFile scripts\kb-pipeline-selftest.ps1" "KB coded pipeline spike selftest detected"))
  }
  if (Test-Path "scripts/skill-surface-report.ps1") {
    $checks.Add((Add-Check "skill-surface-report" "$psFile scripts\skill-surface-report.ps1" "skill loaded-surface report detected"))
  }
  if (Test-Path "scripts/skill-marketplace-firebreak.ps1") {
    $checks.Add((Add-Check "skill-marketplace-firebreak" "$psFile scripts\skill-marketplace-firebreak.ps1" "private marketplace quarantine firebreak detected"))
  }
  if (Test-Path "scripts/skill-marketplace-firebreak-selftest.ps1") {
    $checks.Add((Add-Check "skill-marketplace-firebreak-selftest" "$psFile scripts\skill-marketplace-firebreak-selftest.ps1" "private marketplace quarantine firebreak negative selftest detected"))
  }
  if (Test-Path "scripts/promote-marketplace-skill-selftest.ps1") {
    $checks.Add((Add-Check "marketplace-promotion-selftest" "$psFile scripts\promote-marketplace-skill-selftest.ps1" "private marketplace safe promotion selftest detected"))
  }
  if (Test-Path "scripts/kb-release-gate-selftest.ps1") {
    $checks.Add((Add-Check "kb-release-gate-selftest" "$psFile scripts\kb-release-gate-selftest.ps1" "release gate profile selftest detected"))
  }
  if (Test-Path "scripts/skill-surface-minimality-selftest.ps1") {
    $checks.Add((Add-Check "skill-surface-minimality-selftest" "$psFile scripts\skill-surface-minimality-selftest.ps1" "skill/agent minimality classification selftest detected"))
  }
  if (Test-Path "scripts/skill-surface-minimality.ps1") {
    $checks.Add((Add-Check "skill-surface-minimality" "$psFile scripts\skill-surface-minimality.ps1" "static skill/agent minimality report detected"))
  }
  if (Test-Path "scripts/atv-upstream-delta-selftest.ps1") {
    $checks.Add((Add-Check "atv-upstream-delta-selftest" "$psFile scripts\atv-upstream-delta-selftest.ps1" "read-only ATV upstream delta selftest detected"))
  }
  if (Test-Path "scripts/atv-upstream-delta.ps1") {
    $checks.Add((Add-Check "atv-upstream-delta" "$psFile scripts\atv-upstream-delta.ps1" "read-only ATV upstream delta report detected"))
  }
  if (Test-Path "scripts/skill-sync-report.ps1") {
    $checks.Add((Add-Check "skill-sync-report" "$psFile scripts\skill-sync-report.ps1" "skill sync target config detected"))
  }
}

if ($List -or -not $All) {
  $checks | Format-Table -AutoSize
  exit 0
}

foreach ($check in $checks) {
  Write-Host "==> $($check.Name): $($check.Command)"
  Invoke-Expression $check.Command
  if ($LASTEXITCODE -ne 0) {
    Write-Error "check failed: $($check.Name)"
    exit 1
  }
}

exit 0
