# Workflows Reusables - EduGo

Este directorio contiene workflows reusables que pueden ser consumidos por cualquier proyecto del ecosistema EduGo.

---

## Workflows Disponibles

| Workflow | Archivo | Propósito | Usado por |
|----------|---------|-----------|-----------|
| Go Test | `go-test.yml` | Tests unitarios y de integración | Todas las apps Go |
| Go Lint | `go-lint.yml` | Linter con golangci-lint | Todas las apps Go |
| Sync Branches | `sync-branches.yml` | Sincronización main - dev | Todos los repos |
| Docker Build | `docker-build.yml` | Build de imágenes Docker | APIs y Worker |

---

## Composite Actions

| Action | Directorio | Propósito |
|--------|-----------|-----------|
| Setup EduGo Go | `../actions/setup-edugo-go/` | Setup Go + GOPRIVATE |
| Coverage Check | `../actions/coverage-check/` | Validar cobertura |
| Docker Build | `../actions/docker-build-edugo/` | Build Docker estándar |

---

## Como Usar

### Workflow Reusable

```yaml
name: CI

on: [push, pull_request]

jobs:
  test:
    uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/go-test.yml@main
    with:
      go-version: '1.25'
      coverage-threshold: 33
```

### Composite Action

```yaml
steps:
  - uses: actions/checkout@v4

  - name: Setup Go
    uses: EduGoGroup/edugo-infrastructure/.github/actions/setup-edugo-go@main
```

---

## Versionado

**Recomendaciones:**
- **Produccion:** Usar tag especifico `@v1.0.0`
- **Desarrollo:** Usar `@dev` o `@main`

```yaml
# Produccion (recomendado)
uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/go-test.yml@v1.0.0

# Desarrollo
uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/go-test.yml@dev
```

---

## Documentacion Completa

Ver: [docs/workflows-reusables/](../../../docs/workflows-reusables/)

---

**Ultima actualizacion:** 21 Nov 2025
**Version:** 1.0
