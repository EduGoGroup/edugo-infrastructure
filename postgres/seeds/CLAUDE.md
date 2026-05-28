# Reglas para modificar seeds

## SchemaVersion — cuándo hacer bump

Al cierre de cada fase del rebuild (Fase 2 en adelante), incrementar `SchemaVersion` en `../migrations/version.go`.

Fuera del rebuild, cualquier cambio que altere el output de `ComputeFilesHash()` debe acompañarse de bump:
- Renombre de capa (campo `Name()` en un `Layer`)
- Bump de `SeedVersion` en cualquier capa del sistema o en `demo`
- Agregado o eliminación de capas en `system.Layers()`
- Cambio en `demo.SeedVersion`

**Regla de e2e**: cualquier cambio bajo `e2e/framework/`, `e2e/fixtures/`, `e2e/scenarios/` o `e2e/exports/` también requiere bump — el framework forma parte del contrato que valida el migrator.

## Estructura

### system/ — Datos del sistema (SIEMPRE se aplican)

Datos que deben existir en cualquier ambiente: roles, permisos, resources, screen templates, screen instances, concept types.

Implementado por capas que cumplen la interfaz `Layer` (métodos `Name()`, `SeedVersion()`, `Apply()`).

- `system/layer.go` — interfaz `Layer`.
- `system/system.go` — `Layers()`, `ApplySystem(db, upTo)`. Post-Fase-2: `Layers()` retorna `[]Layer{layers.NewL0()}` (legacy desactivado por ADR-6).
- `system/layers/` — capas L0..L4 del rebuild.
  - `l0_*.go` — capa L0 (mínimo viable: ~17 filas — recurso `announcements` + rol `super_admin` + permisos CRUD + 3 templates + 1 ScreenInstance + user bootstrap).
  - `l1_*.go` — capa L1 readonly (rol `announcement_viewer` + escuela mínima + membership).
  - `l2_*.go` — capa L2 (segunda pantalla `announcement-form`).
  - `l3_*.go` — capa L3 (recurso `materials` con CRUD parcial sin delete + 2 pantallas).
  - `l4_full.go` + `l4_constants.go` — capa L4 (sistema completo, datos por dominio en `system/l4/`).
  - `l*_apply_twice_integration_test.go` — tests integration de idempotencia con testcontainer postgres.
- `system/l4/` — sub-paquete con los datos de L4 por dominio (resources, roles_permissions, screen_templates, screen_instances, resource_screens, concept_types) + accessors públicos para que el cross-checker los consuma.

### demo/ — Datos de prueba (solo en desarrollo)

Datos para testing: escuelas, usuarios, unidades académicas, materias, etc.

- `demo/development.go` — `ApplyDemo(gdb)`, `SeedVersion`.

### e2e/ — Fixtures focalizables (plan E2E + system-data-quality, Fase C)

A partir de Fase C el seed E2E pasa de ser una pila aditiva (`fase0..fase4`) a un framework de fixtures compositivas:

- `e2e/framework/` — interfaces, composer, cleanup selectivo, registry, helpers de upsert seguro y logger estructurado.
- `e2e/fixtures/` — piezas atómicas (role_only, screen_only, readonly_role, partial_crud, menu_subtree, guardian_relation, l0_constants_export).
- `e2e/scenarios/` — recetas canónicas (observer_audits, teacher_grades_only, guardian_views_child, l0_minimal).
- `e2e/exports/` — JSON `fixtures-constants.json` (artefacto de build, no se commitea) consumido por los tests Kotlin del KMP.

> **ADR-6 (Fase 2)**: las fixtures `legacy_*.go`, el scenario `scenario_legacy_e2e.go`, el paquete shim `e2e/e2e.go` + `fase{2,3,4}_*.go` y el test `legacy_compat_integration_test.go` fueron eliminados. Los scenarios `observer_audits`, `teacher_grades_only` y `guardian_views_child` quedan **skip-eados temporalmente** en el integration test porque dependen de permisos del catálogo legacy (`audit:read`, `grades:create`, etc.) que sólo se reintroducen en L4. Reactivar al cierre de Fase 6.

## Cambio pequeño vs fuerte

- **Pequeño**: editar datos en `system/layers/l*_*.go` o `system/l4/*.go`, ajustar `SeedVersion` (o `L4_SEED_VERSION`) de la capa afectada.
- **Fuerte**: refactor de capa, agregar nueva capa, cambiar interfaz `Layer`, cambiar `Layers()` — requiere bump `SchemaVersion` y recrear con `make docker-recreate`.

## Regla de pantallas (post-Fase-6)

La fuente de verdad de pantallas cambia según la capa:

- **L0** (`system/layers/l0_*.go`): 3 templates base (list/detail/form-basic-v1 con `definition.zones` canónico) + pantalla `announcements-list`. Si el render dinámico de KMP rompe en cualquier pantalla list/detail/form, revisar primero `l0_screens.go` (ahí está el JSON de `zones` que el SDUI engine exige).
- **L1..L3** (`system/layers/l{1,2,3}_*.go`): incremental — viewer + memberships (L1), `announcement-form` (L2), recurso `materials` con CRUD parcial (L3).
- **L4** (`system/l4/*.go`): sistema completo reorganizado por dominio. Cualquier cambio en runtime sobre pantallas/recursos/roles/permisos del producto vive aquí.
