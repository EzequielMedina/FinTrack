# Script de verificación post-deploy del FAQ con Soporte

Write-Host "================================================" -ForegroundColor Cyan
Write-Host "   VERIFICACIÓN: FAQ + Modal de Soporte" -ForegroundColor Yellow
Write-Host "================================================" -ForegroundColor Cyan
Write-Host ""

Write-Host "✅ Cambios implementados:" -ForegroundColor Green
Write-Host "   1. Sección FAQ con 20+ preguntas" -ForegroundColor White
Write-Host "   2. Modal de soporte integrado" -ForegroundColor White
Write-Host "   3. Envío de emails via notification-service" -ForegroundColor White
Write-Host "   4. Template EmailJS: template_yst8bd2" -ForegroundColor White
Write-Host ""

Write-Host "================================================" -ForegroundColor Cyan
Write-Host "   PASOS PARA PROBAR" -ForegroundColor Yellow
Write-Host "================================================" -ForegroundColor Cyan
Write-Host ""

Write-Host "1️⃣  Abre el navegador en:" -ForegroundColor Cyan
Write-Host "   http://localhost:4200/faq" -ForegroundColor White
Write-Host ""

Write-Host "2️⃣  Haz Ctrl + Shift + R para limpiar cache" -ForegroundColor Cyan
Write-Host ""

Write-Host "3️⃣  Verifica que veas:" -ForegroundColor Cyan
Write-Host "   ✓ Filtros de categorías (Todas, General, Cuentas, etc.)" -ForegroundColor Gray
Write-Host "   ✓ Preguntas frecuentes con acordeón" -ForegroundColor Gray
Write-Host "   ✓ Card de 'Contactar Soporte' al final" -ForegroundColor Gray
Write-Host ""

Write-Host "4️⃣  Click en 'Contactar Soporte'" -ForegroundColor Cyan
Write-Host "   ✓ Debe abrirse un MODAL (ventana emergente)" -ForegroundColor Gray
Write-Host "   ✓ Con formulario: Nombre, Email, Asunto, Mensaje" -ForegroundColor Gray
Write-Host "   ✓ Datos del usuario pre-rellenados" -ForegroundColor Gray
Write-Host ""

Write-Host "5️⃣  Completa y envía el formulario:" -ForegroundColor Cyan
Write-Host "   ✓ Debe mostrar spinner 'Enviando...'" -ForegroundColor Gray
Write-Host "   ✓ Debe mostrar mensaje de éxito" -ForegroundColor Gray
Write-Host "   ✓ Modal se cierra automáticamente" -ForegroundColor Gray
Write-Host ""

Write-Host "================================================" -ForegroundColor Cyan
Write-Host "   TROUBLESHOOTING" -ForegroundColor Yellow
Write-Host "================================================" -ForegroundColor Cyan
Write-Host ""

Write-Host "❌ Si el modal NO se abre:" -ForegroundColor Red
Write-Host "   1. Abre DevTools (F12)" -ForegroundColor White
Write-Host "   2. Ve a la pestaña Console" -ForegroundColor White
Write-Host "   3. Busca errores en rojo" -ForegroundColor White
Write-Host "   4. Verifica que notification-service esté corriendo:" -ForegroundColor White
Write-Host "      docker-compose ps notification-service" -ForegroundColor Gray
Write-Host ""

Write-Host "❌ Si el email NO se envía:" -ForegroundColor Red
Write-Host "   1. Verifica logs del backend:" -ForegroundColor White
Write-Host "      docker-compose logs notification-service --tail 50" -ForegroundColor Gray
Write-Host "   2. Verifica la configuración de EmailJS en:" -ForegroundColor White
Write-Host "      backend/services/notification-service/config.json" -ForegroundColor Gray
Write-Host ""

Write-Host "================================================" -ForegroundColor Cyan
Write-Host "   ENDPOINT DE SOPORTE" -ForegroundColor Yellow
Write-Host "================================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "URL: POST http://localhost:8084/api/notifications/support" -ForegroundColor White
Write-Host ""
Write-Host "Request Body:" -ForegroundColor Gray
Write-Host '{' -ForegroundColor DarkGray
Write-Host '  "name": "Tu Nombre",' -ForegroundColor DarkGray
Write-Host '  "email": "tu@email.com",' -ForegroundColor DarkGray
Write-Host '  "subject": "Asunto",' -ForegroundColor DarkGray
Write-Host '  "message": "Mensaje..."' -ForegroundColor DarkGray
Write-Host '}' -ForegroundColor DarkGray
Write-Host ""

Write-Host "================================================" -ForegroundColor Green
Write-Host "   ¡Listo para probar!" -ForegroundColor Green
Write-Host "================================================" -ForegroundColor Green
Write-Host ""
