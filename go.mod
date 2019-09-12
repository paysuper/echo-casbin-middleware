module github.com/paysuper/echo-casbin-middleware

require (
	github.com/labstack/echo/v4 v4.1.6
	github.com/micro/go-micro v1.10.0
	github.com/paysuper/casbin-server v0.0.0-20190912100652-3aad9c05777d
//github.com/paysuper/casbin-server v0.0.0-20190817124042-6b2b224a39ab
)

replace (
	github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.0
	github.com/hashicorp/consul => github.com/hashicorp/consul v1.5.2
	github.com/hashicorp/consul/api => github.com/hashicorp/consul/api v1.1.0
	github.com/lucas-clemente/quic-go => github.com/lucas-clemente/quic-go v0.11.0
	gopkg.in/DATA-DOG/go-sqlmock.v1 => github.com/DATA-DOG/go-sqlmock v1.3.3
//github.com/unistack-org/casbin-micro => github.com/paysuper/casbin-server v0.0.0-20190817124042-6b2b224a39ab
)

go 1.13
