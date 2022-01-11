package model

type Action struct {
	InAction string `json:"inAction"`
	Jump     `json:"jump"`
}

type Jump struct {
	SquatDepth   int `json:"squatDepth"`
	ArmAmplitude int `json:"armAmplitude"`
	BodyTilt     int `json:"bodyTilt"`
	RunUp        int `json:"runUp"`
}

type ActionResult struct {
	Name       string `json:"name" db:"name"`
	DamageType `json:"damage_type" db:"damage_type"`
}

type DamageType struct {
	DamageHP int `json:"damage_hp" db:"damage_hp"`
	DamageMP int `json:"damage_mp" db:"damage_mp"`
}
