# Installation

```bash
go get
```

- Create a postgres localhost

```bash
docker run --restart unless-stopped -d -p 5432:5432 --name postgres -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -e POSTGRES_DB=example  postgres
```

- Create a rabbitmq localhost

```bash
docker run --restart unless-stopped -d -p 5432:5432 --name postgres -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -e POSTGRES_DB=example  postgres
```

- Create an elastic search localhost

```bash
docker run --restart unless-stopped -d -p 5432:5432 --name postgres -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -e POSTGRES_DB=example  postgres
```

## Running the app

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
