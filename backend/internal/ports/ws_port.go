package ports

type WSPort interface {
	BroadcastMessage(msg interface{})
}
