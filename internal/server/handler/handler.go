package handler

import (
	"github.com/Aserose/ArchaicReverie/internal/config"
	"github.com/Aserose/ArchaicReverie/internal/service"
	"github.com/Aserose/ArchaicReverie/pkg/logger"
	"github.com/gin-gonic/gin"
	"io"
)

type Handler struct {
	service      *service.Service
	utilitiesStr config.UtilitiesStr
	msgToUser    config.MsgToUser
	log          logger.Logger
	logMsg       config.LogMsg
}

func NewHandler(service *service.Service, utilitiesStr config.UtilitiesStr, msgToUser config.MsgToUser, log logger.Logger, logMsg config.LogMsg) *Handler {
	return &Handler{
		service:      service,
		utilitiesStr: utilitiesStr,
		msgToUser:    msgToUser,
		log:          log,
		logMsg:       logMsg,
	}
}

func (h *Handler) Routes(endpoints config.Endpoints) *gin.Engine {
	router := gin.New()

	auth := router.Group(endpoints.AuthEndpoints.Auth)
	{
		auth.POST(endpoints.AuthEndpoints.SignIn, h.signIn)
		auth.POST(endpoints.AuthEndpoints.SignUp, h.signUp)
		auth.POST(endpoints.AuthEndpoints.SignOut, h.signOut)
		auth.POST(endpoints.AuthEndpoints.NewPassword, h.updPassword)
		auth.DELETE(endpoints.AuthEndpoints.DeleteAccount, h.deleteAccount)
	}

	api := router.Group(endpoints.ApiEndpoints.Api, h.verification)
	{
		api.POST(endpoints.ApiEndpoints.CreateChar, h.createChar)
		api.GET(endpoints.ApiEndpoints.GetAllChar, h.getAllChar)
		api.POST(endpoints.ApiEndpoints.SelectChar, h.selectChar)
		api.PUT(endpoints.ApiEndpoints.UpdChar, h.updChar)
		api.DELETE(endpoints.ApiEndpoints.DelChar, h.delChar)
		api.DELETE(endpoints.ApiEndpoints.DelAllChar, h.delAllChar)

		action := api.Group(endpoints.ActionEndpoints.Action)
		{
			action.GET(endpoints.ActionEndpoints.CharacterMenu, h.characterMenu)
			action.POST(endpoints.ActionEndpoints.ActionScene, h.beginActionScene)
			action.POST(endpoints.ActionEndpoints.Restock, h.restock)
		}
	}

	return router
}

func (h Handler) readRespBody(closer io.ReadCloser) []byte {
	respBody, err := io.ReadAll(closer)
	if err != nil {
		h.log.Errorf(h.logMsg.Format, h.log.CallInfoStr(), err.Error())
	}

	return respBody
}
