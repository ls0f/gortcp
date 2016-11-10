package gortcp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

const (
	authMessage            = 1
	connectMessage         = 2
	connectOKmessage       = 21
	execCmdMessage         = 3
	execCmdResultMessage   = 31
	uploadFileMessage      = 4
	fileInfoMessage        = 41
	uploadDoneMessage      = 42
	replyUploadDoneMessage = 43
	listNodeMessage        = 5
	listNodeResultMessage  = 51
	matchNodeMessage       = 6
	matchOKMessage         = 61
	errorMessage           = 7
	pingMessage            = 8
	pingOKMessage          = 9
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
	_, err = io.ReadFull(r, b1)
	if err != nil {
		logger.Errorf("read message type error: %s", err.Error())
		return
	}
	m.msgType = uint8(b1[0])

	//read length
	b2 := make([]byte, 4)
	_, err = io.ReadFull(r, b2)
	if err != nil {
		logger.Errorf("read message length error: %s", err.Error())
		return
	}

	m.length = binary.BigEndian.Uint32(b2)

	//read content
	if m.length == 0 {
		return
	}
	m.content = make([]byte, m.length)
	_, err = io.ReadFull(r, m.content)
	if err != nil {
		logger.Errorf("read message length error: %s", err.Error())
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
	if n, ok := wrap.rw.(*net.TCPConn); ok {
		n.SetReadDeadline(time.Now().Add(time.Second * ReadTimeOut))
	}
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
	if n, ok := wrap.rw.(*net.TCPConn); ok {
		n.SetWriteDeadline(time.Now().Add(time.Second * WriteTimeOut))
	}
	_, err = m.sendOneMessage(wrap.rw)
	return
}
