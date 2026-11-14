# Workflows y actions reutilizables de EduGo

Lo que otros repos consumen de aquí, medido el 2026-11-14 contra los `.github/` del ecosistema.

## Workflows reutilizables (`workflow_call`)

| Workflow | Consumidores |
|---|---|
| `reusable-go-lint.yml` | Las 5 APIs (`pr-to-dev.yml` y `pr-to-main.yml`) y `edugo-worker` (`ci.yml`): 11 referencias |
| `reusable-go-test.yml` | `edugo-worker` |

```yaml
jobs:
  lint:
    uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable-go-lint.yml@main
```

Los consumen con `@main`: **no se cambia un default ni se retira un input** sin revisar a los seis
consumidores. Un input que el reusable no declara (el caso `use-services` del worker) hace que GitHub
rechace el workflow llamador con `startup_failure`; se arregla del lado del consumidor.

| Workflow | Inputs |
|---|---|
| `reusable-go-lint.yml` | `go-version`, `golangci-lint-version`, `working-directory`, `args`, `skip-cache` |
| `reusable-go-test.yml` | `go-version`, `coverage-threshold`, `working-directory`, `run-race`, `test-flags`, `upload-coverage` |

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
