# 🔌 Guía de Integración - edugo-infrastructure

**Cómo usar infrastructure desde otros proyectos**

---

## 📦 Instalación en Proyectos Consumidores

### api-admin

```bash
cd edugo-api-admin

# Agregar módulo database
go get github.com/EduGoGroup/edugo-infrastructure/database@latest

# go.mod quedará así:
require (
    github.com/EduGoGroup/edugo-infrastructure/database v0.1.0
    github.com/EduGoGroup/edugo-shared/auth v0.7.0
    github.com/EduGoGroup/edugo-shared/logger v0.7.0
    // ...
)
```

### api-mobile

```bash
cd edugo-api-mobile

# Agregar módulos database y schemas
go get github.com/EduGoGroup/edugo-infrastructure/database@latest
go get github.com/EduGoGroup/edugo-infrastructure/schemas@latest

require (
    github.com/EduGoGroup/edugo-infrastructure/database v0.1.0
    github.com/EduGoGroup/edugo-infrastructure/schemas v0.1.0
    github.com/EduGoGroup/edugo-shared/auth v0.7.0
    // ...
)
```

### worker

```bash
cd edugo-worker

# Agregar módulo schemas
go get github.com/EduGoGroup/edugo-infrastructure/schemas@latest

require (
    github.com/EduGoGroup/edugo-infrastructure/schemas v0.1.0
    github.com/EduGoGroup/edugo-shared/logger v0.7.0
    github.com/EduGoGroup/edugo-shared/messaging/rabbit v0.7.0
    // ...
)
```

---

## 🛠️ Uso en Makefile de Proyectos

### api-admin/Makefile

```makefile
INFRA_PATH := ../edugo-infrastructure

.PHONY: dev-setup
dev-setup: ## Setup de desarrollo para api-admin
	@echo "🚀 Setup de api-admin..."
	@cd $(INFRA_PATH) && make dev-up-core
	@echo "⏳ Esperando PostgreSQL..."
	@sleep 5
	@cd $(INFRA_PATH) && make migrate-up
	@cd $(INFRA_PATH) && make seed-minimal
	@echo "✅ Ambiente listo para api-admin"

.PHONY: dev-teardown
dev-teardown: ## Limpiar ambiente
	@cd $(INFRA_PATH) && make dev-teardown

.PHONY: run
run: ## Correr API
	@go run cmd/api/main.go
```

### api-mobile/Makefile

```makefile
INFRA_PATH := ../edugo-infrastructure

.PHONY: dev-setup
dev-setup: ## Setup para api-mobile (necesita RabbitMQ)
	@echo "🚀 Setup de api-mobile..."
	@cd $(INFRA_PATH) && make dev-up-messaging
	@sleep 5
	@cd $(INFRA_PATH) && make migrate-up
	@cd $(INFRA_PATH) && make seed
	@echo "✅ Ambiente listo para api-mobile"

.PHONY: run
run:
	@go run cmd/api/main.go
```

### worker/Makefile

```makefile
INFRA_PATH := ../edugo-infrastructure

.PHONY: dev-setup
dev-setup: ## Setup para worker
	@echo "🚀 Setup de worker..."
	@cd $(INFRA_PATH) && make dev-up-messaging
	@sleep 5
	@cd $(INFRA_PATH) && make migrate-up
	@cd $(INFRA_PATH) && make seed
	@echo "✅ Ambiente listo para worker"

.PHONY: run
run:
	@go run cmd/worker/main.go
```

---

## 📋 Uso de Schemas en Código

### Publisher (api-mobile)

