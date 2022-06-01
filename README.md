# Requirement:
- Go v1.18.
- Docker v20.10.14.
- Docker-compose v2.5.1.
# Installation

```bash
go get
```

- Run Redis, RabbitMQ, ElasticSearch, Postgres database

```bash
cd ./servers && docker-compose up
```
## Running indexer app
The app will help indexing data into elastic search.

```bash
$ go run ./packages/indexer/main.go
```

## Running tool insert fake user
The app will help get users ramdom from "https://randomuser.me/api/", then insert to Postgres database

```bash
$ go run ./packages/random_user/main.go
```

## Running the rest API local

```bash
$ go run ./packages/rest_api/main.go
```


## Build the app

```bash
# build
$ go build

# run
$ ./user_management

```

## Load test

```bash
# install package autocannon with nodejs
$ npm i -g autocannon

# run
$ autocannon localhost:8080/api/v1/user

```
