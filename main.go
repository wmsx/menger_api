package main

import (
	"github.com/micro/go-micro/v2/util/log"
	"github.com/micro/go-micro/v2/web"
	"github.com/wmsx/menger_api/routers"
)

const name = "wm.sx.web.menger"

func main() {
	svc := web.NewService(
		web.Name(name),
	)

	if err := svc.Init(); err != nil {
		log.Fatal("初始化失败", err)
	}

	router := routers.InitRouter(svc.Options().Service.Client())
	svc.Handle("/", router)

	if err := svc.Run(); err != nil {
		log.Fatal("启动失败", err)
	}
}
