# 🔴 Código Duplicado - Mejoras Pendientes

Código duplicado identificado que debe consolidarse para mejorar mantenibilidad.

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
| DUP-002 | getEnv() duplicado | 🟢 Baja | Aceptar |
| DUP-003 | sanitizeName() duplicado | 🟢 Baja | Aceptar |

---

**Última actualización:** Diciembre 2024
