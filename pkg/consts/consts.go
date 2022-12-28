package consts

import "time"

const (
	SyncThreshold         = 0
	MaxMoveSize           = 4
	MaxPlayerSize         = 5
	MaxWSUpdateBufferSize = 10
	MaxGameLogSize        = 100
	SleepDuration         = 500 * time.Millisecond
	TimeoutDuration       = 3 * time.Second
)
