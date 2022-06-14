package common

import "time"

const (
	NotificationNamespace = "notification"
	NotificationRoom      = "notification"
	NotificationEvent     = "notification"
)

type Notification struct {
	Id          string    `json:"id"`
	Event       string    `json:"event"`
	Message     string    `json:"message"`
	Data        any       `json:"data"`
	CreatedTime time.Time `json:"createdTime"`
}
