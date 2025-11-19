# Sprint 4: Workflows Reusables - edugo-infrastructure

**Duración:** 5 días  
**Objetivo:** Crear workflows reusables para todo el ecosistema EduGo  
**Estado:** Listo para Ejecución (REQUIERE Sprint 1 completado)

---

## 📋 Resumen del Sprint

| Métrica | Objetivo |
|---------|----------|
| **Tareas Totales** | 15 |
| **Tiempo Estimado** | 20-25 horas |
| **Prioridad Alta (P0)** | 8 tareas |
| **Workflows Reusables** | 4 |
| **Composite Actions** | 3 |
| **Proyectos a Migrar** | 1+ (api-mobile mínimo) |
| **Commits Esperados** | 10-12 |
| **PRs a Crear** | 1 en infrastructure + 1 por proyecto migrado |

---

## 🏠 ¿Por Qué infrastructure?

### infrastructure es el HOGAR de workflows reusables

```
✅ Conceptualmente correcto (infraestructura CI/CD)
✅ Independiente de lógica de negocio
✅ Puede versionar workflows sin afectar features
✅ Centraliza herramientas y configuraciones
✅ Nombre coherente con propósito
```

### shared NO es el lugar adecuado

```
❌ shared contiene lógica de negocio (Logger, Auth, DB)
❌ Mezclaría concerns (business + tools)
❌ Versionar workflows allí sería confuso
❌ shared se usa como dependencia Go, no como tooling
```

---

## 🎯 Objetivos del Sprint

1. **Crear workflows reusables** que eliminen duplicación (~70% → ~20%)
2. **Crear composite actions** para bloques comunes
3. **Migrar al menos 1 proyecto** (api-mobile) a workflows reusables
4. **Documentar** uso, ejemplos y plan de migración
5. **Establecer** infrastructure como estándar de CI/CD

---

## 🗓️ Cronograma Diario

### Día 1: Setup y Composite Actions (5-6h)
- Tarea 1.1: Crear estructura para workflows reusables
- Tarea 1.2: Composite action - setup-edugo-go
- Tarea 1.3: Composite action - coverage-check
- Tarea 1.4: Composite action - docker-build-edugo

### Día 2: Workflows Reusables Core (5-6h)
- Tarea 2.1: Workflow reusable - go-test.yml
- Tarea 2.2: Workflow reusable - go-lint.yml
- Tarea 2.3: Workflow reusable - sync-branches.yml
- Tarea 2.4: Workflow reusable - docker-build.yml

### Día 3: Testing y Documentación (4-5h)
- Tarea 3.1: Testing exhaustivo de workflows
- Tarea 3.2: Documentación completa de uso
- Tarea 3.3: Ejemplos de integración

### Día 4: Migración de api-mobile (4-5h)
- Tarea 4.1: Migrar ci.yml de api-mobile
- Tarea 4.2: Migrar test.yml de api-mobile
- Tarea 4.3: Validar workflows migrados
- Tarea 4.4: PR en api-mobile

### Día 5: Review y Plan de Migración (2-3h)
- Tarea 5.1: Review completo de infrastructure
- Tarea 5.2: PR en infrastructure
- Tarea 5.3: Plan de migración para otros proyectos

---

## 📝 TAREAS DETALLADAS

---

## DÍA 1: SETUP Y COMPOSITE ACTIONS

---

### ✅ Tarea 1.1: Crear Estructura para Workflows Reusables

**Prioridad:** 🔴 Alta  
**Estimación:** ⏱️ 60 minutos  
**Prerequisitos:** Sprint 1 completado y en producción

#### Objetivo

Crear estructura en infrastructure para alojar workflows reusables y composite actions.

#### Verificar Prerequisitos

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure

# Verificar que Sprint 1 está completado
echo "Verificando estado actual..."

# 1. Check success rate reciente
gh run list --repo EduGoGroup/edugo-infrastructure --limit 10 --json conclusion \
  | jq '[.[] | select(.conclusion == "success")] | length'
# Debe ser >= 8 de 10

