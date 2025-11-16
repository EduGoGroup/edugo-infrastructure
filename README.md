# edugo-infrastructure

Infraestructura compartida del ecosistema EduGo.

## Propósito

Centraliza:
- Migraciones de bases de datos (PostgreSQL + MongoDB)
- Docker Compose con perfiles
- JSON Schemas para validación de eventos
- Scripts de setup y seeds

## Estructura

```
edugo-infrastructure/
├── database/          # Módulo: Migraciones
├── docker/            # Módulo: Docker Compose
├── schemas/           # Módulo: JSON Schemas
├── scripts/           # Scripts de utilidades
├── seeds/             # Datos de prueba
└── Makefile
```

## Quick Start

```bash
make dev-setup    # Setup completo
make dev-up-core  # Solo PostgreSQL + MongoDB
```

Documentación completa próximamente.
