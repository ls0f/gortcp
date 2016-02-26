package gortcp

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessage_ReadOneMessage(t *testing.T) {
	buf := new(bytes.Buffer)
	m := &Message{}
	buf.Write([]byte{1, 0, 0, 0, 2, 97, 98})

	err := m.readOneMessage(buf)
	assert.NoError(t, err)
	assert.Equal(t, uint8(1), m.msgType)
	assert.Equal(t, uint32(2), m.length)
	assert.Equal(t, []byte{'a', 'b'}, m.content)
}

func TestMessage_SendOneMessage(t *testing.T) {

	m := &Message{msgType: 1, length: 2, content: []byte{'a', 'b'}}
	buf := new(bytes.Buffer)
	n, err := m.sendOneMessage(buf)
	assert.NoError(t, err)
	assert.Equal(t, n, 7)
	assert.Equal(t, buf.Bytes(), []byte{1, 0, 0, 0, 2, 97, 98})
}

func TestMessage_ReadOneMessage2(t *testing.T) {

	m1 := &Message{msgType: 1, length: 2, content: []byte{'a', 'b'}}
	buf := new(bytes.Buffer)
	_, err := m1.sendOneMessage(buf)
	assert.NoError(t, err)

	m2 := &Message{}
	err = m2.readOneMessage(buf)
	assert.NoError(t, err)

	assert.Equal(t, m1.msgType, m2.msgType)
	assert.Equal(t, m1.length, m2.length)
	assert.Equal(t, m1.content, m2.content)
}

func TestMessageWrap_RendOneMessage(t *testing.T) {

	buf := new(bytes.Buffer)
	wrap := &MessageWrap{rw: buf}
	wrap.rw.Write([]byte{1, 0, 0, 0, 2, 97, 98})
	m, err := wrap.ReadOneMessage()
	assert.NoError(t, err)
	assert.Equal(t, uint8(1), m.msgType)
	assert.Equal(t, uint32(2), m.length)
	assert.Equal(t, []byte{'a', 'b'}, m.content)
}

func TestMessageWrap_SendOneMessage(t *testing.T) {

	m := &Message{msgType: 1, length: 2, content: []byte{'a', 'b'}}
	buf := new(bytes.Buffer)
	wrap := &MessageWrap{rw: buf}
	err := wrap.SendOneMessage(m)
	assert.NoError(t, err)
	assert.Equal(t, buf.Bytes(), []byte{1, 0, 0, 0, 2, 97, 98})

}

func TestMessageWrap_ReadTheSpecialTypeMessage(t *testing.T) {

	buf := new(bytes.Buffer)
	wrap := &MessageWrap{rw: buf}
	wrap.rw.Write([]byte{1, 0, 0, 0, 2, 97, 98})
	_, err := wrap.ReadTheSpecialTypeMessage(2)
	assert.Error(t, err)
}
