# Guía de Migración: edugo-infrastructure v0.8.0

**Para:** edugo-api-mobile  
**Fecha:** 2025-11-18  
**Prioridad:** ALTA  
**Breaking Changes:** ✅ SÍ

---

## 🎯 RESUMEN EJECUTIVO

edugo-infrastructure v0.8.0 eliminó el módulo `migrations/` y lo movió a `postgres/testing/`. 

**Acción requerida:** Actualizar imports en tus tests de integración.

---

## 🚨 BREAKING CHANGES

### Cambio 1: Módulo `migrations/` eliminado

**ANTES:**
```
github.com/EduGoGroup/edugo-infrastructure/migrations
```

**AHORA:**
```
github.com/EduGoGroup/edugo-infrastructure/postgres/testing
```

---

## 🔧 PASOS DE MIGRACIÓN

### Paso 1: Actualizar Import

**Archivo:** `internal/testing/suite/integration_suite.go`

**ANTES:**
```go
import (
    infrastructureTesting "github.com/EduGoGroup/edugo-infrastructure/migrations"
)
```

**DESPUÉS:**
```go
import (
    pgtesting "github.com/EduGoGroup/edugo-infrastructure/postgres/testing"
)
```

---

### Paso 2: Actualizar Llamadas a Funciones

**Buscar y reemplazar en todo el archivo:**

| ANTES | DESPUÉS |
|-------|---------|
| `infrastructureTesting.CleanDatabase(` | `pgtesting.CleanDatabase(` |
| `infrastructureTesting.ApplySeeds(` | `pgtesting.ApplySeeds(` |
| `infrastructureTesting.ApplyMigrations(` | `pgtesting.ApplyMigrations(` |

**Ejemplo:**

**ANTES:**
```go
func (s *IntegrationTestSuite) SetupTest() {
    if err := infrastructureTesting.CleanDatabase(s.PostgresDB); err != nil {
        s.T().Fatalf("Error limpiando BD: %v", err)
    }
    
    if err := infrastructureTesting.ApplySeeds(s.PostgresDB, s.seedsPath); err != nil {
        s.T().Fatalf("Error aplicando seeds: %v", err)
    }
}
```

**DESPUÉS:**
```go
func (s *IntegrationTestSuite) SetupTest() {
    if err := pgtesting.CleanDatabase(s.PostgresDB); err != nil {
        s.T().Fatalf("Error limpiando BD: %v", err)
    }
    
    if err := pgtesting.ApplySeeds(s.PostgresDB, s.seedsPath); err != nil {
        s.T().Fatalf("Error aplicando seeds: %v", err)
    }
}
```

---

### Paso 3: Actualizar go.mod

```bash
cd edugo-api-mobile

# Actualizar postgres a v0.8.0
go get github.com/EduGoGroup/edugo-infrastructure/postgres@v0.8.0

# Limpiar dependencias obsoletas (migrations/ se eliminará automáticamente)
go mod tidy
```

---

### Paso 4: Verificar Compilación

```bash
# Compilar proyecto
go build ./...

# Ejecutar tests unitarios
go test ./... -short

# Ejecutar tests de integración
make test-integration
```

**Resultado esperado:** Todo debe compilar y los tests deben pasar sin cambios.

---

### Paso 5: Commit

```bash
git add .
git commit -m "chore: actualizar a edugo-infrastructure postgres/v0.8.0

- Cambiar import: migrations → postgres/testing
- Actualizar referencias a helpers de testing
- go mod tidy

Relacionado: edugo-infrastructure v0.8.0 (simplificación de módulos)"

git push origin <tu-rama>
```

---

## ✅ CHECKLIST

- [ ] Import actualizado en `integration_suite.go`
- [ ] Todas las llamadas a `infrastructureTesting.` cambiadas a `pgtesting.`
- [ ] `go get postgres@v0.8.0` ejecutado
- [ ] `go mod tidy` ejecutado
- [ ] `go build ./...` exitoso
- [ ] Tests unitarios: PASS
- [ ] Tests de integración: PASS
- [ ] Commit y push realizados

---

## ❓ FAQ

### ¿Cambiaron las funciones?
No, las funciones son idénticas. Solo cambió el import path.

### ¿Hay nuevas funcionalidades?
No en postgres/. Las nuevas collections son de MongoDB para worker.

### ¿Cuánto tiempo tomará?
~15-20 minutos (cambio simple de imports)

---

## 📞 SOPORTE

Si encuentras problemas:
1. Verifica que usaste `postgres@v0.8.0` (no `migrations@...`)
2. Ejecuta `go mod tidy` para limpiar dependencias
3. Verifica que el import es `postgres/testing` (no `postgres`)

---

**Generado por:** edugo-infrastructure  
**Versión:** v0.8.0  
**Fecha:** 2025-11-18
