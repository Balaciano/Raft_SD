param(
    [ValidateRange(1, 1000)][int]$Repeticoes = 1,
    [switch]$Race
)
$ErrorActionPreference = 'Stop'
$raizProjeto = Split-Path -Parent $PSScriptRoot
$gopathAnterior = $env:GOPATH
$modulesAnterior = $env:GO111MODULE
try {
    if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
        throw 'Go nao encontrado no PATH. Instale o Go e abra um novo terminal.'
    }
    $env:GOPATH = $raizProjeto
    $env:GO111MODULE = 'off'
    $pastaResultados = Join-Path $raizProjeto 'resultados'
    New-Item -ItemType Directory -Path $pastaResultados -Force | Out-Null
    $sufixo = if ($Race) { 'race' } else { 'normal' }
    $arquivoLog = Join-Path $pastaResultados ("2A-{0}-{1}.log" -f $sufixo, (Get-Date -Format 'yyyyMMdd-HHmmss'))
    $argumentosGo = @('test', '-v', '-run', '2A', "-count=$Repeticoes", '-timeout=30m')
    if ($Race) { $argumentosGo += '-race' }
    $argumentosGo += 'raft'
    & go version | Tee-Object -FilePath $arquivoLog
    & go @argumentosGo 2>&1 | Tee-Object -FilePath $arquivoLog -Append
    if ($LASTEXITCODE -ne 0) { throw "Testes falharam. Consulte $arquivoLog" }
    Write-Host "Testes concluidos. Evidencia: $arquivoLog"
} finally {
    $env:GOPATH = $gopathAnterior
    $env:GO111MODULE = $modulesAnterior
}
