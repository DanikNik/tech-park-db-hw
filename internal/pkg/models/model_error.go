package models

type ModelError struct {
	// Текстовое описание ошибки. В процессе проверки API никаких проверок на содерижимое данного описание не делается.
	Message string `json:"message,omitempty"`
}
