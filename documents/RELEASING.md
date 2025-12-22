# 🏷️ Guía de Releases - EduGo Infrastructure

Esta guía documenta el proceso de versionado y creación de releases para los módulos de `edugo-infrastructure`.

---

## 📋 Tabla de Contenidos

- [Visión General](#-visión-general)
- [Estructura de Tags](#-estructura-de-tags)
- [Versionado Semántico](#-versionado-semántico)
- [Proceso de Release](#-proceso-de-release)
- [Comandos Útiles](#-comandos-útiles)
- [Ejemplos por Módulo](#-ejemplos-por-módulo)
- [Troubleshooting](#-troubleshooting)

---

## 🎯 Visión General

El proyecto `edugo-infrastructure` utiliza **versionado por módulo** con tags Git que siguen el patrón:

```
<módulo>/v<SemVer>
```

### Módulos Versionados

| Módulo | Descripción | Último Tag |
|--------|-------------|------------|
| `postgres` | Migraciones y CLI PostgreSQL | `postgres/v0.11.1` |
| `mongodb` | Migraciones y CLI MongoDB | `mongodb/v0.10.1` |
| `schemas` | Schemas compartidos | `schemas/v0.1.2` |
| `messaging` | Utilidades de mensajería | `messaging/v0.1.x` |

### ¿Por qué Versionado por Módulo?

✅ **Independencia**: Cada módulo puede evolucionar a su propio ritmo  
✅ **Precisión**: Los consumidores pueden fijar versiones específicas por módulo  
✅ **Go Modules**: Compatible con `go get <module>@<version>`  
✅ **Rollback**: Fácil volver a versiones anteriores por módulo

---

## 🏗️ Estructura de Tags

### Patrón de Tag

```
<módulo>/v<MAJOR>.<MINOR>.<PATCH>
```

### Ejemplos Reales

```bash
postgres/v0.11.1    # PostgreSQL módulo, versión 0.11.1
mongodb/v0.10.1     # MongoDB módulo, versión 0.10.1
schemas/v0.1.2      # Schemas módulo, versión 0.1.2
```

### Anatomía de un Tag

```
postgres/v0.11.1
│        │ │  │  └─ PATCH: Bug fixes, cambios menores
│        │ │  └──── MINOR: Nuevas features, retrocompatible
│        │ └─────── MAJOR: Cambios breaking
│        └────────── Prefijo semántico obligatorio
└─────────────────── Nombre del módulo
```

---

## 📊 Versionado Semántico

Seguimos [SemVer 2.0.0](https://semver.org/):

### MAJOR (X.0.0)

**Cuándo incrementar:**
- Cambios incompatibles en la API
- Cambios en estructura de migraciones que rompen compatibilidad
- Remoción de funcionalidades públicas

**Ejemplo:**
```go
// v0.11.1
func ApplyMigrations(db *sql.DB) error

// v1.0.0 (breaking change)
func ApplyMigrations(ctx context.Context, db *sql.DB) error
```

### MINOR (0.X.0)

**Cuándo incrementar:**
- Nuevas migraciones agregadas
- Nuevas funciones públicas
- Nuevas features retrocompatibles

**Ejemplo:**
```bash
# Nueva migración agregada
postgres/migrations/
  011_add_user_preferences.sql  # ← Nueva migración
```

### PATCH (0.0.X)

**Cuándo incrementar:**
- Bug fixes
- Mejoras de documentación
- Optimizaciones de rendimiento sin cambios de API

**Ejemplo:**
```go
// v0.11.0
func ValidateSchema(s string) error {
    return nil  // Bug: no valida nada
}

// v0.11.1 (patch)
func ValidateSchema(s string) error {
    if s == "" {
        return errors.New("schema vacío")
    }
    return nil
}
```

---

## 🚀 Proceso de Release

### 1. Verificar Estado del Código

```bash
# Asegurar que estás en la rama correcta
git checkout main
git pull origin main

# Verificar que no hay cambios sin commitear
git status
```

### 2. Ejecutar Tests

```bash
# PostgreSQL
cd postgres && go test ./... && cd ..

# MongoDB
cd mongodb && go test ./... && cd ..

# Schemas
cd schemas && go test ./... && cd ..
```

### 3. Determinar Nueva Versión

```bash
# Ver último tag del módulo
git tag -l "postgres/v*" | sort -V | tail -1
# Salida: postgres/v0.11.1

# Decidir nueva versión según tipo de cambio:
# - Breaking change → v1.0.0
# - Nueva feature   → v0.12.0
# - Bug fix         → v0.11.2
```

### 4. Crear Tag

```bash
# Crear tag anotado (recomendado)
git tag -a postgres/v0.12.0 -m "Release postgres v0.12.0

Nuevas features:
- Agregada migración 012_user_sessions
- Mejorado manejo de errores en CLI

Bug fixes:
- Corregido timeout en migraciones largas
"

# O tag ligero (simple)
git tag postgres/v0.12.0
```

### 5. Publicar Tag

```bash
# Publicar tag específico
git push origin postgres/v0.12.0

# O publicar todos los tags
git push origin --tags
```

### 6. Verificar Publicación

```bash
# Verificar que el tag existe remotamente
git ls-remote --tags origin | grep postgres

# Verificar que consumidores pueden usarlo
go get github.com/edugo/edugo-infrastructure/postgres@v0.12.0
```

### 7. Actualizar Documentación

Actualizar este archivo con el nuevo tag en la tabla de "Módulos Versionados".

---

## 💻 Comandos Útiles

### Listar Tags

```bash
# Todos los tags
git tag -l

# Tags de un módulo específico
git tag -l "postgres/v*"

# Tags ordenados por versión
git tag -l "postgres/v*" | sort -V

# Último tag de un módulo
git tag -l "postgres/v*" | sort -V | tail -1
```

### Crear Tags

```bash
# Tag anotado (recomendado para releases)
git tag -a <módulo>/v<version> -m "Mensaje"

# Tag ligero
git tag <módulo>/v<version>

# Tag en commit específico
git tag -a postgres/v0.12.0 abc123 -m "Mensaje"
```

### Eliminar Tags

```bash
# Eliminar tag local
git tag -d postgres/v0.12.0

# Eliminar tag remoto
git push origin --delete postgres/v0.12.0

# Eliminar ambos (local y remoto)
git tag -d postgres/v0.12.0 && git push origin --delete postgres/v0.12.0
```

### Ver Información de Tag

```bash
# Ver detalles de tag anotado
git show postgres/v0.11.1

# Ver commit asociado
git rev-list -n 1 postgres/v0.11.1

# Ver cambios desde último tag
git log postgres/v0.11.0..postgres/v0.11.1 --oneline
```

### Consumir Versiones

```bash
# En go.mod
go get github.com/edugo/edugo-infrastructure/postgres@v0.11.1
go get github.com/edugo/edugo-infrastructure/mongodb@v0.10.1

# Actualizar a última versión
go get github.com/edugo/edugo-infrastructure/postgres@latest

# Listar versiones disponibles
go list -m -versions github.com/edugo/edugo-infrastructure/postgres
```

---

## 📦 Ejemplos por Módulo

### PostgreSQL

```bash
# Escenario: Agregaste nueva migración 012_add_audit_logs.sql

# 1. Verificar cambios
cd postgres
go test ./...

# 2. Determinar versión (nueva migración = MINOR bump)
git tag -l "postgres/v*" | sort -V | tail -1
# Salida: postgres/v0.11.1
# Nueva versión: postgres/v0.12.0

# 3. Crear tag
git tag -a postgres/v0.12.0 -m "Release postgres v0.12.0

- Agregada migración 012: audit_logs table
- Mejoras en CLI de migraciones
"

# 4. Publicar
git push origin postgres/v0.12.0
```

### MongoDB

```bash
# Escenario: Corregiste bug en ApplySeeds()

# 1. Verificar cambios
cd mongodb
go test ./...

# 2. Determinar versión (bug fix = PATCH bump)
git tag -l "mongodb/v*" | sort -V | tail -1
# Salida: mongodb/v0.10.1
# Nueva versión: mongodb/v0.10.2

# 3. Crear tag
git tag -a mongodb/v0.10.2 -m "Release mongodb v0.10.2

Bug fixes:
- Corregido error en ApplySeeds() con colecciones vacías
"

# 4. Publicar
git push origin mongodb/v0.10.2
```

### Schemas

```bash
# Escenario: Agregaste nuevo schema user_preferences.go

# 1. Verificar cambios
cd schemas
go test ./...

# 2. Determinar versión (nueva feature = MINOR bump)
git tag -l "schemas/v*" | sort -V | tail -1
# Salida: schemas/v0.1.2
# Nueva versión: schemas/v0.2.0

# 3. Crear tag
git tag -a schemas/v0.2.0 -m "Release schemas v0.2.0

Features:
- Agregado schema UserPreferences
- Agregadas validaciones para preferencias de notificación
"

# 4. Publicar
git push origin schemas/v0.2.0
```

---

## 🔧 Troubleshooting

### Problema: Tag ya existe

```bash
# Error
fatal: tag 'postgres/v0.12.0' already exists

# Solución 1: Usar nueva versión
git tag postgres/v0.12.1

# Solución 2: Eliminar y recrear (CUIDADO en producción)
git tag -d postgres/v0.12.0
git push origin --delete postgres/v0.12.0
git tag -a postgres/v0.12.0 -m "Mensaje"
git push origin postgres/v0.12.0
```

### Problema: go get no encuentra versión

```bash
# Error
go: github.com/edugo/edugo-infrastructure/postgres@v0.12.0: invalid version: unknown revision

# Causas posibles:
# 1. Tag no está pusheado
git push origin postgres/v0.12.0

# 2. Proxy de Go no tiene la versión aún (esperar ~10 min)
GOPROXY=direct go get github.com/edugo/edugo-infrastructure/postgres@v0.12.0

# 3. Verificar que el tag existe remotamente
git ls-remote --tags origin | grep postgres
```

### Problema: Versión equivocada

```bash
# Publicaste postgres/v0.13.0 pero debió ser v0.12.1

# Solución:
# 1. Eliminar tag incorrecto
git tag -d postgres/v0.13.0
git push origin --delete postgres/v0.13.0

# 2. Crear tag correcto
git tag -a postgres/v0.12.1 -m "Mensaje"
git push origin postgres/v0.12.1

# 3. Informar a consumidores si ya se distribuyó
```

### Problema: go.mod no se actualiza

```bash
# go.mod sigue usando versión antigua

# Solución:
go get github.com/edugo/edugo-infrastructure/postgres@v0.12.0
go mod tidy
```

---

## 📚 Referencias

- [Go Modules Reference](https://go.dev/ref/mod)
- [Semantic Versioning 2.0.0](https://semver.org/)
- [Git Tagging](https://git-scm.com/book/en/v2/Git-Basics-Tagging)
- [Proyecto edugo-infrastructure](https://github.com/edugo/edugo-infrastructure)

---

## 🤝 Contribución

Para proponer cambios al proceso de versionado, abre un issue en el repositorio.

---

**Última actualización:** Diciembre 2025  
**Responsable:** Equipo de Infraestructura EduGo
