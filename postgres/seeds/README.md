# PostgreSQL Seeds

Datos iniciales del sistema y datos de prueba para desarrollo.

## Estructura

```
version.go             — ComputeFilesHash() dinámico desde system.Layers() (solo el CONTRATO)

system/                — CONTRATO PURO: roles, permisos, pantallas, tipos (L0–L4). SIEMPRE se aplica.
  layer.go             (interfaz Layer: Name, SeedVersion, Apply)
  system.go            (Layers(), ApplySystem(db, upTo))
  layers/              (capas L0..L4 del rebuild)
  l4/                  (datos de L4 por dominio)

playground_v2/         — ÚNICO mundo de datos (fotos inmutables, componibles). NO entra al hash.
  playground_v2.go     (registry: base, n4_evaluacion, onboarding, n1_inscripcion, …)
  base/                (base.Apply: fixture por DEFECTO — 2 escuelas, 9 usuarios @edugo.test, 12345678)
  common/              (helpers compartidos Seed*)

e2e/                   — Framework Fase C de fixtures compositivas + scenarios
```

> **MP-09 (2026-06-14):** se eliminaron `seeds/demo/` y `seeds/playground/` (v1). El default de
> `docker-recreate`/`cloud-migrate` es ahora `system` + `playground_v2/base`. Para un dataset focalizado:
> `make docker-playground-v2 P=<fixture>`.

## Cómo ejecutar

```bash
cd /Users/jhoanmedina/source/EduGo/EduBack/edugo-dev-environment

# Recrear Docker completa (contrato system + playground_v2/base, sin flags)
make docker-recreate

# Recrear con un fixture focalizado
make docker-playground-v2 P=n4_evaluacion

# Aplicar hasta una capa específica del contrato
make docker-seed-layer LAYER=l4

# Aplicar scenario E2E
make cloud-seed-scenario SCENARIO=teacher_grades_only

# Ver version actual
make neon-status
```

## Archivos clave

- `system/layers/l*_*.go` — Capas L0..L4 del contrato (resources, roles, permisos, pantallas).
- `system/l4/*.go` — Datos de L4 por dominio (sistema completo reorganizado post-Fase-6).
- `system/system.go` — `ApplySystem(db, upTo)`: aplica capas hasta `upTo` (vacío = todas).
- `playground_v2/base/base.go` — Fixture por defecto (incluye usuarios de prueba; password: `12345678`).

## Version

Al modificar un seed del **contrato** (`system/`):
1. Ajustar `SeedVersion` en la capa afectada (`system/layers/lN_constants.go`).
2. Incrementar `SchemaVersion` en `../migrations/version.go` si el cambio altera `ComputeFilesHash()` o cierra una fase del rebuild.

Los fixtures de `playground_v2/` (incl. `base`) **no entran al hash** → no exigen bump.

Ver reglas completas en `CLAUDE.md`.
