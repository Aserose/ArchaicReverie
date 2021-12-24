package handler

import (
	"github.com/Aserose/ArchaicReverie/internal/config"
	"github.com/Aserose/ArchaicReverie/internal/service"
	"github.com/Aserose/ArchaicReverie/pkg/logger"
	"github.com/gin-gonic/gin"
	"io"
)

const HandlerPackageName = "handler"

type Handler struct {
	service   *service.Service
	msgToUser config.MsgToUser
	log       logger.Logger
	logMsg    config.LogMsg
}

func NewHandler(service *service.Service, msgToUser config.MsgToUser, log logger.Logger, logMsg config.LogMsg) *Handler {
	return &Handler{
		service:   service,
		msgToUser: msgToUser,
		log:       log,
		logMsg:    logMsg,
	}
}

func (h *Handler) Routes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/signIn", h.signIn)
		auth.POST("/signUp", h.signUp)
		auth.POST("/signOut", h.signOut)
	}

	api := router.Group("/api", h.verification)
	{
		api.POST("/createChar", h.createChar)
		api.GET("/getAllChar", h.getAllChar)
		api.GET("/getOneChar", h.getChar)
		api.PUT("/updChar", h.updChar)
		api.DELETE("/delChar", h.delChar)

		action := api.Group("/action")
		{
			action.GET("/infoChar", h.infoChar)
		}
	}

	return router
}

func (h Handler) readRespBody(closer io.ReadCloser) []byte {
	respBody, err := io.ReadAll(closer)
	if err != nil {
		h.log.Errorf(h.logMsg.FormatErr, HandlerPackageName, h.logMsg.Read, err.Error())
	}

	return respBody
}
