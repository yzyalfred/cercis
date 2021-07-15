package processor

import (
	"github.com/gogo/protobuf/proto"
	"github.com/yzyalfred/cercis/log"
	"github.com/yzyalfred/cercis/utils"
	"reflect"
)

const (
	LITTLE_ENDIAN = true // little_endian default
	MSGID_SIZE    = 4    // msgId size
)

type PBProcessor struct {
	msgInfoList map[uint32]reflect.Type

	littleEndian bool
}

func NewPBProcessor() *PBProcessor {
	return &PBProcessor{
		msgInfoList:  make(map[uint32]reflect.Type),
		littleEndian: LITTLE_ENDIAN,
	}
}

func (this *PBProcessor) SetByteOrder(littleEndian bool) {
	this.littleEndian = littleEndian
}

func (this *PBProcessor) Unmarshal(msgData []byte) (interface{}, interface{}, error) {
	// activity you should check before
	if len(msgData) <= MSGID_SIZE {
		return nil, nil, utils.ERROR_MSGID_SHORT
	}

	msgId := utils.ByteToUin32(msgData, this.littleEndian)
	msgType, ok := this.msgInfoList[msgId]
	if !ok {
		log.Warn("msgId not found: ", msgId)
		return nil, nil, utils.ERROR_NOT_FOUND
	}

	msgValue := reflect.New(msgType.Elem()).Interface()
	msg := msgValue.(proto.Message)
	err := proto.Unmarshal(msgData[MSGID_SIZE:], msg)
	if err != nil {
		log.Warn("unmarshall error, msgId: ", msgId, "data: ", msgData)
		return nil, nil, err
	}

	return msgId, msg, nil
}

// add head: msgId
func (this *PBProcessor) Marshal(msgId interface{}, msg interface{}) ([]byte, error) {
	msgData, err := proto.Marshal(msg.(proto.Message))
	if err != nil {
		return nil, err
	}

	buf := make([]byte, len(msgData)+MSGID_SIZE)
	utils.PutUint32ToByte(buf, msgId.(uint32), this.littleEndian)
	copy(buf[MSGID_SIZE:], msgData)

	return buf, nil
}

func (this *PBProcessor) Register(msgId interface{}, msg interface{}) {
	reflectType := reflect.TypeOf(msg)
	this.msgInfoList[msgId.(uint32)] = reflectType
}
