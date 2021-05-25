### Usage

#### 1. Run the application
```shell
docker run -p 7001:7001 -p 7002:7002 ramone/rusprofile:v0.1
```
Then you can interact using gRPC or HTTP

#### 2. Try gRPC
```shell
grpc_cli call 127.0.0.1:7002 rusprofile.RusProfileService/GetCompanyInfo 'inn: "7704217370"'
```

#### 3. Try HTTP
Visit http://127.0.0.1:7001/doc/ and use Swagger UI to make a request
