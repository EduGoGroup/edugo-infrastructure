# PostgreSQL Seeds

Datos iniciales del sistema y datos de prueba para desarrollo.

## Estructura

```
version.go             — ComputeFilesHash() dinámico desde system.Layers()

system/                — Datos obligatorios (roles, permisos, pantallas, tipos)
  layer.go             (interfaz Layer: Name, SeedVersion, Apply)
  system.go            (Layers(), ApplySystem(db, upTo))
  legacy/
    legacy.go          (Seed(gdb), SeedVersion="production-gorm-v1")
    data.go            (dataset crudo generado)
    accessors.go       (9 funciones read-only)
    legacy_layer.go    (NewLegacyLayer() implementa Layer)

demo/                  — Datos de prueba (escuelas, usuarios, materias, etc.)
  development.go       (ApplyDemo(gdb), SeedVersion)

e2e/                   — Framework Fase C de fixtures compositivas + scenarios
```

## Cómo ejecutar

```bash
cd /Users/jhoanmedina/source/EduGo/EduBack/edugo-dev-environment

# Recrear Docker completa (estructura + seeds)
make docker-recreate

# Aplicar hasta una capa específica del sistema
make docker-seed-layer LAYER=legacy

# Aplicar seeds de sistema en cloud
make cloud-seed-layer LAYER=legacy

# Aplicar scenario E2E
make cloud-seed-scenario SCENARIO=teacher_grades_only

# Ver version actual
make neon-status
```

## Archivos clave

- `system/layers/l*_*.go` — Capas L0..L4 del rebuild (resources, roles, permisos, pantallas).
- `system/l4/*.go` — Datos de L4 por dominio (sistema completo reorganizado post-Fase-6).
- `system/system.go` — `ApplySystem(db, upTo)`: aplica capas hasta `upTo` (vacío = todas).
- `demo/development.go` — Seeds de desarrollo (incluye usuarios de prueba; password: `12345678`).

## Version

Al modificar un seed:
1. Ajustar `SeedVersion` en el paquete afectado (`system/layers/lN_constants.go` o `demo/development.go`).
2. Incrementar `SchemaVersion` en `../migrations/version.go` si el cambio altera `ComputeFilesHash()` o cierra una fase del rebuild.

Ver reglas completas en `CLAUDE.md`.
