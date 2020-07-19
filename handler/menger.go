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

func (m *MengerHandler) Register(c *gin.Context) {
	var (
		registerRes  *menger.RegisterResponse
		loginRes     *menger.LoginResponse
		err          error
		registerInfo RegisterInfo
	)

	if err = c.ShouldBindJSON(&registerInfo); err != nil {
		log.Error("参数校验失败", err)
		c.String(http.StatusBadRequest, "参数不完善")
		return
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
		c.String(http.StatusInternalServerError, "服务器异常")
		return
	}
	if registerRes.ErrorMsg != nil {
		c.String(http.StatusBadRequest, registerRes.ErrorMsg.Msg)
		return
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
		c.String(http.StatusUnauthorized, "自动登陆失败")
		return
	}
	if loginRes.ErrorMsg != nil {
		log.Error("自动登陆失败 errorMsg: ", loginRes.ErrorMsg)
		c.String(http.StatusUnauthorized, "自动登陆失败")
		return
	}

	c.JSON(http.StatusOK, loginRes.MengerInfo)
	return
}

func (m *MengerHandler) Login(c *gin.Context) {
	var (
		loginRes  *menger.LoginResponse
		err       error
		loginInfo *LoginInfo
		s         *mygin.Session
	)

	if err = c.ShouldBindJSON(&loginInfo); err != nil {
		log.Error("参数校验失败", err)
		c.String(http.StatusBadRequest, "参数不完善")
		return
	}
	if loginInfo.Name == "" && loginInfo.Email == "" {
		c.String(http.StatusBadRequest, "邮箱或用户名不能为空")
		return
	}

	loginRequest := &menger.LoginRequest{
		Email:    loginInfo.Email,
		Name:     loginInfo.Name,
		Password: loginInfo.Password,
	}

	if loginRes, err = m.mengerClient.Login(context.Background(), loginRequest); err != nil {
		log.Error("登录失败", err)
		c.String(http.StatusInternalServerError, "服务器异常")
		return
	}
	if loginRes.ErrorMsg != nil {
		c.String(http.StatusBadRequest, loginRes.ErrorMsg.Msg)
		return
	}

	if s, err = mygin.NewSession(c); err != nil {
		log.Error("获取session失败 err: ", err)
		c.String(http.StatusInternalServerError, "服务器异常")
		return
	}
	s.SaveMenger(loginRes.MengerInfo.MengerId, loginRes.MengerInfo.Name)
	if err = s.Save(); err != nil {
		log.Error("保存session失败", err)
		c.String(http.StatusInternalServerError, "服务器异常")
		return
	}
	c.JSON(http.StatusOK, loginRes.MengerInfo)
	return
}

func (m *MengerHandler) Logout(c *gin.Context) {
	var (
		err error
		s   *mygin.Session
	)
	if s, err = mygin.NewSession(c); err != nil {
		log.Error("获取session失败 err: ", err)
		c.String(http.StatusInternalServerError, "服务器异常")
		return
	}
	if err = s.Remove(); err != nil {
		log.Error("保存session失败", err)
		c.String(http.StatusInternalServerError, "服务器异常")
		return
	}
	return
}
