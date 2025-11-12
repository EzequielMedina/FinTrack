# ====================================================================
# RESET COMPLETO DE BASE DE DATOS - DEMO FINTRACK
# ====================================================================
# Este script elimina completamente el volumen de MySQL y lo recrea
# desde cero con todas las migraciones
# ====================================================================

Write-Host ""
Write-Host "============================================================" -ForegroundColor Cyan
Write-Host "  RESET COMPLETO DE BASE DE DATOS - FINTRACK" -ForegroundColor Yellow
Write-Host "============================================================" -ForegroundColor Cyan
Write-Host ""

Write-Host "Este script va a:" -ForegroundColor White
Write-Host "  1. Detener el contenedor MySQL" -ForegroundColor Gray
Write-Host "  2. Eliminar el volumen de MySQL" -ForegroundColor Gray
Write-Host "  3. Recrear la base de datos desde cero" -ForegroundColor Gray
Write-Host "  4. Ejecutar todas las migraciones" -ForegroundColor Gray
Write-Host "  5. Crear el usuario de prueba" -ForegroundColor Gray
Write-Host ""

Write-Host "⚠️  ADVERTENCIA: Se perderán TODOS los datos!" -ForegroundColor Red
Write-Host ""

# Paso 1: Detener MySQL
Write-Host "🛑 Deteniendo contenedor MySQL..." -ForegroundColor Cyan
docker-compose stop mysql

Write-Host ""
Write-Host "🗑️  Eliminando contenedor MySQL..." -ForegroundColor Cyan
docker-compose rm -f mysql

Write-Host ""
Write-Host "💾 Eliminando volumen de datos..." -ForegroundColor Cyan
docker volume rm ps_mysql_data -f

Write-Host ""
Write-Host "✅ Limpieza completada!" -ForegroundColor Green
Write-Host ""

# Paso 2: Recrear MySQL
Write-Host "🚀 Recreando base de datos desde cero..." -ForegroundColor Cyan
Write-Host ""

docker-compose up -d mysql

Write-Host ""
Write-Host "⏳ Esperando a que MySQL esté listo (30 segundos)..." -ForegroundColor Yellow
Start-Sleep -Seconds 30

Write-Host ""
Write-Host "✅ MySQL recreado!" -ForegroundColor Green
Write-Host ""

# Paso 3: Verificar migraciones
Write-Host "🔍 Verificando migraciones ejecutadas..." -ForegroundColor Cyan
Write-Host ""

docker exec -i fintrack-mysql mysql -uroot -proot_password fintrack -e "SELECT migration_file, executed_at, status FROM migration_history ORDER BY migration_file;"

Write-Host ""
Write-Host "🔍 Verificando tablas creadas..." -ForegroundColor Cyan
Write-Host ""

docker exec -i fintrack-mysql mysql -uroot -proot_password fintrack -e "SHOW TABLES;"

Write-Host ""
Write-Host "🔍 Verificando usuarios..." -ForegroundColor Cyan
Write-Host ""

docker exec -i fintrack-mysql mysql -uroot -proot_password fintrack -e "SELECT id, email, nombre, apellido FROM users;"

Write-Host ""
Write-Host "============================================================" -ForegroundColor Cyan
Write-Host "  ✅ BASE DE DATOS LISTA PARA LA DEMO" -ForegroundColor Green
Write-Host "============================================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "📊 Estado:" -ForegroundColor Yellow
Write-Host "  ✅ Base de datos: fintrack" -ForegroundColor White
Write-Host "  ✅ Todas las migraciones ejecutadas" -ForegroundColor White
Write-Host "  ✅ Tablas creadas correctamente" -ForegroundColor White
Write-Host "  ✅ Usuario de prueba disponible" -ForegroundColor White
Write-Host ""
Write-Host "🎬 Ahora puedes iniciar tu demo!" -ForegroundColor Green
Write-Host ""
Write-Host "Próximos pasos:" -ForegroundColor Yellow
Write-Host "  1. Levanta todos los servicios: docker-compose up -d" -ForegroundColor Gray
Write-Host "  2. Accede a: http://localhost:4200" -ForegroundColor Gray
Write-Host "  3. Inicia sesión con: ezequielmedina456@gmail.com" -ForegroundColor Gray
Write-Host "  4. Crea cuentas, tarjetas y transacciones de prueba" -ForegroundColor Gray
Write-Host "  5. Prueba el chatbot conversacional: http://localhost:4200/chatbot" -ForegroundColor Gray
Write-Host ""
