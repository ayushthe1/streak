package util

import "github.com/ayushthe1/streak/models"

var BroadcastKafkaEvent = make(chan *models.Notification)
