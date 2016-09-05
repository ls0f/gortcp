package gortcp

import (
	"bytes"
	"testing"

	"os"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestControl_auth(t *testing.T) {
	buf := new(bytes.Buffer)
	c := &Control{Auth: "123", wrap: &MessageWrap{rw: buf}}
	c.auth()
	assert.Equal(t, buf.Bytes(), []byte{authMessage, 0, 0, 0, 3, 49, 50, 51})
}

func TestControl_print(t *testing.T) {

	buf := new(bytes.Buffer)
	c := &Control{Auth: "123", wrap: &MessageWrap{rw: buf}}
	buf.Write([]byte{authMessage, 0, 0, 0, 3, 49, 50, 51})
	buf2 := new(bytes.Buffer)
	c.print(buf2)
	assert.Equal(t, buf2.Bytes(), []byte{49, 50, 51})

}

func TestControl_UploadFile(t *testing.T) {
	s := &Server{Addr: ":12346", Auth: "123456"}
	go s.Listen()
	time.Sleep(100 * time.Millisecond)
	c := &Control{Auth: "123456", Addr: "127.0.0.1:12346"}
	a := &Client{Addr: "127.0.0.1:12346"}
	go a.Connect()
	time.Sleep(100 * time.Millisecond)
	f, err := os.Create("/tmp/test1")
	assert.NoError(t, err)
	buf := make([]byte, 1024)
	for i := 0; i < 1024; i++ {
		f.Write(buf)
	}
	f.Close()
	c.UploadFile("1", "/tmp/test1", "/tmp/test2")
	m1, err := MD5sum("/tmp/test1")
	assert.NoError(t, err)
	m2, err := MD5sum("/tmp/test2")
	assert.NoError(t, err)
	assert.Equal(t, m1, m2)
	os.Remove("/tmp/test1")
	os.Remove("/tmp/test2")
}

func TestControl_Forward(t *testing.T) {
	s := &Server{Addr: ":12347", Auth: "123456"}
	go s.Listen()
	time.Sleep(100 * time.Millisecond)
	c := &Control{Auth: "123456", Addr: "127.0.0.1:12347"}
	a := &Client{Addr: "127.0.0.1:12347"}
	go a.Connect()
	time.Sleep(100 * time.Millisecond)
	go c.Forward("1", "127.0.0.1:22347", "127.0.0.1:12347")
	time.Sleep(100 * time.Millisecond)
	// request forward port
	c2 := &Control{Auth: "123456", Addr: "127.0.0.1:22347"}
	c2.connect()
	c2.auth()
	c2.wrap.SendOneMessage(&Message{msgType: listNodeMessage})
	buf := new(bytes.Buffer)
	c2.print(buf)
	assert.Contains(t, string(buf.Bytes()), "127.0.0.1")
	assert.Contains(t, string(buf.Bytes()), "ADDRESS")
}

func TestControl_DownloadFile(t *testing.T) {
	s := &Server{Addr: ":12348", Auth: "123456"}
	go s.Listen()
	time.Sleep(100 * time.Millisecond)
	c := &Control{Auth: "123456", Addr: "127.0.0.1:12348"}
	a := &Client{Addr: "127.0.0.1:12348"}
	go a.Connect()
	time.Sleep(100 * time.Millisecond)
	buf := make([]byte, 10241)
	f, _ := os.Create("/tmp/test.bin")
	f.Write(buf)
	c.DownloadFile("1", "/tmp/test.bin", "/tmp/test.bin2")
	md51, _ := MD5sum("/tmp/test.bin")
	md52, _ := MD5sum("/tmp/test.bin2")
	assert.Equal(t, md51, md52)
}