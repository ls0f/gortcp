package gortcp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	authMessage           = 1
	connectMessage        = 2
	execCmdMessage        = 3
	execCmdResultMessage  = 31
	uploadFileMessage     = 4
	listNodeMessage       = 5
	listNodeResultMessage = 51
	matchNodeMessage      = 6
)

type Message struct {
	msgType uint8
	length  uint32
	content []byte
}

type MessageWrap struct {
	rw io.ReadWriter
}

func (m *Message) readOneMessage(r io.Reader) (err error) {
	//read msg type
	b1 := make([]byte, 1)
	n, err := r.Read(b1)
	if err != nil {
		logger.Errorf("read message type error: %s", err.Error())
		return
	}
	if n != len(b1) {
		err = errors.New("read message type error: unexpectd length")
		logger.Error(err)
		return
	}
	m.msgType = uint8(b1[0])

	//read length
	b2 := make([]byte, 4)
	n, err = r.Read(b2)
	if err != nil {
		logger.Errorf("read message length error: %s", err.Error())
		return
	}
	if n != len(b2) {
		err = errors.New("read message length error: unexpectd length")
		logger.Error(err)
		return
	}
	m.length = binary.BigEndian.Uint32(b2)

	//read content
	if m.length == 0 {
		return
	}
	m.content = make([]byte, m.length)
	n, err = r.Read(m.content)
	if err != nil {
		logger.Errorf("read message length error: %s", err.Error())
		return
	}
	if n != len(m.content) {
		err = errors.New("read message length error: unexpectd length")
		return
	}
	return

}

func (m *Message) sendOneMessage(w io.Writer) (n int, err error) {
	m.length = uint32(len(m.content))
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, m.msgType)
	binary.Write(buf, binary.BigEndian, m.length)
	buf.Write(m.content)
	n, err = w.Write(buf.Bytes())
	if err != nil {
		logger.Error(err)
	}
	return
}

func (wrap *MessageWrap) ReadOneMessage() (m *Message, err error) {
	m = &Message{}
	err = m.readOneMessage(wrap.rw)
	return
}

func (wrap *MessageWrap) ReadTheSpecialTypeMessage(msgType uint8) (m *Message, err error) {

	m, err = wrap.ReadOneMessage()
	if err != nil {
		return
	}
	if m.msgType != msgType {
		errStr := fmt.Sprintf("not the expected msgType, expected: %d, actual: %d", msgType, m.msgType)
		err = errors.New(errStr)
		m = nil
		return
	}
	return
}

func (wrap *MessageWrap) SendOneMessage(m *Message) (err error) {
	_, err = m.sendOneMessage(wrap.rw)
	return
}
