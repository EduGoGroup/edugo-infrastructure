# ⚠️ DEPRECADO - Módulo Messaging

> **Este módulo está deprecado. Usar `schemas` en su lugar.**

## Migración

```go
// ❌ Antes (deprecado)
import "github.com/EduGoGroup/edugo-infrastructure/messaging"
validator, _ := messaging.NewEventValidator()

// ✅ Ahora
import "github.com/EduGoGroup/edugo-infrastructure/schemas"
validator, _ := schemas.NewEventValidator()
```

## Instalación del módulo correcto

```bash
go get github.com/EduGoGroup/edugo-infrastructure/schemas
```

## Documentación

Ver: `../schemas/README.md`

---

**Deprecado desde:** Diciembre 2024  
**Razón:** Código duplicado con módulo `schemas`
