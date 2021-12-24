package handler

import (
	"encoding/json"
	"fmt"
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	ActionHandlerPackageName = "ActionHandler"
)

func (h Handler) infoChar(c *gin.Context) {
	character, ok := h.getSelectedCharacter(c)
	if !ok {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	if character.CharId == 0 {
		c.Writer.WriteString("no character selected")
		return
	}

	c.Writer.WriteString(fmt.Sprintf("character growth is %d", character.Growth))
}

func (h Handler) actionScene(c *gin.Context) {
	character, ok := h.getSelectedCharacter(c)
	if !ok {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	if character.CharId == 0 {
		c.Writer.WriteString("no character selected")
		return
	}

	action := h.unmarshalAction(h.readRespBody(c.Request.Body))

	if action.InAction == "jump" {
		//TODO
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
		h.log.Errorf(h.logMsg.FormatErr, ActionHandlerPackageName, h.logMsg.Unmarshal, err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
	}

	return character, true
}
