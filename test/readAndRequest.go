package test

import (
	"encoding/json"
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
	"github.com/Aserose/ArchaicReverie/pkg/logger"
	"io"
	"net/http"
)

type readAndRequest struct {
	log logger.Logger
}

func (r readAndRequest) unmarshalAvailableItems(data []byte) model.Items {
	var items model.Items

	//unmarshal the list of all food and weapons
	if err := json.Unmarshal(data, &items); err != nil {
		r.log.Errorf("%s:%s", r.log.CallInfoStr(), err.Error())
	}

	return items
}

func (r readAndRequest) cheapestOrder(availableItems model.Items) model.Items {
	var (
		orderItems     model.Items
		cheapestWeapon model.Weapon
		cheapestFood   model.Food
	)

	//looking for weapon with the lowest price
	for _, w := range availableItems.Weapons {
		if cheapestWeapon.Price == 0 || cheapestWeapon.Price > w.Price {
			cheapestWeapon = w
		}
	}

	//looking for food with the lowest price
	for _, f := range availableItems.Foods {
		if cheapestFood.Price == 0 || cheapestFood.Price > f.Price {
			cheapestFood = f
		}
	}

	orderItems.Weapons = append(orderItems.Weapons, cheapestWeapon)
	orderItems.Foods = append(orderItems.Foods, cheapestFood)

	return orderItems
}

func (r readAndRequest) doRequest(client http.Client, method string, url string, body io.Reader, cookie []*http.Cookie) (*http.Response, []*http.Cookie) {
	request, _ := http.NewRequest(method, url, body)
	for i := range cookie {
		request.AddCookie(cookie[i])
	}
	resp, _ := client.Do(request)
	return resp, resp.Cookies()
}

func (r readAndRequest) readRespBody(resp *http.Response) []byte {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		r.log.Panicf("%s: %s", r.log.CallInfoStr(), err.Error())
	}
	return bodyBytes
}

func (r readAndRequest) unmarshalInt(data []byte) int {
	var beInt int
	if err := json.Unmarshal(data, &beInt); err != nil {
		r.log.Panicf(`%s: %s: respBody:  %s:"`, r.log.CallInfoStr(), err.Error(), string(data))
	}

	return beInt
}

func (r readAndRequest) unmarshalChar(data []byte) model.Character {
	var char model.Character

	if err := json.Unmarshal(data, &char); err != nil {
		r.log.Errorf("%s: %s: respBody: %s", r.log.CallInfoStr(), err.Error(), string(data))
	}

	return char
}

func (r readAndRequest) unmarshalChars(data []byte) []model.Character {
	var chars []model.Character
	if err := json.Unmarshal(data, &chars); err != nil {
		r.log.Errorf("%s: %s: %s", r.log.CallInfoStr(), err.Error(), string(data))
	}

	return chars
}
