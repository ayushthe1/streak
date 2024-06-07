package models

type Notification struct {
	Type string `json:"type"`
	// username to send the notification
	Username string `json:"username"`
	Message  string `json:"message"`
}
