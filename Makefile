APP = todomvc
PROTO_GEN_DIR = proto-gen
VENDOR_DIR = vendor
GO_GET ?= go get -u
GO_TEST ?= go test -v -race
PID = server.pid

all: build

.PHONY: setup
setup: install-tools

.PHONY: install-tools
install-tools: ## install dependent tools
	cd tools && go list -f='{{ join .Imports "\n" }}' ./tools.go | tr -d [ | tr -d ] | xargs -I{} go install {}

.PHONY: build
build:
	go build -o bin/server github.com/oinume/$(APP)/backend/cmd/server

clean:
	${RM} bin/server

test:
	$(GO_TEST) ./...

.PHONY: go/lint
go/lint:
	golangci-lint version
	golangci-lint run -j 4 --out-format=line-number ./...

.PHONY: test/db/goose/%
test/db/goose/%:
	goose -dir ./db/migration mysql "$(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(MYSQL_HOST):$(MYSQL_PORT))/todomvc_test?charset=utf8mb4&parseTime=true&loc=UTC" $*

.PHONY: test/coverage
test/coverage:
	$(GO_TEST) -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: proto/go
proto/go:
	rm -rf $(PROTO_GEN_DIR)/go && mkdir -p $(PROTO_GEN_DIR)/go
	protoc -I/usr/local/include -I. -I./proto -I./proto/third_party \
  		-I$(GOPATH)/src \
  		--go_out=$(PROTO_GEN_DIR)/go \
  		--go_opt=paths=source_relative \
		--experimental_allow_proto3_optional \
  		./proto/todomvc/*.proto

.PHONY: db/goose/%
db/goose/%:
	goose -dir ./db/migration mysql "$(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(MYSQL_HOST):$(MYSQL_PORT))/$(MYSQL_DATABASE)?charset=utf8mb4&parseTime=true&loc=UTC" $*
#	goose -dir ./db/migration mysql "$(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(MYSQL_HOST):$(MYSQL_PORT_XO))/$(MYSQL_DATABASE)?charset=utf8mb4&parseTime=true&loc=UTC" $*

.PHONY: db/reset
db/reset:
	mysql -h $(MYSQL_HOST) -P $(MYSQL_PORT) -uroot -proot -e "DROP DATABASE IF EXISTS $(MYSQL_DATABASE); DROP DATABASE IF EXISTS $(MYSQL_DATABASE)_test"
	mysql -h $(MYSQL_HOST) -P $(MYSQL_PORT) -uroot -proot < db/docker-entrypoint-initdb.d/create_database.sql

.PHONY: db/connect
db/connect:
	mysql -h $(MYSQL_HOST) -P $(MYSQL_PORT) -u$(MYSQL_USER) -p$(MYSQL_PASSWORD) $(MYSQL_DATABASE)

db/generate:
	go run ./tools/cmd/sqlboiler/main.go > sqlboiler.toml
	sqlboiler -c sqlboiler.toml mysql

run: clean build
	bin/server

kill:
	@kill `cat $(PID)` 2> /dev/null || true

restart: kill clean build
	bin/server & echo $$! > $(PID)

watch: restart
	fswatch -o -e ".*" -e vendor -e node_modules -e .venv -i "\\.go$$" . | xargs -n1 -I{} make restart || make kill

.PHONY: help
help:  ## show this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[\/a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
