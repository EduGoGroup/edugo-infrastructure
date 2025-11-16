# 🤝 Guía de Contribución - edugo-infrastructure

## 🔄 Workflow de Desarrollo

### 1. Crear Feature Branch

```bash
git checkout dev
git pull origin dev
git checkout -b feature/nombre-descriptivo
```

### 2. Hacer Cambios

```bash
# Editar archivos
# Agregar migraciones, schemas, etc.

git add .
git commit -m "feat(modulo): descripción del cambio"
```

### 3. Push y PR a dev

```bash
git push origin feature/nombre-descriptivo

# Crear PR
gh pr create --base dev --head feature/nombre-descriptivo
```

### 4. Después de Merge a dev → PR a main

```bash
# Crear PR de dev a main
gh pr create --base main --head dev --title "Release vX.Y.Z"
```

### 5. Después de Merge a main → Crear Tags

```bash
git checkout main
git pull origin main

# Tag general
git tag -a v0.2.0 -m "Release v0.2.0"

# Tags por módulo (si cambiaron)
git tag -a database/v0.2.0 -m "database v0.2.0"
git tag -a schemas/v0.2.0 -m "schemas v0.2.0"

git push origin --tags
```

### 6. Automático: Sync main → dev

El workflow `sync-main-to-dev.yml` sincroniza automáticamente.

---

## 📝 Convenciones de Commits

```
feat(modulo): agregar nueva funcionalidad
fix(modulo): corregir bug
docs: actualizar documentación
ci: cambios en CI/CD
chore: tareas de mantenimiento
```

---

## ✅ Checklist Antes de PR

- [ ] Tests pasan localmente
- [ ] Código formateado (`gofmt`, `goimports`)
- [ ] CHANGELOG.md actualizado
- [ ] README.md actualizado si aplica
- [ ] Co-Authored-By en commit
