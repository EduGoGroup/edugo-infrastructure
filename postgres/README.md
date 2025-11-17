# Módulo PostgreSQL - edugo-infrastructure

Módulo de migraciones para PostgreSQL del ecosistema EduGo.

## 🎯 Propósito

Gestionar las migraciones de base de datos PostgreSQL de forma centralizada y controlada.

## 🗄️ Tablas Gestionadas

| Migración | Tabla | Descripción |
|-----------|-------|-------------|
| 001 | users | Usuarios del sistema (admin, teacher, student, guardian) |
| 002 | schools | Instituciones educativas |
| 003 | academic_units | Cursos, clases, secciones |
| 004 | memberships | Relación usuario-escuela-curso |
| 005 | materials | Materiales educativos |
| 006 | assessment | Quizzes (referencia a MongoDB) |
| 007 | assessment_attempt | Intentos de estudiantes |
| 008 | assessment_attempt_answer | Respuestas individuales |

## 🚀 Uso

### Ejecutar Migraciones

```bash
cd postgres
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
go run migrate.go create "add_column_to_users"
```

## 🔧 Variables de Entorno

```bash
DB_HOST=localhost
DB_PORT=5432
DB_NAME=edugo_dev
DB_USER=edugo
DB_PASSWORD=changeme
DB_SSL_MODE=disable
```

## 📦 Importar en Proyectos

```go
import "github.com/EduGoGroup/edugo-infrastructure/postgres"
```

## 📚 Documentación

Ver documentación completa de ownership en: `../docs/TABLE_OWNERSHIP.md`

---

**Versión:** 0.5.0  
**Mantenedores:** Equipo EduGo
