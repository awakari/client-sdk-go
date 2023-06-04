.PHONY: test clean
default: build

COVERAGE_FILE_NAME=cover.out
COVERAGE_TMP_FILE_NAME=cover.tmp

proto:
	go install github.com/golang/protobuf/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
	PATH=${PATH}:~/go/bin protoc --go_out=plugins=grpc:. --go_opt=paths=source_relative \
		api/grpc/limits/*.proto \
		api/grpc/permits/*.proto \
		api/grpc/reader/*.proto \
		api/grpc/subject/*.proto \
		api/grpc/subscriptions/*.proto \
		api/grpc/writer/*.proto

vet: proto
	go vet

test: vet
	go test -race -cover -coverprofile=${COVERAGE_FILE_NAME} ./...
	cat ${COVERAGE_FILE_NAME} | grep -v _mock.go | grep -v logging.go | grep -v .pb.go > ${COVERAGE_FILE_NAME}.tmp
	mv -f ${COVERAGE_FILE_NAME}.tmp ${COVERAGE_FILE_NAME}
	go tool cover -func=${COVERAGE_FILE_NAME} | grep -Po '^total\:\h+\(statements\)\h+\K.+(?=\.\d+%)' > ${COVERAGE_TMP_FILE_NAME}
	./scripts/cover.sh
	rm -f ${COVERAGE_TMP_FILE_NAME}

clean:
	go clean
	rm -f ${COVERAGE_FILE_NAME}
