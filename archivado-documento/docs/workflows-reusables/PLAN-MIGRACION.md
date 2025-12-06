# Plan de Migracion - Workflows Reusables

Plan detallado para migrar proyectos EduGo a workflows reusables.

---

## Estado Actual

| Proyecto | Workflows Locales | Lineas Codigo | Duplicacion | Prioridad |
|----------|-------------------|---------------|-------------|-----------|
| api-mobile | ci.yml, test.yml, docker.yml | ~120 | Alta | P0 |
| api-admin | ci.yml, test.yml, docker.yml | ~125 | Alta | P0 |
| worker | ci.yml, test.yml, docker.yml | ~130 | Alta | P0 |
| shared | test.yml, lint.yml | ~70 | Media | P1 |
| infrastructure | ci.yml, sync.yml | ~80 | Media | P1 |

**Total lineas duplicadas: ~525**

---

## Objetivo Post-Migracion

| Proyecto | Workflows Nuevos | Lineas Codigo | Reduccion |
|----------|------------------|---------------|-----------|
| api-mobile | ci.yml | ~25 | 79% |
| api-admin | ci.yml | ~25 | 80% |
| worker | ci.yml | ~25 | 80% |
| shared | ci.yml | ~20 | 71% |
| infrastructure | ci.yml, sync.yml | ~30 | 62% |

**Total lineas post-migracion: ~125** (reduccion 76%)

---

## Fase 1: api-mobile (Prioridad P0)

### Estado Actual

**Archivo:** `.github/workflows/ci.yml`

```yaml
# ~120 lineas con duplicacion
# - Setup Go manual
# - Tests manual
# - Coverage manual
# - Lint manual
# - Docker manual
```

### Estado Objetivo

**Archivo:** `.github/workflows/ci.yml`

```yaml
name: CI

on:
  push:
    branches: [main, dev]
  pull_request:
    branches: [main, dev]

jobs:
  test:
    uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/go-test.yml@v1.0.0
    with:
      go-version: '1.25'
      coverage-threshold: 33
    secrets:
      GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

  lint:
    uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/go-lint.yml@v1.0.0
    secrets:
      GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

  docker:
    needs: [test, lint]
    if: github.ref == 'refs/heads/main'
    uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/docker-build.yml@v1.0.0
    with:
      image-name: 'api-mobile'
    secrets:
      GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### Pasos de Migracion

1. **Backup actual**
   ```bash
   cp .github/workflows/ci.yml .github/workflows/ci.yml.backup
   ```

2. **Crear nuevo ci.yml**
   - Usar plantilla de infrastructure
   - Ajustar coverage-threshold si necesario
   - Validar nombres de imagen Docker

3. **Testing local**
   ```bash
   # Validar sintaxis
   act -l
   ```

4. **Commit y push**
   ```bash
   git checkout -b feat/workflows-reusables
   git add .github/workflows/ci.yml
   git commit -m "feat: migrar a workflows reusables de infrastructure"
   git push origin feat/workflows-reusables
   ```

5. **Crear PR**
   - Titulo: "Migrar a workflows reusables"
   - Descripcion: Comparativa antes/despues
   - Review: Validar 3+ ejecuciones exitosas

6. **Merge y cleanup**
   ```bash
   # Post-merge
   rm .github/workflows/ci.yml.backup
   rm .github/workflows/test.yml
   rm .github/workflows/docker.yml
   ```

### Validacion

- [ ] Tests pasan (3+ ejecuciones)
- [ ] Coverage >= 33%
- [ ] Lint sin errores
- [ ] Docker build exitoso
- [ ] Tiempo de ejecucion similar

### Rollback Plan

Si algo falla:
```bash
git revert HEAD
# O restaurar backup
cp .github/workflows/ci.yml.backup .github/workflows/ci.yml
```

---

## Fase 2: api-admin (Prioridad P0)

Identico a api-mobile, cambiar:
- `image-name: 'api-admin'`
- Ajustar `coverage-threshold` si necesario

---

## Fase 3: worker (Prioridad P0)

Identico a api-mobile, cambiar:
- `image-name: 'worker'`
- Ajustar `coverage-threshold` si necesario
- Considerar `platforms: 'linux/amd64'` si solo necesita amd64

---

## Fase 4: shared (Prioridad P1)

### Diferencias

- NO tiene Docker (es libreria)
- Coverage threshold mas alto (60%+)

```yaml
name: CI

on:
  push:
    branches: [main, dev]
  pull_request:
    branches: [main, dev]

jobs:
  test:
    uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/go-test.yml@v1.0.0
    with:
      go-version: '1.25'
      coverage-threshold: 60
    secrets:
      GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

  lint:
    uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/go-lint.yml@v1.0.0
    secrets:
      GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

---

## Fase 5: infrastructure (Prioridad P1)

### ci.yml

Similar a shared (sin Docker):

```yaml
jobs:
  test:
    uses: ./.github/workflows/reusable/go-test.yml
    with:
      working-directory: './postgres'
```

**Nota:** Usa path local `./.github` porque esta en mismo repo

### sync-main-to-dev.yml

Migrar a workflow reusable:

```yaml
name: Sync Main to Dev

on:
  push:
    branches:
      - main

jobs:
  sync:
    uses: ./.github/workflows/reusable/sync-branches.yml
    with:
      source-branch: 'main'
      target-branch: 'dev'
    secrets:
      GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

---

## Cronograma de Migracion

| Semana | Proyecto | Responsable | Estado |
|--------|----------|-------------|--------|
| 1 | api-mobile | Claude Code | Pendiente |
| 1 | api-admin | Claude Code | Pendiente |
| 2 | worker | Claude Code | Pendiente |
| 2 | shared | Claude Code | Pendiente |
| 3 | infrastructure | Claude Code | Pendiente |

**Duracion total estimada:** 3 semanas

---

## Checklist Pre-Migracion

Antes de migrar cada proyecto:

- [ ] infrastructure workflows reusables estan en `main`
- [ ] Tag `v1.0.0` creado en infrastructure
- [ ] Workflows testeados en infrastructure
- [ ] Documentacion actualizada
- [ ] Equipo notificado

---

## Checklist Post-Migracion

Despues de migrar cada proyecto:

- [ ] 5+ ejecuciones exitosas
- [ ] Coverage mantiene o mejora
- [ ] Tiempo de ejecucion similar
- [ ] Docker images funcionan
- [ ] Cleanup de archivos viejos
- [ ] Documentacion actualizada

---

## Metricas de Exito

### Pre-Migracion
- Duplicacion: ~70%
- Lineas totales: ~525
- Tiempo mantenimiento: Alto

### Post-Migracion
- Duplicacion: ~20%
- Lineas totales: ~125
- Tiempo mantenimiento: Bajo

### ROI
- Reduccion codigo: 76%
- Tiempo ahorrado: ~60% en mantenimiento
- Consistencia: 100%

---

## Soporte

Problemas durante migracion:
1. Revisar docs/workflows-reusables/GUIA-USO.md
2. Revisar docs/workflows-reusables/EJEMPLOS-INTEGRACION.md
3. Abrir issue en infrastructure
4. Contactar equipo EduGo

---

**Mantenido por:** EduGo Team
**Ultima actualizacion:** 21 Nov 2025
**Version:** 1.0
