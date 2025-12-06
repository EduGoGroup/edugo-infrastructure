# 🔴 Código Duplicado - Mejoras Urgentes

Código duplicado identificado que debe consolidarse para mejorar mantenibilidad.

---

## DUP-001: validator.go Duplicado

### Descripción

Los archivos `schemas/validator.go` y `messaging/validator.go` son **100% idénticos**.

### Ubicación

| Archivo | Líneas | Tamaño |
|---------|--------|--------|
| `schemas/validator.go` | 139 | 4092 bytes |
| `messaging/validator.go` | 139 | 4092 bytes |

### Código Duplicado

```go
// Ambos archivos contienen exactamente el mismo código:
package schemas

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

//go:embed events/*.json
var schemasFS embed.FS

type EventValidator struct {
	schemas map[string]*gojsonschema.Schema
}

func NewEventValidator() (*EventValidator, error) {
	// ... implementación idéntica
}

func (v *EventValidator) Validate(event interface{}) error {
	// ... implementación idéntica
}

func (v *EventValidator) ValidateJSON(jsonBytes []byte, eventType, eventVersion string) error {
	// ... implementación idéntica
}
```

### Problema

1. **Mantenibilidad:** Cambios deben hacerse en dos lugares
2. **Confusión:** ¿Cuál usar? `schemas` o `messaging`?
3. **Inconsistencia potencial:** Podrían divergir con el tiempo
4. **Tamaño del módulo:** Código innecesario duplicado

### Solución Propuesta

**Opción A: Eliminar `messaging/validator.go`** (Recomendada)

```bash
# El módulo messaging debe importar desde schemas
rm messaging/validator.go
```

Actualizar `messaging/go.mod`:
```go
require github.com/EduGoGroup/edugo-infrastructure/schemas v0.1.0
```

Actualizar imports en proyectos:
```go
// Antes
import "github.com/EduGoGroup/edugo-infrastructure/messaging"
validator := messaging.NewEventValidator()

// Después
import "github.com/EduGoGroup/edugo-infrastructure/schemas"
validator := schemas.NewEventValidator()
```

**Opción B: Crear módulo compartido `validation`**

```
edugo-infrastructure/
├── validation/           # Nuevo módulo
│   ├── validator.go
│   ├── events/          # Schemas JSON
│   └── go.mod
├── schemas/             # Re-exporta desde validation
└── messaging/           # Re-exporta desde validation
```

### Impacto

| Aspecto | Antes | Después |
|---------|-------|---------|
| Archivos | 2 | 1 |
| Líneas de código | 278 | 139 |
| Puntos de cambio | 2 | 1 |

### Esfuerzo Estimado

- **Complejidad:** Baja
- **Tiempo:** 1-2 horas
- **Riesgo:** Bajo (solo cambio de imports)

### Pasos de Implementación

1. [ ] Verificar que ningún proyecto use `messaging.NewEventValidator()`
2. [ ] Si lo usan, agregar alias de re-exportación temporal
3. [ ] Eliminar `messaging/validator.go`
4. [ ] Actualizar documentación
5. [ ] Actualizar tests
6. [ ] Release nueva versión de módulos

---

## DUP-002: Funciones getEnv() Duplicadas

### Descripción

La función `getEnv()` está duplicada en múltiples archivos CLI.

### Ubicación

| Archivo | Línea |
|---------|-------|
| `postgres/cmd/migrate/migrate.go` | 118-123 |
| `mongodb/cmd/migrate/migrate.go` | 130-135 |

### Código Duplicado

```go
// Misma implementación en ambos archivos
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
```

### Problema

- Función trivial pero duplicada
- Si se necesita cambiar comportamiento (ej: logging), hay que hacerlo en múltiples lugares

### Solución Propuesta

**Opción A: Crear paquete `internal/config`**

```go
// internal/config/env.go
package config

import "os"

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func MustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic("required env var not set: " + key)
	}
	return value
}
```

**Opción B: Aceptar duplicación** (pragmático)

Dado que:
- Es código trivial (6 líneas)
- Está en archivos `main` de CLI
- No afecta APIs públicas

Se puede aceptar la duplicación como costo menor que la complejidad de extraer.

### Recomendación

**Aceptar duplicación por ahora** - El costo de mantener supera el beneficio de extraer para código tan simple en archivos CLI standalone.

### Esfuerzo Estimado

- **Complejidad:** Muy baja
- **Tiempo:** 30 minutos
- **Prioridad:** Baja

---

## DUP-003: Función sanitizeName() Duplicada

### Descripción

La función `sanitizeName()` para limpiar nombres de migraciones está duplicada.

### Ubicación

| Archivo | Línea |
|---------|-------|
| `postgres/cmd/migrate/migrate.go` | 439-452 |
| `mongodb/cmd/migrate/migrate.go` | 508-521 |

### Código Duplicado

```go
func sanitizeName(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")

	var result strings.Builder
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			result.WriteRune(r)
		}
	}

	return result.String()
}
```

### Solución Propuesta

Similar a DUP-002, se puede:
1. Extraer a paquete compartido
2. Aceptar duplicación en CLIs standalone

### Recomendación

**Aceptar duplicación** - Mismo razonamiento que DUP-002.

---

## 📊 Resumen de Acciones

| ID | Descripción | Prioridad | Acción |
|----|-------------|-----------|--------|
| DUP-001 | validator.go duplicado | 🔴 Alta | Eliminar messaging/validator.go |
| DUP-002 | getEnv() duplicado | 🟢 Baja | Aceptar |
| DUP-003 | sanitizeName() duplicado | 🟢 Baja | Aceptar |

---

## ✅ Resueltos

| Fecha | ID | Descripción | PR |
|-------|-----|-------------|-----|
| - | - | - | - |

---

**Última actualización:** Diciembre 2024
