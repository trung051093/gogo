# Requirement:
- Go v1.18.
- Docker v20.10.14.
- Docker-compose v2.5.1.

# Installation

```bash
go get
```

# Running with docker compose

```bash
# build rest_api package
$ docker build -t api:multistage -f Dockerfile.multistage.api .

# tag rest_api to api:v1.0
$ docker image tag api:multistage api:v1.0

# build indexer package
$ docker build -t indexer:multistage -f Dockerfile.multistage.indexer .

# tag rest_api to api:v1.0
$ docker image tag indexer:multistage indexer:v1.0

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
The app will help get users ramdom from "https://randomuser.me/api/", then insert to Postgres database
Require: RabbitMQ, Postgres

```bash
$ go run ./packages/random_user/main.go
```

## Running the rest API local
Require: RabbitMQ, Postgres, ElasticSearch

```bash
$ go run ./packages/rest_api/main.go
```

## Load test

```bash
# install package autocannon with nodejs
$ npm i -g autocannon

# run
$ autocannon localhost:8080/api/v1/user

```
