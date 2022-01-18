package tests

import (
	"github.com/Aserose/ArchaicReverie/internal/config"
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
	"github.com/Aserose/ArchaicReverie/pkg/logger"
	"net/http"
	"strings"
)

type authorization struct {
	log logger.Logger
}

func (a authorization) authorizeUserWithCharsAndSelect(client http.Client, apiScheme config.Endpoints, numberCharLimit int) []*http.Cookie {
	cookie, chars := a.authorizeUserWithChars(client, apiScheme, numberCharLimit)

	var temp = NewTemplates(a.log)

	_, cookie = temp.doRequest(
		client,
		strings.Split(apiScheme.ApiEndpoints.SelectChar, " ")[0],
		strings.Split(apiScheme.ApiEndpoints.SelectChar, " ")[1],
		reqBody(a.log, chars[0]),
		cookie)

	return cookie
}

func (a authorization) authorizeUserWithChars(client http.Client, apiScheme config.Endpoints, numberCharLimit int) ([]*http.Cookie, []model.Character) {
	testUser, cookie := a.authorizeUser(client, apiScheme)

	var (
		chars []model.Character
		resp  *http.Response
		temp  = NewTemplates(a.log)
	)

	for i := 0; i < numberCharLimit; i++ {
		chars = append(chars, generateChar(testUser.Id))
		resp, _ = temp.doRequest(
			client,
			strings.Split(apiScheme.ApiEndpoints.CreateChar, " ")[0],
			strings.Split(apiScheme.ApiEndpoints.CreateChar, " ")[1],
			reqBody(a.log, chars[i]),
			cookie)

		chars[i].CharId = temp.unmarshalChar(temp.readRespBody(resp)).CharId
	}
	return cookie, chars
}

func (a authorization) authorizeUser(client http.Client, apiScheme config.Endpoints) (model.User, []*http.Cookie) {
	var (
		resp     *http.Response
		cookie   []*http.Cookie
		testUser = generateTestUser()
		temp     = NewTemplates(a.log)
	)

	resp, cookie = temp.doRequest(
		client,
		strings.Split(apiScheme.AuthEndpoints.SignUp, " ")[0],
		strings.Split(apiScheme.AuthEndpoints.SignUp, " ")[1], reqBody(a.log, testUser),
		cookie)

	testUser.Id = temp.unmarshalInt(temp.readRespBody(resp))

	return testUser, cookie
}
