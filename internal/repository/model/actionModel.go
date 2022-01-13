package model

type Action struct {
	InAction string `json:"inAction"`
	Jump     `json:"jump"`
	Run      `json:"run"`
	Hit      `json:"hit"`
}

type Jump struct {
	SquatDepth   int `json:"squatDepth"`
	ArmAmplitude int `json:"armAmplitude"`
	BodyTilt     int `json:"bodyTilt"`
	RunUp        int `json:"runUp"`
}

type Run struct {
	BodyTilt int `json:"bodyTilt"`
}

type Hit struct {
	Backswing int `json:"backswing"`
	Straight  int `json:"straight"`
}

type ActionResult struct {
	Name       string `json:"name" db:"name"`
	DamageType `json:"damage_type" db:"damage_type"`
}

type DamageType struct {
	DamageHP int `json:"damage_hp" db:"damage_hp"`
	DamageMP int `json:"damage_mp" db:"damage_mp"`
}
