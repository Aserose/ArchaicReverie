package model

type Inventory struct {
	Weapons    []Weapon `json:"weapon" db:"weapon"`
	CoinAmount int      `json:"coinAmount" db:"coin_amount"`
}

type Weapon struct {
	Name   string `json:"name" db:"name"`
	Sharp  int    `json:"sharp" db:"sharp"`
	Weight int    `json:"weight" db:"weight"`
}

type Food struct {
	Name      string `json:"name" db:"name"`
	Price     int    `json:"price" db:"price"`
	RestoreHp int    `json:"restoreHp" db:"restore_hp"`
}
