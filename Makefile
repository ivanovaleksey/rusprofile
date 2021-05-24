MODULE=github.com/ivanovaleksey/rusprofile
SWAGGER_VERSION = 2.2.10

.PHONY: build-server
build-server:
	go build -o ./bin/server ./cmd/server

.PHONY: test
test:
	go test -v -count 1 ./...

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
		--swagger_out=doc \
		api.proto
	mv doc/api.swagger.json doc/swagger.json

.PHONY: doc
doc:
	mkdir -p doc

.PHONY: swagger-ui
swagger-ui: doc
	curl -s https://codeload.github.com/swagger-api/swagger-ui/tar.gz/v$(SWAGGER_VERSION) | tar xzv -C doc swagger-ui-$(SWAGGER_VERSION)/dist
	mv -f doc/swagger-ui-$(SWAGGER_VERSION)/dist/* doc/
	rm -rf doc/swagger-ui-$(SWAGGER_VERSION) doc/swagger-ui.js
	sed -i_ "s/swagger-ui\.js/swagger-ui\.min\.js/" doc/index.html
	sed -i_ "s/http:\/\/petstore\.swagger\.io\/v2\///" doc/index.html
	rm -f doc/*_

.PHONY: clean
clean:
	rm -rf doc
