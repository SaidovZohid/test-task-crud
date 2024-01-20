### Golang [Golang BLOG app](https://task.zohiddev.me) 🌟

#### 👨‍💻 Full list what has been used:
* [gin](https://github.com/gin-gonic/gin) - Web framework
* [gin-cors](https://github.com/gin-contrib/cors) - Gin middleware/handler to enable CORS support.
* [swag](https://github.com/swaggo/swag) - Swagger
* [sqlx](https://github.com/jmoiron/sqlx) - Extensions to database/sql.
* [pq](https://github.com/lib/pq) - Go postgres driver for Go's database/sql package
* [go-redis](https://github.com/redis/go-redis) - Redis client for Go
* [logrus](https://github.com/sirupsen/logrus) - Logger
* [jwt-go](https://github.com/golang-jwt/jwt) - JSON Web Tokens (JWT)
* [migrate](https://github.com/golang-migrate/migrate) - Database migrations. CLI and Golang library.
* [bluemonday](https://github.com/microcosm-cc/bluemonday) - HTML sanitizer
* [faker](https://github.com/bxcodec/faker) - Faker will generate you a fake data based on your Struct.
* [testify](https://github.com/stretchr/testify) - Testing toolkit
* [Docker](https://www.docker.com/) - Docker

#### Recomendation for local development most comfortable usage:
```
    make up // run all containers
    make run // it's easier way to attach debugger or rebuild/rerun project
```
#### 🙌👨‍💻🚀 Docker-compose files:
    docker-compose.yml - run docker development environment

### Docker development usage:
```
    make docker
```
### Local development usage:
```
    make run
```
### Local test:
```
    make test
```
### SWAGGER UI:

# If you run locally:
https://localhost:8080/swagger/index.html

# Checkout:
https://task.zohiddev.me/swagger/index.html

#### Note for Docker Setup:
Before running `make up`, make sure to copy the `sample.env.docker` file to `.env.docker` and update the necessary secure fields.
#### Note for Local Development:
Before running `make run`, make sure to copy the `sample.env` file to `.env` and update the necessary secure fields. Additionally, you need to run the Redis and PostgreSQL containers.