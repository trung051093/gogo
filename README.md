# Installation

```bash
go get
```

- Create a postgres localhost

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
