# PostgreSQL Seeds

Datos iniciales del sistema y datos de prueba para desarrollo.

## Estructura

```
production/         — Datos obligatorios (roles, permisos, pantallas, tipos)
  001-008_*.sql

development/        — Datos de prueba (escuelas, usuarios, materias, etc.)
  000-013_*.sql
```

## Como ejecutar

```bash
cd /Users/jhoanmedina/source/EduGo/EduBack/edugo-dev-environment

# Recrear Neon completa (estructura + seeds)
make neon-recreate

# Ver version actual
make neon-status
```

## Archivos clave

- `production/006_ui_screen_instances.sql` — Configuracion de TODAS las pantallas del sistema
- `production/008_concept_types.sql` — Tipos de institucion y terminologia
- `development/003_users.sql` — Usuarios de prueba (password: 12345678)

## Version

Al modificar cualquier seed, incrementar `SchemaVersion` en `../migrations/version.go`.
