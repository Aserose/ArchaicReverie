package handler

import (
	"encoding/json"
	"fmt"
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
	"github.com/Aserose/ArchaicReverie/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h Handler) createChar(c *gin.Context) {
	userId, ok := c.Get(UserId)
	if !ok {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	creatingChar := unmarshalRespCharacter(h.readRespBody(c.Request.Body), userId.(int), h.log)

	charId, err := h.service.CreateCharacter(creatingChar)
	if err != nil {
		switch err.Error() {
		case h.logMsg.CharGrowthOutErr:
			if _, err := c.Writer.WriteString(h.msgToUser.CharStatus.CharGrowthRange); err != nil {
				h.log.Panicf(h.logMsg.Format, h.log.CallInfoStr(), err.Error())
			}
			return
		case h.logMsg.CharWeightOutErr:
			if _, err := c.Writer.WriteString(h.msgToUser.CharStatus.CharWeightRange); err != nil {
				h.log.Panicf(h.logMsg.Format, h.log.CallInfoStr(), err.Error())
			}
		case h.logMsg.CharGrowthAndWeightOutErr:
			if _, err := c.Writer.WriteString(
				fmt.Sprintf("%s\n%s", h.msgToUser.CharStatus.CharGrowthRange, h.msgToUser.CharStatus.CharWeightRange)); err != nil {
				h.log.Panicf(h.logMsg.Format, h.log.CallInfoStr(), err.Error())
			}
		case h.logMsg.CharLimitOutErr:
			if _, err := c.Writer.WriteString(
				fmt.Sprint(h.msgToUser.CharStatus.CharCreateLimit)); err != nil {
				h.log.Panicf(h.logMsg.Format, h.log.CallInfoStr(), err.Error())
			}
		default:
			c.Writer.WriteHeader(http.StatusInternalServerError)
			h.log.Panicf(h.logMsg.Format, h.log.CallInfoStr(), err.Error())
			return
		}
		return
	}

	creatingChar.CharId = charId
	creatingChar.OwnerId = userId.(int)

	c.JSON(http.StatusOK, creatingChar)

}

func (h Handler) getAllChar(c *gin.Context) {
	userId, ok := c.Get(UserId)
	if !ok {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	c.JSON(http.StatusOK, h.service.GetAllCharacters(userId.(int)))

}

func (h Handler) selectChar(c *gin.Context) {
	userId, ok := c.Get(UserId)
	if !ok {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	selectedCharId := unmarshalRespCharacter(h.readRespBody(c.Request.Body), userId.(int), h.log).CharId

	h.setCookie(c, h.service.Authorization.UpdateToken(
		userId.(int), h.service.Character.SelectChar(
			userId.(int), selectedCharId)), selectedCharId)

}

func (h Handler) updChar(c *gin.Context) {
	userId, ok := c.Get(UserId)
	if !ok {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := h.service.Character.Update(unmarshalRespCharacter(h.readRespBody(c.Request.Body), userId.(int), h.log)); err != nil {
		switch err.Error() {
		case h.logMsg.CharGrowthOutErr:
			if _, err := c.Writer.WriteString(h.msgToUser.CharStatus.CharGrowthRange); err != nil {
				h.log.Panicf(h.logMsg.Format, h.log.CallInfoStr(), err.Error())
			}
			return
		case h.logMsg.CharWeightOutErr:
			if _, err := c.Writer.WriteString(h.msgToUser.CharStatus.CharWeightRange); err != nil {
				h.log.Panicf(h.logMsg.Format, h.log.CallInfoStr(), err.Error())
			}
		default:
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	_, err := c.Writer.WriteString(h.msgToUser.CharStatus.CharUpdate)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h Handler) delChar(c *gin.Context) {
	userId, ok := c.Get(UserId)
	if !ok {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := h.service.Character.Delete(userId.(int), unmarshalRespCharacter(h.readRespBody(c.Request.Body), userId.(int), h.log).CharId); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err := c.Writer.WriteString(h.msgToUser.CharStatus.CharDelete)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h Handler) delAllChar(c *gin.Context) {
	userId, ok := c.Get(UserId)
	if !ok {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err := h.service.Character.DeleteAll(userId.(int)); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err := c.Writer.WriteString(h.msgToUser.CharStatus.CharAllDelete)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func unmarshalRespCharacter(respBody []byte, userId int, log logger.Logger) model.Character {
	var character model.Character

	if err := json.Unmarshal(respBody, &character); err != nil {
		log.Errorf("%s:%s", log.CallInfoStr(), err.Error())
	}
	character.OwnerId = userId

	return character
}
