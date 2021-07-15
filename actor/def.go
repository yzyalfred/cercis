package actor

const (
	ACTOR_EVENT_CLOSE uint8 = 1		// close actor
)

type IActor interface {
	SendMsg(uint32, uint64, interface{}, interface{})
}