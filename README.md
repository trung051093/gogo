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
wrk -t6 -c150 -d15s https://api.tdo.works/api/v1/users
```

```bash
Running 15s test @ https://api.tdo.works/api/v1/users
  6 threads and 150 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   689.33ms  356.31ms   1.83s    55.45%
    Req/Sec    33.83     18.60   100.00     69.51%
  2838 requests in 15.01s, 4.19MB read
Requests/sec:    189.02
Transfer/sec:    285.93KB
```