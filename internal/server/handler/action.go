package handler

import (
	"encoding/json"
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
	"github.com/Aserose/ArchaicReverie/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h Handler) characterMenu(c *gin.Context) {
	userId, character, ok := h.getSelectedCharacter(c)
	if !ok {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	if character.CharId == 0 {
		if _, err := c.Writer.WriteString(h.msgToUser.CharStatus.CharNotSelect); err != nil {
			h.log.Panicf(h.logMsg.Format, h.log.CallInfoStr(), err.Error())
		}
		return
	}

	items := h.unmarshalItems(h.readRespBody(c.Request.Body), h.log)
	updChar := h.service.Action.CharMenu(character, items)

	if &updChar != &character {
		h.updateCookie(c, h.service.UpdateToken(userId, updChar))
	}

	c.JSON(http.StatusOK, updChar)
}

func (h Handler) beginActionScene(c *gin.Context) {
	userId, character, ok := h.getSelectedCharacter(c)
	if !ok {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	if character.CharId == 0 {
		if _, err := c.Writer.WriteString(h.msgToUser.CharStatus.CharNotSelect); err != nil {
			h.log.Panicf(h.logMsg.Format, h.log.CallInfoStr(), err.Error())
		}
		return
	}

	if !h.service.Action.HealthCheck(character) {
		_, err := c.Writer.WriteString(h.msgToUser.ActionMsg.LowHP)
		if err != nil {
			h.log.Errorf(h.logMsg.Format, h.log.CallInfoStr(), err.Error())
		}
		return
	}

	locationFeatures := h.service.Action.GenerateScene()

	if locationFeatures != nil {
		c.JSON(http.StatusOK, locationFeatures)
		return
	} else {
		h.actionScene(c, character, userId)
	}
}

func (h Handler) restock(c *gin.Context) {
	userId, character, ok := h.getSelectedCharacter(c)
	if !ok {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	if character.CharId == 0 {
		if _, err := c.Writer.WriteString(h.msgToUser.CharStatus.CharNotSelect); err != nil {
			h.log.Panicf(h.logMsg.Format, h.log.CallInfoStr(), err.Error())
		}
		return
	}

	items := h.unmarshalItems(h.readRespBody(c.Request.Body), h.log)

	availableItems, updatedChar, status := h.service.Action.RecreationalMain(character, items)

	if availableItems.Weapons != nil && availableItems.Foods != nil {
		c.JSON(http.StatusOK, availableItems)
		return
	}

	if status != empty {
		_, err := c.Writer.WriteString(status)
		if err != nil {
			h.log.Errorf(h.logMsg.Format, h.log.CallInfoStr(), err.Error())
		}
		return
	}

	if &updatedChar != &character {
		h.updateCookie(c, h.service.UpdateToken(userId, updatedChar))
		return
	}
}

func (h Handler) unmarshalItems(respBody []byte, log logger.Logger) model.Items {
	var items model.Items

	if err := json.Unmarshal(respBody, &items); err != nil {
		h.log.Errorf("%s:%s", log.CallInfoStr(), err.Error)
	}

	return items
}

func (h Handler) actionScene(c *gin.Context, character model.Character, userId int) {
	var actionResult string

	action := h.unmarshalAction(h.readRespBody(c.Request.Body), h.log)

	if action.InAction != empty {
		actionResult, character = h.service.Action.Action(character, action)
		h.updateCookie(c, h.service.UpdateToken(userId, character))

		switch actionResult {
		case h.utilitiesStr.BadRequest:
			c.Writer.WriteHeader(http.StatusBadRequest)
			return
		default:
			if _, err := c.Writer.WriteString(actionResult); err != nil {
				h.log.Panicf(h.logMsg.Format, h.log.CallInfoStr(), err.Error())
			}
			return
		}

	}
}

func (h Handler) unmarshalAction(respBody []byte, log logger.Logger) model.Action {
	var action model.Action

	if err := json.Unmarshal(respBody, &action); err != nil {
		log.Errorf("%s:%s", log.CallInfoStr(), err.Error())
	}

	return action
}

func (h Handler) getSelectedCharacter(c *gin.Context) (int, model.Character, bool) {
	var character model.Character

	userId, ok := c.Get(UserId)
	if !ok {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return 0, model.Character{}, false
	}

	characterMarshal, _ := c.Get(Character)

	err := json.Unmarshal(characterMarshal.([]byte), &character)
	if err != nil {
		h.log.Errorf(h.logMsg.Format, h.log.CallInfoStr(), err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
	}

	return userId.(int), character, true
}
