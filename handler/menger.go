package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/util/log"
	menger "github.com/wmsx/menger_svc/proto/menger"
	"net/http"
)

const mengerSvcName = "wm.sx.svc.menger"

type MengerApiService struct {
	mengerClient menger.MengerService
}

func New(client client.Client) *MengerApiService {
	return &MengerApiService{
		mengerClient: menger.NewMengerService(mengerSvcName, client),
	}
}

func (m *MengerApiService) Login(c *gin.Context) {
	var (
		loginRes *menger.LoginResponse
		err      error
	)

	if loginRes, err = m.mengerClient.Login(context.Background(), &menger.LoginRequest{}); err != nil {
		log.Fatal("调用wm.sx.svc.menger失败", err)
		c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "调用wm.sx.svc.menger失败",
		})
		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"message": loginRes.Msg,
	})
}
