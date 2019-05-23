package models

// Информация о пользователе.
type UserUpdate struct {

	// Полное имя пользователя.
	Fullname string `json:"fullname,omitempty"`

	// Описание пользователя.
	About string `json:"about,omitempty"`

	// Почтовый адрес пользователя (уникальное поле).
	Email string `json:"email,omitempty"`
}
