# Installation

```bash
go get gorm.io/gorm@v1.20.11 
```

```bash
go get gorm.io/driver/postgres
```

- Create a redis localhost

```bash
docker run --restart unless-stopped -d -p 6379:6379 --name redis redis
```

- Create a postgres localhost

```bash
docker run --restart unless-stopped -d -p 5432:5432 --name postgres -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -e POSTGRES_DB=example  postgres
```

- Create a minio localhost

```bash
docker run --restart unless-stopped -d -p 9000:9000 -e MINIO_ROOT_USER={MINIO_ROOT_USER} -e MINIO_ROOT_PASSWORD={MINIO_ROOT_PASSWORD} minio/minio server /data
```

## Running the app

- Create a minio localhost

```bash
# development
$ go build

# watch mode
$ ./user_management

```
