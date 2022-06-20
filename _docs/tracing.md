# Tracing in go:
Distributed tracing is the process of tracking the activity resulting from a request to an application. With this feature, you can:

- Trace the path of a request as it travels across a complex system
- Discover the latency of the components along that path
- Know which component in the path is creating a bottleneck
- Inspect payloads that are being sent between components
- Build execution graph for each component internals and more.

# Dependencies
- [jaegertracing](https://www.jaegertracing.io/)
- [opencensus](https://opencensus.io/)

# Demo:
- [Trace](https://trace.tdo.works/)

# Installation

## Jaeger:
- [docker-compose.yml](https://github.com/trung051093/gogo/blob/main/docker-compose.yml)

```bash
docker-compose up -d
```

## Opencensus

```bash
go get -u go.opencensus.io
```

## Code

This implementation is advanced. There are something you may know before reading it like: decorator, reflection.

My approach is that we have a pre-existing application. Now QA reports some requests very slow and that request is a computational combination of several service. You don't know which service slowly and which don't. 

So I'm try write a decorator function, this will help end the current service, doing exactly what it does.

- [Invoke func decorator](https://github.com/trung051093/gogo/blob/main/decorators/invoke.go)
The decorator function do exactly what it does.

- [Trace decorator](https://github.com/trung051093/gogo/blob/main/decorators/trace.go)
The decorator function wrap current service with a trace span.

- [Service Trace](https://github.com/trung051093/gogo/blob/main/modules/user/user_service_trace.go)
You can see I've just wrapper every User Service with a trace span. 

- [Auth Trace](https://github.com/trung051093/gogo/blob/main/modules/auth/auth_service.go)
Let check that file, the register service is computational combination of `SearchUser` service and `CreateUser` service. (Before create an user, we have to check the user is exists or not)

I've switch current service `SearchUser`, `CreateUser` to `SearchUserTrace`, `CreateUserTrace`.

- [Test](https://trace.tdo.works/trace/392c6a7843ae49667fd6fe558b6928d9)
![trace image](./assets/trace.png)
