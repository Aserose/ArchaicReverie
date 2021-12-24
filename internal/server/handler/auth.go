package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	AuthHandlerPackageName = "AuthHandler"
	username               = "username"
	password               = "password"
	UserId                 = "userId"
	Character              = "character"
	cookiePath             = "/api"
	Empty                  = ""
)

func (h Handler) signIn(c *gin.Context) {
	respBody := unmarshalCredentials(h.readRespBody(c.Request.Body))
	h.setToken(c, h.service.Authorization.SignIn(
		respBody[username],
		respBody[password]),
		h.msgToUser.AuthStatus.SignIn)

}
func (h Handler) signUp(c *gin.Context) {
	respBody := unmarshalCredentials(h.readRespBody(c.Request.Body))
	h.setToken(c, h.service.Authorization.SignUp(
		respBody[username],
		respBody[password]),
		h.msgToUser.AuthStatus.SignUp)
}

func (h Handler) signOut(c *gin.Context) {
	h.setToken(c, Empty, h.msgToUser.AuthStatus.SignOut)
}

func unmarshalCredentials(respBody []byte) map[string]string {
	var unmarshalRespBody map[string]string
	json.Unmarshal(respBody, &unmarshalRespBody)

	return unmarshalRespBody
}

func (h Handler) setToken(c *gin.Context, token, successStatus string) {
	switch token {
	case h.msgToUser.AuthStatus.BusyUsername:
		c.Writer.WriteString(token)
	case h.msgToUser.AuthStatus.UserNotFound:
		c.Writer.WriteString(token)
	case h.msgToUser.AuthStatus.InvalidPassword:
		c.Writer.WriteString(token)
	case h.msgToUser.AuthStatus.InvalidUsername:
		c.Writer.WriteString(token)
	default:
		http.SetCookie(c.Writer, &http.Cookie{
			Name:  h.msgToUser.AuthStatus.CookieName,
			Value: token,
			Path:  cookiePath,
		})
		c.Writer.Write([]byte(successStatus))
	}
}

func (h Handler) verification(c *gin.Context) {
	token, err := c.Request.Cookie(h.msgToUser.AuthStatus.CookieName)
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
		h.log.Errorf(h.logMsg.FormatErr, AuthHandlerPackageName)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if userId != 0 {
		c.Set(UserId, userId)
		c.Set(Character, marshalCharacter)
	}
}
