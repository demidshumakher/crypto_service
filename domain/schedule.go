package domain

import "time"

type ScheduleCfg struct {
	Enabled         bool      `json:"enabled"`
	IntervalSeconds int       `json:"interval_seconds"`
	LastUpdate      time.Time `json:"last_update"`
	NextUpdate      time.Time `json:"next_update"`
}
