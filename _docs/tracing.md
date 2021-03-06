# Distributed tracing:
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
- [trace.tdo.works](https://trace.tdo.works/)

# Installation

## Jaeger:
Deploy with docker: [docker-compose.yml](../docker-compose.yml)

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

So I'm try write a decorator function, this will help wrap our current service, doing exactly what it does and do not change the business logic.

- [Invoke func decorator](../decorators/invoke.go):
The decorator function help invoke a function and do exactly what it does.

- [Trace decorator](../decorators/trace.go):
The decorator function wrap current service with a trace span.

- [Service Trace](../modules/user/user_service_trace.go):
You can see I've just wrapper every User Service with a trace span. 

- [Auth Trace](../modules/auth/auth_service.go):
Let check that file, the register service is computational combination of `SearchUser` service and `CreateUser` service. (Before create an user, we have to check the user is exists or not)

I've switch current service `SearchUser`, `CreateUser` to `SearchUserTrace`, `CreateUserTrace`.

- [Test](https://trace.tdo.works/trace/392c6a7843ae49667fd6fe558b6928d9)
![trace image](../_assets/trace.png)
