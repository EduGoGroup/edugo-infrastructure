# Plantillas de Workflows Reusables

Plantillas listas para copiar y usar en proyectos EduGo.

---

## Plantillas Disponibles

| Plantilla | Uso | Proyectos |
|-----------|-----|-----------|
| `api-con-docker.yml` | APIs con Docker | api-mobile, api-admin, worker |
| `libreria-sin-docker.yml` | Librerias Go | shared, modulos comunes |
| `sync-branches.yml` | Sync automatico | Todos los repos |

---

## Como Usar

### 1. Elegir Plantilla

Segun tipo de proyecto:
- **API con Docker**: `api-con-docker.yml`
- **Libreria Go**: `libreria-sin-docker.yml`
- **Sync branches**: `sync-branches.yml`

### 2. Copiar Plantilla

```bash
# Ejemplo: api-mobile
cp docs/workflows-reusables/plantillas/api-con-docker.yml .github/workflows/ci.yml
```

### 3. Personalizar

Editar el archivo copiado:
- Reemplazar `{{IMAGE_NAME}}` con nombre real
- Ajustar `coverage-threshold` si necesario
- Ajustar `args` de linter si necesario

### 4. Backup Workflows Viejos

```bash
mkdir -p .github/workflows/backup
mv .github/workflows/*.yml .github/workflows/backup/
mv .github/workflows/backup/ci.yml .github/workflows/ci.yml
```

### 5. Commit y Push

```bash
git checkout -b feat/workflows-reusables
git add .github/workflows/ci.yml
git commit -m "feat: migrar a workflows reusables"
git push origin feat/workflows-reusables
```

### 6. Crear PR

- Titulo: "Migrar a workflows reusables"
- Incluir comparativa antes/despues
- Esperar 3+ ejecuciones exitosas
- Merge

### 7. Cleanup

```bash
# Post-merge
rm -rf .github/workflows/backup/
```

---

## Ejemplo: api-mobile

### Antes
```yaml
# .github/workflows/ci.yml (120 lineas)
name: CI
on:
  push:
    branches: [main, dev]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.25'
      # ... 100+ lineas mas
```

### Despues
```yaml
# .github/workflows/ci.yml (25 lineas)
name: CI
on:
  push:
    branches: [main, dev]
jobs:
  test:
    uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/go-test.yml@v1.0.0
    secrets:
      GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
  # ... solo 20 lineas mas
```

**Reduccion: 79%**

---

## Validacion

Antes de mergear, validar:

- [ ] Tests pasan
- [ ] Coverage >= threshold
- [ ] Lint sin errores
- [ ] Docker build exitoso (si aplica)
- [ ] Tiempo similar o mejor

---

## Rollback

Si algo falla:

```bash
# Opcion 1: Revert commit
git revert HEAD

# Opcion 2: Restaurar backup
mv .github/workflows/backup/*.yml .github/workflows/
```

---

## Soporte

Problemas:
1. Revisar [GUIA-USO.md](../GUIA-USO.md)
2. Revisar [EJEMPLOS-INTEGRACION.md](../EJEMPLOS-INTEGRACION.md)
3. Revisar [PLAN-MIGRACION.md](../PLAN-MIGRACION.md)
4. Abrir issue en infrastructure

---

**Mantenido por:** EduGo Team
**Ultima actualizacion:** 21 Nov 2025
