# database - Módulo de Migraciones

Centraliza todas las migraciones de PostgreSQL y MongoDB.

## Uso

```bash
cd database
go run migrate.go up      # Ejecutar migraciones
go run migrate.go down    # Revertir última
go run migrate.go status  # Ver estado
```

## Crear Nueva Migración

```bash
go run migrate.go create "add_avatar_to_users"
```

Genera:
- `migrations/postgres/00X_add_avatar_to_users.up.sql`
- `migrations/postgres/00X_add_avatar_to_users.down.sql`
