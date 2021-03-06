module github.com/paysuper/echo-casbin-middleware

go 1.13

replace github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.0

require (
	cloud.google.com/go v0.44.0 // indirect
	github.com/ProtocolONE/go-core/v2 v2.1.0
	github.com/google/uuid v1.1.1
	github.com/labstack/echo/v4 v4.1.11
	github.com/mattn/go-sqlite3 v1.10.0
	github.com/micro/go-micro v1.18.0
	github.com/paysuper/casbin-server v1.0.1-0.20200203122730-7aefb4994b08 // indirect
	github.com/paysuper/paysuper-proto/go/casbinpb v0.0.0-20200203130641-45056764a1d7
	github.com/vektra/mockery v0.0.0-20181123154057-e78b021dcbb5 // indirect
)
