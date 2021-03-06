package test

import (
	"github.com/Aserose/ArchaicReverie/internal/app"
	"github.com/Aserose/ArchaicReverie/pkg/logger"
	cv "github.com/smartystreets/goconvey/convey"
	"net/http"
	"strings"
	"testing"
)

func TestAuth(t *testing.T) {
	logs := logger.NewLogger()
	_, msgToUser, _ := loadEnv(logs)

	go app.Start(1)

	apiScheme := loadApiScheme(logs)

	cv.Convey("setup", t, func() {
		var (
			client   = http.Client{}
			cookie   []*http.Cookie
			resp     *http.Response
			testUser = generateTestUser()
			temp     = NewTemplates(logs)
		)

		cv.Convey("signUp", func() {
			resp, cookie = temp.readAndReq.doRequest(
				client,
				strings.Split(apiScheme.AuthEndpoints.SignUp, " ")[0],
				strings.Split(apiScheme.AuthEndpoints.SignUp, " ")[1], reqBody(logs, testUser),
				cookie)

			testUser.Id = temp.readAndReq.unmarshalInt(temp.readAndReq.readRespBody(resp))
			cv.So(testUser.Id, cv.ShouldNotBeZeroValue)
			cv.So(cookie[0].Value, cv.ShouldNotBeEmpty)

			cv.Convey("newPassword", func() {
				newPassword := randStr(7)

				resp, _ = temp.readAndReq.doRequest(
					client,
					strings.Split(apiScheme.AuthEndpoints.NewPassword, " ")[0],
					strings.Split(apiScheme.AuthEndpoints.NewPassword, " ")[1],
					reqBody(logs, struct {
						Id          int    `json:"id"`
						Username    string `json:"username"`
						Password    string `json:"password"`
						NewPassword string `json:"newPassword"`
					}{
						Id:          testUser.Id,
						Username:    testUser.Username,
						Password:    testUser.Password,
						NewPassword: newPassword,
					}), cookie)

				testUser.Password = newPassword

				cv.So(string(temp.readAndReq.readRespBody(resp)), cv.ShouldEqual, msgToUser.AuthStatus.PasswordUpdated)

				cv.Convey("signOut", func() {
					resp, cookie = temp.readAndReq.doRequest(
						client,
						strings.Split(apiScheme.AuthEndpoints.SignOut, " ")[0],
						strings.Split(apiScheme.AuthEndpoints.SignOut, " ")[1],
						reqBody(logs, testUser),
						cookie)

					cv.Convey("signIn", func() {
						resp, cookie = temp.readAndReq.doRequest(
							client,
							strings.Split(apiScheme.AuthEndpoints.SignIn, " ")[0],
							strings.Split(apiScheme.AuthEndpoints.SignIn, " ")[1],
							reqBody(logs, testUser),
							cookie)

						testUser.Id = temp.readAndReq.unmarshalInt(temp.readAndReq.readRespBody(resp))
						cv.So(testUser.Id, cv.ShouldNotBeZeroValue)

						cv.Convey("delete", func() {
							resp, cookie = temp.readAndReq.doRequest(
								client,
								strings.Split(apiScheme.AuthEndpoints.DeleteAccount, " ")[0],
								strings.Split(apiScheme.AuthEndpoints.DeleteAccount, " ")[1],
								reqBody(logs, testUser),
								cookie)

							resp, cookie = temp.readAndReq.doRequest(
								client,
								strings.Split(apiScheme.AuthEndpoints.SignIn, " ")[0],
								strings.Split(apiScheme.AuthEndpoints.SignIn, " ")[1],
								reqBody(logs, testUser),
								cookie)

							cv.So(string(temp.readAndReq.readRespBody(resp)), cv.ShouldEqual, msgToUser.AuthStatus.InvalidUsername)
						})
					})
				})
			})
		})
	})
}
