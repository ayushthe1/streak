package models

type ActivityEvent struct {
	Id        uint   `json:"id"`
	Type      string `json:"type"`
	Username  string `json:"username"`
	Action    string `json:"action"`
	Timestamp int64  `json:"timestamp"`
	Details   string `json:"details"`
}

type ChatEvent struct {
	Type    string `json:"type"`
	ChatMsg *Chat  `json:"chat"`
}
