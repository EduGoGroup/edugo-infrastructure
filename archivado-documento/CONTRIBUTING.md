# 🤝 Guía de Contribución

## 🔄 Workflow de Desarrollo

### 1. Crear Feature Branch
```bash
git checkout dev
git pull origin dev
git checkout -b feature/nombre-descriptivo
```

### 2. Hacer Cambios y Commit
```bash
git add .
git commit -m "feat(modulo): descripción

Co-Authored-By: Claude <noreply@anthropic.com>"
```

### 3. PR a dev
```bash
git push origin feature/nombre-descriptivo
gh pr create --base dev --head feature/nombre-descriptivo
```

### 4. Después de merge: PR dev → main
```bash
gh pr create --base main --head dev --title "Release vX.Y.Z"
```

### 5. Crear tags en main
```bash
git checkout main && git pull origin main
git tag -a v0.2.0 -m "Release v0.2.0"
git tag -a database/v0.2.0 -m "database v0.2.0"
git push origin --tags
```

## ✅ Checklist
- [ ] Tests pasan
- [ ] Código formateado
- [ ] CHANGELOG actualizado
- [ ] Co-Authored-By en commit
