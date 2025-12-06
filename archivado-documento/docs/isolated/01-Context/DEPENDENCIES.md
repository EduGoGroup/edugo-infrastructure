# DEPENDENCIES - Infrastructure

## Matriz de Dependencias

```
┌──────────────────────────────────────────────────────────┐
│              INFRASTRUCTURE (v0.1.1)                     │
├──────────────────────────────────────────────────────────┤
│ Dependencias Críticas (runtime)                          │
│ ├─ PostgreSQL 15+                                        │
│ ├─ MongoDB 7.0+                                          │
│ ├─ RabbitMQ 3.12+                                        │
│ └─ Redis 7.0+                                            │
│                                                          │
│ Dependencias de Desarrollo                               │
│ ├─ Docker 20.10+                                         │
│ ├─ Docker Compose 2.0+                                   │
│ └─ Go 1.21+ (solo para migrate.go y validator.go)       │
│                                                          │
│ NO Depende de (Importante)                               │
│ ├─ ❌ edugo-shared                                       │
│ ├─ ❌ edugo-api-mobile                                   │
│ ├─ ❌ edugo-api-administracion                           │
│ └─ ❌ edugo-worker                                       │
└──────────────────────────────────────────────────────────┘
```

---

## Dependencias Críticas

### 1. PostgreSQL 15+

**¿Por qué?** Base de datos relacional principal del ecosistema

**Versión mínima:** 15.0  
**Recomendada:** 15.4+

**Características requeridas:**
- Soporte de CTEs (Common Table Expressions)
- JSON/JSONB columns
- Foreign keys y constraints
- Triggers

**Configuración:**
```yaml
# En docker-compose.yml
services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: edugo_user
      POSTGRES_PASSWORD: secure_password
      POSTGRES_DB: edugo_dev
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
```

**Uso:**
- `database/migrations/*.sql` se ejecutan aquí
- Todos los servicios se conectan a esta instancia

---

### 2. MongoDB 7.0+

**¿Por qué?** Base de datos documental para datos no relacionales

**Versión mínima:** 7.0  
**Recomendada:** 7.0.4+

**Características requeridas:**
- Replica set (opcional en dev)
- Índices compuestos
- Agregaciones

**Configuración:**
```yaml
# En docker-compose.yml
services:
  mongo:
    image: mongo:7.0
    environment:
      MONGO_INITDB_ROOT_USERNAME: root
      MONGO_INITDB_ROOT_PASSWORD: secure_password
    ports:
      - "27017:27017"
    volumes:
      - mongo_data:/data/db
```

**Uso:**
- Worker guarda resúmenes y quizzes
- Api-mobile lee evaluaciones generadas

---

### 3. RabbitMQ 3.12+

**¿Por qué?** Message broker para comunicación asíncrona

**Versión mínima:** 3.12  
**Recomendada:** 3.12.10+

**Características requeridas:**
- Management plugin
- Topic exchanges
- Quorum queues (opcional)

**Configuración:**
```yaml
# En docker-compose.yml
services:
  rabbitmq:
    image: rabbitmq:3.12-management
    environment:
      RABBITMQ_DEFAULT_USER: guest
      RABBITMQ_DEFAULT_PASS: guest
    ports:
      - "5672:5672"     # AMQP
      - "15672:15672"   # Management UI
    volumes:
      - rabbitmq_data:/var/lib/rabbitmq
```

**Uso:**
- `schemas/events/*.json` validan mensajes publicados/consumidos
- Exchanges y queues configurados en startup

---

### 4. Redis 7.0+

**¿Por qué?** Cache y sesiones

**Versión mínima:** 7.0  
**Recomendada:** 7.0.14+

**Características requeridas:**
- Persistencia AOF o RDB
- Pub/Sub (opcional)

**Configuración:**
```yaml
# En docker-compose.yml
services:
  redis:
    image: redis:7.0-alpine
    command: redis-server --requirepass secure_password
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
```

**Uso:**
- Cache de queries frecuentes
- Sesiones de usuario (futuro)

---

## Dependencias de Desarrollo

### 5. Docker 20.10+

**¿Por qué?** Containerización de servicios

**Instalación:**
```bash
# Mac
brew install --cask docker

# Windows
# Descargar Docker Desktop

# Linux
sudo apt-get install docker-ce docker-ce-cli containerd.io
```

**Verificación:**
```bash
docker --version
# Docker version 20.10.21, build baeda1f
```

---

### 6. Docker Compose 2.0+

**¿Por qué?** Orquestación multi-contenedor

**Nota:** Incluido en Docker Desktop (Mac/Windows)

**Verificación:**
```bash
docker-compose --version
# Docker Compose version v2.12.2
```

---

### 7. Go 1.21+ (Opcional)

**¿Por qué?** Solo para `migrate.go` y `validator.go`

**Nota:** Infraestructura funciona sin Go (SQL directo, Docker)

**Uso:**
```bash
# CLI de migraciones (cuando esté implementado)
go run database/migrate.go up

# Validador de eventos (cuando esté implementado)
go run schemas/validator.go validate material.uploaded event.json
```

---

## Consumidores de Infrastructure

### 1. edugo-api-mobile (v0.4.0+)

**Qué usa:**
- `database/migrations/003_create_materials.up.sql`
- `database/migrations/004_create_assessment.up.sql`
- `schemas/events/material.uploaded.json`
- `schemas/events/evaluation.submitted.json`

