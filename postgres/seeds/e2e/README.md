# Seed E2E — Fixtures focalizables (Fase C)

> Framework de fixtures compositivas para el plan E2E del monorepo
> EduGo. Reemplaza la pila aditiva `fase0..fase4` por piezas atómicas
> que se combinan a demanda. Spec: `EduUI/edugo-ui-kmp/e2e-integration-plan/system-data-quality-spec/phase-c-fixtures-refactor/`.

## Estructura

```
seeds/e2e/
├── e2e.go                              # Shim de compat: constantes + Apply([]string)
├── fase{2,3,4}_*.go                    # Sólo constantes públicas legacy (sin lógica)
├── framework/                          # Motor: composer, cleanup, registry, logger
├── fixtures/                           # Piezas atómicas (role_only, screen_only, ...)
├── scenarios/                          # Recetas (observer_audits, ..., legacy_e2e)
└── exports/                            # fixtures-constants.json (artefacto de build)
```

## Modelo conceptual

```
Production seed (inmutable)
  resources, roles, permissions, role_permissions del catálogo,
  screen_templates, screen_instances, resource_screens, concept_*
       │
       ▼
Scenario (overlay con prefijo único)
  E2E-XXXXX-<entity>-NN   ← TenantPrefix (códigos visibles)
  e2eXXXXX-...            ← SchemaPrefix (UUIDs)
       │
  ┌────┴────────────────────────────────┐
  ▼                                     ▼
Fixture A (Provides:["school"])   Fixture B (Requires:["school"])
```

## Catálogo de fixtures

| Nombre              | Provides                                      | Requires        | Tablas tocadas                                                              |
|---------------------|-----------------------------------------------|-----------------|------------------------------------------------------------------------------|
| `role_only`         | `school`, `user`, `user_role`, `membership`   | —               | `academic.schools`, `auth.users`, `iam.user_roles`, `academic.memberships`  |
| `screen_only`       | `screen_data`                                 | `school`        | depende del `ScreenKey`: `assessment.assessment` (`assessments-list`); `academic.{subjects,academic_periods,memberships,grades}` + `auth.users` (`grades-list`) |
| `readonly_role`     | `readonly_role`                               | —               | `iam.roles` (overlay), `iam.role_permissions` (con SchemaPrefix)            |
| `partial_crud`      | `partial_crud_role`                           | —               | `iam.roles`, `iam.role_permissions`                                         |
| `menu_subtree`      | `menu_subtree`                                | `readonly_role` | `iam.role_permissions` (subset según subtree)                               |
| `guardian_relation` | `guardian_relation`                           | `school`, `user`| `auth.users` (student), `academic.memberships` (student), `academic.guardian_relations` |
| `legacy_school_admin`              | `legacy_school`, `legacy_admin`, …  | —               | `academic.schools` + admin                                                  |
| `legacy_student_dual_school`       | `legacy_student`, `legacy_school2`  | `legacy_school` | `auth.users`, `iam.user_roles`, `academic.memberships`, segunda school     |
| `legacy_announcement`              | `legacy_announcement`               | `legacy_school`,`legacy_admin` | `academic.announcements`                                       |
| `legacy_unit_subject`              | `legacy_unit`, `legacy_subject`     | `legacy_school` | `academic.academic_units`, `academic.subjects`                              |
| `legacy_course_material_assessment`| `legacy_learning`                   | `legacy_school`,`legacy_admin`,`legacy_student`,`legacy_unit` | `content.courses`, `content.materials`, `assessment.*` |

## Catálogo de scenarios

| Nombre                   | Tags                       | Composición                                                                          |
|--------------------------|----------------------------|--------------------------------------------------------------------------------------|
| `observer_audits`        | rbac, menu, audit          | `role_only(readonly_auditor)` + `readonly_role(audit-events)` + `menu_subtree(audit-events)` |
| `teacher_grades_only`    | rbac, screen-config        | `role_only(teacher)` + `partial_crud(grades)` + `screen_only(grades-list)`           |
| `guardian_views_child`   | rbac, screen-config        | `role_only(guardian)` + `guardian_relation` + `screen_only(child-progress)` |
| `legacy_e2e`             | legacy, backward-compat    | `legacy_school_admin` + `legacy_student_dual_school` + `legacy_announcement` + `legacy_unit_subject` + `legacy_course_material_assessment` |

## Cómo escribir una fixture nueva

1. Crear `seeds/e2e/fixtures/<nombre>.go` con un struct que implemente
   `framework.Fixture` (`Manifest()`, `Apply()`, `Cleanup()`).
2. Declarar en `Manifest`:
   - `Name` único (snake_case sin prefijo `fixture_`).
   - `Provides`: capacidades que aporta. Ej. `["readonly_role"]`.
   - `Requires`: capacidades que necesita. Ej. `["school"]`.
   - `Tables`: tablas tocadas (formato `schema.table`) en orden de
     creación. El cleanup las recorre al revés.
   - `Constants`: claves que se exportarán al JSON de Kotlin
     (`E2E<Fixture><EntityKind>`).
3. En `Apply`:
   - Usar `framework.MakeUUID(ctx, suffix)` y `framework.MakeCode(ctx, kind, idx)`.
   - Antes de cada INSERT con UUID generado, llamar
     `framework.AssertNotProductionNamespace(uuid)`.
   - Idempotencia: `clause.OnConflict{ DoNothing: true }` por id.
   - Booleanos críticos (F2·H5): forzar con
     `framework.UpsertBool(tx, table, "id", id, "is_active", true)`.
   - Registrar en `ctx.Provide(capability, ProvidedEntity{...})` cada
     entidad que otra fixture pueda querer reusar.
   - Poblar `ctx.SetConstant(...)` con los valores que los tests Kotlin
     vayan a leer.
