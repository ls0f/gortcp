package gortcp

import (
	"testing"

	"os"

	"bytes"

	"github.com/stretchr/testify/assert"
)

func TestControl_ExecCommand(t *testing.T) {
	c := Client{}
	res, err := c.execCommand("pwd")
	assert.NoError(t, err)
	d, _ := os.Getwd()
	assert.Contains(t, string(res), d)
}

func TestControl_ExecCommand2(t *testing.T) {
	c := Client{}
	_, err := c.execCommand("")
	assert.Error(t, err)
}

func TestControl_handlerFileInfo(t *testing.T) {

	c := Client{}
	c.handlerFileInfo([]byte("e10adc3949ba59abbe56e057f20f883e/tmp/test"))
	assert.Equal(t, c.fm.dstPath, "/tmp/test")
	assert.Equal(t, c.fm.md5, "e10adc3949ba59abbe56e057f20f883e")

}

func TestClient_handlerDownload(t *testing.T) {

	buf := new(bytes.Buffer)
	c := &Client{Addr: "123456", wrap: &MessageWrap{rw: buf}}
	c.handlerDownload([]byte("/tmp/notfound"))
	m, _ := c.wrap.ReadOneMessage()
	assert.Contains(t, string(m.content), "no such file or directory\n")
}

func TestClient_handlerDownload2(t *testing.T) {

	buf := new(bytes.Buffer)
	c := &Client{Addr: "123456", wrap: &MessageWrap{rw: buf}}
	c.handlerDownload([]byte("/tmp"))
	m, _ := c.wrap.ReadOneMessage()
	assert.Contains(t, string(m.content), "tmp is a dir")
}

func TestClient_handlerDownload3(t *testing.T) {

	f, _ := os.Create("/tmp/test.txt")
	f.WriteString("123456")
	f.Close()

	buf := new(bytes.Buffer)
	c := &Client{Addr: "123456", wrap: &MessageWrap{rw: buf}}
	c.handlerDownload([]byte("/tmp/test.txt"))
	m, _ := c.wrap.ReadOneMessage()
	assert.Equal(t, m.msgType, uint8(uploadFileMessage))
	assert.Equal(t, string(m.content), "123456")
	m, _ = c.wrap.ReadOneMessage()
	assert.Equal(t, m.msgType, uint8(downloadDoneMessage))
	assert.Contains(t, string(m.content), "e10adc3949ba59abbe56e057f20f883e")
}
