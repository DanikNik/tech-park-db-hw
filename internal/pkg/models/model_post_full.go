package models

// Полная информация о сообщении, включая связанные объекты.
type PostFull struct {
	Post *Post `json:"post,omitempty"`

	Author *User `json:"author,omitempty"`

	Thread *Thread `json:"thread,omitempty"`

	Forum *Forum `json:"forum,omitempty"`
}
