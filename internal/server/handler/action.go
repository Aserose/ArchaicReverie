package handler

import (
	"encoding/json"
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h Handler) infoAboutSelectedChar(c *gin.Context) {
	character, ok := h.getSelectedCharacter(c)
	if !ok {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	if character.CharId == 0 {
		if _, err := c.Writer.WriteString(h.msgToUser.CharStatus.CharNotSelect); err != nil {
			h.log.Panicf(h.logMsg.FormatErr, h.log.CallInfoStr(), h.logMsg.WriterResponse, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, character)

	//formatted response
	//if _, err := c.Writer.WriteString(fmt.Sprintf(h.msgToUser.ActionMsg.InfoCharFormat,
	//	character.Name, character.Growth, character.Weight)); err != nil {
	//	c.Writer.WriteHeader(http.StatusInternalServerError)
	//}
}

func (h Handler) beginActionScene(c *gin.Context) {

	locationFeatures := h.service.Action.GenerateScene()

	if locationFeatures != empty {
		if _, err := c.Writer.WriteString(locationFeatures); err != nil {
			h.log.Panicf(h.logMsg.FormatErr, h.log.CallInfoStr(), h.logMsg.WriterResponse, err.Error())
		}
	} else {
		h.actionScene(c)
	}
}

func (h Handler) actionScene(c *gin.Context) {
	character, ok := h.getSelectedCharacter(c)
	if !ok {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	if character.CharId == 0 {
		if _, err := c.Writer.WriteString(h.msgToUser.CharStatus.CharNotSelect); err != nil {
			h.log.Panicf(h.logMsg.FormatErr, h.log.CallInfoStr(), h.logMsg.WriterResponse, err.Error())
		}
		return
	}

	action := h.unmarshalAction(h.readRespBody(c.Request.Body))

	actionResult := h.service.Action.Jump(character, action.Jump)

	if action.InAction == "jump" {
		switch actionResult {
		case h.utilitiesStr.BadRequest:
			c.Writer.WriteHeader(http.StatusBadRequest)
		default:
			if _, err := c.Writer.WriteString(actionResult); err != nil {
				h.log.Panicf(h.logMsg.FormatErr, h.log.CallInfoStr(), h.logMsg.WriterResponse, err.Error())
			}
			return
		}
	}
}

func (h Handler) unmarshalAction(respBody []byte) model.Action {
	var Action model.Action

	json.Unmarshal(respBody, &Action)

	return Action
}

func (h Handler) getSelectedCharacter(c *gin.Context) (model.Character, bool) {
	var character model.Character

	_, ok := c.Get(UserId)
	if !ok {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return character, false
	}

	characterMarshal, _ := c.Get(Character)

	err := json.Unmarshal(characterMarshal.([]byte), &character)
	if err != nil {
		h.log.Errorf(h.logMsg.FormatErr, h.log.CallInfoStr(), h.logMsg.Unmarshal, err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
	}

	return character, true
}
