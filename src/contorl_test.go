package gortcp

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestControl_UploadFile(t *testing.T) {

}

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
