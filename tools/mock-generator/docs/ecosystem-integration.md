# mock-generator ecosystem integration

## Rol ecosistemico

`tools/mock-generator` aparece en el `go.work` del ecosistema, pero no se observo consumo externo directo desde APIs, worker, frontend ni migrator.

## Lectura correcta del modulo en fase 2

No es una pieza de runtime del ecosistema. Es una herramienta auxiliar para generar datasets Go desde SQL.

## Integracion observada

### Con `postgres`

La integracion es directa:

- el generador importa `postgres/entities`
- su entrada esperada son scripts SQL
- su salida modela datasets tipados alineados con el dominio relacional

### Con el workspace local

`go.work` incluye `./edugo-infrastructure/tools/mock-generator`, lo que confirma que forma parte del toolkit local de desarrollo aunque no aparezca como dependencia de servicios ejecutables.

## Integracion probable

Por su proposito y nombre, la herramienta parece orientada a alimentar mocks o datasets para consumo local de desarrollo, posiblemente en capas de frontend o testing. Esa integracion no se encontro formalizada en los repos escaneados.

## Implicancia

Este modulo debe documentarse como herramienta interna del ecosistema, no como dependencia productiva de APIs o worker.
