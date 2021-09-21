package model

import (
	"time"
)

const (
	DefaultAuctionDuration = 48 * 60 * 60
)

func AuctionIsOpen(startedAt *time.Time, duration int) bool {
	return startedAt.Add(time.Second * time.Duration(duration)).After(time.Now())
}