# 2. Check Go version
grep "^go " */go.mod
# Todos deben tener "go 1.25"

# 3. Check workflows estandarizados
ls -lh .github/workflows/
# Debe tener ci.yml y sync-main-to-dev.yml optimizados

echo "✅ Prerequisitos verificados"
```

#### Crear Estructura

```bash
# Crear rama de trabajo
git checkout dev
git pull origin dev
git checkout -b feature/workflows-reusables

# Crear estructura de directorios
mkdir -p .github/workflows/reusable
mkdir -p .github/actions/setup-edugo-go
mkdir -p .github/actions/coverage-check
mkdir -p .github/actions/docker-build-edugo

# Crear directorio de documentación
mkdir -p docs/workflows-reusables
```

#### Crear README de Workflows Reusables

```bash
cat > .github/workflows/reusable/README.md << 'README'
# Workflows Reusables - EduGo

Este directorio contiene workflows reusables que pueden ser consumidos por cualquier proyecto del ecosistema EduGo.

---

## 📋 Workflows Disponibles

| Workflow | Archivo | Propósito | Usado por |
|----------|---------|-----------|-----------|
| Go Test | `go-test.yml` | Tests unitarios y de integración | Todas las apps Go |
| Go Lint | `go-lint.yml` | Linter con golangci-lint | Todas las apps Go |
| Sync Branches | `sync-branches.yml` | Sincronización main → dev | Todos los repos |
| Docker Build | `docker-build.yml` | Build de imágenes Docker | APIs y Worker |

---

## 🔧 Composite Actions

| Action | Directorio | Propósito |
|--------|-----------|-----------|
| Setup EduGo Go | `../actions/setup-edugo-go/` | Setup Go + GOPRIVATE |
| Coverage Check | `../actions/coverage-check/` | Validar cobertura |
| Docker Build | `../actions/docker-build-edugo/` | Build Docker estándar |

---

## 📚 Cómo Usar

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

## 🎯 Versionado

**Recomendaciones:**
- **Producción:** Usar tag específico `@v1.0.0`
- **Desarrollo:** Usar `@dev` o `@main`

```yaml
# Producción (recomendado)
uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/go-test.yml@v1.0.0

# Desarrollo
uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/go-test.yml@dev
```

---

## 📖 Documentación Completa

Ver: [docs/workflows-reusables/](../../docs/workflows-reusables/)

---

**Última actualización:** $(date)  
**Versión:** 1.0
README

echo "✅ README de workflows reusables creado"
```

#### Crear Archivo de Versiones Centralizadas

```bash
cat > .github/config/versions.yml << 'VERSIONS'
# Versiones centralizadas para workflows reusables
# Actualizar aquí y todos los workflows se actualizan

go:
  version: "1.25"
  versions_matrix:
    - "1.24"
    - "1.25"
    - "1.26"

tools:
  golangci-lint: "v1.64.7"
  
github-actions:
  checkout: "v4"
  setup-go: "v5"
  upload-artifact: "v4"
  download-artifact: "v4"
  github-script: "v7"

docker:
  setup-buildx: "v3"
  login-action: "v3"
  build-push-action: "v5"
  metadata-action: "v5"

coverage:
  default-threshold: 33

# Configuración de EduGo
edugo:
  goprivate: "github.com/EduGoGroup/*"
  docker-registry: "ghcr.io"
  docker-platforms: "linux/amd64,linux/arm64"
VERSIONS

echo "✅ Versiones centralizadas configuradas"
```

#### Commit Inicial

```bash
git add .github/ docs/
git commit -m "feat: estructura para workflows reusables

Preparación para Sprint 4 - Workflows Reusables en infrastructure.

Estructura creada:
- .github/workflows/reusable/ (workflows reusables)
- .github/actions/ (composite actions)
- .github/config/versions.yml (versiones centralizadas)
- docs/workflows-reusables/ (documentación)

Por qué infrastructure:
- Es infraestructura CI/CD (conceptualmente correcto)
- Independiente de lógica de negocio
- Centraliza herramientas y configuraciones
- Nombre coherente con propósito

