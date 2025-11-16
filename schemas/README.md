# schemas - Módulo de JSON Schemas

Validación automática de eventos RabbitMQ usando JSON Schema.

## Schemas Disponibles

| Evento | Versión | Publicado por | Consumido por |
|--------|---------|---------------|---------------|
| `material.uploaded` | v1.0 | api-mobile | worker |
| `assessment.generated` | v1.0 | worker | api-mobile |
| `material.deleted` | v1.0 | api-mobile | worker |
| `student.enrolled` | v1.0 | api-admin | api-mobile |

## Uso en Publisher

```go
import "github.com/EduGoGroup/edugo-infrastructure/schemas"

event := MaterialUploadedEvent{
    EventID:      uuid.New(),
    EventType:    "material.uploaded",
    EventVersion: "1.0",
    Payload:      payload,
}

validator := schemas.NewEventValidator()
if err := validator.Validate(event); err != nil {
    return fmt.Errorf("invalid event: %w", err)
}

publisher.Publish(event)  // ✅ Validado
```

## Uso en Consumer

```go
validator := schemas.NewEventValidator()

if err := validator.ValidateJSON(msg, "material.uploaded", "1.0"); err != nil {
    logger.Error("invalid event", err)
    return sendToDLQ(msg, err)
}

// Procesar evento validado
```

## Versionamiento

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
