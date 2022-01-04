package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Aserose/ArchaicReverie/internal/app"
	"github.com/Aserose/ArchaicReverie/internal/config"
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
	"github.com/Aserose/ArchaicReverie/pkg/logger"
	"github.com/joho/godotenv"
	cv "github.com/smartystreets/goconvey/convey"
	"io"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestApp(t *testing.T) {
	logs := logger.NewLogger()
	_, msgToUser := loadEnv(logs)
	go app.Start(1)
	apiScheme := loadApiScheme(logs)

	cv.Convey("setup", t, func() {
		var (
			client = http.Client{}
			cookie []*http.Cookie
			resp   *http.Response
			numberCharLimit = 3
			testUser = generateTestUser()
		)

		cv.Convey("signUp", func() {
			resp, cookie = doRequest(
				client,
				strings.Split(apiScheme.AuthEndpoints.SignUp, " ")[0],
				strings.Split(apiScheme.AuthEndpoints.SignUp, " ")[1], reqBody(logs, testUser),
				cookie)

			testUser.Id = unmarshalInt(readRespBody(resp), logs)

			cv.So(testUser.Id, cv.ShouldNotBeZeroValue)
			cv.So(cookie[0].Value, cv.ShouldNotBeEmpty)

			cv.Convey("newPassword", func() {
				newPassword := RandStr(7)

				resp, _ = doRequest(
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

				cv.So(string(readRespBody(resp)), cv.ShouldEqual, msgToUser.AuthStatus.PasswordUpdated)

				cv.Convey("signOut", func() {
					resp, cookie = doRequest(
						client,
						strings.Split(apiScheme.AuthEndpoints.SignOut, " ")[0],
						strings.Split(apiScheme.AuthEndpoints.SignOut, " ")[1],
						reqBody(logs, testUser),
						cookie)

					cv.Convey("signIn", func() {
						resp, cookie = doRequest(
							client,
							strings.Split(apiScheme.AuthEndpoints.SignIn, " ")[0],
							strings.Split(apiScheme.AuthEndpoints.SignIn, " ")[1],
							reqBody(logs, testUser),
							cookie)

						testUser.Id = unmarshalInt(readRespBody(resp), logs)
						cv.So(testUser.Id, cv.ShouldNotBeZeroValue)

						cv.Convey("createChars", func() {
							var chars []model.Character

							for i := 0; i < numberCharLimit; i++ {
								chars = append(chars, generateChar(testUser.Id))
								resp, _ = doRequest(
									client,
									strings.Split(apiScheme.ApiEndpoints.CreateChar, " ")[0],
									strings.Split(apiScheme.ApiEndpoints.CreateChar, " ")[1],
									reqBody(logs, chars[i]),
									cookie)

								chars[i].CharId = unmarshalChar(readRespBody(resp), logs).CharId

								cv.So(chars[i].CharId, cv.ShouldNotBeZeroValue)
							}

							cv.Convey("getAllChars", func() {
								resp, _ = doRequest(
									client,
									strings.Split(apiScheme.ApiEndpoints.GetAllChar, " ")[0],
									strings.Split(apiScheme.ApiEndpoints.GetAllChar, " ")[1],
									reqBody(logs, testUser),
									cookie)

								receiveChars := unmarshalChars(readRespBody(resp), logs)

								for i := 0; i <= len(chars)-1; i++ {
									cv.So(receiveChars[i], cv.ShouldResemble, chars[i])
								}

								cv.Convey("updChar", func() {
									resp, _ = doRequest(
										client,
										strings.Split(apiScheme.ApiEndpoints.UpdChar, " ")[0],
										strings.Split(apiScheme.ApiEndpoints.UpdChar, " ")[1],
										reqBody(logs, func() model.Character {
											a := generateChar(chars[0].OwnerId)
											a.CharId = chars[0].CharId
											chars[0] = a
											return a
										}()),
										cookie)

									cv.So(string(readRespBody(resp)), cv.ShouldEqual, msgToUser.CharStatus.CharUpdate)

									cv.Convey("delChar", func() {
										resp, _ = doRequest(
											client,
											strings.Split(apiScheme.ApiEndpoints.DelChar, " ")[0],
											strings.Split(apiScheme.ApiEndpoints.DelChar, " ")[1],
											reqBody(logs, chars[len(chars)-1]),
											cookie)

										resp, _ = doRequest(
											client,
											strings.Split(apiScheme.ApiEndpoints.GetAllChar, " ")[0],
											strings.Split(apiScheme.ApiEndpoints.GetAllChar, " ")[1],
											reqBody(logs, testUser),
											cookie)

										receiveChars = unmarshalChars(readRespBody(resp), logs)
										cv.So(len(chars), cv.ShouldNotEqual, len(receiveChars))

										cv.Convey("selectChar", func() {
											resp, cookie = doRequest(
												client,
												strings.Split(apiScheme.ApiEndpoints.SelectChar, " ")[0],
												strings.Split(apiScheme.ApiEndpoints.SelectChar, " ")[1],
												reqBody(logs, chars[0]),
												cookie)

											cv.So(unmarshalInt(readRespBody(resp), logs), cv.ShouldEqual, chars[0].CharId)

											cv.Convey("infoAboutSelectedChar", func() {
												resp, _ = doRequest(
													client,
													strings.Split(apiScheme.ActionEndpoints.InfoAboutSelectedChar, " ")[0],
													strings.Split(apiScheme.ActionEndpoints.InfoAboutSelectedChar, " ")[1],
													reqBody(logs, testUser),
													cookie)

												cv.So(unmarshalChar(readRespBody(resp), logs).CharId, cv.ShouldEqual, chars[0].CharId)

												cv.Convey("BeginActionScene", func() {
													resp, _ = doRequest(
														client,
														strings.Split(apiScheme.ActionEndpoints.BeginActionScene, " ")[0],
														strings.Split(apiScheme.ActionEndpoints.BeginActionScene, " ")[1],
														reqBody(logs, testUser),
														cookie)

													cv.So(string(readRespBody(resp)), cv.ShouldNotBeEmpty)

													var result string
													for {
														resp, _ = doRequest(
															client,
															strings.Split(apiScheme.ActionEndpoints.BeginActionScene, " ")[0],
															strings.Split(apiScheme.ActionEndpoints.BeginActionScene, " ")[1],
															reqBody(logs, generateAction()),
															cookie)
														result = string(readRespBody(resp))
														if result == msgToUser.ActionMsg.JumpOver {
															cv.So(result, cv.ShouldEqual, msgToUser.ActionMsg.JumpOver)
															break
														}
														cv.So(result, cv.ShouldNotBeEmpty)
													}

													cv.Convey("delete", func() {
														resp, cookie = doRequest(
															client,
															strings.Split(apiScheme.AuthEndpoints.DeleteAccount, " ")[0],
															strings.Split(apiScheme.AuthEndpoints.DeleteAccount, " ")[1],
															reqBody(logs, testUser),
															cookie)

														resp, cookie = doRequest(
															client,
															strings.Split(apiScheme.AuthEndpoints.SignIn, " ")[0],
															strings.Split(apiScheme.AuthEndpoints.SignIn, " ")[1],
															reqBody(logs, testUser),
															cookie)

														cv.So(string(readRespBody(resp)), cv.ShouldEqual, msgToUser.AuthStatus.InvalidUsername)
													})
												})
											})
										})
										cv.Convey("createForbiddenChar", func() {
											resp, _ = doRequest(
												client,
												strings.Split(apiScheme.ApiEndpoints.CreateChar, " ")[0],
												strings.Split(apiScheme.ApiEndpoints.CreateChar, " ")[1],
												reqBody(logs, model.Character{
													Weight: 10,
													Growth: 180,
												}),
												cookie)

											cv.So(string(readRespBody(resp)), cv.ShouldEqual, msgToUser.CharStatus.CharWeightRange)

											resp, _ = doRequest(
												client,
												strings.Split(apiScheme.ApiEndpoints.CreateChar, " ")[0],
												strings.Split(apiScheme.ApiEndpoints.CreateChar, " ")[1],
												reqBody(logs, model.Character{
													Weight: 60,
													Growth: 290,
												}),
												cookie)

											cv.So(string(readRespBody(resp)), cv.ShouldEqual, msgToUser.CharStatus.CharGrowthRange)

											resp, _ = doRequest(
												client,
												strings.Split(apiScheme.ApiEndpoints.CreateChar, " ")[0],
												strings.Split(apiScheme.ApiEndpoints.CreateChar, " ")[1],
												reqBody(logs, model.Character{
													Weight: 10,
													Growth: 290,
												}),
												cookie)

											cv.So(string(readRespBody(resp)), cv.ShouldContainSubstring, msgToUser.CharStatus.CharGrowthRange)
										})
									})
								})
							})
							cv.Convey("errorLimitChar", func() {
								resp, _ = doRequest(
									client,
									strings.Split(apiScheme.ApiEndpoints.CreateChar, " ")[0],
									strings.Split(apiScheme.ApiEndpoints.CreateChar, " ")[1],
									reqBody(logs, model.Character{
										Weight: 10,
										Growth: 180,
									}),
									cookie)

								cv.So(string(readRespBody(resp)), cv.ShouldEqual, msgToUser.CharStatus.CharCreateLimit)
							})
						})
					})
				})
			})
		})
	})
}

func doRequest(client http.Client, method string, url string, body io.Reader, cookie []*http.Cookie) (*http.Response, []*http.Cookie) {
	request, _ := http.NewRequest(method, url, body)
	for i := range cookie {
		request.AddCookie(cookie[i])
	}
	resp, _ := client.Do(request)
	return resp, resp.Cookies()
}

func readRespBody(resp *http.Response) []byte {
	bodyBytes, _ := io.ReadAll(resp.Body)
	return bodyBytes
}

func unmarshalInt(data []byte, log logger.Logger) int {
	var beInt int
	if err := json.Unmarshal(data, &beInt); err != nil {
		log.Error(err.Error())
	}

	return beInt
}

func unmarshalChar(data []byte, log logger.Logger) model.Character {
	var char model.Character

	if err := json.Unmarshal(data, &char); err != nil {
		log.Error(err.Error())
	}

	return char
}

func unmarshalChars(data []byte, log logger.Logger) []model.Character {
	var chars []model.Character
	if err := json.Unmarshal(data, &chars); err != nil {
		log.Error(err.Error())
	}

	return chars
}

func loadApiScheme(log logger.Logger) config.Endpoints {
	var apiScheme config.Endpoints
	resultJson := <-app.Ch
	if err := json.Unmarshal(resultJson, &apiScheme); err != nil {
		log.Error(err.Error())
	}
	return apiScheme
}

func loadEnv(log logger.Logger) (config.LogMsg, config.MsgToUser) {
	re := regexp.MustCompile(`^(.*` + "ArchaicReverie" + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))
	godotenv.Load(string(rootPath) + `/.env`)
	app.YmlFilename = os.Getenv("CONFIG_FILE")
	logMsg, msgToUser, _, _ := config.InitStrSet(app.YmlFilename, log)

	return logMsg, msgToUser
}

func reqBody(log logger.Logger, v interface{}) *bytes.Buffer {
	jsonStr, err := json.Marshal(v)
	if err != nil {
		log.Panicf("%s: %s", log.CallInfoStr(), err.Error())
	}
	return bytes.NewBuffer(jsonStr)
}

func generateTestUser() model.User {
	return model.User{
		Username: RandStr(7),
		Password: RandStr(7),
	}
}

func generateChar(ownedId int) model.Character {
	return model.Character{
		OwnerId: ownedId,
		Name:    RandStr(4),
		Growth:  RandInt(150, 180),
		Weight:  RandInt(40, 90),
	}
}

func generateAction() model.Action {
	return model.Action{
		InAction: "jump",
		Jump: model.Jump{
			SquatDepth:   RandInt(0, 1),
			ArmAmplitude: RandInt(0, 1),
			BodyTilt:     RandInt(0, 1),
			RunUp:        RandInt(0, 1),
		},
	}
}

func RandInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func RandStr(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
