package gortcp

import (
	"testing"

	"os"

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
