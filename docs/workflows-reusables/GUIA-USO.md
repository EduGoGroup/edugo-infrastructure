# Guia de Uso - Workflows Reusables EduGo

Documentacion completa para usar workflows reusables y composite actions de EduGo.

---

## Indice

1. [Workflows Reusables](#workflows-reusables)
2. [Composite Actions](#composite-actions)
3. [Ejemplos por Proyecto](#ejemplos-por-proyecto)
4. [Mejores Practicas](#mejores-practicas)
5. [Troubleshooting](#troubleshooting)

---

## Workflows Reusables

### go-test.yml

**Proposito:** Ejecutar tests unitarios con cobertura

**Uso basico:**
```yaml
jobs:
  test:
    uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/go-test.yml@main
    secrets:
      GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

**Uso avanzado:**
```yaml
jobs:
  test:
    uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/go-test.yml@main
    with:
      go-version: '1.25'
      coverage-threshold: 50
      working-directory: './app'
      run-race: true
      test-flags: '-short -v'
      upload-coverage: true
    secrets:
      GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

**Inputs:**
- `go-version` (default: 1.25)
- `coverage-threshold` (default: 33)
- `working-directory` (default: .)
- `run-race` (default: true)
- `test-flags` (default: -short)
- `upload-coverage` (default: true)

**Outputs:**
- `coverage`: Porcentaje de cobertura
- `test-result`: Resultado de los tests

---

### go-lint.yml

**Proposito:** Ejecutar golangci-lint

**Uso basico:**
```yaml
jobs:
  lint:
    uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/go-lint.yml@main
    secrets:
      GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

**Uso avanzado:**
```yaml
jobs:
  lint:
    uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/go-lint.yml@main
    with:
      go-version: '1.25'
      golangci-lint-version: 'v1.64.7'
      working-directory: './app'
      args: '--timeout=10m --verbose'
      skip-cache: false
    secrets:
      GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

---

### sync-branches.yml

**Proposito:** Sincronizar branches automaticamente

**Uso basico:**
```yaml
jobs:
  sync:
    uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/sync-branches.yml@main
    secrets:
      GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

**Uso avanzado:**
```yaml
jobs:
  sync:
    uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/sync-branches.yml@main
    with:
      source-branch: 'main'
      target-branch: 'dev'
      create-pr-on-conflict: true
      auto-merge: true
    secrets:
      GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

---

### docker-build.yml

**Proposito:** Build y push de imagenes Docker

**Uso basico:**
```yaml
jobs:
  docker:
    uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/docker-build.yml@main
    with:
      image-name: 'api-mobile'
    secrets:
      GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

**Uso avanzado:**
```yaml
jobs:
  docker:
    uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/docker-build.yml@main
    with:
      image-name: 'api-mobile'
      registry: 'ghcr.io'
      context: '.'
      dockerfile: 'Dockerfile'
      platforms: 'linux/amd64,linux/arm64'
      push: true
      tags: |
        latest
        v1.0.0
      build-args: |
        VERSION=1.0.0
        COMMIT=${{ github.sha }}
    secrets:
      GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

---

## Composite Actions

### setup-edugo-go

**Proposito:** Setup Go + GOPRIVATE

**Uso:**
```yaml
- uses: EduGoGroup/edugo-infrastructure/.github/actions/setup-edugo-go@main
  with:
    go-version: '1.25'
    cache: true
```

### coverage-check

**Proposito:** Validar cobertura de tests

**Uso:**
```yaml
- uses: EduGoGroup/edugo-infrastructure/.github/actions/coverage-check@main
  with:
    coverage-file: 'coverage.out'
    threshold: 33
    fail-on-threshold: true
```

### docker-build-edugo

**Proposito:** Build Docker estandarizado

**Uso:**
```yaml
- uses: EduGoGroup/edugo-infrastructure/.github/actions/docker-build-edugo@main
  with:
    image-name: 'api-mobile'
    registry-password: ${{ secrets.GITHUB_TOKEN }}
```

---

## Ejemplos por Proyecto

### api-mobile

```yaml
name: CI

on:
  push:
    branches: [main, dev]
  pull_request:
    branches: [main, dev]

jobs:
  test:
    uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/go-test.yml@main
    with:
      go-version: '1.25'
      coverage-threshold: 33
    secrets:
      GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

  lint:
    uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/go-lint.yml@main
    secrets:
      GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

  docker:
    needs: [test, lint]
    if: github.ref == 'refs/heads/main'
    uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/docker-build.yml@main
    with:
      image-name: 'api-mobile'
    secrets:
      GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### api-admin

Similar a api-mobile, cambiar `image-name` a `'api-admin'`

### worker

Similar a api-mobile, cambiar `image-name` a `'worker'`

---

## Mejores Practicas

### 1. Versionado

**Produccion (recomendado):**
```yaml
uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/go-test.yml@v1.0.0
```

**Desarrollo:**
```yaml
uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/go-test.yml@main
```

### 2. Secrets

Siempre pasar `GITHUB_TOKEN`:
```yaml
secrets:
  GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### 3. Coverage Threshold

Ajustar por proyecto:
- Nuevos proyectos: 60%+
- Proyectos legacy: 33%+
- Proyectos criticos: 80%+

### 4. Working Directory

Si tu proyecto tiene estructura especial:
```yaml
with:
  working-directory: './services/api'
```

---

## Troubleshooting

### Error: "workflow is not accessible"

**Causa:** Branch/tag no existe o repo es privado

**Solucion:**
1. Verificar que infrastructure es publico
2. Usar @main o tag valido
3. Verificar permisos de GITHUB_TOKEN

### Error: "required secret not provided"

**Causa:** Falta pasar GITHUB_TOKEN

**Solucion:**
```yaml
secrets:
  GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### Error: "coverage below threshold"

**Causa:** Cobertura insuficiente

**Solucion:**
1. Agregar mas tests
2. Ajustar threshold temporalmente
3. Excluir archivos generados

### Tests fallan en workflow pero pasan local

**Causa:** Diferencias de ambiente

**Solucion:**
1. Verificar flags (-short, -race)
2. Revisar timeouts
3. Check dependencias externas

---

**Mantenido por:** EduGo Team
**Ultima actualizacion:** 21 Nov 2025
**Version:** 1.0
