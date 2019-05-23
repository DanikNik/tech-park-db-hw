package models

// Информация о голосовании пользователя.
type Vote struct {

	// Идентификатор пользователя.
	Nickname string `json:"nickname"`

	// Отданный голос.
	Voice float32 `json:"voice"`
}
