# 📚 EduGo Infrastructure - Documentación Completa

> Visión 360° del proyecto de infraestructura compartida para el ecosistema EduGo.
> Este documento proporciona una guía exhaustiva para entender, configurar y trabajar con la infraestructura central de EduGo.

---

## 📖 Tabla de Contenidos

1. [¿Qué es edugo-infrastructure?](#-qué-es-edugo-infrastructure)
2. [Contexto del Negocio](#-contexto-del-negocio)
3. [Índice de Documentación](#-índice-de-documentación)
4. [Arquitectura de Alto Nivel](#️-arquitectura-de-alto-nivel)
5. [Estructura del Proyecto](#️-estructura-del-proyecto)
6. [Servicios Requeridos](#-servicios-requeridos)
7. [Quick Start](#-quick-start)
8. [Proyectos Consumidores](#-proyectos-que-consumen-esta-infraestructura)
9. [Casos de Uso Principales](#-casos-de-uso-principales)
10. [Roadmap](#-roadmap)

---

## 🎯 ¿Qué es edugo-infrastructure?

**edugo-infrastructure** es el repositorio central que contiene:

- **Schemas de base de datos** (PostgreSQL + MongoDB)
- **Entities/Models** compartidas entre microservicios
- **Contratos de eventos** para mensajería (RabbitMQ)
- **Validadores JSON Schema** para eventos
- **Migraciones** de bases de datos
- **Configuración Docker** para desarrollo local

---

## 🏢 Contexto del Negocio

### ¿Qué es EduGo?

**EduGo** es una plataforma educativa integral que permite a instituciones educativas:

- **Gestionar materiales educativos** (PDFs, documentos, presentaciones)
- **Generar assessments automáticos** usando Inteligencia Artificial
- **Evaluar estudiantes** con quizzes generados desde el contenido
- **Seguir el progreso** académico de cada estudiante
- **Administrar escuelas** con estructuras jerárquicas flexibles

### Usuarios del Sistema

| Rol | Descripción | Acciones Principales |
|-----|-------------|---------------------|
| **Administrador** | Gestiona la escuela completa | Crear unidades académicas, matricular estudiantes, ver reportes |
| **Coordinador** | Supervisa docentes y cursos | Ver estadísticas, gestionar membresías |
| **Docente** | Crea contenido educativo | Subir materiales, ver assessments generados, revisar progreso |
| **Estudiante** | Consume contenido y rinde evaluaciones | Ver materiales, tomar assessments, ver resultados |
| **Apoderado** | Monitorea el progreso del estudiante | Ver calificaciones, reportes de progreso |

### Propuesta de Valor

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          PROPUESTA DE VALOR                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   ┌──────────────┐    ┌──────────────┐    ┌──────────────┐                  │
│   │   Docente    │    │     IA       │    │  Estudiante  │                  │
│   │   sube PDF   │───▶│   genera     │───▶│   resuelve   │                  │
│   │              │    │   quiz       │    │   quiz       │                  │
│   └──────────────┘    └──────────────┘    └──────────────┘                  │
│                                                                              │
│   💡 El docente ahorra horas de trabajo creando evaluaciones                │
│   📊 El estudiante recibe feedback inmediato                                │
│   📈 La escuela tiene métricas de aprendizaje en tiempo real                │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 📑 Índice de Documentación

| Documento | Descripción |
|-----------|-------------|
| [ARCHITECTURE.md](./ARCHITECTURE.md) | Arquitectura del sistema, componentes y flujos |
| [DATABASE_SCHEMA.md](./DATABASE_SCHEMA.md) | Modelo de datos PostgreSQL y MongoDB |
| [EVENT_CONTRACTS.md](./EVENT_CONTRACTS.md) | Eventos RabbitMQ y sus schemas |
| [SERVICES_SETUP.md](./SERVICES_SETUP.md) | Guía de configuración de servicios |
| [MODULES.md](./MODULES.md) | Descripción detallada de cada módulo |
| [DEVELOPMENT_GUIDE.md](./DEVELOPMENT_GUIDE.md) | Guía para desarrolladores |
| [PROCESS_FLOWS.md](./PROCESS_FLOWS.md) | Diagramas de secuencia detallados |
| [API_REFERENCE.md](./API_REFERENCE.md) | Referencia de endpoints esperados |
| [GLOSSARY.md](./GLOSSARY.md) | Glosario de términos del dominio |

---

## 🏗️ Arquitectura de Alto Nivel

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              CLIENTES                                        │
│                    (Mobile App / Web Admin Panel)                            │
└──────────────────────────────┬──────────────────────────────────────────────┘
                               │
          ┌────────────────────┼────────────────────┐
          │                    │                    │
          ▼                    ▼                    ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│   api-mobile    │  │api-administracion│  │     worker      │
│   (Estudiantes  │  │   (Admin Panel)  │  │   (AI/Tasks)    │
│    & Docentes)  │  │                  │  │                  │
└────────┬────────┘  └────────┬────────┘  └────────┬────────┘
         │                    │                    │
         └────────────────────┼────────────────────┘
                              │
         ┌────────────────────┼────────────────────┐
         │                    │                    │
         ▼                    ▼                    ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│   PostgreSQL    │  │    MongoDB      │  │    RabbitMQ     │
│   (Relacional)  │  │   (Documentos)  │  │   (Mensajería)  │
└─────────────────┘  └─────────────────┘  └─────────────────┘
                              │
                    ┌─────────┴─────────┐
                    │                   │
                    ▼                   ▼
             ┌───────────┐       ┌───────────┐
             │  AWS S3   │       │  OpenAI   │
             │ (Storage) │       │   (IA)    │
             └───────────┘       └───────────┘
```

---

## 🗂️ Estructura del Proyecto

```
edugo-infrastructure/
├── docker/                    # 🐳 Docker Compose para desarrollo
│   └── docker-compose.yml
│
├── postgres/                  # 🐘 Módulo PostgreSQL
│   ├── entities/              # Entities Go (structs)
│   ├── migrations/            # Migraciones SQL
│   └── cmd/                   # CLI de migraciones
│
├── mongodb/                   # 🍃 Módulo MongoDB
│   ├── entities/              # Entities Go (structs)
│   ├── migrations/            # Scripts de índices
│   └── seeds/                 # Datos de prueba
│
├── schemas/                   # 📋 JSON Schemas
│   └── events/                # Schemas de eventos
│
├── messaging/                 # 📬 Validación de eventos
│   ├── events/                # Tipos de eventos Go
│   └── validator.go           # Validador JSON Schema
│
├── seeds/                     # 🌱 Datos de prueba
│   ├── postgres/
│   └── mongodb/
│
├── scripts/                   # 🔧 Scripts de utilidad
├── tools/                     # 🛠️ Herramientas internas
└── documents/                 # 📚 Esta documentación
```

---

## 🔧 Servicios Requeridos

| Servicio | Puerto | Uso |
|----------|--------|-----|
| **PostgreSQL** | 5432 | Base de datos relacional principal |
| **MongoDB** | 27017 | Documentos (assessments, summaries) |
| **RabbitMQ** | 5672 / 15672 | Mensajería entre servicios |
| **Redis** | 6379 | Cache (opcional) |
| **AWS S3** | - | Almacenamiento de archivos |
| **OpenAI API** | - | Generación de contenido con IA |

---

## 🚀 Quick Start

```bash
# 1. Clonar repositorio
git clone git@github.com:EduGoGroup/edugo-infrastructure.git
cd edugo-infrastructure

# 2. Copiar variables de entorno
cp .env.example .env

# 3. Levantar servicios core (PostgreSQL + MongoDB)
make dev-up-core

# 4. Ejecutar migraciones
make migrate-up

# 5. Cargar datos de prueba (opcional)
make seed
```

---

## 📊 Proyectos que Consumen Esta Infraestructura

| Proyecto | Módulos Usados | Descripción |
|----------|----------------|-------------|
| **api-mobile** | postgres/entities, mongodb/entities, messaging | API para app móvil |
| **api-administracion** | postgres/entities, messaging | Panel de administración |
| **worker** | postgres/entities, mongodb/entities, messaging, schemas | Procesamiento con IA |

---

## 📖 Versiones

| Componente | Versión |
|------------|---------|
| Go | 1.22+ |
| PostgreSQL | 15 |
| MongoDB | 7.0 |
| RabbitMQ | 3.12 |
| Redis | 7 |

---

## 🔗 Links Útiles

- **GitHub:** [EduGoGroup/edugo-infrastructure](https://github.com/EduGoGroup/edugo-infrastructure)
- **PgAdmin (local):** http://localhost:5050
- **Mongo Express (local):** http://localhost:8082
- **RabbitMQ Management (local):** http://localhost:15672

---

---

## 🎯 Casos de Uso Principales

### CU-001: Subida de Material Educativo

**Actor:** Docente  
**Precondición:** Docente autenticado y con membresía activa  
**Flujo:**
1. Docente selecciona archivo PDF desde la app
2. Sistema sube archivo a S3
3. Sistema registra material en PostgreSQL
4. Sistema dispara evento `material.uploaded`
5. Worker procesa material con IA
6. Worker genera assessment y resumen
7. Worker actualiza estado a "ready"

**Postcondición:** Material disponible con quiz generado

### CU-002: Toma de Assessment

**Actor:** Estudiante  
**Precondición:** Estudiante matriculado, assessment publicado  
**Flujo:**
1. Estudiante solicita assessment disponible
2. Sistema retorna preguntas (sin respuestas correctas)
3. Estudiante responde cada pregunta
4. Estudiante envía intento completado
5. Sistema calcula score
6. Sistema muestra resultados y explicaciones

**Postcondición:** Intento registrado con score

### CU-003: Matrícula de Estudiante

**Actor:** Administrador  
**Precondición:** Escuela y unidad académica existentes  
**Flujo:**
1. Admin busca estudiante por email
2. Admin selecciona unidad académica destino
3. Sistema crea membership
4. Sistema dispara evento `student.enrolled`
5. Estudiante recibe notificación

**Postcondición:** Estudiante matriculado con acceso a contenido

### CU-004: Generación de Reportes

**Actor:** Coordinador/Administrador  
**Precondición:** Datos de assessments completados  
**Flujo:**
1. Actor selecciona rango de fechas y filtros
2. Sistema agrega datos de intentos
3. Sistema calcula métricas (promedio, desviación, etc.)
4. Sistema genera reporte visual
5. Actor exporta o visualiza reporte

**Postcondición:** Reporte generado

---

## 🗺️ Roadmap

### Fase Actual: Infraestructura Base ✅

- [x] Schema PostgreSQL completo (16 migraciones)
- [x] Schema MongoDB (3 collections)
- [x] Entities Go para PostgreSQL y MongoDB
- [x] Sistema de eventos RabbitMQ (4 eventos)
- [x] Validadores JSON Schema
- [x] Docker Compose para desarrollo
- [x] Documentación completa

### Próxima Fase: Integración

- [ ] Tests de integración end-to-end
- [ ] Migraciones para entities pendientes
- [ ] CI/CD con GitHub Actions
- [ ] Monitoreo con métricas

### Fase Futura: Escalabilidad

- [ ] Sharding de MongoDB
- [ ] Read replicas PostgreSQL
- [ ] Cache distribuido con Redis
- [ ] Rate limiting

---

## 📞 Contacto y Soporte

| Recurso | Descripción |
|---------|-------------|
| **GitHub Issues** | Reportar bugs o solicitar features |
| **Pull Requests** | Contribuir código |
| **Wiki** | Documentación extendida |

---

**Última actualización:** Diciembre 2024  
**Versión del documento:** 2.0  
**Mantenedores:** Equipo EduGo