4. En `Cleanup`: `framework.DeleteByPrefix(tx, table, "id", ctx.SchemaPrefix)`
   en orden inverso al de creación.
5. Tests:
   - **Manifest** (verifica Provides/Requires/Tables/Constants).
   - **Apply nil tx** (confirma que las precondiciones son sin BD).
   - Tests con BD real bajo build tag `integration` (Wave 6).

## Cómo escribir un scenario nuevo

1. Crear `seeds/e2e/scenarios/scenario_<nombre>.go` con un struct que
   implemente `framework.Scenario`.
2. `Manifest`: `Name` único, `Description`, `Tags`, `FixtureNames`
   (orden declarativo, hint para el composer).
3. `BuildFixtures(ctx)` retorna las instancias en el orden esperado.
   El composer resuelve dependencias por `Provides`/`Requires` antes
   de tocar la BD.
4. Registrar en `scenarios.RegisterAll(reg)`:
   `_ = reg.RegisterScenario(&MiScenario{})`.

## Migrando desde el seed E2E legacy

Hasta Fase C el seed se invocaba con el target `cloud-seed-e2e` (eliminado).
El contrato equivalente sigue funcionando vía `scenario_legacy_e2e`: usa el
scenario legacy que reproduce bit-a-bit los mismos UUIDs/códigos/emails.
Para tests nuevos, preferí scenarios atómicos:

```bash
make cloud-seed-scenario SCENARIO=teacher_grades_only
make cloud-seed-scenario SCENARIO=legacy_e2e
```

Cualquier código existente que importe constantes públicas
(`E2ESchoolCode`, `E2EAdminEmail`, …) sigue funcionando — vienen del
shim `e2e.go` y se siguen produciendo bit-a-bit por
`scenario_legacy_e2e`.

## Constantes para tests Kotlin

Tras cada `Apply` el binario `seed_e2e` regenera
`seeds/e2e/exports/fixtures-constants.json` con:

```json
{
  "schemaVersion": "1",
  "generatedAt": "2026-05-08T12:00:00Z",
  "scenarios": {
    "teacher_grades_only": {
      "tenantPrefix": "E2E-A1B2C-",
      "schemaPrefix": "e2ea1b2c-",
      "constants": {
        "E2EFixtureRoleOnlyUserEmail": "teacher-role_only-a1b2c@edugo.test",
        "E2EFixtureRoleOnlySchoolCode": "E2E-A1B2C-SCHOOL-01",
        ...
      }
    }
  }
}
```

Los tests Kotlin lo consumen vía
`com.edugo.kmp.repository.e2e.E2EFixtureConstants` (ubicado en
`EduUI/edugo-ui-kmp/modules/repository/src/commonTest`). El parser
valida `schemaVersion` y falla con mensaje accionable si hay drift
entre el binario y el helper Kotlin.

## SchemaVersion

Cualquier cambio en este directorio (`framework/`, `fixtures/`,
`scenarios/` o `exports/`) requiere bumpear `SchemaVersion` en
`migrations/version.go`. Es la regla mandatada por
`postgres/seeds/CLAUDE.md`.

## Tests

| Categoría    | Build tag       | Cómo correr                                                                                                              |
|--------------|-----------------|--------------------------------------------------------------------------------------------------------------------------|
| Unit (puros) | (sin tag)       | `go test ./seeds/e2e/... -count=1`                                                                                       |
| Integration  | `integration`   | `ENABLE_INTEGRATION_TESTS=true go test -tags=integration -count=1 -timeout=15m ./seeds/e2e/...`                          |
| Benchmarks   | `integration`   | `ENABLE_INTEGRATION_TESTS=true go test -tags=integration -bench=. -count=1 -benchtime=1x -timeout=15m ./seeds/e2e/framework/...` |

Los tests sin tag corren en CI vanilla (no requieren BD). Los que
necesitan BD viven bajo:

- **Build tag** `integration` (gate de compilación).
- **Variable de entorno** `ENABLE_INTEGRATION_TESTS=true` (gate de
  ejecución — sin ella los tests se skipean con `t.Skip(...)`).

### Backend de BD para los tests integration

El helper `seeds/e2e/internal/testdb.StartPostgres(tb)` decide el
backend automáticamente:

- **Modo local (default)** — arranca un contenedor
  `postgres:15-alpine` con `testcontainers-go`, aplica
  `migrations.Migrate(Force=true)` y `system.ApplySystem(db, "")`.
  Necesita Docker corriendo.
- **Modo cloud (override)** — si la variable `POSTGRES_URI` está
  definida, abre conexión directa contra esa URI sin levantar
  contenedor. Útil para los benchmarks que miden la cota real
  contra Neon:

  ```bash
  source EduBack/edugo-dev-environment/migrator/.env.cloud
  ENABLE_INTEGRATION_TESTS=true POSTGRES_URI=$POSTGRES_URI \
      go test -tags=integration -bench=. -count=1 -benchtime=1x \
          -timeout=15m ./seeds/e2e/framework/...
  ```

### Cobertura

Las cotas de cobertura por paquete (≥80%) se miden con
`-coverpkg` cruzado, porque las fixtures se ejercitan principalmente
desde los integration tests del paquete `scenarios`:

```bash
ENABLE_INTEGRATION_TESTS=true go test -tags=integration -count=1 \
    -timeout=15m \
    -coverpkg=./seeds/e2e/fixtures,./seeds/e2e/framework,./seeds/e2e/scenarios \
    -coverprofile=/tmp/e2e-cover.out \
    ./seeds/e2e/...
go tool cover -func=/tmp/e2e-cover.out | grep '^total:'
```

Cierre 2026-05-08: total combinado **80.0%**.
