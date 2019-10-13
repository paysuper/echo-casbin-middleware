module github.com/paysuper/echo-casbin-middleware

go 1.12

replace (
	github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.0
	github.com/hashicorp/consul => github.com/hashicorp/consul v1.5.2
	github.com/hashicorp/consul/api => github.com/hashicorp/consul/api v1.1.0
	github.com/lucas-clemente/quic-go => github.com/lucas-clemente/quic-go v0.12.0
	gopkg.in/DATA-DOG/go-sqlmock.v1 => github.com/DATA-DOG/go-sqlmock v1.3.3
)

require (
	github.com/google/uuid v1.1.1
	github.com/labstack/echo/v4 v4.1.6
	github.com/mattn/go-sqlite3 v1.10.0
	github.com/micro/go-micro v1.8.0
	github.com/miekg/dns v1.1.17 // indirect
	github.com/paysuper/casbin-server v0.0.0-20190918141556-1ea0c75fd623
	golang.org/x/crypto v0.0.0-20190911031432-227b76d455e7 // indirect
	golang.org/x/net v0.0.0-20190909003024-a7b16738d86b // indirect
	golang.org/x/sys v0.0.0-20190911201528-7ad0cfa0b7b5 // indirect
	google.golang.org/genproto v0.0.0-20190911173649-1774047e7e51 // indirect
)
