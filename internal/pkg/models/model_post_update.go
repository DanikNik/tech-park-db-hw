package models

// Сообщение для обновления сообщения внутри ветки на форуме. Пустые параметры остаются без изменений.
type PostUpdate struct {

	// Собственно сообщение форума.
	Message string `json:"message,omitempty"`
}
