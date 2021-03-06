module github.com/wmsx/menger_api

go 1.14

//	github.com/wmsx/menger_svc => /Users/zengqiang96/codespace/sx/menger_svc
//	github.com/wmsx/xconf => /Users/zengqiang96/codespace/xconf
//  替换为v1.26.0版本的gRPC库st
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/gin-gonic/gin v1.6.3
	github.com/micro/cli/v2 v2.1.2
	github.com/micro/go-micro/v2 v2.5.0
	github.com/wmsx/menger_svc v0.0.0-20200903171825-192703a70f26
	github.com/wmsx/pkg v0.0.0-20200904032714-e9c2f482fcbc
	github.com/wmsx/xconf v0.0.0-20200720194624-34a4108a8759
)
