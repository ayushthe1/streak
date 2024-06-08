package models

type ActivityEvent struct {
	Type      string `json:"type"`
	Username  string `json:"username"`
	Action    string `json:"action"`
	Timestamp int64  `json:"timestamp"`
	Details   string `json:"details"`
}
