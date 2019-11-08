package network

type Agent interface {
	Run()
	OnClose()
	GetAutoReconnect() bool
	SetAutoReconnect(v bool)
}
