module github.com/wmsx/menger_api

go 1.14

require (
	github.com/gin-gonic/gin v1.6.3
	github.com/micro/go-micro/v2 v2.9.1
	github.com/wmsx/menger_svc v0.0.0-20200707110412-b3462a19f1b7
	github.com/wmsx/pkg v0.0.0-00010101000000-000000000000
)

// 替换为v1.26.0版本的gRPC库
replace (
	github.com/wmsx/menger_svc => /Users/zengqiang96/codespace/sx/menger_svc
	github.com/wmsx/pkg => /Users/zengqiang96/codespace/sx/pkg
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)
