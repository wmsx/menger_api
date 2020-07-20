module github.com/wmsx/menger_api

go 1.14

require (
	github.com/gin-gonic/gin v1.6.3
	github.com/micro/cli/v2 v2.1.2
	github.com/micro/go-micro/v2 v2.9.1
	github.com/wmsx/menger_svc v0.0.0-20200720105758-ea97af9f3c81
	github.com/wmsx/pkg v0.0.0-20200720153510-e000d75295a3
	github.com/wmsx/xconf v0.0.0-20200710193800-f97c7e3c9e84
)

//  替换为v1.26.0版本的gRPC库
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
