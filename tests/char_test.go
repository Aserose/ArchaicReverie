package tests

import (
	"github.com/Aserose/ArchaicReverie/internal/app"
	"github.com/Aserose/ArchaicReverie/internal/config"
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
	"github.com/Aserose/ArchaicReverie/pkg/logger"
	cv "github.com/smartystreets/goconvey/convey"
	"net/http"
	"strings"
	"testing"
)

func TestChar(t *testing.T) {
	logs := logger.NewLogger()
	_, msgToUser, charConfig := loadEnv(logs)
	go app.Start(1)
	apiScheme := loadApiScheme(logs)

	cv.Convey("setup", t, func() {
		var (
			client          = http.Client{}
			resp            *http.Response
			numberCharLimit = charConfig.NumberCharLimit
			temp            templates
		)

		cv.Convey("authorize", func() {
			testUser, cookie := temp.authorizeUser(client, apiScheme, logs)

			cv.Convey("createChars", func() {
				var chars []model.Character

				for i := 0; i < numberCharLimit; i++ {
					chars = append(chars, generateChar(testUser.Id))
					resp, _ = temp.doRequest(
						client,
						strings.Split(apiScheme.ApiEndpoints.CreateChar, " ")[0],
						strings.Split(apiScheme.ApiEndpoints.CreateChar, " ")[1],
						reqBody(logs, chars[i]),
						cookie)

					chars[i].CharId = temp.unmarshalChar(temp.readRespBody(resp), logs).CharId
					cv.So(chars[i].CharId, cv.ShouldNotBeZeroValue)
				}

				cv.Convey("errorLimitChar", func() {
					resp, _ = temp.doRequest(
						client,
						strings.Split(apiScheme.ApiEndpoints.CreateChar, " ")[0],
						strings.Split(apiScheme.ApiEndpoints.CreateChar, " ")[1],
						reqBody(logs, model.Character{
							Weight: 10,
							Growth: 180,
						}),
						cookie)

					cv.So(string(temp.readRespBody(resp)), cv.ShouldEqual, msgToUser.CharStatus.CharCreateLimit)

					cv.Convey("getAllChars", func() {
						resp, _ = temp.doRequest(
							client,
							strings.Split(apiScheme.ApiEndpoints.GetAllChar, " ")[0],
							strings.Split(apiScheme.ApiEndpoints.GetAllChar, " ")[1],
							reqBody(logs, testUser),
							cookie)

						receiveChars := temp.unmarshalChars(temp.readRespBody(resp), logs)

						for i := 0; i <= len(chars)-1; i++ {
							cv.So(receiveChars[i], cv.ShouldResemble, chars[i])
						}

						cv.Convey("updChar", func() {
							resp, _ = temp.doRequest(
								client,
								strings.Split(apiScheme.ApiEndpoints.UpdChar, " ")[0],
								strings.Split(apiScheme.ApiEndpoints.UpdChar, " ")[1],
								reqBody(logs, func() model.Character {
									a := generateChar(chars[0].OwnerId)
									a.CharId = chars[0].CharId
									chars[0] = a
									return a
								}()), cookie)

							cv.So(string(temp.readRespBody(resp)), cv.ShouldEqual, msgToUser.CharStatus.CharUpdate)

							cv.Convey("delChar", func() {
								resp, _ = temp.doRequest(
									client,
									strings.Split(apiScheme.ApiEndpoints.DelChar, " ")[0],
									strings.Split(apiScheme.ApiEndpoints.DelChar, " ")[1],
									reqBody(logs, chars[len(chars)-1]),
									cookie)

								resp, _ = temp.doRequest(
									client,
									strings.Split(apiScheme.ApiEndpoints.GetAllChar, " ")[0],
									strings.Split(apiScheme.ApiEndpoints.GetAllChar, " ")[1],
									reqBody(logs, testUser),
									cookie)

								receiveChars = temp.unmarshalChars(temp.readRespBody(resp), logs)
								cv.So(len(chars), cv.ShouldNotEqual, len(receiveChars))

								cv.Convey("createForbiddenChar", func() {
									resp, _ = temp.doRequest(
										client,
										strings.Split(apiScheme.ApiEndpoints.CreateChar, " ")[0],
										strings.Split(apiScheme.ApiEndpoints.CreateChar, " ")[1],
										reqBody(logs, model.Character{
											Weight: 10,
											Growth: 180,
										}),
										cookie)

									cv.So(string(temp.readRespBody(resp)), cv.ShouldEqual, msgToUser.CharStatus.CharWeightRange)

									resp, _ = temp.doRequest(
										client,
										strings.Split(apiScheme.ApiEndpoints.CreateChar, " ")[0],
										strings.Split(apiScheme.ApiEndpoints.CreateChar, " ")[1],
										reqBody(logs, model.Character{
											Weight: 60,
											Growth: 290,
										}),
										cookie)

									cv.So(string(temp.readRespBody(resp)), cv.ShouldEqual, msgToUser.CharStatus.CharGrowthRange)

									resp, _ = temp.doRequest(
										client,
										strings.Split(apiScheme.ApiEndpoints.CreateChar, " ")[0],
										strings.Split(apiScheme.ApiEndpoints.CreateChar, " ")[1],
										reqBody(logs, model.Character{
											Weight: 10,
											Growth: 290,
										}),
										cookie)

									cv.So(string(temp.readRespBody(resp)), cv.ShouldContainSubstring, msgToUser.CharStatus.CharGrowthRange)

									cv.Convey("selectChar", func() {
										resp, cookie = temp.doRequest(
											client,
											strings.Split(apiScheme.ApiEndpoints.SelectChar, " ")[0],
											strings.Split(apiScheme.ApiEndpoints.SelectChar, " ")[1],
											reqBody(logs, chars[0]),
											cookie)

										cv.So(temp.unmarshalInt(temp.readRespBody(resp), logs), cv.ShouldEqual, chars[0].CharId)

										cv.Convey("delete", func() {
											resp, cookie = temp.doRequest(
												client,
												strings.Split(apiScheme.AuthEndpoints.DeleteAccount, " ")[0],
												strings.Split(apiScheme.AuthEndpoints.DeleteAccount, " ")[1],
												reqBody(logs, testUser),
												cookie)

											resp, cookie = temp.doRequest(
												client,
												strings.Split(apiScheme.AuthEndpoints.SignIn, " ")[0],
												strings.Split(apiScheme.AuthEndpoints.SignIn, " ")[1],
												reqBody(logs, testUser),
												cookie)

											cv.So(string(temp.readRespBody(resp)), cv.ShouldEqual, msgToUser.AuthStatus.InvalidUsername)
										})
									})
								})
							})
						})
					})
				})
			})
		})
	})
}

func (a authorization) authorizeUser(client http.Client, apiScheme config.Endpoints, logs logger.Logger) (model.User, []*http.Cookie) {
	var (
		resp     *http.Response
		cookie   []*http.Cookie
		testUser = generateTestUser()
		temp     templates
	)

	resp, cookie = temp.doRequest(
		client,
		strings.Split(apiScheme.AuthEndpoints.SignUp, " ")[0],
		strings.Split(apiScheme.AuthEndpoints.SignUp, " ")[1], reqBody(logs, testUser),
		cookie)

	testUser.Id = temp.unmarshalInt(temp.readRespBody(resp), logs)

	return testUser, cookie
}
