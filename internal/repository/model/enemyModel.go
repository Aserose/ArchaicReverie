package model

type Enemy struct {
	Name         string `json:"name" db:"name"`
	Class        int    `json:"class" db:"class"`
	Weight       int    `json:"weight" db:"weight"`
	Growth       int    `json:"growth" db:"growth"`
	RemainHealth int    `json:"remain_health" db:"remain_health"`
	RemainEnergy int    `json:"remain_energy" db:"remain_energy"`
	Weapons      Weapon `json:"weapon" db:"weapon"`
	CoinAmount   int    `json:"-"`
}
