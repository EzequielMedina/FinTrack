# 🤖 Automatización de Backlog - FinTrack

## 📋 Descripción

Este sistema automatiza la carga de tareas del cronograma de FinTrack a diferentes plataformas de gestión de proyectos, eliminando la necesidad de crear manualmente cada tarea en el backlog.

## 🎯 Plataformas Soportadas

- ✅ **Jira** (Atlassian)
- ✅ **GitHub Projects** (v2)
- ✅ **Azure DevOps** (Boards)
- ✅ **Trello**
- ✅ **Linear**
- ✅ **Exportación CSV** (para importación manual)

## 📁 Archivos del Sistema

```
FinTrack/
├── FinTrack_Sprint_Backlog_Cronograma.md  # Cronograma principal
├── automate_backlog_upload.py             # Script de automatización
├── config_example.json                    # Configuración de ejemplo
├── config.json                           # Tu configuración (crear)
└── README_Automatizacion_Backlog.md      # Este archivo
```

## 🚀 Instalación y Configuración

### 1. Requisitos Previos

```bash
# Instalar Python 3.8+
python --version

# Instalar dependencias
pip install requests
```

### 2. Configuración Inicial

#### Opción A: Crear configuración automáticamente
```bash
python automate_backlog_upload.py --create-config
```

#### Opción B: Copiar configuración de ejemplo
```bash
cp config_example.json config.json
```

### 3. Configurar Credenciales

Edita `config.json` con tus credenciales:

#### Para Jira:
1. Ve a: https://id.atlassian.com/manage-profile/security/api-tokens
2. Crea un nuevo API token
3. Actualiza en `config.json`:
   ```json
   {
     "jira": {
       "base_url": "https://tu-dominio.atlassian.net",
       "email": "tu-email@ejemplo.com",
       "api_token": "ATATT3xFfGF0...",
       "project_key": "FINTRACK"
     }
   }
   ```

#### Para GitHub:
1. Ve a: https://github.com/settings/tokens
2. Crea un Personal Access Token con permisos: `repo`, `project`
3. Actualiza en `config.json`:
   ```json
   {
     "github": {
       "token": "ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
       "repo": "usuario/fintrack",
       "owner": "tu-usuario"
     }
   }
   ```

#### Para Azure DevOps:
1. Ve a: https://dev.azure.com/tu-org/_usersSettings/tokens
2. Crea un Personal Access Token con permisos: `Work Items (Read & Write)`
3. Actualiza en `config.json`:
   ```json
   {
     "azure": {
       "organization": "tu-organizacion",
       "project": "FinTrack",
       "personal_access_token": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
     }
   }
   ```

## 🎮 Uso del Sistema

### Comandos Básicos

#### 1. Exportar a CSV (Recomendado para empezar)
```bash
python automate_backlog_upload.py --platform csv
```
**Resultado:** Genera `fintrack_backlog.csv` para importación manual

#### 2. Subir a Jira
```bash
python automate_backlog_upload.py --platform jira --config config.json
```

#### 3. Subir a GitHub Projects
```bash
python automate_backlog_upload.py --platform github --config config.json
```

#### 4. Subir a Azure DevOps
```bash
python automate_backlog_upload.py --platform azure --config config.json
```

### Comandos Avanzados

#### Usar archivo de cronograma personalizado
```bash
python automate_backlog_upload.py --platform jira --markdown mi_cronograma.md
```

#### Usar configuración personalizada
```bash
python automate_backlog_upload.py --platform github --config mi_config.json
```

## 📊 Estructura de Tareas Generadas

Cada tarea incluye:

- **ID único:** TASK-001, TASK-002, etc.
- **Título descriptivo:** Basado en el cronograma
- **Descripción detallada:** Contexto y objetivos
- **Sprint asignado:** 1-6 según cronograma
- **Estimación:** Días de trabajo estimados
- **Prioridad:** Alta, Media, Baja
- **Criterios de aceptación:** Definición de "terminado"
- **Labels/Tags:** Para categorización
- **Fechas:** Inicio y fin del sprint

### Ejemplo de Tarea Generada:

