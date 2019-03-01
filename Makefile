LOCAL_PATH = $(shell pwd)

.PHONY: example proto install gen-tag test

example: proto install
	cd example && protoc -I /usr/local/include \
		-I ${LOCAL_PATH}/example/ \
		-I ${LOCAL_PATH}/ \
		--gotag_out=xxx="graphql+\"-\" bson+\"-\"":${LOCAL_PATH}/example/ ${LOCAL_PATH}/example/example.proto

proto:
	protoc -I /usr/local/include \
	-I ${LOCAL_PATH} \
	--go_out=:. ${LOCAL_PATH}/example/example.proto

install:
	go install .

gen-tag:
	protoc -I /usr/local/include \
	-I ${LOCAL_PATH} \
	--go_out=paths=source_relative:. ${LOCAL_PATH}/tagger/tagger.proto

test:
	go test ./...
