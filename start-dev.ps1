param(
    [string]$ApiKey = 'dev-key-change-me-now-but-long-enough-12345',
    [switch]$BindAll,
    [string]$Port = '3000'
)

$backendDir = Join-Path $PSScriptRoot 'backend'
Set-Location $backendDir

$env:APP_ENV = 'development'
$env:API_KEYS = $ApiKey
$env:PORT = $Port

if ($BindAll) {
    $env:BIND_ADDR = '0.0.0.0'
}

go run -tags sqlite_fts5 .