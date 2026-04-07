param(
    [string]$BaselineRef = "HEAD~1",
    [string]$CandidateRef = "HEAD",
    [string]$InputRelativePath = "test/hello.go",
    [string]$ModuleName = "hello",
    [int]$Iterations = 5000,
    [int]$Warmup = 500,
    [string]$ReportPath = "output/perf-size-report.json",
    [string]$BaselineNodePath = "",
    [string]$CandidateNodePath = "",
    [switch]$SkipBenchmark,
    [switch]$KeepWorktrees
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

function Write-Stage {
    param([string]$Message)
    Write-Host "==> $Message"
}

function Ensure-PathExists {
    param(
        [string]$Path,
        [string]$Label
    )
    # Fail early with clear diagnostics when an expected artifact is missing.
    if (-not (Test-Path -LiteralPath $Path)) {
        throw "$Label not found: $Path"
    }
}

function Get-PESectionSizes {
    param([string]$FilePath)

    # Read PE bytes and parse section table so report can quantify resource payload ratio.
    $bytes = [System.IO.File]::ReadAllBytes($FilePath)
    $peOffset = [BitConverter]::ToInt32($bytes, 0x3C)
    $numSections = [BitConverter]::ToUInt16($bytes, $peOffset + 6)
    $optSize = [BitConverter]::ToUInt16($bytes, $peOffset + 20)
    $sectionOffset = $peOffset + 24 + $optSize

    $map = @{}
    for ($i = 0; $i -lt $numSections; $i++) {
        $off = $sectionOffset + 40 * $i
        $nameBytes = $bytes[$off..($off + 7)]
        $name = ([System.Text.Encoding]::ASCII.GetString($nameBytes)).Trim([char]0)
        $rawSize = [BitConverter]::ToUInt32($bytes, $off + 16)
        $map[$name] = [int64]$rawSize
    }

    return $map
}

function Get-NodeArtifactMetrics {
    param([string]$NodePath)

    # Build comparable size metrics from artifact bytes and PE section details.
    $fileInfo = Get-Item -LiteralPath $NodePath
    $sections = Get-PESectionSizes -FilePath $NodePath
    $resourceBytes = 0
    if ($sections.ContainsKey(".rsrc")) {
        $resourceBytes = [int64]$sections[".rsrc"]
    }
    $wrapperBytes = [int64]$fileInfo.Length - $resourceBytes
    if ($wrapperBytes -lt 0) {
        $wrapperBytes = 0
    }
    $resourcePct = 0.0
    if ($fileInfo.Length -gt 0) {
        $resourcePct = [math]::Round(($resourceBytes * 100.0) / $fileInfo.Length, 2)
    }

    return [ordered]@{
        nodePath = $NodePath
        nodeBytes = [int64]$fileInfo.Length
        resourceBytes = [int64]$resourceBytes
        wrapperBytes = [int64]$wrapperBytes
        resourcePct = $resourcePct
    }
}

function Invoke-NodeBenchmark {
    param(
        [string]$RepoRoot,
        [string]$NodePath,
        [string]$OutputPath,
        [int]$Iterations,
        [int]$Warmup
    )

    # Run the dedicated benchmark script against one module and parse JSON result.
    $benchScript = Join-Path $RepoRoot "scripts/benchmark_module.js"
    Ensure-PathExists -Path $benchScript -Label "Benchmark script"
    Ensure-PathExists -Path $NodePath -Label "Node module"

    & node $benchScript "--module=$NodePath" "--iterations=$Iterations" "--warmup=$Warmup" "--output=$OutputPath" | Out-Null
    Ensure-PathExists -Path $OutputPath -Label "Benchmark output"

    $json = Get-Content -LiteralPath $OutputPath -Raw
    return ($json | ConvertFrom-Json)
}

function Invoke-BuildFromRef {
    param(
        [string]$RepoRoot,
        [string]$Ref,
        [string]$Label,
        [string]$InputRelativePath,
        [string]$ModuleName,
        [string]$WorkspaceRoot
    )

    # Build one git ref in an isolated worktree and return generated module path.
    $worktreePath = Join-Path $WorkspaceRoot "worktree-$Label"
    $outputDir = Join-Path $WorkspaceRoot "$Label-output"
    $toolPath = Join-Path $WorkspaceRoot "$Label-go-node.exe"

    if (Test-Path -LiteralPath $worktreePath) {
        Remove-Item -LiteralPath $worktreePath -Recurse -Force
    }
    New-Item -ItemType Directory -Path $outputDir -Force | Out-Null

    & git -C $RepoRoot worktree add --detach $worktreePath $Ref

    try {
        & go build -o $toolPath . 2>&1 | Out-Host
        if (-not (Test-Path -LiteralPath $toolPath)) {
            throw "Failed to build go-node binary for $Label ($Ref)"
        }

        $inputPath = Join-Path $worktreePath $InputRelativePath
        Ensure-PathExists -Path $inputPath -Label "Input source for $Label"

        & $toolPath "-input=$inputPath" "-name=$ModuleName" "-output=$outputDir" 2>&1 | Out-Host
        $nodePath = Join-Path $outputDir "$ModuleName.node"
        Ensure-PathExists -Path $nodePath -Label "Compiled module for $Label"

        return [ordered]@{
            label = $Label
            source = "ref:$Ref"
            nodePath = $nodePath
            worktreePath = $worktreePath
        }
    }
    catch {
        throw
    }
}

function Get-BenchmarkMap {
    param($BenchmarkReport)

    # Transform benchmark array into a lookup table for stable delta calculations.
    $map = @{}
    foreach ($item in $BenchmarkReport.benchmarks) {
        $map[$item.name] = $item
    }
    return $map
}

function Compare-Benchmarks {
    param(
        $BaselineReport,
        $CandidateReport
    )

    # Compute throughput and latency deltas per benchmark case.
    $baselineMap = Get-BenchmarkMap -BenchmarkReport $BaselineReport
    $candidateMap = Get-BenchmarkMap -BenchmarkReport $CandidateReport
    $rows = @()

    foreach ($name in $baselineMap.Keys) {
        if (-not $candidateMap.ContainsKey($name)) {
            continue
        }
        $b = $baselineMap[$name]
        $c = $candidateMap[$name]
        $opsDeltaPct = 0.0
        if ($b.ops -gt 0) {
            $opsDeltaPct = [math]::Round((($c.ops - $b.ops) * 100.0) / $b.ops, 2)
        }
        $avgDeltaPct = 0.0
        if ($b.avgNs -gt 0) {
            $avgDeltaPct = [math]::Round((($c.avgNs - $b.avgNs) * 100.0) / $b.avgNs, 2)
        }
        $rows += [ordered]@{
            name = $name
            baselineOps = [math]::Round([double]$b.ops, 2)
            candidateOps = [math]::Round([double]$c.ops, 2)
            opsDeltaPct = $opsDeltaPct
            baselineAvgNs = [math]::Round([double]$b.avgNs, 2)
            candidateAvgNs = [math]::Round([double]$c.avgNs, 2)
            avgLatencyDeltaPct = $avgDeltaPct
            baselineP95Ns = [math]::Round([double]$b.p95Ns, 2)
            candidateP95Ns = [math]::Round([double]$c.p95Ns, 2)
        }
    }

    return $rows
}

$repoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
$reportAbsPath = Join-Path $repoRoot $ReportPath
$reportDir = Split-Path -Parent $reportAbsPath
New-Item -ItemType Directory -Path $reportDir -Force | Out-Null

$workspaceRoot = Join-Path $repoRoot ".perf-compare"
New-Item -ItemType Directory -Path $workspaceRoot -Force | Out-Null

$baselineBuild = $null
$candidateBuild = $null

try {
    if ($BaselineNodePath -and $CandidateNodePath) {
        Write-Stage "Using prebuilt modules"
        $baselineNodeAbs = (Resolve-Path $BaselineNodePath).Path
        $candidateNodeAbs = (Resolve-Path $CandidateNodePath).Path
        Ensure-PathExists -Path $baselineNodeAbs -Label "Baseline module"
        Ensure-PathExists -Path $candidateNodeAbs -Label "Candidate module"

        $baselineBuild = [ordered]@{
            label = "baseline"
            source = "prebuilt"
            nodePath = $baselineNodeAbs
            worktreePath = ""
        }
        $candidateBuild = [ordered]@{
            label = "candidate"
            source = "prebuilt"
            nodePath = $candidateNodeAbs
            worktreePath = ""
        }
    }
    else {
        Write-Stage "Building baseline from $BaselineRef"
        $baselineBuild = Invoke-BuildFromRef `
            -RepoRoot $repoRoot `
            -Ref $BaselineRef `
            -Label "baseline" `
            -InputRelativePath $InputRelativePath `
            -ModuleName $ModuleName `
            -WorkspaceRoot $workspaceRoot

        Write-Stage "Building candidate from $CandidateRef"
        $candidateBuild = Invoke-BuildFromRef `
            -RepoRoot $repoRoot `
            -Ref $CandidateRef `
            -Label "candidate" `
            -InputRelativePath $InputRelativePath `
            -ModuleName $ModuleName `
            -WorkspaceRoot $workspaceRoot
    }

    Write-Stage "Collecting size metrics"
    $baselineMetrics = Get-NodeArtifactMetrics -NodePath $baselineBuild.nodePath
    $candidateMetrics = Get-NodeArtifactMetrics -NodePath $candidateBuild.nodePath

    $baselineBench = $null
    $candidateBench = $null
    $benchmarkDiff = @()

    if (-not $SkipBenchmark) {
        Write-Stage "Running benchmark on baseline module"
        $baselineBenchPath = Join-Path $workspaceRoot "baseline-benchmark.json"
        $baselineBench = Invoke-NodeBenchmark `
            -RepoRoot $repoRoot `
            -NodePath $baselineBuild.nodePath `
            -OutputPath $baselineBenchPath `
            -Iterations $Iterations `
            -Warmup $Warmup

        Write-Stage "Running benchmark on candidate module"
        $candidateBenchPath = Join-Path $workspaceRoot "candidate-benchmark.json"
        $candidateBench = Invoke-NodeBenchmark `
            -RepoRoot $repoRoot `
            -NodePath $candidateBuild.nodePath `
            -OutputPath $candidateBenchPath `
            -Iterations $Iterations `
            -Warmup $Warmup

        $benchmarkDiff = Compare-Benchmarks -BaselineReport $baselineBench -CandidateReport $candidateBench
    }

    $sizeDeltaPct = 0.0
    if ($baselineMetrics.nodeBytes -gt 0) {
        $sizeDeltaPct = [math]::Round((($candidateMetrics.nodeBytes - $baselineMetrics.nodeBytes) * 100.0) / $baselineMetrics.nodeBytes, 2)
    }

    $report = [ordered]@{
        generatedAt = (Get-Date).ToString("o")
        repoRoot = $repoRoot
        baseline = [ordered]@{
            source = $baselineBuild.source
            metrics = $baselineMetrics
            benchmark = $baselineBench
        }
        candidate = [ordered]@{
            source = $candidateBuild.source
            metrics = $candidateMetrics
            benchmark = $candidateBench
        }
        summary = [ordered]@{
            nodeSizeDeltaPct = $sizeDeltaPct
            benchmarkSkipped = [bool]$SkipBenchmark
            benchmarkDiff = $benchmarkDiff
        }
    }

    $json = $report | ConvertTo-Json -Depth 8
    Set-Content -LiteralPath $reportAbsPath -Value $json -Encoding UTF8

    Write-Stage "Report written to $reportAbsPath"
    Write-Host ("Baseline node bytes : {0}" -f $baselineMetrics.nodeBytes)
    Write-Host ("Candidate node bytes: {0}" -f $candidateMetrics.nodeBytes)
    Write-Host ("Node size delta (%): {0}" -f $sizeDeltaPct)
}
finally {
    # Clean worktrees unless caller requests retention for debugging.
    if (-not $KeepWorktrees) {
        foreach ($path in @($baselineBuild.worktreePath, $candidateBuild.worktreePath)) {
            if ($path -and (Test-Path -LiteralPath $path)) {
                try {
                    & git -C $repoRoot worktree remove --force $path 2>&1 | Out-Null
                }
                catch {
                }
            }
        }
    }
}
