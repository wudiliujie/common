package baseconf

import (
	"log"
	"time"
)

var (
	// log conf
	LogFlag = log.LstdFlags

	// skeleton conf
	GoLen              = 10000
	TimerDispatcherLen = 10000
	AsyncCallLen       = 10000
	ChanRPCLen         = 10000

	// cluster conf
	HeartBeatInterval = 5

	PendingWriteNum        = 3000
	MaxMsgLen       uint32 = 409600
	HTTPTimeout            = 10 * time.Second
	LenMsgLen              = 2
	LittleEndian           = false

	// agent conf
	AgentGoLen              = 500
	AgentTimerDispatcherLen = 50
	AgentAsyncCallLen       = 50
	AgentChanRPCLen         = 50

	// skeleton conf

)
