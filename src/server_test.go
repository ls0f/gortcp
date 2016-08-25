package gortcp

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
	"bytes"
	"errors"
	"sync"
	"net"
	"os"
)

type test struct{

	sync.Mutex
	run bool
	s *Server
}

func (t *test) runServer(){
	t.Lock()
	if t.run == true{
		return
	}
	t.s = &Server{Addr: ":12345", Auth:"123456"}
	go t.s.Listen()
	time.Sleep(500*time.Millisecond)
	t.run = true
	t.Unlock()

}


var ser *test = &test{}

func TestServer_Listen(t *testing.T) {

	defer func(){
		r:=recover()
		assert.True(t, r != nil)
		_, ok := r.(error)
		assert.False(t, ok)
	}()
	go ser.runServer()
	s := Server{Addr: ":12345", Auth:"123456"}
	time.Sleep(500*time.Millisecond)
	s.Listen()

}

func TestServer_Listen2(t *testing.T) {
	go ser.runServer()
	time.Sleep(500*time.Millisecond)
	c := Client{Addr:"127.0.0.1:12345"}
	go c.Connect()
	time.Sleep(500*time.Millisecond)
	assert.Equal(t, len(ser.s.Pool.pool), 1)
	n := ser.s.Pool.getNode(1)
	assert.True(t, n != nil)
}

func TestServer_ListNode(t *testing.T) {
	go ser.runServer()
	time.Sleep(500*time.Millisecond)
	c := Client{Addr:"127.0.0.1:12345"}
	go c.Connect()
	time.Sleep(500*time.Millisecond)
	buf := new(bytes.Buffer)
	wrap := &MessageWrap{rw: buf}
	ser.s.listNode(wrap)
	assert.Contains(t, string(buf.Bytes()), "127.0.0.1")
}

func TestServer_WriteErrorMessage(t *testing.T){

	s := Server{Addr: ":12345", Auth:"123456"}
	err := errors.New("hello,world")
	buf := new(bytes.Buffer)
	wrap := &MessageWrap{rw: buf}
	s.writeErrorMessage(wrap, err)
	assert.Contains(t, string(buf.Bytes()), "hello,world")
	m, _ := wrap.ReadOneMessage()
	assert.Equal(t, m.msgType, uint8(errorMessage))
}

func TestServer_WriteMatchOkMessage(t *testing.T){
	s := Server{Addr: ":12345", Auth:"123456"}
	buf := new(bytes.Buffer)
	wrap := &MessageWrap{rw: buf}
	s.writeMatchOkMessage(wrap)
	assert.Equal(t, uint8(buf.Bytes()[0]), uint8(matchOKMessage))
}

func TestServer_ConnetAuthError(t *testing.T){
	go ser.runServer()
	conn, err := net.Dial("tcp", "127.0.0.1:12345")
	assert.NoError(t, err)
	wrap := &MessageWrap{rw: conn}
	err = wrap.SendOneMessage(&Message{msgType:authMessage, content: []byte("12345")})
	assert.NoError(t, err)
	m, err := wrap.ReadOneMessage()
	assert.NoError(t, err)
	assert.Equal(t, m.msgType, uint8(errorMessage))
	assert.Contains(t, string(m.content), "auth error")
}

func TestServer_MatchError(t *testing.T){
	go ser.runServer()
	conn, err := net.Dial("tcp", "127.0.0.1:12345")
	assert.NoError(t, err)
	wrap := &MessageWrap{rw: conn}
	err = wrap.SendOneMessage(&Message{msgType:authMessage, content: []byte("123456")})
	assert.NoError(t, err)
	err = wrap.SendOneMessage(&Message{msgType:matchNodeMessage, content: []byte("abcdefg")})
	assert.NoError(t, err)
	m, err := wrap.ReadOneMessage()
	assert.NoError(t, err)
	assert.Equal(t, m.msgType, uint8(errorMessage))
	assert.Contains(t, string(m.content), "strconv.ParseUint")
}

func TestServer_IDNotFoundError(t *testing.T){
	go ser.runServer()
	conn, err := net.Dial("tcp", "127.0.0.1:12345")
	assert.NoError(t, err)
	wrap := &MessageWrap{rw: conn}
	err = wrap.SendOneMessage(&Message{msgType:authMessage, content: []byte("123456")})
	assert.NoError(t, err)
	err = wrap.SendOneMessage(&Message{msgType:matchNodeMessage, content: []byte("1234567")})
	assert.NoError(t, err)
	m, err := wrap.ReadOneMessage()
	assert.NoError(t, err)
	assert.Equal(t, m.msgType, uint8(errorMessage))
	assert.Contains(t, string(m.content), "not found")
}

func TestServer_ExecCommand(t *testing.T){
	s := Server{Addr: ":22345", Auth:"123456"}
	go s.Listen()
	time.Sleep(500*time.Millisecond)
	c := Client{Addr:"127.0.0.1:22345"}
	go c.Connect()
	time.Sleep(200*time.Millisecond)
	conn, err := net.Dial("tcp", "127.0.0.1:22345")
	assert.NoError(t, err)
	wrap := &MessageWrap{rw: conn}
	err = wrap.SendOneMessage(&Message{msgType:authMessage, content: []byte("123456")})
	assert.NoError(t, err)
	err = wrap.SendOneMessage(&Message{msgType:matchNodeMessage, content: []byte("1")})
	assert.NoError(t, err)
	err = wrap.SendOneMessage(&Message{msgType: execCmdMessage, content: []byte("hostname")})
	assert.NoError(t, err)
	m, err := wrap.ReadOneMessage()
	assert.NoError(t, err)
	assert.Equal(t, m.msgType, uint8(matchOKMessage))
	m, err = wrap.ReadOneMessage()
	assert.NoError(t, err)
	assert.Equal(t, m.msgType, uint8(execCmdResultMessage))
	host, err := os.Hostname()
	assert.NoError(t, err)
	assert.Equal(t, string(m.content), host+"\n")
}