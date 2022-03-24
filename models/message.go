package models

import "time"

type ChangeMessage struct {
	Category string    `json:"category,omitempty"`
	Type     string    `json:"type,omitempty"`
	Info     string    `json:"info,omitempty"`
	Time     time.Time `json:"time"`
}
