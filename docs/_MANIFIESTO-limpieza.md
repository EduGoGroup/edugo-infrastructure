# Manifiesto de limpieza — edugo-infrastructure/docs (2026-05-22)

Triage de la documentación interna de `edugo-infrastructure`. Sesgo conservador.
No se tocó código, seeds, docker/, ni .git.

## Resultado

Toda la documentación está vigente, es precisa y consistente con el estado real
del repo. Se verificó que cada módulo, script y Makefile referenciado existe
(postgres, mongodb, schemas, tools/mock-generator, docker; scripts/*; make/*.mk).
No hubo nada que borrar ni rescatar.

## Veredictos

| Archivo | Veredicto | Razón |
| --- | --- | --- |
| `README.md` (raíz) | KEEP | Índice vigente, módulos y scripts coinciden con el repo |
| `CHANGELOG.md` (raíz) | KEEP | Changelog activo (Unreleased de la base documental actual) |
| `docs/README.md` | KEEP | Índice general de docs, enlaces válidos |
| `docs/phase-1-scope.md` | KEEP | Alcance de fase 1, autocontenido y correcto |
| `docs/repository-map.md` | KEEP | Mapa de superficies y conteos reales (33 SQL, 27 entities, etc.) |
| `docs/processes.md` | KEEP | Inventario de procesos por módulo, enlaces válidos |
| `docs/architecture.md` | KEEP | Vista de arquitectura del repo, coincide con estructura |
| `docs/ecosystem-integration.md` | KEEP | Integración fase 2; documenta desalineaciones de ecosistema.md correctamente |
| `docs/automation.md` | KEEP | Superficies de automatización (Makefiles, scripts, workflows) alineadas |
| `docs/releasing.md` | KEEP | Flujo de release por módulo, fuente de verdad operativa |
| `docs/roadmap.md` | KEEP | Estado de fases 1/2/3 + pendiente operativo (primer release real) |

## Nota

`ecosystem-integration.md` cita nombres antiguos de APIs (admin-new/mobile-new) pero
provienen de `ecosistema.md` y el propio documento marca explícitamente esas
desalineaciones como tales. Es documentación correcta de la integración, no obsoleta.
