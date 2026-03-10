# docker ecosystem integration

## Rol ecosistemico

La superficie `docker/` de este repo sirve para trabajo local aislado sobre infraestructura, pero no es el orquestador canonico del ecosistema completo.

## Relacion con `edugo-dev-environment`

`ecosistema.md` posiciona a `edugo-dev-environment` como la herramienta principal para levantar ambiente local y ejecutar migraciones.

Por eso la relacion correcta es:

- `docker/` en este repo: soporte local para validar Postgres, MongoDB y servicios auxiliares de forma aislada
- `edugo-dev-environment`: flujo principal de ambiente compartido y migraciones ecosistemicas

## Cuándo usar esta superficie

### Uso adecuado

- trabajar solo en `edugo-infrastructure`
- levantar dependencias minimas para tests o validacion manual
- inspeccionar Postgres y MongoDB con herramientas visuales

### Uso no canonico

- recrear por completo el ambiente del ecosistema
- reemplazar el flujo de migracion del migrador
- asumir que esta compose file representa todo el ecosistema EduGo

## Integracion interna

`docker/` da soporte directo a:

- `postgres`
- `mongodb`
- pruebas locales de `schemas` si requieren servicios productores/consumidores externos

## Implicancia

Esta superficie debe seguir documentada, pero subordinada a `edugo-dev-environment` en cualquier explicacion de ecosistema.