Próximos pasos:
- Implementar composite actions
- Implementar workflows reusables
- Migrar api-mobile

🤖 Generated with Claude Code

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

### ✅ Tarea 1.2: Composite Action - setup-edugo-go

**Prioridad:** 🔴 Alta  
**Estimación:** ⏱️ 90 minutos  
**Prerequisitos:** Tarea 1.1 completada

#### Objetivo

Crear composite action para setup de Go + GOPRIVATE, reemplazando ~15 líneas duplicadas en cada workflow.

#### Implementación

```bash
cat > .github/actions/setup-edugo-go/action.yml << 'ACTION'
name: 'Setup EduGo Go Environment'
description: 'Configura Go con versión estándar EduGo y acceso a repos privados'
author: 'EduGo Team'

branding:
  icon: 'package'
  color: 'blue'

inputs:
  go-version:
    description: 'Versión de Go a usar'
    required: false
    default: '1.25'
  
  cache:
    description: 'Habilitar cache de Go modules'
    required: false
    default: 'true'
  
  cache-dependency-path:
    description: 'Path a go.sum para cache'
    required: false
    default: 'go.sum'
  
  github-token:
    description: 'GitHub token para acceso a repos privados'
    required: false
    default: ${{ github.token }}

outputs:
  go-version:
    description: 'Versión de Go instalada'
    value: ${{ steps.setup.outputs.go-version }}
  
  cache-hit:
    description: 'Si el cache fue encontrado'
    value: ${{ steps.setup.outputs.cache-hit }}

runs:
  using: 'composite'
  steps:
    - name: Setup Go
      id: setup
      uses: actions/setup-go@v5
      with:
        go-version: ${{ inputs.go-version }}
        cache: ${{ inputs.cache }}
        cache-dependency-path: ${{ inputs.cache-dependency-path }}
    
    - name: Configurar acceso a repos privados
      shell: bash
      run: |
        echo "🔐 Configurando acceso a repos privados de EduGoGroup..."
        git config --global url."https://${{ inputs.github-token }}@github.com/".insteadOf "https://github.com/"
      env:
        GOPRIVATE: github.com/EduGoGroup/*
    
    - name: Verificar configuración
      shell: bash
      run: |
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        echo "✅ Go Environment Configurado"
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        echo "Go version: $(go version)"
        echo "GOROOT: $(go env GOROOT)"
        echo "GOPATH: $(go env GOPATH)"
        echo "GOPRIVATE: $GOPRIVATE"
        echo "Cache enabled: ${{ inputs.cache }}"
        echo "Cache hit: ${{ steps.setup.outputs.cache-hit }}"
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
ACTION

echo "✅ Composite action setup-edugo-go creada"
```

#### Crear README de la Action

```bash
cat > .github/actions/setup-edugo-go/README.md << 'README'
# Setup EduGo Go Environment

Composite action para configurar el entorno Go estándar de EduGo.

---

## ✨ Características

- ✅ Setup de Go con versión configurable
- ✅ Configuración automática de GOPRIVATE
- ✅ Acceso a repos privados de EduGoGroup
- ✅ Cache de Go modules
- ✅ Verificación de configuración
- ✅ Output de versión y cache hit

---

## 📖 Uso Básico

```yaml
steps:
  - uses: actions/checkout@v4
  
  - name: Setup Go
    uses: EduGoGroup/edugo-infrastructure/.github/actions/setup-edugo-go@main
```

---

## 🔧 Uso Avanzado

```yaml
steps:
  - name: Setup Go
    uses: EduGoGroup/edugo-infrastructure/.github/actions/setup-edugo-go@main
    with:
      go-version: '1.25'
      cache: true
      cache-dependency-path: '**/go.sum'
      github-token: ${{ secrets.GITHUB_TOKEN }}
