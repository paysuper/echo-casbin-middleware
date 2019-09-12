module github.com/paysuper/echo-casbin-middleware

require (
	github.com/InVisionApp/go-health v2.1.1-0.20181204182500-21e799eae2a9+incompatible // indirect
	github.com/casbin/casbin/v2 v2.0.2 // indirect
	github.com/go-xorm/xorm v0.7.6 // indirect
	github.com/hashicorp/consul/api v1.2.0 // indirect
	github.com/labstack/echo/v4 v4.1.6
	github.com/marten-seemann/qtls v0.3.3 // indirect
	github.com/mattn/go-sqlite3 v1.11.0 // indirect
	github.com/micro/go-micro v1.10.0
	github.com/miekg/dns v1.1.17 // indirect
	github.com/paysuper/casbin-server v0.0.0-20190912105647-fe8264c2c735
	github.com/prometheus/client_model v0.0.0-20190812154241-14fe0d1b01d4 // indirect
	github.com/prometheus/procfs v0.0.4 // indirect
	golang.org/x/crypto v0.0.0-20190911031432-227b76d455e7 // indirect
	golang.org/x/net v0.0.0-20190909003024-a7b16738d86b // indirect
	golang.org/x/sys v0.0.0-20190911201528-7ad0cfa0b7b5 // indirect
	google.golang.org/genproto v0.0.0-20190911173649-1774047e7e51 // indirect
	google.golang.org/grpc v1.23.1 // indirect
)

replace (
	github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.0
	github.com/hashicorp/consul => github.com/hashicorp/consul v1.5.2
	github.com/hashicorp/consul/api => github.com/hashicorp/consul/api v1.1.0
	github.com/lucas-clemente/quic-go => github.com/lucas-clemente/quic-go v0.12.0
	gopkg.in/DATA-DOG/go-sqlmock.v1 => github.com/DATA-DOG/go-sqlmock v1.3.3
)

go 1.13