**Integración:**
```go
// En api-mobile
import "github.com/EduGoGroup/edugo-infrastructure/schemas"

// Validar evento antes de publicar
func PublishMaterialEvent(material Material) error {
    event := buildMaterialEvent(material)
    if err := schemas.Validate("material.uploaded", event); err != nil {
        return err
    }
    return publisher.Publish(event)
}
```

---

### 2. edugo-api-administracion (v0.2.0)

**Qué usa:**
- `database/migrations/001_create_users.up.sql`
- `database/migrations/002_create_schools.up.sql`
- `database/migrations/005_create_academic_hierarchy.up.sql`
- `database/TABLE_OWNERSHIP.md` (documentación)

**Integración:**
```bash
# Aplicar migraciones
cd infrastructure/database/migrations
psql -U edugo_user -d edugo_admin < 001_create_users.up.sql
psql -U edugo_user -d edugo_admin < 002_create_schools.up.sql
psql -U edugo_user -d edugo_admin < 005_create_academic_hierarchy.up.sql
```

---

### 3. edugo-worker (v0.3.0+)

**Qué usa:**
- `schemas/events/material.uploaded.json` (consume)
- `schemas/events/assessment.generated.json` (publica)
- `schemas/events/summary.completed.json` (publica)

**Integración:**
```go
// En worker
import "github.com/EduGoGroup/edugo-infrastructure/schemas"

// Validar evento recibido
func ConsumeMaterialEvent(msg []byte) error {
    if err := schemas.Validate("material.uploaded", msg); err != nil {
        logger.Error("Invalid event", err)
        return err
    }
    // Procesar...
}
```

---

### 4. edugo-dev-environment (v1.0.0)

**Qué usa:**
- `docker/docker-compose.yml` (puede copiar o referenciar)
- `scripts/init-db.sh`
- `scripts/seed-data.sh`

**Integración:**
```bash
# Opción 1: Copiar docker-compose.yml
cp infrastructure/docker/docker-compose.yml dev-environment/

# Opción 2: Referenciar directamente
cd infrastructure/docker
docker-compose --profile core up -d
```

---

### 5. edugo-shared (v0.7.0)

**Qué usa:**
- ❌ **NO depende** de infrastructure
- Infrastructure es más bajo nivel

**Nota:** shared/evaluation usa modelos que coinciden con schemas/events/, pero no importa directamente.

---

## Dependencias Inversas (Quién Necesita Qué)

### Tabla `users`
- **Dueño:** api-administracion
- **Lectores:** api-mobile, worker
- **Migración:** `001_create_users.up.sql`

### Tabla `schools`
- **Dueño:** api-administracion
- **Lectores:** api-mobile, worker
- **Migración:** `002_create_schools.up.sql`

### Tabla `materials`
- **Dueño:** api-mobile
- **Lectores:** worker
- **Migración:** `003_create_materials.up.sql`

### Tabla `assessment`
- **Dueño:** api-mobile
- **Lectores:** -
- **Migración:** `004_create_assessment.up.sql`

### Tabla `academic_units`
- **Dueño:** api-administracion
- **Lectores:** api-mobile
- **Migración:** `005_create_academic_hierarchy.up.sql`

---

## Matriz de Compatibilidad

| Componente | Versión Mínima | Versión Máxima | Notas |
|-----------|----------------|----------------|-------|
| PostgreSQL | 15.0 | 16.x | Compatible |
| MongoDB | 7.0 | 8.x | Compatible |
| RabbitMQ | 3.12 | 3.13.x | Compatible |
| Redis | 7.0 | 7.2.x | Compatible |
| Docker | 20.10 | Latest | Compatible |
| Docker Compose | 2.0 | Latest | Compatible |
| Go (opcional) | 1.21 | 1.22 | Compatible |

---

## Checklist de Instalación

```markdown
## Servicios Docker
- [ ] Docker instalado (>= 20.10)
- [ ] Docker Compose instalado (>= 2.0)
- [ ] Puertos disponibles (5432, 27017, 5672, 6379)
- [ ] Espacio en disco (>= 10GB)

## Configuración
- [ ] Archivo .env creado desde .env.example
- [ ] Credenciales configuradas
- [ ] Networking configurado (edugo-network)

## Servicios Levantados
- [ ] PostgreSQL respondiendo en 5432
- [ ] MongoDB respondiendo en 27017
- [ ] RabbitMQ respondiendo en 5672 y 15672
- [ ] Redis respondiendo en 6379

## Validación
- [ ] docker-compose ps (todos "healthy")
- [ ] ./scripts/health-check.sh (éxito)
- [ ] Conexión a PostgreSQL exitosa
- [ ] Conexión a MongoDB exitosa
```

---

## Resolución de Problemas

### Puerto 5432 ocupado
```bash
# Ver qué proceso usa el puerto
lsof -i :5432

# Matar proceso o cambiar puerto en docker-compose.yml
```

### Docker sin permisos
```bash
# Linux: agregar usuario a grupo docker
sudo usermod -aG docker $USER
newgrp docker
```

### Servicios no pasan health check
```bash
# Ver logs detallados
docker-compose logs postgres
docker-compose logs mongo
docker-compose logs rabbitmq
```

---

**Última actualización:** 16 de Noviembre, 2025  
**Versión infrastructure:** v0.1.1
