package entities

import "sync/atomic"

type ID = int64

var sequenceNumber ID = 0

const InvalidID ID = -1

func NewID() ID {
	return atomic.AddInt64(&sequenceNumber, 1)
}
