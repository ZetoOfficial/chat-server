package domain

type Message struct {
	User    string `json:"user"`
	Content string `json:"content"`
}

type UserList struct {
	Type string   `json:"type"`
	Data []string `json:"data"`
}

type BroadcastMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}
