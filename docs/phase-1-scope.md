# Fase 1 Scope

## Objetivo

Definir una fuente de verdad documental autocontenida para `edugo-infrastructure`, basada solo en los archivos presentes en este repositorio.

## Que si entra en fase 1

- Modulos existentes en esta carpeta.
- Estructura real de directorios.
- Procesos implementados en SQL, Go, scripts y workflows.
- Estado actual de Makefiles, scripts y automatizaciones.
- Gaps visibles entre intencion y realidad.

## Que no entra en fase 1

- Integracion con otros repositorios de EduGo.
- Dependencias narradas desde el ecosistema global.
- Contratos de colaboracion entre equipos fuera de este repo.
- Reescritura de flujos de release o validacion aun no estandarizados.

## Criterio editorial

- No se backportea documentacion vieja solo por conservar historia.
- Si una automatizacion esta rota o desalineada, se documenta como tal.
- Si dos documentos repiten contenido, la version detallada vive en el modulo y el resto enlaza.

## Fase 2 prevista

La fase 2 debera tomar como insumo principal `/Users/jhoanmedina/source/EduGo/Common/ecosistema.md` y describir como cada modulo de este repo participa en el ecosistema y como se integra con otros modulos o repositorios.

## Fase 3 prevista

La fase 3 debera estandarizar validacion y release por modulo:

- `build`
- `test`
- `lint`
- `fmt`
- releases de GitHub por modulo
- actualizacion disciplinada de changelog por modulo
