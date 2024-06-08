package channels

import "github.com/ayushthe1/streak/models"

var BroadcastKafkaNotification = make(chan *models.Notification)

var BroadcastKafkaActivity = make(chan *models.ActivityEvent)
