<h1>Golang REST API boiterplate</h1>

# Requirement:
- Go v1.18.
- Docker v20.10.14.
- Docker-compose v2.5.1.

# Docs:
- [Distributed tracing](./_docs/tracing.md)
- [Logging](./_docs/logging.md)
- [Swagger API](./_docs/swagger.md)
- [Docker build](./_docs/build.md)
- [Send mail via Gmail SMTP server](./_docs/mail.md)
- [Sign in with google](./_docs/google-login.md)

# Demo:
- [Api](https://api.tdo.works/swagger/index.html#/)
- [Trace](https://trace.tdo.works/)
- [Graylog](https://graylog.tdo.works/):
```bash
user: anonymous
pass: anonymous
```

# Installation

```bash
go mod download
go mod tidy
```

# Running with docker compose

```bash
# build rest_api package
docker build -t api:multistage -f Dockerfile.multistage.api .

# tag api
docker image tag api:multistage api

# build indexer package
docker build -t indexer:multistage -f Dockerfile.multistage.indexer .

# tag indexer
docker image tag indexer:multistage indexer

# run api with Redis, RabbitMQ, ElasticSearch, Postgres database
docker-compose up
```

## Running Redis, RabbitMQ, ElasticSearch, Postgres database

```bash
# run api with Redis, RabbitMQ, ElasticSearch, Postgres database
docker-compose up redis rabbitmq elasticsearch postgres
```

## Running indexer app
The app will help indexing data into elastic search.

Require: RabbitMQ, ElasticSearch

```bash
go run ./packages/indexer/main.go
```

## Running tool insert fake user
The app will help get users ramdom from [randomuser.me](https://randomuser.me/api/), then insert to Postgres database.

Require: RabbitMQ, Postgres

```bash
go run ./packages/random_user/main.go
```

## Running the rest API local
Require: RabbitMQ, Postgres, ElasticSearch

```bash
go run ./packages/rest_api/main.go
```

## Load test: [Wrk](https://github.com/wg/wrk)

### Cache 
```bash
wrk -t6 -c100 -d15s https://api.tdo.works/api/v1/users-cache
```

```bash
Running 15s test @ https://api.tdo.works/api/v1/users-cache
  6 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   227.11ms   17.19ms 318.56ms   76.39%
    Req/Sec    70.48     21.42   141.00     73.44%
  6044 requests in 15.10s, 13.25MB read
Requests/sec:    400.35
Transfer/sec:      0.88MB
```

### No Cache 
```bash
wrk -t6 -c100 -d15s https://api.tdo.works/api/v1/users
```

```bash
Running 15s test @ https://api.tdo.works/api/v1/users
  6 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   581.02ms  217.44ms   1.56s    68.91%
    Req/Sec    27.78     15.62    89.00     65.57%
  2304 requests in 15.03s, 5.05MB read
Requests/sec:    153.25
Transfer/sec:    343.92KB
```
