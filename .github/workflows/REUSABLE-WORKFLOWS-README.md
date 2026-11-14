# Workflows y actions reutilizables de EduGo

Lo que otros repos consumen de aquí, medido el 2026-11-14 contra los `.github/` del ecosistema.

## Workflows reutilizables (`workflow_call`)

| Workflow | Consumidores |
|---|---|
| `reusable-go-lint.yml` | APIs, worker y fronts Go (11 referencias) |
| `reusable-go-test.yml` | `edugo-worker` |

```yaml
jobs:
  lint:
    uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable-go-lint.yml@main
```

## Composite actions

| Action | Para qué |
|---|---|
| `../actions/setup-edugo-go/` | Go, `GOPRIVATE` y caché de módulos |
| `../actions/coverage-check/` | Umbral de cobertura |
| `../actions/docker-build-edugo/` | Build Docker estándar |

```yaml
- uses: EduGoGroup/edugo-infrastructure/.github/actions/setup-edugo-go@main
```

## Retirados (2026-11-14, plan 054)

`reusable-docker-build.yml`, `reusable-sync-branches.yml` y `test-workflows-reusables.yml`. Ningún repo
los consumía. Los dos primeros declaraban `GITHUB_TOKEN` como secreto de `workflow_call`, un nombre
reservado, y el tercero apuntaba a rutas `reusable/*.yml` que no existen. GitHub los daba por inválidos y
generaba un run fallido en cada push.
