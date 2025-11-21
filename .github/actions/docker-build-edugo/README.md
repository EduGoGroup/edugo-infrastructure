# EduGo Docker Build

Composite action para build y push de imagenes Docker con configuracion estandar EduGo.

---

## Caracteristicas

- Build multi-plataforma (amd64, arm64)
- Push a ghcr.io por defecto
- Tags automaticos (branch, PR, semver, sha)
- Cache de layers con GitHub Actions cache
- Configuracion estandar EduGo

---

## Uso Basico

```yaml
steps:
  - uses: actions/checkout@v4

  - name: Build and push
    uses: EduGoGroup/edugo-infrastructure/.github/actions/docker-build-edugo@main
    with:
      image-name: 'api-mobile'
      registry-password: ${{ secrets.GITHUB_TOKEN }}
```

---

## Uso Avanzado

```yaml
steps:
  - name: Build and push
    id: docker
    uses: EduGoGroup/edugo-infrastructure/.github/actions/docker-build-edugo@main
    with:
      image-name: 'api-mobile'
      registry: 'ghcr.io'
      registry-username: ${{ github.actor }}
      registry-password: ${{ secrets.GITHUB_TOKEN }}
      context: '.'
      dockerfile: 'Dockerfile'
      platforms: 'linux/amd64,linux/arm64'
      push: true
      tags: 'latest,v1.0.0'
      build-args: |
        VERSION=1.0.0
        COMMIT=${{ github.sha }}

  - name: Use outputs
    run: |
      echo "Image: ${{ steps.docker.outputs.image }}"
      echo "Digest: ${{ steps.docker.outputs.digest }}"
```

---

## Inputs

| Input | Required | Default | Description |
|-------|----------|---------|-------------|
| `image-name` | Yes | - | Nombre de la imagen |
| `registry` | No | `ghcr.io` | Registry de Docker |
| `registry-username` | No | `github.actor` | Usuario del registry |
| `registry-password` | Yes | - | Token del registry |
| `context` | No | `.` | Contexto de build |
| `dockerfile` | No | `Dockerfile` | Path al Dockerfile |
| `platforms` | No | `linux/amd64,linux/arm64` | Plataformas |
| `push` | No | `true` | Push al registry |
| `tags` | No | - | Tags adicionales |
| `build-args` | No | - | Build arguments |
| `cache-from` | No | `type=gha` | Cache from |
| `cache-to` | No | `type=gha,mode=max` | Cache to |

---

## Outputs

| Output | Description |
|--------|-------------|
| `image` | Full image name with tags |
| `digest` | Image digest |
| `metadata` | Build metadata JSON |

---

## Tags Automaticos

La action genera automaticamente tags basados en:

| Evento | Tag generado |
|--------|--------------|
| Push a branch | `branch-name` |
| Pull Request | `pr-123` |
| Tag semver | `v1.0.0`, `1.0` |
| Cualquier push | `sha-abc1234` |

---

## Equivalencia

**Antes (~40 lineas):**
```yaml
- uses: docker/setup-qemu-action@v3
- uses: docker/setup-buildx-action@v3
- uses: docker/login-action@v3
  with:
    registry: ghcr.io
    username: ${{ github.actor }}
    password: ${{ secrets.GITHUB_TOKEN }}
- uses: docker/metadata-action@v5
  id: meta
  with:
    images: ghcr.io/${{ github.repository }}
    # ... config de tags
- uses: docker/build-push-action@v5
  with:
    # ... muchos parametros
```

**Despues (5 lineas):**
```yaml
- uses: EduGoGroup/edugo-infrastructure/.github/actions/docker-build-edugo@main
  with:
    image-name: 'api-mobile'
    registry-password: ${{ secrets.GITHUB_TOKEN }}
```

---

**Reduccion:** ~87% menos codigo

---

**Mantenido por:** EduGo Team
**Ultima actualizacion:** 21 Nov 2025
