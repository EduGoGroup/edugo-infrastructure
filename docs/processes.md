# Processes

Esta vista resume los procesos que hoy existen en el repositorio. El detalle vive en la documentacion de cada modulo.

## Inventario de procesos

| Proceso | Modulo principal | Resultado | Detalle |
| --- | --- | --- | --- |
| Definir schema relacional | `postgres` | Schemas, tablas, funciones, vistas y FK | [../postgres/docs/processes.md](../postgres/docs/processes.md) |
| Sembrar configuracion canonica | `postgres` | resources, roles, permissions, UI templates, concept types | [../postgres/docs/processes.md](../postgres/docs/processes.md) |
| Sembrar dataset de desarrollo | `postgres` | escuelas, usuarios, memberships, materiales, assessments, intentos | [../postgres/docs/processes.md](../postgres/docs/processes.md) |
| Mantener collections documentales | `mongodb` | `material_summary`, `material_assessment_worker`, `material_event` | [../mongodb/docs/processes.md](../mongodb/docs/processes.md) |
| Poblar documentos canonicos y mock | `mongodb` | fixtures alineados con materiales de Postgres | [../mongodb/docs/processes.md](../mongodb/docs/processes.md) |
| Validar contratos de eventos | `schemas` | validacion de payloads JSON | [../schemas/docs/processes.md](../schemas/docs/processes.md) |
| Generar dataset Go desde SQL | `tools/mock-generator` | codigo generado para consumo local | [../tools/mock-generator/docs/processes.md](../tools/mock-generator/docs/processes.md) |
| Levantar runtime local | `docker` | Postgres, MongoDB, RabbitMQ, Redis y herramientas visuales | [../docker/docs/processes.md](../docker/docs/processes.md) |
| Orquestar calidad y CI | repo raiz | Makefiles, scripts, workflows y actions | [automation.md](automation.md) |

## Orden logico de lectura en fase 1

1. `docker` para entorno local.
2. `postgres` para la base estructural del dominio.
3. `mongodb` para documentos derivados y eventos internos del worker.
4. `schemas` para contratos de mensajeria.
5. `tools/mock-generator` para artefactos de dataset.
6. `automation.md` para entender lo que hoy esta alineado y lo que no.

## Nota de alcance

Aunque algunos procesos se tocan semantica o temporalmente, esta fase no documenta integracion entre modulos ni con el ecosistema global. Solo registra los procesos propios de cada superficie.
