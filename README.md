# todomvc

[![Actions/Build](https://github.com/oinume/todomvc/workflows/ci/badge.svg)](https://github.com/oinume/todomvc/actions?query=workflow%3Aci)
[![Codecov](https://codecov.io/gh/oinume/todomvc/branch/master/graph/badge.svg)](https://codecov.io/gh/oinume/todomvc)

This is a todomvc backend implementation in Go.

## How to run

### Requirements

- Docker for Mac
- docker-compose
- Go 1.14 or later
- Protobuf

### Run server

```shell script
# Setup commands like protoc
$ make setup

# Setup env vars. Specify 127.0.0.1 if you use Docker for Mac.
$ sed 's/<DOCKER_IP>/127.0.0.1/g' .env.sample > .env  

# Start MySQL and Jaeger
$ docker compose up -d

# Create tables
$ make db/goose/up

# Finally, start server
$ make restart

# Kill the server
$ make kill
```

### Create new todo
 
```shell script
$ curl -X POST -d '{"title":"My first task"}' http://127.0.0.1:5001/todos

{"id":"77b34605-8481-471a-9b92-266ec1f36486","title":"My first task"}
```

### List todos

```shell script
$ curl -X GET http://127.0.0.1:5001/todos

{"todos":[{"id":"77b34605-8481-471a-9b92-266ec1f36486","title":"My 1st task"}]}
```

### Update todo

```shell script
$ curl -X PATCH -d '{"todo":{"title":"My 2nd task"}}' http://127.0.0.1:5001/todos/<id>
```

### Delete todo

```shell script
$ curl -X DELETE http://127.0.0.1:5001/todos/<id>
```

## Backend

- Go
- Protocol: HTTP + Protocol Buffers with JSON codec
- Routing: [gorilla/mux](https://github.com/gorilla/mux)
- Database: [MySQL](https://www.mysql.com/)
- ORMapper: [volatiletech/sqlboiler](https://github.com/volatiletech/sqlboiler)
- DI: google/wire (TBD)
- Validation: [go-playground/validator](https://github.com/go-playground/validator)
- Logging: [uber-go/zap](https://github.com/uber-go/zap)
- Tracing: [OpenCensus](https://opencensus.io/) + [Jaeger](https://www.jaegertracing.io/)

## CI

- GitHub Actions
