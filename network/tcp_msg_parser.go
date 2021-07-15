package network

import (
	"errors"
	"github.com/yzyalfred/cercis/utils"
	"io"
)

type TCPMsgParser struct {
	maxMsgLen uint32

	littleEndian bool
}

func (this *TCPMsgParser) Read(client *TCPClient) ([]byte, error) {
	msgLenBuf := make([]byte, MSG_LEN_SIZE)

	// read len
	if _, err := io.ReadFull(client.conn, msgLenBuf); err != nil {
		return nil, err
	}

	msgLen := utils.ByteToUin32(msgLenBuf, this.littleEndian)
	if msgLen > this.maxMsgLen {
		return nil, errors.New("msg too long")
	}

	msgBuf := make([]byte, msgLen)
	if _, err := io.ReadFull(client.conn, msgBuf); err != nil {
		return nil, err
	}

	return msgBuf, nil
}

func (this *TCPMsgParser) Write(client *TCPClient, args ...interface{}) error {
	var msgLen uint32
	for _, v := range args {
		switch v.(type) {
		case uint32:
			msgLen += 4
		case []byte:
			msgLen += uint32(len(v.([]byte)))
		case uint64:
			msgLen += 8
		default:
			return errors.New("error type")
		}
	}

	if msgLen > this.maxMsgLen {
		return errors.New("msg too long")
	}

	buf := make([]byte, MSG_LEN_SIZE+msgLen)
	utils.PutUint32ToByte(buf, msgLen, this.littleEndian)

	offset := MSG_LEN_SIZE
	for _, v := range args {
		switch v.(type) {
		case uint32:
			utils.PutUint32ToByte(buf[offset:], v.(uint32), this.littleEndian)
			offset += 4
		case []byte:
			copy(buf[offset:], v.([]byte))
			offset += len(v.([]byte))
		case uint64:
			utils.PutUint64ToByte(buf[offset:], v.(uint64), this.littleEndian)
			offset += 8
		}
	}

	client.Send(buf)

	return nil
}
