package gate

import (
	"github.com/wudiliujie/common/network"
	"time"
)

func Connect(addr string, tag interface{}) {
	client := new(network.WSClient)
	client.Addr = "ws://" + addr
	client.NewAgent = func(conn *network.WSConn, tag interface{}) network.Agent {
		return newAgent(conn, tag)
	}
	client.ConnNum = 1
	client.ConnectInterval = 5 * time.Second
	client.AutoReconnect = true
	client.PendingWriteNum = 100
	client.Tag = tag
	client.Start()
}
