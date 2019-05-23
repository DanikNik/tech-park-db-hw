package models

// Сообщение для обновления ветки обсуждения на форуме. Пустые параметры остаются без изменений.
type ThreadUpdate struct {

	// Заголовок ветки обсуждения.
	Title string `json:"title,omitempty"`

	// Описание ветки обсуждения.
	Message string `json:"message,omitempty"`
}
