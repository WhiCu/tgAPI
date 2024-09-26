package types

type UpdatesResponse struct {
	Ok     bool     `json : "ok"`
	Result []Update `json : "result"`
}

type Update struct {
	Update_id *int     `json : "update_id"`
	Message   *Message `json : "message"`
}

type Message struct {
	Message_id *int    `json : "message_id"`
	From       *User   `json : "from"`
	Chat       *Chat   `json : "chat"`
	Date       *int    `json : "date"`
	Text       *string `json : "text"`
}
type User struct {
	Id       *int    `json : "id"`
	Is_bot   *bool   `json : "is_bot"`
	UserName *string `json : "username"`
}
type Chat struct {
	Id       *int    `json : "id"`
	UserName *string `json : "username"`
}
