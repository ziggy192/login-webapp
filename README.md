# README #

A Profile management site that integrate with Google Sign In with high security level

### How to run in local ? ###

Start docker compose for database and redis

```shell
make up
```

Start api server

```shell
make api
```

Start frontend server

```shell
make frontend
```

Verify the apps are running

api

```shell
$ curl localhost:8080/ping

{"message":"pong","time":"2022-11-22T19:29:03+07:00","data":{"host":"localhost:8080","port":"8080"}}%  
```

frontend

```shell
»  curl localhost:9090/ping

{"message":"pong","time":"2022-11-22T19:30:18+07:00","data":{"host":"localhost:9090","port":"9090"}}% 
```

### How to run tests

Start docker compose for db and redis

```shell
make up
```

Run tests

```shell
make test
```

### Technical decisions

**Separate backend REST API and frontend**. This clear separation helps the API stand independently if we were to change
the frontend into something entirely different. API is decoupled with Frontend, making them independent to code changes,
feature changes, or just redeployment of the other service.

**Authentication with JWT**. JWT token is a light and simple way to authenticate user. JWT token also match with
microservice architecture (if we were to go this direction in the future). The token is signed by simple HMAC secret,
but in the future we can change to RSA and safely share the public key with other services.

**Clear cookie and no-cache**. This is to prevent the back button, or any weir behaviors on the browser to access a
log-outed resource.

**Block token when logout**. When user with a token log-out. I block the token as a black-list mechanism using redis.
When the auth
middleware verify a token, it also checks redis if this token is blocked. This to make sure a blocked/log-outed token
can
never be valid

**Unit test with real database**. There are 2 styles when doing unit test with a database: mocked database, or real
database. Each style has its own pros and cons. With mocked database, our code have to depend on a database interface
with all the functions of the real database, and generate mocks or stubs with this interface in testing. With real
database, we dont need to declare any interface, or have to regenerate mocks when changing/adding functions. But this
style requires seeding data before tests. I chose to go with real database.

**Context propagation**. Each function support a `context.Context`. This practice
support [Go Concurrency Pattern with Context](https://go.dev/blog/context)
and [OpenTelemetry](https://opentelemetry.io/) tracing using a `request-id`.

**Use format and linter** as a passive way to follow standard Golang coding conventions.

### Contacts ###

nghialuu.it@gmail.com
