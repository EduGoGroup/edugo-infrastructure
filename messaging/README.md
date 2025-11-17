# Módulo Messaging - edugo-infrastructure

Validación automática de eventos RabbitMQ usando JSON Schema.

## 🎯 Propósito

Proporcionar validación centralizada de eventos de mensajería (RabbitMQ) con JSON Schemas versionados.

## 📋 Schemas Disponibles

| Evento | Versión | Publicado por | Consumido por |
|--------|---------|---------------|---------------|
| `material.uploaded` | v1.0 | api-mobile | worker |
| `assessment.generated` | v1.0 | worker | api-mobile |
| `material.deleted` | v1.0 | api-mobile | worker |
| `student.enrolled` | v1.0 | api-admin | api-mobile |

## 🚀 Uso

### Publisher

```go
import "github.com/EduGoGroup/edugo-infrastructure/messaging"

event := MaterialUploadedEvent{
    EventID:      uuid.New(),
    EventType:    "material.uploaded",
    EventVersion: "1.0",
    Payload:      payload,
}

validator := messaging.NewEventValidator()
if err := validator.Validate(event); err != nil {
    return fmt.Errorf("invalid event: %w", err)
}

publisher.Publish(event)  // ✅ Validado
```

### Consumer

```go
validator := messaging.NewEventValidator()

if err := validator.ValidateJSON(msg, "material.uploaded", "1.0"); err != nil {
    logger.Error("invalid event", err)
    return sendToDLQ(msg, err)
}

// Procesar evento validado
```

## 📦 Instalación

```bash
go get github.com/EduGoGroup/edugo-infrastructure/messaging
```

## 🔄 Versionamiento

- **Minor change (1.0 → 1.1):** Agregar campos opcionales
- **Major change (1.0 → 2.0):** Breaking changes

Consumer debe manejar múltiples versiones:

```go
switch event.EventVersion {
case "1.0", "1.1":
    return handleV1(event)
case "2.0":
    return handleV2(event)
default:
    return fmt.Errorf("unsupported version: %s", event.EventVersion)
}
```

## 📚 Documentación

Ver contratos completos en: `../EVENT_CONTRACTS.md`

---

**Versión:** 0.5.0  
**Mantenedores:** Equipo EduGo
