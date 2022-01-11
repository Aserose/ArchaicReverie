package handler

import (
	"encoding/json"
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
	"github.com/Aserose/ArchaicReverie/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h Handler) infoAboutSelectedChar(c *gin.Context) {
	_, character, ok := h.getSelectedCharacter(c)
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

	c.JSON(http.StatusOK, character)
}

func (h Handler) beginActionScene(c *gin.Context) {
	locationFeatures := h.service.Action.GenerateScene()
	if locationFeatures != nil {
		c.JSON(http.StatusOK, locationFeatures)
	} else {
		h.actionScene(c)
	}
}

func (h Handler) beginRepast(c *gin.Context) {
	foods := h.service.Action.GetFoodList()

	if len(foods) != 0 {
		c.JSON(http.StatusOK, foods)
	} else {
		h.repast(c)
	}
}

func (h Handler) repast(c *gin.Context) {
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

	var status string
	order := unmarshalOrder(h.readRespBody(c.Request.Body), h.log)

	status, character = h.service.Action.Eat(character, order)
	h.updateCookie(c, h.service.UpdateToken(userId, character))

	if status != empty {
		_, err := c.Writer.WriteString(status)
		if err != nil {
			h.log.Errorf(h.logMsg.Format, h.log.CallInfoStr(), err.Error())
		}
	} else {
		return
	}

}

func unmarshalOrder(respBody []byte, log logger.Logger) model.Food {
	var order model.Food

	if err := json.Unmarshal(respBody, &order); err != nil {
		log.Errorf("%s:%s", log.CallInfoStr(), err.Error)
	}

	return order
}

func (h Handler) actionScene(c *gin.Context) {
	var actionResult string
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

	action := h.unmarshalAction(h.readRespBody(c.Request.Body), h.log)

	if action.InAction == "jump" {
		actionResult, character = h.service.Action.Jump(character, action.Jump)
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
		return userId.(int), character, false
	}

	characterMarshal, _ := c.Get(Character)

	err := json.Unmarshal(characterMarshal.([]byte), &character)
	if err != nil {
		h.log.Errorf(h.logMsg.Format, h.log.CallInfoStr(), err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
	}

	return userId.(int), character, true
}
