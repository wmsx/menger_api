package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/util/log"
	menger "github.com/wmsx/menger_svc/proto/menger"
	mygin "github.com/wmsx/pkg/gin"
	"net/http"
)

const mengerSvcName = "wm.sx.svc.menger"

type LoginInfo struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password" binding:"required"`
}

type RegisterInfo struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Avatar   string `json:"avatar"`
}

type MengerHandler struct {
	mengerClient menger.MengerService
}

func NewMengerHandler(client client.Client) *MengerHandler {
	return &MengerHandler{
		mengerClient: menger.NewMengerService(mengerSvcName, client),
	}
}

func (m *MengerHandler) Register(c *gin.Context) error {
	var (
		registerRes  *menger.RegisterResponse
		loginRes     *menger.LoginResponse
		err          error
		registerInfo RegisterInfo
	)

	if err = c.ShouldBindJSON(&registerInfo); err != nil {
		log.Error("参数校验失败", err)
		return mygin.LogicError("参数不完善")
	}

	registerRequest := &menger.RegisterRequest{
		Name:     registerInfo.Name,
		Email:    registerInfo.Email,
		Password: registerInfo.Password,
		Avatar:   registerInfo.Avatar,
	}

	cxt := context.Background()
	if registerRes, err = m.mengerClient.Register(cxt, registerRequest); err != nil {
		log.Error("注册失败", err)
		return mygin.ServerError()
	}
	if registerRes.ErrorMsg != nil {
		return mygin.LogicError(registerRes.ErrorMsg.Msg)
	}

	log.Info("注册成功", registerRequest)

	loginReq := &menger.LoginRequest{
		Name:     registerInfo.Name,
		Email:    registerInfo.Email,
		Password: registerInfo.Password,
	}

	loginRes, err = m.mengerClient.Login(cxt, loginReq)
	if err != nil {
		log.Error("自动登陆失败 err:", err)
		return mygin.LogicError("自动登陆失败")
	}
	if loginRes.ErrorMsg != nil {
		log.Error("自动登陆失败 errorMsg: ", loginRes.ErrorMsg)
		return mygin.LogicError("自动登陆失败")
	}

	c.JSON(http.StatusOK, gin.H{
		"token":      loginRes.Token,
		"mengerInfo": loginRes.MengerInfo,
	})
	return nil
}

func (m *MengerHandler) Login(c *gin.Context) error {
	var (
		loginRes  *menger.LoginResponse
		err       error
		loginInfo *LoginInfo
	)

	if err = c.ShouldBindJSON(&loginInfo); err != nil {
		log.Error("参数校验失败", err)
		return mygin.LogicError("参数不完善")
	}
	if loginInfo.Name == "" && loginInfo.Email == "" {
		return mygin.LogicError("邮箱或用户名不能为空")
	}

	loginRequest := &menger.LoginRequest{
		Email:    loginInfo.Email,
		Name:     loginInfo.Name,
		Password: loginInfo.Password,
	}

	if loginRes, err = m.mengerClient.Login(context.Background(), loginRequest); err != nil {
		log.Error("调用wm.sx.svc.menger失败", err)
		return mygin.ServerError()
	}
	if loginRes.ErrorMsg != nil {
		return mygin.LogicError(loginRes.ErrorMsg.Msg)
	}

	c.JSON(http.StatusOK, gin.H{
		"token":      loginRes.Token,
		"mengerInfo": loginRes.MengerInfo,
	})
	return nil
}

func (m *MengerHandler) Logout(c *gin.Context) error {

	return nil
}
