package gortcp

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"bytes"
	"encoding/binary"
)

func TestMessage_ReadOneMessage(t *testing.T) {
	m1 := &Message{msgType:1, length:2, content: []byte{'a', 'b'}}
	m2 := &Message{}

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, m1.msgType)
	binary.Write(buf, binary.BigEndian, m1.length)
	buf.Write(m1.content)

	err := m2.ReadOneMessage(buf)
	assert.NoError(t, err)
	assert.Equal(t, uint8(1), m2.msgType)
	assert.Equal(t, uint32(2), m2.length)
	assert.Equal(t, []byte{'a', 'b'}, m2.content)
}

func TestMessage_SendOneMessage(t *testing.T) {

	m1 := &Message{msgType:1, length:2, content: []byte{'a', 'b'}}
	buf := new(bytes.Buffer)
	n, err := m1.SendOneMessage(buf)
	assert.NoError(t, err)
	assert.Equal(t, n, 7)
	assert.Equal(t, buf.Bytes(), []byte{1,0, 0, 0, 2, 97, 98})
}

func TestMessage_ReadOneMessage2(t *testing.T) {

	m1 := &Message{msgType:1, length:2, content: []byte{'a', 'b'}}
	buf := new(bytes.Buffer)
	_, err := m1.SendOneMessage(buf)
	assert.NoError(t, err)

	m2 := &Message{}
	err = m2.ReadOneMessage(buf)
	assert.NoError(t, err)

	assert.Equal(t, m1.msgType, m2.msgType)
	assert.Equal(t, m1.length, m2.length)
	assert.Equal(t, m1.content, m2.content)
}