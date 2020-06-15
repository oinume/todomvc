APP = todomvc
VENDOR_DIR = vendor
PROTO_GEN_DIR = proto-gen
GRPC_GATEWAY_REPO = github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis
GO_GET ?= go get -u
GO_TEST ?= go test -v -race
VENDOR_DIR = vendor
PID = server.pid

all: build

.PHONY: setup
setup: install-commands install-tools

.PHONY: install-commands
install-commands:
	$(GO_GET) google.golang.org/protobuf/cmd/protoc-gen-go@v1.24.0
#	$(GO_GET) github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
#	$(GO_GET) https://github.com/golang/protobuf/tree/master/protoc-gen-go@v1.4.2

.PHONY: install-tools
install-tools:
	cd tools && $(GO_GET) \
		github.com/pressly/goose/cmd/goose

.PHONY: build
build:
	go build -o bin/server github.com/oinume/$(APP)/backend/cmd/server

clean:
	${RM} bin/server

test:
	$(GO_TEST) ./...

.PHONY: proto/go
proto/go:
	rm -rf $(PROTO_GEN_DIR)/go && mkdir -p $(PROTO_GEN_DIR)/go
	protoc -I/usr/local/include -I. -I./proto -I./proto/third-party \
  		-I$(GOPATH)/src \
  		--go_out=$(PROTO_GEN_DIR)/go \
  		--go_opt=paths=source_relative \
		--experimental_allow_proto3_optional \
  		./proto/todomvc/*.proto
#	protoc -I/usr/local/include -I. -I./proto -I./proto/third-party \
#		-I$(GOPATH)/src \
#		-I$(VENDOR_DIR)/$(GRPC_GATEWAY_REPO) \
#		--grpc-gateway_out=logtostderr=true:$(PROTO_GEN_DIR)/go \
		proto/api/v1/*.proto

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


kill:
	@kill `cat $(PID)` 2> /dev/null || true

restart: kill clean build
	bin/server & echo $$! > $(PID)

watch: restart
	fswatch -o -e ".*" -e vendor -e node_modules -e .venv -i "\\.go$$" . | xargs -n1 -I{} make restart || make kill

