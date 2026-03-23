# Releasing

## Convencion

- Cada modulo se versiona con tags `modulo/vX.Y.Z`.
- Las notas del release salen de `<modulo>/CHANGELOG.md`, seccion `## [X.Y.Z]`.
- `## [Unreleased]` es obligatoria en cada changelog (puede estar vacia entre releases).

## Prerrequisitos

- `git` con acceso de push al repositorio.
- `gh` autenticado (`gh auth status` debe pasar).
- `golangci-lint` instalado para modulos Go.
- Go instalado para `postgres`, `mongodb`, `schemas` y `tools/mock-generator`.
- Docker Compose disponible si se valida el modulo `docker`.

---

## Flujo normal (recomendado)

> El GitHub Release **lo crea el workflow automaticamente** cuando se hace push del tag.
> Tu trabajo es preparar el CHANGELOG y el tag. **No ejecutes `release-github` si el tag ya fue publicado.**

**1. Documentar los cambios**

Abre `<modulo>/CHANGELOG.md` y agrega los cambios bajo `## [Unreleased]`:

```markdown
## [Unreleased]

### Added
- Descripcion del cambio.
```

**2. Validar el modulo**

```bash
make -C <modulo> release-check
```

**3. Congelar el changelog**

```bash
make -C <modulo> release-prepare VERSION=vX.Y.Z
```

Mueve el contenido de `[Unreleased]` a `## [X.Y.Z] - YYYY-MM-DD` y deja `[Unreleased]` vacio.

**4. Revisar el diff**

```bash
git diff <modulo>/CHANGELOG.md
```

Confirma que la nueva seccion tiene el contenido correcto antes de continuar.

**5. Commitear el changelog y crear el tag**

```bash
git add <modulo>/CHANGELOG.md
git commit -m "<modulo>: release vX.Y.Z"
make -C <modulo> release-tag VERSION=vX.Y.Z
```

**6. Publicar el tag — el workflow hace el resto**

```bash
make -C <modulo> release-push-tag VERSION=vX.Y.Z
```

Al detectar el tag `modulo/vX.Y.Z`, el workflow ejecuta `release-check`,
extrae las notas del CHANGELOG y crea el GitHub Release automaticamente.
En modulos Go agrega ademas el bloque de instalacion con `go get`.

**7. Verificar**

```bash
gh release view "<modulo>/vX.Y.Z"
```

Un release correcto tiene:
- Titulo: `<modulo> vX.Y.Z` (con espacio, no slash)
- Notas: contenido del CHANGELOG para esa version
- En modulos Go: bloque `go get github.com/EduGoGroup/edugo-infrastructure/<modulo>@vX.Y.Z`

---

## Cuando el workflow falla o el release no se creo

Solo en este caso ejecuta `release-github` manualmente:

```bash
make -C <modulo> release-github VERSION=vX.Y.Z
```

Si el workflow creo un release incompleto o incorrecto, eliminalo primero:

```bash
gh release delete "<modulo>/vX.Y.Z" --yes
make -C <modulo> release-github VERSION=vX.Y.Z
```

---

## Ejemplo completo — modulo `postgres`

```bash
# 1. Editar postgres/CHANGELOG.md y agregar cambios bajo [Unreleased]

# 2. Validar
make -C postgres release-check

# 3. Congelar changelog
make -C postgres release-prepare VERSION=v0.67.0

# 4. Revisar
git diff postgres/CHANGELOG.md

# 5. Commitear y taggear
git add postgres/CHANGELOG.md
git commit -m "postgres: release v0.67.0"
make -C postgres release-tag VERSION=v0.67.0

# 6. Publicar (el workflow crea el GitHub Release automaticamente)
make -C postgres release-push-tag VERSION=v0.67.0

# 7. Verificar
gh release view "postgres/v0.67.0"
```

Para instalar en otro proyecto:

```bash
go get github.com/EduGoGroup/edugo-infrastructure/postgres@v0.67.0
```

---

## Alcance de validacion por modulo

- `postgres`: build, tests, lint, vet, fmt-check y changelog.
- `mongodb`: build, tests, lint, vet, fmt-check y changelog.
- `schemas`: build, tests, lint, vet, fmt-check y changelog.
- `tools/mock-generator`: build, tests, lint, vet, fmt-check y changelog.
- `docker`: `docker compose config -q` y changelog.

---

## Errores comunes

**"La seccion Unreleased esta vacia"**
Falta documentar los cambios en el CHANGELOG antes de `release-prepare`.

**"La version X.Y.Z ya existe en CHANGELOG"**
Ya se corrio `release-prepare` para esta version. Continua desde el paso 5.

**"El tag modulo/vX.Y.Z ya existe"**
El tag fue creado antes. Si el release no existe aun, crealo manualmente:
`make -C <modulo> release-github VERSION=vX.Y.Z`

**Release con titulo `modulo/vX.Y.Z` (con slash en el titulo)**
Fue creado manualmente desde la UI de GitHub. Editalo:
`gh release edit "<modulo>/vX.Y.Z" --title "<modulo> vX.Y.Z"`

**"Full Changelog" compara desde un modulo distinto**
Ocurre cuando el release anterior de ese modulo no era el ultimo release global en GitHub.
No afecta la funcionalidad ni el `go get`. Para corregirlo hay que recrear el release.
