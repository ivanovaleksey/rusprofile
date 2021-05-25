FROM golang:1.16-alpine AS builder

RUN apk add --update make

WORKDIR /build

COPY go.* .
COPY pkg/pb/rusprofile/go.* pkg/pb/rusprofile/
RUN go mod download

COPY Makefile Makefile
RUN make install-proto-tools

COPY api.proto api.proto
COPY app/ app/
COPY cmd/ cmd/
COPY pkg/ pkg/
RUN make build-server

FROM alpine

WORKDIR /app

COPY --from=builder /build/bin/server .
COPY doc/ doc/

CMD ["/app/server"]
