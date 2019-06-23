package models

type ErrorMessage struct {
	Msg string `json:"message"`
}

func NewErrorMessage() ErrorMessage {
	return ErrorMessage{
		Msg: "SOME SHIT HAPPENED",
	}
}
