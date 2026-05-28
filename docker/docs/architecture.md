# docker architecture

## Mapa interno

```text
docker/
|-- docs/
|-- docker-compose.yml
|-- README.md
`-- CHANGELOG.md
```

## Servicios definidos

| Servicio | Tipo | Perfil |
| --- | --- | --- |
| `postgres` | base de datos relacional | default |
| `mongodb` | base de datos documental | default |
| `rabbitmq` | mensajeria | `messaging` |
| `redis` | cache | `cache` |
| `pgadmin` | herramienta visual | `tools` |
| `mongo-express` | herramienta visual | `tools` |

## Diagrama local

```mermaid
flowchart TB
    ENV[.env values] --> DC[docker-compose.yml]
    DC --> PG[postgres]
    DC --> MG[mongodb]
    DC --> MQ[rabbitmq profile]
    DC --> RD[redis profile]
    DC --> TOOLS[pgadmin and mongo-express profile]
```

## Decisiones estructurales visibles

- Hay un core minimo y perfiles opcionales.
- Los servicios se exponen por puertos locales estables.
- La topologia privilegia desarrollo local y debugging visual.
