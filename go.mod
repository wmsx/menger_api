module github.com/wmsx/menger_api

go 1.14

require (
	github.com/gin-gonic/gin v1.6.3
	github.com/micro/go-micro/v2 v2.4.0
	github.com/wmsx/menger_svc v0.0.0-20200707110412-b3462a19f1b7
	google.golang.org/protobuf v1.25.0 // indirect
)

// 替换为v1.26.0版本的gRPC库
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
