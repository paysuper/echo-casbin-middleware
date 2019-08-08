module github.com/paysuper/echo-casbin-middleware

require (
	github.com/labstack/echo/v4 v4.1.6
	github.com/micro/go-micro v1.7.1-0.20190724203029-7ca8f8f0ab98
	github.com/unistack-org/casbin-micro v0.0.0-20190725231001-e35616531b31
)

replace (
	github.com/hashicorp/consul => github.com/hashicorp/consul v1.5.2
	github.com/hashicorp/consul/api => github.com/hashicorp/consul/api v1.1.0
)