```
ID: TASK-001
Título: Configurar repositorio Git con estructura de microservicios
Descripción: Establecer la estructura base del repositorio con separación frontend/backend/docs y configuración inicial
Sprint: 1
Estimación: 1 día
Prioridad: Alta
Criterios: Repo con estructura frontend/backend/docs
Labels: setup, git, infrastructure
Fechas: 2024-09-10 a 2024-09-24
```

## 🔧 Personalización

### Modificar Tareas

Edita el archivo `automate_backlog_upload.py` en la sección `sprint_tasks` para:

- Agregar nuevas tareas
- Modificar estimaciones
- Cambiar prioridades
- Actualizar criterios de aceptación
- Añadir dependencias

### Agregar Nueva Plataforma

1. Crea un nuevo método `upload_to_nueva_plataforma()`
2. Implementa la lógica de API específica
3. Agrega la opción en el método `run()`

## 📈 Flujo de Trabajo Recomendado

### Para Equipos Nuevos:

1. **Exportar a CSV primero**
   ```bash
   python automate_backlog_upload.py --platform csv
   ```

2. **Revisar el archivo CSV generado**
   - Verificar que las tareas sean correctas
   - Ajustar estimaciones si es necesario

3. **Importar manualmente** en tu plataforma preferida

4. **Una vez validado, usar automatización directa**

### Para Equipos Experimentados:

1. **Configurar credenciales**
2. **Ejecutar directamente** en la plataforma elegida
3. **Verificar resultados** en la plataforma

## 🛠️ Troubleshooting

### Errores Comunes

#### Error de Autenticación
```
❌ Error conectando con Jira: 401 Unauthorized
```
**Solución:** Verificar credenciales en `config.json`

#### Error de Permisos
```
❌ Error creando tarea TASK-001: 403 Forbidden
```
**Solución:** Verificar permisos del token/usuario

#### Error de Conexión
```
❌ Error conectando con GitHub: Connection timeout
```
**Solución:** Verificar conexión a internet y URLs

### Logs y Debugging

El script muestra progreso en tiempo real:
```
🚀 Iniciando automatización de carga de backlog...
📋 Cargando tareas del cronograma...
✅ 33 tareas cargadas
✅ Tarea TASK-001 creada exitosamente en Jira
✅ Tarea TASK-002 creada exitosamente en Jira
...
🎉 ¡Automatización completada exitosamente!
```

## 📋 Checklist de Verificación

Antes de ejecutar la automatización:

- [ ] Python 3.8+ instalado
- [ ] Dependencias instaladas (`pip install requests`)
- [ ] Archivo `config.json` creado y configurado
- [ ] Credenciales válidas y con permisos correctos
- [ ] Proyecto/Board creado en la plataforma destino
- [ ] Conexión a internet estable

Después de ejecutar:

- [ ] Verificar que todas las tareas se crearon
- [ ] Revisar que los sprints estén correctamente asignados
- [ ] Confirmar que las estimaciones son correctas
- [ ] Validar que los criterios de aceptación están completos

## 🔄 Actualizaciones y Mantenimiento

### Agregar Nuevas Tareas

1. Edita `FinTrack_Sprint_Backlog_Cronograma.md`
2. Actualiza `sprint_tasks` en `automate_backlog_upload.py`
3. Ejecuta nuevamente el script

### Modificar Cronograma

1. Cambia las fechas en la variable `base_date`
2. Ajusta la duración de sprints si es necesario
3. Re-ejecuta la automatización

## 🎯 Próximas Mejoras

- [ ] Soporte para Notion
- [ ] Integración con Monday.com
- [ ] Sincronización bidireccional
- [ ] Interface web para configuración
- [ ] Plantillas personalizables
- [ ] Reportes de progreso automáticos

## 📞 Soporte

Si encuentras problemas:

1. Revisa este README
2. Verifica la configuración
3. Consulta los logs de error
4. Contacta al equipo de desarrollo

---

**¡Automatiza tu backlog y enfócate en desarrollar! 🚀**

*Generado para el proyecto FinTrack - Tecnicatura en Programación UNT*