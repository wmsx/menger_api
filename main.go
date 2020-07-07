package main

import (
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/util/log"
	"github.com/micro/go-micro/v2/web"
	"github.com/wmsx/menger_api/handler"
	"net/http"
)

const name = "wm.sx.web.menger"

func main() {
	svc := web.NewService(
		web.Name(name),
	)

	if err := svc.Init(); err != nil {
		log.Fatal("初始化失败", err)
	}

	router := gin.Default()
	mengerHandler := handler.New(svc.Options().Service.Client())
	r := router.Group("/menger")
	r.GET("/login", mengerHandler.Login)
	r.GET("/xx", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{
			"xx": "success",
		})
	})

	svc.Handle("/", router)

	if err := svc.Run(); err != nil {
		log.Fatal("启动失败", err)
	}
}
