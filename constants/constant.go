package constants

import "time"

const (
	GinModelsDepth        = 3
	SyncThreshold         = 0
	MaxMoveSize           = 4
	MaxPlayerSize         = 5
	MaxPoint              = 5
	MaxWSUpdateBufferSize = 10
	MaxGameLogSize        = 100
	SleepDuration         = 500 * time.Millisecond
	TimeoutDuration       = 3 * time.Second
	SaveDuration          = 5 * time.Second
)
