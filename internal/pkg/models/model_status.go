package models

type Status struct {

	// Кол-во пользователей в базе данных.
	User float32 `json:"user"`

	// Кол-во разделов в базе данных.
	Forum float32 `json:"forum"`

	// Кол-во веток обсуждения в базе данных.
	Thread float32 `json:"thread"`

	// Кол-во сообщений в базе данных.
	Post float32 `json:"post"`
}
