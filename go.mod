module github.com/ivanovaleksey/rusprofile

go 1.16

require (
	github.com/PuerkitoBio/goquery v1.6.1
	github.com/go-chi/chi v1.5.4
	github.com/golang/protobuf v1.5.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/hashicorp/golang-lru v0.5.4
	github.com/ivanovaleksey/rusprofile/pkg/pb/rusprofile v0.0.0-00010101000000-000000000000
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
	go.uber.org/zap v1.16.0
	google.golang.org/grpc v1.38.0
)

replace github.com/ivanovaleksey/rusprofile/pkg/pb/rusprofile => ./pkg/pb/rusprofile
