package handler

import (
	"encoding/json"
	"fmt"
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (h Handler) createChar(c *gin.Context) {
	userId, ok := c.Get(UserId)
	if !ok {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := h.service.CreateCharacter(unmarshalCharacter(h.readRespBody(c.Request.Body), userId.(int))); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err := c.Writer.WriteString(h.msgToUser.CharStatus.CharCreate)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (h Handler) getAllChar(c *gin.Context) {
	userId, ok := c.Get(UserId)
	if !ok {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	for _, character := range h.service.GetAllCharacters(userId.(int)) {
		//TODO
		_, err := c.Writer.WriteString(fmt.Sprint(character))
		//TODO
		if err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		} //форматировать/упорядочивать ответ
	}
}

func (h Handler) getChar(c *gin.Context) {
	userId, ok := c.Get(UserId)
	if !ok {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	h.setToken(c, h.service.Authorization.UpdateToken(
		userId.(int), h.service.Character.GetOne(
			userId.(int), unmarshalCharacter(h.readRespBody(c.Request.Body), userId.(int)).CharId)), "ok")

	//log selection character
}

func (h Handler) updChar(c *gin.Context) {
	userId, ok := c.Get(UserId)
	if !ok {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := h.service.Character.Update(unmarshalCharacter(h.readRespBody(c.Request.Body), userId.(int))); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
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

	if err := h.service.Character.Delete(unmarshalCharacter(h.readRespBody(c.Request.Body), userId.(int)).CharId); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err := c.Writer.WriteString(h.msgToUser.CharStatus.CharDelete)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func unmarshalCharacter(respBody []byte, userId int) model.Character {
	var character model.Character

	json.Unmarshal(respBody, &character)
	character.OwnerId = userId
	log.Print(character)

	return character
}
