package processor

const (
	PROCESSOR_TYPE_PB uint8 = 1
)

type IProcessor interface {
	Unmarshal([]byte) (interface{}, interface{}, error)
	Marshal(interface{},interface{}) ([]byte, error)
	Register(interface{}, interface{})
}
