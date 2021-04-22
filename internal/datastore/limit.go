package datastore

import "github.com/AlekSi/pointer"

const (
	DefaultLimit = uint64(50)
	DefaultMaxLimit = uint64(100)
)

type LimitOpts struct {
	Limit  *uint64
	Offset *uint64
}

func NewLimitOpts(offset, limit uint64) *LimitOpts {
	newLimit := limit
	if newLimit == 0 {
		newLimit = DefaultLimit
	}
	if newLimit > DefaultMaxLimit {
		newLimit = DefaultLimit
	}
	return &LimitOpts{
		Offset: pointer.ToUint64(offset),
		Limit: pointer.ToUint64(newLimit),
	}
}
