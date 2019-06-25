package models

type Status struct {

	// Кол-во пользователей в базе данных.
	User int32 `json:"user"`

	// Кол-во разделов в базе данных.
	Forum int32 `json:"forum"`

	// Кол-во веток обсуждения в базе данных.
	Thread int32 `json:"thread"`

	// Кол-во сообщений в базе данных.
	Post int32 `json:"post"`
}
