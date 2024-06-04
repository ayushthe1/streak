package models

type Chat struct {
	Id        uint   `json:"id"`
	From      string `json:"from"`
	To        string `json:"to"`
	Msg       string `json:"message"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

type ContactList struct {
	Username     string `json:"username"`
	LastActivity int64  `json:"last_activity"`
}
