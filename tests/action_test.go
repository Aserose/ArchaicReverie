package tests

import (
	"github.com/Aserose/ArchaicReverie/internal/app"
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
	"github.com/Aserose/ArchaicReverie/pkg/logger"
	cv "github.com/smartystreets/goconvey/convey"
	"net/http"
	"strings"
	"testing"
)

func TestAction(t *testing.T) {
	logs := logger.NewLogger()
	_, msgToUser, charConfig := loadEnv(logs)
	go app.Start(1)
	apiScheme := loadApiScheme(logs)

	cv.Convey("setup", t, func() {
		var (
			temp            = NewTemplates(logs)
			numberCharLimit = charConfig.Restriction.NumberCharLimit
			client          = http.Client{}
			cookie          = temp.auth.authorizeUserWithCharsAndSelect(client, apiScheme, numberCharLimit)
			resp            *http.Response
			testUser        = generateTestUser()
		)

		cv.Convey("beginActionScene", func() {
			resp, _ = temp.doRequest(
				client,
				strings.Split(apiScheme.ActionEndpoints.ActionScene, " ")[0],
				strings.Split(apiScheme.ActionEndpoints.ActionScene, " ")[1],
				reqBody(logs, testUser),
				cookie)

			cv.So(string(temp.readRespBody(resp)), cv.ShouldNotBeEmpty)

			var result string
			for {
				resp, updCookie := temp.doRequest(
					client,
					strings.Split(apiScheme.ActionEndpoints.ActionScene, " ")[0],
					strings.Split(apiScheme.ActionEndpoints.ActionScene, " ")[1],
					reqBody(logs, generateAction()),
					cookie)

				if len(updCookie) > 0 {
					if updCookie[0].Value != "" {
						cookie = updCookie
					}
				}

				result = string(temp.readRespBody(resp))
				logs.Print(result)

				if result == msgToUser.ActionMsg.JumpOver {
					cv.So(result, cv.ShouldEqual, msgToUser.ActionMsg.JumpOver)
					break
				}
				if result == msgToUser.ActionMsg.LowHP {
					cv.So(result, cv.ShouldEqual, "low health points")
					break
				}
				cv.So(result, cv.ShouldNotBeEmpty)
			}

			cv.Convey("getInfoAboutChar", func() {
				resp, _ = temp.doRequest(
					client,
					strings.Split(apiScheme.ActionEndpoints.CharacterMenu, " ")[0],
					strings.Split(apiScheme.ActionEndpoints.CharacterMenu, " ")[1],
					reqBody(logs, testUser),
					cookie)

				selectedChar := temp.unmarshalChar(temp.readRespBody(resp))
				logs.Printf("%d health points left out of %d", selectedChar.RemainHealth, selectedChar.ThresholdHealth)
				resp, _ = temp.doRequest(
					client,
					strings.Split(apiScheme.ActionEndpoints.Restock, " ")[0],
					strings.Split(apiScheme.ActionEndpoints.Restock, " ")[1],
					reqBody(logs, testUser),
					cookie)

				cheapestOrder := temp.cheapestOrder(temp.unmarshalAvailableItems(temp.readRespBody(resp)))

				if selectedChar.RemainHealth < selectedChar.ThresholdHealth {
					logs.Printf("restore health by %d points with %s", cheapestOrder.Foods[0].RestoreHp, cheapestOrder.Foods[0].Name)

					resp, cookie = temp.doRequest(
						client,
						strings.Split(apiScheme.ActionEndpoints.Restock, " ")[0],
						strings.Split(apiScheme.ActionEndpoints.Restock, " ")[1],
						reqBody(logs, model.Items{Foods: cheapestOrder.Foods}),
						cookie)

					resp, _ = temp.doRequest(
						client,
						strings.Split(apiScheme.ActionEndpoints.CharacterMenu, " ")[0],
						strings.Split(apiScheme.ActionEndpoints.CharacterMenu, " ")[1],
						reqBody(logs, testUser),
						cookie)

					selectedChar := temp.unmarshalChar(temp.readRespBody(resp))
					logs.Printf("%d health points left out of %d", selectedChar.RemainHealth, selectedChar.ThresholdHealth)
				}

				cv.Convey("purchase weapon", func() {

					resp, cookie = temp.doRequest(
						client,
						strings.Split(apiScheme.ActionEndpoints.Restock, " ")[0],
						strings.Split(apiScheme.ActionEndpoints.Restock, " ")[1],
						reqBody(logs, model.Items{Weapons: cheapestOrder.Weapons}),
						cookie)

					resp, _ = temp.doRequest(
						client,
						strings.Split(apiScheme.ActionEndpoints.CharacterMenu, " ")[0],
						strings.Split(apiScheme.ActionEndpoints.CharacterMenu, " ")[1],
						reqBody(logs, testUser),
						cookie)

					selectedChar := temp.unmarshalChar(temp.readRespBody(resp))

					logs.Print("character inventory ", selectedChar.Inventory.Weapons," the amount of coins: ", selectedChar.Inventory.CoinAmount)

					cv.Convey("discard weapon", func() {

						resp, cookie = temp.doRequest(
							client,
							strings.Split(apiScheme.ActionEndpoints.CharacterMenu, " ")[0],
							strings.Split(apiScheme.ActionEndpoints.CharacterMenu, " ")[1],
							reqBody(logs, model.Items{Weapons: cheapestOrder.Weapons}),
							cookie)

						selectedChar := temp.unmarshalChar(temp.readRespBody(resp))
						logs.Print("character inventory ", selectedChar.Inventory.Weapons," the amount of coins: ", selectedChar.Inventory.CoinAmount)


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
}