```go
package messaging

import (
	"encoding/json"
	"github.com/EduGoGroup/edugo-infrastructure/schemas"
	"github.com/google/uuid"
	"time"
)

type MaterialPublisher struct {
	validator *schemas.EventValidator
	publisher *RabbitMQPublisher
}

func NewMaterialPublisher(publisher *RabbitMQPublisher) (*MaterialPublisher, error) {
	validator, err := schemas.NewEventValidator()
	if err != nil {
		return nil, err
	}
	
	return &MaterialPublisher{
		validator: validator,
		publisher: publisher,
	}, nil
}

func (p *MaterialPublisher) PublishMaterialUploaded(materialID, schoolID, teacherID uuid.UUID, fileURL string, fileSize int64, fileType string) error {
	event := map[string]interface{}{
		"event_id":      uuid.New().String(),
		"event_type":    "material.uploaded",
		"event_version": "1.0",
		"timestamp":     time.Now().UTC().Format(time.RFC3339),
		"payload": map[string]interface{}{
			"material_id":     materialID.String(),
			"school_id":       schoolID.String(),
			"teacher_id":      teacherID.String(),
			"file_url":        fileURL,
			"file_size_bytes": fileSize,
			"file_type":       fileType,
		},
	}

	// Validar ANTES de publicar
	if err := p.validator.Validate(event); err != nil {
		return fmt.Errorf("invalid event: %w", err)
	}

	eventJSON, _ := json.Marshal(event)
	
	return p.publisher.Publish("edugo.topic", "material.uploaded", eventJSON)
}
```

### Consumer (worker)

```go
package messaging

import (
	"encoding/json"
	"github.com/EduGoGroup/edugo-infrastructure/schemas"
)

type MaterialConsumer struct {
	validator *schemas.EventValidator
	processor *MaterialProcessor
}

func NewMaterialConsumer(processor *MaterialProcessor) (*MaterialConsumer, error) {
	validator, err := schemas.NewEventValidator()
	if err != nil {
		return nil, err
	}
	
	return &MaterialConsumer{
		validator: validator,
		processor: processor,
	}, nil
}

func (c *MaterialConsumer) HandleMessage(msg []byte) error {
	// Validar JSON contra schema
	if err := c.validator.ValidateJSON(msg, "material.uploaded", "1.0"); err != nil {
		logger.Error("invalid event received", "error", err)
		return c.sendToDLQ(msg, err)
	}

	// Deserializar evento validado
	var event MaterialUploadedEvent
	if err := json.Unmarshal(msg, &event); err != nil {
		return err
	}

	// Procesar evento
	return c.processor.ProcessMaterial(event.Payload.MaterialID, event.Payload.FileURL)
}
```

---

## 🔄 Workflow de Desarrollo

### Día 1: Setup Inicial

```bash
# 1. Clonar todos los repos
cd /Users/jhoanmedina/source/EduGo/repos-separados
git clone git@github.com:EduGoGroup/edugo-infrastructure.git
git clone git@github.com:EduGoGroup/edugo-api-admin.git
git clone git@github.com:EduGoGroup/edugo-api-mobile.git

# 2. Setup de infrastructure
cd edugo-infrastructure
cp .env.example .env
make dev-setup

# ✅ PostgreSQL, MongoDB corriendo con datos de prueba
```

### Día 2: Desarrollar api-admin

```bash
# 1. Actualizar dependencias
cd edugo-api-admin
go get github.com/EduGoGroup/edugo-infrastructure/database@v0.1.0
go mod tidy

# 2. Crear Makefile con referencia a infrastructure
# (ver ejemplo arriba)

# 3. Levantar ambiente
make dev-setup

# 4. Desarrollar
make run
```

### Día 3: Desarrollar api-mobile

```bash
cd edugo-api-mobile
go get github.com/EduGoGroup/edugo-infrastructure/database@v0.1.0
go get github.com/EduGoGroup/edugo-infrastructure/schemas@v0.1.0
go mod tidy

make dev-setup
make run
```

---

## 🐛 Troubleshooting

### Error: "schema_migrations table does not exist"

```bash
cd edugo-infrastructure
make migrate-up

# O manualmente:
cd database
go run migrate.go up
```

### Error: "connection refused" al conectar a PostgreSQL

```bash
# Verificar que Docker está corriendo
make dev-ps

# Si no está corriendo:
make dev-up-core
```

### Error: "invalid event" al publicar

```go
// Verificar que evento tiene todos los campos requeridos
// Consultar EVENT_CONTRACTS.md para estructura exacta
```

---

## 📚 Referencias

- **Migraciones:** `database/TABLE_OWNERSHIP.md`
- **Eventos:** `EVENT_CONTRACTS.md`
- **Docker:** `docker/README.md`
- **Schemas:** `schemas/README.md`

---

**Última actualización:** 15 de Noviembre, 2025
