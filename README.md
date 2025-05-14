# Ualá Backend Challenge - Plataforma tipo Twitter 🐦

Este desafío consiste en desarrollar una versión simplificada de una plataforma de microblogging similar a Twitter, donde los usuarios pueden publicar tweets, seguir a otros usuarios y ver un timeline personalizado de los tweets de las personas a las que siguen.

## 🧩 Tecnologías utilizadas

- **Lenguaje**: Go (Golang)
- **Arquitectura**: Hexagonal (Ports & Adapters)
- **Base de Datos**: PostgreSQL
- **Cache**: Redis
- **Documentación de la API**: Swagger (OpenAPI)
- **Contenedores**: Docker
- **Despliegue**: Amazon Web Services (ECS)
- **Testing**: Unit tests

## 📦 Cómo levantar el proyecto

### Opción 1: Docker (recomendado)

```bash
docker-compose up --build
```
Este comando levanta todos los servicios necesarios, crea la base de datos y las tablas (API + PostgreSQL + Redis) y expone el servicio en http://localhost:8080.
Opción 2: Makefile
Usá los siguientes comandos para correr la aplicación localmente con tu entorno:
make start-db     # Inicia PostgreSQL y Redis usando docker-compose
make run          # Corre la aplicación Go en modo local
Base de datos
Crea una base de datos llamada challenge_db.

Tablas
Las tablas se encuentran en challenge_be/postgres-init/ddl-challenge-be.sql.

🧪 Ejecutar tests
```bash
make test
```
📚 Documentación de la API
Una vez levantado el servicio, podés acceder a la documentación Swagger desde:

http://localhost:8080/swagger/index.html

🚀 Funcionalidades implementadas
Publicar un tweet: POST /tweets

Seguir a otro usuario: POST /follow

Ver homeTimeline de un usuario: GET /timeline/{user_id}

📌 Supuestos
No hay login: Se asume que los user_id recibidos son válidos.

Los identificadores de usuario pueden recibirse por header, parámetro o body.

La aplicación fue pensada para escalar a millones de usuarios, priorizando la lectura.

No se contempló unfollow, likes ni replies en esta etapa.

Redis es utilizado para cachear timelines y mejorar la velocidad de lectura.

Para más detalles, ver el archivo business.txt.

🏗️ Arquitectura de Alto Nivel
La solución está basada en una arquitectura Hexagonal (Ports & Adapters). Esta arquitectura permite separar claramente el dominio de la aplicación (lógica de negocio) de las interfaces externas (como HTTP, bases de datos y otros servicios). De esta forma, la aplicación es flexible y fácil de escalar.

Componentes principales
Dominio: Contiene las entidades y las reglas de negocio.

Aplicación: Define los casos de uso y coordina la interacción entre el dominio y las interfaces.

Infraestructura: Implementa los detalles concretos, como las conexiones a bases de datos, Redis y los adaptadores externos.

Interfaces: Exponen la API HTTP para interactuar con la aplicación.

⚙️ Elección de Tecnología
Go (Golang)

PostgreSQL: Base de datos relacional para almacenar tweets y follow.

Redis: Usado como cache para optimizar la lectura.

Docker y AWS ECS: Utilizados para facilitar el despliegue y la escalabilidad de la aplicación.

☁️ Despliegue en AWS
El proyecto está dockerizado y preparado para ser desplegado en AWS ECS. Puede adaptarse fácilmente a EC2 o EKS según necesidades. También puede integrarse con servicios como:

RDS (PostgreSQL)

ElastiCache (Redis)

CloudWatch para logs y métricas

📂 Estructura del proyecto
```bash
├── cmd/                # Entrada principal de la aplicación
├── internal/
│   ├── domain/         # Entidades y contratos del dominio
│   ├── application/    # Casos de uso
│   ├── infrastructure/ # Repositorios, Redis, adaptadores externos
│   └── interfaces/     # Handlers HTTP
├── docs/               # Swagger y documentación
├── docker/             # Dockerfiles, compose y configuraciones
├── Makefile
└── README.md
```

