package gortcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestControl_ExecCommand(t *testing.T) {
	c := Client{}
	res, err := c.execCommand("ifconfig")
	assert.NoError(t, err)
	assert.Contains(t, string(res), "127.0.0.1")
}

func TestControl_handlerFileInfo(t *testing.T) {

	c := Client{}
	c.handlerFileInfo([]byte("e10adc3949ba59abbe56e057f20f883e/tmp/test"))
	assert.Equal(t, c.fm.dstPath, "/tmp/test")
	assert.Equal(t, c.fm.md5, "e10adc3949ba59abbe56e057f20f883e")

}
