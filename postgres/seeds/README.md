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
> `docker-recreate`/`cloud-migrate` es ahora `system` + `playground_v2/base`. Para un fixture autónomo:
> `make docker-playground-v2 P=<fixture>`. Los escenarios on-top (plan 053 T5.2) no van por el migrador:
> ver «Cómo ejecutar».

## Cómo ejecutar

```bash
cd /Users/jhoanmedina/source/EduGo/EduBack/edugo-dev-environment

# Recrear Docker completa (contrato system + playground_v2/base, sin flags)
make docker-recreate

# Recrear con un fixture autónomo (base, mp10_material, f6_representante)
make docker-playground-v2 P=mp10_material

# demo (HU-DEMO.5) es autónomo pero queda fuera de "all": solo va en la base de demo,
# recién recreada, con --playground-v2=demo (el envoltorio de la AR-DEMO.1). No trunca y
# se niega si encuentra otra escuela. No es el viejo seeds/demo/, que se eliminó.

# Los escenarios on-top (n4_evaluacion, onboarding, n1_inscripcion, n17_secciones,
# multi_unidad, n0n1_escuelas) NO crean usuarios desde el plan 053 T5.2: van encima
# de base, sus usuarios (playground_v2.ScenarioUsers) los crea antes la API de
# identidad y luego se aplica playground_v2.ApplyOnTop. Pedirlos por el migrador
# sobre una BD recién recreada falla listando los usuarios ausentes.

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
