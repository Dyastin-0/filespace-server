package types

type expiration struct {
	Value int64  `json:"value"`
	Str   string `json:"text"`
}

type ShareBody struct {
	Email string     `json:"email"`
	File  string     `json:"file"`
	Exp   expiration `json:"expiration"`
}
