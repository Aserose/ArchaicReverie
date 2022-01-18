package model

type Inventory struct {
	Weapons    []Weapon `json:"weapon" db:"weapon"`
	CoinAmount int      `json:"coinAmount" db:"coin_amount"`
}
