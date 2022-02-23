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
$ npm run start

# watch mode
$ npm run start:dev

# production mode
$ npm run start:prod
```

## Running app with docker

```push to docker cloud
# development
$ gulp

# development mode
$ docker rm -f whoscan_api_dev && docker image rm -f dotrung051093/baam:whoscan_api_dev && docker run -i -t --restart unless-stopped -d -p 5001:5001 --name whoscan_api_dev --env NODE_ENV=development --env HOST_HOSTNAME=$(hostname) --log-driver none --cpus="2.0" dotrung051093/baam:whoscan_api

# production mode
$ docker rm -f whoscan_api_prod && docker image rm -f dotrung051093/baam:whoscan_api && docker run -i -t --restart unless-stopped -d -p 5000:5000 --name whoscan_api_prod --env NODE_ENV=production --env HOST_HOSTNAME=$(hostname) --log-driver none --cpus="2.0" dotrung051093/baam:whoscan_api
```

## Test

```bash
# unit tests
$ npm run test

# e2e tests
$ npm run test:e2e

# test coverage
$ npm run test:cov
```

## Support

Nest is an MIT-licensed open source project. It can grow thanks to the sponsors and support by the amazing backers. If you'd like to join them, please [read more here](https://docs.nestjs.com/support).

## Stay in touch

- Author - [Kamil Myśliwiec](https://twitter.com/kammysliwiec)
- Website - [https://nestjs.com](https://nestjs.com/)
- Twitter - [@nestframework](https://twitter.com/nestframework)

## License

Nest is [MIT licensed](https://github.com/nestjs/nest/blob/master/LICENSE).
