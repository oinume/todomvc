# todomvc

[![Actions/Build](https://github.com/oinume/todomvc/workflows/ci/badge.svg)](https://github.com/oinume/todomvc/actions?query=workflow%3Aci)
[![Codecov](https://codecov.io/gh/oinume/todomvc/branch/master/graph/badge.svg)](https://codecov.io/gh/oinume/todomvc)

This is a todomvc backend implementation in Go.

## How to run

TBD

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
