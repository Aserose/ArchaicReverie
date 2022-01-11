package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const (
	username    = "username"
	password    = "password"
	newPassword = "newPassword"
	UserId      = "userId"
	Character   = "character"
	cookiePath  = "/"
	empty       = ""
)

func (h Handler) signIn(c *gin.Context) {
	obtainedCookie, _ := c.Request.Cookie(h.utilitiesStr.CookieName)

	if obtainedCookie != nil {
		if obtainedCookie.Value != empty {
			if _, err := c.Writer.WriteString(h.msgToUser.AuthStatus.SignAlready); err != nil {
				h.log.Panicf(h.logMsg.Format, h.log.CallInfoStr(), err.Error())
			}
			return
		}
	}

	respBody := unmarshalCredentials(h.readRespBody(c.Request.Body))
	token, id := h.service.Authorization.SignIn(
		respBody[username],
		respBody[password])

	h.setCookie(c, token, id)
}

func (h Handler) signUp(c *gin.Context) {
	obtainedCookie, _ := c.Request.Cookie(h.utilitiesStr.CookieName)

	if obtainedCookie != nil {
		if obtainedCookie.Value != empty {
			if _, err := c.Writer.WriteString(h.msgToUser.AuthStatus.SignAlready); err != nil {
				h.log.Panicf(h.logMsg.Format, h.log.CallInfoStr(), err.Error())
			}
			return
		}
	}

	respBody := unmarshalCredentials(h.readRespBody(c.Request.Body))
	token, id := h.service.Authorization.SignUp(
		respBody[username],
		respBody[password])

	h.setCookie(c, token, id)
}

func (h Handler) signOut(c *gin.Context) {
	h.setCookie(c, empty, 0)
}

func (h Handler) updPassword(c *gin.Context) {
	respBody := unmarshalCredentials(h.readRespBody(c.Request.Body))
	if _, err := c.Writer.WriteString(h.service.Authorization.UpdatePassword(
		respBody[username],
		respBody[password],
		respBody[newPassword])); err != nil {
		h.log.Panicf(h.logMsg.Format, h.log.CallInfoStr(), err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
	}
}

func (h Handler) deleteAccount(c *gin.Context) {
	respBody := unmarshalCredentials(h.readRespBody(c.Request.Body))
	h.setCookie(c, empty, 0)
	if _, err := c.Writer.WriteString(h.service.Authorization.DeleteAccount(
		respBody[username],
		respBody[password])); err != nil {
		h.log.Panicf(h.logMsg.Format, h.log.CallInfoStr(), err.Error())
	}
}

func unmarshalCredentials(respBody []byte) map[string]string {
	var unmarshalRespBody map[string]string
	json.Unmarshal(respBody, &unmarshalRespBody)

	return unmarshalRespBody
}

func (h Handler) setCookie(c *gin.Context, token string, id int) {
	switch token {
	case h.msgToUser.AuthStatus.BusyUsername:
		if _, err := c.Writer.WriteString(token); err != nil {
			h.log.Panicf(h.logMsg.Format, h.log.CallInfoStr(), err.Error())
		}
	case h.msgToUser.AuthStatus.UserNotFound:
		if _, err := c.Writer.WriteString(token); err != nil {
			h.log.Panicf(h.logMsg.Format, h.log.CallInfoStr(), err.Error())
		}
	case h.msgToUser.AuthStatus.InvalidPassword:
		if _, err := c.Writer.WriteString(token); err != nil {
			h.log.Panicf(h.logMsg.Format, h.log.CallInfoStr(), err.Error())
		}
	case h.msgToUser.AuthStatus.InvalidUsername:
		if _, err := c.Writer.WriteString(token); err != nil {
			h.log.Panicf(h.logMsg.Format, h.log.CallInfoStr(), err.Error())
		}
	case h.msgToUser.CharStatus.CharNotSelect:
		if _, err := c.Writer.WriteString(token); err != nil {
			h.log.Panicf(h.logMsg.Format, h.log.CallInfoStr(), err.Error())
		}
	default:
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     h.utilitiesStr.CookieName,
			Value:    token,
			Path:     cookiePath,
			Expires:  time.Now().AddDate(0, 0, 1),
			HttpOnly: true,
		})
		c.JSON(http.StatusOK, id)
	}
}

func (h Handler) updateCookie(c *gin.Context, token string) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     h.utilitiesStr.CookieName,
		Value:    token,
		Path:     cookiePath,
		Expires:  time.Now().AddDate(0, 0, 1),
		HttpOnly: true,
	})
}

func (h Handler) verification(c *gin.Context) {
	token, err := c.Request.Cookie(h.utilitiesStr.CookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			c.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	userId, character, err := h.service.Verification(token.Value)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	marshalCharacter, err := json.Marshal(character)
	if err != nil {
		h.log.Errorf(h.logMsg.Format, h.log.CallInfoStr(), err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if userId != 0 {
		c.Set(UserId, userId)
		c.Set(Character, marshalCharacter)
	}
}
