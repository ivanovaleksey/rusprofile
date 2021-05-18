MODULE=github.com/ivanovaleksey/rusprofile

.PHONY: build-server
build-server:
	go build -o ./bin/server ./cmd/server

.PHONY: install-proto
install-proto:
	go install \
		github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
		github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger \
		github.com/golang/protobuf/protoc-gen-go

.PHONY: generate-proto
generate-proto:
	protoc \
		-I . \
		-I `go list -f '{{ .Dir }}' -m github.com/grpc-ecosystem/grpc-gateway`/third_party/googleapis \
		--go_out=. --go_opt=module=${MODULE} \
		--go-grpc_out=. --go-grpc_opt=module=${MODULE} \
		--grpc-gateway_out=. --grpc-gateway_opt logtostderr=true --grpc-gateway_opt module=${MODULE} --grpc-gateway_opt generate_unbound_methods=true \
		--swagger_out=pkg/pb/rusprofile \
		api.proto
