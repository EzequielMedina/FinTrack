# Reset completo de base de datos - FinTrack Demo
# Este script elimina y recrea el volumen de MySQL

Write-Host ""
Write-Host "===========================================================" -ForegroundColor Cyan
Write-Host "  RESET COMPLETO DE BASE DE DATOS - FINTRACK" -ForegroundColor Yellow
Write-Host "===========================================================" -ForegroundColor Cyan
Write-Host ""

Write-Host "Este script va a:" -ForegroundColor White
Write-Host "  1. Detener el contenedor MySQL" -ForegroundColor Gray
Write-Host "  2. Eliminar el volumen de MySQL" -ForegroundColor Gray
Write-Host "  3. Recrear la base de datos desde cero" -ForegroundColor Gray
Write-Host "  4. Ejecutar todas las migraciones" -ForegroundColor Gray
Write-Host ""

Write-Host "ADVERTENCIA: Se perderan TODOS los datos!" -ForegroundColor Red
Write-Host ""

# Paso 1: Detener MySQL
Write-Host "[1/4] Deteniendo contenedor MySQL..." -ForegroundColor Cyan
docker-compose stop mysql

Write-Host ""
Write-Host "[2/4] Eliminando contenedor y volumen..." -ForegroundColor Cyan
docker-compose rm -f mysql
docker volume rm ps_mysql_data -f

Write-Host ""
Write-Host "[3/4] Recreando MySQL desde cero..." -ForegroundColor Cyan
docker-compose up -d mysql

Write-Host ""
Write-Host "Esperando a que MySQL este listo (30 segundos)..." -ForegroundColor Yellow
Start-Sleep -Seconds 30

Write-Host ""
Write-Host "[4/4] Verificando migraciones y tablas..." -ForegroundColor Cyan
Write-Host ""

Write-Host "--- Migraciones ejecutadas ---" -ForegroundColor Yellow
docker exec -i fintrack-mysql mysql -uroot -proot_password fintrack -e "SELECT migration_file, executed_at FROM migration_history ORDER BY migration_file;"

Write-Host ""
Write-Host "--- Tablas creadas ---" -ForegroundColor Yellow
docker exec -i fintrack-mysql mysql -uroot -proot_password fintrack -e "SHOW TABLES;"

Write-Host ""
Write-Host "--- Usuarios en la base ---" -ForegroundColor Yellow
docker exec -i fintrack-mysql mysql -uroot -proot_password fintrack -e "SELECT id, email, nombre, apellido FROM users;"

Write-Host ""
Write-Host "===========================================================" -ForegroundColor Cyan
Write-Host "  BASE DE DATOS LISTA PARA LA DEMO" -ForegroundColor Green
Write-Host "===========================================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Ahora puedes:" -ForegroundColor Yellow
Write-Host "  1. Levantar todos los servicios: docker-compose up -d" -ForegroundColor Gray
Write-Host "  2. Acceder a: http://localhost:4200" -ForegroundColor Gray
Write-Host "  3. Iniciar sesion con: ezequielmedina456@gmail.com" -ForegroundColor Gray
Write-Host "  4. Probar el chatbot: http://localhost:4200/chatbot" -ForegroundColor Gray
Write-Host ""
