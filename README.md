# todomvc
[
[![Actions/Build](https://github.com/oinume/todomvc/workflows/ci/badge.svg)](https://github.com/oinume/todomvc/actions?query=workflow%3Aci)

This is a todomvc backend implementation in Go.

## Backend

- Go
- Protocol: HTTP + Protocol Buffers with JSON codec
- Routing: [gorilla/mux](https://github.com/gorilla/mux)
- Database: [MySQL](https://www.mysql.com/)
- ORMapper: [volatiletech/sqlboiler](https://github.com/volatiletech/sqlboiler)
- DI: google/wire (TBD)
- Logging: [uber-go/zap](https://github.com/uber-go/zap)
- Tracing: [OpenCensus](https://opencensus.io/) + [Jaeger](https://www.jaegertracing.io/)

## CI

- GitHub Actions
