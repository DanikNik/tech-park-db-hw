package models

import "fmt"

type NotFoundMessage struct {
	Msg string `json:"message"`
}

func GenerateUserNotFoundMessage(id int64) NotFoundMessage {
	return NotFoundMessage{
		Msg: fmt.Sprintf("Can't find user with id #%c\n", string(id)),
	}
}
