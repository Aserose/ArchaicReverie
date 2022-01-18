package model

type Items struct {
	Weapons []Weapon `json:"weapon,omitempty"`
	Foods   []Food   `json:"food,omitempty"`
}

type Weapon struct {
	Name   string `json:"name" db:"name"`
	Price  int    `json:"price" db:"price"`
	Sharp  int    `json:"sharp" db:"sharp"`
	Weight int    `json:"weight" db:"weight"`
}

type Food struct {
	Name      string `json:"name" db:"name"`
	Price     int    `json:"price" db:"price"`
	RestoreHp int    `json:"restoreHp" db:"restore_hp"`
}
