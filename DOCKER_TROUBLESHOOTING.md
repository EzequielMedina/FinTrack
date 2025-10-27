# 🐳 Docker Desktop - Solución de Problemas

## ❌ Error Actual

```
error during connect: Get "http://%2F%2F.%2Fpipe%2FdockerDesktopLinuxEngine/v1.51/...": 
open //./pipe/dockerDesktopLinuxEngine: El sistema no puede encontrar el archivo especificado.
```

## 🔍 Causa

El **motor de Docker (Docker Engine)** no está disponible, aunque Docker Desktop está ejecutándose. Esto puede pasar cuando:

1. Docker Desktop está iniciándose pero el motor aún no está listo
2. El motor se detuvo inesperadamente
3. Hay un problema con WSL2 (Windows Subsystem for Linux)
4. Docker Desktop necesita reiniciarse

---

## ✅ Solución Aplicada

He reiniciado Docker Desktop con este comando:

```powershell
Stop-Process -Name "Docker Desktop" -Force
Start-Process "C:\Program Files\Docker\Docker\Docker Desktop.exe"
```

---

## ⏱️ Espera de Inicialización

Docker Desktop tarda **1-2 minutos** en iniciar completamente. Debes esperar hasta que:

1. El ícono de Docker Desktop en la barra de tareas esté **verde** ✅
2. Al hacer hover sobre el ícono diga: **"Docker Desktop is running"**
3. No diga "Starting..." o "Starting the Docker Engine..."

---

## 🧪 Verificar que Docker Está Listo

### **Paso 1: Espera 2 minutos**

Mientras tanto, puedes ver el estado en la barra de tareas.

### **Paso 2: Verifica con este comando**

```powershell
docker info
```

**Resultado esperado:**
```
Server:
 Containers: X
 Running: X
 Paused: 0
 Stopped: X
 Images: X
```

**Si aún da error:**
```
error during connect: Get "http://...": open //./pipe/dockerDesktopLinuxEngine: El sistema no puede encontrar el archivo especificado.
```

→ **Espera 1 minuto más** y vuelve a intentar.

---

## 🚀 Luego de que Docker Esté Listo

### **1. Verificar contenedores actuales**

```powershell
docker ps -a
```

### **2. Limpiar contenedores viejos (opcional)**

```powershell
docker-compose down
```

### **3. Reconstruir sin caché**

```powershell
docker-compose build --no-cache
```

### **4. Levantar todos los servicios**

```powershell
docker-compose up
```

---

## 🔧 Soluciones Alternativas

### **Si el reinicio no funciona:**

#### **Opción 1: Reiniciar WSL2**

```powershell
# Detener WSL2
wsl --shutdown

# Esperar 10 segundos
Start-Sleep -Seconds 10

# Reiniciar Docker Desktop
Start-Process "C:\Program Files\Docker\Docker\Docker Desktop.exe"
```

#### **Opción 2: Verificar WSL2**

```powershell
# Ver distribuciones WSL
wsl --list --verbose

# Debería mostrar:
# NAME                   STATE           VERSION
# docker-desktop         Running         2
# docker-desktop-data    Running         2
```

#### **Opción 3: Reiniciar el servicio de Docker**

1. Abre Docker Desktop (interfaz gráfica)
2. Ve a: **Settings** → **General**
3. Marca/desmarca **"Start Docker Desktop when you log in"**
4. Click **"Restart"**

#### **Opción 4: Reiniciar Windows (última opción)**

Si nada funciona, reinicia tu PC.

---

## ⚠️ Problemas Comunes

### **Error: "Docker Desktop is stopping..."**

```powershell
# Forzar detención
taskkill /F /IM "Docker Desktop.exe"

# Esperar 5 segundos
Start-Sleep -Seconds 5

# Iniciar nuevamente
Start-Process "C:\Program Files\Docker\Docker\Docker Desktop.exe"
```

### **Error: "WSL 2 installation is incomplete"**

1. Instala el kernel de WSL2:
   ```powershell
   wsl --update
   ```

2. Reinicia Docker Desktop

### **Error: "Docker daemon is not running"**

1. Abre Docker Desktop manualmente (doble clic en el ícono del escritorio)
2. Espera a que diga "Docker Desktop is running"
3. Intenta de nuevo

---

## 📊 Estados de Docker Desktop

| Ícono | Estado | Acción |
|-------|--------|--------|
| 🟢 Verde | Running | ✅ Listo para usar |
| 🟡 Amarillo | Starting | ⏱️ Espera 1-2 minutos |
| 🔴 Rojo | Stopped/Error | ❌ Necesita reiniciarse |
| 🔵 Azul animado | Updating | ⏱️ Espera a que termine |

---

## ✅ Checklist de Verificación

Antes de ejecutar `docker-compose up`:

- [ ] Docker Desktop está corriendo (proceso visible)
- [ ] Ícono de Docker en barra de tareas está **verde** ✅
- [ ] `docker info` devuelve información del servidor (sin errores)
- [ ] `docker ps` funciona (muestra lista de contenedores)
- [ ] WSL2 está corriendo (`wsl --list --verbose`)

---

## 🎯 Comandos Útiles

```powershell
# Ver si Docker Desktop está corriendo
Get-Process "Docker Desktop" -ErrorAction SilentlyContinue

# Ver contenedores activos
docker ps

# Ver todos los contenedores (incluso detenidos)
docker ps -a

# Ver imágenes
docker images

# Ver uso de espacio
docker system df

# Limpiar todo (¡CUIDADO! Elimina todo)
docker system prune -a --volumes
```

---

## 📞 Próximos Pasos

### **Una vez que Docker esté listo (ícono verde):**

1. **Para el frontend Angular** (desarrollo local):
   ```powershell
   cd frontend
   ng serve --host 0.0.0.0 --port 4200
   ```
   
   → Accede a: `http://localhost:4200`

2. **Para los servicios backend** (Docker):
   ```powershell
   docker-compose up mysql user-service account-service transaction-service
   ```

3. **O todos los servicios a la vez**:
   ```powershell
   docker-compose up
   ```

---

## 🌐 Frontend (Angular) - Sin Docker

**El frontend NO necesita Docker para desarrollo.** Es más rápido ejecutarlo directamente:

```powershell
cd c:\Facultad\Alumno\PS\frontend
ng serve --host 0.0.0.0 --port 4200
```

**Ventajas:**
- ✅ Hot reload instantáneo
- ✅ Compilación más rápida
- ✅ Menos uso de recursos
- ✅ Mejor experiencia de desarrollo

**Usa Docker solo para backend (microservicios).**

---

## 🔄 Workflow Recomendado

```powershell
# 1. Inicia Docker Desktop (espera a que esté verde)
# 2. Levanta solo los servicios backend
docker-compose up mysql user-service account-service transaction-service

# 3. En otra terminal, inicia el frontend
cd frontend
ng serve --host 0.0.0.0 --port 4200

# 4. Abre el navegador
# http://localhost:4200
```

---

**⏱️ RECUERDA: Docker Desktop tarda 1-2 minutos en iniciar completamente.**

**Espera a que el ícono esté verde antes de ejecutar comandos de Docker.**
