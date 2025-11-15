# 🤖 Prueba del Chatbot FinTrack

Write-Host "`n╔══════════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║     🤖 PRUEBA DEL CHATBOT FINTRACK 🤖       ║" -ForegroundColor Cyan
Write-Host "╚══════════════════════════════════════════════╝`n" -ForegroundColor Cyan

$ErrorActionPreference = "Continue"

$backend = "http://localhost:8090"
$frontend = "http://localhost:4200"
$userId = "018c3f3e-51fc-7d7e-8f2a-2d3e4f5a6b7c"
$today = Get-Date -Format "yyyy-MM-dd"

# Test 1: Health
Write-Host "📊 Test 1: Health Check" -ForegroundColor Yellow
try {
    $h = Invoke-RestMethod -Uri "$backend/health"
    Write-Host "✅ Backend OK - $($h.status)" -ForegroundColor Green
} catch {
    Write-Host "❌ Backend ERROR" -ForegroundColor Red
}

# Test 2: Query Backend
Write-Host "`n💬 Test 2: Query Backend" -ForegroundColor Yellow
$headers = @{ "Content-Type" = "application/json"; "X-User-ID" = $userId }
$body = @{ message = "gastos de hoy"; period = @{ from = $today; to = $today }; filters = @{ contextFocus = "expenses" } } | ConvertTo-Json -Depth 3

try {
    $r = Invoke-RestMethod -Uri "$backend/api/chat/query" -Method Post -Headers $headers -Body $body
    Write-Host "✅ $($r.reply)" -ForegroundColor Green
} catch {
    Write-Host "❌ ERROR" -ForegroundColor Red
}

# Test 3: Query Proxy
Write-Host "`n🌐 Test 3: Query Proxy" -ForegroundColor Yellow
try {
    $r = Invoke-RestMethod -Uri "$frontend/api/chat/query" -Method Post -Headers $headers -Body $body
    Write-Host "✅ $($r.reply)" -ForegroundColor Green
} catch {
    Write-Host "❌ ERROR" -ForegroundColor Red
}

Write-Host "`n🎉 Chatbot funcionando!" -ForegroundColor Yellow
Write-Host "📝 Navegador: http://localhost:4200/chatbot" -ForegroundColor Cyan
Write-Host ""
