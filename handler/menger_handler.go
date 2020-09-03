package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/util/log"
	mengerProto "github.com/wmsx/menger_svc/proto/menger"
	mygin "github.com/wmsx/pkg/gin"
	"net/http"
	"strconv"
)

type LoginInfo struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterInfo struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
	Avatar   string `json:"avatar"`
}

type MengerHandler struct {
	mengerClient mengerProto.MengerService
}

func NewMengerHandler(client client.Client) *MengerHandler {
	return &MengerHandler{
		mengerClient: mengerProto.NewMengerService(mengerSvcName, client),
	}
}

func (m *MengerHandler) Register(c *gin.Context) {
	var (
		registerRes  *mengerProto.RegisterResponse
		loginRes     *mengerProto.LoginResponse
		err          error
		registerInfo RegisterInfo
	)

	app := mygin.Gin{C: c}

	if err = c.ShouldBindJSON(&registerInfo); err != nil {
		log.Error("参数校验失败", err)
		app.LogicErrorResponse("参数不完善")
		return
	}

	registerRequest := &mengerProto.RegisterRequest{
		Name:     registerInfo.Name,
		Password: registerInfo.Password,
		Avatar:   registerInfo.Avatar,
	}

	cxt := context.Background()
	if registerRes, err = m.mengerClient.Register(cxt, registerRequest); err != nil {
		log.Error("注册失败", err)
		app.ServerErrorResponse()
		return
	}
	if registerRes.ErrorMsg != nil {
		c.String(http.StatusBadRequest, registerRes.ErrorMsg.Msg)
		return
	}

	log.Info("注册成功", registerRes)

	loginReq := &mengerProto.LoginRequest{
		Name:     registerInfo.Name,
		Password: registerInfo.Password,
	}

	loginRes, err = m.mengerClient.Login(cxt, loginReq)
	if err != nil {
		log.Error("自动登陆失败 err:", err)
		app.LogicErrorResponse("自动登陆失败")
		return
	}
	if loginRes.ErrorMsg != nil {
		log.Error("自动登陆失败 errorMsg: ", loginRes.ErrorMsg)
		app.LogicErrorResponse("自动登陆失败")
		return
	}

	app.Response(loginRes.MengerInfo)
	return
}

func (m *MengerHandler) Login(c *gin.Context) {
	var (
		loginRes  *mengerProto.LoginResponse
		err       error
		loginInfo *LoginInfo
		s         *mygin.Session
	)
	app := mygin.Gin{C: c}

	if err = c.ShouldBindJSON(&loginInfo); err != nil {
		log.Error("参数校验失败", err)
		app.LogicErrorResponse("参数不完善")
		return
	}
	if loginInfo.Username == "" {
		app.LogicErrorResponse("用户名不能为空")
		return
	}

	loginRequest := &mengerProto.LoginRequest{
		Name:     loginInfo.Username,
		Password: loginInfo.Password,
	}

	if loginRes, err = m.mengerClient.Login(context.Background(), loginRequest); err != nil {
		log.Error("登录失败", err)
		app.ServerErrorResponse()
		return
	}
	if loginRes.ErrorMsg != nil {
		app.LogicErrorResponse(loginRes.ErrorMsg.Msg)
		return
	}

	if s, err = mygin.NewSession(c); err != nil {
		log.Error("获取session失败 err: ", err)
		app.ServerErrorResponse()
		return
	}
	s.SaveMenger(loginRes.MengerInfo.Id, loginRes.MengerInfo.Name)
	if err = s.Save(); err != nil {
		log.Error("保存session失败", err)
		app.ServerErrorResponse()
		return
	}
	app.Response(loginRes.MengerInfo)
	return
}

func (m *MengerHandler) Logout(c *gin.Context) {
	var (
		err error
		s   *mygin.Session
	)
	app := mygin.Gin{C: c}
	if s, err = mygin.NewSession(c); err != nil {
		log.Error("获取session失败 err: ", err)
		app.ServerErrorResponse()
		return
	}
	if err = s.Remove(); err != nil {
		log.Error("保存session失败", err)
		app.ServerErrorResponse()
		return
	}
	app.Response(nil)
	return
}

func (h *MengerHandler) GetMengerInfo(c *gin.Context) {
	var (
		mengerId     int64
		err          error
		getMengerRes *mengerProto.GetMengerResponse
	)
	app := mygin.Gin{C: c}

	if mengerId, err = strconv.ParseInt(c.Param("id"), 10, 64); err != nil {
		app.ServerErrorResponse()
		return
	}

	getMengerRequest := &mengerProto.GetMengerRequest{MengerId: mengerId}
	if getMengerRes, err = h.mengerClient.GetMenger(c, getMengerRequest); err != nil {
		log.Error("GetMenger 获取用户信息失败 err: ", err)
		app.ServerErrorResponse()
		return
	}
	app.Response(getMengerRes.MengerInfo)
	return
}
