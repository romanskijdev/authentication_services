package variables

import (
	"time"
)

// context timeout
const (
	ContextTimeoutShort = 10 * time.Second
	ContextTimeoutBase  = 15 * time.Second
	ContextTimeoutLong  = 30 * time.Second
)

const MaxMsgGRPCSize = 100 * 1024 * 1024 // 100MB в байтах
