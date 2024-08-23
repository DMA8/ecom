CURDIR=$(shell pwd)
BINDIR=${CURDIR}/bin
GOVER=$(shell go version | perl -nle '/(go\d\S+)/; print $$1;')
SMARTIMPORTS=${BINDIR}/smartimports_${GOVER}
LINTVER=v1.51.1
LINTBIN=${BINDIR}/lint_${GOVER}_${LINTVER}
PACKAGE=ecom/cart/cmd/app
PSQL_DNS=postgres://postgres:our_password@localhost:5432/postgres

all: format build test lint

build: bindir
	GOOS=linux GOARCH=amd64 go build -o ${BINDIR}/app ${PACKAGE}

test:
	go test ./...

run:
	go run ${PACKAGE}

lint: install-lint
	${LINTBIN} run

precommit: format build test lint
	echo "OK"

bindir:
	mkdir -p ${BINDIR}

format: install-smartimports
	${SMARTIMPORTS} -exclude internal/mocks

install-lint: bindir
	test -f ${LINTBIN} || \
		(GOBIN=${BINDIR} go install github.com/golangci/golangci-lint/cmd/golangci-lint@${LINTVER} && \
		mv ${BINDIR}/golangci-lint ${LINTBIN})

install-smartimports: bindir
	test -f ${SMARTIMPORTS} || \
		(GOBIN=${BINDIR} go install github.com/pav5000/smartimports/cmd/smartimports@latest && \
		mv ${BINDIR}/smartimports ${SMARTIMPORTS})

protogen:
	protoc  --go_out=pkg --go-grpc_out=pkg api/cart/v1/cart.proto
	go mod tidy

# Используем bin в текущей директории для установки плагинов protoc
LOCAL_BIN:=$(CURDIR)/bin

# Добавляем bin в текущей директории в PATH при запуске protoc
PROTOC = PATH="$$PATH:$(LOCAL_BIN)" protoc

# Установка всех необходимых зависимостей
.PHONY: .bin-deps
.bin-deps:
	$(info Installing binary dependencies...)

	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest && \
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest && \
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@latest && \
	GOBIN=$(LOCAL_BIN) go install github.com/bufbuild/buf/cmd/buf@v1.21.0


# Вендоринг внешних proto файлов
.vendor-proto: .vendor-rm vendor-proto/google/api vendor-proto/google/protobuf vendor-proto/protoc-gen-openapiv2/options vendor-proto/validate 

.PHONY: .vendor-rm
.vendor-rm:
	rm -rf vendor-proto

# Устанавливаем proto описания google/googleapis
vendor-proto/google/api:
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/googleapis/googleapis vendor-proto/googleapis && \
 	cd vendor-proto/googleapis && \
	git sparse-checkout set --no-cone google/api && \
	git checkout
	mkdir -p  vendor-proto/google
	mv vendor-proto/googleapis/google/api vendor-proto/google
	rm -rf vendor-proto/googleapis

# Устанавливаем proto описания protoc-gen-openapiv2/options
vendor-proto/protoc-gen-openapiv2/options:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/grpc-ecosystem/grpc-gateway vendor-proto/grpc-ecosystem && \
 	cd vendor-proto/grpc-ecosystem && \
	git sparse-checkout set --no-cone protoc-gen-openapiv2/options && \
	git checkout
	mkdir -p vendor-proto/protoc-gen-openapiv2
	mv vendor-proto/grpc-ecosystem/protoc-gen-openapiv2/options vendor-proto/protoc-gen-openapiv2
	rm -rf vendor-proto/grpc-ecosystem

	
# Устанавливаем proto описания google/protobuf
vendor-proto/google/protobuf:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf vendor-proto/protobuf &&\
	cd vendor-proto/protobuf &&\
	git sparse-checkout set --no-cone src/google/protobuf &&\
	git checkout
	mkdir -p vendor-proto/google
	mv vendor-proto/protobuf/src/google/protobuf vendor-proto/google
	rm -rf vendor-proto/protobuf

# Устанавливаем proto описания validate
vendor-proto/validate:
	git clone -b main --single-branch --depth=2 --filter=tree:0 \
		https://github.com/bufbuild/protoc-gen-validate vendor-proto/tmp && \
		cd vendor-proto/tmp && \
		git sparse-checkout set --no-cone validate &&\
		git checkout
		mkdir -p vendor-proto/validate
		mv vendor-proto/tmp/validate vendor-proto/
		rm -rf vendor-proto/tmp


SERVICE_PROTO_PATH:="api/cart/v1"
PRODUCTSERVICE_PROTO_PATH="api/productService"

# Генерация протофайлов с использованием protoc
PHONY: .protoc-generate
#  .bin-deps .vendor-proto
.protoc-generate:
	mkdir -p pkg/${SERVICE_PROTO_PATH}
	$(PROTOC) -I ${SERVICE_PROTO_PATH} -I vendor-proto \
	--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go --go_out pkg/${SERVICE_PROTO_PATH} --go_opt paths=source_relative \
	--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc --go-grpc_out pkg/${SERVICE_PROTO_PATH} --go-grpc_opt paths=source_relative \
	--plugin=protoc-gen-validate=$(LOCAL_BIN)/protoc-gen-validate --validate_out="lang=go,paths=source_relative:pkg/api/cart/v1" \
	api/cart/v1/cart.proto

protoc-productservice:
	mkdir -p pkg/${PRODUCTSERVICE_PROTO_PATH}
	$(PROTOC) -I ${PRODUCTSERVICE_PROTO_PATH} -I vendor-proto \
	--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go --go_out pkg/${PRODUCTSERVICE_PROTO_PATH} --go_opt paths=source_relative \
	--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc --go-grpc_out pkg/${PRODUCTSERVICE_PROTO_PATH} --go-grpc_opt paths=source_relative \
	--plugin=protoc-gen-validate=$(LOCAL_BIN)/protoc-gen-validate --validate_out="lang=go,paths=source_relative:pkg/api/productService" \
	api/productService/productService.proto

# Генерация протофайлов с использованием protoc
PHONY: generate
generate: .bin-deps .vendor-proto .protoc-generate

.PHONY: fast-generate
fast-generate: .protoc-generate

# Генерация протофайлов с использованием buf
PHONY: generate-buf
generate-buf: .bin-deps .vendor-rm .buf-generate

# Генерация протофайлов с использованием buf
.PHONY: .buf-generate
.buf-generate: 
	PATH=$(LOCAL_BIN):$(PATH) $(LOCAL_BIN)/buf mod update
	PATH=$(LOCAL_BIN):$(PATH) $(LOCAL_BIN)/buf generate

migrate-up:
	goose -dir migrations postgres $(PSQL_DNS) up

migrate-down:
	goose -dir migrations postgres $(PSQL_DNS) down