```

---

## 📥 Inputs

| Input | Required | Default | Description |
|-------|----------|---------|-------------|
| `go-version` | No | `1.25` | Versión de Go |
| `cache` | No | `true` | Habilitar cache |
| `cache-dependency-path` | No | `go.sum` | Path para cache |
| `github-token` | No | `github.token` | Token para repos privados |

---

## 📤 Outputs

| Output | Description |
|--------|-------------|
| `go-version` | Versión de Go instalada |
| `cache-hit` | Si el cache fue encontrado (`true`/`false`) |

---

## 🔄 Equivalencia

**Antes (15+ líneas):**
```yaml
- uses: actions/setup-go@v5
  with:
    go-version: '1.25'
    cache: true

- name: Configurar repos privados
  run: |
    git config --global url."https://${{ secrets.GITHUB_TOKEN }}@github.com/".insteadOf "https://github.com/"
  env:
    GOPRIVATE: github.com/EduGoGroup/*

- name: Verificar
  run: |
    echo "Go: $(go version)"
    echo "GOPRIVATE: $GOPRIVATE"
```

**Después (1 línea):**
```yaml
- uses: EduGoGroup/edugo-infrastructure/.github/actions/setup-edugo-go@main
```

---

**Reducción:** ~93% menos código (15 líneas → 1 línea)

---

**Mantenido por:** EduGo Team  
**Última actualización:** $(date)
README

echo "✅ README de setup-edugo-go creado"
```

#### Crear Workflow de Testing

```bash
cat > .github/workflows/test-setup-go-action.yml << 'WORKFLOW'
name: Test - Setup Go Action

on:
  workflow_dispatch:
  push:
    paths:
      - '.github/actions/setup-edugo-go/**'

jobs:
  test-action:
    name: Test Setup Go Action
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Test con defaults
        id: test-default
        uses: ./.github/actions/setup-edugo-go
      
      - name: Verificar outputs
        run: |
          echo "Go version: ${{ steps.test-default.outputs.go-version }}"
          echo "Cache hit: ${{ steps.test-default.outputs.cache-hit }}"
      
      - name: Test compilación en infrastructure
        run: |
          cd postgres
          go build ./...
      
      - name: Test con versión específica
        uses: ./.github/actions/setup-edugo-go
        with:
          go-version: '1.24'
          cache: false
      
      - name: Verificar versión
        run: |
          go version | grep -q "go1.24" || exit 1
WORKFLOW

echo "✅ Workflow de testing creado"
```

#### Commit

```bash
git add .github/actions/setup-edugo-go/
git add .github/workflows/test-setup-go-action.yml
git commit -m "feat: composite action setup-edugo-go

Composite action para setup de Go + GOPRIVATE.

Reemplaza ~15 líneas de código repetido en cada workflow.

Características:
- Setup Go con versión configurable
- GOPRIVATE automático para EduGoGroup
- Cache de módulos incluido
- Outputs para uso posterior
- Testing automático

Uso:
  uses: EduGoGroup/edugo-infrastructure/.github/actions/setup-edugo-go@main

Reducción de código: ~93% (15 líneas → 1 línea)

🤖 Generated with Claude Code

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

### ✅ Tareas 1.3 y 1.4: Composite Actions Adicionales

**[Debido al límite de longitud, proporcionaré estructura similar a 1.2]**

#### Tarea 1.3: coverage-check (90 min)
- Valida cobertura de tests
- Genera reportes HTML
- Outputs para integración con PRs
- Script: Similar a shared

#### Tarea 1.4: docker-build-edugo (90 min)
- Build multi-platform
- Push a ghcr.io
- Tags automáticos
- Cache layers

---

## DÍAS 2-5: ESTRUCTURA RESUMIDA

### Día 2: Workflows Reusables Core (5-6h)

#### Tarea 2.1: go-test.yml (120 min)
```yaml
# Workflow reusable para tests
# Inputs: go-version, coverage-threshold, working-directory
# Jobs: test, coverage, upload-artifacts
```

#### Tarea 2.2: go-lint.yml (90 min)
```yaml
# Workflow reusable para linting
# Inputs: go-version, golangci-lint-version
# Jobs: lint con golangci-lint
```

#### Tarea 2.3: sync-branches.yml (90 min)
```yaml
# Workflow reusable para sync main → dev
# Inputs: target-branch, source-branch
# Jobs: merge automático con conflictos manejados
```

#### Tarea 2.4: docker-build.yml (90 min)
```yaml
# Workflow reusable para Docker build
# Inputs: context, dockerfile, platforms, tags
# Jobs: build multi-platform y push
```

---

### Día 3: Testing y Documentación (4-5h)

#### Tarea 3.1: Testing exhaustivo (120 min)
- Probar cada workflow reusable
- Validar inputs/outputs
- Casos edge

#### Tarea 3.2: Documentación completa (90 min)
- Guía de uso
- Ejemplos por proyecto
- Troubleshooting

#### Tarea 3.3: Ejemplos de integración (60 min)
- Ejemplos para api-mobile
- Ejemplos para api-admin
- Ejemplos para worker

---

### Día 4: Migración de api-mobile (4-5h)

#### Tarea 4.1: Migrar ci.yml (90 min)
```yaml
# ANTES: ~80 líneas
# DESPUÉS: ~20 líneas usando workflows reusables
jobs:
  test:
    uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/go-test.yml@main
```

#### Tarea 4.2: Migrar test.yml (90 min)
- Reemplazar con workflow reusable
- Validar coverage

#### Tarea 4.3: Validar migración (90 min)
- Tests en GitHub Actions
- Comparar con versión anterior

#### Tarea 4.4: PR en api-mobile (30 min)
- Crear PR con migración
- Documentar cambios

---

### Día 5: Review y Plan (2-3h)

#### Tarea 5.1: Review completo (60 min)
- Self-review de todos los cambios
- Checklist de completitud

#### Tarea 5.2: PR en infrastructure (45 min)
- PR con todos los workflows reusables
- Documentación completa

#### Tarea 5.3: Plan de migración (45 min)
- Documento para api-admin
- Documento para worker
- Priorización

---

## 📊 Métricas del Sprint 4

### Pre-Sprint 4
```yaml
code_duplication: "70%"
workflows_centralized: 0
composite_actions: 0
projects_using_reusables: 0
maintenance_effort: "Alto"
```

### Post-Sprint 4 (Objetivo)
```yaml
code_duplication: "20%"
workflows_centralized: 4
composite_actions: 3
projects_using_reusables: 1+
maintenance_effort: "Reducido 50%"
documentation: "Completa"
```

---

## 🎯 Entregables del Sprint 4

1. ✅ 4 Workflows Reusables funcionando
2. ✅ 3 Composite Actions funcionando
3. ✅ api-mobile migrado exitosamente
4. ✅ Documentación completa con ejemplos
5. ✅ Plan de migración para otros proyectos

---

## ✅ Checklist de Completitud

### Workflows Reusables
- [ ] go-test.yml funcional y documentado
- [ ] go-lint.yml funcional y documentado
- [ ] sync-branches.yml funcional y documentado
- [ ] docker-build.yml funcional y documentado

### Composite Actions
- [ ] setup-edugo-go funcional y documentado
- [ ] coverage-check funcional y documentado
- [ ] docker-build-edugo funcional y documentado

### Migración
- [ ] api-mobile usando workflows reusables
- [ ] Validado en GitHub Actions (5+ ejecuciones exitosas)
- [ ] PR mergeado en api-mobile

### Documentación
- [ ] README en .github/workflows/reusable/
- [ ] README en cada action
- [ ] Guía de uso en docs/workflows-reusables/
- [ ] Ejemplos de integración
- [ ] Plan de migración para otros proyectos

---

**¡Sprint 4 Completado!**

**Resultado:** infrastructure es ahora el hogar estándar de workflows reusables para todo EduGo

---

**Generado por:** Claude Code  
**Fecha:** 19 de Noviembre, 2025  
**Versión:** 1.0  
**Basado en:** Plan de shared v1.0 (Sprint 4)  
**Prerequisito:** Sprint 1 de infrastructure completado
