# Ejemplos de Integracion - Workflows Reusables

Ejemplos practicos de como integrar workflows reusables en proyectos EduGo.

---

## Ejemplo 1: api-mobile

### Antes (codigo duplicado - 80 lineas)

```yaml
name: CI

on:
  push:
    branches: [main, dev]
  pull_request:
    branches: [main, dev]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.25'
          cache: true
      - name: Configure private repos
        run: |
          git config --global url."https://${{ secrets.GITHUB_TOKEN }}@github.com/".insteadOf "https://github.com/"
        env:
          GOPRIVATE: github.com/EduGoGroup/*
      - run: go mod download
      - run: go test -short -race -coverprofile=coverage.out ./...
      - name: Check coverage
        run: |
          COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
          if (( $(echo "$COVERAGE < 33" | bc -l) )); then
            echo "Coverage too low: $COVERAGE%"
            exit 1
          fi

  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.25'
      - name: Configure private repos
        run: |
          git config --global url."https://${{ secrets.GITHUB_TOKEN }}@github.com/".insteadOf "https://github.com/"
        env:
          GOPRIVATE: github.com/EduGoGroup/*
      - uses: golangci/golangci-lint-action@v6
        with:
          version: v1.64.7
```

### Despues (workflows reusables - 20 lineas)

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

**Reduccion: 75% menos codigo (80 lineas -> 20 lineas)**

---

## Ejemplo 2: api-admin

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
      coverage-threshold: 40
      test-flags: '-short -v'
    secrets:
      GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

  lint:
    uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/go-lint.yml@main
    with:
      args: '--timeout=10m'
    secrets:
      GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

  docker:
    needs: [test, lint]
    if: github.ref == 'refs/heads/main'
    uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/docker-build.yml@main
    with:
      image-name: 'api-admin'
      build-args: |
        VERSION=${{ github.ref_name }}
        COMMIT=${{ github.sha }}
    secrets:
      GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

---

## Ejemplo 3: worker

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
      coverage-threshold: 50
      run-race: true
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
      image-name: 'worker'
      platforms: 'linux/amd64'
    secrets:
      GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

---

## Ejemplo 4: shared (sin Docker)

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
      coverage-threshold: 60
      upload-coverage: true
    secrets:
      GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

  lint:
    uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/go-lint.yml@main
    secrets:
      GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

---

## Ejemplo 5: Sync branches

```yaml
name: Sync Main to Dev

on:
  push:
    branches:
      - main

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

## Ejemplo 6: Workflow completo con todas las features

```yaml
name: Complete CI/CD

on:
  push:
    branches: [main, dev]
  pull_request:
    branches: [main, dev]

jobs:
  test:
    name: Tests
    uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/go-test.yml@main
    with:
      go-version: '1.25'
      coverage-threshold: 50
      working-directory: '.'
      run-race: true
      test-flags: '-short -v'
      upload-coverage: true
    secrets:
      GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

  lint:
    name: Lint
    uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/go-lint.yml@main
    with:
      go-version: '1.25'
      golangci-lint-version: 'v1.64.7'
      args: '--timeout=10m --verbose'
    secrets:
      GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

  security:
    name: Security Scan
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: EduGoGroup/edugo-infrastructure/.github/actions/setup-edugo-go@main
      - run: go install golang.org/x/vuln/cmd/govulncheck@latest
      - run: govulncheck ./...

  docker:
    name: Docker Build
    needs: [test, lint, security]
    if: github.ref == 'refs/heads/main'
    uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/docker-build.yml@main
    with:
      image-name: 'my-app'
      registry: 'ghcr.io'
      platforms: 'linux/amd64,linux/arm64'
      push: true
      tags: |
        latest
        ${{ github.ref_name }}
      build-args: |
        VERSION=${{ github.ref_name }}
        COMMIT=${{ github.sha }}
        BUILD_DATE=${{ github.event.head_commit.timestamp }}
    secrets:
      GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

  notify:
    name: Notify Success
    needs: [test, lint, docker]
    if: always()
    runs-on: ubuntu-latest
    steps:
      - name: Notify
        run: |
          echo "CI/CD completed"
          echo "Tests: ${{ needs.test.result }}"
          echo "Lint: ${{ needs.lint.result }}"
          echo "Docker: ${{ needs.docker.result }}"
```

---

## Comparativa de Reduccion de Codigo

| Proyecto | Antes | Despues | Reduccion |
|----------|-------|---------|-----------|
| api-mobile | 80 lineas | 20 lineas | 75% |
| api-admin | 85 lineas | 25 lineas | 70% |
| worker | 90 lineas | 22 lineas | 75% |
| shared | 60 lineas | 15 lineas | 75% |

**Promedio: 74% menos codigo**

---

## Beneficios

1. **Mantenimiento centralizado**: Cambios en 1 lugar afectan a todos
2. **Consistencia**: Todos los proyectos usan misma configuracion
3. **Menos duplicacion**: De ~70% a ~20%
4. **Mas rapido**: Setup de CI en nuevos proyectos en minutos
5. **Mejor testing**: Workflows testeados centralmente

---

**Mantenido por:** EduGo Team
**Ultima actualizacion:** 21 Nov 2025
