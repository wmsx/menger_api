package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/client"
	"github.com/wmsx/menger_api/handler"
	mygin "github.com/wmsx/pkg/gin"
)

/**
 * 初始化路由信息
 */
func InitRouter(c client.Client) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	mengerHandler := handler.NewMengerHandler(c)
	mengerRouter := r.Group("/menger")

	mengerRouter.POST("/Register", mygin.Wrapper(mengerHandler.Register))
	mengerRouter.POST("/login", mygin.Wrapper(mengerHandler.Login))
	mengerRouter.POST("/logout",  mygin.Wrapper(mengerHandler.Logout))

	return r
}
