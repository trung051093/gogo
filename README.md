# Requirement:
- Go v1.18.
- Docker v20.10.14.
- Docker-compose v2.5.1.

# Docs:
- [Distributed tracing](./_docs/tracing.md)
- [Logging](./_docs/logging.md)
- [Swagger API](./_docs/swagger.md)
- [Docker build](./_docs/build.md)

# Demo:
- [Api](https://api.tdo.works/swagger/index.html#/)
- [Trace](https://trace.tdo.works/)
- [Graylog](https://graylog.tdo.works/)
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
$ docker build -t api:multistage -f Dockerfile.multistage.api .

# tag api
$ docker image tag api:multistage api

# build indexer package
$ docker build -t indexer:multistage -f Dockerfile.multistage.indexer .

# tag indexer
$ docker image tag indexer:multistage indexer

# run api with Redis, RabbitMQ, ElasticSearch, Postgres database
$ docker-compose up
```

## Running Redis, RabbitMQ, ElasticSearch, Postgres database

```bash
# run api with Redis, RabbitMQ, ElasticSearch, Postgres database
$ docker-compose up redis rabbitmq elasticsearch postgres
```

## Running indexer app
The app will help indexing data into elastic search.

Require: RabbitMQ, ElasticSearch

```bash
$ go run ./packages/indexer/main.go
```

## Running tool insert fake user
The app will help get users ramdom from "https://randomuser.me/api/", then insert to Postgres database.

Require: RabbitMQ, Postgres

```bash
$ go run ./packages/random_user/main.go
```

## Running the rest API local
Require: RabbitMQ, Postgres, ElasticSearch

```bash
$ go run ./packages/rest_api/main.go
```

## Load test:
[Wrk](https://github.com/wg/wrk)

```bash
wrk -t6 -c200 -d30s https://api.tdo.works/api/v1/users
```

```bash
Running 30s test @ https://api.tdo.works/api/v1/users
  6 threads and 200 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   242.63ms   94.15ms 805.17ms   88.87%
    Req/Sec   137.92     58.51   316.00     70.20%
  24025 requests in 30.08s, 49.76MB read
Requests/sec:    798.62
Transfer/sec:      1.65MB
```
