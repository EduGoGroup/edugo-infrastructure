# Módulo MongoDB - edugo-infrastructure

Módulo de migraciones para MongoDB del ecosistema EduGo.

## 🎯 Propósito

Gestionar las migraciones de MongoDB con schemas, índices y validaciones de forma centralizada.

## 🗄️ Colecciones Gestionadas

| Migración | Colección | Descripción |
|-----------|-----------|-------------|
| 001 | material_assessment | Contenido de assessments/quizzes generados por IA |
| 002 | material_content | Contenido procesado de materiales educativos |
| 003 | assessment_attempt_result | Resultados detallados de intentos |
| 004 | audit_logs | Logs de auditoría del sistema (TTL: 90 días) |
| 005 | notifications | Notificaciones para usuarios (TTL: 30 días) |
| 006 | analytics_events | Eventos de analítica y comportamiento (TTL: 365 días) |

## 🚀 Uso

### Ejecutar Migraciones

```bash
cd mongodb
go run migrate.go up
```

### Ver Estado

```bash
go run migrate.go status
```

### Revertir Última Migración

```bash
go run migrate.go down
```

### Crear Nueva Migración

```bash
go run migrate.go create "add_new_collection"
```

## 🔧 Variables de Entorno

```bash
MONGO_HOST=localhost
MONGO_PORT=27017
MONGO_DB_NAME=edugo
MONGO_USER=     # opcional
MONGO_PASSWORD= # opcional
```

## 📋 Requisitos

- **mongosh** instalado (para ejecutar scripts JavaScript)

## 📦 Importar en Proyectos

```go
import "github.com/EduGoGroup/edugo-infrastructure/mongodb"
```

## 📚 Documentación

Ver documentación completa de schemas en: `../docs/MONGODB_SCHEMA.md`

---

**Versión:** 0.5.0  
**Mantenedores:** Equipo EduGo
